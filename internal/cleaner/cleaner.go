package cleaner

import (
	"sync"
	"time"

	"github.com/Mario-F/cert-manager-selfservice/internal/kube"
	log "github.com/sirupsen/logrus"
)

type Cleaner struct {
	runningMux   sync.Mutex
	cleanupHours int64
}

func (c *Cleaner) Run(hours int64) error {
	c.runningMux.Lock()
	defer c.runningMux.Unlock()
	c.cleanupHours = hours

	log.Infof("Starting cleanup run, delete older than %d hours", c.cleanupHours)

	certs, err := kube.GetCertificates()
	if err != nil {
		return err
	}
	log.Debugf("Found %d certificates managed by this instance", len(certs))

	expireTime := time.Now().Unix() - (c.cleanupHours * 3600)
	for _, cert := range certs {
		if cert.LastAccess < expireTime {
			log.Infof("Certificate %s is expired, deleting", cert.Certificate.Name)
			err := cert.Delete()
			if err != nil {
				log.Errorf("Error deleting certificate %s: %v+", cert.Certificate.Name, err)
			}
		}
	}

	return nil
}
