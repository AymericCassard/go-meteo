package main

import (
	"fmt"
	"github.com/a-h/templ"
	"go-meteo/controller"
	"net/http"
)

func main() {

	index := index(controller.Default())
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images/"))))
	http.Handle("/", templ.Handler(index))
	http.HandleFunc("/ville", controller.ReturnVilles)
	fmt.Println("Listening on http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
