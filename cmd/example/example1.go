package main

import (
	"fmt"
	uni "github.com/inkmi/unicornus/pkg"
	"net/http"
)

type data1 struct {
	Name string
}

func example1(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<!DOCTYPE html><html><body>")

	// The data of the form
	d := data1{
		Name: "Unicornus",
	}
	// Create a FormLayout
	// describing the form
	ui := uni.NewFormLayout().
		Add("Name", "Name Label")

	// Render form layout with data
	// to html
	html := ui.RenderForm(d)
	fmt.Fprintf(w, html)

	fmt.Fprintf(w, "</body></html>")

}
