package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"tinyUrl/src/api_service/server/cache"
	"tinyUrl/src/api_service/server/database"
	"tinyUrl/src/api_service/server/handler"
	"tinyUrl/src/api_service/server/types"
)

type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

type Routes []Route

func AllRoutes() Routes {
	routes := Routes{
		Route{"fetchUrl", "GET", "/", fetchUrl},
		Route{"createUrl", "POST", "/books", createUrl},
	}
	return routes
}

//Reads from the routes slice to translate the values to httprouter.Handle
func NewRouter(routes Routes) *httprouter.Router {

	router := httprouter.New()
	for _, route := range routes {
		var handle httprouter.Handle

		handle = route.HandlerFunc

		router.Handle(route.Method, route.Path, handle)
	}

	return router
}

func createUrl(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the HomePage!")
	app.Handler.CreateTinyUrl()
	fmt.Println("Endpoint Hit: homePage")
}

func fetchUrl(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the HomePage!")
	app.Handler.GetOriginalUrl()
	fmt.Println("Endpoint Hit: homePage")
}

var app types.App

func initializeApp(app *types.App) {

	var err error
	app.DB, err = database.New()
	if err != nil {
		log.Fatal(err.Error())
	}
	app.Cache = cache.NewCache(1000)

	app.Handler, err = handler.NewHandler(app.DB, app.Cache)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	log.Println("Server will start at http://localhost:8000/")
	initializeApp(&app)
	router := NewRouter(AllRoutes())
	log.Fatal(http.ListenAndServe(":8080", router))
}
