package config

const (
	// basic auth creds
	Basic_Auth_Username = "BASIC_AUTH_USERNAME"
	Basic_Auth_Password = "BASIC_AUTH_PASSWORD"

	// external geoIP lookup api creds
	Maxmind_User_Id    = "MAXMIND_USER_ID"
	Maxind_License_Key = "MAXIND_LICENSE_KEY"

	Environment = "ENVIRONMENT"
	TestEnv     = "test"

	Address = ":8080"

	// routes
	Api_Group = "/api"
	V1_Group  = "/v1"

	Home            = "/"
	Alive           = "/alive"
	Inform          = "/inform"
	Teapot          = "/teapot"
	CheckIPLocation = "/checkiplocation"

	English       = "en"
	Serve_Command = "serve"

	ConfigFileType     = "env"
	ConfigFileLocation = "."
	CobraCommandName   = "avoxi-demo-app"

	GeoliteUrl = "https://geolite.info/geoip/v2.1/country/%s?pretty"
)
