package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"inet.af/netaddr"

	"github.com/mjevans93308/avoxi-demo-app/config"
	"github.com/mjevans93308/avoxi-demo-app/models"
)

func (a *App) aliveHandler(c *gin.Context) {
	logger.Info("It's...ALIVE!!!")
	logger.Info(viper.GetString(config.Basic_Auth_Username))
	c.String(http.StatusOK, "It's...ALIVE!!!")
}

type payload struct {
	Ip_address    string `json:"ip_address"`
	Country_names string `json:"country_names"`
}

type response struct {
	Location_check_pass bool   `json:"location_check_pass"`
	Country_name        string `json:"country_name"`
}

func (a *App) CheckGeoLocation(c *gin.Context) {
	logger.Info("Received call to check geolocation")
	var payload payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Error("Could not bind json")
		c.AbortWithError(http.StatusBadRequest, err)
	}

	if payload.Ip_address == "" || payload.Country_names == "" {
		logger.Error("Could not complete request due to missing parameters")
		c.AbortWithError(http.StatusInternalServerError, errors.New("could not complete request due to missing parameters"))
	}

	// unpack payload
	// netaddr.ParseIP() handles both IPv4 and IPv6 addr schemas
	ip_address, err := netaddr.ParseIP(payload.Ip_address)
	if err != nil {
		logger.Error("Could not parse IP from payload to native IP object")
		c.AbortWithError(http.StatusInternalServerError, errors.New("could not complete request due to missing parameters"))
	}

	geolocation := models.GeolocationPackage{
		Ip_address:    ip_address,
		Country_Names: strings.Split(payload.Country_names, ","),
	}

	country_name, err := SendOutboundRequest(&geolocation.Ip_address)
	if err != nil {
		logger.Error("Unable to complete lookup request")
		c.AbortWithStatusJSON(http.StatusNotFound, err)
	}

	// check if country_name resides in list of supplied countries
	resp := response{
		Location_check_pass: false,
		Country_name:        "",
	}
	for _, c_name := range geolocation.Country_Names {
		if c_name == country_name {
			resp.Location_check_pass = true
			resp.Country_name = c_name
			break
		}
	}

	// respond with either a 302 or 404 and a json body
	if resp.Location_check_pass {
		c.JSON(http.StatusFound, resp)
	} else {
		c.JSON(http.StatusNotFound, resp)
	}
}
