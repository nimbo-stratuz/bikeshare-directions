package models

import (
	"errors"
	"net/http"
	"time"
)

// FromTo represents incoming requests to /v1/directions. Example:
// {
// 	 "from": "Ljubljana, Faculty of Computer and Information Science",
// 	 "to": "Medvode"
// }
type FromTo struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

// Bind ensures both fields are set in a FromTo
func (ft *FromTo) Bind(r *http.Request) error {
	if ft.From == "" {
		return errors.New("Missing 'from' (start) field")
	}
	if ft.To == "" {
		return errors.New("Missing 'to' (destination) field")
	}
	return nil
}

// DirectionsRequest is sent to the MapQuest API
type DirectionsRequest struct {
	Locations []string                 `json:"locations"`
	Options   DirectionsRequestOptions `json:"options"`
}

// DirectionsRequestOptions is the Options part of DirectionsRequest
// Unused fields are commented out.
type DirectionsRequestOptions struct {
	// Avoids               []string `json:"avoids"`
	// AvoidTimedConditions bool     `json:"avoidTimedConditions"`
	// DoReverseGeocode     bool     `json:"doReverseGeocode"`
	// ShapeFormat          string   `json:"shapeFormat"`
	// Generalize           int      `json:"generalize"`
	RouteType string `json:"routeType"`
	// TimeType             int      `json:"timeType"`
	// Locale               string   `json:"locale"`
	Unit string `json:"unit"`
	// EnhancedNarrative    bool     `json:"enhancedNarrative"`
	// DrivingStyle         int      `json:"drivingStyle"`
	// HighwayEfficiency    int      `json:"highwayEfficiency"`
}

// Directions is recieved as a response from the MapQuest API
// Unused fields are commented out.
type Directions struct {
	Route struct {
		// 	HasTollRoad       bool          `json:"hasTollRoad"`
		// 	ComputedWaypoints []interface{} `json:"computedWaypoints"`
		// 	FuelUsed          float64       `json:"fuelUsed"`
		// 	Shape             struct {
		// 		ManeuverIndexes []int     `json:"maneuverIndexes"`
		// 		ShapePoints     []float64 `json:"shapePoints"`
		// 		LegIndexes      []int     `json:"legIndexes"`
		// 	} `json:"shape"`
		// 	HasUnpaved  bool `json:"hasUnpaved"`
		// 	HasHighway  bool `json:"hasHighway"`
		// 	RealTime    int  `json:"realTime"`
		// 	BoundingBox struct {
		// 		Ul struct {
		// 			Lng float64 `json:"lng"`
		// 			Lat float64 `json:"lat"`
		// 		} `json:"ul"`
		// 		Lr struct {
		// 			Lng float64 `json:"lng"`
		// 			Lat float64 `json:"lat"`
		// 		} `json:"lr"`
		// 	} `json:"boundingBox"`
		Distance float64 `json:"distance"`
		// 	Time               int     `json:"time"`
		// 	LocationSequence   []int   `json:"locationSequence"`
		// 	HasSeasonalClosure bool    `json:"hasSeasonalClosure"`
		// 	SessionID          string  `json:"sessionId"`
		Locations []struct {
			LatLng struct {
				Lng float64 `json:"lng"`
				Lat float64 `json:"lat"`
			} `json:"latLng"`
			AdminArea1 string `json:"adminArea1"`
			// 		AdminArea1Type     string `json:"adminArea1Type"`
			AdminArea3 string `json:"adminArea3"`
			// 		AdminArea3Type     string `json:"adminArea3Type"`
			AdminArea4 string `json:"adminArea4"`
			// 		AdminArea4Type string `json:"adminArea4Type"`
			AdminArea5 string `json:"adminArea5"`
			// 		AdminArea5Type string `json:"adminArea5Type"`
			// 		Street         string `json:"street"`
			Type string `json:"type"`
			// 		DisplayLatLng  struct {
			// 			Lng float64 `json:"lng"`
			// 			Lat float64 `json:"lat"`
			// 		} `json:"displayLatLng"`
			// 		LinkID             int    `json:"linkId"`
			// 		PostalCode         string `json:"postalCode"`
			// 		SideOfStreet       string `json:"sideOfStreet"`
			// 		DragPoint          bool   `json:"dragPoint"`
			// 		GeocodeQuality     string `json:"geocodeQuality"`
			// 		GeocodeQualityCode string `json:"geocodeQualityCode"`
		} `json:"locations"`
		// 	HasCountryCross bool `json:"hasCountryCross"`
		Legs []struct {
			// 		HasTollRoad        bool            `json:"hasTollRoad"`
			// 		Index              int             `json:"index"`
			// 		RoadGradeStrategy  [][]interface{} `json:"roadGradeStrategy"`
			// 		HasHighway         bool            `json:"hasHighway"`
			// 		HasUnpaved         bool            `json:"hasUnpaved"`
			// 		Distance           float64         `json:"distance"`
			// 		Time               int             `json:"time"`
			// 		OrigIndex          int             `json:"origIndex"`
			// 		HasSeasonalClosure bool            `json:"hasSeasonalClosure"`
			// 		OrigNarrative      string          `json:"origNarrative"`
			// 		HasCountryCross    bool            `json:"hasCountryCross"`
			// 		FormattedTime      string          `json:"formattedTime"`
			// 		DestNarrative      string          `json:"destNarrative"`
			// 		DestIndex          int             `json:"destIndex"`
			Maneuvers []struct {
				// 			Signs         []interface{} `json:"signs"`
				// 			Index         int           `json:"index"`
				// 			ManeuverNotes []interface{} `json:"maneuverNotes"`
				// 			Direction     int           `json:"direction"`
				Narrative string `json:"narrative"`
				// 			IconURL       string        `json:"iconUrl"`
				// 			Distance      float64       `json:"distance"`
				// 			Time          int           `json:"time"`
				// 			LinkIds       []interface{} `json:"linkIds"`
				// 			Streets       []string      `json:"streets"`
				// 			Attributes    int           `json:"attributes"`
				// 			TransportMode string        `json:"transportMode"`
				// 			FormattedTime string        `json:"formattedTime"`
				// 			DirectionName string        `json:"directionName"`
				// 			MapURL        string        `json:"mapUrl,omitempty"`
				StartPoint struct {
					Lng float64 `json:"lng"`
					Lat float64 `json:"lat"`
				} `json:"startPoint"`
				// 			TurnType int `json:"turnType"`
			} `json:"maneuvers"`
			// 		HasFerry bool `json:"hasFerry"`
		} `json:"legs"`
		// 	FormattedTime string `json:"formattedTime"`
		// 	RouteError    struct {
		// 		Message   string `json:"message"`
		// 		ErrorCode int    `json:"errorCode"`
		// 	} `json:"routeError"`
		// 	Options struct {
		// 		MustAvoidLinkIds           []interface{} `json:"mustAvoidLinkIds"`
		// 		DrivingStyle               int           `json:"drivingStyle"`
		// 		CountryBoundaryDisplay     bool          `json:"countryBoundaryDisplay"`
		// 		Generalize                 int           `json:"generalize"`
		// 		NarrativeType              string        `json:"narrativeType"`
		// 		Locale                     string        `json:"locale"`
		// 		AvoidTimedConditions       bool          `json:"avoidTimedConditions"`
		// 		DestinationManeuverDisplay bool          `json:"destinationManeuverDisplay"`
		// 		EnhancedNarrative          bool          `json:"enhancedNarrative"`
		// 		FilterZoneFactor           int           `json:"filterZoneFactor"`
		// 		TimeType                   int           `json:"timeType"`
		// 		MaxWalkingDistance         int           `json:"maxWalkingDistance"`
		// 		RouteType                  string        `json:"routeType"`
		// 		TransferPenalty            int           `json:"transferPenalty"`
		// 		StateBoundaryDisplay       bool          `json:"stateBoundaryDisplay"`
		// 		WalkingSpeed               int           `json:"walkingSpeed"`
		// 		MaxLinkID                  int           `json:"maxLinkId"`
		// 		ArteryWeights              []interface{} `json:"arteryWeights"`
		// 		TryAvoidLinkIds            []interface{} `json:"tryAvoidLinkIds"`
		// 		Unit                       string        `json:"unit"`
		// 		RouteNumber                int           `json:"routeNumber"`
		// 		ShapeFormat                string        `json:"shapeFormat"`
		// 		ManeuverPenalty            int           `json:"maneuverPenalty"`
		// 		UseTraffic                 bool          `json:"useTraffic"`
		// 		ReturnLinkDirections       bool          `json:"returnLinkDirections"`
		// 		AvoidTripIds               []interface{} `json:"avoidTripIds"`
		// 		Manmaps                    string        `json:"manmaps"`
		// 		HighwayEfficiency          int           `json:"highwayEfficiency"`
		// 		SideOfStreetDisplay        bool          `json:"sideOfStreetDisplay"`
		// 		CyclingRoadFactor          int           `json:"cyclingRoadFactor"`
		// 		UrbanAvoidFactor           int           `json:"urbanAvoidFactor"`
		// 	} `json:"options"`
		// 	HasFerry bool `json:"hasFerry"`
	} `json:"route"`
	Info struct {
		Copyright struct {
			Text         string `json:"text"`
			ImageURL     string `json:"imageUrl"`
			ImageAltText string `json:"imageAltText"`
		} `json:"copyright"`
		Statuscode int           `json:"statuscode"`
		Messages   []interface{} `json:"messages"`
	} `json:"info"`
}

// Render ...
func (dirs *Directions) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type Bicycle struct {
	Available bool      `json:"available"`
	DateAdded time.Time `json:"dateAdded"`
	ID        int       `json:"id"`
	Location  struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	OwnerID       int    `json:"ownerId"`
	SmartLockUUID string `json:"smartLockUUID"`
}

// Render ...
func (b *Bicycle) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type DirectionsWithBicycle struct {
	Bicycle    *Bicycle    `json:"bicycle"`
	Directions *Directions `json:"directions"`
}

// Render ...
func (dwb *DirectionsWithBicycle) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
