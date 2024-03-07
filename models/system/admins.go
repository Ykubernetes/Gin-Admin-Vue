package system

import (
	"gitee.com/go-server/global"
	"gitee.com/go-server/models/basemodel"
	"gitee.com/go-server/service/response"
	"gorm.io/gorm"
	"time"
)

// 后台用户
type Admins struct {
	basemodel.Model
	Memo         string `gorm:"column:memo;size:64;" json:"memo" form:"memo"`
	UserName     string `gorm:"column:user_name;size:32;unique_index:uk_admins_user_name;not null;" json:"user_name" form:"user_name"` // 用户名
	RealName     string `gorm:"column:real_name;size:32;" json:"real_name" form:"real_name"`                                           // 真实姓名
	Password     string `gorm:"column:password;type:char(150);not null;" json:"password" form:"password"`                              // 密码(sha1(md5(明文))加密)
	Email        string `gorm:"column:email;size:64;" json:"email" form:"email"`                                                       // 邮箱
	Phone        string `gorm:"column:phone;type:char(20);" json:"phone" form:"phone"`                                                 // 手机号
	Status       uint8  `gorm:"column:status;type:tinyint(1);not null;" json:"status" form:"status"`                                   // 状态(1:正常 2:未激活 3:暂停使用)
	Avatar       string `gorm:"not null;"json:"avatar"`                                                                                // 头像
	Introduction string `gorm:"not null;"json:"introduction"`                                                                          // 介绍
	DeptId       int    `gorm:"type:int(11)" json:"deptId"`                                                                            //部门编码
	PostId       int    `gorm:"type:int(11)" json:"postId"`                                                                            //职位编码 	// 角色编码
}

func (a Admins) TableName() string {
	return "admins"
}

// 添加前
func (m *Admins) BeforeCreate(tx *gorm.DB) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

func (m *Admins) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}

// 删除用户及关联数据
func (Admins) Delete(adminsids []uint64) error {
	tx := global.DB.Begin() // 手动开启事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	// 删除Admins
	if err := tx.Where("id in (?)", adminsids).Delete(&Admins{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 再删除AdminsRole
	if err := tx.Where("admins_id in (?)", adminsids).Delete(&AdminsRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func FindByUsername(username string) (Admins, error) {
	var user Admins
	err := global.DB.Where("user_name = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		global.Log.Fatal("数据库中没有查询到此用户...")
		return user, err
	} else if err != nil {
		global.Log.Fatalf("数据库查询用户发生错误,%s", err)
		return user, err
	}
	return user, nil
}

func FindByUserId(userid uint64) (Admins, error) {
	if userid == response.SUPER_ADMIN_ID {
		userid = 1
	}
	var user Admins
	err := global.DB.Where("id = ?", userid).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		global.Log.Fatal("数据库中没有查询到此用户...")
		return user, err
	} else if err != nil {
		global.Log.Fatalf("数据库查询用户发生错误,%s", err)
		return user, err
	}
	return user, nil
}

// 返回用户和用户角色 连接查询的结果
func (a *Admins) GetUserListAndRoleId() (u []Admins, err error) {
	table := global.DB.Table(a.TableName()).Select([]string{"admins.*", "role.name"})
	table = table.Joins("left join admins_role on admins.id = admins_role.admins_id").Joins("left join role ON  admins_role.role_id = role.id")
	if a.ID != 0 {
		table = table.Where("id = ?", a.ID)
	}
	if a.UserName != "" {
		table = table.Where("user_name = ?", a.UserName)
	}
	if a.Password != "" {
		table = table.Where("password = ?", a.Password)
	}
	if a.DeptId != 0 {
		table = table.Where("dept_id = ?", a.DeptId)
	}
	if a.PostId != 0 {
		table = table.Where("post_id = ?", a.PostId)
	}
	if err = table.Find(&u).Error; err != nil {
		return
	}
	return
}
