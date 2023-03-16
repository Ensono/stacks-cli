package config

import (
	_ "embed"

	"github.com/amido/stacks-cli/internal/config/staticFiles"
)

type Internal struct {
	Files map[string]interface{}
}

type internalFileInfo struct {
	filename string
	data     []byte
}

func (i *Internal) AddFiles() {

	// initialise the map
	i.Files = make(map[string]interface{})

	// Add the files that are required internally
	i.Files["banner"] = internalFileInfo{
		filename: "banner.txt",
		data:     []byte(staticFiles.IntFile_Banner),
	}

	i.Files["config"] = internalFileInfo{
		filename: "internal_config.yml",
		data:     []byte(staticFiles.IntFile_config),
	}

	i.Files["azdo"] = internalFileInfo{
		filename: "azdo_variable_template.yml",
		data:     []byte(staticFiles.Azdo_Variable_Template_Tmpl),
	}

}

func (i *Internal) GetFileContentString(name string) string {
	return string(i.GetFileContent(name))
}

func (i *Internal) GetFileContent(name string) []byte {
	return i.Files[name].(internalFileInfo).data
}

func (i *Internal) GetFilename(name string) string {
	return i.Files[name].(internalFileInfo).filename
}
