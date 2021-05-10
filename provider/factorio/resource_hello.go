package factorio

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceHello() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHelloCreate,
		ReadContext:   resourceHelloRead,
		UpdateContext: resourceHelloUpdate,
		DeleteContext: resourceHelloDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_as_ghost": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Used only at time of creation. Whether to create a ghost or the actual entity.",
			},
		},
	}
}

func resourceHelloCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*factorioClient)
	create_config := make(map[string]interface{})
	create_config["create_as_ghost"] = d.Get("create_as_ghost")
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
	c := m.(*factorioClient)
	hello := make(map[string]interface{})
	err := c.Read("hello", map[string]interface{}{"id": d.Id()}, &hello)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(hello["id"].(string))
	return diags
}

func resourceHelloUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceHelloRead(ctx, d, m)
}

func resourceHelloDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
