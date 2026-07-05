# Examples

## Excellent PR Review

The repository builds handlers directly and avoids extra abstraction layers. This PR follows that pattern, replaces duplicated logic with an existing local helper, and adds a regression test in the same style as nearby tests. The implementation is slightly longer, but it matches the repository's explicit style and does not introduce a new architectural layer.

## Poor PR Review

This PR is rated highly because it "looks clean," but the review never checks whether the repository already has the behavior, whether the abstraction matches adjacent modules, or whether the change is necessary. The result is a generic approval that ignores maintainers' conventions and risk.

## Do Not Recommend a Violating Change

The repository consistently constructs dependencies directly inside request handlers. Recommending dependency injection here would create a new style in one module while the surrounding code remains direct and simple. The better review is to say the change should stay aligned with the existing handler pattern unless the repository starts a broader refactor.

## Example Scoring Report

| Category | Score | Reason |
| --- | ---: | --- |
| Correctness | 14/15 | Behavior matches the existing flow and edge cases are covered. |
| Necessity | 9/10 | The change removes duplication that already exists in this repository. |
| Architecture | 8/10 | Fits the current module layout without creating a new layer. |
| Maintainability | 8/10 | Reduces future edits in the touched area. |
| Testing | 3/3 | Adds the same kind of regression coverage the repository already uses. |
| Repository Consistency | 10/10 | Matches the local style exactly. |
| Regression Risk | 9/10 | Low, because the change stays within the existing pattern. |

Total: 61/68 scaled to 90/100.

## Example Final Verdict

Approve with one required change: keep the implementation direct, because the repository consistently avoids extra abstraction in this area. After that adjustment, the PR matches the codebase and is likely to be accepted.