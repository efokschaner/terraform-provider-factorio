package factorio

import (
	"encoding/json"
	"fmt"
	"sync"

	rcon "github.com/gtaylor/factorio-rcon"
)

type factorioClient struct {
	conn *rcon.RCON
	// rcon.RCON is not threadsafe, quick and dirty mutex
	// TODO rewrite RCON to handle parallel / interleaved calls
	mutex sync.Mutex
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

func (client *factorioClient) DoHandShake() error {
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

func (client *factorioClient) Read(resource_type string, query interface{}, result_out interface{}) error {
	query_bytes, err := json.Marshal(query)
	if err != nil {
		return err
	}
	return client.doCall(result_out, "read", resource_type, string(query_bytes))
}

func (client *factorioClient) Create(resource_type string, create_config interface{}, result_out interface{}) error {
	create_config_bytes, err := json.Marshal(create_config)
	if err != nil {
		return err
	}
	return client.doCall(result_out, "create", resource_type, string(create_config_bytes))
}

func (client *factorioClient) Update(resource_type string, resource_id string, update_config map[string]interface{}) error {
	update_config_bytes, err := json.Marshal(update_config)
	if err != nil {
		return err
	}
	var ignore interface{}
	return client.doCall(&ignore, "update", resource_type, resource_id, string(update_config_bytes))
}

func (client *factorioClient) Delete(resource_type string, resource_id string) error {
	result := struct {
		IsDeleted bool `json:"is_deleted"`
	}{
		IsDeleted: false,
	}
	err := client.doCall(&result, "delete", resource_type, resource_id)
	if err != nil {
		return err
	}
	if !result.IsDeleted {
		return fmt.Errorf("reportedly failed to delete")
	}
	return nil
}

func (client *factorioClient) doCall(result_out interface{}, params ...string) error {
	client.mutex.Lock()
	response, err := client.conn.Execute(formatRconCommand(params))
	client.mutex.Unlock()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(response.Body), result_out)
	if err != nil {
		err = fmt.Errorf("unmarshalling \"%v\": %v", response.Body, err)
	}
	return err
}

func formatRconCommand(params []string) string {
	// TODO: Replace the api with something a bit more jsonrpc-ish
	// to simplify method / param / response delivery
	result := "/c rcon.print(remote.call(\"terraform-crud-api\""
	for _, p := range params {
		// Use single quotes to avoid conflict with json double quotes
		result += ",'" + p + "'"
	}
	result += "))"
	return result
}
