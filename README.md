# Sirius Agent Modules

Official repository for Sirius vulnerability detection templates and scripts.

## Overview

This repository contains vulnerability detection templates and scripts for the Sirius security scanning system. It provides a centralized location for community-contributed detection modules that can be automatically downloaded and updated by Sirius agents.

## Repository Structure

```
sirius-agent-modules/
├── templates/           # YAML-based vulnerability detection templates
│   ├── file-content/   # File content matching
│   ├── file-hash/      # File hash checks
│   ├── file-search/    # File searches
│   └── version-cmd/    # Installed version checks
├── scripts/            # Executable vulnerability detection scripts
│   ├── windows/        # PowerShell scripts for Windows
│   ├── linux/          # Bash scripts for Linux
│   └── cross-platform/ # Python scripts for all platforms
└── docs/               # Development and contribution documentation
```

## Components

### Templates

Templates use the agent's `detection.steps[]` format with `file_content`,
`file_hash`, `file_search`, and `version_cmd` modules.

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
3. [Regenerate and validate the manifests](docs/template-development.md#regenerate-and-validate-the-manifests)
4. Submit a pull request

## Contributing

See the documentation in the `docs/` directory:

- [Template Development Guide](docs/template-development.md)
- [Script Development Guide](docs/script-development.md)
- [Contribution Guidelines](docs/contribution-guidelines.md)

## License

MIT License - see LICENSE file for details.

## Support

For questions or issues, please open an issue in this repository.
