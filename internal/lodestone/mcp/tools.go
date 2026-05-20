package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/audit"
	"github.com/jmt-labs/lodestone/internal/lodestone/fingerprint"
	"github.com/jmt-labs/lodestone/internal/lodestone/planning"
	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
	"github.com/jmt-labs/lodestone/internal/lodestone/scoring"
	"github.com/jmt-labs/lodestone/internal/lodestone/store"
)

type ToolHandler func(ctx context.Context, args json.RawMessage) (*CallToolResult, error)

type ToolRegistry struct {
	tools map[string]registeredTool
	order []string
}

type registeredTool struct {
	def     Tool
	handler ToolHandler
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{tools: map[string]registeredTool{}}
}

func (r *ToolRegistry) Register(t Tool, h ToolHandler) {
	r.tools[t.Name] = registeredTool{def: t, handler: h}
	r.order = append(r.order, t.Name)
}

func (r *ToolRegistry) List() []Tool {
	out := make([]Tool, 0, len(r.order))
	for _, name := range r.order {
		out = append(out, r.tools[name].def)
	}
	return out
}

func (r *ToolRegistry) Call(ctx context.Context, name string, args json.RawMessage) (*CallToolResult, error) {
	tool, ok := r.tools[name]
	if !ok {
		return ErrorResult(fmt.Sprintf("unknown tool: %s", name)), nil
	}
	return tool.handler(ctx, args)
}

type ToolDeps struct {
	StoreRoot string
	RepoRoot  string
	Planning  *planning.Engine
	Now       func() time.Time
}

func DefaultDeps(repoRoot, storeRoot string) ToolDeps {
	return ToolDeps{
		StoreRoot: storeRoot,
		RepoRoot:  repoRoot,
		Planning:  planning.New(),
		Now:       func() time.Time { return time.Now().UTC() },
	}
}

func RegisterBuiltins(reg *ToolRegistry, deps ToolDeps) {
	reg.Register(Tool{
		Name:        "list_signals",
		Description: "List signals from the lodestone store, optionally filtered.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"source": map[string]any{"type": "string", "description": "optional source filter (e.g. github_trending, hackernews)"},
				"since":  map[string]any{"type": "string", "description": "RFC3339 cutoff timestamp"},
				"top":    map[string]any{"type": "integer", "description": "limit to top-N by stars"},
			},
		},
	}, listSignalsHandler(deps))

	reg.Register(Tool{
		Name:        "query_trends",
		Description: "Aggregate statistics over the stored signals.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"since": map[string]any{"type": "string", "description": "RFC3339 cutoff timestamp"},
			},
		},
	}, queryTrendsHandler(deps))

	reg.Register(Tool{
		Name:        "score_repo",
		Description: "Recompute fingerprint + score and return the top recommendations.",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
	}, scoreRepoHandler(deps))

	reg.Register(Tool{
		Name:        "generate_plan",
		Description: "Invoke the planning engine for a specific recommendation ID.",
		InputSchema: map[string]any{
			"type":     "object",
			"required": []string{"rec_id"},
			"properties": map[string]any{
				"rec_id": map[string]any{"type": "string"},
				"model":  map[string]any{"type": "string"},
			},
		},
	}, generatePlanHandler(deps))

	reg.Register(Tool{
		Name:        "record_decision",
		Description: "Append an entry to .lodestone/decisions.log.",
		InputSchema: map[string]any{
			"type":     "object",
			"required": []string{"verb", "outcome"},
			"properties": map[string]any{
				"verb":    map[string]any{"type": "string"},
				"outcome": map[string]any{"type": "string"},
				"detail":  map[string]any{"type": "string"},
				"args":    map[string]any{"type": "object"},
			},
		},
	}, recordDecisionHandler(deps))
}

func listSignalsHandler(deps ToolDeps) ToolHandler {
	type listArgs struct {
		Source string `json:"source"`
		Since  string `json:"since"`
		Top    int    `json:"top"`
	}
	return func(_ context.Context, raw json.RawMessage) (*CallToolResult, error) {
		var a listArgs
		if len(raw) > 0 {
			if err := json.Unmarshal(raw, &a); err != nil {
				return ErrorResult("invalid arguments: " + err.Error()), nil
			}
		}
		s, err := store.New(deps.StoreRoot)
		if err != nil {
			return ErrorResult("store: " + err.Error()), nil
		}
		var since time.Time
		if a.Since != "" {
			t, err := time.Parse(time.RFC3339, a.Since)
			if err != nil {
				return ErrorResult("since must be RFC3339: " + err.Error()), nil
			}
			since = t
		}
		sigs, err := s.ListSince(since)
		if err != nil {
			return ErrorResult("list: " + err.Error()), nil
		}
		if a.Source != "" {
			filtered := sigs[:0]
			for _, sig := range sigs {
				if sig.Source == a.Source {
					filtered = append(filtered, sig)
				}
			}
			sigs = filtered
		}
		sort.Slice(sigs, func(i, j int) bool { return sigs[i].Stars > sigs[j].Stars })
		if a.Top > 0 && len(sigs) > a.Top {
			sigs = sigs[:a.Top]
		}
		raw2, _ := json.MarshalIndent(sigs, "", "  ")
		return TextResult(string(raw2)), nil
	}
}

func queryTrendsHandler(deps ToolDeps) ToolHandler {
	type trendsArgs struct {
		Since string `json:"since"`
	}
	type trendsOut struct {
		CountBySource map[string]int     `json:"count_by_source"`
		AvgStars      map[string]float64 `json:"avg_stars"`
		Total         int                `json:"total"`
	}
	return func(_ context.Context, raw json.RawMessage) (*CallToolResult, error) {
		var a trendsArgs
		if len(raw) > 0 {
			_ = json.Unmarshal(raw, &a)
		}
		s, err := store.New(deps.StoreRoot)
		if err != nil {
			return ErrorResult("store: " + err.Error()), nil
		}
		var since time.Time
		if a.Since != "" {
			t, err := time.Parse(time.RFC3339, a.Since)
			if err == nil {
				since = t
			}
		}
		sigs, err := s.ListSince(since)
		if err != nil {
			return ErrorResult("list: " + err.Error()), nil
		}
		counts := map[string]int{}
		stars := map[string]int{}
		for _, sig := range sigs {
			counts[sig.Source]++
			stars[sig.Source] += sig.Stars
		}
		avg := map[string]float64{}
		for src, n := range counts {
			if n > 0 {
				avg[src] = float64(stars[src]) / float64(n)
			}
		}
		out := trendsOut{CountBySource: counts, AvgStars: avg, Total: len(sigs)}
		raw2, _ := json.MarshalIndent(out, "", "  ")
		return TextResult(string(raw2)), nil
	}
}

func scoreRepoHandler(deps ToolDeps) ToolHandler {
	type out struct {
		FingerprintSummary string `json:"fingerprint_summary"`
		TopRecommendations []any  `json:"top_recommendations"`
	}
	return func(_ context.Context, _ json.RawMessage) (*CallToolResult, error) {
		fp, err := fingerprint.New(deps.RepoRoot).Analyze()
		if err != nil {
			return ErrorResult("fingerprint: " + err.Error()), nil
		}
		s, err := store.New(deps.StoreRoot)
		if err != nil {
			return ErrorResult("store: " + err.Error()), nil
		}
		if err := s.Write(fp); err != nil {
			return ErrorResult("store write: " + err.Error()), nil
		}
		sigs, err := s.ListSince(time.Time{})
		if err != nil {
			return ErrorResult("signals: " + err.Error()), nil
		}
		recs, err := scoring.Score(fp, sigs)
		if err != nil {
			return ErrorResult("score: " + err.Error()), nil
		}
		if err := s.Replace(recs); err != nil {
			return ErrorResult("store replace: " + err.Error()), nil
		}
		top := recs
		if len(top) > 5 {
			top = top[:5]
		}
		topAny := make([]any, len(top))
		for i, r := range top {
			topAny[i] = r
		}
		summary := fmt.Sprintf("languages=%v frameworks=%v deps=%d", fp.Languages, fp.Frameworks, len(fp.Deps))
		raw, _ := json.MarshalIndent(out{FingerprintSummary: summary, TopRecommendations: topAny}, "", "  ")
		return TextResult(string(raw)), nil
	}
}

func generatePlanHandler(deps ToolDeps) ToolHandler {
	type planArgs struct {
		RecID string `json:"rec_id"`
		Model string `json:"model"`
	}
	type planOut struct {
		Spec     string `json:"spec_md"`
		Plan     string `json:"plan_md"`
		SpecPath string `json:"spec_path"`
		PlanPath string `json:"plan_path"`
		Model    string `json:"model"`
	}
	return func(ctx context.Context, raw json.RawMessage) (*CallToolResult, error) {
		var a planArgs
		if err := json.Unmarshal(raw, &a); err != nil {
			return ErrorResult("invalid arguments: " + err.Error()), nil
		}
		if a.RecID == "" {
			return ErrorResult("rec_id required"), nil
		}
		s, err := store.New(deps.StoreRoot)
		if err != nil {
			return ErrorResult("store: " + err.Error()), nil
		}
		fp, err := s.Read()
		if err != nil {
			return ErrorResult("read fingerprint: " + err.Error()), nil
		}
		recs, err := s.List()
		if err != nil {
			return ErrorResult("read recommendations: " + err.Error()), nil
		}
		var match *schema.Recommendation
		for i := range recs {
			if recs[i].ID == a.RecID || recs[i].SignalID == a.RecID {
				match = &recs[i]
				break
			}
		}
		if match == nil {
			return ErrorResult("recommendation not found: " + a.RecID), nil
		}
		engine := deps.Planning
		if a.Model != "" {
			engine = planning.New(planning.WithModel(a.Model))
		}
		res, err := engine.Plan(ctx, fp, *match)
		if err != nil {
			return ErrorResult("plan: " + err.Error()), nil
		}
		if err := res.Persist(deps.RepoRoot); err != nil {
			return ErrorResult("persist: " + err.Error()), nil
		}
		out := planOut{
			Spec:     res.Spec,
			Plan:     res.Plan,
			SpecPath: res.SpecPath,
			PlanPath: res.PlanPath,
			Model:    res.Model,
		}
		raw2, _ := json.MarshalIndent(out, "", "  ")
		return TextResult(string(raw2)), nil
	}
}

func recordDecisionHandler(deps ToolDeps) ToolHandler {
	type recArgs struct {
		Verb    string            `json:"verb"`
		Outcome string            `json:"outcome"`
		Detail  string            `json:"detail"`
		Args    map[string]string `json:"args"`
	}
	return func(_ context.Context, raw json.RawMessage) (*CallToolResult, error) {
		var a recArgs
		if err := json.Unmarshal(raw, &a); err != nil {
			return ErrorResult("invalid arguments: " + err.Error()), nil
		}
		if a.Verb == "" || a.Outcome == "" {
			return ErrorResult("verb and outcome required"), nil
		}
		log, err := audit.New(deps.StoreRoot)
		if err != nil {
			return ErrorResult("audit: " + err.Error()), nil
		}
		if err := log.Record(audit.Entry{
			Verb:    a.Verb,
			Outcome: a.Outcome,
			Detail:  a.Detail,
			Args:    a.Args,
		}); err != nil {
			return ErrorResult("record: " + err.Error()), nil
		}
		return TextResult(`{"ok":true}`), nil
	}
}
