package cmd

import "testing"

var (
	initMessageOriginal = "CertManager SelfService!"
)

func TestInitMessage(t *testing.T) {
	t.Run("InitMessage valid", func(t *testing.T) {
		resMes := initMessage()
		if resMes != initMessageOriginal {
			t.Errorf("initMessage %s is not what it should be %s", resMes, initMessageOriginal)
		}
	})
}
