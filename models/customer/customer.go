package customer

import (
	"gitee.com/go-server/global"
	"gitee.com/go-server/models/basemodel"
	"gorm.io/gorm"
	"time"
)

// 客户管理
type Customer struct {
	basemodel.Model
	UserName       string `gorm:"column:user_name;size:32;unique_index:uk_customer_user_name;not null;" json:"user_name" form:"user_name"`     // 客户姓名
	CompanyName    string `gorm:"column:company_name;unique_index:uk_customer_company_name;not null;" json:"company_name" form:"company_name"` //公司名称
	Position       string `gorm:"column:position;size:32;" json:"position" form:"position"`                                                    // 职位名称                                                      // 职位名称
	Birthday       string `gorm:"column:birthday;" json:"birthday" form:"birthday"`
	Gender         string `gorm:"column:gender;not null;"json:"gender" form:"gender"`
	Province       string `gorm:"column:province;size:32;" json:"province" form:"province"` // 省份
	City           string `gorm:"column:city;size:32;" json:"city" form:"province"`         // 城市
	County         string `gorm:"column:county;size:32;" json:"county" form:"county"`       // 乡镇
	Address        string `gorm:"column:address;not null;" json:"address" form:"address"`   // 门牌号
	Email          string `gorm:"column:email;size:64;" json:"email" form:"email"`
	Phone          string `gorm:"column:phone;type:char(20);" json:"phone" form:"phone"`
	Founder        string `gorm:"column:founder;not null;"json:"founder" form:"founder"`                                      // 创建人
	CustomerSoruce string `gorm:"column:customer_soruce;not null;comment:'来源';"json:"customer_soruce" form:"customer_soruce"` // 来源
	Status         uint8  `gorm:"column:status;type:tinyint(1);not null;" json:"status" form:"status"`
	Remark         string `gorm:"comment:'备注';" json:"remark" form:"remark"`
}

// 添加前
func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Customer) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}

func FindCustomerByUserName(username string) (c *Customer, err error) {
	var c1 Customer
	err = global.DB.First(&c1, "user_name=?", username).Error
	if err == gorm.ErrRecordNotFound { // record not found
		global.Log.Warnln("数据库中没有查询到结果record not found")
		return &Customer{}, err
	} else if err != nil {
		global.Log.Warnln("err里面有其他错误")
		return &Customer{}, err
	}
	return &c1, nil
}

func FindCustomerByUserID(Id uint64) (c *Customer, err error) {
	var c1 Customer
	err = global.DB.Debug().Where("id = ?", Id).First(&c1).Error
	if err == gorm.ErrRecordNotFound { // record not found
		global.Log.Warnln("数据库中没有查询到结果record not found")
		return &Customer{}, err
	} else if err != nil {
		global.Log.Warnln("err里面有其他错误")
		return &Customer{}, err
	}
	return &c1, nil
}

func DelCustomerByUserName(username string) bool {
	var c1 Customer
	err := global.DB.Where("user_name = ?", username).First(&c1).Error
	if err == gorm.ErrRecordNotFound { // record not found
		global.Log.Warnln("数据库中没有查询到结果record not found")
		return false
	} else if err != nil {
		global.Log.Warnln("err里面有其他错误")
		return false
	}
	global.DB.Delete(&c1)
	return true
}
