<#

.SYNOPSIS
Build the different formats of the documentation

#>

[CmdletBinding()]
param (

    [string]
    # Set the base path
    $BasePath = $null,

    [string]
    # Output directory for reports
    $OutputDir = "outputs",

    [string]
    # Set the build number
    $BuildNumber = $env:BUILDNUMBER,

    [string[]]
    # operating system that binary should be build for
    $targets = @("pdf", "md"),

    [string[]]
    # list of files that should be exluded from markdown conversion
    $excludes = @("manual.adoc")
)

# Get the script dir so that functions can be imported
$scriptDir = Split-Path -Path $MyInvocation.MyCommand.Path -Parent

# Determine the path to the project
if ([string]::IsNullOrEmpty($BasePath)) {
    $BasePath = Get-Location
}

# Define where the raw documents exist
$DocsDir = [IO.Path]::Combine($BasePath, "docs")

# Iterate around the targets as the format of each doc
$targets | ForEach-Object -Parallel {

    $format = $_

    . $using:scriptDir/functions/ConvertTo-MDX.ps1

    # Set the output directory
    $OutputDir = [IO.Path]::Combine($using:BasePath, $using:OutputDir, "docs", $format)

    # Create a temporary directory to be used to store transitional files
    $TempDir = [IO.Path]::Combine($using:BasePath, $using:OutputDir, "docs", "temp")

    if (!(Test-Path -Path $OutputDir)) {
        Write-Output ("Creating output dir: {0}" -f $OutputDir)
        New-Item -ItemType Directory -Path $OutputDir | Out-Null
    }

    if (!(Test-Path -Path $TempDir)) {
        Write-Output ("Creating temporary output dir: {0}" -f $TempDir)
        New-Item -ItemType Directory -Path $TempDir | Out-Null
    }

    if ($format -eq "pdf") {

        # Build the command to generate the PDF
        $cmd_parts = @(
            "asciidoctor-pdf",
            ("-a pdf-theme={0}/styles/theme.yml" -f $using:DocsDir),
            ('-a pdf-fontsdir="{0}/styles/fonts;GEM_FONTS_DIR"' -f $using:DocsDir),
            '-a doctype="book"',
            ('-a stackscli_version="{0}"' -f $using:BuildNumber),
            ('-o "Stacks CLI Manual - {0}.pdf"' -f $using:BuildNumber),
            "-a toc",
            ("-D {0}" -f $OutputDir),
            ("{0}/manual.adoc" -f $using:DocsDir)
        )

        # Run the command
        $cmd = $cmd_parts -join " "
        Write-Output $cmd
        Invoke-Expression -Command $cmd
    }

    if ($format -eq "md") {

        # get a list of the docs
        $list = Get-ChildItem -Path $using:DocsDir/* -Attributes !Directory -Include *.adoc -Exclude $using:excludes

        # iterate around the files
        foreach ($file in $list) {

            # define the paths for the xml and md files
            $xml_file = [IO.Path]::Combine($TempDir, ("{0}.xml" -f [System.IO.Path]::GetFileNameWithoutExtension($file.FullName)))
            $md_file = [IO.Path]::Combine($OutputDir, ("{0}.md" -f[System.IO.Path]::GetFileNameWithoutExtension($file.FullName)))
            $mdx_file = [IO.Path]::Combine($using:BasePath, $using:OutputDir, "docs", "mdx", ("{0}.mdx" -f[System.IO.Path]::GetFileNameWithoutExtension($file.FullName)))

            # Create the cmd to run
            Write-Information -MessageData "Converting:"
            # -- convert to xml
            Write-Information -MessageData "`tDockbook XML"
            Invoke-Expression -Command ("asciidoctor -b docbook -o {0} {1}" -f $xml_file, $file.FullName)

            # -- escape any characters that are going to cause issues, e.g. <# and #> in powershell
            Write-Information -MessageData "`tEscaping known problem characters"

            # create hash of strings to look for and their replacements
            $patterns = @{
                "<#" = "&lt;#"
                "#>" = "#&gt;"
                "<version>" = "&lt;version&gt;"
            }

            $xml_content = Get-Content -Path $xml_file -Raw

            # iterate around the patterns and perform the replacements as required
            foreach ($item in $patterns.GetEnumerator()) {
                $xml_content = $xml_content.Replace($item.Name, $item.Value)
            }
            
            Set-Content -Path $xml_file -Value $xml_content

            # -- convert to markdown
            Write-Information -MessageData "`tMarkdown"
            $resp = Invoke-Expression -Command ("pandoc -f docbook -t gfm --wrap none {0} -o {1}" -f $xml_file, $md_file) -ErrorVariable errors

            if (!(Test-Path -Path $md_file) -or $LASTEXITCODE -gt 0) {
                Write-Error -Message ("Error creating Markdown file: {0}`n{1}" -f $md_file, $errors[0])
                continue
            }

            # -- convert to MDX format to handle Docusaurus
            ConvertTo-MDX -Path $md_file -Destination $mdx_file

            # Remove the xml_file
            if (Test-Path -Path $xml_file) {
                Remove-Item -Path $xml_file
            }
        }

        # Copy the images into the md output dir
        Copy-Item -Path $using:DocsDir/images -Destination $OutputDir/ -Recurse -ErrorAction SilentlyContinue
        Copy-Item -Path $using:DocsDir/images -Destination "$(Split-Path -Path $mdx_file -Parent)/" -Recurse -ErrorAction SilentlyContinue
    }


}