package provider

import (
	"testing"
)

func TestAccProvider(t *testing.T) {
	// This test simply verifies that the provider can be instantiated
	// without errors. More comprehensive tests would require API access.
	provider := New("test")()
	if provider == nil {
		t.Fatal("Expected provider to be instantiated")
	}
}