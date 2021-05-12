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
  end
  -- Wrapped with () to prevent gsub returning a multival
  return (string.gsub(game.table_to_json(value), '{}', '[]'))
end

exports = {
  ping = function() return '"pong"' end,

  read = function(resource_type, query_string)
    local query = game.json_to_table(query_string)
    return serialize(resources[resource_type].read(query))
  end,

  create = function(resource_type, create_config_string)
    local create_config = game.json_to_table(create_config_string)
    return serialize(resources[resource_type].create(create_config))
  end,

  update = function(resource_type, resource_id, update_config_string)
    local update_config = game.json_to_table(update_config_string)
    return serialize(resources[resource_type].update(resource_id, update_config))
  end,

  delete = function(resource_type, resource_id)
    return serialize(resources[resource_type].delete(resource_id))
  end,
}

remote.add_interface('terraform-crud-api', exports)

script.on_init(function()
  global.resource_db = {}
end)
