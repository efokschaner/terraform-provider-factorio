package factorio

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"rcon_host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FACTORIO_RCON_HOST", nil),
			},
			"rcon_pw": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("FACTORIO_RCON_PW", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"factorio_players": dataSourcePlayers(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	rcon_host := d.Get("rcon_host").(string)
	rcon_pw := d.Get("rcon_pw").(string)
	if rcon_host == "" {
		return nil, diag.Errorf("rcon_host was empty")
	}
	if rcon_pw == "" {
		return nil, diag.Errorf("rcon_pw was empty")
	}
	client, err := NewFactorioClient(rcon_host, rcon_pw)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return client, nil
}
