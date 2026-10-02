# Issue Standards

An issue is the **single source of truth for *intent***: why a change is wanted and the acceptance criteria that define "done" (see the change lifecycle in AGENTS.md for how issues, commits, and PRs relate).

## Issue Classification & Selection

When the user's request involves creating a GitHub issue, you MUST classify it into exactly one of the following types and select the corresponding Issue Type when prompted:

| Type | When to use |
|------|-------------|
| **DEFECT** | A bug, regression, error, failure, or any deviation from expected behavior |
| **FEATURE** | A new piece of functionality or enhancement that delivers value to the end-user |
| **TASK** | Internal work that is neither a defect nor a feature (maintenance, docs, refactoring, research, setup, etc.) |

### Default to small, focused tickets

Every ticket is a **DEFECT**, **FEATURE**, or **TASK**. Keep each one narrow and well-defined.

- **Start small, grow as needed.** A ticket that tries to capture too much scope upfront tends to stall, accumulate ambiguity, and resist estimation. Split broad work into several focused tickets rather than one sprawling one.
- **When in doubt, create a FEATURE or TASK** with a narrow, well-defined scope. It is much easier to split related tickets out later than to break an oversized one back down into actionable pieces.

### No meta-tickets

Since an issue's only job is to carry *intent*, every issue must add intent of its own. A **meta-ticket** (an issue whose only purpose is to group or point at other issues without contributing its own description or acceptance criteria) adds nothing and is just noise. Do not create them. Before creating any "umbrella" or "tracking" ticket, ask whether it carries its own substance (a distinct description, its own acceptance criteria). If not, drop it and let each ticket stand on its own.

- **Each ticket must stand on its own.** A handful of related fixes do not need a parent to tie them together: create them as independent DEFECT/FEATURE/TASK tickets, each with a well-defined scope.
- **Don't fragment issues to mirror PRs or commits.** One PR may close several issues, and one issue may span several commits (see the change lifecycle in AGENTS.md). Never invent an extra parent/umbrella issue just to bundle work that will ship together: the PR is the bundle, and the commits are the record. Conversely, don't split work into separate PRs just because it touches separate issues.

### Parent / sub-issue relationships

Parent ↔ sub-issue links are fine, but **only** when the parent carries real intent of its own and the split serves an organizational purpose:

- **The parent is the SSOT of intent for the grouped work.** The parent ticket holds the description and defines the acceptance criteria. Sub-issues exist to break the work into chunks that can be implemented incrementally or in parallel by multiple devs, and to make that progress visible on the project board.
- **Sub-issues are organizational proxies and carry no Issue Type.** They exist purely for visibility on the project board; the parent (typically a FEATURE) remains the source of truth that defines the ACs. Do not duplicate or fragment the acceptance criteria across the sub-issues.
- **If the parent adds nothing, remove it.** When the would-be parent has no intent of its own and every sub-ticket can stand alone, delete the parent and keep the standalone tickets.

### Classify by the actual problem

Choose the Issue Type that matches what the ticket really describes, so the form's sections fit the content:

- A ticket describing a **correctness issue** (something behaving incorrectly) is a **DEFECT**, even if it was first framed as a generic task. As a DEFECT, its Observed/Expected Behavior sections capture the problem cleanly.
- Watch for a misclassified ticket whose "Why?" section is really just describing observed-vs-expected behavior: that is a signal it should be a DEFECT rather than a FEATURE/TASK.

### Other rules

- Do not add manual labels (e.g., `bug`, `priority`) unless the user specifically requests you to do so. Issue Types already apply their own labels automatically.
- Triage is signaled by milestone, not a label. An issue with no milestone assigned is implicitly untriaged; assigning a milestone is the triage action itself. Do not add a `needs triage` (or similar) label to the issue forms or to individual issues.

---

## Issue form templates

Each issue type has a GitHub issue form. The canonical sources live under
[`src/`](../src/); GitHub uses the generated copies under
`.github/ISSUE_TEMPLATE/`:

| Type | Canonical source | Generated form |
|------|------------------|----------------|
| DEFECT | `src/issue-defect.yml` | `.github/ISSUE_TEMPLATE/1-defect.yml` |
| FEATURE | `src/issue-feature.yml` | `.github/ISSUE_TEMPLATE/2-feature.yml` |
| TASK | `src/issue-task.yml` | `.github/ISSUE_TEMPLATE/3-task.yml` |

The fields, their guidance, and the worked examples are defined in those sources; treat them as authoritative. When filing an issue programmatically (e.g. via `gh issue create`), derive the body from the matching form's fields rather than inventing a structure.

To change a form, edit its source and run `task generate`; never edit the generated file (see [Contributing](./contributing.md)).

---

## Determining the version for defect tickets

When creating a DEFECT ticket, you MUST resolve the version before drafting the issue:

1. **Fetch the latest release** from the GitHub API:
   ```
   gh api repos/{owner}/{repo}/releases/latest --jq '.tag_name + " (" + .target_commitish + ")"'
   ```
2. **If `target_commitish` is a branch name** (not a SHA), resolve it to a short SHA:
   ```
   gh api repos/{owner}/{repo}/commits/{branch} --jq '.sha[:8]'
   ```
3. **Ask the user to confirm** the resolved version is correct before submitting: the defect may have been observed on an older release.

Default to the latest release. Only use a different version if the user specifies one.
