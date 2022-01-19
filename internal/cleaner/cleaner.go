package cleaner

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Cleaner struct {
	runningMux sync.Mutex
}

func (c *Cleaner) Run(ctx context.Context) error {
	c.runningMux.Lock()
	defer c.runningMux.Unlock()

	log.Infof("Starting cleanup run")

	return nil
}
