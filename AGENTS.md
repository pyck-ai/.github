# Agent Protocols

## The change lifecycle

Every change flows through the same pipeline, and each artifact has exactly **one job**:

```
idea → issue → commit → PR → review → merge
```

- **Issue**: the single source of truth for *intent*: why the change is wanted and the acceptance criteria that define "done". Issues are the mechanism that produces commits (and therefore PRs).
- **Commit**: the single source of truth for the *change itself*: what shipped and why it was done this way. It must stand entirely on its own (see docs/git-standards.md).
- **PR**: purely a *mechanism for review*. Its description and discussion serve the review; they are not the historical record. A PR is **never** an SSOT: anything that must outlive the review belongs in a commit.

Two consequences fall out of this:

- **Trivial changes may skip the issue** and go straight to a commit/PR.
- **One PR does not mean one issue.** A single PR may bundle several commits and close several issues; conversely a single issue may span several commits. The mapping is flexible: don't fragment issues, commits, or PRs just to keep them one-to-one.

## Documentation

Docs must never drift from the tree on disk. When a change alters what a doc
describes, fix that doc **in the same change**; this covers every doc artifact:
a repo's root `README.md`, every per-directory `README.md`, and the in-code
documentation (a module/package doc comment or header docblock) that introduces a
unit of code.

- **In-code documentation**: the doc comment that introduces a module, package, or
  public API: update it whenever that unit's purpose, exported API, invariants, or
  documented conventions change; it is the first thing a reader sees, so a stale one
  is worse than none.
- **Module / directory `README.md`**: update whenever its layout, commands, ports,
  entrypoints, or documented behavior change. A directory's own `README.md` may carry
  local maintenance rules that own its scope; follow those there.
- **Generated docs** (a `DO NOT EDIT` header): regenerate them; never hand-edit,
  and never let them drift from their source.

Two rules keep docs usable:

- **Dense, not marketing.** Short and agent-oriented: facts, tables over prose;
  prefer editing a line to adding a section.
- **A pointer, not a duplicate.** Each fact lives once; link to it instead of copying.

**Before you finish:** touched a module or package? Re-read its in-code docs and
`README.md` and repair whatever your change made stale: confirm every module, path,
command, and port a touched doc names still exists, and that any path or command you
added resolves/runs from a clean checkout.

## Standards files

Do NOT read these files upfront. Read them on demand when the task requires it:

- **[github-issue-standards.md](docs/github-issue-standards.md)**: Read when creating, editing, or reviewing issues.
- **[github-pr-standards.md](docs/github-pr-standards.md)**: Read when creating, editing, or reviewing pull requests.
- **[git-standards.md](docs/git-standards.md)**: Read when creating, reviewing, or editing commits or branches.
- **These standards are org defaults; the repo you're changing wins.** If a repo defines its own convention (its `git log`, a local `CONTRIBUTING`, or a repo-level standards doc), follow that for changes made in it, and don't impose one repo's convention on another.
