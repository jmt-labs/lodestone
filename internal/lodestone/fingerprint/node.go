package fingerprint

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var nodeFrameworks = map[string]string{
	"react":                     "react",
	"vue":                       "vue",
	"next":                      "next",
	"@anthropic-ai/sdk":         "anthropic-sdk",
	"svelte":                    "svelte",
	"@modelcontextprotocol/sdk": "mcp-sdk",
}

type packageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func (a *Analyzer) parsePackageJSON() (map[string]string, []string, error) {
	path := filepath.Join(a.root, "package.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	var pkg packageJSON
	if err := json.Unmarshal(raw, &pkg); err != nil {
		return nil, nil, err
	}

	deps := map[string]string{}
	for k, v := range pkg.Dependencies {
		deps[k] = v
	}
	for k, v := range pkg.DevDependencies {
		if _, exists := deps[k]; exists {
			continue
		}
		deps[k] = v
	}

	var frameworks []string
	for dep := range deps {
		if fw, ok := nodeFrameworks[dep]; ok {
			frameworks = append(frameworks, fw)
		}
	}
	return deps, frameworks, nil
}
