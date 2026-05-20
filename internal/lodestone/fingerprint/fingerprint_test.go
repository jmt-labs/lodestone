package fingerprint

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"
)

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func TestAnalyzeGoMinimal(t *testing.T) {
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	a := New(filepath.Join("testdata", "go_minimal"), WithNow(func() time.Time { return now }))

	fp, err := a.Analyze()
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	if !contains(fp.Languages, "Go") {
		t.Errorf("Languages missing Go: %v", fp.Languages)
	}
	if !contains(fp.Frameworks, "cobra") {
		t.Errorf("Frameworks missing cobra: %v", fp.Frameworks)
	}
	if got, ok := fp.Deps["github.com/spf13/cobra"]; !ok || got != "v1.10.2" {
		t.Errorf("Deps[github.com/spf13/cobra] = %q (ok=%v), want v1.10.2", got, ok)
	}
	if _, ok := fp.Deps["github.com/spf13/pflag"]; !ok {
		t.Errorf("expected pflag in deps (block require), got %v", fp.Deps)
	}
	if fp.LOCPerLanguage["Go"] == 0 {
		t.Errorf("expected non-zero Go LOC")
	}
	if fp.TestRatio <= 0 {
		t.Errorf("expected positive TestRatio, got %v", fp.TestRatio)
	}
	if !fp.GeneratedAt.Equal(now) {
		t.Errorf("GeneratedAt = %v, want %v", fp.GeneratedAt, now)
	}
	if fp.SchemaVersion == 0 {
		t.Errorf("SchemaVersion should be set")
	}
}

func TestAnalyzeNodeReact(t *testing.T) {
	a := New(filepath.Join("testdata", "node_react"))

	fp, err := a.Analyze()
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	if !contains(fp.Languages, "JavaScript") {
		t.Errorf("Languages missing JavaScript: %v", fp.Languages)
	}
	if !contains(fp.Frameworks, "react") {
		t.Errorf("Frameworks missing react: %v", fp.Frameworks)
	}
	if !contains(fp.Frameworks, "next") {
		t.Errorf("Frameworks missing next: %v", fp.Frameworks)
	}
	if !contains(fp.Frameworks, "anthropic-sdk") {
		t.Errorf("Frameworks missing anthropic-sdk: %v", fp.Frameworks)
	}
	if _, ok := fp.Deps["react"]; !ok {
		t.Errorf("Deps missing react: %v", fp.Deps)
	}
	if _, ok := fp.Deps["vitest"]; !ok {
		t.Errorf("Deps missing vitest (devDependency): %v", fp.Deps)
	}
}

func TestAnalyzeSkipsExcludedDirs(t *testing.T) {
	root := t.TempDir()

	writeFile(t, root, "main.go", "package m\n\nfunc Foo() {}\n")
	writeFile(t, root, "vendor/big/lib.go", "package big\n\nfunc Big() {}\nfunc Bigger() {}\nfunc Biggest() {}\n")
	writeFile(t, root, "node_modules/react/index.js", "module.exports = {};\nmodule.exports.x = 1;\nmodule.exports.y = 2;\n")
	writeFile(t, root, ".git/HEAD", "ref: refs/heads/main\n")

	fp, err := New(root).Analyze()
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	if fp.LOCPerLanguage["Go"] != 3 {
		t.Errorf("Go LOC = %d, want 3 (vendor must be skipped)", fp.LOCPerLanguage["Go"])
	}
	if got, ok := fp.LOCPerLanguage["JavaScript"]; ok && got > 0 {
		t.Errorf("JavaScript LOC = %d, want 0 (node_modules must be skipped)", got)
	}
}

func TestAnalyzeTestRatio(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "foo.go", "package x\n\nfunc A() {}\nfunc B() {}\nfunc C() {}\n")            // 5 lines, non-test
	writeFile(t, root, "foo_test.go", "package x\n\nimport \"testing\"\nfunc TestA(t *testing.T) {}\n") // 4 lines, test

	fp, err := New(root).Analyze()
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	want := 4.0 / 5.0
	if got := fp.TestRatio; got < want-0.0001 || got > want+0.0001 {
		t.Errorf("TestRatio = %v, want %v", got, want)
	}
	if fp.LOCPerLanguage["Go"] != 9 {
		t.Errorf("Go LOC = %d, want 9", fp.LOCPerLanguage["Go"])
	}
}

func TestAnalyzeCIDetection(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".github/workflows/ci.yml", "name: ci\non: [push]\n")
	writeFile(t, root, "main.go", "package m\n")

	fp, err := New(root).Analyze()
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if !fp.HasCI {
		t.Errorf("HasCI = false, want true")
	}
	if fp.CIProvider != "github_actions" {
		t.Errorf("CIProvider = %q, want github_actions", fp.CIProvider)
	}
}

func TestAnalyzeMCPDetection(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".mcp.json", `{"servers":{}}`)

	fp, err := New(root).Analyze()
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if len(fp.MCPServers) == 0 {
		t.Errorf("expected non-empty MCPServers")
	}
}

func TestAnalyzeEmptyRepo(t *testing.T) {
	fp, err := New(t.TempDir()).Analyze()
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if len(fp.Languages) != 0 {
		t.Errorf("expected no languages in empty repo, got %v", fp.Languages)
	}
}

func TestAnalyzeDeterministicOrdering(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "z.go", "package x\nfunc Z() {}\n")
	writeFile(t, root, "a.go", "package x\nfunc A() {}\n")
	writeFile(t, root, "src/m.jsx", "import React from 'react'\n")
	writeFile(t, root, "package.json", `{"dependencies":{"react":"^18","vue":"^3"}}`)

	var languages, frameworks []string
	for i := 0; i < 3; i++ {
		fp, err := New(root).Analyze()
		if err != nil {
			t.Fatalf("Analyze[%d]: %v", i, err)
		}
		langs := append([]string(nil), fp.Languages...)
		sort.Strings(langs)
		fws := append([]string(nil), fp.Frameworks...)
		sort.Strings(fws)
		if i == 0 {
			languages = langs
			frameworks = fws
			continue
		}
		if !equalSlices(languages, langs) {
			t.Errorf("non-deterministic Languages between run 0 and %d", i)
		}
		if !equalSlices(frameworks, fws) {
			t.Errorf("non-deterministic Frameworks between run 0 and %d", i)
		}
	}
}

func writeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	full := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
