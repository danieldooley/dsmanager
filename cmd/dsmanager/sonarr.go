package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	sonarrHost   = "dooley-server.local:8989"
	sonarrApiKey = "40347991623e4fc3bddb56575f72e227"

	radarrHost   = "dooley-server.local:7878"
	radarrApiKey = "b3b2121fc53f4401b44009861f446f2f"

	apiEndpoint = "/api/"
)

type targetEnum int

const (
	SONARR targetEnum = iota
	RADARR
)

func generateUrl(relPath string, extraQ map[string]string, t targetEnum) (string, error) {
	u := &url.URL{
		Scheme: "http",
		Path:   apiEndpoint,
	}

	q := u.Query()

	switch t {
	case SONARR:
		u.Host = sonarrHost
		q.Add("apikey", sonarrApiKey)
	case RADARR:
		u.Host = radarrHost
		q.Add("apikey", radarrApiKey)
	default:
		return "", fmt.Errorf("provided targetEnum cannot be handled")
	}

	u, err := u.Parse(relPath)
	if err != nil {
		return "", fmt.Errorf("could not resolve relPath on target: %v", err)
	}

	for k, v := range extraQ {
		q.Add(k, v)
	}

	u.RawQuery = q.Encode()

	return u.String(), nil
}

func makeRequest(relPath string, extraQ map[string]string, t targetEnum, out interface{}) error {
	u, err := generateUrl(relPath, extraQ, t)
	if err != nil {
		return fmt.Errorf("failed to generateUrl: %v", err)
	}

	res, err := http.Get(u)
	if err != nil {
		return fmt.Errorf("failed to get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("non-200 return code: %d", res.StatusCode)
	}

	dec := json.NewDecoder(res.Body)

	err = dec.Decode(out)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return nil
}

func isActive() (bool, error) {
	var sonarrQueue Queue

	err := makeRequest("queue", map[string]string{}, SONARR, &sonarrQueue)
	if err != nil {
		return false, fmt.Errorf("couldn't fetch Sonarr queue: %v", err)
	}

	for _, i := range sonarrQueue {
		if strings.EqualFold(i.Status, "downloading") {
			return true, nil
		}
	}

	var radarrQueue Queue

	err = makeRequest("queue", map[string]string{}, RADARR, &radarrQueue)
	if err != nil {
		return false, fmt.Errorf("couldn't fetch Radarr queue: %v", err)
	}

	for _, i := range radarrQueue {
		if strings.EqualFold(i.Status, "downloading") {
			return true, nil
		}
	}

	return false, nil
}

func nextActivity() (time.Time, error) {

	var sonarrCal Calendar

	err := makeRequest("calendar", map[string]string{
		"start": time.Now().Format(time.RFC3339),
		"end":   time.Now().AddDate(0, 0, 7).Format(time.RFC3339),
	}, SONARR, &sonarrCal)
	if err != nil {
		return time.Time{}, fmt.Errorf("couldn't fetch Sonarr calendar: %v", err)
	}

	var next = time.Date(9999, 0, 0, 0, 0, 0, 0, time.UTC)

	for _, v := range sonarrCal {
		if v.Monitored {
			if v.AirDateUtc.Before(next) {
				next = v.AirDateUtc
			}
		}
	}

	return next, err
}

// Types

type Calendar []struct {
	AirDateUtc time.Time `json:"airDateUtc"`
	Monitored  bool      `json:"monitored"`
}

type Queue []struct {
	Status string `json:"status"`
}
