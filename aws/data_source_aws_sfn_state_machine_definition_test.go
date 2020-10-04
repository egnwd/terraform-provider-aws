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

func TestAccAWSDataSourceSfnDefinition_wait(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSfnStateMachineDefinitionWaitConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_sfn_state_machine_definition.test", "json",
						testAccAWSSfnStateMachineDefinitionWaitJSON,
					),
				),
			},
		},
	})
}

func TestAccAWSDataSourceSfnDefinition_task(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSfnStateMachineDefinitionTaskConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_sfn_state_machine_definition.test", "json",
						testAccAWSSfnStateMachineDefinitionTaskJSON,
					),
				),
			},
		},
	})
}

func TestAccAWSDataSourceSfnDefinition_taskPath(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSfnStateMachineDefinitionTaskPathConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_sfn_state_machine_definition.test", "json",
						testAccAWSSfnStateMachineDefinitionTaskPathJSON,
					),
				),
			},
		},
	})
}

func TestAccAWSDataSourceSfnDefinition_parallel(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSfnStateMachineDefinitionParallelConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_sfn_state_machine_definition.test", "json",
						testAccAWSSfnStateMachineDefinitionParallelJSON,
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

        output_path = ""
        next        = "State2"
    }

    pass {
        name    = "State2"

        input_path  = "$.data"
        end         = true
    }
}
`

var testAccAWSSfnStateMachineDefinitionCommonJSON = `{
  "Version": "1.0",
  "Comment": "Foo Bar",
  "StartAt": "State1",
  "States": {
    "State1": {
      "Type": "Pass",
      "Next": "State2",
      "Comment": "Doesn't do anything",
      "OutputPath": null
    },
    "State2": {
      "Type": "Pass",
      "End": true,
      "InputPath": "$.data"
    }
  }
}`

var testAccAWSSfnStateMachineDefinitionPassConfig = `
data "aws_sfn_state_machine_definition" "test" {
    comment   = "Foo Bar"
    start_at  = "State1"

    pass {
        name    = "State1"

        output_path = ""
        end         = true

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
  "Version": "1.0",
  "Comment": "Foo Bar",
  "StartAt": "State1",
  "States": {
    "State1": {
      "Type": "Pass",
      "End": true,
      "OutputPath": null,
      "Result": {
        "a": 123,
        "foo": {
          "b": [
            true,
            false
          ]
        }
      }
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
  "Version": "1.0",
  "StartAt": "State1",
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
  "Version": "1.0",
  "StartAt": "State1",
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
  "Version": "1.0",
  "StartAt": "Should Fail?",
  "States": {
    "Should Fail?": {
      "Type": "Choice",
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

var testAccAWSSfnStateMachineDefinitionWaitConfig = `
data "aws_sfn_state_machine_definition" "test" {
    start_at  = "State1"

    wait {
        name         = "State1"
        seconds_path = "$.seconds"

        next        = "State2"
    }

    wait {
        name           = "State2"
        timestamp_path = "$.timestamp"

        next        = "State3"
    }

    wait {
        name      = "State3"
        timestamp = "2016-08-18T17:33:00Z"

        next        = "State4"
    }

    wait {
        name    = "State4"
        seconds = 4

        next        = "State5"
    }

    wait {
        name    = "State5"
        seconds = 0

        end         = true
    }
}
`

var testAccAWSSfnStateMachineDefinitionWaitJSON = `{
  "Version": "1.0",
  "StartAt": "State1",
  "States": {
    "State1": {
      "Type": "Wait",
      "Next": "State2",
      "SecondsPath": "$.seconds"
    },
    "State2": {
      "Type": "Wait",
      "Next": "State3",
      "TimestampPath": "$.timestamp"
    },
    "State3": {
      "Type": "Wait",
      "Next": "State4",
      "Timestamp": "2016-08-18T17:33:00Z"
    },
    "State4": {
      "Type": "Wait",
      "Next": "State5",
      "Seconds": 4
    },
    "State5": {
      "Type": "Wait",
      "End": true,
      "Seconds": 0
    }
  }
}`

var testAccAWSSfnStateMachineDefinitionTaskConfig = `
data "aws_sfn_state_machine_definition" "test" {
    start_at  = "State1"

    task {
      name        = "State1"
      output_path = ""
      next        = "State2"

      resource   = "arn:aws:states:::batch:submitJob.sync"
      parameters = <<-EOF
      {
        "JobDefinition": "preprocessing",
        "JobName": "PreprocessingBatchJob",
        "JobQueue": "SecondaryQueue",
        "Parameters.$": "$.batchjob.parameters",
        "RetryStrategy": {
          "attempts": 5
        }
      }
      EOF

      retry {
        errors       = ["States.Timeout", "ErrorA"]
        interval     = 1
        max_attempts = 2
        backoff      = 2.0 
      }

      retry {
        errors       = ["States.ALL"]
        interval     = 5
        max_attempts = 1
        backoff      = 1.0 
      }
    }

    task {
      name        = "State2"
      result_path = ""
      end         = true
      
      resource = "arn:aws:lambda:us-east-1:123456789012:function:HelloWorld"

      catch {
        errors = ["States.ALL"]
        next   = "Failed"
      }

      timeout   = 4
      heartbeat = 1
    }
}
`

var testAccAWSSfnStateMachineDefinitionTaskJSON = `{
  "Version": "1.0",
  "StartAt": "State1",
  "States": {
    "State1": {
      "Type": "Task",
      "Next": "State2",
      "OutputPath": null,
      "Resource": "arn:aws:states:::batch:submitJob.sync",
      "Parameters": {
        "JobDefinition": "preprocessing",
        "JobName": "PreprocessingBatchJob",
        "JobQueue": "SecondaryQueue",
        "Parameters.$": "$.batchjob.parameters",
        "RetryStrategy": {
          "attempts": 5
        }
      },
      "Retry": [
        {
          "ErrorEquals": [
            "States.Timeout",
            "ErrorA"
          ],
          "IntervalSeconds": 1,
          "MaxAttempts": 2,
          "BackoffRate": 2
        },
        {
          "ErrorEquals": [
            "States.ALL"
          ],
          "IntervalSeconds": 5,
          "MaxAttempts": 1,
          "BackoffRate": 1
        }
      ]
    },
    "State2": {
      "Type": "Task",
      "End": true,
      "Resource": "arn:aws:lambda:us-east-1:123456789012:function:HelloWorld",
      "ResultPath": null,
      "Catch": [
        {
          "ErrorEquals": [
            "States.ALL"
          ],
          "Next": "Failed"
        }
      ],
      "TimeoutSeconds": 4,
      "HeartbeatSeconds": 1
    }
  }
}`

var testAccAWSSfnStateMachineDefinitionTaskPathConfig = `
data "aws_sfn_state_machine_definition" "test" {
    start_at  = "State1"

    task {
      name        = "State1"
      output_path = ""
      end        = true

      resource   = "arn:aws:states:::batch:submitJob.sync"

      timeout_path   = "$.timeout"
      heartbeat_path = "$.heartbeat"
    }
}
`

var testAccAWSSfnStateMachineDefinitionTaskPathJSON = `{
  "Version": "1.0",
  "StartAt": "State1",
  "States": {
    "State1": {
      "Type": "Task",
      "End": true,
      "OutputPath": null,
      "Resource": "arn:aws:states:::batch:submitJob.sync",
      "TimeoutSecondsPath": "$.timeout",
      "HeartbeatSecondsPath": "$.heartbeat"
    }
  }
}`

var testAccAWSSfnStateMachineDefinitionParallelConfig = `
data "aws_sfn_state_machine_definition" "test" {
    start_at  = "State1"

    parallel {
      name        = "State1"
      output_path = ""
      end         = true

      branch {
        start_at  = "SubState1a"

        pass {
          name    = "SubState1a"
  
          next        = "SubState1b"
        }

        task {
          name        = "SubState1b"
          end         = true
          
          resource = "arn:aws:lambda:us-east-1:123456789012:function:HelloWorld"
        }
      }

      branch {
        start_at  = "SubState2"

        parallel {
          name    = "SubState2"
          output_path = ""
          end         = true

          branch {
            start_at = "SubState2a"

            succeed {
              name = "SubState2a"
            }
          }

          branch {
            start_at = "SubState2a"

            fail {
              name = "SubState2a"
            }
          }
        }
      }
    }
}
`

var testAccAWSSfnStateMachineDefinitionParallelJSON = `{
  "Version": "1.0",
  "StartAt": "State1",
  "States": {
    "State1": {
      "Type": "Parallel",
      "End": true,
      "OutputPath": null,
      "Branches": [
        {
          "StartAt": "SubState1a",
          "States": {
            "SubState1a": {
              "Type": "Pass",
              "Next": "SubState1b"
            },
            "SubState1b": {
              "Type": "Task",
              "End": true,
              "Resource": "arn:aws:lambda:us-east-1:123456789012:function:HelloWorld"
            }
          }
        },
        {
          "StartAt": "SubState2",
          "States": {
            "SubState2": {
              "Type": "Parallel",
              "End": true,
              "OutputPath": null,
              "Branches": [
                {
                  "StartAt": "SubState2a",
                  "States": {
                    "SubState2a": {
                      "Type": "Succeed"
                    }
                  }
                },
                {
                  "StartAt": "SubState2a",
                  "States": {
                    "SubState2a": {
                      "Type": "Fail"
                    }
                  }
                }
              ]
            }
          }
        }
      ]
    }
  }
}`
