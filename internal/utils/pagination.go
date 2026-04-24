package utils

import (
	"math"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

type PaginationParams struct {
	Page     int
	PageSize int
	Sort     string
}

func GetPaginationParams(c *gin.Context) PaginationParams {
	page := 1
	pageSize := DefaultPageSize

	if p := c.Query("page"); p != "" {
		if val, err := parseInt(p); err == nil && val > 0 {
			page = val
		}
	}

	if ps := c.Query("pageSize"); ps != "" {
		if val, err := parseInt(ps); err == nil && val > 0 {
			if val > MaxPageSize {
				val = MaxPageSize
			}
			pageSize = val
		}
	}

	sort := c.Query("sort")

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
	}
}

func (p PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func Paginate(p PaginationParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.GetOffset()).Limit(p.PageSize)
	}
}

func CalculatePaginationMeta(total int64, page, pageSize int) PaginationMeta {
	totalPages := int64(math.Ceil(float64(total) / float64(pageSize)))
	return PaginationMeta{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

func parseInt(s string) (int, error) {
	var v int
	_, err := parseUint([]byte(s), &v)
	return v, err
}

func parseUint(b []byte, v *int) (int, error) {
	var u uint64
	for _, c := range b {
		if c < '0' || c > '9' {
			return 0, ErrSyntax
		}
		u = u*10 + uint64(c-'0')
	}
	*v = int(u)
	return len(b), nil
}

var ErrSyntax error = nil // Placeholder, should be defined properly
