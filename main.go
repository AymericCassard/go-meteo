package main

import (
	"fmt"
	"github.com/a-h/templ"
	"go-meteo/view/components"
	"net/http"
)

func main() {
	component := components.Hello("Go-Meteo")

	http.Handle("/", templ.Handler(component))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
