package utils

import (
	"html/template"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// Markdown will sanitize the content and return markdown html
func Markdown(content string) template.HTML {
	// make the post formatted with markdown
	unsafe := blackfriday.MarkdownCommon([]byte(content))
	// sanitize the input
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	// convert to template format
	return template.HTML(html)
}
