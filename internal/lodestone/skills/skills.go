package skills

import (
	"embed"
	"io/fs"
)

//go:embed data/*.md
var raw embed.FS

func FS() fs.FS {
	sub, err := fs.Sub(raw, "data")
	if err != nil {
		panic(err)
	}
	return sub
}

func List() ([]string, error) {
	entries, err := raw.ReadDir("data")
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			out = append(out, e.Name())
		}
	}
	return out, nil
}

func Read(name string) ([]byte, error) {
	return raw.ReadFile("data/" + name)
}
