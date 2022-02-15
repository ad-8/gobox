package strava

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ad-8/gobox/net"
	goboxtime "github.com/ad-8/gobox/time"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	StravaOAuth                 = "https://www.strava.com/oauth/token"
	StravaActivitiesEndpoint    = "https://www.strava.com/api/v3/athlete/activities"
	MaxActivitiesAllowedPerPage = 200 // to minimize number of requests, got only 1000 per day
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

	body, statusCode, err := net.MakePOSTRequest(StravaOAuth, params)

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

// GetAllActivities queries the Strava API for all activities of a user using an access token
// and returns the data structured in a *SafeMap.
// The access token must have read_all scope.
func GetAllActivities(info TokenInfo, maxPages int) (*SafeMap, error) {
	allActivities := SafeMap{V: make(map[int][]StravaActivity)}

	var wg sync.WaitGroup
	for pageNum := 1; pageNum <= maxPages; pageNum++ {
		fmt.Printf("starting goroutine for page %d\n", pageNum)
		wg.Add(1)
		go getPage(info.AccessToken, pageNum, &allActivities, &wg)
	}
	wg.Wait()

	return &allActivities, nil
}

// getPage queries the Strava API for all activities on the specified page.
// This function runs in multiple goroutines and writes its result to the provided SafeMap m.
func getPage(accessToken string, pageNum int, m *SafeMap, wg *sync.WaitGroup) {
	start := time.Now()
	fmt.Printf("inside getPage %d\n", pageNum)

	var activitiesOnPage []StravaActivity
	body, _ := requestActivitiesFromPage(accessToken, pageNum)

	if string(body) == "[]" {
		fmt.Printf("goroutine %d done in %v -- response body is %v\n", pageNum, time.Since(start), string(body))
		wg.Done()
		return
	}

	if err := json.Unmarshal(body, &activitiesOnPage); err != nil {
		fmt.Println(string(body))
		log.Fatal(err)
	}

	m.Add(activitiesOnPage, pageNum)

	fmt.Printf("goroutine %d done in %v\n", pageNum, time.Since(start))
	wg.Done()
}

// requestActivitiesFromPage makes an HTTP GET request to get all user activities from one page and
// returns the data and nil if successful. Because only a maximum of 200 activities can be requested
// at once, one may need to run this function multiple times while incrementing pageNum from 1 to n, until the response
// data equals "[]", so the un-marshaled slice of type StravaActivity is empty.
func requestActivitiesFromPage(accessToken string, pageNum int) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, StravaActivitiesEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + accessToken},
	}
	q := req.URL.Query()
	q.Add("page", strconv.Itoa(pageNum))
	q.Add("per_page", strconv.Itoa(MaxActivitiesAllowedPerPage))
	req.URL.RawQuery = q.Encode()

	resp, _, err := net.MakeGETRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
