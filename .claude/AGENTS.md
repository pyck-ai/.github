# Agent Protocols

## The change lifecycle

Every change flows through the same pipeline, and each artifact has exactly **one job**:

```
idea → issue → commit → PR → review → merge
```

- **Issue** — the single source of truth for *intent*: why the change is wanted and the acceptance criteria that define "done". Issues are the mechanism that produces commits (and therefore PRs).
- **Commit** — the single source of truth for the *change itself*: what shipped and why it was done this way. It must stand entirely on its own (see GIT_STANDARDS.md).
- **PR** — purely a *mechanism for review*. Its description and discussion serve the review; they are not the historical record. A PR is **never** an SSOT — anything that must outlive the review belongs in a commit.

Two consequences fall out of this:

- **Trivial changes may skip the issue** and go straight to a commit/PR.
- **One PR does not mean one issue.** A single PR may bundle several commits and close several issues; conversely a single issue may span several commits. The mapping is flexible — don't fragment issues, commits, or PRs just to keep them one-to-one.

## Standards files

Do NOT read these files upfront. Read them on demand when the task requires it:

- **[ISSUE_STANDARDS.md](./ISSUE_STANDARDS.md)** — Read when creating, editing, or reviewing issues.
- **[PR_STANDARDS.md](./PR_STANDARDS.md)** — Read when creating, editing, or reviewing pull requests.
- **[GIT_STANDARDS.md](./GIT_STANDARDS.md)** — Read when creating, reviewing, or editing commits or branches.
