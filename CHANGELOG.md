# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [v3.3.0] - 2021-07-21

### Added

- ci: Add gofmt action (#50)
- ci: Add diff check action (#53)
- ci: Add dependabot auto build support (#54)

### Fixed

- ci: Fix auto-build not work correctly
- storage: Fix copy and move behavior (#57)

## [v3.2.0] - 2021-06-29

### Changed

- *: Implement GSP-109 Redesign Features (#48)
- *: Implement GSP-117 Rename Service to System as the Opposite to Global (#48)

## [v3.1.0] - 2021-06-11

### Added

- *: Implement GSP-87 Feature Gates (#44)
- storage: Create dir (#45)

## [v3.0.0] - 2021-05-24

### Added

- storage: Implement GSP-49 Add CreateDir Operation (#39)
- *: Implement GSP-47 & GSP-51 (#40)
- storage: Implement GSP-61 Add object mode check for operations (#41)

### Changed

- storage: Idempotent storager delete operation (#38)
- *: Implement GSP-73 Organization rename (#42)

## [v2.1.0] - 2021-04-24

### Added

- storage: Implement proposal unify obejct metadata (#29)
- *: Implement default pair support for service (#30)
- storage: Add Mkdir support (#31)
- storage: Implement Create API (#32)
- *: Add UnimplementedStub (#33)
- storage: Implement Appender support (#34)
- tests: Introduce STORAGE_FS_INTEGRATION_TEST (#35)

### Changed

- ci: Only run Integration Test while push to master

## [v2.0.0] - 2021-01-21

### Added

- storage: Implement Fetcher (#26)

### Changed

- Migrate to go-storage v3 (#27)

## v1.0.0 - 2020-11-12

### Added

- Implement fs services.

[v3.3.0]: https://github.com/beyondstorage/go-service-fs/compare/v3.2.0...v3.3.0
[v3.2.0]: https://github.com/beyondstorage/go-service-fs/compare/v3.1.0...v3.2.0
[v3.1.0]: https://github.com/beyondstorage/go-service-fs/compare/v3.0.0...v3.1.0
[v3.0.0]: https://github.com/beyondstorage/go-service-fs/compare/v2.1.0...v3.0.0
[v2.1.0]: https://github.com/beyondstorage/go-service-fs/compare/v2.0.0...v2.1.0
[v2.0.0]: https://github.com/beyondstorage/go-service-fs/compare/v1.0.0...v2.0.0
