<#

.SYNOPSIS
Compile the CLI binary and tests

#>

[CmdletBinding()]
param (

    [string]
    # Base directory for all binary builds
    $BasePath = "/app/outputs/bin",

    [string]
    # Set the build number
    $BuildNumber = $env:BUILDNUMBER,

    [string[]]
    # operating system that binary should be build for
    $targets = @("windows", "linux", "darwin"),

    [string[]]
    # arhictectire that should be targeted
    $arch = @("amd64")
)

# If the base directory does not exist, create it
if (!(Test-Path -Path $BasePath)) {
    Write-Output ("Creating output dir: {0}" -f $BasePath)
    New-Item -ItemType Directory -Path $BasePath
}

# Run command to get the packages required for the build
$cmd = "go get"
Invoke-Expression -Command $cmd

# iterate around each of the architectures
foreach ($_arch in $arch) {

    # iterate around each of the target os
    foreach ($os in $targets) {

        # set the filename of the CLI and intest
        $cli_filename = "{0}/stacks-cli-{1}-{2}-{3}" -f $BasePath, $os, $_arch, $BuildNumber
        $inttest_filename = "{0}/stacks-cli-{1}-inttest-{2}-{3}" -f $BasePath, $os, $_arch, $BuildNumber

        # Add the extension if it has been set
        if ($os -eq "windows") {
            $cli_filename += ".exe"
            $inttest_filename += ".exe"
        }

        $env:GOARCH = $_arch
        $env:GOOS = $os

        Write-Output ("Building for '{0}'" -f $os)

        # Build up the command to create the CLI binary
        $cmd = 'go build -ldflags "-X github.com/amido/stacks-cli/cmd.version={0}" -o {1}' -f
                    $BuildNumber,
                    $cli_filename

        Invoke-Expression -Command $cmd

        # Build up the command to build the test
        $cmd = 'go test -ldflags "-X github.com/amido/stacks-cli/testing/integration.version={0}" -tags integration -o {1} -c github.com/amido/stacks-cli/testing/integration/...' -f
                    $BuildNumber,
                    $inttest_filename

        Invoke-Expression -Command $cmd

        Write-Output ("End build for '{0}': {1}" -f $os, $cli_filename)

    }
}



