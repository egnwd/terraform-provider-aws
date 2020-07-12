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

func TestAccAWSDataSourceSfnDefinition_succeed(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSfnStateMachineDefinitionSucceedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_sfn_state_machine_definition.test", "json",
						testAccAWSSfnStateMachineDefinitionSucceedJSON,
					),
				),
			},
		},
	})
}

func TestAccAWSDataSourceSfnDefinition_fail(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSfnStateMachineDefinitionFailConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_sfn_state_machine_definition.test", "json",
						testAccAWSSfnStateMachineDefinitionFailJSON,
					),
				),
			},
		},
	})
}

func TestAccAWSDataSourceSfnDefinition_choice(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSfnStateMachineDefinitionChoiceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_sfn_state_machine_definition.test", "json",
						testAccAWSSfnStateMachineDefinitionChoiceJSON,
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
        result      = <<-EOF
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

var testAccAWSSfnStateMachineDefinitionSucceedConfig = `
data "aws_sfn_state_machine_definition" "test" {
    start_at  = "State1"

    succeed {
        name    = "State1"
        comment = "Yay! Success!"
    }
}
`

var testAccAWSSfnStateMachineDefinitionSucceedJSON = `{
  "StartAt": "State1",
  "Version": "1.0",
  "States": {
    "State1": {
      "Type": "Succeed",
      "Comment": "Yay! Success!"
    }
  }
}`

var testAccAWSSfnStateMachineDefinitionFailConfig = `
data "aws_sfn_state_machine_definition" "test" {
    start_at  = "State1"

    fail {
        name  = "State1"
        cause = "Invalid response."
        error = "ErrorA" 
    }
}
`

var testAccAWSSfnStateMachineDefinitionFailJSON = `{
  "StartAt": "State1",
  "Version": "1.0",
  "States": {
    "State1": {
      "Type": "Fail",
      "Cause": "Invalid response.",
      "Error": "ErrorA"
    }
  }
}`

var testAccAWSSfnStateMachineDefinitionChoiceConfig = `
data "aws_sfn_state_machine_definition" "test" {
    start_at  = "Should Fail?"

    choice {
        name    = "Should Fail?"
        default = "No"

        input_path  = "$"
        output_path = "$"

        option {
            next       = "Yes"
            comparison = <<-EOF
            {
                "Variable": "$.value",
                "NumericEquals": 0
            }
            EOF
        }

        option {
            next       = "Yes"
            comparison = <<-EOF
            {
                "Variable": "$.value",
                "NumericGreaterThanEquals": 10
            }
            EOF
        }
    }
}
`

var testAccAWSSfnStateMachineDefinitionChoiceJSON = `{
  "StartAt": "Should Fail?",
  "Version": "1.0",
  "States": {
    "Should Fail?": {
      "Type": "Choice",
      "InputPath": "$",
      "OutputPath": "$",
      "Choices": [
        {
          "Next": "Yes",
          "NumericEquals": 0,
          "Variable": "$.value"
        },
        {
          "Next": "Yes",
          "NumericGreaterThanEquals": 10,
          "Variable": "$.value"
        }
      ],
      "Default": "No"
    }
  }
}`
