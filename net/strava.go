package net

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

// GetAllStravaActivitiesFromOnePage makes an HTTP GET request to get all user activities from one page and
// returns the data and nil if successful. Because only a maximum of 200 activities can be requested
// at once, one may need to run this function multiple times while incrementing pageNum from 1 to n, until the data
// returned equals "[]". Returns nil and the error if one occurs.
func GetAllStravaActivitiesFromOnePage(accessToken string, pageNum int) ([]byte, error) {
	const StravaActivitiesEndpoint = "https://www.strava.com/api/v3/athlete/activities"
	const MaxAllowedPerPage = 200

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
