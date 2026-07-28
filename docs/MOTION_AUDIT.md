# Motion Audit: Re-writing-History (Git History Telemetry)

## 1. Existing Motion Grep Analysis
* **Grep Hits**: 0 files reference animation or transition properties.
* **Missing Animations**: Hero banner entrance, commit history timeline card staggered reveal, commit SHA hover scale, live engine status pulse indicator.

## 2. High-Value Targeted Additions
* **Library / Strategy**: Pure CSS `@keyframes` + Vanilla JS IntersectionObserver (no external bundle size overhead).
* **Target Interactions**:
  1. Hero section entrance reveal (`@keyframes slideFadeIn`).
  2. Commit log timeline cards staggered entrance (`animation-delay: calc(var(--i) * 50ms)`).
  3. Commit SHA pill scale & highlight hover state (`transition: transform 0.15s ease, box-shadow 0.15s ease`).
  4. Active engine status light pulse (`@keyframes enginePulse`).
* **Accessibility**: Respect `prefers-reduced-motion: reduce`.
