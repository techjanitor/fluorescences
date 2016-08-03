package utils

import "html/template"

// TemplateFuncs holds our template functions
var TemplateFuncs = template.FuncMap{
	"add": add,
}

// add will add numbers with an offset
func add(offset, x, y int) int {
	return offset + x + y
}
