<#

.SYNOPSIS
Compile the CLI binary and tests

#>

[CmdletBinding()]
param (

    [string]
    # Base directory for all binary builds
    $BasePath = "/eirctl/outputs/bin",

    [string]
    # Set the build number
    $BuildNumber = $env:BUILDNUMBER,

    [string[]]
    # operating system that binary should be build for
    $targets = @("windows", "linux", "darwin:arm64"),

    [string[]]
    # architecture that should be targeted
    $arch = @("amd64"),

    [switch]
    # specify if the VCS status should be built into the go binaries
    $NoVCS
)

# If the base directory does not exist, create it
if (!(Test-Path -Path $BasePath)) {
    Write-Output ("Creating output dir: {0}" -f $BasePath)
    New-Item -ItemType Directory -Path $BasePath | Out-Null
}

# Run command to get the packages required for the build
$cmd = "go get"
Invoke-Expression -Command $cmd

# Set up the go build argument for vcs
$build_vcs = ""
if ($NoVCS.IsPresent) {
    $build_vcs = "-buildvcs=false"
}

# iterate around each of the target os
foreach ($os in $targets) {

    # get a list of the architectures to build for this OS
    $archs = $arch

    # split the OS using the : character and see if there are any additional architectures
    # that should be built for that OS
    $manifest = $os -split ":"
    if ($manifest.count -eq 2) {

        $os = $manifest[0]

        # now split the second element of the manifest using a comma, to determine what extra
        # arhces should be built
        $additional = , $manifest[1] -split ","

        if ($additional) {
            $archs = $archs + $additional
        }
    }

    # iterate around each of the architectures
    foreach ($_arch in $archs) {

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

        Write-Output ("Building for '{0}' ({1})" -f $os, $_arch)

        # Build up the command to create the CLI binary
        $cmd = 'go build -ldflags "-X github.com/Ensono/stacks-cli/cmd.version={0}" -o {1} {2}' -f
        $BuildNumber,
        $cli_filename,
        $build_vcs

        Invoke-Expression -Command $cmd

        # Build up the command to build the test
        $cmd = 'go test -ldflags "-X github.com/Ensono/stacks-cli/testing/integration.version={0}" -tags integration -o {1} -c github.com/Ensono/stacks-cli/testing/integration/...' -f
        $BuildNumber,
        $inttest_filename

        Invoke-Expression -Command $cmd

        Write-Output ("End build for '{0}': {1}" -f $os, $cli_filename)

    }
}
