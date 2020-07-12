package aws

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAWSDataSourceSfnDefinition_common(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSfnStateMachineDefinitionCommonConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_sfn_state_machine_definition.test", "json",
						testAccAWSSfnStateMachineDefinitionCommonJSON,
					),
				),
			},
		},
	})
}

func TestAccAWSDataSourceSfnDefinition_pass(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSfnStateMachineDefinitionPassConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_sfn_state_machine_definition.test", "json",
						testAccAWSSfnStateMachineDefinitionPassJSON,
					),
				),
			},
		},
	})
}

var testAccAWSSfnStateMachineDefinitionCommonConfig = `
data "aws_sfn_state_machine_definition" "test" {
    comment   = "Foo Bar"
    start_at  = "State1"

    pass {
        name    = "State1"
        comment = "Doesn't do anything"

        input_path  = "$"
        output_path = ""
        result_path = "$"
        next        = "State2"
    }

    pass {
        name    = "State2"

        input_path  = "$.data"
        output_path = "$"
        result_path = "$"
        end         = true
    }
}
`

var testAccAWSSfnStateMachineDefinitionCommonJSON = `{
  "Comment": "Foo Bar",
  "StartAt": "State1",
  "Version": "1.0",
  "States": {
    "State1": {
      "Type": "Pass",
      "Next": "State2",
      "Comment": "Doesn't do anything",
      "InputPath": "$",
      "OutputPath": null,
      "ResultPath": "$"
    },
    "State2": {
      "Type": "Pass",
      "End": true,
      "InputPath": "$.data",
      "OutputPath": "$",
      "ResultPath": "$"
    }
  }
}`

var testAccAWSSfnStateMachineDefinitionPassConfig = `
data "aws_sfn_state_machine_definition" "test" {
    comment   = "Foo Bar"
    start_at  = "State1"

    pass {
        name    = "State1"

        input_path  = "$"
        output_path = ""
        end         = true

        result_path = "$"
        result      = <<EOF
{
  "a": 123,
  "foo": {
    "b": [true, false]
  }
}
EOF
    }
}
`

var testAccAWSSfnStateMachineDefinitionPassJSON = `{
  "Comment": "Foo Bar",
  "StartAt": "State1",
  "Version": "1.0",
  "States": {
    "State1": {
      "Type": "Pass",
      "End": true,
      "InputPath": "$",
      "OutputPath": null,
      "Result": {
        "a": 123,
        "foo": {
          "b": [
            true,
            false
          ]
        }
      },
      "ResultPath": "$"
    }
  }
}`
