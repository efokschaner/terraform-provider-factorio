package factorio

import (
	"encoding/json"
	"fmt"

	rcon "github.com/gtaylor/factorio-rcon"
)

type factorioClient struct {
	conn *rcon.RCON
}

func NewFactorioClient(rcon_host string, rcon_password string) (*factorioClient, error) {
	r, err := rcon.Dial(rcon_host)
	if err != nil {
		return nil, err
	}
	err = r.Authenticate(rcon_password)
	if err != nil {
		return nil, err
	}
	c := new(factorioClient)
	c.conn = r
	err = c.DoHandShake()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (client factorioClient) DoHandShake() error {
	var result string
	// Ignore the first error.
	// Execute the ping twice in order to skip past the warning about
	// how Lua console commands will disable achievements
	client.doCall(&result, "ping")
	err := client.doCall(&result, "ping")
	if err != nil {
		return err
	}
	if result != "pong" {
		return fmt.Errorf("expected \"pong\" from handshake but got \"%s\"", result)
	}
	return nil
}

func (client factorioClient) Read(resource_type string, result_out interface{}) error {
	return client.doCall(result_out, "read", resource_type)
}

func (client factorioClient) Create(resource_type string, result_out interface{}) error {
	return client.doCall(result_out, "create", resource_type)
}

func (client factorioClient) Update(resource_type string, result_out interface{}) error {
	return client.doCall(result_out, "update", resource_type)
}

func (client factorioClient) Delete(resource_type string) error {
	ignore := make([]interface{}, 0)
	return client.doCall(&ignore, "delete", resource_type)
}

func (client factorioClient) doCall(result_out interface{}, params ...string) error {
	response, err := client.conn.Execute(formatRconCommand(params))
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(response.Body), result_out)
	return err
}

func formatRconCommand(params []string) string {
	// TODO: Replace the api with something a bit more jsonrpc-ish
	// to simplify method / param / response delivery
	result := "/c rcon.print(remote.call(\"terraform-crud-api\""
	for _, p := range params {
		result += ",\"" + p + "\""
	}
	result += "))"
	return result
}
