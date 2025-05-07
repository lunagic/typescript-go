package typescript

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"time"
)

func New(
	configFuncs ...ServiceConfigFunc,
) *Service {
	s := &Service{
		outputNamespace: "GoGenerated",
		kindMap: map[reflect.Kind]string{
			reflect.String:     "string",
			reflect.Bool:       "boolean",
			reflect.Int:        "number",
			reflect.Int8:       "number",
			reflect.Int16:      "number",
			reflect.Int32:      "number",
			reflect.Int64:      "number",
			reflect.Uint:       "number",
			reflect.Uint8:      "number",
			reflect.Uint16:     "number",
			reflect.Uint32:     "number",
			reflect.Uint64:     "number",
			reflect.Uintptr:    "number",
			reflect.Float32:    "number",
			reflect.Float64:    "number",
			reflect.Complex64:  "number",
			reflect.Complex128: "number",
		},
		typeMap: map[reflect.Type]string{
			reflect.TypeFor[any]():           "any",
			reflect.TypeFor[time.Time]():     "string",
			reflect.TypeFor[time.Duration](): "number",
		},
	}

	for _, configFunc := range configFuncs {
		configFunc(s)
	}

	return s
}

type Service struct {
	outputNamespace string
	typeMap         map[reflect.Type]string
	kindMap         map[reflect.Kind]string
	outputTypes     map[string]reflect.Type
	outputData      map[string]any
	outputRoutes    map[string]Route
}

func (s *Service) Generate(writer io.Writer) error {
	typeScriptFileData := tsFile{
		namespace: s.outputNamespace,
		items:     []typescriptGenerator{},
	}

	// Types
	if len(s.outputTypes) > 0 {
		typeNames := []string{}
		for typeName, entry := range s.outputTypes {
			// Keep track of the keys so they can be sorted and used later
			typeNames = append(typeNames, typeName)

			// This maps the give structs to what they should be converted to when encountered later
			s.typeMap[entry] = typeName
		}

		sort.Strings(typeNames)

		for _, typeName := range typeNames {
			entryType := s.outputTypes[typeName]
			kind := entryType.Kind()

			if foundType, found := s.kindMap[kind]; found {
				typeScriptFileData.items = append(typeScriptFileData.items, tsType{
					Name: typeName,
					Type: foundType,
				})

				continue
			}

			if kind == reflect.Map {
				typeScriptFileData.items = append(typeScriptFileData.items, tsType{
					Name: typeName,
					Type: s.convertGoTypeToTypeScriptType(entryType),
				})

				continue
			}

			if kind == reflect.Slice && entryType.Elem().Kind() == reflect.Uint8 {
				typeScriptFileData.items = append(typeScriptFileData.items, tsType{
					Name: typeName,
					Type: "string",
				})
				continue
			}

			fields := []tsField{}

			loopOverStructFields(entryType, func(fieldDefinition reflect.StructField) {
				tag := parseJSONFieldTag(fieldDefinition.Tag.Get("json"))
				fieldName := fieldDefinition.Name

				if tag.NameOverride != "" {
					fieldName = tag.NameOverride
				}

				if tag.Ignored {
					return
				}

				tsType := s.convertGoTypeToTypeScriptType(fieldDefinition.Type)

				fields = append(fields, tsField{
					Name:     fieldName,
					Type:     tsType,
					Optional: tag.Omitempty,
				})
			})

			inter := tsInterface{
				Name:   typeName,
				Fields: fields,
			}

			typeScriptFileData.items = append(typeScriptFileData.items, inter)
		}
	}

	// Routes
	if len(s.outputRoutes) > 0 {
		// Add route endpoints
		routeNames := []string{}
		for routeName := range s.outputRoutes {
			routeNames = append(routeNames, routeName)
		}

		sort.Strings(routeNames)

		for _, routeName := range routeNames {
			route := s.outputRoutes[routeName]
			responseBodyType := s.convertGoTypeToTypeScriptType(route.ResponseBody)
			requestBodyType := ""

			if route.RequestBody != nil {
				requestBodyType = s.convertGoTypeToTypeScriptType(route.RequestBody)
			}

			params := []tsRouteParam{}

			paramKeys := []string{}
			for key := range route.QueryParameters {
				paramKeys = append(paramKeys, key)
			}

			sort.Strings(paramKeys)

			for _, key := range paramKeys {
				params = append(params, tsRouteParam{
					Name: key,
					Type: s.convertGoTypeToTypeScriptType(route.QueryParameters[key]),
				})
			}

			typeScriptFileData.items = append(typeScriptFileData.items, tsRoute{
				Name:            routeName,
				Path:            route.Path,
				Method:          route.Method,
				Params:          params,
				RequestBodyType: requestBodyType,
				ResponseType:    responseBodyType,
			})
		}
	}

	// Data
	if len(s.outputData) > 0 {
		dataVarNames := []string{}
		for dataVarName := range s.outputData {
			dataVarNames = append(dataVarNames, dataVarName)
		}
		sort.Strings(dataVarNames)

		for _, dataVarName := range dataVarNames {
			data := s.outputData[dataVarName]
			typeScriptFileData.items = append(typeScriptFileData.items, tsData{
				Name: dataVarName,
				Type: s.convertGoTypeToTypeScriptType(reflect.ValueOf(data).Type()),
				Data: data,
			})
		}
	}

	typeScriptToOutput := typeScriptFileData.GenerateTypeScript()

	_, _ = writer.Write([]byte(typeScriptToOutput))

	return nil
}

func (s *Service) convertGoTypeToTypeScriptType(item reflect.Type) string {
	isPointer := item.Kind() == reflect.Pointer
	if isPointer {
		item = item.Elem()
	}

	isSlice := item.Kind() == reflect.Slice
	isMap := item.Kind() == reflect.Map
	stringerTyp := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	isStringer := item.Implements(stringerTyp)

	// Read through the pointer/slice/map
	if isSlice {
		item = item.Elem()
	}

	if isMap {
		return fmt.Sprintf(
			"{ [key: %s]: %s } | null",
			s.convertGoTypeToTypeScriptType(item.Key()),
			s.convertGoTypeToTypeScriptType(item.Elem()),
		)
	}

	typeFromMapping, found := s.getFromMap(item)

	if !found {
		if isStringer {
			return "string"
		}

		return "unknown"
	}

	if isSlice {
		typeFromMapping += "[]"
	}

	if isSlice || isPointer {
		typeFromMapping += " | null"
	}

	return typeFromMapping
}

func (s *Service) getFromMap(item reflect.Type) (string, bool) {
	if typeString, found := s.typeMap[item]; found {
		return typeString, true
	}

	if typeString, found := s.kindMap[item.Kind()]; found {
		return typeString, true
	}

	return "", false
}

func loopOverStructFields(fieldType reflect.Type, fieldHandler func(fieldDefinition reflect.StructField)) {
	if fieldType.Kind() == reflect.Pointer {
		fieldType = fieldType.Elem()
	}

	for i := range fieldType.NumField() {
		subFieldType := fieldType.Field(i).Type
		fieldDefinition := fieldType.Field(i)

		if !fieldDefinition.IsExported() {
			continue
		}

		if fieldDefinition.Type.Kind() == reflect.Struct && fieldDefinition.Anonymous {
			loopOverStructFields(subFieldType, fieldHandler)

			continue
		}

		fieldHandler(fieldDefinition)
	}
}
