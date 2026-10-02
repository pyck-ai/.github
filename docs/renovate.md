# Renovate

Dependency updates for the whole org run from a dedicated private repository,
[`pyck-ai/renovate`](https://github.com/pyck-ai/renovate). Nothing Renovate-related
lives here any more: no workflow, no credentials, no environment.

**To onboard a repository, commit a config at exactly this path:**

```
.github/renovate.json5
```

That is the entire opt-in. The Renovate GitHub App is installed org-wide, so no
installation change and no admin is needed. The runner matches that path
literally: a repo without it is never processed, and other locations Renovate
would normally accept (`renovate.json`, `.renovaterc`, and their variants) do
**not** onboard a repo here.

Everything else (how a run works, the dispatch cadence, credentials, the
automerge policy) is documented in the
[runner's README](https://github.com/pyck-ai/renovate), which is the source of
truth. Don't duplicate it here.
