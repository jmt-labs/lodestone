package scoring

import (
	"strings"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

const (
	languageWeight  = 1.5
	frameworkWeight = 1.0
)

func Compatibility(sig schema.Signal, fp schema.Fingerprint) float64 {
	signalSet := buildSignalSet(sig)
	if len(signalSet) == 0 {
		return 0
	}
	langSet := buildLowerSet(fp.Languages)
	fwSet := buildLowerSet(fp.Frameworks)

	union := map[string]struct{}{}
	for k := range signalSet {
		union[k] = struct{}{}
	}
	for k := range langSet {
		union[k] = struct{}{}
	}
	for k := range fwSet {
		union[k] = struct{}{}
	}
	if len(union) == 0 {
		return 0
	}

	var num float64
	for elem := range signalSet {
		if _, ok := langSet[elem]; ok {
			num += languageWeight
			continue
		}
		if _, ok := fwSet[elem]; ok {
			num += frameworkWeight
		}
	}

	score := num / float64(len(union))
	if score > 1.0 {
		score = 1.0
	}
	return score
}

func buildSignalSet(sig schema.Signal) map[string]struct{} {
	set := map[string]struct{}{}
	for _, t := range sig.TopicTags {
		k := strings.ToLower(strings.TrimSpace(t))
		if k != "" {
			set[k] = struct{}{}
		}
	}
	if k := strings.ToLower(strings.TrimSpace(sig.Language)); k != "" {
		set[k] = struct{}{}
	}
	return set
}

func buildLowerSet(items []string) map[string]struct{} {
	set := map[string]struct{}{}
	for _, item := range items {
		k := strings.ToLower(strings.TrimSpace(item))
		if k != "" {
			set[k] = struct{}{}
		}
	}
	return set
}
