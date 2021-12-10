package api

import (
	"../entities"
	"../logger"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func GetWeather(city string) (entities.Weather, error) {
	rawUrl, exists := os.LookupEnv("API_URL")
	if !exists {
		logger.LogData("No API_URL line in env file")
	}
	rawUrl += "/weather"
	apiKey, exists := os.LookupEnv("API_KEY")
	if !exists {
		logger.LogData("No API_KEY line in env file, create one")
	}

	urlValues := url.Values{}
	urlValues.Add("q", city)
	urlValues.Add("appid", apiKey)
	urlValues.Add("units", "metric")

	url, err := url.Parse(rawUrl + "?" + urlValues.Encode())
	if logger.LogErr(err) {
		return entities.Weather{}, err
	}

	urlStr := url.String()
	resp, err := http.Get(urlStr)
	if logger.LogErr(err) {
		return entities.Weather{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	weather := entities.Weather{}
	err = json.Unmarshal(body, &weather)
	if logger.LogErr(err) {
		return entities.Weather{}, err
	}

	return weather, nil
}
