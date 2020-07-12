package aws

import (
	"encoding/json"
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
		state := &SfnStateMachineState{}

		stateName := cfgState["name"].(string)

		states[stateName] = state
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
