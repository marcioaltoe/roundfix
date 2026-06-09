#!/bin/sh
# Export a styled HTML document to PDF via headless Chrome/Chromium.
# Mutating helper: writes the output PDF file; changes nothing else.
#
# Usage: sh export_pdf.sh <input.html> <output.pdf>
#
# Exits 0 and prints "OK: <output.pdf>" on success.
# Exits 1 with an ERROR line on stderr when no browser is found or render fails.

set -eu

if [ "$#" -ne 2 ]; then
  echo "ERROR: usage: sh export_pdf.sh <input.html> <output.pdf>" >&2
  exit 1
fi

INPUT=$1
OUTPUT=$2

if [ ! -f "$INPUT" ]; then
  echo "ERROR: input file not found: $INPUT" >&2
  exit 1
fi

find_browser() {
  for candidate in \
    "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" \
    "/Applications/Chromium.app/Contents/MacOS/Chromium" \
    "/Applications/Brave Browser.app/Contents/MacOS/Brave Browser" \
    "/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge" \
    google-chrome google-chrome-stable chromium chromium-browser brave-browser microsoft-edge; do
    case "$candidate" in
      /*) [ -x "$candidate" ] && { printf '%s' "$candidate"; return 0; } ;;
      *) command -v "$candidate" >/dev/null 2>&1 && { command -v "$candidate"; return 0; } ;;
    esac
  done
  return 1
}

BROWSER=$(find_browser) || {
  echo "ERROR: no Chrome/Chromium/Brave/Edge found." >&2
  echo "Deliver the HTML instead — it prints to PDF from any browser (File > Print > Save as PDF)." >&2
  exit 1
}

# Resolve absolute paths (Chrome needs a file:// URL).
INPUT_ABS=$(cd "$(dirname "$INPUT")" && pwd)/$(basename "$INPUT")
OUTPUT_ABS=$(cd "$(dirname "$OUTPUT")" 2>/dev/null && pwd || pwd)/$(basename "$OUTPUT")

"$BROWSER" --headless --disable-gpu --no-pdf-header-footer \
  --print-to-pdf="$OUTPUT_ABS" "file://$INPUT_ABS" >/dev/null 2>&1

if [ -s "$OUTPUT_ABS" ]; then
  # Page count from the PDF page tree (/Count N on the root Pages object).
  PAGES=$(grep -ao "/Count [0-9]*" "$OUTPUT_ABS" | head -1 | grep -o "[0-9]*" || echo "?")
  echo "OK: $OUTPUT_ABS (pages: ${PAGES:-?})"
else
  echo "ERROR: PDF render failed (empty or missing output)" >&2
  exit 1
fi
