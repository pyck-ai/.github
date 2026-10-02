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

## The PR template

The PR body uses the org template. Its canonical source is
[`src/pull_request.md`](../src/pull_request.md); GitHub injects the generated
copy at `.github/PULL_REQUEST_TEMPLATE.md`. Use the template verbatim, keeping
every heading and blank line, replacing only the HTML comments with real
content.

The template's own comments are the point-of-use guidance for filling it in:

- **Summary**: the short executive summary described above; one or two
  paragraphs explaining the change and why it matters, **not** a copy of the
  commit messages or a per-commit changelog.
- **Linked issues**: a separate line per reference, one of the three trailers
  defined in git-standards.md's footer rule.
- **Review readiness**: there are no compliance or testing checkboxes. A PR
  stays in draft while the work is in progress; the author marks it ready for
  review only after testing the changes (manually, where applicable) and adding
  or updating the tests the change requires. The ready-for-review transition
  **is** the testing attestation, and reviewers should treat it that way.

To change the template, edit the source and run `task generate`; never edit the
generated file (see [Contributing](./contributing.md)).

## Classification

Change classification (additive Feature/Task change, fix of something broken
(Defect), not backwards compatible (Breaking Change)) is automated, not
self-reported. A central, scheduled workflow in
[`pyck-ai/pyck-review`](https://github.com/pyck-ai/pyck-review) runs it for
the whole org, the same model Renovate uses from `pyck-ai/renovate`, using
the pyck-review GitHub App. A repo opts in by setting `classify: true` in its
`.github/pyck-review.json5`. On a ready (non-draft) PR it posts its
assessment as a sticky review comment, and it re-classifies whenever the head
commit changes, so the reviewer can see an up-to-date assessment at a glance.
It is advisory input for reviewers, not a required check, and it never
blocks a merge.
