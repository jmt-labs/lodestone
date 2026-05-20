package fingerprint

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var goFrameworks = map[string]string{
	"github.com/spf13/cobra":                 "cobra",
	"github.com/anthropics/anthropic-sdk-go": "anthropic-sdk",
	"github.com/gin-gonic/gin":               "gin",
	"github.com/labstack/echo/v4":            "echo",
	"github.com/go-chi/chi/v5":               "chi",
}

var goRequireLineRE = regexp.MustCompile(`^([^\s]+)\s+(v[^\s]+)`)

func (a *Analyzer) parseGoMod() (map[string]string, []string, error) {
	path := filepath.Join(a.root, "go.mod")
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	deps := map[string]string{}
	inBlock := false

	for _, line := range strings.Split(string(raw), "\n") {
		if idx := strings.Index(line, "//"); idx >= 0 {
			line = line[:idx]
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if trimmed == "require (" {
			inBlock = true
			continue
		}
		if inBlock && trimmed == ")" {
			inBlock = false
			continue
		}

		candidate := trimmed
		if !inBlock {
			if !strings.HasPrefix(trimmed, "require ") {
				continue
			}
			candidate = strings.TrimSpace(strings.TrimPrefix(trimmed, "require "))
		}

		if m := goRequireLineRE.FindStringSubmatch(candidate); m != nil {
			deps[m[1]] = m[2]
		}
	}

	var frameworks []string
	for dep := range deps {
		if fw, ok := goFrameworks[dep]; ok {
			frameworks = append(frameworks, fw)
		}
	}
	return deps, frameworks, nil
}
