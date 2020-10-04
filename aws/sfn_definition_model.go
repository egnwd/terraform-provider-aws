package aws

import (
	"encoding/json"
)

type sfnStateMachinePath struct {
	s      string
	isNull bool
}

type SfnStateMachineDefinition struct {
	Version string                 `json:",omitempty"`
	Comment string                 `json:",omitempty"`
	StartAt string                 `json:""`
	Timeout int                    `json:"TimeoutSeconds,omitempty"`
	States  map[string]interface{} `json:""`
}

type SfnStateMachineState struct {
	Type       string               `json:""`
	Next       string               `json:",omitempty"`
	End        bool                 `json:",omitempty"`
	Comment    string               `json:",omitempty"`
	InputPath  *sfnStateMachinePath `json:",omitempty"`
	OutputPath *sfnStateMachinePath `json:",omitempty"`
}

type SfnStateMachinePassState struct {
	SfnStateMachineState
	Result     map[string]interface{} `json:",omitempty"`
	ResultPath *sfnStateMachinePath   `json:",omitempty"`
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
	Resource             string                    `json:""`
	Parameters           map[string]interface{}    `json:",omitempty"`
	ResultPath           *sfnStateMachinePath      `json:",omitempty"`
	Retry                []*SfnStateMachineRetrier `json:",omitempty"`
	Catch                []*SfnStateMachineCatcher `json:",omitempty"`
	TimeoutSeconds       int                       `json:",omitempty"`
	HeartbeatSeconds     int                       `json:",omitempty"`
	TimeoutSecondsPath   string                    `json:",omitempty"`
	HeartbeatSecondsPath string                    `json:",omitempty"`
}

type SfnStateMachineParallelState struct {
	SfnStateMachineState
	Branches   []*SfnStateMachineStates  `json:""`
	ResultPath *sfnStateMachinePath      `json:",omitempty"`
	Retry      []*SfnStateMachineRetrier `json:",omitempty"`
	Catch      []*SfnStateMachineCatcher `json:",omitempty"`
}

type SfnStateMachineStates struct {
	StartAt string                 `json:""`
	States  map[string]interface{} `json:""`
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

func (p sfnStateMachinePath) MarshalJSON() ([]byte, error) {
	s := p.s

	if p.isNull {
		return []byte(`null`), nil
	}

	return json.Marshal(s)
}

func sfnStateMachineDefinitionConfigStringList(lI []interface{}) []string {
	ret := make([]string, len(lI))
	for i, vI := range lI {
		ret[i] = vI.(string)
	}
	return ret
}

func sfnStateMachineDefinitionPath(s string) *sfnStateMachinePath {
	if len(s) == 0 {
		return &sfnStateMachinePath{isNull: true}
	}

	if s == "$" {
		return nil
	}

	return &sfnStateMachinePath{s: s}
}
