package models

import (
	"inet.af/netaddr"
)

type GeolocationPackage struct {
	Ip_address    netaddr.IP
	Country_Names []string
}
