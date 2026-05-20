package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"
)

const DefaultFilename = ".lodestone.yaml"

type Config struct {
	Goals         []string        `yaml:"goals,omitempty"`
	TechInterests []string        `yaml:"tech_interests,omitempty"`
	Lodestone     LodestoneConfig `yaml:"lodestone,omitempty"`
}

type LodestoneConfig struct {
	MinStars             int  `yaml:"min_stars"`
	MinAgeDays           int  `yaml:"min_age_days"`
	MaxLastCommitAgeDays int  `yaml:"max_last_commit_age_days"`
	RequireLicense       bool `yaml:"require_license"`
}

func Defaults() Config {
	return Config{
		Lodestone: LodestoneConfig{
			MinStars:             50,
			MinAgeDays:           30,
			MaxLastCommitAgeDays: 180,
			RequireLicense:       true,
		},
	}
}

type rawConfig struct {
	Goals         []string            `yaml:"goals,omitempty"`
	TechInterests []string            `yaml:"tech_interests,omitempty"`
	Lodestone     *rawLodestoneConfig `yaml:"lodestone,omitempty"`
}

type rawLodestoneConfig struct {
	MinStars             *int  `yaml:"min_stars,omitempty"`
	MinAgeDays           *int  `yaml:"min_age_days,omitempty"`
	MaxLastCommitAgeDays *int  `yaml:"max_last_commit_age_days,omitempty"`
	RequireLicense       *bool `yaml:"require_license,omitempty"`
}

func Load(path string) (Config, error) {
	cfg := Defaults()
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("read %s: %w", path, err)
	}
	if len(raw) == 0 {
		return cfg, nil
	}

	var overlay rawConfig
	if err := yaml.Unmarshal(raw, &overlay); err != nil {
		return cfg, fmt.Errorf("parse %s: %w", path, err)
	}

	if overlay.Goals != nil {
		cfg.Goals = overlay.Goals
	}
	if overlay.TechInterests != nil {
		cfg.TechInterests = overlay.TechInterests
	}
	if overlay.Lodestone != nil {
		if overlay.Lodestone.MinStars != nil {
			cfg.Lodestone.MinStars = *overlay.Lodestone.MinStars
		}
		if overlay.Lodestone.MinAgeDays != nil {
			cfg.Lodestone.MinAgeDays = *overlay.Lodestone.MinAgeDays
		}
		if overlay.Lodestone.MaxLastCommitAgeDays != nil {
			cfg.Lodestone.MaxLastCommitAgeDays = *overlay.Lodestone.MaxLastCommitAgeDays
		}
		if overlay.Lodestone.RequireLicense != nil {
			cfg.Lodestone.RequireLicense = *overlay.Lodestone.RequireLicense
		}
	}
	return cfg, nil
}
