package aws

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAWSDataSourceSfnDefinition_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSISfnStateMachineDefintionConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_sfn_state_machine_definition.test", "json",
						testAccAWSISfnStateMachineDefintionJSON,
					),
				),
			},
		},
	})
}

var testAccAWSISfnStateMachineDefintionConfig = `
data "aws_sfn_state_machine_definition" "test" {
	comment   = "Foo Bar"
	start_at  = "Baz"

	state {
		name = "Baz"
	}
}
`

var testAccAWSISfnStateMachineDefintionJSON = `{
  "Comment": "Foo Bar",
  "StartAt": "Baz",
  "Version": "1.0",
  "States": {
    "Baz": {}
  }
}`
