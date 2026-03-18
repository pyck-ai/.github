# Pull Request Standards

## Pull Request Template

When creating a pull request, the PR body MUST be copied **verbatim** from the template below — including all section headings, checkbox lists (`- [ ]`), and formatting. Only the HTML comments are replaced with actual content:

~~~markdown
<!-- ---------------------------------------------------------------------------
Briefly describe what the Pull Request changes or adds. If the pull request is
related to any GitHub issues, link them like this:

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

- Write a clear description at the top explaining what the PR changes and why.
- Link related issues using `Closes: #NNN` syntax.
- Check the appropriate compliance boxes based on the nature of the changes.
- Check the testing confirmation boxes only if testing was actually performed.
- Add meaningful additional information, or "N/A" if none is needed.
