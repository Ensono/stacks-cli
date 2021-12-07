<#

.SYNOPSIS
Checks whether the Build Number env var exists and sets default if not

#>

[CmdletBinding()]
param (

    [string]
    # Set the default build number
    $default = "100.98.99"
)

if ([String]::IsNullOrEmpty($env:BUILDNUMBER)) {
    $result = $default
} else {
    $result = $env:BUILDNUMBER
}

# output the build number
Write-Output $result.Trim()
