package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type Server struct {
	Registry *ToolRegistry
	Name     string
	Version  string
}

func NewServer(name, version string, reg *ToolRegistry) *Server {
	return &Server{Registry: reg, Name: name, Version: version}
}

func (s *Server) Serve(ctx context.Context, in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	enc := json.NewEncoder(out)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		resp := s.handle(ctx, line)
		if resp == nil {
			continue
		}
		if err := enc.Encode(resp); err != nil {
			return fmt.Errorf("encode response: %w", err)
		}
	}
	return scanner.Err()
}

func (s *Server) handle(ctx context.Context, raw []byte) *Response {
	var req Request
	if err := json.Unmarshal(raw, &req); err != nil {
		return &Response{
			JSONRPC: "2.0",
			Error:   &RPCError{Code: ErrParse, Message: "parse error: " + err.Error()},
		}
	}
	if req.JSONRPC != "2.0" {
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &RPCError{Code: ErrInvalidRequest, Message: "jsonrpc must be \"2.0\""},
		}
	}

	switch req.Method {
	case "initialize":
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: InitializeResult{
				ProtocolVersion: ProtocolVersion,
				ServerInfo:      ServerInfo{Name: s.Name, Version: s.Version},
				Capabilities:    ServerCapabilities{Tools: &struct{}{}},
			},
		}
	case "initialized", "notifications/initialized":
		// notification, no response
		return nil
	case "tools/list":
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  ListToolsResult{Tools: s.Registry.List()},
		}
	case "tools/call":
		var p CallToolParams
		if err := json.Unmarshal(req.Params, &p); err != nil {
			return &Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &RPCError{Code: ErrInvalidParams, Message: "invalid params: " + err.Error()},
			}
		}
		res, err := s.Registry.Call(ctx, p.Name, p.Arguments)
		if err != nil {
			return &Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &RPCError{Code: ErrInternal, Message: err.Error()},
			}
		}
		return &Response{JSONRPC: "2.0", ID: req.ID, Result: res}
	default:
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &RPCError{Code: ErrMethodNotFound, Message: "method not found: " + req.Method},
		}
	}
}
