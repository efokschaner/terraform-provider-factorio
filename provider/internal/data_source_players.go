package internal

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-factorio/client"
)

func dataSourcePlayers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePlayersRead,
		Schema: map[string]*schema.Schema{
			"players": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"position": positionSchema(&schema.Schema{
							Computed: true,
						}),
					},
				},
			},
		},
	}
}

func dataSourcePlayersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.FactorioClient)
	players := make([]map[string]interface{}, 0)
	err := c.Read("players", nil, &players)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, player := range players {
		// "flattening" as terraform calls it.
		// It seems in Terraform, all nested objects are arrays of length one.
		player["position"] = []interface{}{player["position"]}
	}

	if err := d.Set("players", players); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
