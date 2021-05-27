package deploy

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/wperron/terraform-deploy-provider/client"
)

const testAccProjectConfig_basic = `
resource "deploy_project" "test" {
  name       = "terraform-test"
  source_url = "https://dash.deno.com/examples/hello.js"
}
`

const testAccProjectConfig_update = `
resource "grafana_user" "test" {
  name       = "terraform-test"
  source_url = "https://dash.deno.com/examples/hello.js"
}
`

func TestAccUser_basic(t *testing.T) {
	var project client.Project
	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    testAccProviders,
		CheckDestroy: testAccProjectCheckDestroy(&project),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists("grafana_user.test", &project),
					resource.TestCheckResourceAttr(
						"grafana_user.test", "email", "terraform-test@localhost",
					),
					resource.TestCheckResourceAttr(
						"grafana_user.test", "name", "Terraform Test",
					),
					resource.TestCheckResourceAttr(
						"grafana_user.test", "login", "tt",
					),
					resource.TestCheckResourceAttr(
						"grafana_user.test", "password", "abc123",
					),
					resource.TestMatchResourceAttr(
						"grafana_user.test", "id", regexp.MustCompile(`\d+`),
					),
				),
			},
			{
				Config: testAccProjectConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists("grafana_user.test", &project),
					resource.TestCheckResourceAttr(
						"grafana_user.test", "email", "terraform-test-update@localhost",
					),
					resource.TestCheckResourceAttr(
						"grafana_user.test", "name", "Terraform Test Update",
					),
					resource.TestCheckResourceAttr(
						"grafana_user.test", "login", "ttu",
					),
					resource.TestCheckResourceAttr(
						"grafana_user.test", "password", "zyx987",
					),
					resource.TestCheckResourceAttr(
						"grafana_user.test", "is_admin", "true",
					),
				),
			},
		},
	})
}

func testAccProjectCheckExists(rn string, a *client.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// TODO(wperron) wat?
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}
		tmp, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		id := int64(tmp)
		if err != nil {
			return fmt.Errorf("resource id is malformed")
		}
		client := testAccProvider.Meta().(*client.Client)
		project, err := client.GetProject(id)
		if err != nil {
			return fmt.Errorf("error getting data source: %s", err)
		}
		*a = project
		return nil
	}
}

func testAccProjectCheckDestroy(a *client.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		project, err := client.GetProject(a.ID)
		if err == nil && project.Name != "" {
			return fmt.Errorf("project still exists")
		}
		return nil
	}
}
