package server

import (
	"../api"
	"../logger"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

var templates map[string]template.Template

func StartServer() {
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

func startPage(w http.ResponseWriter, _ *http.Request) {
	t := templates["index"]

	err := t.Execute(w, nil)
	if logger.LogErr(err) {
		return
	}
}

func handleWeatherRequest(w http.ResponseWriter, r *http.Request) {

	cities, _ := r.URL.Query()["city"]
	city := cities[0]
	weather, _ := api.GetWeather(city)

	t := templates["index"]
	err := t.ExecuteTemplate(w, "update",
		struct {
			Temp      float64
			Feels     float64
			Humidity  int
			WindSpeed float64
		}{
			weather.Main.Temp,
			weather.Main.FeelsLike,
			weather.Main.Humidity,
			weather.Wind.Speed,
		})
	if logger.LogErr(err) {
		return
	}
}
