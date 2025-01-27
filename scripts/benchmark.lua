-- scripts/benchmark.lua
--
-- Usage:
--   wrk -t4 -c100 -d10s --script scripts/benchmark.lua http://backend:8080/api/shorten
--

local counter = 0

function init(args)
  math.randomseed(os.time() + tonumber(tostring({}):sub(8)))
end

local function randomDomainName()
  local letters = "abcdefghijklmnopqrstuvwxyz"
  local domainLength = math.random(5, 50)
  local nameParts = {}

  for i = 1, domainLength do
    local index = math.random(#letters)
    table.insert(nameParts, letters:sub(index, index))
  end

  local tlds = {".com", ".org", ".gov"}
  local tld = tlds[math.random(#tlds)]

  return table.concat(nameParts) .. tld
end

local function randomPath()
  local letters = "abcdefghijklmnopqrstuvwxyz0123456789"
  local pathLength = math.random(10, 100)
  local pathParts = {}

  for i = 1, pathLength do
    local index = math.random(#letters)
    table.insert(pathParts, letters:sub(index, index))
  end

  return table.concat(pathParts)
end

function request()
  counter = counter + 1

  local domainName = randomDomainName()
  local path = randomPath() .. tostring(counter)

  local body = string.format('{"url": "https://%s/%s"}', domainName, path)
  local headers = {}
  headers["Content-Type"] = "application/json"

  return wrk.format("POST", nil, headers, body)
end
