package fingerprint

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

type Analyzer struct {
	root string
	now  func() time.Time
}

type Option func(*Analyzer)

func WithNow(fn func() time.Time) Option {
	return func(a *Analyzer) { a.now = fn }
}

func New(root string, opts ...Option) *Analyzer {
	a := &Analyzer{
		root: root,
		now:  func() time.Time { return time.Now().UTC() },
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

var skipDirs = map[string]struct{}{
	".git":         {},
	"vendor":       {},
	"node_modules": {},
	"dist":         {},
	"build":        {},
}

var extLanguage = map[string]string{
	".go":  "Go",
	".js":  "JavaScript",
	".jsx": "JavaScript",
	".ts":  "JavaScript",
	".tsx": "JavaScript",
	".mjs": "JavaScript",
	".cjs": "JavaScript",
}

type langStats struct {
	totalLOC int
	testLOC  int
}

func (a *Analyzer) Analyze() (schema.Fingerprint, error) {
	fp := schema.Fingerprint{
		SchemaVersion:  schema.FingerprintSchemaVersion,
		GeneratedAt:    a.now(),
		Deps:           map[string]string{},
		LOCPerLanguage: map[string]int{},
	}

	stats, err := a.walk()
	if err != nil {
		return schema.Fingerprint{}, err
	}

	var totalTest, totalNon int
	langs := make([]string, 0, len(stats))
	for lang, s := range stats {
		langs = append(langs, lang)
		fp.LOCPerLanguage[lang] = s.totalLOC
		totalTest += s.testLOC
		totalNon += s.totalLOC - s.testLOC
	}
	sort.Strings(langs)
	fp.Languages = langs
	if totalNon > 0 {
		fp.TestRatio = float64(totalTest) / float64(totalNon)
	}

	frameworks := map[string]struct{}{}

	if goDeps, goFrameworks, err := a.parseGoMod(); err != nil {
		return schema.Fingerprint{}, err
	} else {
		for k, v := range goDeps {
			fp.Deps[k] = v
		}
		for _, fw := range goFrameworks {
			frameworks[fw] = struct{}{}
		}
	}

	if nodeDeps, nodeFrameworks, err := a.parsePackageJSON(); err != nil {
		return schema.Fingerprint{}, err
	} else {
		for k, v := range nodeDeps {
			fp.Deps[k] = v
		}
		for _, fw := range nodeFrameworks {
			frameworks[fw] = struct{}{}
		}
	}

	fwList := make([]string, 0, len(frameworks))
	for fw := range frameworks {
		fwList = append(fwList, fw)
	}
	sort.Strings(fwList)
	fp.Frameworks = fwList

	fp.HasCI, fp.CIProvider = a.detectCI()
	fp.MCPServers = a.detectMCPServers()

	return fp, nil
}

func (a *Analyzer) walk() (map[string]*langStats, error) {
	stats := map[string]*langStats{}

	err := filepath.WalkDir(a.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path == a.root {
				return nil
			}
			if _, skip := skipDirs[d.Name()]; skip {
				return filepath.SkipDir
			}
			return nil
		}
		ext := strings.ToLower(filepath.Ext(d.Name()))
		lang, known := extLanguage[ext]
		if !known {
			return nil
		}
		rel, _ := filepath.Rel(a.root, path)
		rel = filepath.ToSlash(rel)
		loc, err := countLines(path)
		if err != nil {
			return fmt.Errorf("count %s: %w", rel, err)
		}
		s, ok := stats[lang]
		if !ok {
			s = &langStats{}
			stats[lang] = s
		}
		s.totalLOC += loc
		if isTestFile(rel, d.Name(), lang) {
			s.testLOC += loc
		}
		return nil
	})
	return stats, err
}

func countLines(path string) (int, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	if len(raw) == 0 {
		return 0, nil
	}
	n := 0
	for _, b := range raw {
		if b == '\n' {
			n++
		}
	}
	if raw[len(raw)-1] != '\n' {
		n++
	}
	return n, nil
}

func isTestFile(rel, base, lang string) bool {
	switch lang {
	case "Go":
		return strings.HasSuffix(base, "_test.go")
	case "JavaScript":
		if strings.Contains(rel, "__tests__/") {
			return true
		}
		for _, suffix := range []string{
			".test.js", ".test.jsx", ".test.ts", ".test.tsx",
			".spec.js", ".spec.jsx", ".spec.ts", ".spec.tsx",
		} {
			if strings.HasSuffix(base, suffix) {
				return true
			}
		}
	}
	return false
}

func (a *Analyzer) detectCI() (bool, string) {
	ghDir := filepath.Join(a.root, ".github", "workflows")
	if entries, err := os.ReadDir(ghDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
				return true, "github_actions"
			}
		}
	}
	if _, err := os.Stat(filepath.Join(a.root, ".gitlab-ci.yml")); err == nil {
		return true, "gitlab_ci"
	}
	if _, err := os.Stat(filepath.Join(a.root, ".circleci", "config.yml")); err == nil {
		return true, "circleci"
	}
	return false, ""
}

func (a *Analyzer) detectMCPServers() []string {
	if _, err := os.Stat(filepath.Join(a.root, ".mcp.json")); err == nil {
		return []string{"mcp_config_present"}
	}
	return nil
}
