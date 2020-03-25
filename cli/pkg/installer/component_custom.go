package installer

import (
	"sync"
	"sync/atomic"
)

type CustomComponent struct {
	Name         string
	Dependencies []string
	CRDs         []string
	Install      func(s *Installer) error
	Uninstall    func(s *Installer) error
	done         <-chan error
	isHandled    int32
	l            sync.Mutex
}

func (c *CustomComponent) IsHandled() bool {
	return atomic.LoadInt32(&c.isHandled) == 1
}

func (c *CustomComponent) GetName() string {
	return c.Name
}

func (c *CustomComponent) GetDependencies() []string {
	return c.Dependencies
}

func (c *CustomComponent) GetCRDs() []string {
	return c.CRDs
}

func (c *CustomComponent) Lock() {
	c.l.Lock()
}

func (c *CustomComponent) Unlock() {
	c.l.Unlock()
}

func (c *CustomComponent) RunInstall(s *Installer) error {
	if atomic.SwapInt32(&c.isHandled, 1) == 1 {
		return nil
	}
	return c.Install(s)
}

func (c *CustomComponent) RunUninstall(s *Installer) error {
	if atomic.SwapInt32(&c.isHandled, 1) == 1 {
		return nil
	}
	return c.Uninstall(s)
}

func (c *CustomComponent) init() {
	if c.Install == nil {
		panic("CustomComponent has no Install method")
	}
	if c.Uninstall == nil {
		panic("CustomComponent has no Uninstall method")
	}
	c.done = make(chan error, 1)
}

func (c *CustomComponent) reuse() {
	atomic.StoreInt32(&c.isHandled, 0)
}
