package domain

import "math"

type Pagination struct {
	Limit     int64 `json:"limit"`
	Page      int64 `json:"page"`
	TotalPage int64 `json:"totalPage"`
	TotalData int64 `json:"totalData"`
}

func (p *Pagination) GetTotalPage() int64 {
	if p.TotalPage > 0 {
		return p.TotalPage
	}

	return int64(math.Ceil(float64(p.TotalData) / float64(p.Limit)))
}
