# Changelog

All notable changes to this module are documented in this file. The format is based on [Common Changelog](https://common-changelog.org/), and this module adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.5.1] - 2026-07-16

### Changed

- Simplify the encoder internals by consolidating attribute writes and error handling.

## [0.5.0] - 2026-07-14

### Changed

- **Breaking:** Rename the `ErrInvalidPlaylist` error type to `InvalidPlaylistError`.
- Encode extra attributes in a deterministic, sorted order.

### Fixed

- Parse channel names that contain commas.
- Preserve fractional track durations when encoding.

## [0.4.0] - 2025-04-04

### Changed

- Improve error handling.

## [0.3.0] - 2025-03-24

### Added

- Add support for `#EXTM3U` attributes.

## [0.2.0] - 2025-03-10

### Fixed

- Recognize a playlist as complete when its final line is a URL without a trailing newline.

## [0.1.0] - 2025-03-10

_:seedling: Initial release._

[0.5.1]: https://github.com/sherif-fanous/m3u/releases/tag/v0.5.1
[0.5.0]: https://github.com/sherif-fanous/m3u/releases/tag/v0.5.0
[0.4.0]: https://github.com/sherif-fanous/m3u/releases/tag/v0.4.0
[0.3.0]: https://github.com/sherif-fanous/m3u/releases/tag/v0.3.0
[0.2.0]: https://github.com/sherif-fanous/m3u/releases/tag/v0.2.0
[0.1.0]: https://github.com/sherif-fanous/m3u/releases/tag/v0.1.0
