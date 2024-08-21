package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Owm struct {
	apiKey string
}

type OwmDataMain struct {
	Temp    float32 `json:"temp"`
	TempMin float32 `json:"temp_min"`
	TempMax float32 `json:"temp_max"`
}

type OwmData struct {
	Name string      `json:"name"`
	Main OwmDataMain `json:"main"`
}

func NewOwm(apiKey string) *Owm {
	return &Owm{
		apiKey: apiKey,
	}
}

func (o *Owm) Query(location string) (*OwmData, error) {
	url := fmt.Sprintf(
		"http://api.openweathermap.org/data/2.5/weather?q=%s&APPID=%s&units=metric",
		url.QueryEscape(location),
		url.QueryEscape(o.apiKey),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var record OwmData
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		return &record, err
	}

	return &record, nil
}
