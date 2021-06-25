package controller

import (
	"github.com/openshift/hive-health-operator/pkg/controller/clustersync"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, clustersync.Add)
}
