---
name: ui-unslop
description: Guide new frontend design or review existing interfaces to prevent redundant copy, generic composition, decorative clutter, weak hierarchy and unstable interactions. Use for UI anti-slop guidance, audits and cleanup across frontend stacks while respecting each project's audience, style, conventions and functionality. Not for prose-only editing.
---

# UI unslop

Make the interface easier to understand and use without turning it into another generic template. Treat "AI slop" as shorthand for observable problems, not evidence of authorship. A dot, card, gradient, font or animation is not a defect by itself. Identify the cost to the user before changing it.

## Work within the request

For a new interface, use these principles while choosing its structure, controls and states. An audit-only request produces findings; when fixes are also authorized, make the scoped improvements. Scale the work to the scope; a control fix does not require redesigning the entire page.

Identify the audience, task, information density, current design system and applicable constraints. A portfolio, checkout, reading page and professional editor need different treatments. Consult relevant existing references; research a gap when it would affect the decision.

Stay compatible with the chosen stack and platform. Use the project's components, tokens and conventions when available; otherwise use suitable platform primitives and establish coherent choices for the task. No framework, library, component preset or CSS method is required. Do not replace dependencies or import a favorite design system just to follow this skill.

Separate product requirements from visual choices and current implementation. Existing behavior is evidence, not automatically a requirement. Consult the project's designated sources when a behavior decision depends on them; migration history is relevant only when the task involves it. Reference images suggest design approaches, not product requirements or exact CSS tokens.

Do not carry defaults from another project into this one. Typography, density, palette, radii, button emphasis, tab treatments, state alignment and motion belong to the current context. Respect them when established; proposing better defaults is appropriate when that scope is open. Ask about unresolved consequential decisions, not every routine styling choice. Where a convention conflicts with a functional or accessibility requirement, identify the specific conflict and resolve it within the authorized scope rather than silently preserving or replacing it.

## Find the actual problems

Read the relevant sections of [the pattern catalogue](references/patterns.md). It covers copy, composition, visual hierarchy, controls, feedback, responsive behavior, data and media. Use only applicable patterns; it is not a checklist every screen must fail.

Actively inspect for recurring agent-output shortcuts: dot/dash-separated metadata chains, redundant instructions, repeated justification, sentence-length buttons, indiscriminate pills, nested card frames, formulaic hero/bento layouts, invented metrics, decorative icon overload and gratuitous motion. Prevent these defaults during generation and flag them during review. Adapting to a project's style must not dilute these concrete checks into vague advice to "make it cleaner"; fix the pattern using that project's own visual language.

For each material finding, identify its location, the observed problem, its effect on the user and an appropriate remedy. Distinguish visible defects from hypotheses and preferences. Prioritize lost functionality and unclear consequences, then task structure and hierarchy, then repetitive detail. A few cosmetic edits do not resolve a structural problem.

When building something new, apply the same reasoning to proposed choices rather than inventing findings about an interface that does not exist. Let real content determine hierarchy and component choice. If project conventions already work, continue them. Challenge a convention only for a concrete problem, explaining the tradeoff rather than substituting personal taste.

Use [the research reference](references/evidence.md) for source-backed principles and their limits. The catalogue also contains design synthesis. Do not present every recommendation as a research finding or assign unsupported design scores.

## Use a deliberate refinement workflow

For work needing several design decisions, use [the refinement workflow](references/refinement-workflow.md) to select useful activities, not to impose a fixed sequence. It covers Impeccable command selection, visual alternatives and bounded review. When Impeccable is requested or its commands fit the task, follow its installed instructions, within the current user's scope and project constraints. Do not run every command, add mandatory dependencies or treat its general aesthetic guidance as authority over the project's visual identity. Small, settled fixes can use the relevant parts directly.

The workflow complements the concrete anti-slop checks above. Command execution, a polished screenshot or a generic quality score does not replace resolving the actual problems.

## Replace clutter with structure

Let alignment, proximity, useful defaults, concise labels and clear states explain ordinary interactions. Remove instructions that merely repeat an obvious control. Keep help that prevents a plausible mistake, explains an unfamiliar concept or communicates a meaningful consequence.

Avoid using dot/dash chains as the default way to assemble unrelated metadata. Use a natural phrase, independently spaced values, aligned label/value pairs, a description list or columns. Retain an intentional compact metadata line when its meaning and grouping are already clear. Do not replace every separator with a badge, rule or label. Preserve real ranges, negative signs, decimals, identifiers, legal punctuation and meaningful missing-value markers. Never apply a blanket punctuation replacement.

Choose components by their job. Commands, destinations, content panels, value choices and disclosure need appropriate semantics even when they look similar. Use existing suitable components and variants; a cleanup does not require expanding the component library. The catalogue's examples describe behavior and relationships, not required component names or APIs.

Fix repeated problems in shared styles or components when the meanings match and the scope permits. Check affected callers. Prefer a focused local fix when a shared change would introduce unrelated risk. Avoid both duplicated overrides and unnecessary abstraction.

## Preserve meaning and control

Before removing an element, distinguish redundancy from unique content or capability. Preserve a useful home for fields, units, statuses, actions and consequences unless their removal is explicitly intended. For larger changes, use a compact before/after accounting map; small fixes need no extra documentation ceremony.

Keep accessible names, semantic relationships, keyboard behavior, visible focus and non-color cues. Preserve full-value access, precision, permissions, confirmations and draft protection where applicable. A shorter label must still identify its action. A compact control must remain operable.

Model loading, empty, filtered-empty, invalid, unavailable and failed states accurately. Preserve usable content during a recoverable refresh failure when it remains safe to show. Selected, previewed, submitted and completed are different outcomes; do not imply service success without evidence.

Keep the acted-on element, its neighbors and content below anchored through routine action and result states unless repositioning is an intended part of the interaction. Treat unexpected movement from pending indicators, changed labels, validation, selection or asynchronous feedback as a defect. Follow the catalogue's [layout-stability guidance](references/patterns.md#unintended-layout-shifts-after-actions); do not preserve position by clipping content or reserving excessive empty space.

Use color and motion purposefully, preserve media fidelity and respect reduced motion. Do not achieve minimalism by hiding useful content, suppressing errors or making everything faint.

## Check the result honestly

Exercise the changed states and relevant content, themes, widths and inputs. Use authorized browser verification where available; source inspection and unit tests cannot prove visual alignment, responsive layout or animation smoothness. Run the project's required checks. The [workflow's review and closeout stages](references/refinement-workflow.md#6-inspect-fix-as-a-batch-and-confirm) explain how to batch this work without repeated polishing loops.

Report material improvements, preserved behavior, evidence and remaining gaps. Keep unverified claims distinct from passes. This skill adds no tool permissions or approval gates, and does not authorize unrelated redesigns, service changes or publication.
