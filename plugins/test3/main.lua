local resp = get("https://www.timeapi.io/api/Time/current/zone?timeZone=Europe/Amsterdam")
print(resp)

local obj = json_decode(resp)
print(obj.date, obj.time)