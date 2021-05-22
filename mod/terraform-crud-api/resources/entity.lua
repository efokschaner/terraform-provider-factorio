local resource_db = require('resource_db')

local function table_invert(t)
  local s={}
  for k,v in pairs(t) do
    s[v]=k
  end
  return s
end

local direction_to_string = table_invert(defines.direction)

local function get_direction(direction_name)
  local direction = defines.direction[direction_name]
  if direction == nil then
    error(string.format('Expected valid direction but got "%s"', direction_name))
  end
  return direction
end

local function entity_to_resource(e)
  return {
    unit_number = e.unit_number,
    surface = e.surface.name,
    name = e.name,
    position = e.position,
    direction = direction_to_string[e.direction],
    force = e.force.name
  }
end

return {
  read = function(query)
    local unit_number = query.unit_number
    local entity = resource_db.get('entity', unit_number)
    if entity == nil then
      return nil
    end
    if not entity.valid then
      resource_db.put('entity', unit_number, nil)
      return nil
    end
    return entity_to_resource(entity)
  end,

  create = function(config)
    local surface = game.surfaces[config.surface]
    if surface == nil then
      error(string.format('Could not find surface with id "%s"', config.surface))
    end
    local entity_creation_params = {
      name = config.name,
      position = config.position,
      direction = get_direction(config.direction),
      force = config.force,
      target = nil,
      source = nil,
      fast_replace = false,
      player = nil,
      spill = true,
      raise_built = true,
      create_build_effect_smoke = false,
      spawn_decorations = true,
      move_stuck_players = true,
      item = nil,
    }
    if config.entity_specific_parameters ~= nil then
      for k,v in pairs(config.entity_specific_parameters) do
        entity_creation_params[k] = v
      end
    end
    local e = surface.create_entity(entity_creation_params)
    if e == nil then
      error(string.format('Failed to create "%s"', config.name))
    end
    resource_db.put('entity', e.unit_number, e)
    return entity_to_resource(e)
  end,

  update = function(resource_id, update_config)
    local unit_number = tonumber(resource_id)
    local entity = resource_db.get('entity', unit_number)
    if entity == nil then
      return nil
    end
    if not entity.valid then
      resource_db.put('entity', unit_number, nil)
      return nil
    end
    if update_config.direction ~= nil then
      entity.direction = get_direction(update_config.direction)
    end
    if update_config.force ~= nil then
      entity.force = update_config.force
    end

    return entity_to_resource(entity)
  end,

  delete = function(resource_id)
    local unit_number = tonumber(resource_id)
    local entity = resource_db.get('entity', unit_number)
    if entity == nil then
      return {
        resource_exists = false
      }
    end
    if not entity.valid then
      resource_db.put('entity', unit_number, nil)
      return {
        resource_exists = false
      }
    end
    local destroyed = entity.destroy()
    if destroyed then
      resource_db.put('entity', unit_number, nil)
    end
    return {
      resource_exists = not destroyed
    }
  end,
}