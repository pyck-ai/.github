# Community-health files

GitHub automatically applies the files under `.github/` to every repository in
the `pyck-ai` org that doesn't ship its own copy; no setup required by
consumers. Any repo can override a default by adding the same file at its own
root.

The defaults are:

- **Issue forms**: `.github/ISSUE_TEMPLATE/{1-defect,2-feature,3-task}.yml`
- **PR template**: `.github/PULL_REQUEST_TEMPLATE.md`
- **Security policy**: `.github/SECURITY.md`
- **Code of conduct**: `.github/CODE_OF_CONDUCT.md`

The issue forms and the PR template are **generated** from canonical sources in
[`src/`](../src/). See [Contributing](./contributing.md) before editing them.

## Badges

[`badges/`](../badges/README.md) is not a community-health file. It holds static
SVG badges that comments link by raw URL. It is **generated** by
[`scripts/generate-badges/`](../scripts/generate-badges/main.go) (Mona Sans
glyph outlines; the font is downloaded and cached, not committed). Regenerate
with `task generate:badges`; CI-style check: `task generate:badges:check`.
