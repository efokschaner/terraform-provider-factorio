return {
  get = function(type, id)
    local type_query = global.resource_db[type]
    if type_query == nil then
      return nil
    end
    return type_query[id]
  end,

  put = function(type, id, resource)
    local type_query = global.resource_db[type]
    if type_query == nil then
      type_query = {}
      global.resource_db[type] = type_query
    end
    type_query[id] = resource
  end,
}