package server

import (
	"../logger"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

var app *Application
var templates map[string]string

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
	//buy, err := ioutil.ReadFile("html/templates/buy.html")
	//if logger.LogErr(err) { return }
	//purchaseSuccess, err := ioutil.ReadFile("html/templates/purchase_success.html")
	//if logger.LogErr(err) { return }

	templates = make(map[string]string)

	templates["index"] = string(index)
	//templates["buy"] = string(buy)
	//templates["purchaseSuccess"] = string(purchaseSuccess)
}

func addPageListeners() {
	port, _ := os.LookupEnv("SITE_PORT")
	//http.HandleFunc("/buy", handleBuyRequest)
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
	//idCookie, err := r.Cookie("uuid")

	//var userId string
	//if err != nil {
	//	idCookie = new(http.Cookie)
	//	idCookie.Name = "uuid"
	//	idCookie.Value = userId
	//	idCookie.Expires = time.Now().Add(30 * 24 * time.Hour)
	//	http.SetCookie(w, idCookie)
	//	if logger.LogErr(err) { return }
	//} else {
	//	userId = idCookie.Value
	//}

	tpl, err := template.New("index").Parse(templates["index"])
	if logger.LogErr(err) {
		return
	}

	err = tpl.Execute(w, nil)
	if logger.LogErr(err) {
		return
	}
}
