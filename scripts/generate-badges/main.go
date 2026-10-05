// Command generate-badges generates the static SVG badges under badges/ and
// the gallery badges/README.md.
//
// The badges look like GitHub Copilot's severity badges (62x18 outline pill,
// text as glyph outlines, so they render identically everywhere and need no
// font on the viewer's machine). Comments link them through
//
//	https://raw.githubusercontent.com/pyck-ai/.github/main/badges/<group>/<name>-<light|dark>.svg
//
// The glyph outlines come from Mona Sans SemiBold (GitHub's font, OFL-1.1),
// downloaded at generate time from a pinned commit of github/mona-sans, sha256
// verified, and cached under the user cache dir. The font is not committed.
// Size 12 px with 0.17 px tracking reproduces Copilot's "High", "Low" and
// "Medium" ink boxes to within 0.04 px.
//
// Run it from the repo root (the Taskfile does this):
//
//	go run ./scripts/generate-badges            # write outputs
//	go run ./scripts/generate-badges --check    # verify; exit 1 if stale
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

const generator = "scripts/generate-badges/main.go"

const (
	outDir = "badges"

	// Pinned font: github/mona-sans tag v2.0.27.
	fontCommit = "0f7dc66ddd766605eb0e75c3f47bf9d1dd38ceca"
	fontPath   = "fonts/static/ttf/MonaSans-SemiBold.ttf"
	fontSHA256 = "67636cd715fc8637df3d0f215f478af37924fa3f2c9ee0d3ed8f3147211811f2"

	rawBase = "https://raw.githubusercontent.com/pyck-ai/.github/main/badges"

	fontSize  = 12.0 // px; SemiBold at 12 px matches Copilot's cap height (8.748) and stem
	tracking  = 0.17 // px added between glyphs, fits Copilot's ink widths
	height    = 18   // badge height in px
	padding   = 8.5  // px between the pill edge and the text's ink box
	minWidth  = 18   // smallest badge width
	centreY   = 9.0  // Copilot centres the ink box vertically in the pill
	gridWidth = 10   // confidence badges per gallery row
	darkBG    = "#0d1117"
	lightBG   = "#ffffff"
	white     = "#ffffff"
)

// ---------------------------------------------------------------------------
// Colours

// palette holds a Primer functional colour in both modes.
type palette struct{ light, dark string }

var (
	red    = palette{"#d1242f", "#f85149"}
	orange = palette{"#bc4c00", "#db6d28"}
	yellow = palette{"#9a6700", "#d29922"}
	blue   = palette{"#0969da", "#4493f8"}
	green  = palette{"#1a7f37", "#3fb950"}
	grey   = palette{"#59636e", "#9198a1"}
)

func (p palette) mode(m string) string {
	if m == "dark" {
		return p.dark
	}
	return p.light
}

// confidenceAnchors are the gradient stops: percent and colour.
var confidenceAnchors = []struct {
	pct int
	col palette
}{{0, red}, {50, orange}, {70, yellow}, {85, green}, {100, green}}

func parseHex(s string) [3]float64 {
	var c [3]float64
	for i := range c {
		v, err := strconv.ParseUint(s[1+2*i:3+2*i], 16, 8)
		if err != nil {
			panic(err)
		}
		c[i] = float64(v) / 255
	}
	return c
}

func toLinear(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func fromLinear(v float64) float64 {
	if v <= 0.0031308 {
		return 12.92 * v
	}
	return 1.055*math.Pow(v, 1/2.4) - 0.055
}

// luminance is the WCAG relative luminance of a hex colour.
func luminance(hexc string) float64 {
	c := parseHex(hexc)
	return 0.2126*toLinear(c[0]) + 0.7152*toLinear(c[1]) + 0.0722*toLinear(c[2])
}

// contrast is the WCAG contrast ratio between two hex colours.
func contrast(a, b string) float64 {
	la, lb := luminance(a), luminance(b)
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

type lch struct{ l, c, h float64 } // OKLCH, h in degrees

func hexToLCH(hexc string) lch {
	c := parseHex(hexc)
	r, g, b := toLinear(c[0]), toLinear(c[1]), toLinear(c[2])
	l := math.Cbrt(0.4122214708*r + 0.5363325363*g + 0.0514459929*b)
	m := math.Cbrt(0.2119034982*r + 0.6806995451*g + 0.1073969566*b)
	s := math.Cbrt(0.0883024619*r + 0.2817188376*g + 0.6299787005*b)
	L := 0.2104542553*l + 0.7936177850*m - 0.0040720468*s
	A := 1.9779984951*l - 2.4285922050*m + 0.4505937099*s
	B := 0.0259040371*l + 0.7827717662*m - 0.8086757660*s
	h := math.Atan2(B, A) * 180 / math.Pi
	if h < 0 {
		h += 360
	}
	return lch{L, math.Hypot(A, B), h}
}

// lchLinear converts OKLCH to linear sRGB (unclamped).
func lchLinear(x lch) [3]float64 {
	a := x.c * math.Cos(x.h*math.Pi/180)
	b := x.c * math.Sin(x.h*math.Pi/180)
	l := x.l + 0.3963377774*a + 0.2158037573*b
	m := x.l - 0.1055613458*a - 0.0638541728*b
	s := x.l - 0.0894841775*a - 1.2914855480*b
	l, m, s = l*l*l, m*m*m, s*s*s
	return [3]float64{
		4.0767416621*l - 3.3077115913*m + 0.2309699292*s,
		-1.2684380046*l + 2.6097574011*m - 0.3413193965*s,
		-0.0041960863*l - 0.7034186147*m + 1.7076147010*s,
	}
}

func inGamut(c [3]float64) bool {
	const eps = 1e-6
	for _, v := range c {
		if v < -eps || v > 1+eps {
			return false
		}
	}
	return true
}

// lchToHex converts OKLCH to a hex colour, reducing chroma until it is inside
// the sRGB gamut.
func lchToHex(x lch) string {
	lo, hi := 0.0, x.c
	if !inGamut(lchLinear(x)) {
		for i := 0; i < 40; i++ {
			mid := (lo + hi) / 2
			if inGamut(lchLinear(lch{x.l, mid, x.h})) {
				lo = mid
			} else {
				hi = mid
			}
		}
		x.c = lo
	}
	c := lchLinear(x)
	var out [3]int
	for i, v := range c {
		out[i] = int(math.Round(math.Min(1, math.Max(0, fromLinear(math.Min(1, math.Max(0, v))))) * 255))
	}
	return fmt.Sprintf("#%02x%02x%02x", out[0], out[1], out[2])
}

func mix(a, b lch, t float64) lch {
	dh := b.h - a.h
	for dh > 180 {
		dh -= 360
	}
	for dh < -180 {
		dh += 360
	}
	return lch{a.l + (b.l-a.l)*t, a.c + (b.c-a.c)*t, math.Mod(a.h+dh*t+360, 360)}
}

// confidenceColour interpolates the gradient in OKLCH for a percentage and
// mode. If the result has less than 3:1 contrast against the page background
// it is darkened (light) or lightened (dark) until it reaches 4.5:1.
func confidenceColour(pct int, mode string) string {
	bg := lightBG
	if mode == "dark" {
		bg = darkBG
	}
	i := 0
	for i < len(confidenceAnchors)-2 && pct >= confidenceAnchors[i+1].pct {
		i++
	}
	a, b := confidenceAnchors[i], confidenceAnchors[i+1]
	t := float64(pct-a.pct) / float64(b.pct-a.pct)
	x := mix(hexToLCH(a.col.mode(mode)), hexToLCH(b.col.mode(mode)), t)
	col := lchToHex(x)
	if contrast(col, bg) >= 3 {
		return col
	}
	step := -0.005
	if mode == "dark" {
		step = 0.005
	}
	for contrast(col, bg) < 4.5 && x.l > 0 && x.l < 1 {
		x.l += step
		col = lchToHex(x)
	}
	return col
}

// ---------------------------------------------------------------------------
// Font and glyph outlines

var fontData []byte

// loadFont returns the pinned font, from the user cache or downloaded.
func loadFont() (*sfnt.Font, error) {
	cacheRoot, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(cacheRoot, "pyck-github-badges")
	file := filepath.Join(dir, fontCommit[:12]+"-"+filepath.Base(fontPath))

	ok := func(b []byte) bool {
		sum := sha256.Sum256(b)
		return hex.EncodeToString(sum[:]) == fontSHA256
	}
	b, err := os.ReadFile(file)
	if err != nil || !ok(b) {
		url := "https://raw.githubusercontent.com/github/mona-sans/" + fontCommit + "/" + fontPath
		fmt.Fprintln(os.Stderr, "downloading", url)
		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("download %s: %s", url, resp.Status)
		}
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if !ok(b) {
			return nil, fmt.Errorf("download %s: sha256 mismatch", url)
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
		tmp := file + ".tmp"
		if err := os.WriteFile(tmp, b, 0o644); err != nil {
			return nil, err
		}
		if err := os.Rename(tmp, file); err != nil {
			return nil, err
		}
	}
	fontData = b
	return sfnt.Parse(b)
}

// seg is one path segment; coordinates are in px, y down, baseline y=0.
type seg struct {
	op  byte // 'M', 'L', 'Q', 'C'
	pts [][2]float64
}

// text is laid-out text: segments placed along the baseline, and the ink box.
type text struct {
	segs       []seg
	minX, maxX float64
	minY, maxY float64 // y down, baseline 0
}

func (t text) inkWidth() float64 { return t.maxX - t.minX }

// quadExtreme returns the t in (0,1) where a quadratic Bezier with 1-D
// control values p0,p1,p2 has its extremum, or -1.
func quadExtreme(p0, p1, p2 float64) float64 {
	d := p0 - 2*p1 + p2
	if d == 0 {
		return -1
	}
	t := (p0 - p1) / d
	if t <= 0 || t >= 1 {
		return -1
	}
	return t
}

func layout(ft *sfnt.Font, s string) (text, error) {
	var buf sfnt.Buffer
	upm := fixedPPEM(ft)
	scale := fontSize / float64(ft.UnitsPerEm())
	t := text{minX: math.Inf(1), minY: math.Inf(1), maxX: math.Inf(-1), maxY: math.Inf(-1)}
	grow := func(x, y float64) {
		t.minX, t.maxX = math.Min(t.minX, x), math.Max(t.maxX, x)
		t.minY, t.maxY = math.Min(t.minY, y), math.Max(t.maxY, y)
	}
	pen := 0.0
	var prev sfnt.GlyphIndex
	for i, r := range s {
		g, err := ft.GlyphIndex(&buf, r)
		if err != nil || g == 0 {
			return t, fmt.Errorf("no glyph for %q", r)
		}
		if i > 0 {
			k, err := ft.Kern(&buf, prev, g, upm, 0)
			if err == nil {
				pen += float64(k) / 64 * scale
			}
			pen += tracking
		}
		segs, err := ft.LoadGlyph(&buf, g, upm, nil)
		if err != nil {
			return t, err
		}
		var cx, cy float64
		for _, sg := range segs {
			n := 1
			op := byte('M')
			switch sg.Op {
			case sfnt.SegmentOpMoveTo:
			case sfnt.SegmentOpLineTo:
				op = 'L'
			case sfnt.SegmentOpQuadTo:
				op, n = 'Q', 2
			case sfnt.SegmentOpCubeTo:
				op, n = 'C', 3
			}
			out := seg{op: op}
			for j := 0; j < n; j++ {
				out.pts = append(out.pts, [2]float64{
					pen + float64(sg.Args[j].X)/64*scale,
					float64(sg.Args[j].Y) / 64 * scale, // sfnt already returns y down
				})
			}
			end := out.pts[n-1]
			if op == 'Q' {
				// Exact extrema of the curve, not just its control points.
				for axis := 0; axis < 2; axis++ {
					p0 := [2]float64{cx, cy}[axis]
					if tt := quadExtreme(p0, out.pts[0][axis], out.pts[1][axis]); tt > 0 {
						u := 1 - tt
						x := u*u*cx + 2*u*tt*out.pts[0][0] + tt*tt*out.pts[1][0]
						y := u*u*cy + 2*u*tt*out.pts[0][1] + tt*tt*out.pts[1][1]
						grow(x, y)
					}
				}
			} else if op == 'C' {
				// Not used by the pinned TrueType font; control-point box.
				grow(out.pts[0][0], out.pts[0][1])
				grow(out.pts[1][0], out.pts[1][1])
			}
			grow(end[0], end[1])
			cx, cy = end[0], end[1]
			t.segs = append(t.segs, out)
		}
		adv, err := ft.GlyphAdvance(&buf, g, upm, 0)
		if err != nil {
			return t, err
		}
		pen += float64(adv) / 64 * scale
		prev = g
	}
	return t, nil
}

// fixedPPEM returns a ppem equal to the font's units per em, in 26.6 fixed
// point, so glyph coordinates come back in font units (full precision).
func fixedPPEM(ft *sfnt.Font) fixed.Int26_6 { return fixed.Int26_6(ft.UnitsPerEm()) << 6 }

func num(v float64) string {
	v = math.Round(v*1000) / 1000
	if v == 0 {
		return "0"
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// pathData renders the text as SVG path data shifted by (dx, dy).
func (t text) pathData(dx, dy float64) string {
	var sb strings.Builder
	open := false
	for _, s := range t.segs {
		if s.op == 'M' {
			if open {
				sb.WriteString(" Z ")
			}
			open = true
		} else {
			sb.WriteByte(' ')
		}
		sb.WriteByte(s.op)
		for i, p := range s.pts {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(num(p[0]+dx) + " " + num(p[1]+dy))
		}
	}
	if open {
		sb.WriteString(" Z")
	}
	return strings.TrimSpace(sb.String())
}

// ---------------------------------------------------------------------------
// Badges

// badge is one badge in both modes.
type badge struct {
	group, name string // output: badges/<group>/<name>-<mode>.svg
	label       string // visible text
	title       string // <title> and aria-label
	colour      func(mode string) string
	filled      bool // filled pill with white text (Critical)
}

type group struct {
	name, heading, about string
	badges               []badge
	// shareWith, when non-empty, is a label whose width every badge shares
	// (at least).
	shareWith string
}

func fixedColour(p palette) func(string) string { return func(m string) string { return p.mode(m) } }

func buildGroups() []group {
	sev := func(name, label string, p palette, filled bool) badge {
		return badge{"severity", name, label, label + " severity", fixedColour(p), filled}
	}
	simple := func(grp, prefix, name, label string, p palette) badge {
		return badge{grp, strings.ToLower(name), label, prefix + label, fixedColour(p), false}
	}
	var conf []badge
	for i := 0; i <= 100; i++ {
		pct := i
		conf = append(conf, badge{"confidence", strconv.Itoa(pct), fmt.Sprintf("%d%%", pct),
			fmt.Sprintf("Confidence %d%%", pct), func(m string) string { return confidenceColour(pct, m) }, false})
	}
	conf = append(conf,
		badge{"confidence", "rule", "Rule", "Confidence: rule check", fixedColour(green), false},
		badge{"confidence", "ai", "AI", "Confidence: AI, no score", fixedColour(grey), false})

	return []group{
		{name: "severity", heading: "Severity", about: "Critical is filled with white text; the others are outlines.",
			badges: []badge{
				sev("critical", "Critical", red, true),
				sev("high", "High", red, false),
				sev("medium", "Medium", yellow, false),
				sev("low", "Low", blue, false),
				sev("info", "Info", grey, false),
			}},
		{name: "status", heading: "Status", about: "Lifecycle of a finding.",
			badges: []badge{
				simple("status", "Status: ", "fixed", "Fixed", green),
				simple("status", "Status: ", "new", "New", blue),
				simple("status", "Status: ", "open", "Open", yellow),
			}},
		{name: "kind", heading: "Kind", about: "Defect, Feature and Task use the org's Issue Type colours.",
			badges: []badge{
				simple("kind", "Kind: ", "defect", "Defect", red),
				simple("kind", "Kind: ", "feature", "Feature", blue),
				simple("kind", "Kind: ", "task", "Task", yellow),
				simple("kind", "Kind: ", "breaking", "Breaking", orange),
			}},
		{name: "confidence", heading: "Confidence", shareWith: "100%",
			about: "`0` to `100` colour by value: a continuous OKLCH gradient, red at 0%, orange at 50%, yellow at 70%, green from 85%. " +
				"85% and up is treated as sure, 50 to 85% as a suggestion. `rule` (green) is a fixed rule check, certain; `ai` (grey) has no number.",
			badges: conf},
	}
}

// groupWidth is the shared width of a group's badges: the widest ink box plus
// padding, rounded up to an integer, at least minWidth. For a group with
// shareWith, that label's width is the floor.
func groupWidth(ft *sfnt.Font, g group) (int, error) {
	labels := []string{}
	for _, b := range g.badges {
		labels = append(labels, b.label)
	}
	if g.shareWith != "" {
		labels = append(labels, g.shareWith)
	}
	w := 0.0
	for _, l := range labels {
		t, err := layout(ft, l)
		if err != nil {
			return 0, err
		}
		w = math.Max(w, t.inkWidth())
	}
	return max(minWidth, int(math.Ceil(w+2*padding))), nil
}

func xmlEscape(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;").Replace(s)
}

// renderSVG returns one badge for a mode at the given width.
func renderSVG(ft *sfnt.Font, b badge, mode string, width int, fill string) (string, error) {
	t, err := layout(ft, b.label)
	if err != nil {
		return "", err
	}
	// Centre the ink box in the pill, horizontally and vertically.
	dx := float64(width)/2 - (t.minX+t.maxX)/2
	dy := centreY - (t.minY+t.maxY)/2
	colour := b.colour(mode)
	rectFill, textFill := "none", colour
	if b.filled {
		rectFill, textFill = fill, white
		colour = fill
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" role="img" aria-label="%s">`+"\n", width, height, xmlEscape(b.title))
	fmt.Fprintf(&sb, "  <!-- Code generated by %s; DO NOT EDIT. -->\n", generator)
	fmt.Fprintf(&sb, "  <title>%s</title>\n", xmlEscape(b.title))
	fmt.Fprintf(&sb, `  <rect x="0.5" y="0.5" width="%d" height="17" rx="8.5" fill="%s" stroke="%s"/>`+"\n", width-1, rectFill, colour)
	fmt.Fprintf(&sb, `  <path fill="%s" d="%s"/>`+"\n", textFill, t.pathData(dx, dy))
	sb.WriteString("</svg>\n")
	return sb.String(), nil
}

// filledBackground returns the fill for a filled badge: the mode's colour if
// white text on it reaches 4.5:1, else the light-mode colour.
func filledBackground(b badge, mode string) string {
	c := b.colour(mode)
	if contrast(white, c) >= 4.5 {
		return c
	}
	return b.colour("light")
}

// ---------------------------------------------------------------------------
// Gallery

func picture(rel string, b badge, width int) string {
	u := rel + "/" + b.group + "/" + b.name
	return fmt.Sprintf(`<picture><source media="(prefers-color-scheme: dark)" srcset="%s-dark.svg"><source media="(prefers-color-scheme: light)" srcset="%s-light.svg"><img alt="%s" src="%s-light.svg" width="%d" height="%d" align="texttop"></picture>`,
		u, u, xmlEscape(b.title), u, width, height)
}

func gallery(groups []group, widths map[string]int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "<!-- Code generated by %s; DO NOT EDIT. -->\n\n", generator)
	sb.WriteString("# Badges\n\n")
	sb.WriteString("Static outline badges in the look of GitHub Copilot's severity badges, for comments and docs.\n")
	sb.WriteString("Text is glyph outlines (Mona Sans SemiBold), so they render the same everywhere. Each badge has a light and a dark SVG.\n")
	sb.WriteString("Generated by [`" + generator + "`](../" + generator + "); regenerate with `task generate:badges`.\n\n")
	sb.WriteString("## Use\n\n")
	sb.WriteString("URL pattern (`<mode>` is `light` or `dark`):\n\n")
	sb.WriteString("```text\n" + rawBase + "/<group>/<name>-<mode>.svg\n```\n\n")
	sb.WriteString("Snippet (picks the SVG for the viewer's theme, falls back to light; `width` is the group's width below):\n\n")
	sb.WriteString("```html\n")
	sb.WriteString("<picture>\n")
	fmt.Fprintf(&sb, "  <source media=\"(prefers-color-scheme: dark)\" srcset=\"%s/severity/high-dark.svg\">\n", rawBase)
	fmt.Fprintf(&sb, "  <source media=\"(prefers-color-scheme: light)\" srcset=\"%s/severity/high-light.svg\">\n", rawBase)
	fmt.Fprintf(&sb, "  <img alt=\"High severity\" src=\"%s/severity/high-light.svg\" width=\"%d\" height=\"%d\" align=\"texttop\">\n", rawBase, widths["severity"], height)
	sb.WriteString("</picture>\n```\n\n")
	sb.WriteString("The raw URLs only resolve once the files are on `main`. The tables below use relative paths so the gallery renders on any branch.\n\n")
	for _, g := range groups {
		fmt.Fprintf(&sb, "## %s\n\n%s Width %d, height %d.\n\n", g.heading, g.about, widths[g.name], height)
		if g.name == "confidence" {
			sb.WriteString("|" + strings.Repeat(" |", gridWidth) + "\n|" + strings.Repeat("---|", gridWidth) + "\n")
			for i := 0; i < len(g.badges); i += gridWidth {
				sb.WriteString("|")
				for j := i; j < i+gridWidth; j++ {
					if j < len(g.badges) {
						sb.WriteString(" " + picture(".", g.badges[j], widths[g.name]))
					}
					sb.WriteString(" |")
				}
				sb.WriteString("\n")
			}
			sb.WriteString("\n")
			continue
		}
		sb.WriteString("| Badge | Path (without `-<mode>.svg`) | Title |\n|---|---|---|\n")
		for _, b := range g.badges {
			fmt.Fprintf(&sb, "| %s | `%s/%s` | %s |\n", picture(".", b, widths[g.name]), b.group, b.name, b.title)
		}
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n") + "\n"
}

// ---------------------------------------------------------------------------
// Main

func main() {
	check := false
	for _, a := range os.Args[1:] {
		if a == "--check" {
			check = true
		}
	}
	if err := run(check); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}
}

func run(check bool) error {
	ft, err := loadFont()
	if err != nil {
		return err
	}
	groups := buildGroups()
	files := map[string]string{} // path -> content
	widths := map[string]int{}
	minContrast := map[string]float64{"light": math.Inf(1), "dark": math.Inf(1)}
	minName := map[string]string{}
	bgs := map[string]string{"light": lightBG, "dark": darkBG}

	for _, g := range groups {
		w, err := groupWidth(ft, g)
		if err != nil {
			return err
		}
		widths[g.name] = w
		for _, b := range g.badges {
			for _, mode := range []string{"light", "dark"} {
				fill := ""
				if b.filled {
					fill = filledBackground(b, mode)
					if c := contrast(white, fill); c < 4.5 {
						return fmt.Errorf("%s/%s %s: white text contrast %.2f < 4.5", b.group, b.name, mode, c)
					}
				}
				svg, err := renderSVG(ft, b, mode, w, fmt.Sprint(fill))
				if err != nil {
					return fmt.Errorf("%s/%s: %w", b.group, b.name, err)
				}
				files[filepath.Join(outDir, b.group, b.name+"-"+mode+".svg")] = svg
				if b.filled {
					continue // its pill colour is a fill, checked against its text above
				}
				if c := contrast(b.colour(mode), bgs[mode]); c < minContrast[mode] {
					minContrast[mode], minName[mode] = c, b.group+"/"+b.name
				}
			}
		}
	}
	files[filepath.Join(outDir, "README.md")] = gallery(groups, widths)

	// Existing generated files that are no longer expected are stale.
	var stale []string
	_ = filepath.WalkDir(outDir, func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			if _, ok := files[p]; !ok {
				stale = append(stale, p)
			}
		}
		return nil
	})

	var paths []string
	for p := range files {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	var outOfSync []string
	wrote := 0
	for _, p := range paths {
		cur, err := os.ReadFile(p)
		if err == nil && bytes.Equal(cur, []byte(files[p])) {
			continue
		}
		if check {
			outOfSync = append(outOfSync, p)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(p, []byte(files[p]), 0o644); err != nil {
			return err
		}
		wrote++
	}
	sort.Strings(stale)
	if check {
		outOfSync = append(outOfSync, stale...)
	} else {
		for _, p := range stale {
			if err := os.Remove(p); err != nil {
				return err
			}
			fmt.Println("removed", p)
		}
	}

	fmt.Printf("badges: %d files, min contrast against page: light %.2f:1 (%s), dark %.2f:1 (%s)\n", len(paths),
		minContrast["light"], minName["light"], minContrast["dark"], minName["dark"])
	if !check {
		fmt.Printf("wrote %d changed files\n", wrote)
		return nil
	}
	if len(outOfSync) > 0 {
		sort.Strings(outOfSync)
		fmt.Println("Generated badges are out of sync with the generator:")
		for i, p := range outOfSync {
			if i == 10 {
				fmt.Printf("  ... and %d more\n", len(outOfSync)-10)
				break
			}
			fmt.Println("  -", p)
		}
		fmt.Println("Run: task generate:badges")
		os.Exit(1)
	}
	fmt.Println("All generated badges are in sync.")
	return nil
}
