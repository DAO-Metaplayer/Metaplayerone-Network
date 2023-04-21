package main

import (
	_ "embed"

	"github.com/vishnushankarsg/metad/command/root"
	"github.com/vishnushankarsg/metad/licenses"
)

var (
	//go:embed LICENSE
	license string
)

func main() {
	licenses.SetLicense(license)

	root.NewRootCommand().Execute()
}
