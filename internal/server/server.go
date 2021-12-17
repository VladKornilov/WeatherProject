package server

import (
	"../api"
	"../entities"
	"../logger"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"os"
	"time"
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
	if logger.LogErr(err) {
		return
	}
}

func SetupRedis(url string, port string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     url + ":" + port,
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

	var weather entities.Weather

	weatherJson, err := redisClient.Get(ctx, city).Bytes()
	if logger.LogErr(err) {
		weather, err = api.GetWeather(city)
		json, err := json.Marshal(weather)
		redisClient.Set(ctx, city, json, time.Minute)
		if logger.LogErr(err) {
			return
		}
	} else {
		err = json.Unmarshal(weatherJson, &weather)
		if logger.LogErr(err) {
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
