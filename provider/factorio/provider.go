package factorio

import (
	"context"
	"strings"

	rcon "github.com/gtaylor/factorio-rcon"
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

	var diags diag.Diagnostics

	if rcon_host == "" {
		return nil, diag.Errorf("rcon_host was empty")
	}
	if rcon_pw == "" {
		return nil, diag.Errorf("rcon_pw was empty")
	}
	r, err := rcon.Dial(rcon_host)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	err = r.Authenticate(rcon_pw)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	_, err = r.Execute(`/c rcon.print(remote.call("terraform-crud-api", "ping"))`)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	// Execute the ping twice in order to skip past the warning about
	// how Lua console commands will disable achievements
	response, err := r.Execute(`/c rcon.print(remote.call("terraform-crud-api", "ping"))`)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	trimmed_response := strings.TrimSpace(response.Body)
	if trimmed_response != "pong" {
		return nil, diag.Errorf("Expected \"pong\" from handshake but got \"%s\"", response.Body)
	}
	return r, diags
}
