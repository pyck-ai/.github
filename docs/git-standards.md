# Git Standards

The commit message is the **single source of truth for the change**: what shipped and why it was done this way (see the change lifecycle in AGENTS.md). It must stand entirely on its own: a reader of `git log` years from now must understand the rationale without opening the linked issue or PR. Issue/PR links are conveniences, never required reading; if a piece of context matters for the future, it belongs in the commit body.

When creating commits, follow [Tim Pope's commit message guidelines](https://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html). The key rules are:

## Subject line

- **Aim for 50 characters, at most 72**, imperative mood, capitalized first letter, no trailing period
- No prefixes (e.g., `feat:`, `fix:`): just a plain imperative sentence
- A commit that fixes a bug starts its subject with "Fix" (e.g., "Fix duplicate starts from state-change events"), still with no `fix:` prefix
- Summarize **what** the change does

## Body

- Separated from the subject by a **single blank line**
- **Wrap at 72 characters**
- Focus on **why** the change was made, not what (the diff shows what)
- Make it **self-contained**: include enough context that the reader never has to open the linked issue or PR to understand why. Distill the relevant intent from the issue into the body rather than pointing at it.
- If the explanation needs more than two paragraphs, add a bulleted TLDR summary at the top of the body

## Footer

The footer is **optional**: include it only when there is something concrete to reference or flag. A commit with no related issue and no breaking change correctly has no footer at all. Never invent a reference to satisfy a template, and never attach an issue the commit is not genuinely part of: shipping on the same branch as an issue's work does not make a commit part of that issue.

- Separated from the body by a **single blank line**
- **If** the commit relates to an issue, PR, commit, ADR, or external link, reference it on its own line, one issue id per line (repeat the trailer for several issues: `Closes: #1` then `Closes: #2`, never `Closes: #1, #2`, because GitHub does not handle several ids on one line correctly), using exactly one of three trailers: `Closes: #123` **only** when the commit actually resolves the issue (it fulfills all of the issue's acceptance criteria and is intended to close it in GitHub); `Part-of: #123` when this commit is one of several that together make up a larger change or issue, without closing it; or `See-also: #123` for a plain non-closing reference, including a reference to an issue that is **already closed** (so it isn't re-closed). No other spelling is allowed: not `Fixes`, `Fix`, `Resolves`, `Resolve`, `Implements`, `Refs`, `Ref`, `Related`, `See:`, or `See also`. See [GitHub's docs on linking a pull request to an issue](https://docs.github.com/en/issues/tracking-your-work-with-issues/using-issues/linking-a-pull-request-to-an-issue).
- If the change is not backwards compatible, add a `BREAKING CHANGE: <explanation>` section with a migration path

## Example

~~~
Centralize build caches and Docker base images

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

- **Atomic commits**: each commit represents one logical change
- **Reference related issues when there are any**: one per line in the commit footer, using the right trailer (see [Footer](#footer) above). A commit with no related issue has no footer; that is correct, not an omission
- **No WIP commits on main**: clean up before merging
