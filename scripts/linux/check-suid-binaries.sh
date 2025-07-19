#!/bin/bash

# Sirius Agent Module: SUID Binary Check
# Detects SUID binaries that may pose security risks

# Output JSON structure
cat << 'EOF'
{
  "vulnerable": false,
  "vulnerability_id": "CUSTOM-SUID-001",
  "severity": "medium",
  "confidence": 0.0,
  "evidence": [],
  "metadata": {
    "checked_binaries": [],
    "vulnerable_binaries": []
  }
}
EOF

# Initialize variables
VULNERABLE=false
CONFIDENCE=0.0
EVIDENCE=()
CHECKED_BINARIES=()
VULNERABLE_BINARIES=()

# Function to check if binary is in dangerous SUID list
check_dangerous_suid() {
    local binary="$1"
    local dangerous_suid=(
        "/bin/bash"
        "/bin/sh"
        "/bin/cat"
        "/bin/less"
        "/bin/more"
        "/bin/nano"
        "/bin/vi"
        "/bin/vim"
        "/usr/bin/editor"
        "/usr/bin/view"
    )
    
    for dangerous in "${dangerous_suid[@]}"; do
        if [[ "$binary" == "$dangerous" ]]; then
            return 0  # Found dangerous SUID
        fi
    done
    return 1  # Not dangerous
}

# Function to check if binary is writable
check_writable() {
    local binary="$1"
    if [[ -w "$binary" ]]; then
        return 0  # Writable
    fi
    return 1  # Not writable
}

# Find all SUID binaries
echo "Scanning for SUID binaries..." >&2

# Use find to locate SUID binaries
while IFS= read -r -d '' binary; do
    CHECKED_BINARIES+=("$binary")
    
    # Check if binary is dangerous
    if check_dangerous_suid "$binary"; then
        VULNERABLE=true
        VULNERABLE_BINARIES+=("$binary")
        EVIDENCE+=("Dangerous SUID binary found: $binary")
        CONFIDENCE=0.9
    fi
    
    # Check if binary is writable
    if check_writable "$binary"; then
        VULNERABLE=true
        VULNERABLE_BINARIES+=("$binary")
        EVIDENCE+=("Writable SUID binary found: $binary")
        CONFIDENCE=0.8
    fi
    
    # Check for unusual SUID binaries in user directories
    if [[ "$binary" =~ ^/home/ ]] || [[ "$binary" =~ ^/tmp/ ]] || [[ "$binary" =~ ^/var/tmp/ ]]; then
        VULNERABLE=true
        VULNERABLE_BINARIES+=("$binary")
        EVIDENCE+=("SUID binary in user/temp directory: $binary")
        CONFIDENCE=0.7
    fi
    
done < <(find / -type f -perm -4000 -print0 2>/dev/null)

# If no vulnerabilities found
if [[ "$VULNERABLE" == "false" ]]; then
    EVIDENCE+=("No dangerous SUID binaries found")
fi

# Output final JSON result
cat << EOF
{
  "vulnerable": $VULNERABLE,
  "vulnerability_id": "CUSTOM-SUID-001",
  "severity": "$([[ "$VULNERABLE" == "true" ]] && echo "high" || echo "medium")",
  "confidence": $CONFIDENCE,
  "evidence": $(printf '%s\n' "${EVIDENCE[@]}" | jq -R . | jq -s .),
  "metadata": {
    "checked_binaries": $(printf '%s\n' "${CHECKED_BINARIES[@]}" | jq -R . | jq -s .),
    "vulnerable_binaries": $(printf '%s\n' "${VULNERABLE_BINARIES[@]}" | jq -R . | jq -s .)
  }
}
EOF 