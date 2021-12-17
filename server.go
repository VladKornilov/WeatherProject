package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client
var ctx = context.Background()

func StartServer(router *gin.Engine) {
	router.LoadHTMLGlob("html/templates/*.html")
	addPageListeners(router)

	url, _ := os.LookupEnv("SITE_URL")
	port, _ := os.LookupEnv("SITE_PORT")
	redisClient = SetupRedis(url, port)
	err := router.Run(port)
	if LogErr(err) {
		return
	}
}

func SetupRedis(url string, port string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	return client
}

func addPageListeners(router *gin.Engine) {
	router.GET("/weather", handleWeatherRequest)
	router.GET("/", startPage)
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("./html"))))
}

func startPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func handleWeatherRequest(c *gin.Context) {
	city := c.Query("city")

	var weather Weather

	weatherJson, err := redisClient.Get(ctx, city).Bytes()
	if LogErr(err) {
		weather, err = GetWeather(city)
		if LogErr(err) {
			return
		}
		json, err := json.Marshal(weather)
		redisClient.Set(ctx, city, json, time.Minute)
		if LogErr(err) {
			return
		}
	} else {
		LogData("Found info about " + city + " in Redis cache")
		err = json.Unmarshal(weatherJson, &weather)
		if LogErr(err) {
			return
		}
	}

	c.HTML(http.StatusOK, "update", gin.H{
		"Temp":      weather.Main.Temp,
		"Feels":     weather.Main.FeelsLike,
		"Humidity":  weather.Main.Humidity,
		"WindSpeed": weather.Wind.Speed,
	})
}
