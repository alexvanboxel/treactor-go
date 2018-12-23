package resource

import (
	"github.com/alexvanboxel/reactor/pkg/chem"
)

var Atoms *chem.Atoms

func Init() {
	GoogleCloudInit()
	clientInit()

	Atoms = chem.NewAtoms()
}
