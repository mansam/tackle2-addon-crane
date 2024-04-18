package main

import (
	"path"

	"github.com/konveyor/tackle2-addon/command"
)

type Export struct {
	Data *Data
}

func (r *Export) Run() (err error) {
	cmd := command.New("/usr/bin/crane")
	cmd.Dir = RepoDir
	cmd.Options.Add("export")
	cmd.Options.Add("-n", r.Data.Namespace)
	addon.Activity("[Crane] Exporting resources from namespace %s.", r.Data.Namespace)
	err = cmd.Run()
	return
}

type Transform struct{}

func (r *Transform) Run() (err error) {
	cmd := command.New("/usr/bin/crane")
	cmd.Dir = RepoDir
	cmd.Options.Add("transform")
	addon.Activity("[Crane] Generating transforms.")
	err = cmd.Run()
	return
}

type Apply struct{}

func (r *Apply) Run() (err error) {
	cmd := command.New("/usr/bin/crane")
	cmd.Dir = RepoDir
	cmd.Options.Add("apply")
	addon.Activity("[Crane] Applying transforms.")
	err = cmd.Run()
	return
}

type Rm struct {
	Data *Data
}

func (r *Rm) Run() (err error) {
	cmd := command.New("rm")
	cmd.Options.Add("-rf", path.Join(RepoDir, "*"))
	err = cmd.Run()
	return
}
