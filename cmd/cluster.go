package main

import (
	"os"
	"path"

	liberr "github.com/jortel/go-utils/error"
	"github.com/konveyor/tackle2-addon/command"
	"github.com/konveyor/tackle2-hub/api"
	"github.com/konveyor/tackle2-hub/nas"
)

const EKS = "eks"
const Kubernetes = "kubernetes"
const OpenShift = "openshift"

func NewClusterDeployment(d *api.Deployment, i *api.Identity, p *api.Platform) (c ClusterDeployment) {
	switch p.Kind {
	case EKS:
		c = &EKSDeployment{
			Deployment: d,
			Identity:   i,
			Platform:   p,
		}
	case Kubernetes:
		c = &KubernetesDeployment{
			Deployment: d,
			Identity:   i,
			Platform:   p,
		}
	case OpenShift:
	default:
	}
	return
}

type ClusterDeployment interface {
	WriteKubeConfig() (err error)
}

type EKSDeployment struct {
	Deployment *api.Deployment
	Identity   *api.Identity
	Platform   *api.Platform
}

func (r *EKSDeployment) WriteKubeConfig() (err error) {
	err = r.writeCredentials()
	if err != nil {
		return
	}

	cmd := command.New("/usr/bin/eksctl")
	cmd.Options.Add("utils", "write-kubeconfig")
	cmd.Options.Add("--cluster", r.Platform.Name)
	cmd.Options.Add("--region", r.Platform.Region)
	cmd.Options.Add("--kubeconfig", KubeConfigPath)
	addon.Activity("[EKS] Generating EKS kubeconfig.")
	err = cmd.Run()

	return
}

func (r *EKSDeployment) writeCredentials() (err error) {
	addon.Activity("[EKS] Writing AWS credentials.")
	p := path.Join(HomeDir, ".aws", "credentials")
	found, err := nas.Exists(p)
	if found || err != nil {
		return
	}
	err = nas.MkDir(path.Join(HomeDir, ".aws"), 0755)
	if err != nil {
		return
	}
	f, err := os.Create(p)
	if err != nil {
		err = liberr.Wrap(
			err,
			"path",
			p)
		return
	}
	_, err = f.Write([]byte(r.Identity.Settings + "\n"))
	if err != nil {
		err = liberr.Wrap(
			err,
			"path",
			p)
		return
	}
	return
}

type KubernetesDeployment struct {
	Deployment *api.Deployment
	Identity   *api.Identity
	Platform   *api.Platform
}

func (r *KubernetesDeployment) WriteKubeConfig() (err error) {
	found, err := nas.Exists(KubeConfigPath)
	if found || err != nil {
		return
	}
	f, err := os.Create(KubeConfigPath)
	if err != nil {
		err = liberr.Wrap(
			err,
			"path",
			KubeConfigPath)
		return
	}

	_, err = f.Write([]byte(r.Identity.Settings + "\n"))
	if err != nil {
		err = liberr.Wrap(
			err,
			"path",
			KubeConfigPath)
		return
	}

	_ = f.Close()
	addon.Activity("[FILE] Created %s.", KubeConfigPath)
	return
}
