package controller

import (
	"github.com/foldy-project/foldy/controller/pkg/controller/backend"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, backend.Add)
}
