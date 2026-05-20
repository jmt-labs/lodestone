#!/usr/bin/env bash
# Verifies that every relative Markdown link within docs/ and root-level
# markdown files points to an existing path. Skips http(s), mailto: and
# pure fragments (#anchor).

set -euo pipefail

missing=0

shopt -s globstar nullglob

files=(
  README.md
  CONTRIBUTING.md
  CLAUDE.md
  AGENTS.md
  CHANGELOG.md
)
for f in docs/**/*.md; do
  files+=("$f")
done

link_re='\[[^]]*\]\(([^)]+)\)'

for file in "${files[@]}"; do
  [ -f "$file" ] || continue
  dir=$(dirname "$file")
  # Extract every Markdown link target.
  while IFS= read -r target; do
    [ -z "$target" ] && continue
    # Skip external + mailto + pure-anchor links.
    case "$target" in
      http://*|https://*|mailto:*|'#'*) continue ;;
    esac
    # Strip fragment portion.
    target_noanchor="${target%%#*}"
    [ -z "$target_noanchor" ] && continue
    resolved="$dir/$target_noanchor"
    if [ ! -e "$resolved" ]; then
      echo "::error::$file: broken link → $target"
      missing=1
    fi
  done < <(grep -oE "$link_re" "$file" 2>/dev/null | sed -E 's/.*\(([^)]+)\)$/\1/')
done

exit "$missing"
