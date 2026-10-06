# Pull Request Standards

## What a PR is for

A pull request is a **mechanism for review**, nothing more; see the change
lifecycle in AGENTS.md for how it relates to issues and commits. Its
description and discussion exist to help reviewers understand and evaluate the
change; they are not the historical record. The durable rationale for *why* a
change was made lives in the commit message(s) (see git-standards.md). A PR is
**never** an SSOT: if something matters for the future, it belongs in a commit,
not the PR description.

Because the PR only serves the review, its description is just a **short
executive summary** for reviewers. A PR may bundle several commits and close
several issues; one PR does not mean one issue.

Most of the rules below are checked by the
[pyck-review](https://github.com/pyck-ai/pyck-review) GitHub App; see
[Classification](#classification) and [Guideline check](#guideline-check).

## The PR template

The PR body uses the org template. Its canonical source is
[`src/pull_request.md`](../src/pull_request.md); GitHub injects the generated
copy at `.github/PULL_REQUEST_TEMPLATE.md`. Use the template as is: keep its
two-part layout (see [Description structure](#description-structure)), replace
the HTML comments with real content, and remove every instruction comment
before the PR is ready. Add no other top-level headings; structure extra
information with sub-headings under `## Additional Information`.

The template's own comments are the point-of-use guidance for filling it in:

- **Summary**: see [Summary](#summary).
- **Linked issues**: in the summary part, one issue per line, using one of the
  three trailers defined in git-standards.md's
  [Footer](./git-standards.md#footer) rule.
- **Review readiness**: there are no compliance or testing checkboxes. A PR
  stays in draft while the work is in progress; the author marks it ready for
  review only after testing the changes (manually, where applicable) and adding
  or updating the tests the change requires. The ready-for-review transition
  **is** the testing attestation, and reviewers should treat it that way.

To change the template, edit the source and run `task generate`; never edit the
generated file (see [Contributing](./contributing.md)).

## Description structure

The structure is always **summary + additional information**:

1. The **summary** comes first and has no headings (see [Summary](#summary)).
   Linked issues go here, one `Closes:` / `Part-of:` / `See-also:` per line.
2. Then exactly one `## Additional Information`. Below it is free text: the
   author may put whatever helps the review there, structured only with `###`
   or deeper sub-headings. `N/A` if there is nothing to add.
3. No other `#` or `##` headings anywhere in the description.

## Summary

The summary is a short executive summary, not a copy of the commit messages and
not a per-commit changelog. Summarize, don't copy:

- It describes everything the commits change; one sentence may cover several
  commits.
- It does not claim changes that no commit makes.
- Stay on the shorter side: aim for 2 or 3 paragraphs, and add more only when
  the PR really needs them to be understood or really makes that many changes.
  **5 paragraphs at most.** A list counts as one paragraph; issue-link lines do
  not count. Details reviewers might want go under `## Additional Information`.

Exception: a PR with a **single commit** may use that commit's message 1:1 as
its summary (GitHub pre-fills the description with it). The length limit does
not apply to such a copy.

## Size

Keep a PR small enough to review in one sitting: **aim for 5 commits or
fewer, and 10 at most** (merge commits do not count). Split a larger change
into several PRs and stack them: the first targets `main`, each next one
targets the branch of the PR below it, so each review shows only its own
layer. See GitHub's
[About stacked pull requests](https://docs.github.com/en/pull-requests/get-started/about-stacked-prs)
(public preview); plain PRs based on each other's branches work the same way.

When a lower PR merges and its branch is deleted, GitHub retargets the PR
above it to the lower PR's base. After a squash or rebase merge the landed
commits get new ids, so rebase the upper branch onto the new base
(`git rebase --onto main <lower-branch>`) before it can merge cleanly.

## Classification

Change classification is automated, not self-reported, so the PR has no
classification checkboxes (the old "PR Compliance" ones are gone). It is done by
[`pyck-ai/pyck-review`](https://github.com/pyck-ai/pyck-review) (a GitHub App,
run centrally for the org). A repo turns it on with
`pullRequests: { commits: true }` in its `.github/pyck-review.json5`.

| Aspect | Behavior |
|---|---|
| Unit | Each commit separately, from its message and its own diff |
| Kind | Exactly one of Defect fix, Feature or Task per commit (every issue and every commit is exactly one kind; a PR may hold several commits) |
| Also per commit | Whether it breaks backwards compatibility (a `BREAKING CHANGE:` footer line decides it), and whether a revert cannot undo it (Irreversible: it deletes or rewrites data, removes infrastructure, or has an effect outside the system) |
| Context | PR title, description, other commit subjects, and the types of the linked issues; a linked issue's type is strong evidence |
| Where | A collapsed commits section in the PR's one pyck-review review |
| Mismatch | A commit whose kind matches none of the linked issues' types is a tip in the review |
| Status | Advice; never blocks merging |

Implementation details live in pyck-review's README, not here.

## Guideline check

Another pyck-review check. A repo turns it on with
`pullRequests: { problems: true }` (pull requests) and `issues: { problems: true }`
(issues) in its `.github/pyck-review.json5`. It runs on every ready (non-draft)
PR not opened by a bot, when it is opened, edited or pushed to, and every 10
minutes. Text inside code blocks is ignored. Commit messages are checked
against [git-standards](./git-standards.md).

| Rule | Defined in |
|---|---|
| Summary present | [Description structure](#description-structure) |
| `## Additional Information` present | [Description structure](#description-structure) |
| Only those two top-level parts | [Description structure](#description-structure) |
| Template instruction comments removed | [The PR template](#the-pr-template) |
| No old "PR Compliance" checkboxes | [Classification](#classification) |
| Issue links use only `Closes:` / `Part-of:` / `See-also:` | [git-standards: Footer](./git-standards.md#footer) |
| One issue per link line | [git-standards: Footer](./git-standards.md#footer) |
| Branch is `<issue>-<desc>` or `u/<user>/<desc>` | [git-standards: Branch naming](./git-standards.md#branch-naming) |
| Summary covers what the commits do and claims nothing they don't | [Summary](#summary) |
| Summary mentions what an Irreversible commit changes for good | [Summary](#summary) |
| Summary has at most 5 paragraphs (not for a single-commit copy) | [Summary](#summary) |
| At most 10 commits (a tip above 5, a follow-up above 7) | [Size](#size) |

The summary rules are judged by an AI model: a problem counts when the model is
at least 85% sure; from 50% it is listed as a follow-up.

What it does:

- Posts one review per PR with a table of the problems, each with a concrete
  fix. High problems add the label `guidelines-not-met` and request changes,
  unless the repo sets `pullRequests: { requestChanges: false }`; then the
  review is advice only. Whether a request for changes blocks merging depends
  on the repo's rules.
- On a later check with a different result: posts a new review listing what
  was fixed, hides the older one, and dismisses its request for changes.
- An unchanged result posts nothing.
- Issues get the same check as a comment: an issue type is set, and the
  sections of that type's form are filled in.
