package installer

import (
	"sync"
	"sync/atomic"
)

type ApplicationComponent struct {
	Name          string
	RepoURL       string
	Path          string
	Dependencies  []string
	CRDs          []string
	ExtraRepos    []*Repository
	PreInstall    func(s *Installer) error
	PostInstall   func(s *Installer) error
	PreUninstall  func(s *Installer) error
	PostUninstall func(s *Installer) error
	done          <-chan error
	isHandled     int32
	l             sync.Mutex
}

func (c *ApplicationComponent) init() {
	c.done = make(chan error, 1)
	hasArgoCDDep := false
	for _, dep := range c.Dependencies {
		if dep == "argocd" {
			hasArgoCDDep = true
			break
		}
	}
	if !hasArgoCDDep {
		c.Dependencies = append([]string{"argocd"}, c.Dependencies...)
	}
}

func (c *ApplicationComponent) IsHandled() bool {
	return atomic.LoadInt32(&c.isHandled) == 1
}

func (c *ApplicationComponent) GetName() string {
	return c.Name
}

func (c *ApplicationComponent) GetDependencies() []string {
	return c.Dependencies
}

func (c *ApplicationComponent) GetCRDs() []string {
	return c.CRDs
}

func (c *ApplicationComponent) reuse() {
	atomic.StoreInt32(&c.isHandled, 0)
}

func (c *ApplicationComponent) Lock() {
	c.l.Lock()
}

func (c *ApplicationComponent) Unlock() {
	c.l.Unlock()
}

func (c *ApplicationComponent) RunInstall(s *Installer) error {
	if atomic.SwapInt32(&c.isHandled, 1) == 1 {
		return nil
	}
	if c.PreInstall != nil {
		if err := c.PreInstall(s); err != nil {
			return err
		}
	}
	//for _, repo := range c.ExtraRepos {
	//	if err := s.RunCommandInArgoCDServer("argocd repo add %s %s", repo.Name, repo.URL); err != nil {
	//		return err
	//	}
	//}
	if err := s.CreateApplication(c.Name, c.RepoURL, c.Path); err != nil {
		return err
	}
	if c.PostInstall != nil {
		if err := c.PostInstall(s); err != nil {
			return err
		}
	}
	return nil
}

func (c *ApplicationComponent) RunUninstall(s *Installer) error {
	if atomic.SwapInt32(&c.isHandled, 1) == 1 {
		return nil
	}
	if c.PreUninstall != nil {
		if err := c.PreUninstall(s); err != nil {
			return err
		}
	}
	if err := s.DeleteApplication(c.Name); err != nil {
		return err
	}
	if c.PostUninstall != nil {
		if err := c.PostUninstall(s); err != nil {
			return err
		}
	}
	return nil
}
