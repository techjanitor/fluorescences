package utils

import "html/template"

func add(x, y int) int {
	return x + y
}

// TemplateFuncs holds our template functions
var TemplateFuncs = template.FuncMap{
	"add": add,
}
