package cleaner

import (
	"sync"

	"github.com/Mario-F/cert-manager-selfservice/internal/kube"
	log "github.com/sirupsen/logrus"
)

type Cleaner struct {
	runningMux sync.Mutex
}

func (c *Cleaner) Run() error {
	c.runningMux.Lock()
	defer c.runningMux.Unlock()

	log.Infof("Starting cleanup run")

	certs, err := kube.GetCertificates()
	if err != nil {
		return err
	}
	log.Debugf("Found %d certificates managed by this instance", len(certs.Items))

	return nil
}
