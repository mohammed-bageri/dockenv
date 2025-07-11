# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2025-01-09

### Added

- **Service Management**: `dockenv add` and `dockenv remove` commands for dynamic service management
- **Auto-start Support**: `dockenv autostart` commands for systemd integration
- **Enhanced Profiles**: Additional profiles for common development stacks
- **Environment File Generation**: Automatic `.env` file creation with service credentials
- **Project Detection**: Auto-detect Laravel, Node.js, Django, Rails, and Spring projects
- **Extended Service Support**: Added Elasticsearch and RabbitMQ
- **Custom Port Configuration**: Support for custom ports during initialization and service addition
- **Comprehensive Logging**: `dockenv logs` command with follow and filtering options
- **Enhanced Status Display**: Improved `dockenv status` with detailed service information
- **List Command**: `dockenv list` to show available services, profiles, and current configuration

### Enhanced

- **Interactive Setup**: Improved `dockenv init` with better prompts and validation
- **Docker Integration**: Enhanced Docker and Docker Compose detection and error handling
- **Configuration Management**: Persistent configuration with better defaults
- **Data Management**: Improved volume handling and data persistence
- **Error Handling**: Better error messages and troubleshooting guidance

### Fixed

- **Template Processing**: Fixed Docker Compose template generation
- **Port Conflicts**: Better handling of port conflicts and validation
- **Service Dependencies**: Proper handling of service dependencies (e.g., Kafka + Zookeeper)

## [0.1.0] - 2025-01-09

### Added

- **Core CLI Framework**: Basic command structure with Cobra
- **Interactive Initialization**: `dockenv init` command with service selection
- **Service Management**: `dockenv up`, `dockenv down`, `dockenv restart` commands
- **Docker Integration**: Docker and Docker Compose validation and execution
- **Service Templates**: Docker Compose templates for MySQL, PostgreSQL, Redis, MongoDB, Kafka
- **Volume Management**: Persistent data storage with Docker volumes
- **Configuration System**: YAML-based configuration management
- **Basic Profiles**: Initial support for Laravel, Node.js, Django profiles
- **Installer Script**: One-liner installation script for multiple platforms
- **Documentation**: Comprehensive README with usage examples

### Supported Services

- MySQL 8.0
- PostgreSQL 15
- Redis 7
- MongoDB 7
- Apache Kafka with Zookeeper

### Supported Platforms

- Linux (AMD64, ARM64)
- macOS (Intel, Apple Silicon)
- Windows (AMD64)

## [Unreleased]

### Planned Features

- **TUI Mode**: Terminal UI for service management
- **GUI Dashboard**: Web-based dashboard for visual service management
- **Custom Templates**: Support for user-defined service templates
- **Backup/Restore**: Data backup and restore functionality
- **Service Health Checks**: Advanced health monitoring and alerts
- **Multi-Project Support**: Manage multiple project environments
- **Cloud Integration**: Support for cloud-based development environments
- **Plugin System**: Extensible plugin architecture
- **IDE Integration**: VS Code and other IDE extensions

### Known Issues

- Systemd auto-start requires sudo privileges
- Windows support is experimental
- Custom Docker images not yet supported
- No built-in backup solution for data volumes

---

## Release Notes

### v0.2.0 Release Highlights

This release focuses on advanced service management and developer productivity features:

**🚀 Dynamic Service Management**

- Add and remove services without recreating your entire environment
- Smart configuration updates and Docker Compose regeneration

**⚡ Auto-start Integration**

- Services can now start automatically on system boot
- Full systemd integration for Linux systems

**🎯 Enhanced Project Detection**

- Automatic detection of project types with suggested service configurations
- Smart defaults based on your development stack

**📊 Improved Monitoring**

- Better status reporting and logging capabilities
- Real-time log following and service health checks

**🔧 Developer Experience**

- Comprehensive help system and error messages
- Better validation and conflict resolution
- Enhanced configuration management

### Migration from v0.1.0

No breaking changes - existing configurations will continue to work. To take advantage of new features:

1. Update dockenv: `curl -s https://raw.githubusercontent.com/mohammed-bageri/dockenv/main/install.sh | bash`
2. Run `dockenv status` to verify your existing setup
3. Use `dockenv list` to see new available services and profiles
4. Try `dockenv autostart enable` for automatic startup

### Contributing

We welcome contributions! See our [Contributing Guide](CONTRIBUTING.md) for details.

### Support

- 📚 [Documentation](README.md)
- 🐛 [Issue Tracker](https://github.com/mohammed-bageri/dockenv/issues)
- 💬 [Discussions](https://github.com/mohammed-bageri/dockenv/discussions)
