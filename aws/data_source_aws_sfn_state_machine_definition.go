package aws

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAwsSfnStateMachineDefinition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsSfnStateMachineDefinitionRead,

		Schema: map[string]*schema.Schema{
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_at": {
				Type:     schema.TypeString,
				Required: true,
			},
			"timeout_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "1.0",
			},
			"pass": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: dataSourceAwsSfnStateMachineDefinitionPassState(),
				},
			},
			"succeed": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: dataSourceAwsSfnStateMachineDefinitionSucceedState(),
				},
			},
			"fail": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: dataSourceAwsSfnStateMachineDefinitionFailState(),
				},
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
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
		partialStates, err := dataSourceAwsSfnStateMachineDefinitionStateRead(d, typ)
		if err != nil {
			fmt.Errorf("error reading %s states: %s", typ, err)
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

func dataSourceAwsSfnStateMachineDefinitionStateRead(d *schema.ResourceData, typ string) (map[string]interface{}, error) {
	states := make(map[string]interface{})
	cfgStates := d.Get(typ).([]interface{})

	fn, ok := dataSourceAwsSfnStateMachineDefinitionStateReadFns[typ]
	if !ok {
		return states, fmt.Errorf("error reading type %s", typ)
	}

	for _, stateI := range cfgStates {
		cfgState := stateI.(map[string]interface{})
		n := cfgState["name"].(string)

		state, err := fn(cfgState)
		if err != nil {
			return states, fmt.Errorf("error reading state (%s): %s", n, err)
		}

		states[n] = state
	}

	return states, nil
}

// State Read Functions
type stateMachineReadFunc = func(cfgState map[string]interface{}) (interface{}, error)

var dataSourceAwsSfnStateMachineDefinitionStateReadFns map[string]stateMachineReadFunc = map[string]stateMachineReadFunc{
	"pass":    dataSourceAwsSfnStateMachineDefinitionStatePassRead,
	"succeed": dataSourceAwsSfnStateMachineDefinitionStateSucceedRead,
	"fail":    dataSourceAwsSfnStateMachineDefinitionStateFailRead,
}

func dataSourceAwsSfnStateMachineDefinitionStateKeys() []string {
	keys := make([]string, 0, len(dataSourceAwsSfnStateMachineDefinitionStateReadFns))
	for k := range dataSourceAwsSfnStateMachineDefinitionStateReadFns {
		keys = append(keys, k)
	}

	return keys
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

	if cfgInputPath := cfgState["input_path"].(string); len(cfgInputPath) > 0 {
		state.InputPath = &cfgInputPath
	}

	if cfgOutputPath := cfgState["output_path"].(string); len(cfgOutputPath) > 0 {
		state.OutputPath = &cfgOutputPath
	}

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
		result := make(map[string]interface{})
		err := json.Unmarshal([]byte(resultJson), &result)
		if err != nil {
			return nil, fmt.Errorf("invalid result JSON: %s", err)
		}
		state.Result = result
	}

	if cfgResultPath := cfgState["result_path"].(string); len(cfgResultPath) > 0 {
		state.ResultPath = &cfgResultPath
	}

	if parametersJson := cfgState["parameters"].(string); len(parametersJson) > 0 {
		parameters := make(map[string]interface{})
		err := json.Unmarshal([]byte(parametersJson), &parameters)
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

// Schemas

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
			Type:     schema.TypeString,
			Required: true,
		},
		"output_path": {
			Type:     schema.TypeString,
			Required: true,
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
			Type:     schema.TypeString,
			Required: true,
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
