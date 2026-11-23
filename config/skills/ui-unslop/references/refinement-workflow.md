# UI refinement workflow

Use these activities for a page, component, flow or broader frontend effort. Combine, omit or revisit them as the task requires; the numbering organizes the reference rather than enforcing seven sequential phases. There is no required number of commands, documents, interviews or approval rounds. Apply the relevant [anti-slop checks](patterns.md) throughout, even when a formal command is unnecessary.

## 1. Establish context and a useful baseline

Use the task's known audience, content, conventions and constraints; investigate only gaps that affect the design. Distinguish new work, refinement within an existing identity and an authorized redesign. Different surfaces in one product may prioritize action, operation, reading or expression differently.

For existing UI, inspect relevant source and current visual evidence. Check whether screenshots represent the current code, theme and state. For new UI, start with the intended content and task rather than inventing a before-state. Preserve required behavior and data while keeping presentation decisions open within scope.

Consult suitable existing references. Borrow a useful relationship, hierarchy or interaction, not an unrelated product's requirements. Record the quality being pursued, such as clearer action priority or better reading rhythm, rather than "make it like this app".

When using Impeccable, read its installed skill and follow its current setup and selected-command references. Reuse valid session context; do not copy a launcher path from another machine. Its general taste preferences do not replace explicit project decisions. If a command's operational requirements conflict with the current instructions or authorization, name that conflict before invoking it; choose a narrower compatible approach when one meets the task. Do not initialize documents, install tools or alter dependencies solely to satisfy this workflow. If Impeccable was requested but is unavailable, explain that limitation and use the applicable reasoning directly without claiming a command ran.

## 2. Diagnose implementation and design separately

Distinguish two questions:

- Does the interface work correctly across its relevant states and inputs?
- Is its hierarchy, composition, content and interaction appropriate for this task?

Impeccable `audit` addresses technical quality; `critique` addresses design judgment. Select either or both when that depth of assessment is useful or requested. A narrow defect with a settled solution does not need a formal critique. Do not label an informal inspection as a completed Impeccable command.

Apply the anti-slop catalogue explicitly. Separate a symptom from its cause: repeated text may compensate for unclear controls; nested frames may reflect duplicated composition; a cramped field may come from global spacing rather than its own padding. Preserve what already works.

Treat automated flags as evidence to investigate, not commands to delete every matching pattern. Check actual accessibility requirements and project context before accepting detector assumptions. Keep technical findings, visual judgment and unverified hypotheses distinguishable. If a command produces heuristic scores, report their limited meaning rather than presenting them as usability measurements or proof of compliance.

## 3. Settle the brief and important decisions

For substantial work, identify the task, desired outcomes, permitted changes, protected content/behavior, relevant states and how improvement will be judged. Reuse a sufficient existing brief; otherwise capture only what is needed. Use `shape` when formal discovery or a confirmed design brief is useful or requested. Do not reopen decisions already settled.

Ask focused questions when an answer materially changes the result. Explain the recommendation and tradeoff. Visual alternatives can help settle subjective choices about spacing, proportions, color or hierarchy; they are not a prerequisite for every refinement. `generate` or `live` can support this when compatible with the project's browser and review rules. Keep exploratory settings distinct from adopted project defaults.

Respect the user's or project's existing brief and acceptance gates. Do not invent a new approval requirement for every component or minor change. When explicitly invoking `shape` or another command with its own confirmation protocol, follow that protocol rather than claiming it was completed informally.

## 4. Select passes for the demonstrated problems

Choose a sequence proportional to the task, with a reason for each pass. Resolve task structure and dependencies before finishing details. These are selection cues, not fixed recipes:

| Need | Relevant Impeccable commands | Intended outcome |
| --- | --- | --- |
| Unclear task or scope | `shape` | A buildable brief with unresolved decisions exposed |
| Technical defects | `audit` | Evidence-backed implementation findings |
| Weak design direction or hierarchy | `critique` | Prioritized design findings and what works |
| Redundant content and controls | `distill`, `clarify` | Less clutter, clearer labels and consequences |
| Poor grouping, density or composition | `layout` | Relationships and priorities become easier to see |
| Weak reading or scanning | `typeset` | Coherent type roles and readable content |
| Confused surface or state colors | `colorize` | Useful color roles within the project's identity |
| Excessive noise or insufficient emphasis | `quieter` or `bolder` | Appropriate emphasis, chosen for the actual problem |
| Unclear transitions or feedback | `animate` | Motion that explains change without delaying work |
| A useful opportunity for expression | `delight`, `overdrive` | Intentional character appropriate to the brief |
| Errors, long content and edge cases | `harden` | Usable failure, recovery and boundary states |
| First-use or activation problems | `onboard` | A clearer path to the first useful outcome |
| Another device, width or input context | `adapt` | A composition suited to that context |
| Demonstrated interface slowness | `optimize` | A measured improvement to the identified bottleneck |
| Final inconsistencies | `polish` | A coherent finish after structural work |
| Repeated patterns or system drift | `extract`, `document` | Appropriate reuse and durable system guidance |
| A meaningful visual choice | `generate`, `live` | Comparable alternatives with an explicit selection |
| Missing durable product context | `init` | Useful context when the task calls for establishing it |

Read each selected command's installed reference and obey its actual scope. Some commands diagnose or plan without implementing; others have review, confirmation or tool requirements. Select commands for their useful result, not to check off workflow stages or introduce unnecessary ceremonies. Do not treat every name as a CLI verb or invent syntax.

Use `quieter` and `bolder` where the relevant region needs less or more emphasis; neither is a universal direction. Do not add motion to fill a stage or interpret `delight`/`overdrive` as permission to abandon the brief. Likewise, `document` does not mean every cleanup needs new documentation, and `optimize` is not a reason to chase unrelated benchmarks.

## 5. Implement at the correct level

Apply the selected passes using the project's stack and style. When using Impeccable, read its current craft-floor guidance before UI edits, as its setup requires. Preserve the product brief when broader taste guidance conflicts with explicit choices.

If many surfaces share a demonstrated foundation problem, consider a coordinated token or component fix within the authorized scope. Follow the entrypoint's shared-versus-local guidance; discovering repetition does not automatically authorize global changes. Do not impose a palette-first or component-first phase on every task.

For shared changes, inspect consumers and preserve intentional exceptions. Use realistic examples in an existing gallery or local preview when that helps comparison; building a new gallery is not a prerequisite. Keep required content and capabilities accounted for as the composition changes.

Record consequential reasoning briefly in the existing task evidence when useful. The interface itself should not explain why its designer chose a component, color or layout.

## 6. Inspect, fix as a batch and confirm

Use the inspection points and review limits established by the project and the selected command. They may distinguish early exploration/debugging from final verification. Where no sequence is prescribed, batch related inspection and fixes rather than repeatedly reopening the whole design.

Check the agreed outcomes with realistic content and the changed interaction states. Inspect relevant themes, widths, keyboard/focus behavior, loading/error transitions and motion. Include the catalogue's [action-result layout-stability check](patterns.md#unintended-layout-shifts-after-actions), including transient and delayed states. Use actual references at comparable scales, and state what the evidence proves. A screenshot does not prove timing, playback, scrolling or full keyboard support.

Group the findings, fix related causes together and perform a focused confirmation. Avoid an endless sequence of tiny edits and full rechecks. Further work should respond to a new change, failure or unresolved material concern. Do not mark a blocker complete because the review budget ended.

Use independent review when requested, required or proportionate to the risk, within the available authorization, rather than automatically adding duplicate reviewers. Reviewer roles, model choice, browser permissions and service isolation come from the current task/environment, not this reusable skill. If a selected command requires an independent review, perform it as specified or disclose why the formal command cannot be completed; do not silently substitute self-review.

## 7. Close against outcomes

Explain what materially improved, what behavior and content were preserved, what was exercised and what remains partial or unverified. Before/after comparisons help larger redesigns; a focused fix may need only direct evidence of the corrected behavior. Obtain acceptance where the existing agreement requires it.

Do not substitute completed commands, green tests, detector scores or superficial similarity for an effective interface. Carry only approved reusable decisions into the project's design system; do not promote experimental settings or local exceptions into global defaults.
