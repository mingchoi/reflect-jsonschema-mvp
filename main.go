package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/swaggest/jsonschema-go"
	_ "github.com/swaggest/jsonschema-go"
)

type Family struct {
	Children []string `json:"children" examples:"[\"A\",\"B\",\"C\"]" description:"This is a string array"`
}

type Payload struct {
	Name        string  `json:"name" example:"John" description:"A string"`
	Amount      float64 `json:"amount" example:"20.6" description:"A number"`
	Success     bool    `json:"success" example:"true" description:"A boolean"`
	BoolArray   []bool  `json:"boolArray" examples:"[true,false]" description:"A bool array"`
	StringArray []bool  `json:"stringArray" examples:"[\"A\",\"B\",\"C\"]" description:"A string array"`
	NumberArray []bool  `json:"numberArray" examples:"[1,2,3]" description:"A number array"`
	Family      Family  `json:"family" description:"A family"`
}

func main() {
	reflector := jsonschema.Reflector{}
	schema, err := reflector.Reflect(Payload{})
	if err != nil {
		log.Fatal(err)
	}

	builder := strings.Builder{}
	renderJson(&builder, schema, 1)
	fmt.Println(builder.String())
}

func renderJson(builder *strings.Builder, schema jsonschema.Schema, indent int) error {
	builder.WriteString("{\n")
	for k, v := range schema.Properties {
		if v.TypeObject != nil {
			for i := 0; i < indent; i++ {
				builder.WriteString("\t")
			}

			// object type handle
			if v.TypeObject.Type == nil {
				childInstance := reflect.New(v.TypeObject.ReflectType).Elem()
				//fmt.Printf("%+v\n", childInstance)
				reflector := jsonschema.Reflector{}
				childSchema, err := reflector.Reflect(childInstance)
				if err != nil {
					return err
				}

				builder.WriteString(fmt.Sprintf("\"%v\": ", k))
				err = renderJson(builder, childSchema, indent+1)
				if err != nil {
					return err
				}
				builder.WriteString(fmt.Sprintf("\t\t// %v\n", *v.TypeObject.Description))
				continue
			}
			valueType := v.TypeObject.Type.SimpleTypes
			if valueType == nil {
				valueType = &v.TypeObject.Type.SliceOfSimpleTypeValues[0]
			}
			description := v.TypeObject.Description

			switch *valueType {
			case jsonschema.Array:
				exampleValue := v.TypeObject.Examples
				b, _ := json.Marshal(exampleValue)
				builder.WriteString(fmt.Sprintf("\"%v\": %v\t\t// %v\n", k, string(b), *description))
			case jsonschema.Boolean:
				exampleValue := v.TypeObject.Examples[0]
				builder.WriteString(fmt.Sprintf("\"%v\": %v\t\t// %v\n", k, exampleValue, *description))
			case jsonschema.Integer:
				exampleValue := v.TypeObject.Examples[0]
				builder.WriteString(fmt.Sprintf("\"%v\": %v\n", k, exampleValue))
			case jsonschema.Null:
				exampleValue := "null"
				builder.WriteString(fmt.Sprintf("\"%v\": %v\t\t// %v\n", k, exampleValue, *description))
			case jsonschema.Number:
				exampleValue := v.TypeObject.Examples[0]
				//fmt.Printf("\"%v\": %v\n", k, exampleValue)
				builder.WriteString(fmt.Sprintf("\"%v\": %v\t\t// %v\n", k, exampleValue, *description))
			case jsonschema.Object:
				//builder.WriteString("\"%v\": ")
				//renderJson(builder, schema, indent+1)
				//builder.WriteString("\t\t// %v\n")
			case jsonschema.String:
				exampleValue := v.TypeObject.Examples[0]
				builder.WriteString(fmt.Sprintf("\"%v\": %v\t\t// %v\n", k, exampleValue, *description))
			}
		}
	}

	for i := 0; i < indent-1; i++ {
		builder.WriteString("\t")
	}
	builder.WriteString("}")
	return nil
}
