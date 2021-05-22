package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-factorio/client"
)

func resourceEntity() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEntityCreate,
		ReadContext:   resourceEntityRead,
		UpdateContext: resourceEntityUpdate,
		DeleteContext: resourceEntityDelete,
		Description:   "A LuaEntity in Factorio (https://lua-api.factorio.com/latest/LuaEntity.html), see LuaSurface.create_entity for creation reference (https://lua-api.factorio.com/latest/LuaSurface.html#LuaSurface.create_entity) ",

		Schema: map[string]*schema.Schema{
			"unit_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "unit_number is Factorio's concept of an ID",
			},
			"surface": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "nauvis",
				Description: "The LuaSurface on which the LuaEntity is placed (https://lua-api.factorio.com/latest/LuaSurface.html)",
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The prototype name of the LuaEntity (https://wiki.factorio.com/Prototype_definitions)",
				ForceNew:    true,
			},
			"position": integerPositionSchema(&schema.Schema{
				Required:    true,
				Description: "The position of the LuaEntity.",
				ForceNew:    true,
			}),
			// TODO force 'north' for entities with the 'not-rotatable' flag
			"direction": directionSchema(&schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "north",
				Description: "Which direction the LuaEntity faces.",
			}),
			"force": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "player",
				Description: "The force of this LuaEntity, eg. \"player\", \"enemy\", \"neutral\" (https://lua-api.factorio.com/latest/LuaControl.html#LuaControl.force)",
			},
			"entity_specific_parameters": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A map of additional entity-specific parameters to be passed to create_entity (https://lua-api.factorio.com/latest/LuaSurface.html#LuaSurface.create_entity)",
				ForceNew:    true,
			},
		},
	}
}

func resourceEntityCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.FactorioClient)
	var opts client.EntityCreateOptions

	direction, err := client.ParseDirection(d.Get("direction").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	opts.Surface = d.Get("surface").(string)
	opts.Name = d.Get("name").(string)
	opts.Position.X = float32(d.Get("position.0.x").(float64))
	opts.Position.Y = float32(d.Get("position.0.y").(float64))
	opts.Direction = direction
	opts.Force = d.Get("force").(string)
	opts.EntitySpecificParameters = d.Get("entity_specific_parameters").(map[string]interface{})

	e, err := c.EntityCreate(&opts)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(e.UnitNumber.String())
	return resourceEntityRead(ctx, d, m)
}

func writeAttributeToResource(diagOut *diag.Diagnostics, d *schema.ResourceData, key string, attr interface{}) {
	err := d.Set(key, attr)
	if err != nil {
		*diagOut = append(*diagOut, diag.FromErr(err)...)
	}
}

func flattenPosition(pos client.Position) []map[string]float64 {
	flat := make(map[string]float64)
	flat["x"] = float64(pos.X)
	flat["y"] = float64(pos.Y)
	return []map[string]float64{flat}
}

func writeEntityToResourceData(e *client.Entity, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	writeAttributeToResource(&diags, d, "unit_number", e.UnitNumber)
	writeAttributeToResource(&diags, d, "surface", e.Surface)
	writeAttributeToResource(&diags, d, "name", e.Name)
	writeAttributeToResource(&diags, d, "position", flattenPosition(e.Position))
	writeAttributeToResource(&diags, d, "direction", e.Direction.String())
	writeAttributeToResource(&diags, d, "force", e.Force)
	return diags
}

func resourceEntityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.FactorioClient)
	unitNumber, err := client.ParseUnitNumber(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	entity, err := c.EntityGet(unitNumber)
	if err != nil {
		return diag.FromErr(err)
	}
	if entity == nil {
		d.SetId("")
		return nil
	}
	d.SetId(entity.UnitNumber.String())
	diags := writeEntityToResourceData(entity, d)
	return diags
}

func resourceEntityUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.FactorioClient)
	unitNumber, err := client.ParseUnitNumber(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	var opts client.EntityUpdateOptions
	if d.HasChange("direction") {
		direction, err := client.ParseDirection(d.Get("direction").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		opts.Direction = &direction
	}
	if d.HasChange("force") {
		force := d.Get("force").(string)
		opts.Force = &force
	}
	_, err = c.EntityUpdate(unitNumber, &opts)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceEntityRead(ctx, d, m)
}

func resourceEntityDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.FactorioClient)
	unitNumber, err := client.ParseUnitNumber(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	err = c.EntityDelete(unitNumber)
	return diag.FromErr(err)
}
