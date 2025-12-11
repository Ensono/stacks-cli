
[CmdletBinding()]
param (
    [string]$VariablesFile = "$PSScriptRoot/../variables.yml",
    [string]$ConfigFile = "$PSScriptRoot/../../docs/conf/docs.json",
    [string]$OutputConfigFile = "$PSScriptRoot/../../docs/conf/docs.gen.json"
)

$ErrorActionPreference = "Stop"

if (-not (Test-Path $VariablesFile)) {
    Write-Error "Variables file not found at $VariablesFile"
}

if (-not (Test-Path $ConfigFile)) {
    Write-Error "Config file not found at $ConfigFile"
}

Write-Host "Reading variables from $VariablesFile"
$content = Get-Content -Path $VariablesFile -Raw

if ($content -match "name:\s*EirctlVersion\s+value:\s*([^\s]+)") {
    $version = $matches[1]
    Write-Host "Found EirctlVersion: $version"

    Write-Host "Reading config from $ConfigFile"
    $jsonContent = Get-Content -Path $ConfigFile -Raw
    $config = $jsonContent | ConvertFrom-Json

    $attrString = "eirctl_version=$version"

    # Update PDF attributes
    if ($config.pdf -and $config.pdf.attributes -and $config.pdf.attributes.asciidoctor) {
        $config.pdf.attributes.asciidoctor += $attrString
    }

    # Update HTML attributes
    if ($config.html -and $config.html.attributes) {
        $config.html.attributes += $attrString
    }

    # Update DOCX attributes
    if ($config.docx -and $config.docx.attributes -and $config.docx.attributes.asciidoctor) {
        $config.docx.attributes.asciidoctor += $attrString
    }

    Write-Host "Writing updated config to $OutputConfigFile"
    $config | ConvertTo-Json -Depth 10 | Set-Content -Path $OutputConfigFile

} else {
    Write-Warning "Could not find EirctlVersion in $VariablesFile"
    Copy-Item -Path $ConfigFile -Destination $OutputConfigFile -Force
}
