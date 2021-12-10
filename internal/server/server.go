package server

import (
	"../api"
	"../logger"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

var app *Application
var templates map[string]template.Template

func (a Application) StartServer() {
	app = &a
	addPageTemplates()
	addPageListeners()
}

func addPageTemplates() {
	index, err := ioutil.ReadFile("html/templates/index.html")
	if logger.LogErr(err) {
		return
	}

	templates = make(map[string]template.Template)
	indexTpl, _ := template.New("index").Parse(string(index))
	templates["index"] = *indexTpl
}

func addPageListeners() {
	port, _ := os.LookupEnv("SITE_PORT")
	http.HandleFunc("/weather", handleWeatherRequest)
	//http.HandleFunc("/purchase/", handleFondyRedirect)
	//http.HandleFunc(response, handleResponse)
	http.HandleFunc("/", startPage)
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("./html"))))

	err := http.ListenAndServeTLS(port, "ssl/server.crt", "ssl/server.key", nil)
	if logger.LogErr(err) {
		return
	}
}

func startPage(w http.ResponseWriter, r *http.Request) {
	t := templates["index"]

	err := t.Execute(w, nil)
	if logger.LogErr(err) {
		return
	}
}

func handleWeatherRequest(w http.ResponseWriter, r *http.Request) {
	city := "Murmansk"
	weather, _ := api.GetWeather(city)

	t := templates["index"]
	err := t.ExecuteTemplate(w, "update",
		struct {
			Temp      float64
			Feels     float64
			Pressure  int
			Humidity  int
			WindSpeed float64
		}{
			weather.Main.Temp,
			weather.Main.FeelsLike,
			weather.Main.Pressure,
			weather.Main.Humidity,
			weather.Wind.Speed,
		})
	if logger.LogErr(err) {
		return
	}
}
