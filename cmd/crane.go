package main

import (
	"github.com/konveyor/tackle2-addon/command"
)

type Export struct {
	Namespace string
}

func (r *Export) Run() (err error) {
	cmd := command.New("/usr/bin/crane")
	cmd.Dir = OutputDir
	cmd.Options.Add("export")
	cmd.Options.Add("--kubeconfig", KubeConfigPath)
	cmd.Options.Add("-n", r.Namespace)
	addon.Activity("[Crane] Exporting resources from namespace %s.", r.Namespace)
	err = cmd.Run()
	return
}

type Transform struct{}

func (r *Transform) Run() (err error) {
	cmd := command.New("/usr/bin/crane")
	cmd.Dir = OutputDir
	cmd.Options.Add("transform")
	addon.Activity("[Crane] Generating transforms.")
	err = cmd.Run()
	return
}

type Apply struct{}

func (r *Apply) Run() (err error) {
	cmd := command.New("/usr/bin/crane")
	cmd.Dir = OutputDir
	cmd.Options.Add("apply")
	addon.Activity("[Crane] Applying transforms.")
	err = cmd.Run()
	return
}
