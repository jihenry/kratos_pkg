package page

import (
	"errors"
	"sync"
)

type Page struct {
	PageIndex int `json:"pageIndex" form:"pageIndex"` //页码
	PageSize  int `json:"pageSize" form:"pageSize"`   //分页大小
	Total     int `json:"total"`                      // 符合条件的总记录数（仅用于响应结构中）
}

func New(pageIndex, pageSize int) *Page {
	pool := sync.Pool{
		New: func() interface{} {
			return &Page{
				PageIndex: pageIndex,
				PageSize:  pageSize,
				Total:     0,
			}
		},
	}
	obj := pool.Get().(*Page)
	pool.Put(obj)
	return obj
}

// Check 检查 Page 的参数
func (p *Page) Check() error {
	return Check(p.PageIndex, p.PageSize)
}

// Offset 获取分页偏移量
func (p *Page) Offset() int {
	return GetOffset(p.PageIndex, p.PageSize)
}

// SetValidVal 设置Page参数的有效值
func (p *Page) SetValidVal(defaultSize int) {
	p.PageIndex, p.PageSize = GetValidPage(p.PageIndex, p.PageSize, defaultSize)
}

// Check  检查页码参数是否有错误
func Check(pageIndex, pageSize int) error {
	if pageIndex < 0 {
		return errors.New("invalid pageIndex")
	}
	if pageSize < 0 {
		return errors.New("invalid pageSize")
	}
	return nil
}

// GetOffset 根据页码和分页大小获取查询偏移量 offset
func GetOffset(pageIndex, pageSize int) (offset int) {
	if pageIndex < 1 {
		return 0
	}
	offset = (pageIndex - 1) * pageSize
	return offset
}

// GetValidPage 获取有效的页码和分页大小
func GetValidPage(pageIndex, pageSize, defaultSize int) (page int, size int) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 1 {
		if defaultSize > 0 {
			pageSize = defaultSize
		} else {
			pageSize = 1
		}
	}
	return pageIndex, pageSize
}

// SlicePage 按分页大小分页，获取长度为 total 的切片 pageIndex 页的开始位置的结束位置；exist 表示分页位置是否存在
func SlicePage(total, pageIndex, pageSize int) (start, end int, exist bool) {
	if total < 1 {
		return -1, -1, false
	}

	start = GetOffset(pageIndex, pageSize)
	if total <= start {
		return -1, -1, false
	}
	end = start + pageSize
	if total <= end {
		end = total
	}
	return start, end, true
}
