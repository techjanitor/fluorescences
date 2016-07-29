package models

import (
	"html/template"
	"time"
)

// CommissionType holds the commission page
type CommissionType struct {
	Open        bool
	Content     string
	ContentOut  template.HTML
	HumanTime   string
	UpdatedTime time.Time
}
