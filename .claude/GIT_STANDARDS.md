# Git Standards

When creating commits, follow [Tim Pope's commit message guidelines](https://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html). The key rules are:

## Subject line

- **50 characters max**, imperative mood, capitalized first letter, no trailing period
- No prefixes (e.g., `feat:`, `fix:`) — just a plain imperative sentence
- Summarize **what** the change does

## Body

- Separated from the subject by a **single blank line**
- **Wrap at 72 characters**
- Focus on **why** the change was made, not what (the diff shows what)
- If the explanation needs more than two paragraphs, add a bulleted TLDR summary at the top of the body

## Footer

- Separated from the body by a **single blank line**
- Reference related issues: `Closes: #123`, `Fixes: #456`, `Implements: #789`
- If the change is not backwards compatible, add a `BREAKING CHANGE: <explanation>` section with a migration path

## Example

~~~
Centralize build caches and update Docker base images

- Adopt shared base images for consistent Go and tool provisioning.
- Move Go and module caches to unified Docker volumes.
- Simplifies image maintenance and enables faster incremental builds.

The previous per-Dockerfile Go/tool provisioning led to version drift
and duplicated setup across various builder images (renovate, backend
builders, gateway). This makes builds inconsistent and harder to maintain
at scale.

This change moves all service definitions to derive from a single set
of base images, ensuring all environments use the same Go version and
tooling. It also unifies cache paths to enable faster, more reliable
incremental builds in CI.

BREAKING CHANGE: Go cache locations moved from /root/.cache to
/var/cache/go and /go/pkg/mod. Update Docker volumes accordingly.

Closes: #456
~~~

## Branch naming

When creating branches, use one of the following formats:

- **With an issue:** `[issue-id]-[short-description]` (e.g., `123-implement-dark-mode`)
- **Without an issue:** `u/[username]/[short-description]` (e.g., `u/michael/fix-urgent-bug`)

## Best practices

- **Atomic commits** — each commit represents one logical change
- **Always reference related issues** in the commit footer
- **No WIP commits on main** — clean up before merging
