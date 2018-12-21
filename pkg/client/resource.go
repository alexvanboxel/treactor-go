package client

import (
	"github.com/alexvanboxel/reactor/pkg/chem"
)

var Atoms *chem.Atoms

func Init() {
	GoogleCloudInit()
	ClientInit()

	Atoms = chem.NewAtoms()
}
