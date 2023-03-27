
#SingleInstance Force

; Include the variables in the script
#Include variables.ahk

; Create an array of the inputs that need to be sent
inputs := [
    Company,
    Scope,
    "backend",
    "",
    CloudPlatform,
    TFGroup,
    TFStorage,
    TFContainer,
    DomainExternal,
    "",
    "Command Log",
    ProjectWorkingDir,
    "1",
    "my-webapi",
    "java",
    "cqrs",
    "",
    "",
    "github",
    GitHubOrg . "/my-webapi",
    "westeurope",
    "mywebapi-resources"
]

+^j::
{
    ; Clear the screen beforehand so that a screenshot can be taken afterwards
    ; Send "clear{Enter}"

    ; Run the command and then wait for it to be ready
    Send Format("{1}{Enter}", Command)
    Sleep 5000

    ; Iterate around the inputs and send each one
    for input in inputs
        Send Format("{1}{Enter}", input)
        Sleep 10

}