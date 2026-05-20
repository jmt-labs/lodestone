package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func writeYAML(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".lodestone.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	return path
}

func TestLoadMissingFileReturnsDefaults(t *testing.T) {
	cfg, err := Load(filepath.Join(t.TempDir(), "absent.yaml"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	want := Defaults()
	if !reflect.DeepEqual(cfg, want) {
		t.Errorf("got %+v, want %+v", cfg, want)
	}
}

func TestLoadEmptyFileReturnsDefaults(t *testing.T) {
	cfg, err := Load(writeYAML(t, ""))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	want := Defaults()
	if !reflect.DeepEqual(cfg, want) {
		t.Errorf("got %+v, want %+v", cfg, want)
	}
}

func TestLoadGoalsOnlyKeepsLodestoneDefaults(t *testing.T) {
	cfg, err := Load(writeYAML(t, "goals:\n  - reliability\n  - speed\n"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !reflect.DeepEqual(cfg.Goals, []string{"reliability", "speed"}) {
		t.Errorf("Goals = %v", cfg.Goals)
	}
	if cfg.Lodestone.MinStars != 50 {
		t.Errorf("MinStars = %d, want default 50", cfg.Lodestone.MinStars)
	}
	if !cfg.Lodestone.RequireLicense {
		t.Errorf("RequireLicense should default to true")
	}
}

func TestLoadFullLodestoneBlockOverrides(t *testing.T) {
	yaml := `
goals:
  - shipping
tech_interests:
  - mcp
  - llm-tools
lodestone:
  min_stars: 200
  min_age_days: 7
  max_last_commit_age_days: 30
  require_license: false
`
	cfg, err := Load(writeYAML(t, yaml))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	want := Config{
		Goals:         []string{"shipping"},
		TechInterests: []string{"mcp", "llm-tools"},
		Lodestone: LodestoneConfig{
			MinStars:             200,
			MinAgeDays:           7,
			MaxLastCommitAgeDays: 30,
			RequireLicense:       false,
		},
	}
	if !reflect.DeepEqual(cfg, want) {
		t.Errorf("got %+v\nwant %+v", cfg, want)
	}
}

func TestLoadPartialLodestoneBlockKeepsOtherDefaults(t *testing.T) {
	cfg, err := Load(writeYAML(t, "lodestone:\n  min_stars: 1000\n"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Lodestone.MinStars != 1000 {
		t.Errorf("MinStars = %d, want 1000", cfg.Lodestone.MinStars)
	}
	if cfg.Lodestone.MinAgeDays != 30 {
		t.Errorf("MinAgeDays = %d, want default 30", cfg.Lodestone.MinAgeDays)
	}
	if cfg.Lodestone.MaxLastCommitAgeDays != 180 {
		t.Errorf("MaxLastCommitAgeDays = %d, want default 180", cfg.Lodestone.MaxLastCommitAgeDays)
	}
	if !cfg.Lodestone.RequireLicense {
		t.Errorf("RequireLicense should default to true when unset")
	}
}

func TestLoadInvalidYAMLReturnsError(t *testing.T) {
	_, err := Load(writeYAML(t, "goals: [\nunterminated"))
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestLoadExplicitRequireLicenseFalse(t *testing.T) {
	cfg, err := Load(writeYAML(t, "lodestone:\n  require_license: false\n"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Lodestone.RequireLicense {
		t.Errorf("RequireLicense = true, want false (explicit override)")
	}
}

func TestDefaults(t *testing.T) {
	d := Defaults()
	if d.Lodestone.MinStars != 50 || d.Lodestone.MinAgeDays != 30 ||
		d.Lodestone.MaxLastCommitAgeDays != 180 || !d.Lodestone.RequireLicense {
		t.Errorf("defaults wrong: %+v", d.Lodestone)
	}
	if len(d.Goals) != 0 || len(d.TechInterests) != 0 {
		t.Errorf("expected empty Goals/TechInterests, got %+v / %+v", d.Goals, d.TechInterests)
	}
}
