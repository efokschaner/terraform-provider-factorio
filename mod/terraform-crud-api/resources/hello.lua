local resource_db = require('resource_db')

local text = {
  {
    'X  X  XXX  X    X     XX ',
    'X  X  X    X    X    X  X',
    'XXXX  XXX  X    X    X  X',
    'X  X  X    X    X    X  X',
    'X  X  XXX  XXX  XXX   XX ',
  },
  {
    'XXX  XXX    XX   X   X',
    'X    X  X  X  X  XX XX',
    'XXX  XXX   X  X  X X X',
    'X    X X   X  X  X   X',
    'X    X  X   XX   X   X',
  },
  {
    'XXX  XXX  XXX   XXX    XX   XXX   XX   XXX   X   X',
    ' X   X    X  X  X  X  X  X  X    X  X  X  X  XX XX',
    ' X   XXX  XXX   XXX   XXXX  XXX  X  X  XXX   X X X',
    ' X   X    X X   X X   X  X  X    X  X  X X   X   X',
    ' X   XXX  X  X  X  X  X  X  X     XX   X  X  X   X',
  },
}

local function table_invert(t)
  local s={}
  for k,v in pairs(t) do
    s[v]=k
  end
  return s
end

local direction_to_string = table_invert(defines.direction)

local function create_belt(spawn_pos, create_as_ghost, direction)
  local entity_creation_params = {
    name = 'transport-belt',
    position = spawn_pos,
    direction = direction,
    force = 'player',
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
  if create_as_ghost then
    entity_creation_params.name = 'entity-ghost'
    -- Specific to 'entity-ghost'
    entity_creation_params.inner_name = 'transport-belt'
    entity_creation_params.expires = false
  end
  return game.surfaces[1].create_entity(entity_creation_params)
end

return {
  create = function(config)
    local create_as_ghost = config.create_as_ghost == nil or config.create_as_ghost
    local direction = defines.direction[config.direction]
    local created_entities = {}
    local spawn_pos = {x=0,y=0}
    local player = game.get_player(1)
    if player ~= nil then
      spawn_pos.x = player.position.x - 13
      spawn_pos.y = player.position.y - 10
    end
    for _, word in ipairs(text) do
      for _, line in ipairs(word) do
        local x_orig = spawn_pos.x
        for i = 1, #line do
          local c = line:sub(i,i)
          if c == 'X' then
            local spawn_pos = {spawn_pos.x, spawn_pos.y}
            local belt = create_belt(spawn_pos, create_as_ghost, direction)
            table.insert(created_entities, belt)
          end
          spawn_pos.x = spawn_pos.x + 1
        end
        spawn_pos.x = x_orig
        spawn_pos.y = spawn_pos.y + 1
      end
      spawn_pos.y = spawn_pos.y + 3
    end
    -- use first belt id as our id
    local first_belt = created_entities[1]
    resource_db.put("hello", first_belt.unit_number, created_entities)
    return {id = tostring(first_belt.unit_number)}
  end,

  read = function(query)
    local belts = resource_db.get("hello", tonumber(query.id))
    if belts == nil then
      return nil
    end
    local first_valid_entity = nil
    for _, belt in ipairs(belts) do
      if belt.valid then
        first_valid_entity = belt
        break
      end
    end
    if first_valid_entity == nil then
      return nil
    end
    return {
      id = query.id,
      direction = direction_to_string[first_valid_entity.direction],
    }
  end,

  update = function(resource_id, update_config)
    local belts = resource_db.get("hello", tonumber(resource_id))
    if belts == nil then
      return nil
    end
    if update_config.direction then
      for _, belt in ipairs(belts) do
        if belt.valid then
          belt.direction = defines.direction[update_config.direction]
        end
      end
    end
    return nil
  end,

  delete = function(resource_id)
    local belts = resource_db.get('hello', tonumber(resource_id))
    if belts == nil then
      return {
        is_deleted = true
      }
    end
    local all_succeeded = true
    for _, belt in ipairs(belts) do
      all_succeeded = all_succeeded and belt.destroy()
    end
    if all_succeeded then
      resource_db.put('hello', tonumber(resource_id), nil)
    end
    return {
      is_deleted = all_succeeded
    }
  end,
}