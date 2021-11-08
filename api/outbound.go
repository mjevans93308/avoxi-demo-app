package api

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/savaki/geoip2"
	"inet.af/netaddr"

	"github.com/mjevans93308/avoxi-demo-app/config"
)

// declare once for reuse
var api *geoip2.Api

func init() {
	userId := os.Getenv(config.Maxmind_User_Id)
	licenseKey := os.Getenv(config.Maxind_License_Key)
	api = geoip2.New(userId, licenseKey)
}

func SendOutboundRequest(ip_address *netaddr.IP) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	resp, err := api.Country(ctx, ip_address.String())
	if err != nil {
		logger.Error(err)
		return "", err
	}

	if len(resp.Country.Names) == 0 {
		err := errors.New("geoIP lookup processed successfully but no information returned")
		logger.Error(err)
		return "", err
	}

	name := resp.Country.Names[config.English]
	if name == "" {
		err := errors.New("geoIP lookup processed successfully but no english name found")
		logger.Error(err)
		return "", err
	}

	return name, nil
}
