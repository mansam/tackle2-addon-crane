package main

import (
	"errors"
	"os"
	"path"

	liberr "github.com/jortel/go-utils/error"
	"github.com/konveyor/tackle2-hub/api"
	"github.com/konveyor/tackle2-hub/nas"
)

type Kubernetes struct {
	Identities []api.Ref
}

func (r *Kubernetes) findKubeConfig() (matched *api.Identity, found bool, err error) {
	for _, ref := range r.Identities {
		identity, nErr := addon.Identity.Get(ref.ID)
		if nErr != nil {
			err = nErr
			return
		}
		if identity.Kind == "kubeconfig" {
			found = true
			matched = identity
			break
		}
	}
	return
}

func (r *Kubernetes) WriteKubeConfig() (err error) {
	id, found, err := r.findKubeConfig()
	if err != nil {
		return
	}
	if found {
		addon.Activity(
			"[Kubernetes] Using credentials (id=%d) %s.",
			id.ID,
			id.Name)
	} else {
		err = errors.New("could not find kubeconfig")
		return
	}

	p := path.Join(HomeDir, ".kubeconfig")
	found, err = nas.Exists(p)
	if found || err != nil {
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

	_, err = f.Write([]byte(id.Settings + "\n"))
	if err != nil {
		err = liberr.Wrap(
			err,
			"path",
			p)
		return
	}

	_ = f.Close()
	addon.Activity("[FILE] Created %s.", p)
	return
}
