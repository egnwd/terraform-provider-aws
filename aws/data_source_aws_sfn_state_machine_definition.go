package aws

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAwsSfnStateMachineDefinition() *schema.Resource {
	topSchema := map[string]*schema.Schema{
		"comment": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"version": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "1.0",
		},
		"timeout_seconds": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},
		"json": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	for k, v := range dataSourceAwsSfnStateMachineDefinitionStates(0) {
		topSchema[k] = v
	}

	return &schema.Resource{
		Read:   dataSourceAwsSfnStateMachineDefinitionRead,
		Schema: topSchema,
	}
}

func dataSourceAwsSfnStateMachineDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	dfn := &SfnStateMachineDefinition{
		Version: d.Get("version").(string),
	}

	if dfnComment, hasDfnComment := d.GetOk("comment"); hasDfnComment {
		dfn.Comment = dfnComment.(string)
	}

	dfn.StartAt = d.Get("start_at").(string)

	if dfnTimeout, hasDfnTimeout := d.GetOk("timeout_seconds"); hasDfnTimeout {
		dfn.Timeout = dfnTimeout.(int)
	}

	states := make(map[string]interface{})
	for _, typ := range dataSourceAwsSfnStateMachineDefinitionStateKeys() {
		cfgStates := d.Get(typ).([]interface{})
		partialStates, err := dataSourceAwsSfnStateMachineDefinitionStateRead(cfgStates, typ)
		if err != nil {
			return fmt.Errorf("error reading %s states: %s", typ, err)
		}
		for k, state := range partialStates {
			states[k] = state
		}
	}

	dfn.States = states

	jsonDfn, err := json.MarshalIndent(dfn, "", "  ")
	if err != nil {
		return err
	}
	jsonString := string(jsonDfn)

	d.Set("json", jsonString)
	d.SetId(strconv.Itoa(hashcode.String(jsonString)))

	return nil
}

func dataSourceAwsSfnStateMachineDefinitionStateRead(cfgStates []interface{}, typ string) (map[string]interface{}, error) {
	states := make(map[string]interface{})

	fn, ok := dataSourceAwsSfnStateMachineDefinitionStateReadFns[typ]
	if !ok {
		return states, fmt.Errorf("could not read type %s", typ)
	}

	for _, stateI := range cfgStates {
		cfgState := stateI.(map[string]interface{})
		n := cfgState["name"].(string)

		state, err := fn(cfgState)
		if err != nil {
			return states, fmt.Errorf("(%s): %s", n, err)
		}

		states[n] = state
	}

	return states, nil
}

// State Read Functions

type stateMachineReadFunc = func(cfgState map[string]interface{}) (interface{}, error)

var dataSourceAwsSfnStateMachineDefinitionStateReadFns map[string]stateMachineReadFunc

func init() {
	dataSourceAwsSfnStateMachineDefinitionStateReadFns = map[string]stateMachineReadFunc{
		"pass":     dataSourceAwsSfnStateMachineDefinitionStatePassRead,
		"succeed":  dataSourceAwsSfnStateMachineDefinitionStateSucceedRead,
		"fail":     dataSourceAwsSfnStateMachineDefinitionStateFailRead,
		"choice":   dataSourceAwsSfnStateMachineDefinitionStateChoiceRead,
		"wait":     dataSourceAwsSfnStateMachineDefinitionStateWaitRead,
		"task":     dataSourceAwsSfnStateMachineDefinitionStateTaskRead,
		"parallel": dataSourceAwsSfnStateMachineDefinitionStateParallelRead,
	}
}

func dataSourceAwsSfnStateMachineDefinitionStateKeys() []string {
	keys := make([]string, 0, len(dataSourceAwsSfnStateMachineDefinitionStateReadFns))
	for k := range dataSourceAwsSfnStateMachineDefinitionStateReadFns {
		keys = append(keys, k)
	}

	return keys
}

func dataSourceAwsSfnStateMachineDefinitionSubStatesRead(cfgStates map[string]interface{}) (*SfnStateMachineStates, error) {
	state := &SfnStateMachineStates{}

	state.StartAt = cfgStates["start_at"].(string)

	states := make(map[string]interface{})

	for _, typ := range dataSourceAwsSfnStateMachineDefinitionStateKeys() {
		cfgSubStates := cfgStates[typ].([]interface{})
		partialStates, err := dataSourceAwsSfnStateMachineDefinitionStateRead(cfgSubStates, typ)
		if err != nil {
			return nil, fmt.Errorf("error reading %s states: %s", typ, err)
		}
		for k, state := range partialStates {
			states[k] = state
		}
	}

	state.States = states

	return state, nil
}

func dataSourceAwsSfnStateMachineDefinitionStateCommonRead(cfgState map[string]interface{}) (interface{}, error) {
	state := &SfnStateMachineState{}

	cfgNext := cfgState["next"].(string)
	hasCfgNext := len(cfgNext) > 0

	cfgEnd := cfgState["end"].(bool)

	if !hasCfgNext && !cfgEnd {
		return state, fmt.Errorf("state has neither next nor end set, exactly one must be specified")
	}

	if hasCfgNext && cfgEnd {
		return state, fmt.Errorf("state has both next and end set, exactly one must be specified, %v", cfgState)
	}

	if hasCfgNext {
		state.Next = cfgNext
	}

	if cfgEnd {
		state.End = cfgEnd
	}

	if cfgComment, hasCfgComment := cfgState["comment"]; hasCfgComment {
		state.Comment = cfgComment.(string)
	}

	state.InputPath = sfnStateMachineDefinitionPath(cfgState["input_path"].(string))
	state.OutputPath = sfnStateMachineDefinitionPath(cfgState["output_path"].(string))

	return state, nil
}

func dataSourceAwsSfnStateMachineDefinitionStatePassRead(cfgState map[string]interface{}) (interface{}, error) {
	commonState, err := dataSourceAwsSfnStateMachineDefinitionStateCommonRead(cfgState)
	if err != nil {
		return nil, err
	}

	state := &SfnStateMachinePassState{
		SfnStateMachineState: *commonState.(*SfnStateMachineState),
	}

	state.Type = "Pass"

	if resultJson := cfgState["result"].(string); len(resultJson) > 0 {
		result, err := structure.ExpandJsonFromString(resultJson)
		if err != nil {
			return nil, fmt.Errorf("invalid result JSON: %s", err)
		}
		state.Result = result
	}

	state.ResultPath = sfnStateMachineDefinitionPath(cfgState["result_path"].(string))

	if parametersJson := cfgState["parameters"].(string); len(parametersJson) > 0 {
		parameters, err := structure.ExpandJsonFromString(parametersJson)
		if err != nil {
			return nil, fmt.Errorf("invalid parameters JSON: %s", err)
		}
		state.Parameters = parameters
	}

	return state, nil
}

func dataSourceAwsSfnStateMachineDefinitionStateSucceedRead(cfgState map[string]interface{}) (interface{}, error) {
	state := &SfnStateMachineSucceedState{}

	state.Type = "Succeed"

	if cfgComment, hasCfgComment := cfgState["comment"]; hasCfgComment {
		state.Comment = cfgComment.(string)
	}

	return state, nil
}

func dataSourceAwsSfnStateMachineDefinitionStateFailRead(cfgState map[string]interface{}) (interface{}, error) {
	state := &SfnStateMachineFailState{}

	state.Type = "Fail"

	if cfgComment, hasCfgComment := cfgState["comment"]; hasCfgComment {
		state.Comment = cfgComment.(string)
	}

	if cfgCause, hasCfgCause := cfgState["cause"]; hasCfgCause {
		state.Cause = cfgCause.(string)
	}

	if cfgError, hasCfgError := cfgState["error"]; hasCfgError {
		state.Error = cfgError.(string)
	}

	return state, nil
}

func dataSourceAwsSfnStateMachineDefinitionStateChoiceRead(cfgState map[string]interface{}) (interface{}, error) {
	state := &SfnStateMachineChoiceState{} // Cannot use CommonRead because Choice ignores Next & End

	state.Type = "Choice"

	if cfgComment, hasCfgComment := cfgState["comment"]; hasCfgComment {
		state.Comment = cfgComment.(string)
	}

	state.InputPath = sfnStateMachineDefinitionPath(cfgState["input_path"].(string))
	state.OutputPath = sfnStateMachineDefinitionPath(cfgState["output_path"].(string))

	cfgChoices := cfgState["option"].([]interface{})
	choices := make([]*SfnStateMachineChoiceRule, len(cfgChoices))

	for i, choiceI := range cfgChoices {
		choice := &SfnStateMachineChoiceRule{}
		cfgChoice := choiceI.(map[string]interface{})

		comparison, err := structure.ExpandJsonFromString(cfgChoice["comparison"].(string))
		if err != nil {
			// Shouldn't happen due to validation
			return nil, fmt.Errorf("invalid comparison JSON: %s", err)
		}
		choice.Comparison = comparison

		if cfgNext, hasCfgNext := cfgChoice["next"]; hasCfgNext {
			choice.Next = cfgNext.(string)
		}

		choices[i] = choice
	}

	state.Choices = choices

	if cfgDefault, hasCfgDefault := cfgState["default"]; hasCfgDefault {
		state.Default = cfgDefault.(string)
	}

	return state, nil
}

func dataSourceAwsSfnStateMachineDefinitionStateWaitRead(cfgState map[string]interface{}) (interface{}, error) {
	commonState, err := dataSourceAwsSfnStateMachineDefinitionStateCommonRead(cfgState)
	if err != nil {
		return nil, err
	}

	state := &SfnStateMachineWaitState{
		SfnStateMachineState: *commonState.(*SfnStateMachineState),
	}

	state.Type = "Wait"

	if cfgWait := cfgState["timestamp"].(string); len(cfgWait) > 0 {
		state.Timestamp = &cfgWait
	} else if cfgWait := cfgState["seconds_path"].(string); len(cfgWait) > 0 {
		state.SecondsPath = &cfgWait
	} else if cfgWait := cfgState["timestamp_path"].(string); len(cfgWait) > 0 {
		state.TimestampPath = &cfgWait
	} else {
		cfgWait := cfgState["seconds"].(int)
		state.Seconds = &cfgWait
	}

	return state, nil
}

func dataSourceAwsSfnStateMachineDefinitionStateTaskRead(cfgState map[string]interface{}) (interface{}, error) {
	commonState, err := dataSourceAwsSfnStateMachineDefinitionStateCommonRead(cfgState)
	if err != nil {
		return nil, err
	}

	state := &SfnStateMachineTaskState{
		SfnStateMachineState: *commonState.(*SfnStateMachineState),
	}

	state.Type = "Task"

	state.Resource = cfgState["resource"].(string)

	if cfgParameters := cfgState["parameters"].(string); len(cfgParameters) > 0 {
		parameters, err := structure.ExpandJsonFromString(cfgParameters)
		if err != nil {
			// Shouldn't happen due to validation
			return nil, fmt.Errorf("invalid parameters JSON: %s", err)
		}

		state.Parameters = parameters
	}

	state.ResultPath = sfnStateMachineDefinitionPath(cfgState["result_path"].(string))
	state.ResultSelector = sfnStateMachineDefinitionPath(cfgState["result_selector"].(string))

	if cfgRetriers, hasCfgRetriers := cfgState["retry"]; hasCfgRetriers {
		state.Retry = dataSourceAwsSfnStateMachineDefinitionRetriersRead(cfgRetriers.([]interface{}))
	}

	if cfgCatchers, hasCfgCatchers := cfgState["catch"]; hasCfgCatchers {
		state.Catch = dataSourceAwsSfnStateMachineDefinitionCatchersRead(cfgCatchers.([]interface{}))
	}

	cfgTimeout := cfgState["timeout"]
	cfgTimeoutPath := cfgState["timeout_path"]
	hasCfgTimeout := cfgTimeout != 0
	hasCfgTimeoutPath := cfgTimeoutPath != ""

	if hasCfgTimeout && hasCfgTimeoutPath {
		return nil, fmt.Errorf("both timeout and timeout_path are set (%s, %s)", cfgTimeout, cfgTimeoutPath)
	}

	if hasCfgTimeout {
		state.TimeoutSeconds = cfgTimeout.(int)
	}

	if hasCfgTimeoutPath {
		state.TimeoutSecondsPath = cfgTimeoutPath.(string)
	}

	cfgHeartbeat := cfgState["heartbeat"]
	cfgHeartbeatPath := cfgState["heartbeat_path"]
	hasCfgHeartbeat := cfgHeartbeat != 0
	hasCfgHeartbeatPath := cfgHeartbeatPath != ""

	if hasCfgHeartbeat && hasCfgHeartbeatPath {
		return nil, errors.New("both heartbeat and heartbeat_path are set")
	}

	if hasCfgHeartbeat {
		state.HeartbeatSeconds = cfgHeartbeat.(int)
	}

	if hasCfgHeartbeatPath {
		state.HeartbeatSecondsPath = cfgHeartbeatPath.(string)
	}

	return state, nil
}

func dataSourceAwsSfnStateMachineDefinitionStateParallelRead(cfgState map[string]interface{}) (interface{}, error) {
	commonState, err := dataSourceAwsSfnStateMachineDefinitionStateCommonRead(cfgState)
	if err != nil {
		return nil, err
	}

	state := &SfnStateMachineParallelState{
		SfnStateMachineState: *commonState.(*SfnStateMachineState),
	}

	state.Type = "Parallel"

	state.ResultPath = sfnStateMachineDefinitionPath(cfgState["result_path"].(string))
	state.ResultSelector = sfnStateMachineDefinitionPath(cfgState["result_selector"].(string))

	if cfgRetriers, hasCfgRetriers := cfgState["retry"]; hasCfgRetriers {
		state.Retry = dataSourceAwsSfnStateMachineDefinitionRetriersRead(cfgRetriers.([]interface{}))
	}

	if cfgCatchers, hasCfgCatchers := cfgState["catch"]; hasCfgCatchers {
		state.Catch = dataSourceAwsSfnStateMachineDefinitionCatchersRead(cfgCatchers.([]interface{}))
	}

	cfgBranches := cfgState["branch"].([]interface{})
	branches := make([]*SfnStateMachineStates, len(cfgBranches))

	for i, branchI := range cfgBranches {
		branch, err := dataSourceAwsSfnStateMachineDefinitionSubStatesRead(branchI.(map[string]interface{}))
		if err != nil {
			return nil, fmt.Errorf("error reading branch: %s\n", err)
		}

		branches[i] = branch
	}

	state.Branches = branches

	return state, nil
}

func dataSourceAwsSfnStateMachineDefinitionRetriersRead(in []interface{}) []*SfnStateMachineRetrier {
	retriers := make([]*SfnStateMachineRetrier, len(in))

	for i, retrierI := range in {
		retrier := &SfnStateMachineRetrier{}
		cfgRetrier := retrierI.(map[string]interface{})

		retrier.ErrorEquals = sfnStateMachineDefinitionConfigStringList(cfgRetrier["errors"].(*schema.Set).List())

		if cfgInterval, hasCfgInterval := cfgRetrier["interval"]; hasCfgInterval {
			retrier.IntervalSeconds = cfgInterval.(int)
		}

		if cfgMaxAttempts, hasCfgMaxAttempts := cfgRetrier["max_attempts"]; hasCfgMaxAttempts {
			retrier.MaxAttempts = cfgMaxAttempts.(int)
		}

		if cfgBackoff, hasCfgInterval := cfgRetrier["backoff"]; hasCfgInterval {
			retrier.BackoffRate = cfgBackoff.(float64)
		}

		retriers[i] = retrier
	}

	return retriers
}

func dataSourceAwsSfnStateMachineDefinitionCatchersRead(in []interface{}) []*SfnStateMachineCatcher {
	catchers := make([]*SfnStateMachineCatcher, len(in))

	for i, catcherI := range in {
		catcher := &SfnStateMachineCatcher{}
		cfgCatcher := catcherI.(map[string]interface{})

		catcher.ErrorEquals = sfnStateMachineDefinitionConfigStringList(cfgCatcher["errors"].(*schema.Set).List())

		if cfgNext, hasCfgNext := cfgCatcher["next"]; hasCfgNext {
			catcher.Next = cfgNext.(string)
		}

		catchers[i] = catcher
	}

	return catchers
}

// Schemas

func dataSourceAwsSfnStateMachineDefinitionStateSchema(fn func() map[string]*schema.Schema) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: fn(),
		},
	}
}

const maxStatesDepth = 5

func validateStateMachineDefinitionPath() schema.SchemaValidateFunc {
	re := regexp.MustCompile(`^\$.*`)
	isEmpty := validation.StringLenBetween(0, 0)
	orIsValidPath := validation.StringMatch(re, "JSON Path must begin with '$'")
	return validation.Any(isEmpty, orIsValidPath)
}

func dataSourceAwsSfnStateMachineDefinitionStates(d int) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"start_at": {
			Type:     schema.TypeString,
			Required: true,
		},
		"pass":    dataSourceAwsSfnStateMachineDefinitionStateSchema(dataSourceAwsSfnStateMachineDefinitionPassState),
		"succeed": dataSourceAwsSfnStateMachineDefinitionStateSchema(dataSourceAwsSfnStateMachineDefinitionSucceedState),
		"fail":    dataSourceAwsSfnStateMachineDefinitionStateSchema(dataSourceAwsSfnStateMachineDefinitionFailState),
		"choice":  dataSourceAwsSfnStateMachineDefinitionStateSchema(dataSourceAwsSfnStateMachineDefinitionChoiceState),
		"wait":    dataSourceAwsSfnStateMachineDefinitionStateSchema(dataSourceAwsSfnStateMachineDefinitionWaitState),
		"task":    dataSourceAwsSfnStateMachineDefinitionStateSchema(dataSourceAwsSfnStateMachineDefinitionTaskState),
	}

	if d < maxStatesDepth {
		s["parallel"] = dataSourceAwsSfnStateMachineDefinitionStateSchema(dataSourceAwsSfnStateMachineDefinitionParallelState(d + 1))
	}

	return s
}

func dataSourceAwsSfnStateMachineDefinitionCommonState() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"next": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"end": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"comment": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"input_path": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "$",
			ValidateFunc: validateStateMachineDefinitionPath(),
		},
		"output_path": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "$",
			ValidateFunc: validateStateMachineDefinitionPath(),
		},
	}
}

func dataSourceAwsSfnStateMachineDefinitionPassState() map[string]*schema.Schema {
	passSchema := map[string]*schema.Schema{
		"result": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsJSON,
		},
		"result_path": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "$",
			ValidateFunc: validateStateMachineDefinitionPath(),
		},
		"parameters": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsJSON,
		},
	}

	for k, v := range dataSourceAwsSfnStateMachineDefinitionCommonState() {
		passSchema[k] = v
	}

	return passSchema
}

func dataSourceAwsSfnStateMachineDefinitionSucceedState() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"comment": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func dataSourceAwsSfnStateMachineDefinitionFailState() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"comment": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"cause": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"error": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func dataSourceAwsSfnStateMachineDefinitionChoiceState() map[string]*schema.Schema {
	passSchema := map[string]*schema.Schema{
		"option": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"next": {
						Type:     schema.TypeString,
						Required: true,
					},
					"comparison": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsJSON,
					},
				},
			},
		},
		"default": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	for k, v := range dataSourceAwsSfnStateMachineDefinitionCommonState() {
		passSchema[k] = v
	}

	return passSchema
}

func dataSourceAwsSfnStateMachineDefinitionWaitState() map[string]*schema.Schema {
	waitSchema := map[string]*schema.Schema{
		"seconds": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},
		"timestamp": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},
		"seconds_path": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateStateMachineDefinitionPath(),
		},
		"timestamp_path": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateStateMachineDefinitionPath(),
		},
	}

	for k, v := range dataSourceAwsSfnStateMachineDefinitionCommonState() {
		waitSchema[k] = v
	}

	return waitSchema
}

func dataSourceAwsSfnStateMachineDefinitionTaskState() map[string]*schema.Schema {
	taskSchema := map[string]*schema.Schema{
		"resource": {
			Type:     schema.TypeString,
			Required: true,
		},
		"parameters": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsJSON,
		},
		"result_path": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "$",
			ValidateFunc: validateStateMachineDefinitionPath(),
		},
		"result_selector": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "$",
			ValidateFunc: validateStateMachineDefinitionPath(),
		},
		"retry": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"errors": {
						Type:     schema.TypeSet,
						Required: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"interval": {
						Type:         schema.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"max_attempts": {
						Type:         schema.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"backoff": {
						Type:         schema.TypeFloat,
						Optional:     true,
						ValidateFunc: validation.FloatAtLeast(0),
					},
				},
			},
		},
		"catch": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"errors": {
						Type:     schema.TypeSet,
						Required: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"next": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"timeout": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		"timeout_path": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateStateMachineDefinitionPath(),
		},
		"heartbeat": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		"heartbeat_path": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateStateMachineDefinitionPath(),
		},
	}

	for k, v := range dataSourceAwsSfnStateMachineDefinitionCommonState() {
		taskSchema[k] = v
	}

	return taskSchema
}

func dataSourceAwsSfnStateMachineDefinitionParallelState(d int) func() map[string]*schema.Schema {
	return func() map[string]*schema.Schema {
		taskSchema := map[string]*schema.Schema{
			"branch": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: dataSourceAwsSfnStateMachineDefinitionStates(d),
				},
			},
			"result_path": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "$",
				ValidateFunc: validateStateMachineDefinitionPath(),
			},
			"result_selector": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "$",
				ValidateFunc: validateStateMachineDefinitionPath(),
			},
			"retry": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"errors": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
						"max_attempts": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
						"backoff": {
							Type:         schema.TypeFloat,
							Optional:     true,
							ValidateFunc: validation.FloatAtLeast(0),
						},
					},
				},
			},
			"catch": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"errors": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"next": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		}

		for k, v := range dataSourceAwsSfnStateMachineDefinitionCommonState() {
			taskSchema[k] = v
		}

		return taskSchema
	}
}
