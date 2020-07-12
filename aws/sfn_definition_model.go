package aws

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
