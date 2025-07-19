# Sirius Agent Modules

Official repository for Sirius vulnerability detection templates and scripts.

## Overview

This repository contains vulnerability detection templates and scripts for the Sirius security scanning system. It provides a centralized location for community-contributed detection modules that can be automatically downloaded and updated by Sirius agents.

## Repository Structure

```
sirius-agent-modules/
├── templates/           # YAML-based vulnerability detection templates
│   ├── hash-based/     # File hash-based detection templates
│   ├── registry-based/ # Windows registry-based detection templates
│   └── config-based/   # Configuration file-based detection templates
├── scripts/            # Executable vulnerability detection scripts
│   ├── windows/        # PowerShell scripts for Windows
│   ├── linux/          # Bash scripts for Linux
│   └── cross-platform/ # Python scripts for all platforms
└── docs/               # Development and contribution documentation
```

## Components

### Templates

- **Hash-based**: Detect vulnerable software versions via file hashes
- **Registry-based**: Detect Windows registry vulnerabilities and misconfigurations
- **Config-based**: Detect security misconfigurations in configuration files

### Scripts

- **Windows**: PowerShell scripts for Windows-specific vulnerability detection
- **Linux**: Bash scripts for Linux-specific vulnerability detection
- **Cross-platform**: Python scripts that work across all platforms

## Usage

### For Sirius Agents

Agents automatically download and update templates and scripts from this repository. No manual intervention required.

### For Developers

1. Fork this repository
2. Add your templates or scripts following the guidelines in `docs/`
3. Submit a pull request
4. Your contribution will be reviewed and merged

## Contributing

See the documentation in the `docs/` directory:

- [Template Development Guide](docs/template-development.md)
- [Script Development Guide](docs/script-development.md)
- [Contribution Guidelines](docs/contribution-guidelines.md)

## License

MIT License - see LICENSE file for details.

## Support

For questions or issues, please open an issue in this repository.
