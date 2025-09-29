# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - TBD

### Added
- Initial release of the ICS Terraform Provider
- Support for managing SSH keys (`ics_ssh_key` resource)
- Support for managing bare metal servers (`ics_bare_metal_server` resource)
- Provider configuration with API token authentication
- Complete documentation for all resources
- Automated testing and release workflows

### Features
- **SSH Key Management**: Create, read, update, and delete SSH keys
- **Bare Metal Server Management**: Provision and manage bare metal servers
- **Multi-platform Support**: Windows, macOS, Linux (x86_64, ARM64, etc.)
- **Comprehensive Documentation**: Full resource documentation with examples

[Unreleased]: https://github.com/UK2Group/terraform-provider-ics/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/UK2Group/terraform-provider-ics/releases/tag/v1.0.0