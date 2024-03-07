## Gin-Admin-Vue
基于vue和gin开发的全栈前后端分离的开发基础平台，集成了jwt鉴权、动态路由、动态菜单等功能，总之有助于我们更专注与业务开发，可以快速搭建一套中小型项目。

### 基础架构

- 前端：用基于VUE2的 Element-UI 构建基础页面。
- 后端：用GIN快速搭建,Gin 是一个go语言编写的Web框架。
- 数据库：采用`MySql` > (5.7) 版本 数据库引擎 InnoDB，使用 [gorm](http://gorm.cn/) 实现对数据库的基本操作。
- 缓存：使用`Redis ` > 6+ 实现记录当前活跃用户的`jwt`令牌并实现多点登录限制。
- 配置文件：实现`yaml`格式的配置文件。
- 日志：使用`Logrus`实现日志记录。

## 系统登录
![](https://github.com/Ykubernetes/Gin-Admin-Vue/blob/main/doc/login.png?raw=true)

## 重要提示

本开源项目需要您有一定的golang和vue基础

```TXT
默认账号:admin  默认密码：123456

前端nodejs版本：v16.16.0

后端goland版本：go1.20.11

前端编辑器VSCODE 1.84.2 后端编辑器GOLAND 2023.1.4
```

## 系统管理

### 个人中心

![](https://github.com/Ykubernetes/Gin-Admin-Vue/blob/main/doc/profile.png?raw=true)

### 菜单管理

![](https://github.com/Ykubernetes/Gin-Admin-Vue/blob/main/doc/menu-manager.png?raw=true)

### 角色管理

![](https://github.com/Ykubernetes/Gin-Admin-Vue/blob/main/doc/role-manager.png?raw=true)

#### 角色权限

![](https://github.com/Ykubernetes/Gin-Admin-Vue/blob/main/doc/role-manager-set.png?raw=true)

### 部门管理

![](https://github.com/Ykubernetes/Gin-Admin-Vue/blob/main/doc/department-manager.png?raw=true)

#### 部门设置

![](https://github.com/Ykubernetes/Gin-Admin-Vue/blob/main/doc/department-manager-add.png?raw=true)

### 岗位管理

![](https://github.com/Ykubernetes/Gin-Admin-Vue/blob/main/doc/post-manager.png?raw=true)

### 用户管理

![](https://github.com/Ykubernetes/Gin-Admin-Vue/blob/main/doc/user-manager.png?raw=true)

### 模块管理

#### 模块添加

![](https://github.com/Ykubernetes/Gin-Admin-Vue/blob/main/doc/add-new-model.png?raw=true)

#### 菜单添加

![](https://github.com/Ykubernetes/Gin-Admin-Vue/blob/main/doc/add-new-menu.png?raw=true)

### 客户管理

> 以下客户信息内容为虚拟客户信息

![](https://github.com/Ykubernetes/Gin-Admin-Vue/blob/main/doc/customer-manager.png?raw=true)

#### 信息添加

![](https://github.com/Ykubernetes/Gin-Admin-Vue/blob/main/doc/customer-manager-add.png?raw=true)

#### 信息查看

![](https://github.com/Ykubernetes/Gin-Admin-Vue/blob/main/doc/customer-manager-view.png?raw=true)

#### 信息编辑

![](https://github.com/Ykubernetes/Gin-Admin-Vue/blob/main/doc/customer-manager-edit.png?raw=true)

