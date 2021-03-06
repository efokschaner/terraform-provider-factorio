package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-factorio/client"
)

func resourceHello() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHelloCreate,
		ReadContext:   resourceHelloRead,
		UpdateContext: resourceHelloUpdate,
		DeleteContext: resourceHelloDelete,
		Description:   "A message made of conveyor belts. Not very useful.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_as_ghost": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Used only at time of creation. Whether to create a ghost or the actual entity.",
				// Updating this "creation-time" value necessitates re-creation
				ForceNew: true,
			},
			"direction": directionSchema(&schema.Schema{
				Optional: true,
				Default:  "east",
			}),
		},
	}
}

func resourceHelloCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.FactorioClient)
	create_config := make(map[string]interface{})
	create_config["create_as_ghost"] = d.Get("create_as_ghost")
	create_config["direction"] = d.Get("direction")
	hello := make(map[string]interface{})
	err := c.Create("hello", create_config, &hello)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(hello["id"].(string))
	return resourceHelloRead(ctx, d, m)
}

func resourceHelloRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	c := m.(*client.FactorioClient)
	hello := make(map[string]interface{})
	err := c.Read("hello", map[string]interface{}{"id": d.Id()}, &hello)
	if err != nil {
		return diag.FromErr(err)
	}
	id, id_exists := hello["id"]
	if !id_exists {
		d.SetId("")
		return diags
	}
	d.SetId(id.(string))
	d.Set("direction", hello["direction"])
	return diags
}

func resourceHelloUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.FactorioClient)
	hello_updates := make(map[string]interface{})
	if d.HasChange("direction") {
		hello_updates["direction"] = d.Get("direction")
	}
	err := c.Update(
		"hello",
		d.Id(),
		hello_updates,
		nil)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceHelloRead(ctx, d, m)
}

func resourceHelloDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.FactorioClient)
	var diags diag.Diagnostics
	err := c.Delete("hello", d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
