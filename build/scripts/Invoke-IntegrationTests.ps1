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
    $report,

    [string]
    # Project directory to use to write out files
    # and new projects
    # Will be created if does not exist. Relative path will be appended to the
    # current directory
    $projectDir = "inttest"
)

# If the output directory does not exist, create it
if (!(Test-Path -Path $output_dir)) {
    Write-Output ("Creating output dir: {0}" -f $output_dir)
    New-Item -ItemType Directory -Path $output_dir | Out-Null
}

# Define paths for outputs
$temp_report_file = [IO.Path]::Combine("outputs", "tests", "integration_test_report.out")
$test_report_file = [IO.Path]::Combine("outputs", "tests", "integration_test_report.xml")

# Build up the path to the test_binary and the cli_binary
$test_binary = [IO.Path]::Combine("/", "eirctl", "outputs", "bin", $("stacks-cli-linux-inttest-amd64-{0}" -f $build_number))
$cli_binary = [IO.Path]::Combine("/", "eirctl", "outputs", "bin", $("stacks-cli-linux-amd64-{0}" -f $build_number))

# If running on Linux ensure that the binaries have the correct permissions set
if ($IsLinux) {
    $cmd = "chmod +x {0}/*" -f (Split-Path -Parent $cli_binary)
    Invoke-Expression -Command $cmd
}

# Run the tests if they have been specified
if ($runtests) {

    # determine the path to the projectDir
    if (![System.IO.Path]::IsPathRooted($projectDir)) {
        $projectDir = [IO.Path]::Combine($pwd, $projectDir)
    }

    # ensure that the project dir exists
    if (!(Test-Path -Path $projectDir)) {
        Write-Output ("Creating project dir: {0}" -f $projectDir)
        New-Item -ItemType Directory -Path $projectDir | Out-Null
    }

    $cmd = "{0} --test.v --projectdir {1} --binarycmd {2} | Tee-Object -FilePath {3}" -f
    $test_binary,
    $projectDir,
    $cli_binary,
    $temp_report_file

    Write-Output $cmd

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
