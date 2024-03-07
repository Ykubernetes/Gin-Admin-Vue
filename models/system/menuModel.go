package system

type MenuMeta struct {
	Title   string `json:"title"`   // 标题
	Icon    string `json:"icon"`    // 图标
	NoCache bool   `json:"noCache"` // 是不是缓存
}

type MenuModel struct {
	Path      string      `json:"path"`      // 路由
	Component string      `json:"component"` // 对应vue中的map name
	Name      string      `json:"name"`      // 菜单名称
	Hidden    bool        `json:"hidden"`    // 是否隐藏
	Meta      MenuMeta    `json:"meta"`      // 菜单信息
	Children  []MenuModel `json:"children"`  // 子级菜单
}

type UserData struct {
	Menus        []MenuModel `json:"menus"`        // 菜单
	Introduction string      `json:"introduction"` // 介绍
	Avatar       string      `json:"avatar"`       // 图标
	Name         string      `json:"name"`         // 姓名
	Email        string      `json:"email"`
	//Roles        []string    `json:"roles"`
}
