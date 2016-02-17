package geocode

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	errLatLngOrPlaceIdRequire = errors.New("Lat,Lng or PlaceId is required")

)


// ReverseGeocode is the process of converting geographic coordinates into a human-readable address
func (p *Service) ReverseGeocode(lat, lng float64) *ReverseGeocodeCall {
	return &ReverseGeocodeCall{
		service: p,
		lat:     lat,
		lng:     lng,
	}
}

type ReverseGeocodeCall struct {
	service *Service

	// The latitude/longitude around which to retrieve place information
	lat, lng float64

	// The place ID of the place for which you wish to obtain the human-readable address
	PlaceId string

	// The language code, indicating in which language the results should be returned, if possible.
	Language string

	// One or more address types, Examples of address types: country, street_address, postal_code
	ResultType []string
	// One or more location types, Specifying a type will restrict the results to this type
	LocationType []string
}

func (n *ReverseGeocodeCall) validate() error {

	if n.lat != 0 || n.lng != 0 {
		return nil
	}
	if n.PlaceId != "" {
		return nil
	}

	return errLatLngOrPlaceIdRequire
}

func (n *ReverseGeocodeCall) Do() (*GeocodeResponse, error) {
	if err := n.validate(); err != nil {
		return nil, err
	}

	searchURL := baseURL + "/json?" + n.query()

	resp, err := n.service.client.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad resp %d: %s", resp.StatusCode, body)
	}

	data := &GeocodeResponse{}
	if err := json.Unmarshal(body, data); err != nil {
		return nil, err
	}

	if data.Status != "OK" {
		return nil, &apiError{
			Status:  data.Status,
			Message: data.ErrorMessage,
		}
	}

	return data, nil
}

func (r *ReverseGeocodeCall) query() string {
	query := make(url.Values)
	query.Add("key", r.service.key)

	if r.lat != 0 || r.lng != 0 {
		query.Add("latlng", fmt.Sprintf("%f,%f", r.lat, r.lng))
	}
	if r.PlaceId != "" {
		query.Add("place_id", r.Language)
	}
	if r.Language != "" {
		query.Add("language", r.Language)
	}

	var resultTypes []string
	for _, t := range r.ResultType {
		resultTypes = append(resultTypes, string(t))
	}
	if len(resultTypes) > 0 {
		query.Add("result_type", strings.Join(resultTypes, "|"))
	}

	var locationTypes []string
	for _, t := range r.LocationType {
		locationTypes = append(locationTypes, string(t))
	}
	if len(locationTypes) > 0 {
		query.Add("location_type", strings.Join(locationTypes, "|"))
	}

	return query.Encode()
}

type GeocodeResponse struct {
	// A list of results matching the query
	Results []GeocodeDetail `json:"results"`
	// Contains debugging information to help you track down why the request failed
	Status string `json:"status"`
	// More detailed information about the reasons behind the given status code.
	ErrorMessage string `json:"error_message,omitempty"`
	// A set of attributions about this listing which must be displayed to the user.
	HTMLAttributions []string `json:"html_attributions"`
}

// An AddressComponent is a component used to compose a given address
type AddressComponent struct {
	// An array indicating the type of the address component.
	Types []string `json:"types"`
	// The full text description or name of the address component.
	LongName string `json:"long_name"`
	// An abbreviated textual name for the address component, if available. For example, an address component for the state of Alaska may have a long_name of "Alaska" and a short_name of "AK" using the 2-letter postal abbreviation.
	ShortName string `json:"short_name"`
}

// LatLng contains the geocoded latitude and longitude value for a place.
type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Geometry contains a place's location
type Geometry struct {
	Location LatLng `json:"location"`
	LocationType string `json:"location_type"`
	Viewport LatLngArea `json:"view_port"`
	Bounds LatLngArea `json:"bounds,omitempty"`
}

// LatLngArea contains southwest and northeast corner of the viewport bounding box
type LatLngArea struct {
	SoutWest LatLng `json:"southwest"`
	NorthEase LatLng `json:"northeast"`
}

// FeatureType is a feature type describing a place.
type FeatureType string

// GeocodeDetail is the information returned by a geocode request
// https://developers.google.com/maps/documentation/geocoding/intro#GeocodingResponses
type GeocodeDetail struct {
	// An array of feature types describing the given result.
	Types []FeatureType `json:"types"`
	// A string containing the human-readable address of this place. Often this address is equivalent to the "postal address," which sometimes differs from country to country.
	FormattedAddress string `json:"formatted_address"`
	// An array of separate address components used to compose a given address
	AddressComponents []AddressComponent `json:"address_components"`
	// The place's phone number in its local format.
	FormattedPhoneNumber string `json:"formatted_phone_number"`
	// Geometry contains a place's location
	Geometry Geometry `json:"geometry"`
	// PartialMatch indicates that the geocoder did not return an exact match for the original request
	PartialMatch string `json:"partial_match,omitempty"`
	// A textual identifier that uniquely identifies a place.
	PlaceID string `json:"place_id"`

}
