package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page   int64 `form:"page"`
	Limit  int64 `form:"limit"`
	Offset int64 `form:"offset"`
}

const (
	defaultLimit = 10
	defaultPage  = 1
)

func SetPagination(c *gin.Context) (*Pagination, error) {
	pagination := new(Pagination)

	err := c.BindQuery(pagination)
	if err != nil {
		return nil, err
	}

	err = pagination.setLimit()
	if err != nil {
		return nil, err
	}

	err = pagination.setPage()
	if err != nil {
		return nil, err
	}

	pagination.Offset = pagination.ConvertPageToOffset()

	return pagination, nil
}

func (p *Pagination) setLimit() error {
	if p.Limit > 100 {
		return errors.New("maximum.limit.is.100")
	}

	if p.Limit == 0 {
		p.Limit = defaultLimit
		return nil
	}

	return nil
}

func (p *Pagination) setPage() error {
	if p.Page == 0 {
		p.Page = defaultPage
		return nil
	}
	return nil
}

func (p *Pagination) ConvertPageToOffset() int64 {
	return (p.Page - 1) * p.Limit
}
