package utils

import "html/template"

func add(offset, x, y int) int {
	return offset + x + y
}

// TemplateFuncs holds our template functions
var TemplateFuncs = template.FuncMap{
	"add": add,
}
