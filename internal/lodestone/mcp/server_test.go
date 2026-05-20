package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func runRequest(t *testing.T, s *Server, req string) Response {
	t.Helper()
	var out bytes.Buffer
	in := strings.NewReader(req + "\n")
	if err := s.Serve(context.Background(), in, &out); err != nil {
		t.Fatalf("Serve: %v", err)
	}
	var resp Response
	if err := json.NewDecoder(&out).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v (out=%s)", err, out.String())
	}
	return resp
}

func emptyServer(t *testing.T) *Server {
	t.Helper()
	reg := NewToolRegistry()
	return NewServer("lodestone-mcp-test", "0.0.0", reg)
}

func TestInitialize(t *testing.T) {
	s := emptyServer(t)
	resp := runRequest(t, s, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`)
	if resp.Error != nil {
		t.Fatalf("error: %v", resp.Error)
	}
	raw, _ := json.Marshal(resp.Result)
	if !strings.Contains(string(raw), "lodestone-mcp-test") {
		t.Errorf("result missing serverInfo: %s", raw)
	}
	if !strings.Contains(string(raw), `"tools":{}`) {
		t.Errorf("result missing tools capability: %s", raw)
	}
}

func TestToolsList(t *testing.T) {
	reg := NewToolRegistry()
	reg.Register(Tool{
		Name:        "noop",
		Description: "no-op",
		InputSchema: map[string]any{"type": "object"},
	}, func(_ context.Context, _ json.RawMessage) (*CallToolResult, error) {
		return TextResult("ok"), nil
	})
	s := NewServer("test", "0", reg)
	resp := runRequest(t, s, `{"jsonrpc":"2.0","id":2,"method":"tools/list"}`)
	if resp.Error != nil {
		t.Fatalf("error: %v", resp.Error)
	}
	raw, _ := json.Marshal(resp.Result)
	if !strings.Contains(string(raw), "noop") {
		t.Errorf("tools/list missing registered tool: %s", raw)
	}
}

func TestToolsCallSuccess(t *testing.T) {
	reg := NewToolRegistry()
	reg.Register(Tool{Name: "echo", Description: "echo", InputSchema: map[string]any{"type": "object"}},
		func(_ context.Context, args json.RawMessage) (*CallToolResult, error) {
			return TextResult("got: " + string(args)), nil
		})
	s := NewServer("test", "0", reg)
	resp := runRequest(t, s, `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"echo","arguments":{"x":1}}}`)
	if resp.Error != nil {
		t.Fatalf("error: %v", resp.Error)
	}
	raw, _ := json.Marshal(resp.Result)
	var ctr CallToolResult
	if err := json.Unmarshal(raw, &ctr); err != nil {
		t.Fatalf("re-decode CallToolResult: %v", err)
	}
	if len(ctr.Content) != 1 || ctr.Content[0].Text != `got: {"x":1}` {
		t.Errorf("unexpected content: %+v", ctr.Content)
	}
}

func TestToolsCallUnknownTool(t *testing.T) {
	s := emptyServer(t)
	resp := runRequest(t, s, `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"missing","arguments":{}}}`)
	if resp.Error != nil {
		t.Fatalf("expected isError result, got rpc error: %v", resp.Error)
	}
	raw, _ := json.Marshal(resp.Result)
	if !strings.Contains(string(raw), `"isError":true`) {
		t.Errorf("expected isError flag: %s", raw)
	}
}

func TestUnknownMethod(t *testing.T) {
	s := emptyServer(t)
	resp := runRequest(t, s, `{"jsonrpc":"2.0","id":5,"method":"foo"}`)
	if resp.Error == nil || resp.Error.Code != ErrMethodNotFound {
		t.Errorf("expected method-not-found: %+v", resp.Error)
	}
}

func TestParseError(t *testing.T) {
	s := emptyServer(t)
	resp := runRequest(t, s, "{not-json")
	if resp.Error == nil || resp.Error.Code != ErrParse {
		t.Errorf("expected parse error: %+v", resp.Error)
	}
}

func TestNotificationProducesNoResponse(t *testing.T) {
	s := emptyServer(t)
	var out bytes.Buffer
	in := strings.NewReader(`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n")
	if err := s.Serve(context.Background(), in, &out); err != nil {
		t.Fatalf("Serve: %v", err)
	}
	if out.Len() != 0 {
		t.Errorf("expected no output for notification, got: %s", out.String())
	}
}
