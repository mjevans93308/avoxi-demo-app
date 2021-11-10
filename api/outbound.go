package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"inet.af/netaddr"

	"github.com/mjevans93308/avoxi-demo-app/config"
	"github.com/mjevans93308/avoxi-demo-app/types"
)

type Outbound struct {
	UserId     string
	LicenseKey string
	Client     *http.Client
}

func initOutboundApi() *Outbound {
	userId := viper.GetString(config.Maxmind_User_Id)
	licenseKey := viper.GetString(config.Maxind_License_Key)
	return &Outbound{
		UserId:     userId,
		LicenseKey: licenseKey,
		Client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (o *Outbound) sendOutboundRequest(ip_address *netaddr.IP) (string, error) {

	if o == nil {
		o = initOutboundApi()
	}

	url := fmt.Sprintf(config.GeoliteUrl, ip_address.String())
	logger.Infof("Making request to %s", url)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(o.UserId, o.LicenseKey)

	resp, err := o.Client.Do(req)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		eError := types.GeoIpError{}
		err := json.NewDecoder(resp.Body).Decode(&eError)
		if err != nil {
			return "", err
		}

		logger.Errorf("Received error when querying geolite: %d - %s", eError.Code, eError.ErrorString)
		return "", errors.New(eError.ErrorString)
	}

	logger.Infof("Received response code %d when querying geolite", resp.StatusCode)

	geoIpResp := types.GeoIpResponse{}
	err = json.NewDecoder(resp.Body).Decode(&geoIpResp)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	return buildName(&geoIpResp)
}

func buildName(geoIpResp *types.GeoIpResponse) (string, error) {

	if len(geoIpResp.Country.Names) == 0 {
		err := errors.New("geoIP lookup processed successfully but no information returned")
		logger.Error(err)
		return "", err
	}

	name := geoIpResp.Country.Names[config.English]
	if name == "" {
		err := errors.New("geoIP lookup processed successfully but no english name found")
		logger.Error(err)
		return "", err
	}

	return name, nil
}
