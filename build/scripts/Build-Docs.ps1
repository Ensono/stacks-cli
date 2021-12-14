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
    $targets = @("pdf", "md")
)

# Determine the path to the project
if ([string]::IsNullOrEmpty($BasePath)) {
    $BasePath = Get-Location
}

# Define where the raw documents exist
$DocsDir = [IO.Path]::Combine($BasePath, "docs")

# Iterate around the targets as the format of each doc
$targets | ForEach-Object -Parallel {

    $format = $_

    # Set the output directory
    $OutputDir = [IO.Path]::Combine($using:BasePath, $using:OutputDir, "docs", $format)

    if (!(Test-Path -Path $OutputDir)) {
        Write-Output ("Creating output dir: {0}" -f $OutputDir)
        New-Item -ItemType Directory -Path $OutputDir | Out-Null
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
        $list = Get-ChildItem -Path $using:DocsDir/* -Attributes !Directory -Include *.adoc

        # iterate around the files
        foreach ($file in $list) {

            # define the paths for the xml and md files
            $xml_file = [IO.Path]::Combine($env:TEMP, ("{0}.xml" -f [System.IO.Path]::GetFileNameWithoutExtension($file.FullName)))
            $md_file = [IO.Path]::Combine($OutputDir, ("{0}.md" -f[System.IO.Path]::GetFileNameWithoutExtension($file.FullName)))

            # Create the cmd to run
            # -- convert to xml
            Invoke-Expression -Command ("asciidoctor -b docbook -o {0} {1}" -f $xml_file, $file.FullName)

            # -- convert to markdown
            Invoke-Expression -Command ("pandoc -f docbook -t markdown_strict {0} -o {1}" -f $xml_file, $md_file)

            # Remove the xml_file
            if (Test-Path -Path $xml_file) {
                Remove-Item -Path $xml_file
            }
        }

        # Copy the images into the md output dir
        Copy-Item -Path $using:DocsDir/images -Destination $OutputDir/ -Recurse
    }


}