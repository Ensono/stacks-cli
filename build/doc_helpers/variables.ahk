
; This file conatins the variables that all scripts use
; These are things like the company name etc so that all scripts
; appear to run in the same way

; Command := Format("{1}\..\..\outputs\bin\stacks-cli-windows-amd64-100.98.99.exe interactive", A_WorkingDir)
Command := "stacks-cli interactive"
Company := "Ensono Digital"
Scope := "core"
Pipeline := "azdo"
CloudPlatform := "azure"
TFGroup := "stacks_ancillary-ressources"
TFStorage := "kjh56sdfnjnkjn"
TFContainer := "tfstate"
DomainExternal := "example-stacks.com"

ProjectWorkingDir := "projects"

GitHubOrg := "https://github.com/ensonodigital"
