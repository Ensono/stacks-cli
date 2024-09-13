
#function Set-Options() {

    <#
    .SYNOPSIS
    Generates help documentation detailing the options for each command from a configuration file.

    .DESCRIPTION
    The `Set-Options` function reads a configuration file in YAML format that specifies various commands and their associated Go source files. It extracts command options from these Go source files using a regular expression pattern and generates help documentation in AsciiDoc format.

    The function begins by checking if the provided configuration file path exists. If the file is not found, it outputs an error message and terminates. It then sets up a template for generating tables in the documentation, including a table header and row format.

    The function reads the configuration file and iterates over each command specified. For each command, it verifies the existence of the Go source file and reads its content. Using a regular expression pattern, it extracts options from the source file and processes them to determine their names, default values, descriptions, and whether they have shortcuts.

    The extracted information is formatted using the row template and added to a hash table. Finally, the function generates a table for each sub-command and writes the table to an output file in AsciiDoc format. Informational messages are output during processing and writing to provide feedback on the function's progress.
    #>

    [CmdletBinding()]
    param (

        [string]
        # Path tot he configutation file detaling the files to analyse for options
        $Path = "/app/build/conf/docs/options.yml"

    )

    # Check that the Path exists
    if (!(Test-Path -Path $Path)) {
        Write-Error "Configuration file cannot be found: $Path"
        return
    }

    # Set the table template
    $template = @{
        "table" = @"
[cols="2a,1,2,1,1",options="header"]
|===
2+| Parameter | Environment Variable | Default | Permitted Values

{0}
|===
"@
        "row" = @"
.2+^| {0} ^| {1} | {2} | {3} |
4+| {4}
"@
    }

    # Set the Regex to get the options from the file
    # $pattern = '^(\s+)(?<subcommand>.*)\.(?:Persistent|)Flags\(\)\.(?<type>.*?)\((?<definition>.*?)\)'
    $pattern = '^(\s+)(?<subcommand>(?!.*BindPFlag).*?)\.(?:Persistent|)Flags\(\)\.(?<type>.*?)\((?<definition>.*?)\)'

    # Load the configuration file
    $config = Get-Content -Path $Path -Raw | ConvertFrom-Yaml

    # Iterate around the commands and generate each file
    foreach ($command in $config.commands) {

        # Create a hash table to hold the sub-commands and the options
        # This needs to be cleared on each iteration, and setup for each sub-command
        $rows = @{}

        # Ensure that the path to the code to read exists
        if (!(Test-Path -Path $command.path)) {
            Write-Error "Code file cannot be found: $($command.path)"
            continue
        }

        Write-Information -MessageData "Processing: $($command.path)"

        # Read in the file
        $gocode = Get-Content -Path $command.path -Raw

        # Use the regex pattern to get the options from the file
        $result = [regex]::Matches($gocode,
                       $pattern,
                       [System.Text.RegularExpressions.RegexOptions]::Multiline
                       )

        if ([String]::IsNullOrEmpty($result)) {
            Write-Warning "No options found for: $($command.path)"
            continue
        }

        # Iterate around the matchend
        foreach ($match in $result) {

            # Ignore if this a lookup type
            if ($type -ieq "lookup") {
                continue
            }

            # get the subcommand name
            $subcommand = $match.Groups["subcommand"].Value

            # If the hashtable does not contain an entry for the command, add it
            if (-not $rows.ContainsKey($subcommand)) {
                $rows[$subcommand] = @()
            }

            # determine the the type of the input
            $type = $match.Groups["type"].Value

            # get the values of the definition
            $data = ($match.Groups["definition"] -split ",") -replace "`"", ""

            # Determine whether the input has a shortcut or not, this is done by
            # checking the last character of the type
            $short = ""
            $field_name = 1
            if ($type[-1] -eq 'P') {
                $short = $data[2]
                $field_default = 3
                $field_description = 4
            } else {
                $field_default = 2
                $field_description = 3
            }

            $name = $data[$field_name].Trim()
            $default = $data[$field_default].Trim()
            $description = [String]::IsNullOrEmpty($data[$field_description]) ? "" : $data[$field_description].Trim()

            # using the row template set the row data
            # define the array to be used for the substition in the template
            $subs = @()

            if ([String]::IsNullOrEmpty($short)) {
                $subs += "``--{0}``" -f $name
            } else {
                $subs += "``--{0}``, ``-{1}``" -f $name, $short.Trim()
            }

            # determine if the paramneter is required, based on whether it has a default value or not
            $required_indicator = "icon:times[fw]"
            if ([String]::IsNullOrEmpty($default)) {
                $required_indicator = "icon:check[fw]"
            }

            $subs += $required_indicator
            $subs += $name.ToUpper()
            $subs += $default
            $subs += $description

            $row = $template["row"] -f $subs

            $rows[$subcommand] += $row

        }

        # Iterate around the keys in the rows to get the subcommands
        foreach ($sub in $rows.Keys) {

            # set the filename of the output file
            $output_file = [IO.Path]::Combine($config.output, ("{0}.adoc" -f $sub))

            # Create the table for this command
            $table = $template["table"] -f ($rows[$sub] -join "`n")

            # Finally write the table out to the file
            Write-Information -MessageData "Writing: $($output_file)"
            Set-Content -Path $output_file -Value $table
        }

    }

# }
