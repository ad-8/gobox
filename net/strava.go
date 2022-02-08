package net

import (
	"encoding/json"
	"errors"
	"fmt"
	goboxtime "github.com/ad-8/gobox/time"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	StravaOAuth              = "https://www.strava.com/oauth/token"
	StravaActivitiesEndpoint = "https://www.strava.com/api/v3/athlete/activities"
	MaxAllowedPerPage        = 200
)

// TokenInfo represents the response that contains information about the Strava access token.
// Thank you https://mholt.github.io/json-to-go/
type TokenInfo struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	ExpiresHours int    `json:"-"`
	ExpiresMin   int    `json:"-"`
	ExpiresSec   int    `json:"-"`
}

func (t *TokenInfo) ParseTime() {
	simpleTime, err := goboxtime.SecondsToHrsMinSec(t.ExpiresIn)
	if err != nil {
		log.Fatal(err)
	}
	t.ExpiresHours = simpleTime.H
	t.ExpiresMin = simpleTime.M
	t.ExpiresSec = simpleTime.S
}

func (t *TokenInfo) Print() {
	//fmt.Printf("access token = %s\n\n", t.AccessToken)
	fmt.Printf("the token expires in %02d:%02d:%02d (will be automatically refreshed)\n\n", t.ExpiresHours, t.ExpiresMin, t.ExpiresSec)
}

// NewTokenInfo gets information about the access token - because the access token expires every 6 hours - and returns
// a *TokenInfo and nil if successful. Returns nil and the error if one occurs.
func NewTokenInfo(clientId, clientSecret, refreshToken string) (*TokenInfo, error) {
	params := map[string]interface{}{
		"client_id":     clientId,
		"client_secret": clientSecret,
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
	}

	body, statusCode, err := MakePOSTRequest(StravaOAuth, params)

	if err != nil {
		return nil, err
	}

	tokenInfo := new(TokenInfo)
	if err := json.Unmarshal(body, tokenInfo); err != nil {
		return nil, errors.New(
			fmt.Sprintf("Error: %v\nstatus code is %d. cannot parse this response:\n%v\n",
				err, statusCode, string(body)))
	}

	tokenInfo.ParseTime()

	return tokenInfo, nil
}

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
