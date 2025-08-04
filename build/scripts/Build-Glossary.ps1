

[CmdletBinding()]
param (

    [string]
    # Path to the glossary terms file
    $Path = "$PSScriptRoot/../conf/glossary.json",

    [string]
    # Path to the docs directory
    $DocsPath = "$PSScriptRoot/../../docs",

    [string[]]
    # List of files that should be included in the search
    $IncludeFiles = @("*.adoc", "*.md"),

    [string]
    # Path that the glossary file should be written to
    $OutputPath = "$PSScriptRoot/../../local/glossary_table.adoc"
)

# If the glossary file exists, read it in
if (Test-Path -Path $Path) {
    $glossary = Get-Content -Path $Path | ConvertFrom-Json
}
else {
    Write-Error "Glossary file not found at path: $Path"
    return
}

# Ensure that the parent dir for the output exists
$parentOutputDir = Split-Path -Path $OutputPath -Parent
if (!(Test-Path -Path $parentOutputDir)) {
    Write-Output ("Creating output dir: {0}" -f $parentOutputDir)
    New-Item -ItemType Directory -Path $parentOutputDir | Out-Null
}

# Find all the files in the specified directory so that each file
# can be searched for the glossary terms. Doing this here prevents the system from looking
# for the files on each iteration
Write-Information ("Finding documentation files: {0}" -f ($IncludeFiles -join ", "))
$docs = Get-ChildItem -Path $DocsPath -Recurse -Include $IncludeFiles
Write-Information ("`t{0} files found." -f $docs.Count)

# Create variable to hold the terms that have been found in the documentation
$foundTerms = @()

# Iterate around all the terms in the glossary
# and search for each term in the documentation files. As soon as one term is found, the loop
# will move onto the next term.
foreach ($term in $glossary) {

    # Create the regular expression to use to search for the term in the files
    $pattern = "(?:^|\s|\(){0}(?:$|\s|\))" -f $term.term

    # iterate around all the files that have been found
    foreach ($doc in $docs) {

        # Match the pattern in the file
        $result = Select-String -Path $doc.FullName -Pattern $pattern

        # If the term is found then break from this loop and move onto the next term
        if (-not [String]::IsNullOrEmpty($result) -and $foundTerms -notcontains $term) {
            $foundTerms += $term
        }
    }
}

$sortedFoundTerms = $foundTerms | Sort-Object -Property term

Write-Information ("Number of terms found: {0}" -f $sortedFoundTerms.Count)

# Now build up the asciidoc table for the glossary
$sb = [System.Text.StringBuilder]::new()

# Add in the header row
[void]$sb.AppendLine("|===")
[void]$sb.AppendLine("| Term | Expansion | Description")

# Iterate around the sorted terms and add them to the table
foreach ($term in $sortedFoundTerms) {
    [void]$sb.AppendLine(("| {0} | {1} | {2}" -f $term.term, $term.expansion, $term.description))
}

[void]$sb.AppendLine("|===")

Write-Information ("Creating glossary table: {0}" -f $OutputPath)
$sb.ToString() | Out-File -FilePath $OutputPath
