package internal

import (
	"strconv"
	"terraform-provider-factorio/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func shouldSuppressDiffPosition(k, old, new string, d *schema.ResourceData) bool {
	oldF, _ := strconv.ParseFloat(old, 64)
	newF, _ := strconv.ParseFloat(new, 64)
	return int64(oldF) == int64(newF)
}

func integerPositionSchema(base *schema.Schema) *schema.Schema {
	posSchema := positionSchema(base)
	innerSchema := posSchema.Elem.(*schema.Resource).Schema
	innerSchema["x"].DiffSuppressFunc = shouldSuppressDiffPosition
	innerSchema["y"].DiffSuppressFunc = shouldSuppressDiffPosition
	return posSchema
}

func positionSchema(base *schema.Schema) *schema.Schema {
	base.Type = schema.TypeList
	base.Elem = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"x": {
				Type:     schema.TypeFloat,
				Optional: base.Optional,
				Required: base.Required,
				Computed: base.Computed,
				ForceNew: base.ForceNew,
			},
			"y": {
				Type:     schema.TypeFloat,
				Optional: base.Optional,
				Required: base.Required,
				Computed: base.Computed,
				ForceNew: base.ForceNew,
			},
		},
	}
	if !base.Computed {
		base.MinItems = 1
		base.MaxItems = 1
	}
	return base
}

func validateDirection(i interface{}, s string) ([]string, []error) {
	_, err := client.ParseDirection(i.(string))
	if err != nil {
		return nil, []error{err}
	}
	return nil, nil
}

func directionSchema(base *schema.Schema) *schema.Schema {
	base.Type = schema.TypeString
	base.ValidateFunc = validateDirection
	return base
}
