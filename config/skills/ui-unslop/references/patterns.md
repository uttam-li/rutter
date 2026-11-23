# Interface pattern catalogue

Use the relevant sections, not every entry on every screen. Each pattern describes a possible defect and a remedy, not a forbidden aesthetic. Preserve a treatment when it serves the audience, task or brand. Being a library default or appearing elsewhere in the project is not, by itself, a reason to keep redundant copy, unclear metadata chains or competing decoration. Resolve those problems within the project's visual language.

## Copy and metadata

### Dot- and dash-separated metadata chains

**Signal:** `Invoice 1042 · Acme — Overdue · $240 · 12 May` treats identity, status, amount and date as one undifferentiated string.

**Repair:** Prefer meaningful structure over automatically joining values with dots or dashes. Give those facts stable positions using a phrase, spaced values, label/value pairs or columns. Use labels where meaning is ambiguous, not on every obvious value. Retain an intentional compact line when its relationships are already clear, not merely because it uses fewer characters. Preserve ranges, minus signs, identifiers, decimals and intentional missing-value markers; never replace punctuation globally.

### Narration and repeated orientation

**Signal:** "Use the search box to search" repeats a visible label; breadcrumb, heading, subtitle and empty-state body repeat the same information. Copy explains why the interface is modern instead of what users need to know.

**Repair:** Give each text element a distinct job. Remove repetition and implementation justification. Improve unclear controls rather than explaining them indefinitely. Retain navigation context, meaningful headings, unfamiliar requirements, format help and consequential service limitations.

### Labels that are long, vague or inconsistent

**Signal:** "Click here to add a new item to your collection" is verbose; "Proceed" hides a commitment; the same operation is called Save, Apply and Submit on different screens.

**Repair:** Use the shortest label that accurately distinguishes the action and target. Retain necessary consequence words, and use consistent terms for equivalent outcomes. Keep visible labels in accessible names. Do not impose a word-count limit or collapse genuinely different operations into one label.

### Reassurance or outcomes without evidence

**Signal:** "100% secure", invented testimonials, arbitrary improvement percentages or "Sent" after creating only a local preview.

**Repair:** Use verified facts and outcomes. Distinguish draft, accepted, pending and completed where the system does. Label demonstration content appropriately. Do not fabricate evidence to fill a layout or make an incomplete feature sound finished.

### Every value becomes a pill

**Signal:** Ordinary identifiers and amounts appear in code-style capsules; all categories and counts get saturated badges.

**Repair:** Use plain selectable text for ordinary values, code styling for technical content and badges when category or status benefits from recognition. Preserve copy/link behavior. A consistent badge system can remain useful; the issue is indiscriminate emphasis, not the shape itself.

### Placeholder and tooltip dependence

**Signal:** A field loses its only label while typing; mandatory format rules or consequences exist only in a tooltip; helper text is repeated on every field regardless of need.

**Repair:** Keep persistent labels and place necessary instructions where they can be read during the task. Use optional help for occasional detail. Preserve programmatic associations without making every hint a long screen-reader announcement. Hover/focus behavior is addressed under Interaction and responsive behavior.

## Composition and task structure

### Framing without a separate responsibility

**Signal:** A page card contains a collection card, record cards, framed media and boxed subsections, with redundant borders and dividers between already distinct groups.

**Repair:** Remove the layer that contributes no useful grouping or interaction boundary. Use proximity and alignment where sufficient. Keep independent record cards, meaningful section boundaries and overlay elevation when they clarify structure. There is no universal ban on nested containers.

### A template determines the page instead of its task

**Signal:** An oversized hero delays an ordinary task; unrelated information is forced into equal bento tiles; every section repeats eyebrow, heading, paragraph and three cards.

**Repair:** Let the content and user journey determine the structure. Borrow useful relationships from references instead of copying their whole composition. Marketing pages may need expressive heroes; dashboards may need grids. Change them only when they obscure priority or create repetition, not because they are recognizable patterns.

### Spacing ignores content and context

**Signal:** Short numeric fields span a monitor, unrelated groups look connected, or every control inherits the same large gap. Conversely, labels and targets collide in the pursuit of density.

**Repair:** Size for expected content and distinguish within-group spacing from between-group spacing. Keep room for errors, localization and interaction. Consider inline labels or more columns when usable, not as an automatic compactness rule. Do not force all forms into one spacing scale or all mobile content into cards.

### Everything competes for first place

**Signal:** Every metric is headline-sized; all cards use equal emphasis despite unequal importance; a secondary status competes with the main task.

**Repair:** Establish priority through order, grouping, size, weight and contrast. Comparable figures often benefit from a common size and aligned numerals, but a leading metric may deserve emphasis. Preserve meaningful heading structure and the product's typography system.

### Minimalism hides useful content

**Signal:** A one-item menu conceals a frequent command; every document clause requires expansion; useful secondary content becomes difficult to find.

**Repair:** Keep common actions and decision-critical information discoverable. Use disclosure for optional depth, with clear labels. Secondary can mean lower visual emphasis rather than hidden. Account for sensitive data, limited space and novice/expert differences before exposing or removing controls.

### Steps and confirmations become ceremony

**Signal:** Ordinary reversible actions require repeated confirmation; onboarding blocks the first useful task; an editor forces an unrelated wizard structure.

**Repair:** Remove steps that add no decision, protection or necessary information. Consider contextual guidance or reversible actions where supported. Retain required consent, destructive safeguards, financial confirmation and genuine dependencies. Do not replace a required confirmation with Undo unless the operation is actually reversible.

## Visual hierarchy and character

### Typography sacrifices reading to appearance

**Signal:** Tiny muted labels, excessive uppercase or letter spacing, headings broken into arbitrary lines, or ellipsis hiding values users must compare.

**Repair:** Choose type roles for the reading task, with usable contrast and line length. Remove forced line breaks where responsive wrapping works better. Keep complete information available when truncation is necessary. Support text enlargement and spacing overrides; do not assume one font, size or weight suits every product.

### Neutrality erases boundaries and states

**Signal:** Canvas, sidebar, fields, selected content and disabled controls occupy nearly the same tonal range.

**Repair:** Define distinguishable surfaces, text and interaction states in each supported theme. Separate hover, selected, focus and disabled treatments. Borders, fills and elevation are available tools; avoid either flattening everything or compensating with random bright accents.

### Color has no consistent meaning

**Signal:** Arbitrary pastels decorate each statistic, one color means conflicting states, or dark mode is a mechanical inversion.

**Repair:** Use semantic and brand roles with suitable light/dark values. Preserve established domain meanings. Pair meaningful color with text, icons or other cues and verify actual contrast combinations. Do not make selection and a primary command indistinguishable merely because both use the brand color.

### Decoration and icons compete with content

**Signal:** Glass, gradients, glows, floating shapes, hover lifts and decorative icons appear everywhere; inconsistent strokes make related actions look unrelated.

**Repair:** Keep effects that serve identity, affordance or a meaningful moment, and reduce competing layers. Use consistent icon geometry and labels for uncertain meanings. Do not replace the design with mandatory grayscale, remove all personality or ban a font/effect as inherently AI-generated.

## Actions and controls

### Button emphasis and variants drift

**Signal:** Primary commands and utilities all use strong fills, or the same action gets a different treatment on every page.

**Repair:** Establish priority within each task and reuse variants for equivalent purposes. Preserve meaningful destructive and page-specific action distinctions. A low-emphasis Back link or secondary utility may fit, but no particular color or variant is universal. Avoid a new style for every action.

### Similar-looking controls do different jobs

**Signal:** Commands appear as tabs, navigation looks disabled because it uses button styling poorly, or a single-choice value is implemented as unrelated buttons with no selected state.

**Repair:** Distinguish destinations, commands, associated content panels, value choices and disclosure. Use appropriate semantics and keyboard behavior, including when an action is styled like a text link. A segmented appearance alone does not determine whether tabs, radio choices or toggle buttons are appropriate.

### Navigation, status and display modes compete

**Signal:** Main sections, record status and presentation preferences all appear as equal tab strips; every available operation is added to one growing tab list.

**Repair:** Establish a primary navigation level, then subordinate filters, display controls and commands. Group related variants near their content when useful. Choose overflow, scrolling or another composition based on discoverability and input access; do not require underline tabs, segmented tabs or hidden scrollbars.

### Choice controls ignore scale or value precision

**Signal:** Hundreds of options require scrolling without search; selected chips expand a filter row indefinitely; an exact quantity relies only on dragging.

**Repair:** Use searchable choices, stable summaries, visible selections or synchronized value entry according to the task. Keep full selection access and relevant bulk controls. Preserve distinctions between no filter, all values and no values. Do not replace continuously compared values with a count that hides needed information.

### Visual consistency changes commit or action scope

**Signal:** A new switch silently introduces autosave; Load executes a saved query; a duplicated summary implies a different selection scope; a redesigned bulk action affects hidden records unexpectedly.

**Repair:** Preserve or deliberately specify edit, apply, run and save behavior. Make the action target and commit point clear. Remove duplicate controls only when their responsibilities match. Distinct group totals or allocation counts may be necessary. A styling task does not authorize new validation or eligibility rules.

## Feedback and layout stability

### Feedback is duplicated or placed in the wrong channel

**Signal:** A permanent banner repeats ordinary page information; one operation emits several notifications; field errors live only in a toast; partial-failure details disappear with it.

**Repair:** Use local feedback for local problems, brief notifications for outcomes users might miss and persistent recovery for unresolved failures. Update one operation's notification where appropriate. Keep genuine service-wide alerts or required warnings prominent. Avoid duplicating visible and announced feedback unnecessarily; do not ban all banners or prefer toasts for every error.

### One state stands in for several different outcomes

**Signal:** "No data" means idle, loading, no matches and service failure; a spinner blocks the whole screen during one local action.

**Repair:** Represent the actual outcome with concise context and a useful next step. Use skeletons when the layout is predictable and local progress when an action is scoped. Imagery can help a large state region but is not mandatory everywhere. Choose centered, inline or other alignment to fit the container and product.

### Validation moves controls or offers no way forward

**Signal:** Error text enlarges a center-aligned field container and moves its input; untouched fields start red; a disabled Submit gives no explanation of what is missing.

**Repair:** Anchor controls to their own row and let associated messages flow without clipping. Choose validation timing appropriate to the input and consequence. Make errors discoverable and correctable without relying on a toast. Preserve legitimate pending, permission and eligibility restrictions; do not universally require an enabled Submit or one error placement.

### Unintended layout shifts after actions

**Signal:** After an action or its asynchronous result, the acted-on element, neighboring controls or content below unexpectedly shifts. Examples include a longer button label moving its neighbors, a spinner changing button width, a checkmark changing row height, validation moving the input itself, or a temporary status line pushing results down and pulling them back up.

**Repair:** Keep routine idle, pending, success, failure and reset states geometrically compatible. Use suitable shared label/icon slots, stable media dimensions, anchored control rows and local feedback regions. Reserve space only where it serves a predictable state; scope it to the relevant control or region. Preserve readable errors, full content and longer translated labels. Do not use blanket fixed heights, clipping, absolute-positioned overlaps or large empty areas to conceal the problem. Animating an accidental shift does not make it intentional. Coordinate related state painting too, so pinned and ordinary cells do not highlight at different times.

**Intentional movement:** Expansion/collapse, explicit sorting or reordering, inserted/deleted records, navigation and responsive reflow can legitimately move content. Keep that movement tied to the intended change, with surrounding controls and the user's reading/focus position preserved where appropriate. Do not infer that every incidental shift is acceptable merely because a user action preceded it.

**Verify:** Compare relevant element positions before, during and after the action, including delayed completion and failure. Check the changed labels, indicators and feedback at representative widths. Watch the transition itself; stable initial and final screenshots can miss a temporary jump.

### Overlays waste space or obscure the task

**Signal:** A tiny form occupies a huge blank panel, nested scrolling hides actions, a popover cuts its own focus ring, or several identical dismissal controls crowd the footer.

**Repair:** Size the overlay for its purpose and manage scrolling, insets, focus entry/return and dismissal coherently. A long overlay may keep its title/actions visible; a side inspector may legitimately fill the height. Remove redundant dismissal only when it has no distinct purpose. Preserve unsaved edits and required modal boundaries.

### Progress and attachment UI show ceremony instead of work

**Signal:** An indefinite spinner or fake percentage supplies no useful status; an oversized upload zone lacks file status, rejection reasons or Remove.

**Repair:** Show real progress only when measurable, otherwise use honest pending feedback. Support file selection by keyboard and actionable file rows. Choose a compact attachment area or larger drop zone by task. Preserve restrictions and distinguish local acceptance from completed transfer; do not invent previews or cancellation support.

## Interaction and responsive behavior

### Important interactions exist only on hover

**Signal:** Essential actions appear only under a mouse; hover cards vanish when approached; drag is the only way to reorder or choose a value.

**Repair:** Make the function reachable through relevant keyboard and touch paths, with visible or discoverable alternatives. Ensure hover/focus content remains usable and dismissible as applicable. Optional hover shortcuts can remain. Do not turn a tooltip into the only accessible name or place an entire workflow inside one.

### A static screenshot is mistaken for responsiveness

**Signal:** Fixed widths, absolute offsets or forced text breaks work at one size but collide at narrow widths, zoom, larger text or translated content.

**Repair:** Check content-based resizing, wrapping and sensible breakpoints. Account for text direction and longer labels where supported. Preserve two-dimensional tasks such as tables, maps and editors with appropriate local navigation/scrolling rather than squeezing away their meaning. Do not mandate one mobile composition for every frontend.

### Scroll behavior conceals or traps content

**Signal:** A page and several small inner panels scroll independently without clear ownership; hidden overflow cuts options; a fixed header covers anchor targets; the main canvas loses useful space to permanent panels.

**Repair:** Make scroll ownership and overflow discoverable. Contain scrolling where it supports the task, preserve target visibility and keep important tools reachable. Choose scrollbar treatment by platform and context. Do not force either a single page scroller or invisible scrollbars universally.

### Focus and reading position are discarded

**Signal:** A popup returns focus to the page start, filtering steals focus from typing, a refresh loses selection, or new messages pull readers away from earlier content.

**Repair:** Preserve focus and working position when the task remains continuous; deliberately move focus when a new context or error requires it. Clean up stale async changes. Follow new content automatically only when appropriate. Do not persist private state across accounts or promise draft recovery the system does not implement.

## Data, media and motion

### Tables trade useful density for accidental inconsistency

**Signal:** Identifiers and amounts wrap unpredictably, comments inflate every row, columns waste space, or read-only values resemble editable controls.

**Repair:** Set alignment, wrapping and column sizing by purpose and the existing contract. Retain accessible full values, sorting, resizing and selection behavior. Flexible columns can absorb space where suitable; fixed/user-sized columns may stay fixed. A compact table and a table with wrapped descriptions can both be valid.

### Charts and metric tiles are decorative evidence

**Signal:** Generic sparklines, gauges or progress rings lack labels, units, time periods, a denominator or genuine source data. Every number becomes a tile although a sentence or table would answer the question better.

**Repair:** Start with the comparison or decision the data should support. Use an understandable representation with source-backed values and accessible labels or alternatives. Preserve scale, uncertainty and missing-data meaning. Do not invent trends, percentages or chart data merely to make the screen look complete.

### Media states jump or claim unsupported success

**Signal:** Missing media changes frame height, empty slides look available, meaningful imagery is tinted/cropped for decoration, or an iframe load is treated as proof of video playback.

**Repair:** Give available, loading and failed media coherent geometry and truthful controls. Preserve aspect, subject identity and provider restrictions. Offer recovery when supported. Separate primary media from its variants and keep viewer chrome proportionate. Decorative crops can be intentional; inspection, evidence and documents need fidelity.

### Motion obstructs ordinary interaction

**Signal:** Repeated entrances delay content, hover effects fight focus, selection fills flash, or accordions stutter during reversal.

**Repair:** Identify the actual property/state owner and remove competing effects. Use suitable component behavior, measured geometry or a simpler transition rather than arbitrary timing patches. Check opening, closing, rapid interaction and reduced motion. Preserve intentional expression when it serves the product and does not delay essential work.
