package plex

import "encoding/xml"

type PlexResourceDeviceConnection struct {
	Protocol string `xml:"protocol,attr" json:"protocol,omitempty"`
	Address  string `xml:"address,attr" json:"address,omitempty"`
	Port     string `xml:"port,attr" json:"port,omitempty"`
	URI      string `xml:"uri,attr" json:"uri,omitempty"`
	Local    string `xml:"local,attr" json:"local,omitempty"`
}

type PlexResourceDevice struct {
	Name          string                         `xml:"name,attr" json:"name,omitempty"`
	Product       string                         `xml:"product,attr" json:"product,omitempty"`
	Platform      string                         `xml:"platform,attr" json:"platform,omitempty"`
	Provides      string                         `xml:"provides,attr" json:"provides,omitempty"`
	Owned         string                         `xml:"owned,attr" json:"owned,omitempty"`
	PublicAddress string                         `xml:"publicAddress,attr" json:"publicaddress,omitempty"`
	AccessToken   string                         `xml:"accessToken,attr" json:"accessToken,omitempty"`
	Connection    []PlexResourceDeviceConnection `xml:"Connection" json:"connection,omitempty"`
}

type PlexResourcesResponse struct {
	XMLName xml.Name             `xml:"MediaContainer" json:"mediacontainer,omitempty"`
	Size    string               `xml:"size,attr" json:"size,omitempty"`
	Devices []PlexResourceDevice `xml:"Device" json:"device,omitempty"`
}
