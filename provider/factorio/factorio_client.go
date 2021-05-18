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
	return client.doCall(result_out, "read", resource_type, query)
}

func (client *factorioClient) Create(resource_type string, create_config interface{}, result_out interface{}) error {
	return client.doCall(result_out, "create", resource_type, create_config)
}

func (client *factorioClient) Update(resource_type string, resource_id string, update_config map[string]interface{}) error {
	return client.doCall(nil, "update", resource_type, resource_id, update_config)
}

func (client *factorioClient) Delete(resource_type string, resource_id string) error {
	result := struct {
		ResourceExists bool `json:"resource_exists"`
	}{
		ResourceExists: true,
	}
	err := client.doCall(&result, "delete", resource_type, resource_id)
	if err != nil {
		return err
	}
	if result.ResourceExists {
		return fmt.Errorf("resource still exists")
	}
	return nil
}

type RpcRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

type RpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RpcResponse struct {
	Result *json.RawMessage `json:"result"`
	Error  *RpcError        `json:"error"`
}

func (client *factorioClient) doCall(result interface{}, method string, params ...interface{}) error {
	if params == nil {
		params = []interface{}{}
	}
	req := RpcRequest{
		Method: method,
		Params: params,
	}
	request_bytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	// Use single quotes around request_bytes
	// to avoid conflict with json double quotes
	// TODO: Escape single quotes in request_bytes
	command := fmt.Sprintf("/c rcon.print(remote.call('terraform-crud-api', 'call', '%s'))", request_bytes)
	client.mutex.Lock()
	executeResponse, err := client.conn.Execute(command)
	client.mutex.Unlock()
	if err != nil {
		return err
	}
	var response RpcResponse
	err = json.Unmarshal([]byte(executeResponse.Body), &response)
	if err != nil {
		return fmt.Errorf("unmarshalling \"%v\": %v", executeResponse.Body, err)
	}
	if response.Error != nil {
		return fmt.Errorf(
			"error from api, code: %d, message: \"%s\", details: %+v",
			response.Error.Code,
			response.Error.Message,
			response.Error.Data)
	}
	if response.Result == nil {
		return nil
	}
	return json.Unmarshal(*response.Result, result)
}
