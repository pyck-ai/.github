# pyck-ai/.github

Org-wide GitHub defaults plus the badge generator. See [README.md](README.md).

## Commands

- `task generate`: regenerate `badges/`.
- `task ci`: check that everything is in order, i.e. the badges are in sync (CI runs this).

## Gotchas

- The issue forms and the PR template are edited in place; nothing generates them.
- Everything under `.github/` applies to every pyck-ai repo without its own copy, so review a change there as an org-wide change.
- `badges/*.svg` and `badges/README.md` are generated; never hand-edit them.
