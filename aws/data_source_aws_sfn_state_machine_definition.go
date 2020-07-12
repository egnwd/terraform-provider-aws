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
			"state": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Pass",
							}, false),
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
					},
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

	cfgStates := d.Get("state").([]interface{})
	states := make(map[string]*SfnStateMachineState)

	for _, stateI := range cfgStates {
		cfgState := stateI.(map[string]interface{})
		n := cfgState["name"].(string)

		state, err := dataSourceAwsSfnStateMachineDefinitionStateRead(cfgState)
		if err != nil {
			return fmt.Errorf("error reading state (%s): %s", n, err)
		}

		states[n] = state
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

func dataSourceAwsSfnStateMachineDefinitionStateRead(cfgState map[string]interface{}) (*SfnStateMachineState, error) {
	state := &SfnStateMachineState{}
	state.Type = cfgState["type"].(string)

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

	cfgInputPath := cfgState["input_path"].(string)
	cfgOutputPath := cfgState["output_path"].(string)

	if len(cfgInputPath) > 0 {
		state.InputPath = &cfgInputPath
	}

	if len(cfgOutputPath) > 0 {
		state.OutputPath = &cfgOutputPath
	}

	return state, nil
}
