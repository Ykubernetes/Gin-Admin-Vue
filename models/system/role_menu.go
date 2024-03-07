package system

import (
	"gitee.com/go-server/global"
	"gitee.com/go-server/models/basemodel"
	"gorm.io/gorm"
	"time"
)

// 角色-菜单
type RoleMenu struct {
	basemodel.Model
	RoleID uint64 `gorm:"column:role_id;unique_index:uk_role_menu_role_id;not null;"` // 角色ID
	MenuID uint64 `gorm:"column:menu_id;unique_index:uk_role_menu_role_id;not null;"` // 菜单ID
}

// 添加前
func (m *RoleMenu) BeforeCreate(tx *gorm.DB) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

// 更新前
func (m *RoleMenu) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}

// 设置角色菜单权限
func (RoleMenu) SetRole(roleid uint64, menuids []uint64) error {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where(&RoleMenu{RoleID: roleid}).Delete(&RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if len(menuids) > 0 {
		for _, mid := range menuids {
			rm := new(RoleMenu)
			rm.RoleID = roleid
			rm.MenuID = mid
			if err := tx.Create(rm).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}
