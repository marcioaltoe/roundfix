#!/usr/bin/env python3
"""Generate self-contained SVG charts for business documents.

Bootstrap helper: creates an SVG file (or prints to stdout); modifies nothing else.
Stdlib only — no dependencies, deterministic output, prints perfectly in PDF.

Usage:
  python3 make_chart.py --type bar   --title "Hours saved per week" \
      --data "Support:14,Finance:9,Operations:6" --out chart.svg
  python3 make_chart.py --type line  --title "Tickets per month" \
      --data "Jan:120,Feb:95,Mar:60,Apr:42" --out chart.svg
  python3 make_chart.py --type donut --title "Where time goes today" \
      --data "Manual entry:55,Rework:25,Analysis:20" --out chart.svg

Data format: comma-separated "Label:value" pairs. Values: numbers from the
source document only — never invented.
"""

import argparse
import math
import sys

ACCENT = "#1a5fb4"
PALETTE = ["#1a5fb4", "#26a269", "#e5a50a", "#c01c28", "#613583", "#63452c"]
INK = "#1f2430"
MUTED = "#5c6470"
RULE = "#e3e6ea"
FONT = 'font-family="Avenir Next, Segoe UI, Inter, sans-serif"'


def parse_data(raw):
    pairs = []
    for chunk in raw.split(","):
        chunk = chunk.strip()
        if not chunk:
            continue
        label, _, value = chunk.rpartition(":")
        if not label:
            sys.exit(f"ERROR: bad data chunk '{chunk}' — expected 'Label:value'")
        try:
            pairs.append((label.strip(), float(value)))
        except ValueError:
            sys.exit(f"ERROR: value '{value}' in '{chunk}' is not a number")
    if not pairs:
        sys.exit("ERROR: no data points parsed")
    return pairs


def fmt(value):
    return f"{value:g}"


def svg_open(width, height, title):
    parts = [
        f'<svg xmlns="http://www.w3.org/2000/svg" width="{width}" height="{height}" '
        f'viewBox="0 0 {width} {height}" role="img" aria-label="{title}">',
        f'<text x="{width / 2}" y="26" text-anchor="middle" {FONT} '
        f'font-size="16" font-weight="600" fill="{INK}">{title}</text>',
    ]
    return parts


def bar_chart(pairs, title, width, height, accent):
    top, left, right, bottom = 48, 150, 70, 20
    plot_w = width - left - right
    row_h = (height - top - bottom) / len(pairs)
    bar_h = min(28, row_h * 0.62)
    max_v = max(v for _, v in pairs) or 1

    out = svg_open(width, height, title)
    for i, (label, value) in enumerate(pairs):
        y = top + i * row_h + (row_h - bar_h) / 2
        w = (value / max_v) * plot_w
        out.append(
            f'<text x="{left - 10}" y="{y + bar_h / 2 + 4}" text-anchor="end" '
            f'{FONT} font-size="13" fill="{INK}">{label}</text>'
        )
        out.append(
            f'<rect x="{left}" y="{y:.1f}" width="{w:.1f}" height="{bar_h:.1f}" '
            f'rx="4" fill="{accent}"/>'
        )
        out.append(
            f'<text x="{left + w + 8:.1f}" y="{y + bar_h / 2 + 4:.1f}" {FONT} '
            f'font-size="13" font-weight="600" fill="{MUTED}">{fmt(value)}</text>'
        )
    out.append("</svg>")
    return "\n".join(out)


def line_chart(pairs, title, width, height, accent):
    top, left, right, bottom = 48, 60, 30, 46
    plot_w, plot_h = width - left - right, height - top - bottom
    max_v = max(v for _, v in pairs) or 1
    step = plot_w / max(len(pairs) - 1, 1)

    def xy(i, v):
        return left + i * step, top + plot_h - (v / max_v) * plot_h

    out = svg_open(width, height, title)
    for frac in (0, 0.5, 1):
        gy = top + plot_h - frac * plot_h
        out.append(
            f'<line x1="{left}" y1="{gy:.1f}" x2="{left + plot_w}" y2="{gy:.1f}" '
            f'stroke="{RULE}" stroke-width="1"/>'
        )
        out.append(
            f'<text x="{left - 8}" y="{gy + 4:.1f}" text-anchor="end" {FONT} '
            f'font-size="11" fill="{MUTED}">{fmt(max_v * frac)}</text>'
        )
    points = " ".join(f"{x:.1f},{y:.1f}" for x, y in (xy(i, v) for i, (_, v) in enumerate(pairs)))
    out.append(
        f'<polyline points="{points}" fill="none" stroke="{accent}" '
        f'stroke-width="2.5" stroke-linejoin="round" stroke-linecap="round"/>'
    )
    for i, (label, value) in enumerate(pairs):
        x, y = xy(i, value)
        out.append(f'<circle cx="{x:.1f}" cy="{y:.1f}" r="4" fill="{accent}"/>')
        out.append(
            f'<text x="{x:.1f}" y="{y - 10:.1f}" text-anchor="middle" {FONT} '
            f'font-size="12" font-weight="600" fill="{MUTED}">{fmt(value)}</text>'
        )
        out.append(
            f'<text x="{x:.1f}" y="{top + plot_h + 20}" text-anchor="middle" {FONT} '
            f'font-size="12" fill="{INK}">{label}</text>'
        )
    out.append("</svg>")
    return "\n".join(out)


def donut_chart(pairs, title, width, height, accent):
    cx, cy = width * 0.34, (height + 28) / 2
    radius, ring = min(cx, height - cy) - 16, 30
    total = sum(v for _, v in pairs) or 1

    out = svg_open(width, height, title)
    angle = -90.0
    for i, (label, value) in enumerate(pairs):
        sweep = 360.0 * value / total
        a0, a1 = math.radians(angle), math.radians(angle + min(sweep, 359.99))
        x0, y0 = cx + radius * math.cos(a0), cy + radius * math.sin(a0)
        x1, y1 = cx + radius * math.cos(a1), cy + radius * math.sin(a1)
        large = 1 if sweep > 180 else 0
        color = PALETTE[i % len(PALETTE)] if i else accent
        out.append(
            f'<path d="M {x0:.1f} {y0:.1f} A {radius} {radius} 0 {large} 1 '
            f'{x1:.1f} {y1:.1f}" fill="none" stroke="{color}" '
            f'stroke-width="{ring}"/>'
        )
        angle += sweep
    legend_x, legend_y = width * 0.62, cy - (len(pairs) * 24) / 2 + 8
    for i, (label, value) in enumerate(pairs):
        color = PALETTE[i % len(PALETTE)] if i else accent
        y = legend_y + i * 24
        pct = 100.0 * value / total
        out.append(f'<rect x="{legend_x}" y="{y - 10}" width="12" height="12" rx="3" fill="{color}"/>')
        out.append(
            f'<text x="{legend_x + 20}" y="{y}" {FONT} font-size="13" fill="{INK}">'
            f"{label} — {fmt(value)} ({pct:.0f}%)</text>"
        )
    out.append("</svg>")
    return "\n".join(out)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--type", required=True, choices=["bar", "line", "donut"])
    parser.add_argument("--title", required=True)
    parser.add_argument("--data", required=True, help='Comma-separated "Label:value" pairs')
    parser.add_argument("--out", help="Output SVG path (default: stdout)")
    parser.add_argument("--width", type=int, default=640)
    parser.add_argument("--height", type=int, default=0, help="Default: auto by type")
    parser.add_argument("--accent", default=ACCENT, help="Accent hex color")
    args = parser.parse_args()

    pairs = parse_data(args.data)
    height = args.height or {"bar": 60 + 44 * len(pairs), "line": 320, "donut": 280}[args.type]
    render = {"bar": bar_chart, "line": line_chart, "donut": donut_chart}[args.type]
    svg = render(pairs, args.title, args.width, height, args.accent)

    if args.out:
        with open(args.out, "w", encoding="utf-8") as handle:
            handle.write(svg + "\n")
        print(f"OK: wrote {args.out} ({len(pairs)} data points)")
    else:
        print(svg)


if __name__ == "__main__":
    main()
