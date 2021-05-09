package factorio

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	rcon "github.com/gtaylor/factorio-rcon"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePlayers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePlayersRead,
		Schema: map[string]*schema.Schema{
			"players": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePlayersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	r := m.(*rcon.RCON)
	response, err := r.Execute(`/c rcon.print(remote.call("terraform-crud-api", "read", "players"))`)
	if err != nil {
		return diag.FromErr(err)
	}
	players := make([]map[string]interface{}, 0)
	err = json.Unmarshal([]byte(response.Body), &players)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("players", players); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
