package types

type GeoIpResponse struct {
	Country Country `json:"country,omitempty"`
}

type Country struct {
	Confidence int               `json:"confidence,omitempty"`
	GeoNameId  int               `json:"geoname_id,omitempty"`
	IsoCode    string            `json:"iso_code,omitempty"`
	Names      map[string]string `json:"names,omitempty"`
}

type GeoIpError struct {
	Code        int    `json:"code,omitempty"`
	ErrorString string `json:"error,omitempty"`
}
