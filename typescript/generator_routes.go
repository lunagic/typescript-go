package typescript

import (
	"fmt"
	"reflect"
	"strings"
)

type Route struct {
	Path            string
	Method          string
	QueryParameters map[string]reflect.Type
	RequestBody     reflect.Type
	ResponseBody    reflect.Type
}

type tsRoute struct {
	Name            string
	Path            string
	Method          string
	Params          []tsRouteParam
	RequestBodyType string
	ResponseType    string
}

type tsRouteParam struct {
	Name string
	Type string
}

func (ts tsRoute) GenerateTypeScript() string {
	arguments := []string{}
	for _, param := range ts.Params {
		arguments = append(
			arguments,
			fmt.Sprintf("%s: %s", param.Name, param.Type),
		)
	}

	if ts.RequestBodyType != "" {
		arguments = append(
			arguments,
			fmt.Sprintf("payload: %s", ts.RequestBodyType),
		)
	}

	output := fmt.Sprintf("\texport const %s = async (%s) => {\n", ts.Name, strings.Join(arguments, ", "))

	if len(ts.Params) > 0 {
		output += "\t\tconst params = {\n"
		for _, param := range ts.Params {
			output += fmt.Sprintf("\t\t\t%s: %s,\n", param.Name, param.Name)
		}

		output += "\t\t}\n\n"

		output += "\t\tconst queryString = Object.keys(params).map((key) => {\n"
		output += "\t\t\treturn encodeURIComponent(key) + \"=\" + encodeURIComponent(params[key])\n"
		output += "\t\t}).join(\"&\")\n\n"

		output += fmt.Sprintf("\t\tconst response = await fetch(`%s?${queryString}`, {\n", ts.Path)
	} else {
		output += fmt.Sprintf("\t\tconst response = await fetch(\"%s\", {\n", ts.Path)
	}

	output += fmt.Sprintf("\t\t\tmethod: \"%s\",\n", ts.Method)

	if ts.RequestBodyType != "" {
		output += "\t\t\tbody: JSON.stringify(payload),\n"
	}

	output += "\t\t})\n"
	output += fmt.Sprintf("\n\t\treturn await response.json() as %s\n", ts.ResponseType)
	output += "\t}"

	return output
}
