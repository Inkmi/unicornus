package builder

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unicornus/internal/ui"
)

type TestA struct {
	A string
}

type TestMulti struct {
	A []string `choices:"A1|A2|A3"`
}

type TestBool struct {
	A bool
}

type TestAB struct {
	A string
	B string
}

type TestB struct {
	B int `validate:"int|in:1,2,3" choices:"B1|B2|B3"`
}

func TestRenderForm(t *testing.T) {
	f := NewFormLayout().
		Add("A", "A")
	tdata := TestA{
		A: "b",
	}
	html := f.RenderForm(tdata)
	assert.Equal(t, ui.Clean(`
<label>A</label>
<input name="A" value="b"/>
`), html)
}

func TestRenderCheckbox(t *testing.T) {
	f := NewFormLayout().
		Add("A", "A")
	tdata := TestBool{
		A: true,
	}
	html := f.RenderForm(tdata)
	assert.Equal(t, ui.Clean(`
<label>A</label>
<input type="checkbox" name="A" />
`), html)
}

func TestRenderMulti(t *testing.T) {
	f := NewFormLayout().
		Add("A", "A")
	tdata := TestMulti{
		A: []string{"A", "B"},
	}
	html := f.RenderForm(tdata)
	assert.Equal(t, ui.Clean(`
<fieldset>
<input type="checkbox" name="A#A1">
<label>A1</label>
<input type="checkbox" name="A#A2">
<label>A2</label>
<input type="checkbox" name="A#A3">
<label>A3</label>
</fieldset>
`), html)
}

func TestRenderMultiGroup(t *testing.T) {
	f := NewFormLayout().
		Add("A", "A", ElementConfig{
			Choices: []ui.Choice{
				{
					Label: "A",
					Value: "A",
					Group: "G1",
				},
				{
					Label: "B",
					Value: "B",
					Group: "G1",
				},
				{
					Label: "C",
					Value: "C",
					Group: "G2",
				},
			},
			Groups: []string{"G1", "G2"},
		})
	tdata := TestMulti{
		A: []string{"A", "B"},
	}
	html := f.RenderForm(tdata)
	assert.Equal(t, ui.Clean(`
<fieldset>
<input type="checkbox" name="A#A">
<label>A</label>
<input type="checkbox" name="A#B">
<label>B</label>
<fieldset>
<input type="checkbox" name="A#C">
<label>C</label>
</fieldset>
`), html)
}

func TestTwoElementRenderForm(t *testing.T) {
	f := NewFormLayout().
		Add("A", "A").
		Add("B", "B")
	tdata := TestAB{
		A: "a",
		B: "b",
	}
	html := f.RenderForm(tdata)
	assert.Equal(t, ui.Clean(`
<label>A</label>
<input name="A" value="a"/>
<label>B</label>
<input name="B" value="b"/>
`), html)
}

func TestHeaderRenderForm(t *testing.T) {
	f := NewFormLayout().
		AddHeader("A")
	tdata := TestA{
		A: "a",
	}
	html := f.RenderForm(tdata)
	assert.Equal(t, ui.Clean(`
<h2>A</h2>
`), html)
}

func TestRenderSelectForm(t *testing.T) {
	f := NewFormLayout().
		Add("B", "B")
	tdata := TestB{
		B: 3,
	}
	html := f.RenderForm(tdata)
	assert.Equal(t, ui.Clean(`
<label>B</label>
<select name="B">
   <option value="0">-</option>
   <option value="1">B1</option>
   <option value="2">B2</option>
   <option value="3" selected="selected">B3</option>
</select>
`), html)
}

func TestRenderSelectWithChoicesForm(t *testing.T) {
	f := NewFormLayout().
		Add("A", "A", ElementConfig{
			Choices: []ui.Choice{
				{
					Label:    "A",
					Value:    "A",
					Selected: false,
				},
				{
					Label:    "B",
					Value:    "B",
					Selected: true,
				},
				{
					Label:    "C",
					Value:    "C",
					Selected: false,
				},
			}})
	tdata := TestA{
		A: "B",
	}
	html := f.RenderForm(tdata)
	assert.Equal(t, ui.Clean(`
<label>A</label>
<select name="A">
   <option value="0">-</option>
   <option value="A">A</option>
   <option value="B" selected="selected">B</option>
   <option value="C">C</option>
</select>
`), html)
}

func TestConfigToHtml(t *testing.T) {
	configHtml := configToHtml(ElementConfig{
		Id:          "id",
		Placeholder: "p",
		Label:       "l",
	})
	assert.Equal(t, " id=\"id\" placeholder=\"p\"", configHtml)
}

func TestRenderFormPlaceHolder(t *testing.T) {
	f := NewFormLayout().
		Add("A", "A", ElementConfig{
			Placeholder: "c",
		})
	tdata := TestA{
		A: "b",
	}
	html := f.RenderForm(tdata)
	assert.Equal(t, ui.Clean(`
<label>A</label>
<input name="A" value="b" placeholder="c"/>
`), html)
}

func TestRenderFormPlaceId(t *testing.T) {
	f := NewFormLayout().
		Add("A", "A", ElementConfig{
			Id: "c",
		})
	tdata := TestA{
		A: "b",
	}
	html := f.RenderForm(tdata)
	assert.Equal(t, ui.Clean(`
<label>A</label>
<input name="A" value="b" id="c"/>
`), html)
}