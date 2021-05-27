package deploy

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"deploy": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

// TODO(wperron) maybe uncomment?
// func testAccPreCheck(t *testing.T) {
// 	if v := os.Getenv("GRAFANA_URL"); v == "" {
// 		t.Fatal("GRAFANA_URL must be set for acceptance tests")
// 	}
// 	if v := os.Getenv("GRAFANA_AUTH"); v == "" {
// 		t.Fatal("GRAFANA_AUTH must be set for acceptance tests")
// 	}
// 	if v := os.Getenv("GRAFANA_ORG_ID"); v == "" {
// 		t.Fatal("GRAFANA_ORG_ID must be set for acceptance tests")
// 	}
// }
