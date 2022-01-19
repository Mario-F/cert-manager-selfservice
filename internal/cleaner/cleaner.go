package cleaner

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

type Cleaner struct {
	runningMux sync.Mutex
}

func (c *Cleaner) Run() error {
	c.runningMux.Lock()
	defer c.runningMux.Unlock()

	log.Infof("Starting cleanup run")

	return nil
}
