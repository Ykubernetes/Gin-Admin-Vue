package system

import (
	"gorm.io/gorm"
)

// 分页条件
type PageOrder struct {
	Order string
	Where string
	Value []interface{}
}

// GetPage
func GetPage(table *gorm.DB, where interface{}, out interface{}, pageIndex, pageSize int, totalCount *int64, whereOrder ...PageOrder) error {
	res := table.Where(where)
	if len(whereOrder) > 0 {
		for _, whereOr := range whereOrder {
			if whereOr.Order != "" {
				res = res.Order(whereOr.Order)
			}
			if whereOr.Where != "" {
				res = res.Where(whereOr.Where, whereOr.Value...)
			}
		}
	}
	err := res.Count(totalCount).Error // SELECT count(*) FROM `menu`
	if err != nil {
		return err
	}
	if *totalCount == 0 {
		return err
	}
	return res.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(out).Error
}
