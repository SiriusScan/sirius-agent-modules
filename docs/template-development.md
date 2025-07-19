# Template Development Guide

This guide explains how to create YAML vulnerability detection templates for the Sirius agent.

## Template Types

### Hash-Based Detection
Detects vulnerabilities by matching file hashes against known vulnerable binaries.

#### Template Structure

```yaml
id: "UNIQUE-ID-001"
info:
  name: "Descriptive Template Name"
  author: "your-name"
  severity: "critical|high|medium|low"
  description: "Detailed description of what this template detects"
  references:
    - "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-XXXX"
  cve: "CVE-2024-XXXX"
  created: "2024-01-15T10:30:00Z"
  updated: "2024-01-15T10:30:00Z"

detection:
  type: "file-hash"
  method: "sha256"  # sha256, sha1, or md5
  targets:
    - path: "/path/to/vulnerable/file"
      hash: "actual_hash_value"
      description: "Description of this specific file"
    - path: "C:\\Windows\\path\\to\\file.exe"
      hash: "windows_hash_value"
      description: "Windows version of vulnerable file"

  conditions:
    - file_exists: true
    - hash_match: true
    - file_executable: true  # optional

  metadata:
    confidence: 0.95
    impact: "Description of potential impact"
    affected_versions: ["1.0", "1.1"]
    fixed_versions: ["1.2+"]

remediation:
  description: "How to fix this vulnerability"
  commands:
    linux: "sudo apt update && sudo apt install package"
    windows: "Update instructions for Windows"
  verification:
    command: "command to verify fix"
    expected_pattern: "regex pattern for expected output"
```

### Registry-Based Detection (Windows)
Detects vulnerabilities by checking Windows registry keys and values.

#### Template Structure

```yaml
id: "REGISTRY-001"
info:
  name: "Windows Registry Vulnerability"
  author: "your-name"
  severity: "medium"
  description: "Detects vulnerable software via registry entries"
  cve: "CVE-2024-XXXX"

detection:
  type: "registry"
  keys:
    - path: "HKLM\\SOFTWARE\\Vendor\\Product"
      value: "Version"
      pattern: "^1\\.[0-3]\\."  # regex for vulnerable versions
      description: "Vulnerable version in registry"

  conditions:
    - key_exists: true
    - value_matches_pattern: true

remediation:
  description: "Update the software to latest version"
  commands:
    windows: "Download from vendor website"
```

### Config-Based Detection
Detects misconfigurations by scanning configuration files.

#### Template Structure

```yaml
id: "CONFIG-001"
info:
  name: "Configuration Vulnerability"
  author: "your-name"
  severity: "high"
  description: "Detects insecure configuration settings"
  cve: "CUSTOM-CONFIG-001"

detection:
  type: "config-file"
  files:
    - path: "/etc/service/config.conf"
      patterns:
        - regex: "^\\s*debug\\s*=\\s*true"
          description: "Debug mode enabled in production"
          severity: "medium"
        - regex: "^\\s*password\\s*=\\s*\\w+"
          description: "Hardcoded password in config"
          severity: "high"

  conditions:
    - file_exists: true
    - pattern_found: true

remediation:
  description: "Review and secure configuration file"
  commands:
    linux: "Edit /etc/service/config.conf and restart service"
```

## Field Reference

### Required Fields

#### info section
- **id**: Unique identifier for the template
- **name**: Human-readable template name
- **author**: Template author
- **severity**: One of critical, high, medium, low
- **description**: Detailed description of vulnerability
- **cve**: CVE number or custom identifier

#### detection section
- **type**: One of file-hash, registry, config-file
- **conditions**: Array of conditions that must be met

### Optional Fields

#### info section
- **references**: Array of reference URLs
- **created**: ISO 8601 timestamp
- **updated**: ISO 8601 timestamp

#### detection section
- **metadata**: Additional information about the vulnerability

#### remediation section
- **description**: How to fix the vulnerability
- **commands**: Platform-specific fix commands
- **verification**: How to verify the fix

## Hash-Based Templates

### Generating File Hashes

To create hash-based templates, you need the actual hashes of vulnerable files:

```bash
# SHA256
sha256sum /path/to/file

# SHA1  
sha1sum /path/to/file

# MD5
md5sum /path/to/file
```

### Multi-Platform Support

Include hashes for the same vulnerability across different platforms:

```yaml
targets:
  - path: "/usr/bin/vulnerable-app"
    hash: "linux_hash_here"
    description: "Linux version"
  - path: "/usr/local/bin/vulnerable-app"
    hash: "macos_hash_here" 
    description: "macOS version"
  - path: "C:\\Program Files\\App\\vulnerable-app.exe"
    hash: "windows_hash_here"
    description: "Windows version"
```

## Registry Templates

### Registry Path Format

Use double backslashes for registry paths:

```yaml
keys:
  - path: "HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\{GUID}"
    value: "DisplayVersion"
```

### Common Registry Locations

- **HKLM\\SOFTWARE\\**: Machine-wide software installations
- **HKCU\\SOFTWARE\\**: User-specific software
- **HKLM\\SYSTEM\\**: System configuration
- **HKLM\\SECURITY\\**: Security settings

## Config-Based Templates

### Regex Patterns

Use proper regex escaping in YAML:

```yaml
patterns:
  - regex: "^\\s*password\\s*=\\s*['\"]?\\w+['\"]?"
    description: "Hardcoded password"
```

### Common Config Files

- **Linux**: `/etc/`, `/usr/local/etc/`, `~/.config/`
- **Windows**: `C:\\ProgramData\\`, `%APPDATA%`
- **macOS**: `/etc/`, `/usr/local/etc/`, `~/Library/`

## Testing Templates

### Validation

Before submitting, validate your YAML:

```bash
# Using yq (install with: pip install yq)
yq . your-template.yaml

# Using Python
python -c "import yaml; yaml.safe_load(open('your-template.yaml'))"
```

### Testing Detection

Test your templates against known vulnerable and non-vulnerable systems.

## Best Practices

1. **Unique IDs**: Use descriptive, unique identifiers
2. **Accurate Hashes**: Verify all file hashes are correct
3. **Cross-Platform**: Include paths for all relevant platforms
4. **Clear Descriptions**: Write clear, informative descriptions
5. **Proper Severity**: Assign appropriate severity levels
6. **Remediation**: Always include fix instructions
7. **Testing**: Test templates thoroughly before submission

## Common Patterns

### Version Detection

```yaml
# Registry version check
pattern: "^([0-9]+)\\.([0-9]+)\\.([0-9]+)"

# Config file version check  
regex: "version\\s*[=:]\\s*([0-9\\.]+)"
```

### File Permission Checks

```yaml
# Usually handled in conditions
conditions:
  - file_exists: true
  - file_readable: true
  - file_executable: false  # Should not be executable
```

### Multi-Condition Logic

```yaml
conditions:
  - file_exists: true
  - hash_match: true
  - AND:
    - file_size_greater_than: 1024
    - file_modified_before: "2023-01-01"
```

## Submission Process

1. Create template in appropriate type directory
2. Validate YAML syntax and structure
3. Test against known vulnerable systems
4. Update `templates/manifest.json` 
5. Submit pull request with:
   - Template description
   - Test results
   - Vulnerability details
   - Fix verification 