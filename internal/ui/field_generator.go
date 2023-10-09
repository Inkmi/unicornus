package ui

import (
	"reflect"
)
import . "github.com/moznion/go-optional"

func FieldGenerator(obj interface{}) []DataField {
	vals := make([]DataField, 0, 20)
	original := reflect.ValueOf(obj)
	return translateRecursive(vals, "", original)
	// return translateRecursive(vals, original.Type().Name(), original)
}

func translateRecursive(vals []DataField, prefix string, original reflect.Value) []DataField {
	switch original.Kind() {
	case reflect.Struct:
		vals = translateStruct(prefix, vals, original)
		return vals
	case reflect.Ptr:
		originalValue := original.Elem()
		if !originalValue.IsValid() {
			return vals
		}
		return translateRecursive(vals, prefix, originalValue)
	case reflect.Interface:
		originalValue := original.Elem()
		translateRecursive(vals, prefix, originalValue)
		return vals
	default:
		return vals
	}
}

func translateStruct(prefix string, vals []DataField, original reflect.Value) []DataField {
	for i := 0; i < original.NumField(); i += 1 {
		df := DataField{}
		if len(prefix) == 0 {
			df.Name = original.Type().Field(i).Name
		} else {
			df.Name = prefix + "." + original.Type().Field(i).Name
		}
		optionalValue := hasOptional(original.Type().Field(i).Type.Name())
		// Handle setFromStructPtr of type
		// those are handled as optional
		// *bool
		// *int
		if !optionalValue && original.Type().Field(i).Type.Kind() == reflect.Ptr {
			df.Kind = original.Type().Field(i).Type.Elem().Name()
			df.Optional = true
			valid := original.Field(i).Elem().IsValid()
			if original.Type().Field(i).Type.Elem().ConvertibleTo(reflect.TypeOf(0)) {
				df.Kind = "int"
				if valid {
					df.Value = original.Field(i).Elem().Int()
				}
			} else if original.Type().Field(i).Type.Elem().ConvertibleTo(reflect.TypeOf(true)) {
				df.Kind = "bool"
				if valid {
					df.Value = original.Field(i).Elem().Bool()
				}
			}
		} else {
			typ := ""
			// it seems that the Kind of the Option is Slice
			// so check with the type we got with Name
			if !optionalValue && original.Type().Field(i).Type.Kind() == reflect.Slice {
				df.Multi = true
				typ = original.Type().Field(i).Type.Elem().Name()
			} else {
				typ = original.Type().Field(i).Type.Name()
			}
			df.Optional = hasOptional(typ)
			df.Kind = removeOptional(typ)
			if df.Kind == "string" {
				setString(original, i, &df)
			} else if df.Kind == "bool" {
				setBool(original, i, &df)
			} else if df.Kind == "int" {
				setInt(&df, original, i)
			}
		}

		// parse Tag for validation, choices and error messages
		tagS := string(original.Type().Field(i).Tag)
		if tagS != "" {
			t := ParseTag(tagS)
			if t.Validation != nil {
				df.Validation = *t.Validation
			}
			if t.ErrorMessage != nil {
				df.ErrorMessage = *t.ErrorMessage
			}
			if t.Choices != nil && len(t.Choices) > 0 {
				df.Choices = t.Choices
			}
		}
		vals = append(vals, df)
		newPrefix := prefix + " ." + original.Type().Field(i).Name
		if len(prefix) == 0 {
			newPrefix = original.Type().Field(i).Name
		}
		vals = translateRecursive(vals, newPrefix, original.Field(i))
	}
	return vals
}

func setString(original reflect.Value, i int, f *DataField) {
	if f.Multi {
		// https://stackoverflow.com/questions/32890137/how-to-get-slice-underlying-value-via-reflect-value
		f.Value = original.Field(i).Interface().([]string)
	} else {
		f.Value = original.Field(i).String()
	}
}

func setBool(original reflect.Value, i int, f *DataField) {
	if !f.Optional {
		f.Value = original.Field(i).Bool()
	} else {
		o := original.Field(i).Interface().(Option[bool])
		if o.IsSome() {
			val, _ := o.Take()
			f.Value = val
		} else {
			f.Value = false
		}
	}
}

func setInt(f *DataField, original reflect.Value, i int) {
	if !f.Optional {
		f.Value = original.Field(i).Int()
	} else {
		o := original.Field(i).Interface().(Option[int])
		if o.IsSome() {
			val, _ := o.Take()
			f.Value = int64(val)
		} else {
			f.Value = nil
		}
	}
}