package utils

import "math"

// Paged contains the fields for pagination
type Paged struct {
	Total       int    `json:"total"`
	PerPage     int    `json:"per_page"`
	Pages       int    `json:"pages"`
	CurrentPage int    `json:"current_page"`
	Skip        int    `json:"skip"`
	Start       int    `json:"range_start"`
	End         int    `json:"range_end"`
	Prev        int    `json:"page_prev"`
	Next        int    `json:"page_next"`
	Min         bool   `json:"page_min"`
	Max         bool   `json:"page_max"`
	Path        string `json:"page_path"`
	Key         string `json:"key"`
}

// Asc calculates the pagination numbers as ascending
func (p *Paged) Asc() {
	p.Pages = int(math.Ceil(float64(p.Total) / float64(p.PerPage)))

	p.End = ((p.CurrentPage - 1) * p.PerPage)
	p.Start = (p.End + p.PerPage) - 1

	if p.Start > p.Total {
		p.Start = p.Total
	}

	p.Next = p.CurrentPage + 1
	if p.Next > p.Pages {
		p.Next = p.Pages
		p.Max = true
	}

	p.Prev = p.CurrentPage - 1
	if p.Prev < 1 {
		p.Prev = 1
		p.Min = true
	}

	p.Skip = ((p.CurrentPage - 1) * p.PerPage)

	return

}

// Desc calculates the pagination numbers as descending
func (p *Paged) Desc() {
	p.Pages = int(math.Ceil(float64(p.Total) / float64(p.PerPage)))

	p.Start = (p.Total - (p.CurrentPage * p.PerPage) + p.PerPage)
	p.End = (p.Total) - (p.CurrentPage * p.PerPage)

	if p.End < 0 {
		p.End = 0
	}

	p.Next = p.CurrentPage + 1
	if p.Next > p.Pages {
		p.Next = p.Pages
		p.Max = true
	}

	p.Prev = p.CurrentPage - 1
	if p.Prev < 1 {
		p.Prev = 1
		p.Min = true
	}

	p.Skip = ((p.CurrentPage - 1) * p.PerPage)

	return

}
