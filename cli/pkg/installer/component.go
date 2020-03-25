package installer

import (
	"fmt"
	"sync"
)

type Repository struct {
	Name string
	URL  string
}

type Component interface {
	GetName() string
	GetDependencies() []string
	GetCRDs() []string
	IsHandled() bool
	RunInstall(s *Installer) error
	RunUninstall(s *Installer) error

	init()
	reuse()
	sync.Locker
}

var componentsL sync.Mutex
var components []Component

func AddComponent(comp Component) {
	comp.init()
	componentsL.Lock()
	defer componentsL.Unlock()
	components = append(components, comp)
}

func GetComponentByName(name string) Component {
	for _, comp := range components {
		if comp.GetName() == name {
			return comp
		}
	}
	return nil
}

func GetComponentsByName(names []string) ([]Component, error) {
	filtered := make([]Component, len(names), len(names))
	for i, name := range names {
		comp := GetComponentByName(name)
		if comp == nil {
			return nil, fmt.Errorf("unknown component '%s'", name)
		}
		filtered[i] = comp
	}
	return filtered, nil
}

func GetDirectDependees(name string) []string {
	var dependees []string
	for _, comp := range components {
		if comp.GetName() == name {
			// No need to scan self
			continue
		}
		if len(comp.GetDependencies()) > 0 {
			for _, dependency := range comp.GetDependencies() {
				if dependency == name {
					dependees = append(dependees, comp.GetName())
					break
				}
			}
		}
	}
	return dependees
}
