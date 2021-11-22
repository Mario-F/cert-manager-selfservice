package cert

import "testing"

func TestCaCreation(t *testing.T) {
	var caCertRes *CaCert

	t.Run("Test getCa first time", func(t *testing.T) {
		caCert, err := getCA()
		if err != nil {
			t.Error(err)
		}
		if !caCert.created {
			t.Errorf("CaCert is empty")
		}
		caCertRes = caCert
	})

	t.Run("Test getCa second time", func(t *testing.T) {
		caCert, err := getCA()
		if err != nil {
			t.Error(err)
		}
		if !caCert.created {
			t.Errorf("CaCert is empty")
		}
		if caCertRes != caCert {
			t.Errorf("CaCert was not reused")
		}
	})
}
