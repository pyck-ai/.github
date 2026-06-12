# Pull Request Standards

## What a PR is for

A pull request is a **mechanism for review**, nothing more. Its description and discussion exist to help reviewers understand and evaluate the change — they are not the historical record. The durable rationale for *why* a change was made lives in the commit message(s), which are the single source of truth for the change (see GIT_STANDARDS.md and the change lifecycle in AGENTS.md). A PR is **never** an SSOT: if something matters for the future, it belongs in a commit, not the PR description.

Because the PR only serves the review, its description is a **short executive summary** aimed at reviewers — not a record to be mined later, and not a copy of the commit messages. A PR may bundle several commits and close several issues; one PR does not mean one issue.

## Pull Request Template

When creating a pull request, the PR body MUST be copied **verbatim** from the template below — including all section headings, checkbox lists (`- [ ]`), and formatting. Only the HTML comments are replaced with actual content:

~~~markdown
<!-- ---------------------------------------------------------------------------
Write a SHORT executive summary of what the Pull Request changes or adds — ideally
one paragraph, at most two. This is a high-level summary for reviewers, NOT a copy
of the commit messages or a per-commit changelog. Explain the overall change and
why it matters, not every individual step.

If the pull request is related to any GitHub issues, link them like this:

Closes: #123, #456
---------------------------------------------------------------------------- -->

## PR Compliance

**Please select all that apply (zero or more):**

- [ ] This PR **adds new or improves** existing functionality *(Feature, Task)*
- [ ] This PR **fixes** something that was not working correctly *(Defect)*
- [ ] The changes in this PR are **not** backwards compatible *(Breaking Change)*


**Required: Please confirm you have completed the following testing steps:**

- [ ] I hereby confirm that **I have tested the changes manually** and they work as expected
- [ ] I hereby confirm that **I have added/updated tests** that cover the changes in this PR where necessary


## Additional Information

<!-- ---------------------------------------------------------------------------
Add any additional context, screenshots, notes, or test data that will help
reviewers better understand and test the changes.

Add "N/A" if there is no additional information.
---------------------------------------------------------------------------- -->
~~~

When filling out the PR template:

- Write a SHORT executive summary at the top — one or two paragraphs max — explaining what the PR changes and why. It must be a concise high-level summary, NOT a copy of the commit messages or a per-commit breakdown.
- Link related issues using `Closes: #NNN` syntax.
- Check the appropriate compliance boxes based on the nature of the changes.
- Check the testing confirmation boxes only if testing was actually performed.
- Add meaningful additional information, or "N/A" if none is needed.
