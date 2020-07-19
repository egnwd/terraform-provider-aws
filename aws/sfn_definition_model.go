package aws

import (
	"encoding/json"
)

type SfnStateMachineDefinition struct {
	Comment string                 `json:",omitempty"`
	StartAt string                 `json:""`
	Timeout int                    `json:"TimeoutSeconds,omitempty"`
	Version string                 `json:",omitempty"`
	States  map[string]interface{} `json:""`
}

type SfnStateMachineState struct {
	Type       string  `json:""`
	Next       string  `json:",omitempty"`
	End        bool    `json:",omitempty"`
	Comment    string  `json:",omitempty"`
	InputPath  *string `json:""`
	OutputPath *string `json:""`
}

type SfnStateMachinePassState struct {
	SfnStateMachineState
	Result     map[string]interface{} `json:",omitempty"`
	ResultPath *string                `json:""`
	Parameters map[string]interface{} `json:",omitempty"`
}

type SfnStateMachineSucceedState struct {
	Type    string `json:""`
	Comment string `json:",omitempty"`
}

type SfnStateMachineFailState struct {
	Type    string `json:""`
	Comment string `json:",omitempty"`
	Cause   string `json:",omitempty"`
	Error   string `json:",omitempty"`
}

type SfnStateMachineChoiceState struct {
	SfnStateMachineState
	Choices []*SfnStateMachineChoiceRule `json:""`
	Default string                       `json:",omitempty"`
}

type SfnStateMachineChoiceRule struct {
	Comparison map[string]interface{} `json:"-"`
	Next       string                 `json:""`
}

type SfnStateMachineWaitState struct {
	SfnStateMachineState
	SecondsPath   *string `json:",omitempty"`
	TimestampPath *string `json:",omitempty"`
	Timestamp     *string `json:",omitempty"`
	Seconds       *int    `json:",omitempty"`
}

type SfnStateMachineTaskState struct {
	SfnStateMachineState
	Resource         string                    `json:""`
	Parameters       map[string]interface{}    `json:",omitempty"`
	ResultPath       *string                   `json:""`
	Retry            []*SfnStateMachineRetrier `json:",omitempty"`
	Catch            []*SfnStateMachineCatcher `json:",omitempty"`
	TimeoutSeconds   int                       `json:",omitempty"`
	HeartbeatSeconds int                       `json:",omitempty"`
}

type SfnStateMachineRetrier struct {
	ErrorEquals     []string `json:""`
	IntervalSeconds int      `json:""`
	MaxAttempts     int      `json:""`
	BackoffRate     float64  `json:""`
}

type SfnStateMachineCatcher struct {
	ErrorEquals []string `json:""`
	Next        string   `json:""`
}

func (cr SfnStateMachineChoiceRule) MarshalJSON() ([]byte, error) {
	type SfnStateMachineChoiceRule_ SfnStateMachineChoiceRule // prevent recursion
	b, err := json.Marshal(SfnStateMachineChoiceRule_(cr))
	if err != nil {
		return nil, err
	}

	var m map[string]json.RawMessage
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	for k, v := range cr.Comparison {
		b, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
		m[k] = b
	}

	return json.Marshal(m)
}

func sfnStateMachineDefinitionConfigStringList(lI []interface{}) []string {
	ret := make([]string, len(lI))
	for i, vI := range lI {
		ret[i] = vI.(string)
	}
	return ret
}
