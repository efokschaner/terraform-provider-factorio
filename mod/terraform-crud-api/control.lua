resources = {
  players = {
    read = function()
      result = {}
      for k,v in pairs(game.players) do
        table.insert(result, {
          name = v.name,
          position = v.position,
        })
      end
      return result
    end,
  },
  entity = require('resources.entity'),
  hello = require('resources.hello'),
}

-- In addition to json serialization, this code also replaces empty objects with empty arrays
-- due to the fact that lua cannot differentiate.
-- This means we cannot have empty objects, however it is easier to add an unused key
-- to an object to avoid empty objects than to add an unused element to every array.
function serialize(value)
  if value == nil then
    return 'null'
  elseif type(value) == 'string' then
    return '"' .. value .. '"'
  end
  -- Wrapped with () to prevent gsub returning a multival
  return (string.gsub(game.table_to_json(value), '{}', '[]'))
end

exports = {
  ping = function() return 'pong' end,

  read = function(resource_type, query)
    return resources[resource_type].read(query)
  end,

  create = function(resource_type, create_config)
    return resources[resource_type].create(create_config)
  end,

  update = function(resource_type, resource_id, update_config)
    return resources[resource_type].update(resource_id, update_config)
  end,

  delete = function(resource_type, resource_id)
    return resources[resource_type].delete(resource_id)
  end,
}

local function handle_rpc(request_string)
  local deserialize_succeeded, deserialize_result = xpcall(
    game.json_to_table,
    debug.traceback,
    request_string)
  if not deserialize_succeeded then
    return {
      error = {
        code = 400,
        message = 'Failed to deserialize request_string',
        data = deserialize_result
      }
    }
  end
  local request = deserialize_result
  local method = exports[request.method]
  if method == nil then
    return {
      error = {
        code = 404,
        message = string.format('No method named "%s"', request.method),
      }
    }
  end
  local method_succeeded, result = xpcall(method, debug.traceback, table.unpack(request.params))
  if method_succeeded then
    return {
      result = result,
      -- _preserve_table is a throwaway key to prevent
      -- crazy lua empty object == empty array shenanigans
      -- in the event that result == nil/null
      _preserve_table = true
    }
  else
    return {
      error = {
        code = 500,
        message = string.format('Error during "%s"', request.method),
        data = result
      }
    }
  end
end

local function call_and_serialize_result(request_string)
  return serialize(handle_rpc(request_string))
end

local function call_and_handle_unhandled_errors(request_string)
  local suceeded, response = xpcall(
    call_and_serialize_result,
    debug.traceback,
    request_string)
  if suceeded then
    return response
  else
    return serialize({
      error = {
        code = 500,
        message = 'Unhandled error',
        data = response
      }
    })
  end
end

local function call(request_string)
  local response = call_and_handle_unhandled_errors(request_string)
  print(string.format('terraform-crud-api response: %s', response))
  return response
end

remote.add_interface('terraform-crud-api', {call = call})

script.on_init(function()
  global.resource_db = {}
end)
