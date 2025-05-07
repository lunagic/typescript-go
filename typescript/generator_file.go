package typescript

type tsFile struct {
	namespace string
	items     []typescriptGenerator
}

func (ts tsFile) GenerateTypeScript() string {
	output := ""
	output += "// NOTE: This file was auto-generated\n// and should NOT be edited manually.\n\n"
	output += "export namespace " + ts.namespace + " {\n"

	// Write all the items to the Writer
	for i, tsItem := range ts.items {
		output += tsItem.GenerateTypeScript()
		output += "\n"

		if i != len(ts.items)-1 {
			output += "\n"
		} else {
			output += "}\n"
		}
	}

	return output
}
