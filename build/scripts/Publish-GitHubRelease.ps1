<#

.SYNOPSIS
Publishes a GitHub release using arguments and environment variables

.DESCRIPTION
Using the GitHub API this script will publish a release on specified repository using the
name commit id an version number.

All parameters are passed using arguments, apart from the GitHub token that is passed using an environment
variable

#>

[CmdletBinding()]
param (

    [string]
    # Version number of the release
    $version = $env:VERSION_NUMBER,

    [string]
    # Commit ID to be release
    $commitId = $env:COMMIT_ID,

    [string]
    # Release notes. This can include helpful notes about installation for example
    # that will be specific to the release
    $notes = $env:NOTES,

    [string]
    # Artifacts directory, items in this folder will be added to the release
    $artifactsDir = $env:ARTIFACTS_DIR,

    [string[]]
    # List of files that will be uploaded to the release
    $artifactsList = @(),

    [string]
    # Username to use to perform the release under
    $owner = $env:OWNER,

    [string]
    # API Key to use to authenticate to perform the release
    $apikey = $env:API_KEY,

    [string]
    # GithUB repository that the release is for
    $repository = $env:REPOSITORY,

    [bool]
    # Set if the release is a Draft, e.g. not visible to users
    $draft = $false,

    [bool]
    # Pre-release of an upcoming major release
    $preRelease = $true

)

# As environment variables cannot be easily used for the boolean values
# check to see if they have been set and overwite the values if they have
if ([string]::IsNullOrEmpty($env:DRAFT)) {
    $draft = $false
} else {
    $draft = $true
}

if ([string]::IsNullOrEmpty($env:PRERELEASE)) {
    $preRelease = $false
} else {
    $preRelease = $true
}

# if the artifactsList is empty, get all the files in the specified artifactsDir
# otherwise find the files that have been specified
if ($artifactsList.Count -eq 0) {
    $artifactsList = Get-ChildItem -Path $artifactsDir -Recurse
} else {
    $files = $artifactsList
    $artifactsList = @()

    foreach ($file in $files) {
        $artifactsList += , (Get-ChildItem -Path $artifactsDir -Recurse -Filter $file)
    }
}

# Create an object to be used as the body of the request
$requestBody = @{
    tag_name = ("v{0}" -f $version)
    target_commitsh = $commitId
    name = ("v{0}" -f $version)
    body = $notes
    draft = $draft
    prerelease = $preRelease
}

# Create the Base64encoded string for the APIKey to be used in the header of the API call
$base64key = [Convert]::ToBase64String(
    [Text.Encoding]::Ascii.GetBytes($("{0}:x-oauth-basic" -f $apikey))
)

# Now create the header
$header = @{
    Authorization = ("Basic {0}" -f $base64key)
}

# Create the splat hashtable to be used as the arguments for the Invoke-RestMethod cmdlet
$releaseArgs = @{
    Uri = ("https://api.github.com/repos/{0}/{1}/releases" -f $owner, $repository)
    Method = "POST"
    Headers = $header
    ContentType = "application/json"
    Body = (ConvertTo-JSON -InputObject $requestBody -Compress)
}

# Create the release by making the API call, artifacts will be uploaded afterwards
Write-Output ("Creating release for: {0}" -f $version)
$result = Invoke-RestMethod @releaseArgs

# Get the uploadUri that has been returned by the initial call
$uploadUri = $result | Select-Object -ExpandProperty upload_url

# Iterate around all of the artifacts that are to be uploaded
foreach ($uploadFile in $artifactsList) {

    # get the name of the artifact
    $artifact = Get-Item -Path $uploadFile

    Write-Output ("Adding asset to release: {0}" -f $artifact.Name)

    # Use the uploadUri to create a URI for the artifact
    $artifactUri = $uploadUri -replace "\{\?name,label\}", ("?name={0}" -f $artifact.Name)

    # Create the argument hash to perform the upload
    $uploadArgs = @{
        Uri = $artifactUri
        Method = "POST"
        Headers = $header
        ContentType = "application/octet-stream"
        InFile = $uploadFile
    }

    # Perform the upload of the artifact
    $result = Invoke-WebRequest @uploadArgs
}

