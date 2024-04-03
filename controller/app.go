package controller

import (
	"github.com/a-h/templ"
	"go-meteo/view/components"
)

func Default() templ.Component {
	return components.Hello("Go-Meteo")
}
