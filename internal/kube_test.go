package kube

import "testing"

func TestSelfSignedCreating(t *testing.T) {
	t.Run("Test creating kube client", func(t *testing.T) {
		Client()
	})
}
