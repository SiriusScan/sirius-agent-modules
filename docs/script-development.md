# Script Development Guide

This guide explains how to create custom vulnerability detection scripts for the Sirius agent.

## Script Types

### PowerShell Scripts (.ps1)

For Windows-specific vulnerability detection.

#### Template Structure

```powershell
<#
.SYNOPSIS
Brief description of what the script detects

.VULNERABILITY
CVE-2024-XXXX or CUSTOM-ID

.DESCRIPTION
Detailed description of the vulnerability being detected

.SEVERITY
critical|high|medium|low

.AUTHOR
your-name

.VERSION
1.0
#>

param(
    [Parameter(Mandatory=$false)]
    [string]$ConfigPath = ""
)

function Test-Vulnerability {
    try {
        $result = @{
            "vulnerability_id" = "CVE-2024-XXXX"
            "vulnerable" = $false
            "confidence" = 0.0
            "evidence" = @()
            "metadata" = @{}
            "error" = $null
        }

        # Your detection logic here

        return $result | ConvertTo-Json -Depth 4

    } catch {
        return @{
            "vulnerability_id" = "CVE-2024-XXXX"
            "vulnerable" = $null
            "confidence" = 0.0
            "evidence" = @()
            "metadata" = @{}
            "error" = $_.Exception.Message
        } | ConvertTo-Json -Depth 2
    }
}

Test-Vulnerability
```

### Bash Scripts (.sh)

For Linux/macOS vulnerability detection.

#### Template Structure

```bash
#!/bin/bash

# Script metadata
VULNERABILITY_ID="CVE-2024-XXXX"
SEVERITY="medium"
DESCRIPTION="Description of vulnerability"
AUTHOR="your-name"
VERSION="1.0"

detect_vulnerability() {
    local vulnerable=false
    local confidence=0.0
    local evidence=()
    local error=""

    # Your detection logic here

    # Return JSON result
    cat << EOF
{
    "vulnerability_id": "$VULNERABILITY_ID",
    "vulnerable": $vulnerable,
    "confidence": $confidence,
    "evidence": [$(printf '%s,' "${evidence[@]}" | sed 's/,$//')],
    "metadata": {},
    "error": $(if [[ -n "$error" ]]; then echo "\"$error\""; else echo "null"; fi)
}
EOF
}

detect_vulnerability
```

### Python Scripts (.py)

For cross-platform vulnerability detection.

#### Template Structure

```python
#!/usr/bin/env python3

import json
import sys
import platform

# Script metadata
VULNERABILITY_ID = "CVE-2024-XXXX"
SEVERITY = "medium"
DESCRIPTION = "Description of vulnerability"
AUTHOR = "your-name"
VERSION = "1.0"

def detect_vulnerability():
    """Main vulnerability detection function"""
    result = {
        "vulnerability_id": VULNERABILITY_ID,
        "vulnerable": False,
        "confidence": 0.0,
        "evidence": [],
        "metadata": {
            "platform": platform.system(),
            "python_version": platform.python_version()
        },
        "error": None
    }

    try:
        # Your detection logic here
        pass

    except Exception as e:
        result["error"] = str(e)
        result["vulnerable"] = None

    return result

if __name__ == "__main__":
    result = detect_vulnerability()
    print(json.dumps(result, indent=2))
```

## Output Format Requirements

All scripts must output valid JSON with the following structure:

```json
{
  "vulnerability_id": "CVE-2024-XXXX",
  "vulnerable": true,
  "confidence": 0.9,
  "evidence": [
    {
      "type": "file_permission",
      "location": "/path/to/file",
      "description": "Description of evidence",
      "details": "Additional details"
    }
  ],
  "metadata": {
    "key": "value"
  },
  "error": null
}
```

### Field Descriptions

- **vulnerability_id**: CVE number or custom identifier
- **vulnerable**: Boolean indicating if vulnerability was detected (null for errors)
- **confidence**: Float between 0.0 and 1.0 indicating detection confidence
- **evidence**: Array of evidence objects supporting the detection
- **metadata**: Additional context information
- **error**: Error message if script execution failed

## Best Practices

1. **Error Handling**: Always wrap detection logic in try/catch blocks
2. **Timeouts**: Keep script execution under 30 seconds
3. **Cross-Platform**: Consider platform differences when writing scripts
4. **Documentation**: Include clear comments explaining detection logic
5. **Testing**: Test scripts on multiple systems before submission
6. **Validation**: Validate all inputs and handle edge cases

## Testing Your Scripts

Before submitting, test your scripts using:

```bash
# PowerShell
powershell -File your-script.ps1

# Bash
bash your-script.sh

# Python
python3 your-script.py
```

Ensure the output is valid JSON and follows the required format.

## Submission Process

1. Create script in appropriate platform directory
2. Test thoroughly on target platforms
3. Update `scripts/manifest.json` with script metadata
4. Submit pull request with description of vulnerability detected
5. Include test cases and example outputs

## Example Evidence Types

- **file_permission**: File permission issues
- **registry_key**: Windows registry vulnerabilities
- **service_config**: Service configuration problems
- **network_config**: Network misconfigurations
- **certificate**: Certificate validation issues
- **process**: Running process vulnerabilities
- **user_account**: User account security issues
