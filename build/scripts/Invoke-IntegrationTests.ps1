<#

.SYNOPSIS
Run integration tests

#>

[CmdletBinding()]
param (

    [string]
    # Output directory for reports
    $output_dir = "outputs/tests",

    [string]
    # Build number of the test to run
    $build_number = $env:BUILDNUMBER,

    [switch]
    # Specify if the tests should be run
    $runtests,

    [switch]
    # If set generate the report
    $report
)

# If the output directory does not exist, create it
if (!(Test-Path -Path $output_dir)) {
    Write-Output ("Creating output dir: {0}" -f $output_dir)
    New-Item -ItemType Directory -Path $output_dir
}

# Define paths for outputs
$temp_report_file = [IO.Path]::Combine("outputs", "tests", "integration_test_report.out")
$test_report_file = [IO.Path]::Combine("outputs", "tests", "integration_test_report.xml")

# Build up the path to the test_binary and the cli_binary
$test_binary = [IO.Path]::Combine("/", "app", "outputs", "bin", $("stacks-cli-linux-inttest-amd64-{0}" -f $build_number))
$cli_binary = [IO.Path]::Combine("/", "app", "outputs", "bin", $("stacks-cli-linux-amd64-{0}" -f $build_number))

# Run the tests if they have been specified
if ($runtests) {
    $cmd = "{0} --test.v --projectdir /app/local/inttest --binarycmd {1} | Tee-Object -FilePath {2}" -f
                $test_binary,
                $cli_binary,
                $temp_report_file

    Invoke-Expression -Command $cmd
}

# Generate the report
if ($report) {
    $cmd = "Get-Content -Path {0} -Raw | go-junit-report > {1}" -f $temp_report_file, $test_report_file

    Invoke-Expression -Command $cmd

    if (Test-Path -Path $temp_report_file) {
        Remove-Item -Path $temp_report_file
    }
}