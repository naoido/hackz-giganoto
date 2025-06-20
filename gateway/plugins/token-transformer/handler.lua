local http = require("resty.http")
local cjson = require("cjson.safe")

local TokenTransformerHandler = {}

TokenTransformerHandler.PRIORITY = 1000
TokenTransformerHandler.VERSION = "0.1.0"

function TokenTransformerHandler:access(conf)
  local auth_header = kong.request.get_header("authorization")
  if not auth_header then
    return kong.response.exit(401, { message = "Authorization header missing" })
  end

  local opaque_token = string.match(auth_header, "^[Bb]earer%s+(.+)")
  if not opaque_token then
    return kong.response.exit(401, { message = "Invalid Authorization header format" })
  end

  local auth_server_url = "http://auth:8000/introspect"

  -- resty.http を正しく使用
  local httpc = http.new()
  httpc:set_timeout(5000) -- 5秒のタイムアウト
  
  local res, err = httpc:request_uri(auth_server_url, {
    method = "POST",
    body = '{"token":"' .. opaque_token .. '"}',
    headers = {
      ["Content-Type"] = "application/json"
    }
  })

  if err then
    kong.log.err("Auth server request failed: ", err)
    return kong.response.exit(500, { message = "Internal Server Error" })
  end

  if res.status ~= 200 then
    kong.log.info("Introspection failed with status: ", res.status)
    return kong.response.exit(401, { message = "Invalid token" })
  end

  local body, dec_err = cjson.decode(res.body)
  if dec_err or not body or not body.jwt then
    kong.log.err("Failed to decode JWT from auth server response: ", dec_err or "JWT field missing")
    return kong.response.exit(500, { message = "Internal Server Error" })
  end

  local internal_jwt = body.jwt

  kong.service.request.set_header("authorization", "Bearer " .. internal_jwt)
end

return TokenTransformerHandler