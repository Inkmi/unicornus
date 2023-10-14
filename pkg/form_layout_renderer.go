package pkg

import (
	"fmt"
	"strings"
)

func (f *FormLayout) RenderForm(data any) string {
	var sb strings.Builder
	f.renderFormToBuilder(&sb, data, "")
	return sb.String()
}

func (f *FormLayout) renderFormToBuilder(sb *strings.Builder, data any, prefix string) {
	m := FieldsToMap(FieldGenerator(data))
	for _, e := range f.elements {
		switch e.Kind {
		case "header":
			sb.WriteString(fmt.Sprintf("<h2>%s</h2>", e.Name))
		case "group":
			sb.WriteString("<div>")
			sb.WriteString(e.Name)
			newPrefix := e.Name
			if len(prefix) > 0 {
				newPrefix = prefix + "." + newPrefix
			}
			e.Config.SubLayout.renderFormToBuilder(sb, data, newPrefix)
			sb.WriteString("</div>")
		case "input":
			// take value string from MAP of name -> DataField
			// take type if no type is given from DataField
			fieldName := e.Name
			if len(prefix) > 0 {
				fieldName = prefix + "." + fieldName
			}
			field, ok := m[fieldName]
			if ok {
				if len(e.Config.Choices) > 0 {
					field.Choices = e.Config.Choices
				}
				if field.Multi {
					renderMulti(sb, field, e.Config, prefix)
				} else {
					sb.WriteString("<div>")
					if len(e.Config.Label) > 0 {
						sb.WriteString(fmt.Sprintf("<label>%s</label>", e.Config.Label))
					}
					if field.Kind == "bool" {
						renderCheckbox(sb, field, e.Config, prefix)
					} else if !field.Multi && len(field.Choices) > 0 {
						renderSelect(sb, field, e.Config, prefix)
					} else {
						renderTextInput(sb, field, field.Val(), e.Config, prefix)
					}
					sb.WriteString("</div>")
				}
			}
		}
	}
}

func renderCheckbox(sb *strings.Builder, f DataField, config ElementConfig, prefix string) {
	checked := ""
	v, ok := f.Val().(bool)
	if ok {
		if v {
			checked = "checked"
		}
		name := f.Name
		sb.WriteString(fmt.Sprintf("<input type=\"checkbox\" name=\"%s\" %s%s/>", name, checked, configToHtml(config)))
	}
}

func renderMulti(sb *strings.Builder, f DataField, config ElementConfig, prefix string) {
	// Should this move to Field generation?
	values := f.Value.([]string)
	for i := 0; i < len(f.Choices); i++ {
		choice := &f.Choices[i]
		if containsString(values, choice.Value) {
			choice.Checked = true
		}
	}
	if len(config.Groups) > 0 {
		for _, group := range config.Groups {
			sb.WriteString("<div>")
			sb.WriteString("<fieldset>")
			// range copies slice
			for _, c := range f.Choices {
				if c.Group == group {
					name := f.Name + "#" + c.Val()
					sb.WriteString("<div>")
					if c.Checked {
						sb.WriteString(fmt.Sprintf("<input type=\"checkbox\" name=\"%s\" checked>", name))
					} else {
						sb.WriteString(fmt.Sprintf("<input type=\"checkbox\" name=\"%s\">", name))
					}
					sb.WriteString(fmt.Sprintf(`<label>%s</label>`, c.L()))
					sb.WriteString("</div>")
				}
			}
			sb.WriteString("</fieldset>")
			sb.WriteString("</div>")
		}

	} else {
		sb.WriteString("<div>")
		sb.WriteString("<fieldset>")
		for _, c := range f.Choices {
			name := f.Name + "#" + c.Val()
			sb.WriteString("<div>")
			if c.Checked {
				sb.WriteString(fmt.Sprintf("<input type=\"checkbox\" name=\"%s\" checked>", name))
			} else {
				sb.WriteString(fmt.Sprintf("<input type=\"checkbox\" name=\"%s\">", name))
			}
			sb.WriteString(fmt.Sprintf(`<label>%s</label>`, c.L()))
			sb.WriteString("</div>")
		}
		sb.WriteString("</fieldset>")
		sb.WriteString("</div>")
	}
}

func renderSelect(sb *strings.Builder, f DataField, config ElementConfig, prefix string) {
	sb.WriteString(fmt.Sprintf("<select name=\"%s\"><option value=\"0\">-</option>", f.Name))
	for _, c := range f.Choices {
		if c.IsSelected(f.Value) {
			sb.WriteString(fmt.Sprintf("<option value=\"%s\" selected=\"selected\">%s</option>", c.Val(), c.L()))
		} else {
			sb.WriteString(fmt.Sprintf("<option value=\"%s\">%s</option>", c.Val(), c.L()))
		}
	}
	sb.WriteString("</select>")
}

func renderTextInput(sb *strings.Builder, f DataField, val any, config ElementConfig, prefix string) {
	sb.WriteString(fmt.Sprintf("<input name=\"%s\" value=\"%s\"%s/>", f.Name, val, configToHtml(config)))
}

func configToHtml(config ElementConfig) string {
	id := ""
	if len(config.Id) > 0 {
		id = fmt.Sprintf(" id=\"%s\"", config.Id)
	}
	placeholder := ""
	if len(config.Placeholder) > 0 {
		placeholder = fmt.Sprintf(" placeholder=\"%s\"", config.Placeholder)
	}
	configStr := fmt.Sprintf("%s%s", id, placeholder)
	return configStr
}

func containsString(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

/*

func SetChoices(setKey string, fields []FieldV, allValues []string) {
	for i := range fields {
		if fields[i].Name == setKey {
			var choices []Choice
			values := fields[i].Value.([]string)
			for _, p := range allValues {
				choices = append(choices, Choice{
					Label:    p,
					Value:    p,
					Selected: lo.Contains(values, p),
				})
			}

			fields[i].Choices = choices
			fields[i].Kind = "string"
		}
	}
}

func SetKey(
	setKey string,
	fields []FieldV,
	allValues []string,
	group func(k string) string,
	label func(l string) string,
) {
	for i := range fields {
		if fields[i].Name == setKey {
			var choices []Choice
			values := fields[i].Value.([]string)
			for _, p := range allValues {
				choices = append(choices, Choice{
					Group:    group(p),
					Label:    label(p),
					Value:    p,
					Selected: lo.Contains(values, p),
				})
			}

			fields[i].Choices = choices
			fields[i].Kind = "string"
		}
	}
}

*/