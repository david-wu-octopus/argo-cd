local hs = {}
if obj.status ~= nil and obj.status.conditions ~= nil then
  for i, condition in ipairs(obj.status.conditions) do
    if condition.type == "Ready" then
      if condition.status == "True" and condition.reason == "SuccessfulCreateOrUpdate" then
        hs.status = "Healthy"
        hs.message = "RabbitMQ policy ready"
        return hs
      end 

      if condition.status == "False" and condition.reason == "FailedCreateOrUpdate" then
        hs.status = "Degraded"
        hs.message = "RabbitMQ policy failed to be created or updated"
        return hs
      end
    end
  end
end

hs.status = "Unknown"
hs.message = "RabbitMQ policy status is unknown"
return hs