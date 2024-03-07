package system

import (
	"errors"
	"fmt"
	"gitee.com/go-server/global"
	"gitee.com/go-server/models/basemodel"
	"strconv"
)

type Dept struct {
	basemodel.Model
	ParentId int    `json:"parentId" gorm:"type:int(11);"`      // 上级部门
	DeptPath string `json:"deptPath" gorm:"type:varchar(255);"` //
	DeptName string `json:"deptName" gorm:"type:varchar(128);"` //部门名称
	Sort     int    `json:"sort" gorm:"type:int(4);"`           // 排序
	Leader   int    `json:"leader" gorm:"type:int(11);"`        // 部门负责人
	Phone    string `json:"phone" gorm:"type:varchar(11);"`     // 手机号码
	Email    string `json:"email" gorm:"type:varchar(64);"`     //邮箱
	Status   string `json:"status" gorm:"type:int(1);"`         //状态
	Params   string `json:"params" gorm:"-"`
	Children []Dept `json:"children" gorm:"-"`
}

func (d Dept) TableName() string {
	return "department"
}

type DeptLable struct {
	Id       int         `gorm:"-" json:"id"`
	Label    string      `gorm:"-" json:"label"`
	Children []DeptLable `gorm:"-" json:"children"`
}

func (d *Dept) Create() (Dept, error) {
	var dep Dept
	result := global.DB.Table(d.TableName()).Create(&d)
	if result.Error != nil {
		err := result.Error
		return dep, err
	}
	deptId := strconv.FormatUint(d.ID, 10) // uint64 转 字符串
	deptPath := "/" + deptId
	fmt.Println("第一次", deptPath)
	if int(d.ParentId) != 0 {
		fmt.Println(">>>>>>>>>>>>>>", d.ParentId)
		var deptP Dept
		// 上级部门id
		global.DB.Table(d.TableName()).Where("id = ?", d.ParentId).First(&deptP)
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>deptP.DeptPath>>>>>>", deptP.DeptPath)
		deptPath = deptP.DeptPath + deptPath
		fmt.Println("有上级部门", deptPath)
	} else {
		deptPath = "/0" + deptPath
		fmt.Println("无上级部门", deptPath)
	}
	var mp = map[string]string{}
	mp["dept_path"] = deptPath
	fmt.Println("准备更新的", deptPath)
	if err := global.DB.Debug().Table(d.TableName()).Where("id = ?", d.ID).Update("dept_path", deptPath).Error; err != nil {
		err = result.Error
		return dep, err
	}
	dep = *d
	dep.DeptPath = deptPath
	fmt.Println("返回的", dep.DeptPath)
	return dep, nil
}

func (d *Dept) Get() (Dept, error) {
	var dep Dept
	table := global.DB.Table(d.TableName())
	if d.ID != 0 {
		table = table.Where("id = ?", d.ID)
	}
	if d.DeptName != "" {
		table = table.Where("dept_name = ?", d.DeptName)
	}
	if err := table.First(&dep).Error; err != nil {
		return dep, err
	}
	return dep, nil
}

func (d *Dept) GetList() ([]Dept, error) {
	var dep []Dept
	table := global.DB.Table(d.TableName())
	if d.ID != 0 {
		table = table.Where("id = ?", d.ID)
	}
	if d.DeptName != "" {
		table = table.Where("dept_name = ?", d.DeptName)
	}
	if d.Status != "" {
		table = table.Where("status = ?", d.Status)
	}
	if err := table.Order("sort").Find(&dep).Error; err != nil {
		return dep, err
	}
	return dep, nil
}

func (d *Dept) GetPage(bl bool) ([]Dept, error) {
	var dep []Dept
	table := global.DB.Select("*").Table(d.TableName())
	if d.ID != 0 {
		table = table.Where("id = ?", d.ID)
	}
	if d.DeptName != "" {
		table = table.Where("dept_name like ?", "%"+d.DeptName+"%")
	}
	if d.Status != "" {
		table = table.Where("status = ?", d.Status)
	}
	if d.DeptPath != "" {
		table = table.Where("dept_path like %?%", d.DeptPath)
	}
	if err := table.Order("sort").Find(&dep).Error; err != nil {
		return nil, err
	}
	return dep, nil
}

func (d *Dept) SetDept(bl bool) ([]Dept, error) {
	list, err := d.GetPage(bl)
	m := make([]Dept, 0)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		info := DiGui(&list, list[i])
		m = append(m, info)
	}
	return m, err
}

func DiGui(deptlist *[]Dept, menu Dept) Dept {
	list := *deptlist

	md := make([]Dept, 0)
	for j := 0; j < len(list); j++ {
		if int(menu.ID) != list[j].ParentId {
			continue
		}
		mi := Dept{}
		mi.ID = list[j].ID
		mi.ParentId = list[j].ParentId
		mi.DeptPath = list[j].DeptPath
		mi.DeptName = list[j].DeptName
		mi.Sort = list[j].Sort
		mi.Leader = list[j].Leader
		mi.Phone = list[j].Phone
		mi.Email = list[j].Email
		mi.Status = list[j].Status
		mi.CreatedAt = list[j].CreatedAt
		mi.UpdatedBy = list[j].UpdatedBy
		mi.CreatedBy = list[j].CreatedBy
		mi.UpdatedAt = list[j].UpdatedAt
		mi.Children = []Dept{}
		ms := DiGui(deptlist, mi)
		md = append(md, ms)
	}
	menu.Children = md
	return menu
}

func (d *Dept) Update(id uint64) (update Dept, err error) {
	if err = global.DB.Table(d.TableName()).Where("id = ?", id).First(&update).Error; err != nil {
		return
	}
	deptId := strconv.FormatUint(d.ID, 10) // uint64 转 字符串
	deptPath := "/" + deptId
	fmt.Println("第一次:", deptPath)
	if int(d.ParentId) != 0 {
		var deptP Dept
		global.DB.Table(d.TableName()).Where("id = ?", d.ParentId).First(&deptP)
		deptPath = deptP.DeptPath + deptPath
		fmt.Println("有上级部门:", deptPath)
	} else {
		deptPath = "/0" + deptPath
		fmt.Println("没有上级部门:", deptPath)
	}
	d.DeptPath = deptPath

	if d.DeptPath != "" && d.DeptPath != update.DeptPath {
		return update, errors.New("上级部门不允许修改！")
	}
	fmt.Println("更新前:", d.DeptPath)
	if err = global.DB.Table(d.TableName()).Model(&update).Updates(&d).Error; err != nil {
		return
	}
	return
}

func (d *Dept) Delete(id int) (b bool, err error) {
	user := Admins{}
	user.DeptId = id
	userlist, err := user.GetUserListAndRoleId()
	HasError(err, "", 500)
	if !(len(userlist) <= 0) {
		return false, errors.New("当前部门存在用户，不能删除")
	}
	if err = global.DB.Table(d.TableName()).Where("id = ?", id).Delete(&Dept{}).Error; err != nil {
		global.Log.Info(err)
		return false, errors.New("删除部门失败")
	}
	return true, nil
}

// HasError 错误断言
// 当 error 不为 nil 时触发 panic
// 对于当前请求不会再执行接下来的代码，并且返回指定格式的错误信息和错误码
// 若 msg 为空，则默认为 error 中的内容
func HasError(err error, msg string, code ...int) {
	if err != nil {
		statusCode := 200
		if len(code) > 0 {
			statusCode = code[0]
		}
		if msg == "" {
			msg = err.Error()
		}
		global.Log.Info(err)
		panic("CustomError#" + strconv.Itoa(statusCode) + "#" + msg)
	}
}

func (d *Dept) SetDeptLable() (m []DeptLable, err error) {
	deptlist, err := d.GetList()
	m = make([]DeptLable, 0)
	for i := 0; i < len(deptlist); i++ {
		if deptlist[i].ParentId != 0 {
			continue
		}
		e := DeptLable{}
		e.Id = int(deptlist[i].ID)
		e.Label = deptlist[i].DeptName
		deptsInfo := DiGuiDeptLable(&deptlist, e)
		m = append(m, deptsInfo)
	}
	return
}

func DiGuiDeptLable(deptlist *[]Dept, dept DeptLable) DeptLable {
	list := *deptlist
	md := make([]DeptLable, 0)
	for j := 0; j < len(list); j++ {
		if dept.Id != list[j].ParentId {
			continue
		}
		mi := DeptLable{int(list[j].ID), list[j].DeptName, []DeptLable{}}
		ms := DiGuiDeptLable(deptlist, mi)
		md = append(md, ms)
	}
	dept.Children = md
	return dept
}
