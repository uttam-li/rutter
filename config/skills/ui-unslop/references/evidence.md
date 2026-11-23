# Research and interpretation

## How to use the evidence

Primary guidance reviewed on 3 November 2026. These sources support the narrow principles below. They do not define "AI slop", endorse every catalogue suggestion or supply universal styling tokens. The catalogue combines these principles with design judgment; distinguish that synthesis from a source's actual requirements.

Recheck a source when exact current guidance matters. No browser research is required on every ordinary cleanup. Prefer an existing relevant reference over repeatedly searching for aesthetic authority.

## Primary guidance

| Source | Supported principle | Boundary |
| --- | --- | --- |
| [NN/g, Aesthetic and minimalist design](https://www.nngroup.com/articles/aesthetic-minimalist-design/) | Unnecessary information competes with relevant information. Minimalism concerns usefulness, not an obligatory visual style. | This does not require removing useful detail, imagery or brand expression. |
| [GOV.UK, Button](https://design-system.service.gov.uk/components/button/) | Labels should name the action in sentence case; competing primary actions dilute hierarchy. | Additional words can clarify a consequence. Its government-service layout is not a universal product template. |
| [GOV.UK, Text input](https://design-system.service.gov.uk/components/text-input/) | Keep labels; placeholders cannot replace them. Field widths can reflect expected input, and hints should be concise and useful. | Do not import its exact dimensions or all label-placement decisions into another system. |
| [GOV.UK, Details](https://design-system.service.gov.uk/components/details/) | Optional detailed help can be disclosed to improve scanning. | Do not hide information most people need to complete the task. |
| [GOV.UK, Error message](https://design-system.service.gov.uk/components/error-message/) | Validation should explain the problem and how to correct it, in connection with the affected input. | Do not confuse service failure with invalid input. |
| [IBM Carbon, Empty states](https://www.carbondesignsystem.com/building-blocks/core/patterns/empty-states) | Empty states need context and useful next steps; imagery and actions depend on the situation. | A large illustration or identical alignment is not required for every empty region. |
| [IBM Carbon, Notification guidelines](https://www.carbondesignsystem.com/building-blocks/core/components/notification/guidelines) | Notifications should be relevant and used sparingly; variant and status serve different purposes. | This does not support replacing every inline failure with a transient toast. |
| [WAI-ARIA APG, Tabs](https://www.w3.org/WAI/ARIA/apg/patterns/tabs/) | Tabs represent associated panels with defined roles and keyboard behavior. Automatic activation should avoid noticeable loading latency. | APG is implementation guidance, not a mandate for segmented or underline styling. |
| [W3C, Target size minimum](https://www.w3.org/WAI/WCAG22/Understanding/target-size-minimum.html) | WCAG 2.2 AA specifies 24 by 24 CSS pixel targets, with defined exceptions including spacing. | Visible face size, pointer area and spacing differ. This is not a universal 44px control-height requirement. |
| [W3C, Use of color](https://www.w3.org/WAI/WCAG22/Understanding/use-of-color.html) | Information and state cannot depend on color alone. | Additional cues do not remove separate contrast obligations. |
| [W3C, Focus not obscured minimum](https://www.w3.org/WAI/WCAG22/Understanding/focus-not-obscured-minimum.html) | At AA, focused components must not be entirely obscured by authored content. | A clipped ring still merits a usability fix; do not misstate this criterion as requiring every ring pixel to be visible. |
| [W3C, Animation from interactions](https://www.w3.org/WAI/WCAG22/Understanding/animation-from-interactions.html) | Nonessential interaction-triggered motion can be disabled under this AAA criterion. | Respecting reduced motion is useful; do not label this particular criterion AA or forbid all animation. |
| [W3C, Reflow](https://www.w3.org/WAI/WCAG22/Understanding/reflow.html) | Content should remain usable at the specified narrow dimensions without two-dimensional scrolling, with exceptions for content that requires it. | Tables, maps and similar two-dimensional regions can need local scrolling; the exception is not a blanket exemption for an entire page. |
| [W3C, Content on hover or focus](https://www.w3.org/WAI/WCAG22/Understanding/content-on-hover-or-focus.html) | Additional hover/focus content must meet dismissibility, hoverability and persistence conditions where applicable. | This does not forbid tooltips or imply every action should be permanently visible. |
| [W3C, Text spacing](https://www.w3.org/WAI/WCAG22/Understanding/text-spacing.html) | Supported markup must tolerate the specified user text-spacing overrides without losing content or functionality. | Those override values are not mandatory default typography tokens. |
| [W3C, Label in name](https://www.w3.org/WAI/WCAG22/Understanding/label-in-name.html) | A control's accessible name contains its visible text label. | Shortening visible copy must not leave a conflicting programmatic name. |
| [USWDS, Data visualizations](https://designsystem.digital.gov/components/data-visualizations/) | Use understandable representations, textual labels and access to underlying data; do not make comprehension depend on mouse interaction. | The source's examples are not a required chart library or dashboard style. |

## Boundaries for applying the guidance

These are web guidance sources. Apply the relevant semantics and accessibility obligations for the target platform; do not copy a web implementation contract into a different platform without checking it.

Keep project-specific decisions in that project's own design or product guidance. Empty-state alignment, notification appearance, destructive emphasis, tab styling, scrollbar visibility, control sizes, typography and animation timing are contextual choices. This skill supplies no preferred values or visual preset.

A design system can express these principles through different libraries, native elements or custom components. Judge the resulting reading, interaction and state behavior rather than its package names. A distinctive style is worth keeping when it supports the project and remains usable.
