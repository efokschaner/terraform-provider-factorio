local resource_db = require('resource_db')


return {
  create = function(config)
    local e = game.surfaces[1].create_entity(config)
    resource_db.put('entity', e.unit_number, e)
    return {
      unit_number = tostring(e.unit_number)
    }
  end,

  read = function(query)
    local unit_number = tonumber(query.unit_number)
    local entity = resource_db.get('entity', unit_number)
    if entity == nil then
      return nil
    end
    if entity.valid then
      return {
        unit_number = tostring(entity.unit_number),
      }
    else
      resource_db.put('entity', unit_number, nil)
      return nil
    end
  end
}