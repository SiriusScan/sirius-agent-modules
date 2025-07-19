#!/usr/bin/env pwsh

# Sirius Agent Module: Windows Service Permission Check
# Detects Windows services with weak security configurations

param(
    [string]$ServiceName = ""
)

# Output JSON structure
$result = @{
    vulnerable = $false
    vulnerability_id = "CUSTOM-WINDOWS-SERVICE-001"
    severity = "medium"
    confidence = 0.0
    evidence = @()
    metadata = @{
        checked_services = @()
        vulnerable_services = @()
    }
}

try {
    # Get all services if no specific service provided
    if ($ServiceName -eq "") {
        $services = Get-Service | Where-Object {$_.Status -eq "Running"}
    } else {
        $services = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
    }

    $vulnerableServices = @()
    $checkedServices = @()

    foreach ($service in $services) {
        $checkedServices += $service.Name
        
        try {
            # Check service permissions using WMI
            $serviceObj = Get-WmiObject -Class Win32_Service -Filter "Name='$($service.Name)'"
            
            if ($serviceObj) {
                # Check if service runs under LocalSystem (potentially dangerous)
                if ($serviceObj.StartName -eq "LocalSystem") {
                    $vulnerableServices += $service.Name
                    $result.evidence += "Service '$($service.Name)' runs under LocalSystem account"
                }
                
                # Check if service has weak permissions
                $acl = Get-Acl "HKLM:\SYSTEM\CurrentControlSet\Services\$($service.Name)"
                $everyoneAccess = $acl.Access | Where-Object {$_.IdentityReference -eq "Everyone"}
                
                if ($everyoneAccess) {
                    $vulnerableServices += $service.Name
                    $result.evidence += "Service '$($service.Name)' has weak permissions (Everyone access)"
                }
            }
        } catch {
            # Service might not be accessible, skip
            continue
        }
    }

    # Update result based on findings
    if ($vulnerableServices.Count -gt 0) {
        $result.vulnerable = $true
        $result.confidence = 0.8
        $result.severity = "high"
        $result.metadata.vulnerable_services = $vulnerableServices
    } else {
        $result.evidence += "No vulnerable service configurations found"
    }

    $result.metadata.checked_services = $checkedServices

} catch {
    $result.evidence += "Error checking services: $($_.Exception.Message)"
    $result.confidence = 0.0
}

# Output JSON result
$result | ConvertTo-Json -Depth 10 