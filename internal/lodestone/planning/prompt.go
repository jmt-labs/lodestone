package planning

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

const promptTemplate = `Du bist ein erfahrener Software-Architekt im Team von ` + "`lodestone`" + `.
Für die folgende Recommendation gegen das beschriebene Repo produzierst du:

1. Eine Spec im superpowers-Format (Markdown).
2. Einen ausführbaren Plan mit Checkbox-Tasks (Markdown).

WICHTIG: Antworte deutsch. Verwende exakt diese drei Block-Marker, ohne
Variationen, jeweils auf eigener Zeile:

===SPEC===
<Spec-Inhalt>
===PLAN===
<Plan-Inhalt>
===END===

Repo-Fingerprint (JSON):
%s

Recommendation (JSON):
%s

Konventionen:
- YAGNI. Keine spekulativen Features.
- Plan-Tasks im Format ` + "`- [ ] T<N>: <Beschreibung>`" + `, atomar und einzeln testbar.
- Spec listet Tradeoffs explizit (mindestens build-vs-buy oder lokal-vs-extern).
- Spec und Plan sollen jeweils unter 250 Zeilen bleiben.`

func BuildPrompt(fp schema.Fingerprint, rec schema.Recommendation) (string, error) {
	fpJSON, err := json.MarshalIndent(fp, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal fingerprint: %w", err)
	}
	recJSON, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal recommendation: %w", err)
	}
	return fmt.Sprintf(promptTemplate, fpJSON, recJSON), nil
}

func SplitResponse(out string) (specMD, planMD string, err error) {
	specMarker := "===SPEC==="
	planMarker := "===PLAN==="
	endMarker := "===END==="

	specIdx := strings.Index(out, specMarker)
	planIdx := strings.Index(out, planMarker)
	endIdx := strings.Index(out, endMarker)

	if specIdx < 0 || planIdx < 0 || planIdx < specIdx {
		return "", "", fmt.Errorf("missing or out-of-order SPEC/PLAN markers in claude output")
	}

	specStart := specIdx + len(specMarker)
	specEnd := planIdx
	planStart := planIdx + len(planMarker)
	planEnd := len(out)
	if endIdx > planIdx {
		planEnd = endIdx
	}

	specMD = strings.TrimSpace(out[specStart:specEnd])
	planMD = strings.TrimSpace(out[planStart:planEnd])
	if specMD == "" || planMD == "" {
		return "", "", fmt.Errorf("empty SPEC or PLAN section")
	}
	return specMD, planMD, nil
}
