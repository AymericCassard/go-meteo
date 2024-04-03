package controller

import (
	"github.com/a-h/templ"
	"go-meteo/view/components"
)

func Default(temp string) templ.Component {
	return components.Hello(temp)
}
