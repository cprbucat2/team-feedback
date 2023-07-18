# TeamFeedback Contribution Guidelines
Hi, and thank you for taking the time and energy to contribute to TeamFeedback. We have a few guidelines for you to follow,
as well as a few tips to help you save your valuable time and ours.

- [Code of Conduct](#code-of-conduct)
- [PR titles](#pr-titles)
- [Code style](#code-style)

## Code of Conduct
Help keep our community open, inclusive, and healthy. Read the [Code of Conduct](CODE_OF_CONDUCT.md).

## PR titles
Consistent PR titles help us manage features and enforce [Semantic Versioning](https://semver.org).
Please use the following PR title format:
```
[<tag>] <short description>
```
Make sure the total length **including the tag** is 80 characters or less.
In the case of reversion, titles might get a little funky. Move the "Reverts" after the tag in the PR title.

### `<tag>`
Must be one of the following:
- Breaking: Changes to the public API that are backward incompatible. (Note: We are still pre-1.0.0 so no changes are breaking.)
- Feature: Changes to the public API that are backward compatible.
- Fix: Bug fixes.
- Refactor: Changes to code that do not fix bugs or add features.
- Docs: Changes to documentation only.
- Dev: Changes to build system, CI, or other developer only items.

## Code style
We use a few different languages and our goal is to use a uniform style. Code style is currently enforced in the following ways:
- Go: `gofmt`, `go vet`, and `golangci-lint` run in our GitHub Actions CI.
