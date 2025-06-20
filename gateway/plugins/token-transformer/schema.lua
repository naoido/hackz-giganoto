return {
  name = "token-transformer",
  fields = {
    { consumer = { type = "foreign", reference = "consumers" } },
    { route = { type = "foreign", reference = "routes" } },
    { service = { type = "foreign", reference = "services" } },
  },
}