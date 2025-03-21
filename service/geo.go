package service

import (
	"context"
	"github.com/projectxpolaris/youphoto/plugins"
	"googlemaps.github.io/maps"
)

type GeoAddress struct {
	Full                string
	Country             string
	AdministrativeArea1 string
	AdministrativeArea2 string
	Locality            string
	Route               string
	StreetNumber        string
	Premise             string
	Address             string
}

func GetGeoToAddress(Lat float64, Lng float64) (*GeoAddress, error) {
	if !plugins.DefaultGeoPlugin.IsEnable() {
		return nil, nil
	}
	client := plugins.DefaultGeoPlugin.Client
	req := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: Lat,
			Lng: Lng,
		},
	}
	//timeoutContext, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	//defer cancelFunc()

	resp, err := client.Geocode(context.Background(), req)
	if err != nil {
		return nil, err

	}
	if len(resp) == 0 {
		return nil, nil
	}
	addrResult := resp[0]
	address := &GeoAddress{
		Full: addrResult.FormattedAddress,
	}
	for _, comp := range addrResult.AddressComponents {
		for _, compType := range comp.Types {
			switch compType {
			case "country":
				address.Country = comp.LongName
			case "administrative_area_level_1":
				address.AdministrativeArea1 = comp.LongName
			case "administrative_area_level_2":
				address.AdministrativeArea2 = comp.LongName
			case "locality":
				address.Locality = comp.LongName
			case "route":
				address.Route = comp.LongName
			case "street_number":
				address.StreetNumber = comp.LongName
			case "premise":
				address.Premise = comp.LongName
			}
		}

	}
	return address, nil
}
