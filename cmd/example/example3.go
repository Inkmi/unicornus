package main

// S:1
import (
	// E:1
	"fmt"
	"net/http"

	// S:1
	"github.com/inkmi/unicornus/uni"
)

type subData3 struct {
	SubName string
}
type data3 struct {
	Name string
	Sub  subData3
}

// E:1

func example3(w http.ResponseWriter, req *http.Request) {
	// S:1
	// The data of the form
	d := data3{
		Name: "Unicornus",
		Sub: subData3{
			SubName: "Ha my name!",
		},
	}

	// Create a FormLayout
	// describing the form
	ui := uni.NewFormLayout().
		Add("Name", "Name Label").
		AddGroup("Sub", "Group", "Group Description", func(f *uni.FormLayout) {
			f.
				Add("SubName", "Sub Label")
		})

	// Render form layout with data
	// to html
	html := ui.RenderForm(d)
	// E:1
	fmt.Fprintf(w, html)
}
