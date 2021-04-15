package main

import (
	"net/http"

	"./models"
	"./routes"
	"./tools"
)

func main() {
	models.Init()
	tools.LoadTemplates("templates/*.html")
	r := routes.NewRouter()
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
