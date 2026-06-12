# Issue Standards

An issue is the **single source of truth for *intent*** — why a change is wanted and the acceptance criteria that define "done". It is the mechanism that produces commits and PRs; it does not record the change itself (that is the commit's job). See the change lifecycle in AGENTS.md for how issues, commits, and PRs relate.

## Issue Classification & Selection

When the user's request involves creating a GitHub issue, you MUST classify it into exactly one of the following types and select the corresponding Issue Type when prompted:

| Type | When to use |
|------|-------------|
| **DEFECT** | A bug, regression, error, failure, or any deviation from expected behavior |
| **FEATURE** | A new piece of functionality or enhancement that delivers value to the end-user |
| **TASK** | Internal work that is neither a defect nor a feature (maintenance, docs, refactoring, research, setup, etc.) |
| **EPIC** | A large body of work spanning multiple features/tasks, representing a significant goal |
| **INITIATIVE** | The highest-level planning artifact comprising multiple epics, representing a major strategic objective |

### Default to small, focused tickets

The vast majority of tickets should be **DEFECT**, **FEATURE**, or **TASK**. Default to these types unless the user explicitly requests an Epic or Initiative.

- **Start small, grow as needed.** It is much easier to promote a cluster of related Features/Tasks into an Epic later than it is to break an oversized Epic back down into actionable pieces. A ticket that tries to capture too much scope upfront tends to stall, accumulate ambiguity, and resist estimation.
- **Use EPIC only** when the work clearly spans multiple distinct Features/Tasks that need to be tracked together under a shared goal.
- **Use INITIATIVE only** when coordinating multiple Epics toward a strategic objective — this is rare and typically driven by leadership or long-range planning.
- When in doubt, create a FEATURE or TASK with a narrow, well-defined scope. If it later becomes clear that the work is part of something bigger, promote it then.

### No meta-tickets

Since an issue's only job is to carry *intent*, every issue must add intent of its own. A **meta-ticket** — an issue whose only purpose is to group or point at other issues without contributing its own description or acceptance criteria — adds nothing and is just noise. Do not create them. Before creating any "umbrella" or "tracking" ticket, ask whether it carries its own substance (a distinct description, its own acceptance criteria). If not, drop it and let each ticket stand on its own.

- **Each ticket must stand on its own.** A handful of related fixes do not need a parent to tie them together — create them as independent DEFECT/FEATURE/TASK tickets, each with a well-defined scope.
- **Don't fragment issues to mirror PRs or commits.** One PR may close several issues, and one issue may span several commits (see the change lifecycle in AGENTS.md). Never invent an extra parent/umbrella issue just to bundle work that will ship together — the PR is the bundle, and the commits are the record. Conversely, don't split work into separate PRs just because it touches separate issues.

### Parent / sub-issue relationships

Parent ↔ sub-issue links are fine, but **only** when the parent carries real intent of its own and the split serves an organizational purpose:

- **The parent is the SSOT of intent for the grouped work.** The parent ticket holds the description and defines the acceptance criteria. Sub-issues exist to break the work into chunks that can be implemented incrementally or in parallel by multiple devs, and to make that progress visible on the kanban board.
- **Sub-issues are organizational proxies and carry no Issue Type.** They exist purely for visibility on the kanban board; the parent (typically a FEATURE or EPIC) remains the source of truth that defines the ACs. Do not duplicate or fragment the acceptance criteria across the sub-issues.
- **If the parent adds nothing, remove it.** When the would-be parent has no intent of its own and every sub-ticket can stand alone, delete the parent and keep the standalone tickets.

### Classify by the actual problem

Choose the Issue Type that matches what the ticket really describes, so the template's sections fit the content:

- A ticket describing a **correctness issue** (something behaving incorrectly) is a **DEFECT**, even if it was first framed as a generic task. As a DEFECT, its Observed/Expected Behavior sections capture the problem cleanly.
- Watch for a misclassified ticket whose "Why?" section is really just describing observed-vs-expected behavior — that is a signal it should be a DEFECT rather than a FEATURE/TASK.

### Other rules

- Do not add manual labels (e.g., `bug`, `priority`) unless the user specifically requests you to do so. Issue Types already apply their own labels automatically.

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
3. **Ask the user to confirm** the resolved version is correct before submitting — the defect may have been observed on an older release.

Default to the latest release. Only use a different version if the user specifies one.

---

## Defect Ticket Template

Whenever you generate a DEFECT ticket, use the following structure exactly. The HTML comments in the template are **placeholders** — replace them with the actual ticket content. Keep the section headings intact:

~~~markdown
### Description

<!-- --------------------------------------------------------------------
Provide a clear summary of the defect. What is wrong, missing, or not
working as expected? Include any relevant context that helps explain the
issue.

Example:

    The createInventoryItem mutation fails with a validation error when
    using valid SKU formats. According to the GraphQL API documentation
    (https://docs.pyck.cloud/graphql-api/mutations/create-inventory-item),
    SKUs should support alphanumeric characters with optional hyphens and
    underscores, but the mutation consistently rejects SKUs containing
    hyphens with "Invalid SKU format" errors.

    This affects our warehouse onboarding process since many clients use
    hyphenated SKU formats for their existing inventory systems.
-------------------------------------------------------------------- -->

### Observed Behavior

<!-- --------------------------------------------------------------------
Describe what actually happens — the incorrect, incomplete,
inconsistent, or missing behavior you're experiencing.

Example:

    When I call the createInventoryItem mutation with a valid SKU, the
    GraphQL API returns a validation error "Invalid SKU format". The
    mutation fails even though the SKU follows the documented format
    requirements. This happens consistently with any SKU containing hyphens
-------------------------------------------------------------------- -->

### Expected Behavior

<!-- --------------------------------------------------------------------
Describe what should happen — the correct behavior according to
specifications, documentation, established patterns in the product, or
reasonable expectations.

Example:

    According to the GraphQL API documentation, the createInventoryItem
    mutation should accept SKUs in alphanumeric format with optional hyphens
    and underscores. The mutation should successfully create an inventory
    item and return the new item's ID when provided with valid input data
    including a properly formatted SKU.
-------------------------------------------------------------------- -->

### Reproduction Steps

<!-- --------------------------------------------------------------------
List the specific steps needed to observe this defect. Clear
reproduction steps help us identify and fix the issue faster.

Example:

    1. Open GraphQL playground or API client
    2. Execute `createInventoryItem` mutation with these variables:
       - name: `Test Product`
       - sku: `WH-ABC-001`
       - repositoryId: `valid-repo-id`
    3. Submit the mutation request
    4. Observe the validation error response
-------------------------------------------------------------------- -->

### Version

<!-- --------------------------------------------------------------------
Provide the version number and commit SHA of the affected release.

To determine the version, use the GitHub API to fetch the latest release:

    gh api repos/{owner}/{repo}/releases/latest --jq '.tag_name + " (" + .target_commitish + ")"'

If target_commitish is a branch name rather than a SHA, resolve it:

    gh api repos/{owner}/{repo}/commits/{branch} --jq '.sha[:8]'

Always confirm the resolved version with the user before submitting
the issue — the defect may exist on an older release.

Example: v0.17.1 (d34db33f)
-------------------------------------------------------------------- -->

### Suggested Solution (Optional)

<!-- --------------------------------------------------------------------
If you have ideas about how to fix this defect, share them here.

Example:

    The createInventoryItem mutation should properly validate SKU format and
    accept hyphens as documented. The validation logic should check for
    alphanumeric characters plus hyphens/underscores instead of rejecting
    valid SKU formats that contain hyphens.
-------------------------------------------------------------------- -->
~~~

---

## Feature Ticket Template

Whenever you generate a FEATURE ticket, use the following structure exactly. The HTML comments in the template are **placeholders** — replace them with the actual ticket content. Keep the section headings intact:

~~~markdown
### What?

<!-- --------------------------------------------------------------------
Describe exactly what functionality should be added or enhanced. What
will users be able to do with this feature?

Example:

    Add a new SerialNumber type that links to InventoryItem with a unique
    serial number string. Update createInventoryItem mutation to optionally
    accept serial numbers. Add new mutations "assignSerialNumber",
    "removeSerialNumber".

    Add serial number filtering to inventory queries to find items by
    specific serial numbers.
-------------------------------------------------------------------- -->

### Why?

<!-- --------------------------------------------------------------------
Explain why this feature is needed and what value it delivers to
end-users. What problem does it solve? What benefit does it provide?

Example:

    Many warehouse customers need to track individual serial numbers for
    high-value items, electronics, or compliance requirements. Currently,
    they can only track quantities at the SKU level, which doesn't meet
    regulatory needs for industries like pharmaceuticals or electronics.

    Serial number tracking would enable customers to trace individual items
    through receiving, storage, picking, and shipping, supporting warranty
    tracking, recalls, and compliance reporting.
-------------------------------------------------------------------- -->

### Acceptance Criteria

<!-- --------------------------------------------------------------------
List the criteria that must be fulfilled for this feature to be
considered complete. Focus on user-facing outcomes and functionality.

Example:

    - [ ] createInventoryItem mutation accepts optional serialNumbers array parameter
    - [ ] assignSerialNumber mutation created and accepts inventoryItemId and serialNumber
    - [ ] removeSerialNumber mutation created and accepts serialNumberId
    - [ ] Serial number uniqueness constraint enforced (returns error for duplicates)
    - [ ] inventoryItems query accepts serialNumber filter parameter
    - [ ] InventoryItem.serialNumbers field returns array of SerialNumber objects
    - [ ] API documentation updated with serial number field descriptions and examples
-------------------------------------------------------------------- -->

### How? (Optional)

<!-- --------------------------------------------------------------------
If you have ideas about how this feature should be implemented or how
users should interact with it, describe them here.

Example:

    Extend the existing GraphQL schema with a SerialNumber type containing
    fields "id", "serialNumber" (unique string), "inventoryItemId", and
    "status". Add serial number array input to createInventoryItem mutation.
    Create new mutations for managing serial numbers post-creation. Update
    the InventoryItem type to include a serialNumbers field that returns
    associated serial numbers. Add serial number filtering to existing
    inventory queries.
-------------------------------------------------------------------- -->

### Impact Assessment (Optional)

<!-- --------------------------------------------------------------------
If known, describe what systems, components, and/or processes will be
affected by this feature.

Example:

    - Affected systems: GraphQL API, Database schema, existing inventory queries
    - Database changes: New serial_numbers table with unique constraints
    - API changes: New SerialNumber type, extended InventoryItem type
    - Dependencies: Existing InventoryItem type, createInventoryItem mutation
    - Breaking changes: None (additive only)
    - Performance considerations: Additional database joins for serial number queries
-------------------------------------------------------------------- -->
~~~

---

## Task Ticket Template

Whenever you generate a TASK ticket, use the following structure exactly. The HTML comments in the template are **placeholders** — replace them with the actual ticket content. Keep the section headings intact:

~~~markdown
### What?

<!-- --------------------------------------------------------------------
Describe exactly what work needs to be done. Be specific about the scope
and deliverables.

Example:

    Document all GraphQL mutations and queries related to inventory
    movements, including createInventoryItemMovement,
    executeInventoryItemMovement, and repositoryMovement operations. Create
    detailed examples showing how to move items between locations, update
    quantities, and track movement history. Include error handling scenarios
    and best practices for warehouse operations.
-------------------------------------------------------------------- -->

### Why?

<!-- --------------------------------------------------------------------
Explain why this task is necessary. What problem does it solve or what
value does it provide to the project or team?

Example:

    Customer onboarding teams are spending 3-4 hours explaining inventory
    movement concepts to new warehouse managers. Existing API documentation
    focuses on individual mutations but doesn't explain the complete
    workflow for moving inventory between locations. Clear documentation
    would reduce onboarding time and help customers implement pyck more
    efficiently in their warehouses.
-------------------------------------------------------------------- -->

### Acceptance Criteria

<!-- --------------------------------------------------------------------
List the criteria that define when this task is complete. What
deliverables or outcomes indicate successful completion?

Example:

    - [ ] Documentation page created at docs.pyck.cloud/api/inventory-movements
    - [ ] All 3 inventory movement mutations documented with parameter descriptions
    - [ ] At least 5 code examples provided covering common scenarios
    - [ ] Workflow diagram included showing movement state transitions
    - [ ] Customer success team has reviewed and approved the content
    - [ ] Documentation is indexed and searchable on docs site
-------------------------------------------------------------------- -->

### How? (Optional)

<!-- --------------------------------------------------------------------
If you have ideas about the approach or methodology for completing this
task, describe them here.

Example:

    Review existing GraphQL schema documentation and identify gaps in
    movement workflow coverage. Interview customer success team about common
    questions. Create step-by-step guides with code examples for typical
    scenarios like receiving inventory, moving between locations, and
    picking for orders. Add diagrams showing the relationship between
    InventoryItem, Repository, and Movement types.
-------------------------------------------------------------------- -->
~~~

---

## Epic Ticket Template

Whenever you generate an EPIC ticket, use the following structure exactly. The HTML comments in the template are **placeholders** — replace them with the actual ticket content. Keep the section headings intact:

~~~markdown
### Description & Context

<!-- --------------------------------------------------------------------
Describe the epic and provide context. What problem does this solve? Why
is this important? Include relevant background information, user needs,
technical considerations, or dependencies.

Example:

    Implement a comprehensive order verification system to reduce picking
    errors in the warehouse. This epic addresses the increasing error rates
    observed in Q1 2024, where picking errors rose from 2.1% to 3.8%.

    The system will include:
    - Visual confirmation steps for pickers
    - Barcode scanning validation
    - Real-time error detection and alerts
    - Analytics dashboard for error tracking

    This work depends on the warehouse scanner hardware upgrade (Epic #123)
    and will integrate with the existing WMS system.
-------------------------------------------------------------------- -->

### Success Criteria

<!-- --------------------------------------------------------------------
Define what "done" looks like for this epic. Include measurable
outcomes, acceptance criteria, or key milestones that indicate
successful completion.

Example:

    - [ ] Visual confirmation workflow deployed to all 5 warehouse zones
    - [ ] Barcode scanner integration tested with 99.9% success rate across 1000 scans
    - [ ] Error alert system delivers notifications within 5 seconds (measured via logging)
    - [ ] Analytics dashboard accessible at /analytics/picking-errors with 6 key metrics
    - [ ] Picking error rate measured below 2.5% for 2 consecutive weeks
    - [ ] System successfully handles 500 concurrent users during load testing
    - [ ] All features documented in user manual and training materials published
-------------------------------------------------------------------- -->
~~~

---

## Initiative Ticket Template

Whenever you generate an INITIATIVE ticket, use the following structure exactly. The HTML comments in the template are **placeholders** — replace them with the actual ticket content. Keep the section headings intact:

~~~markdown
### Strategic Objective

<!-- --------------------------------------------------------------------
Describe the major strategic objective this initiative aims to achieve.
What is the overarching goal? How does this align with organizational
strategy and priorities?

Example:

    Transform our warehouse management platform into a comprehensive supply
    chain optimization solution by expanding from basic inventory tracking
    to intelligent demand forecasting, automated replenishment, and
    predictive analytics.

    This initiative aligns with our 2029-2030 strategy to move upmarket and
    serve enterprise customers requiring end-to-end supply chain visibility
    and optimization. Success will position pyck as a complete supply chain
    platform rather than just a WMS.
-------------------------------------------------------------------- -->

### Scope & Expected Impact

<!-- --------------------------------------------------------------------
Describe the scope of this initiative and the expected business impact.
What will be achieved? What value will this deliver? Include market
opportunities, revenue targets, customer benefits, or competitive
advantages.

Example:

    **Scope:** This initiative covers the development of demand forecasting,
    automated replenishment, predictive analytics, and real-time inventory
    optimization capabilities. Individual epics will be created for each
    major component and linked to this initiative.

    Expected Impact:
    - Enable entry into enterprise segment (500+ employees)
    - Target addressable market expansion from $50M to $500M
    - Project $5M ARR from new enterprise customers in Year 1
    - Increase average contract value by 40%
    - Reduce customer churn by 25% through increased product stickiness
    - Help customers reduce inventory carrying costs by 15-20%
    - Improve customer forecast accuracy from 65% to 85%
    - Enable single platform for end-to-end supply chain operations
-------------------------------------------------------------------- -->

### Success Metrics & KPIs

<!-- --------------------------------------------------------------------
Define the key performance indicators and success metrics for this
initiative. How will you measure whether the strategic objective has
been achieved?

Example:

    - [ ] Platform processes 10M+ inventory transactions per day
    - [ ] Forecast accuracy reaches 85% across pilot customers
    - [ ] At least 5 enterprise customers (>$100K ARR) live in production
    - [ ] Net Revenue Retention reaches 120% in enterprise segment
    - [ ] Customer satisfaction score (CSAT) above 4.5/5 for new features
    - [ ] Platform uptime maintains 99.9% SLA
    - [ ] Time-to-value reduced from 6 months to 3 months for enterprise onboarding
    - [ ] Product recognized in analyst reports as comprehensive supply chain platform
-------------------------------------------------------------------- -->
~~~
