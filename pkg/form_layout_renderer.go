package pkg

import (
	"fmt"
	"strings"
)

/*
6:47PM DBG logger.go:48 > People.Min
6:47PM DBG logger.go:48 > The unrealistic maximum is 10.000
*/

func (f *FormLayout) RenderForm(data any) string {
	errors := make(map[string]string)
	return f.RenderFormWithErrors(data, errors)
}

func (f *FormLayout) RenderFormWithErrors(data any, errors map[string]string) string {
	m := FieldsToMap(FieldGenerator(data))
	var sb strings.Builder
	f.renderFormToBuilder(&sb, "", errors, m)
	return sb.String()
}

func (f *FormLayout) renderFormToBuilder(sb *strings.Builder, prefix string, errors map[string]string, m map[string]DataField,
) {
	for _, e := range f.elements {
		f.renderElement(e, sb, prefix, errors, m)
	}
}

func (f *FormLayout) RenderElement(
	name string,
	data any,
	errors map[string]string,
) string {
	e := f.findByName(name)
	m := FieldsToMap(FieldGenerator(data))
	var sb strings.Builder
	// probably prefix taken from name instead of ""
	f.renderElement(*e, &sb, "", errors, m)
	return sb.String()
}

func (f *FormLayout) renderElement(
	e FormElement,
	sb *strings.Builder,
	prefix string,
	errors map[string]string,
	m map[string]DataField,
) {
	switch e.Kind {
	case "hidden":
		fieldName := e.Name
		if len(prefix) > 0 {
			fieldName = prefix + "." + fieldName
		}
		field, ok := m[fieldName]
		if ok {
			sb.WriteString(fmt.Sprintf("<input type=\"hidden\" name=\"%s\" value=\"%v\" />", fieldName, field.Val()))
		}
	case "header":
		f.Theme.themeRenderHeader(sb, e)
	case "group":
		newPrefix := e.Name
		if len(prefix) > 0 {
			newPrefix = prefix + "." + newPrefix
		}
		f.Theme.themeRenderGroup(sb, m, newPrefix, e, errors)
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
				values := field.Value.([]string)
				for i := 0; i < len(field.Choices); i++ {
					choice := &field.Choices[i]
					if containsString(values, choice.Value) {
						choice.Checked = true
					}
				}
				f.Theme.themeRenderMulti(sb, field, e, prefix)
			} else {
				description := e.Description
				if len(e.Config.Description) > 0 {
					description = e.Config.Description
				}
				if field.Kind == "bool" {
					f.Theme.themeRenderCheckbox(sb, e, field, description, prefix, errors)
				} else if !field.Multi && len(field.Choices) > 0 {
					f.Theme.themeRenderSelect(sb, e, field, description, prefix, errors)
				} else {
					f.Theme.themeRenderInput(sb, e, field, prefix, errors)
				}
			}
		}
	case "yesno":
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
			description := e.Description
			if len(e.Config.Description) > 0 {
				description = e.Config.Description
			}
			f.Theme.themeRenderYesNo(sb, e, field, description, prefix, errors)
		}
	case "dropdown":
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
			//if field.Multi {
			//	values := field.Value.([]string)
			//	for i := 0; i < len(field.Choices); i++ {
			//		choice := &field.Choices[i]
			//		if containsString(values, choice.Value) {
			//			choice.Checked = true
			//		}
			//	}
			//	f.Theme.themeRenderMulti(sb, field, e, prefix)
			//} else {
			description := e.Description
			if len(e.Config.Description) > 0 {
				description = e.Config.Description
			}
			f.Theme.themeRenderSelect(sb, e, field, description, prefix, errors)
			//}
		}
	}
}

func renderCheckbox(sb *strings.Builder, f DataField, config ElementOpts, prefix string, class string) {
	checked := ""
	v, ok := f.Val().(bool)
	if ok {
		if v {
			checked = "checked"
		}
	}
	name := f.Name
	sb.WriteString(fmt.Sprintf("<input type=\"checkbox\" name=\"%s\" class=\"%s\" %s%s/>", name, class, checked, configToHtml(config)))
}

func renderSelect(sb *strings.Builder, f DataField, config ElementOpts, prefix string, class string, e FormElement) {
	name := f.Name
	if f.Kind == "int" {
		name = name + ":int"
	}
	// optgroup https://developer.mozilla.org/en-US/docs/Web/HTML/Element/optgroup
	if len(e.Config.Groups) > 0 {
		sb.WriteString(fmt.Sprintf("<select name=\"%s\" class=\"%s\"><option value=\"0\">-</option>", name, class))
		for group, name := range e.Config.Groups {
			sb.WriteString(fmt.Sprintf("<optgroup label=\"%s\">", name))
			for _, c := range f.Choices {
				if len(group) == 0 || c.Group == group {
					if c.IsSelected(f.Value) {
						sb.WriteString(fmt.Sprintf("<option value=\"%s\" selected=\"selected\">%s</option>", c.Val(), c.L()))
					} else {
						sb.WriteString(fmt.Sprintf("<option value=\"%s\">%s</option>", c.Val(), c.L()))
					}
				}
			}
			sb.WriteString("</optgroup>")
		}
		sb.WriteString("</select>")
	} else {
		sb.WriteString(fmt.Sprintf("<select name=\"%s\" class=\"%s\"><option value=\"0\">-</option>", name, class))
		for _, c := range f.Choices {
			if c.IsSelected(f.Value) {
				sb.WriteString(fmt.Sprintf("<option value=\"%s\" selected=\"selected\">%s</option>", c.Val(), c.L()))
			} else {
				sb.WriteString(fmt.Sprintf("<option value=\"%s\">%s</option>", c.Val(), c.L()))
			}
		}
		sb.WriteString("</select>")
	}
}

// TODO: DOES THIS CREATE A COPY?
func renderTextInput(sb *strings.Builder, f DataField, val any, config ElementOpts, prefix string, class string, errors map[string]string) {
	inputConstraints := ""

	inputType := "text"
	name := f.Name
	if f.Kind == "int" {
		name = name + ":int"
	}
	if f.SubKind == "email" {
		inputType = "email"
	} else if f.SubKind == "url" {
		inputType = "url"
	}
	sb.WriteString(fmt.Sprintf("<input name=\"%s\" type=\"%s\"%s value=\"%v\"%s class=\"%s\"/>",
		name, inputType,
		strings.TrimSpace(inputConstraints), val, configToHtml(config), class))

	if errorMsg, hasError := errors[name]; hasError {
		//	sb.WriteString(`
		//<div class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3">
		//  <svg class="h-5 w-5 text-red-500" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
		//    <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-8-5a.75.75 0 01.75.75v4.5a.75.75 0 01-1.5 0v-4.5A.75.75 0 0110 5zm0 10a1 1 0 100-2 1 1 0 000 2z" clip-rule="evenodd" />
		//  </svg>
		//</div>`)
		sb.WriteString(fmt.Sprintf("<p class=\"mt-2 text-sm text-red-600\">%s</p>", errorMsg))
	}
}

func configToHtml(config ElementOpts) string {
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
