package net

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	StravaActivitiesEndpoint = "https://www.strava.com/api/v3/athlete/activities"
	MaxAllowedPerPage        = 200
)


// StravaActivity represents https://developers.strava.com/docs/reference/#api-models-SummaryActivity
// Thank you https://mholt.github.io/json-to-go/
type StravaActivity struct {
	ResourceState int `json:"resource_state"`
	Athlete       struct {
		ID            int `json:"id"`
		ResourceState int `json:"resource_state"`
	} `json:"athlete"`
	Name               string      `json:"name"`
	Distance           float64     `json:"distance"`
	MovingTime         int         `json:"moving_time"`
	ElapsedTime        int         `json:"elapsed_time"`
	TotalElevationGain float64     `json:"total_elevation_gain"`
	Type               string      `json:"type"`
	WorkoutType        interface{} `json:"workout_type"`
	ID                 int64       `json:"id"`
	ExternalID         string      `json:"external_id"`
	UploadID           int64       `json:"upload_id"`
	StartDate          time.Time   `json:"start_date"`
	StartDateLocal     time.Time   `json:"start_date_local"`
	Timezone           string      `json:"timezone"`
	UtcOffset          float64     `json:"utc_offset"`
	StartLatlng        interface{} `json:"start_latlng"`
	EndLatlng          interface{} `json:"end_latlng"`
	LocationCity       interface{} `json:"location_city"`
	LocationState      interface{} `json:"location_state"`
	LocationCountry    string      `json:"location_country"`
	AchievementCount   int         `json:"achievement_count"`
	KudosCount         int         `json:"kudos_count"`
	CommentCount       int         `json:"comment_count"`
	AthleteCount       int         `json:"athlete_count"`
	PhotoCount         int         `json:"photo_count"`
	Map                struct {
		ID              string      `json:"id"`
		SummaryPolyline interface{} `json:"summary_polyline"`
		ResourceState   int         `json:"resource_state"`
	} `json:"map"`
	Trainer              bool    `json:"trainer"`
	Commute              bool    `json:"commute"`
	Manual               bool    `json:"manual"`
	Private              bool    `json:"private"`
	Flagged              bool    `json:"flagged"`
	GearID               string  `json:"gear_id"`
	FromAcceptedTag      bool    `json:"from_accepted_tag"`
	AverageSpeed         float64 `json:"average_speed"`
	MaxSpeed             float64 `json:"max_speed"`
	AverageCadence       float64 `json:"average_cadence"`
	AverageWatts         float64 `json:"average_watts"`
	WeightedAverageWatts int     `json:"weighted_average_watts"`
	Kilojoules           float64 `json:"kilojoules"`
	DeviceWatts          bool    `json:"device_watts"`
	HasHeartrate         bool    `json:"has_heartrate"`
	AverageHeartrate     float64 `json:"average_heartrate"`
	MaxHeartrate         float64 `json:"max_heartrate"`
	MaxWatts             int     `json:"max_watts"`
	PrCount              int     `json:"pr_count"`
	TotalPhotoCount      int     `json:"total_photo_count"`
	HasKudoed            bool    `json:"has_kudoed"`
	SufferScore          int     `json:"suffer_score"`
}

// GetAllStravaActivitiesFromOnePage makes an HTTP GET request to get all user activities from one page and
// returns the data and nil if successful. Because only a maximum of 200 activities can be requested
// at once, one may need to run this function multiple times while incrementing pageNum from 1 to n, until the response
// data equals "[]", so the un-marshaled slice of type StravaActivity is empty.
// The access token must have read_all scope. Returns nil and the error if one occurs.
func GetAllStravaActivitiesFromOnePage(accessToken string, pageNum int) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, StravaActivitiesEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Authorization": []string{"Bearer " + accessToken},
	}
	q := req.URL.Query()
	q.Add("page", strconv.Itoa(pageNum))
	q.Add("per_page", strconv.Itoa(MaxAllowedPerPage))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
