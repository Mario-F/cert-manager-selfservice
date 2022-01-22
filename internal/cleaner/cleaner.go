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
	stop         chan bool
}

func (c *Cleaner) Start(hours int64) {
	log.Infof("Starting cleaner")
	c.cleanupHours = hours
	c.stop = make(chan bool)

	ticker := time.NewTicker(time.Minute * 30)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := c.run()
				if err != nil {
					log.Errorf("Cleanup run error: %v+", err)
				}
			case <-c.stop:
				log.Infof("Stopping cleaner")
				return
			}
		}
	}()
}

func (c *Cleaner) Stop() {
	close(c.stop)
}

func (c *Cleaner) run() error {
	c.runningMux.Lock()
	defer c.runningMux.Unlock()

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
