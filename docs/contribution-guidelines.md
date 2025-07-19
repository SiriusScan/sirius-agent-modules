# Contribution Guidelines

Thank you for your interest in contributing to Sirius Agent Modules! This document provides guidelines for contributing templates and scripts to the repository.

## Getting Started

1. **Fork the repository** to your GitHub account
2. **Clone your fork** locally
3. **Create a feature branch** for your contribution
4. **Make your changes** following the guidelines below
5. **Test your changes** thoroughly
6. **Submit a pull request** with a clear description

## Template Contributions

### Template Structure

Templates should be YAML files with the following structure:

```yaml
id: "unique-template-id"
info:
  name: "Human-readable name"
  description: "Brief description of what this template detects"
  author: "Your name or organization"
  severity: "critical|high|medium|low"
  platforms: ["windows", "linux", "macos"]
  vulnerability_ids: ["CVE-2023-1234", "CUSTOM-001"]
  created: "2024-01-15T10:30:00Z"
  updated: "2024-01-15T10:30:00Z"

detection:
  type: "file-hash|registry|config-file"
  targets:
    - path: "/path/to/file"
      hash: "sha256:abc123..."
      description: "Description of this target"

  conditions:
    - type: "file_exists"
      path: "/path/to/file"
    - type: "hash_matches"
      path: "/path/to/file"
      hash: "sha256:abc123..."

remediation:
  description: "How to fix this vulnerability"
  steps:
    - "Step 1: Update the software"
    - "Step 2: Apply security patch"
  references:
    - "https://example.com/security-advisory"
```

### Template Guidelines

- Use descriptive, unique IDs
- Include comprehensive vulnerability information
- Provide clear remediation steps
- Test templates on target platforms
- Follow existing naming conventions

## Script Contributions

### Script Structure

Scripts should output JSON in the following format:

```json
{
  "vulnerable": true,
  "vulnerability_id": "CVE-2023-1234",
  "severity": "high",
  "confidence": 0.95,
  "evidence": [
    "Found vulnerable version 1.2.3 of software X",
    "Security patch available in version 1.2.4"
  ],
  "metadata": {
    "detected_version": "1.2.3",
    "patch_version": "1.2.4",
    "cve_references": [
      "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2023-1234"
    ]
  }
}
```

### Script Guidelines

- Output valid JSON to stdout
- Handle errors gracefully
- Include comprehensive error messages
- Test on target platforms
- Follow security best practices
- Use appropriate exit codes

## Quality Standards

### Security

- Scripts should not perform destructive operations
- Templates should not access sensitive data unnecessarily
- All code should be reviewed for security implications

### Performance

- Scripts should complete within reasonable time limits
- Templates should be efficient and not cause system slowdown
- Avoid unnecessary file system or network operations

### Reliability

- Test on multiple platforms when applicable
- Handle edge cases and error conditions
- Provide clear error messages and logging

## Testing Requirements

### For Templates

- Test on target platforms
- Verify detection accuracy
- Test with both vulnerable and non-vulnerable systems
- Validate YAML syntax

### For Scripts

- Test on target platforms
- Verify JSON output format
- Test error handling
- Validate security controls

## Submission Process

1. **Create your contribution** following the guidelines above
2. **Add to appropriate directory**:
   - Templates: `templates/{type}/`
   - Scripts: `scripts/{platform}/`
3. **Update manifests** if adding new files
4. **Test thoroughly** on target platforms
5. **Submit pull request** with:
   - Clear description of changes
   - Testing results
   - Platform compatibility information

## Review Process

All contributions are reviewed for:

- **Functionality**: Does it work as intended?
- **Security**: Are there any security concerns?
- **Quality**: Does it meet our standards?
- **Documentation**: Is it well-documented?

## Code of Conduct

- Be respectful and constructive
- Focus on technical merit
- Help improve the project
- Follow security best practices

## Questions?

If you have questions about contributing, please:

1. Check the existing documentation
2. Look at existing examples
3. Open an issue for clarification

Thank you for contributing to Sirius Agent Modules!
