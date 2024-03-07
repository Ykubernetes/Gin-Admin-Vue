package system

import (
	"gitee.com/go-server/global"
	"gitee.com/go-server/models/basemodel"
	"time"
)

// Post 岗位
type Post struct {
	basemodel.Model
	PostName string `gorm:"type:varchar(128);" json:"postName"` //岗位名称
	PostCode string `gorm:"type:varchar(128);" json:"postCode"` //岗位代码
	Sort     int    `gorm:"type:int(4);" json:"sort"`           //岗位排序
	Status   string `gorm:"type:int(1);" json:"status"`         //状态
	Remark   string `gorm:"type:varchar(255);" json:"remark"`   //描述
	CreateBy string `gorm:"type:varchar(128);" json:"createBy"`
	UpdateBy string `gorm:"type:varchar(128);" json:"updateBy"`
	Params   string `gorm:"-" json:"params"`
	// 新增一条删除时间字段
	DeletedAt *time.Time `gorm:"column:delete_time" sql:"index" json:"-"`
}

func (p *Post) TableName() string {
	return "post"
}

// Create 创建岗位
func (p *Post) Create() (Post, error) {
	var pdata Post
	result := global.DB.Table(p.TableName()).Create(&p)
	if result.Error != nil {
		err := result.Error
		return pdata, err
	}
	pdata = *p
	return pdata, nil
}

// Get 获取
func (p *Post) Get() (Post, error) {
	var pdata Post
	table := global.DB.Table(p.TableName())
	if p.ID != 0 {
		table = table.Where("id = ?", p.ID)
	}
	if p.PostName != "" {
		table = table.Where("post_name = ?", p.PostName)
	}
	if p.PostCode != "" {
		table = table.Where("post_code = ?", p.PostCode)
	}
	if p.Status != "" {
		table = table.Where("status = ?", p.Status)
	}
	if err := table.First(&pdata).Error; err != nil {
		return pdata, err
	}
	return pdata, nil
}

// GetList
func (p *Post) GetList() ([]Post, error) {
	var pdata []Post
	table := global.DB.Table(p.TableName())
	if p.ID != 0 {
		table = table.Where("id = ?", p.ID)
	}
	if p.PostName != "" {
		table = table.Where("post_name = ?", p.PostName)
	}
	if p.PostCode != "" {
		table = table.Where("post_code = ?", p.PostCode)
	}
	if p.Status != "" {
		table = table.Where("status = ?", p.Status)
	}
	if err := table.Find(&pdata).Error; err != nil {
		return pdata, err
	}
	return pdata, nil
}

// GetPage
func (p *Post) GetPage(pageSize, pageIndex int) ([]Post, int64, error) {
	var (
		count int64
		pdata []Post
	)

	table := global.DB.Debug().Select("*").Table(p.TableName())
	if p.ID != 0 {
		table = table.Where("id = ?", p.ID)
	}
	if p.PostName != "" {
		table = table.Where("post_name like ?", "%"+p.PostName+"%")
	}
	if p.PostCode != "" {
		table = table.Where("post_code like ?", "%"+p.PostCode+"%")
	}
	if p.Status != "" {
		table = table.Where("status = ?", p.Status)
	}
	if err := table.Debug().Order("sort").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&pdata).Error; err != nil {
		return nil, 0, err
	}
	table.Debug().Where("`delete_time` IS NULL").Count(&count)
	return pdata, count, nil
}

// Update 更新
func (p *Post) Update(id uint64) (update Post, err error) {
	if err = global.DB.Table(p.TableName()).First(&update, id).Error; err != nil {
		return
	}
	if err = global.DB.Table(p.TableName()).Model(&update).Updates(&p).Error; err != nil {
		return
	}
	return
}

func (p *Post) Delete(id int) (success bool, err error) {
	if err = global.DB.Table(p.TableName()).Where("id = ?", id).Delete(&Post{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

func (p *Post) BatchDelete(id []int) (result bool, err error) {
	if err = global.DB.Table(p.TableName()).Where("id in (?)", id).Delete(&Post{}).Error; err != nil {
		return
	}
	result = true
	return
}
