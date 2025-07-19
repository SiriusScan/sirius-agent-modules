#!/usr/bin/env python3

"""
Sirius Agent Module: Certificate Validation
Detects SSL/TLS certificate vulnerabilities across platforms
"""

import json
import ssl
import socket
import subprocess
import sys
import os
from datetime import datetime, timedelta
from pathlib import Path

def check_certificate_expiry(hostname, port=443):
    """Check if certificate is expired or expiring soon"""
    try:
        context = ssl.create_default_context()
        with socket.create_connection((hostname, port), timeout=10) as sock:
            with context.wrap_socket(sock, server_hostname=hostname) as ssock:
                cert = ssock.getpeercert()
                not_after = datetime.strptime(cert['notAfter'], '%b %d %H:%M:%S %Y %Z')
                days_until_expiry = (not_after - datetime.now()).days
                return {
                    'expired': days_until_expiry < 0,
                    'expiring_soon': days_until_expiry < 30,
                    'days_until_expiry': days_until_expiry,
                    'subject': dict(x[0] for x in cert['subject']),
                    'issuer': dict(x[0] for x in cert['issuer'])
                }
    except Exception as e:
        return {'error': str(e)}

def check_weak_ciphers(hostname, port=443):
    """Check for weak SSL/TLS ciphers"""
    weak_ciphers = [
        'RC4', 'DES', '3DES', 'MD5', 'NULL', 'EXPORT'
    ]
    
    try:
        context = ssl.create_default_context()
        with socket.create_connection((hostname, port), timeout=10) as sock:
            with context.wrap_socket(sock, server_hostname=hostname) as ssock:
                cipher = ssock.cipher()
                cipher_name = cipher[0]
                return {
                    'weak_cipher': any(weak in cipher_name for weak in weak_ciphers),
                    'cipher_name': cipher_name,
                    'cipher_version': cipher[1],
                    'cipher_bits': cipher[2]
                }
    except Exception as e:
        return {'error': str(e)}

def check_certificate_store():
    """Check system certificate store for issues"""
    issues = []
    
    # Check common certificate store locations
    cert_paths = [
        '/etc/ssl/certs',  # Linux
        '/System/Library/Keychains',  # macOS
        'C:\\Windows\\System32\\config\\systemprofile\\AppData\\Roaming\\Microsoft\\Crypto\\RSA\\MachineKeys'  # Windows
    ]
    
    for path in cert_paths:
        if os.path.exists(path):
            try:
                # Check for expired certificates
                for cert_file in Path(path).rglob('*.pem'):
                    # Basic check - in production would use proper certificate parsing
                    if cert_file.stat().st_mtime < (datetime.now() - timedelta(days=365)).timestamp():
                        issues.append(f"Potentially expired certificate: {cert_file}")
            except Exception as e:
                issues.append(f"Error checking {path}: {e}")
    
    return issues

def main():
    """Main function to run certificate validation checks"""
    
    # Initialize result structure
    result = {
        'vulnerable': False,
        'vulnerability_id': 'CUSTOM-CERT-001',
        'severity': 'medium',
        'confidence': 0.0,
        'evidence': [],
        'metadata': {
            'checked_hosts': [],
            'vulnerable_hosts': [],
            'certificate_issues': []
        }
    }
    
    # List of hosts to check (could be configurable)
    hosts_to_check = [
        'google.com',
        'github.com',
        'microsoft.com'
    ]
    
    vulnerable_hosts = []
    checked_hosts = []
    
    # Check each host
    for hostname in hosts_to_check:
        checked_hosts.append(hostname)
        
        # Check certificate expiry
        expiry_info = check_certificate_expiry(hostname)
        if 'error' not in expiry_info:
            if expiry_info['expired']:
                result['vulnerable'] = True
                result['evidence'].append(f"Expired certificate for {hostname}")
                vulnerable_hosts.append(hostname)
            elif expiry_info['expiring_soon']:
                result['vulnerable'] = True
                result['evidence'].append(f"Certificate expiring soon for {hostname} ({expiry_info['days_until_expiry']} days)")
                vulnerable_hosts.append(hostname)
        
        # Check for weak ciphers
        cipher_info = check_weak_ciphers(hostname)
        if 'error' not in cipher_info:
            if cipher_info['weak_cipher']:
                result['vulnerable'] = True
                result['evidence'].append(f"Weak cipher detected for {hostname}: {cipher_info['cipher_name']}")
                vulnerable_hosts.append(hostname)
    
    # Check certificate store
    cert_issues = check_certificate_store()
    if cert_issues:
        result['vulnerable'] = True
        result['evidence'].extend(cert_issues)
        result['metadata']['certificate_issues'] = cert_issues
    
    # Update confidence and severity based on findings
    if result['vulnerable']:
        result['confidence'] = 0.8
        result['severity'] = 'high'
    else:
        result['evidence'].append("No certificate vulnerabilities found")
    
    result['metadata']['checked_hosts'] = checked_hosts
    result['metadata']['vulnerable_hosts'] = vulnerable_hosts
    
    # Output JSON result
    print(json.dumps(result, indent=2))

if __name__ == '__main__':
    main() 