package controller

import (
	"github.com/foldy-project/foldy/controller/pkg/controller/experiment"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, experiment.Add)
}
