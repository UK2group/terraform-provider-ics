package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"ics": providerserver.NewProtocol6WithError(New("test")()),
}

func TestAccProvider(t *testing.T) {
	// This test simply verifies that the provider can be instantiated
	// without errors. More comprehensive tests would require API access.
	provider := New("test")()
	if provider == nil {
		t.Fatal("Expected provider to be instantiated")
	}
}