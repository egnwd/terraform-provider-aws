package aws

type SfnStateMachineDefinition struct {
	Comment string                           `json:",omitempty"`
	StartAt string                           `json:""`
	Timeout int                              `json:"TimeoutSeconds,omitempty"`
	Version string                           `json:",omitempty"`
	States  map[string]*SfnStateMachineState `json:""`
}

type SfnStateMachineState struct {
}
