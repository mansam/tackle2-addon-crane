package main

import (
	"os"
	"path"

	liberr "github.com/jortel/go-utils/error"
	hub "github.com/konveyor/tackle2-hub/addon"
	"github.com/konveyor/tackle2-hub/nas"
)

var (
	addon          = hub.Addon
	Dir            = ""
	OutputDir      = ""
	HomeDir        = ""
	KubeConfigPath = ""
)

type Data struct {
	Namespace string
}

func init() {
	Dir, _ = os.Getwd()
	OutputDir = path.Join(Dir, "output")
	HomeDir, _ = os.UserHomeDir()
	KubeConfigPath = path.Join(HomeDir, ".kubeconfig")
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

		if application.Deployment == nil {
			err = liberr.New("Application does not have a cluster deployment.")
			return
		}
		deployment, err := addon.Deployment.Get(application.Deployment.ID)
		if err != nil {
			return
		}
		identity, err := addon.Identity.Get(deployment.Identity.ID)
		if err != nil {
			return
		}
		platform, err := addon.Platform.Get(deployment.Platform.ID)
		if err != nil {
			return
		}

		c := NewClusterDeployment(deployment, identity, platform)
		err = c.WriteKubeConfig()
		if err != nil {
			return
		}

		addon.Activity("Creating output directory (%s)", OutputDir)
		err = nas.MkDir(OutputDir, 0755)
		if err != nil {
			return
		}

		addon.Activity("Prepare to run crane.")
		export := Export{Namespace: deployment.Locator}
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

		addon.Activity("Uploading assets to application bucket.")
		bucket := addon.Application.Bucket(application.ID)
		err = bucket.Put(OutputDir, "crane")
		if err != nil {
			return
		}

		return
	})
}
