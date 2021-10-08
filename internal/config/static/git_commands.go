package static

// declare a slice of git commands that are used at the end of
// provisioning the project to initialise the repository and
// configure the remote, based on the configuration given to the CLI
var GitCmds = []string{
	"git init",
	"git remote add origin {{ .Project.SourceControl.URL }}",
}
