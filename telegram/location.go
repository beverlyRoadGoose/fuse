package telegram // import "heytobi.dev/fuse/telegram"

// Venue represents a venue.
// See https://core.telegram.org/bots/api#venue
type Venue struct {
	Location        *Location `json:"location"`
	Title           string    `json:"title"`
	Address         string    `json:"address"`
	FoursquareID    string    `json:"foursquare_id"`
	FoursquareType  string    `json:"foursquare_type"`
	GooglePlaceID   string    `json:"google_place_id"`
	GooglePlaceType string    `json:"google_place_type"`
}

// Location represents a point on the map.
// See https://core.telegram.org/bots/api#location
type Location struct {
	Longitude            float32 `json:"longitude"`
	Latitude             float32 `json:"latitude"`
	HorizontalAccuracy   float32 `json:"horizontal_accuracy"`
	LivePeriod           int     `json:"live_period"`
	Heading              int     `json:"heading"`
	ProximityAlertRadius int     `json:"proximity_alert_radius"`
}
