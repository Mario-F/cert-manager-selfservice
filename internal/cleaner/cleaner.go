package cleaner

import (
	"sync"
	"time"

	"github.com/Mario-F/cert-manager-selfservice/internal/kube"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

var (
	promCleanupCalledCounter prometheus.Counter
)

type Cleaner struct {
	runningMux   sync.Mutex
	cleanupHours int64
	stop         chan bool
	isStarted    bool
}

func init() {
	promCleanupCalledCounter = promauto.NewCounter(prometheus.CounterOpts{Name: "cms_cleanup_total", Help: "Count of cleanup routines executed"})
}

func (c *Cleaner) Start(hours int64) error {
	log.Infof("Starting cleaner")
	c.cleanupHours = hours

	// exec a initial cleanup to return on error early
	err := c.run()
	if err != nil {
		return err
	}

	c.stop = make(chan bool)
	c.isStarted = true

	// TODO: Add configurable interval in go time format
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

	return nil
}

func (c *Cleaner) Stop() {
	if c.isStarted {
		close(c.stop)
	}
}

func (c *Cleaner) run() error {
	c.runningMux.Lock()
	defer c.runningMux.Unlock()

	promCleanupCalledCounter.Inc()
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
				// TODO: Add prom metrics (certificate delete error)
			}
			// TODO: Add prom metrics (certificate expired count)
		}
	}

	return nil
}
