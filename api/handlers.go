package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjevans93308/avoxi-demo-app/config"
	"github.com/spf13/viper"
	"inet.af/netaddr"
)

var Ip_Country_Mapping = make(map[netaddr.IP]string)

// aliveHandler serves as healthcheck endpoint
func (a *App) aliveHandler(c *gin.Context) {
	logger.Info("Received call to IsAliveHandler")
	c.String(http.StatusOK, "It's...ALIVE!!!")
}

// informHandler is gin's polling for aliveness
func (a *App) informHandler(c *gin.Context) {
	logger.Info("Received call to inform")
	c.Status(http.StatusOK)
}

// informHandler is gin's polling for aliveness
func (a *App) teapotHandler(c *gin.Context) {
	logger.Info("Am I a teapot?")
	c.Status(http.StatusTeapot)
}

type payload struct {
	Ip_address    string   `json:"ip_address"`
	Country_names []string `json:"country_names"`
}

// CheckGeoLocation is our main endpoint
// We expect a POST from an outside endpoint with basic auth header and a json payload
// format:
// {
//     "ip_address": "X.X.X.X",
//     "country_names": [
// 	"Mexico",
// 	"Canada",
// ]
// }
// This endpoint will return a 302 if the IP was mapped to a country in the country_names list
// or a 404 if the IP lookup fails
// Other error codes are 400 if the json payload is malformed or the IP does not conform to IPv4 or IPv6 standards
func (a *App) CheckGeoLocation(c *gin.Context) {
	logger.Info("Received call to check geolocation")

	payload := validatePayload(c)

	// netaddr.ParseIP() handles both IPv4 and IPv6 addr schemas
	// https://pkg.go.dev/inet.af/netaddr#ParseIP
	ip_address, err := netaddr.ParseIP(payload.Ip_address)
	if err != nil {
		logger.Error("Could not parse IP from payload to native IP object")
		c.AbortWithError(http.StatusBadRequest, errors.New("could not complete request due to malformed IP address"))
	}

	var mapped_country_name string
	// Check internal mapping table to see if we've received a request for this
	if Ip_Country_Mapping[ip_address] != "" {
		mapped_country_name = Ip_Country_Mapping[ip_address]
	} else {
		country_name, err := a.Outbound.sendOutboundRequest(&ip_address)
		if err != nil {
			logger.Error("Unable to complete lookup request")
			c.AbortWithStatusJSON(http.StatusNotFound, err)
		}
		mapped_country_name = country_name
		Ip_Country_Mapping[ip_address] = country_name
	}

	// check if country_name resides in list of supplied countries
	found := false
	for _, c_name := range payload.Country_names {
		if c_name == mapped_country_name {
			found = true
			break
		}
	}

	// respond with either a 302 or 404
	// log the IP with the request success/fail if we are in test mode
	if found {
		if viper.GetString(config.Environment) == config.TestEnv {
			logger.Infof("Found country match for ip: %s", ip_address)
		} else {
			logger.Info("Found country match")
		}
		c.Status(http.StatusFound)
	} else {
		if viper.GetString(config.Environment) == config.TestEnv {
			logger.Infof("Found no country match for ip: %s", ip_address)
		} else {
			logger.Info("No match found")
		}
		c.Status(http.StatusNotFound)
	}
}

// validatePayload is a helper function to check that the JSON payload sent in the request was not malformed
// and all the fieds necessary are present and have values
// if any validation fails, we produce a 400 Bad Request response and abort the handler chain
// otherwise return a payload struct
func validatePayload(c *gin.Context) *payload {
	logger.Info("Received call to check geolocation")

	if c.Request.Header.Get("Content-Type") != "application/json" {
		logger.Warn("Received request without Content-Type=application/json header. Will attempt to parse anyways.")
	}

	var payload payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
	}

	if payload.Ip_address == "" || payload.Country_names == nil {
		logger.Error("Could not complete request due to missing parameters")
		c.AbortWithError(http.StatusBadRequest, errors.New("could not complete request due to missing parameters"))
	}

	return &payload
}
