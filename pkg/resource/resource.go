package resource

import (
	"github.com/alexvanboxel/treactor-go/pkg/chem"
)

var Atoms *chem.Atoms

func Init() {
	GoogleCloudInit()
	clientInit()

	Atoms = chem.NewAtoms()
}
