resources = {
  players = {
    read = function()
      result = {}
      for k,v in pairs(game.players) do
        table.insert(result, {name = v.name})
     end
     return result
    end,
  }
}

-- In addition to json serialization, this code also replaces empty objects with empty arrays
-- due to the fact that lua cannot differentiate.
-- This means we cannot have empty objects, however it is easier to add an unused key
-- to an object to avoid empty objects than to add an unused element to every array.
function serialize(value)
  -- Wrapped with () to prevent gsub returning a multival
  return (string.gsub(game.table_to_json(value), "{}", "[]"))
end

exports = {
  ping = function() return "pong" end,

  read = function(resource_type)
    return serialize(resources[resource_type].read())
  end,

  create = function(resource_type)
  end,

  update = function(resource_type)
  end,

  delete = function(resource_type)
  end,
}

remote.add_interface("terraform-crud-api", exports)
