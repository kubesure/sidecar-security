package routing

//Config defines configuration for routing
type Config struct {
	Method string
}

//Configurations define routing config for reverse proxies
var Configurations = []Config{
	Config{
		Method: "GET",
	},
	Config{
		Method: "POST",
	},
	Config{
		Method: "PUT",
	},
	Config{
		Method: "PATCH",
	},
	Config{
		Method: "DELETE",
	},
}
