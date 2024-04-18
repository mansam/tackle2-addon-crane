package main

import (
	"errors"
	"os"
	"path"

	"github.com/konveyor/tackle2-addon/repository"
	"github.com/konveyor/tackle2-addon/ssh"
	hub "github.com/konveyor/tackle2-hub/addon"
)

var (
	addon   = hub.Addon
	Dir     = ""
	RepoDir = ""
	HomeDir = ""
)

type Data struct {
	Namespace string
}

func init() {
	Dir, _ = os.Getwd()
	RepoDir = path.Join(Dir, "repository")
	HomeDir, _ = os.UserHomeDir()
}

func main() {
	addon.Run(func() (err error) {
		d := &Data{}
		err = addon.DataWith(d)
		if err != nil {
			return
		}
		//
		// Fetch application.
		addon.Activity("Fetching application.")
		application, err := addon.Task.Application()
		if err != nil {
			return
		}

		if application.Repository == nil {
			err = errors.New("application repository not defined")
			return
		}

		//
		// SSH
		agent := ssh.Agent{}
		err = agent.Start()
		if err != nil {
			return
		}

		rp, err := repository.New(RepoDir, application.Repository, application.Identities)
		if err != nil {
			return
		}
		err = rp.Fetch()
		if err != nil {
			return
		}
		err = rp.Branch("crane")
		if err != nil {
			return
		}
		rm := Rm{}
		err = rm.Run()
		if err != nil {
			return
		}

		kube := Kubernetes{Identities: application.Identities}
		err = kube.WriteKubeConfig()
		if err != nil {
			return
		}

		export := Export{Data: d}
		err = export.Run()
		if err != nil {
			return
		}

		transform := Transform{}
		err = transform.Run()
		if err != nil {
			return
		}

		apply := Apply{}
		err = apply.Run()
		if err != nil {
			return
		}

		err = rp.Commit([]string{"."}, "transformed by crane")
		if err != nil {
			return
		}

		println(application)
		return
	})
}
