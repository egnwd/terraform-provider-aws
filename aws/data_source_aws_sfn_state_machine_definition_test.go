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
    start_at  = "State1"

    state {
        name    = "State1"
        type    = "Pass"
        comment = "Doesn't do anything"

        input_path  = "$"
        output_path = ""
        next        = "State2"
    }

    state {
        name    = "State2"
        type    = "Pass"

        input_path  = "$.data"
        output_path = "$"
        end         = true
    }
}
`

var testAccAWSISfnStateMachineDefintionJSON = `{
  "Comment": "Foo Bar",
  "StartAt": "State1",
  "Version": "1.0",
  "States": {
    "State1": {
      "Type": "Pass",
      "Next": "State2",
      "Comment": "Doesn't do anything",
      "InputPath": "$",
      "OutputPath": null
    },
    "State2": {
      "Type": "Pass",
      "End": true,
      "InputPath": "$.data",
      "OutputPath": "$"
    }
  }
}`
