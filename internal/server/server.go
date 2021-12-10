package server

import (
	"../api"
	"../logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func StartServer(router *gin.Engine) {
	router.LoadHTMLGlob("html/templates/*.html")
	addPageListeners(router)

	port, _ := os.LookupEnv("SITE_PORT")
	err := router.Run(port)
	if logger.LogErr(err) {
		return
	}
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
	weather, _ := api.GetWeather(city)

	c.HTML(http.StatusOK, "update", gin.H{
		"Temp":      weather.Main.Temp,
		"Feels":     weather.Main.FeelsLike,
		"Humidity":  weather.Main.Humidity,
		"WindSpeed": weather.Wind.Speed,
	})
}
