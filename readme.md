# OrangePi 后台管理系统 - HTTP Service API 文档

## � 项目简介

本项目是一个基于 Go 语言开发的 OrangePi 设备后台管理系统的 HTTP 服务端,提供设备管理、建筑信息管理、录像文件管理、公网配置等功能。

## �📋 完整API接口列表

### 核心管理接口

| 编号 | 接口 | 方法 | 功能 | 使用场景 | 权限要求 |
|------|------|------|------|----------|----------|
| 1 | `/health` | GET | 健康检查 | 应用启动时检查服务状态 | 无 |
| 2 | `/api/auth/public` | GET | 获取公开访问Token | 第三方系统集成 | 无(有速率限制) |
| 3 | `/api/auth/login` | POST | 生成测试Token | 开发测试环境使用 | 管理员权限 |

### 管理员账户管理 (Adminer)

| 编号 | 接口 | 方法 | 功能 | 权限要求 |
|------|------|------|------|----------|
| 4 | `/api/admin` | GET | 查询管理员列表/详情 | 超级管理员 |
| 5 | `/api/admin` | POST | 创建管理员账户 | 超级管理员 |
| 6 | `/api/admin` | PUT | 更新管理员信息 | 超级管理员 |
| 7 | `/api/admin` | DELETE | 删除管理员账户 | 超级管理员 |

### OrangePi设备管理

| 编号 | 接口 | 方法 | 功能 | 权限要求 |
|------|------|------|------|----------|
| 8 | `/api/orangepi` | GET | 查询OrangePi设备列表/详情 | 管理员 |
| 9 | `/api/orangepi` | POST | 创建OrangePi设备记录 | 管理员 |
| 10 | `/api/orangepi` | PUT | 更新OrangePi设备信息 | 管理员 |
| 11 | `/api/orangepi` | DELETE | 删除OrangePi设备记录 | 管理员 |

### 建筑信息管理 (Building)

| 编号 | 接口 | 方法 | 功能 | 权限要求 |
|------|------|------|------|----------|
| 12 | `/api/building` | GET | 查询建筑列表/详情 | 管理员 |
| 13 | `/api/building` | POST | 创建建筑信息 | 管理员 |
| 14 | `/api/building` | PUT | 更新建筑信息 | 管理员 |
| 15 | `/api/building` | DELETE | 删除建筑信息 | 管理员 |

### 设备与网络配置

| 编号 | 接口 | 方法 | 功能 | 权限要求 |
|------|------|------|------|----------|
| 16 | `/api/device/info` | GET | 获取设备信息 | 操作员 |
| 17 | `/api/device/ports` | POST | 远程更新FRPC端口并重启 | 管理员 |
| 19 | `/api/publicnet/config` | PUT | 修改公网配置 | 管理员 |



## 📊 数据模型规范

### 0. 通用字段 ModelFields

所有业务模型统一继承 `ModelFields`，用于存储主键与时间戳：

```go
type ModelFields struct {
    ID        int64          `gorm:"primary_key" json:"id"`           // 主键ID
    CreatedAt time.Time      `json:"createdAt"`                       // 创建时间
    UpdatedAt time.Time      `json:"updatedAt"`                       // 更新时间
    DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`// 软删除时间
}
```

后续的 `Adminer`、`OrangePi`、`Building` 均通过匿名字段方式引入：

```go
type ExampleModel struct {
    ModelFields
    // ...业务字段...
}
```

### 1. Adminer(管理员) 模型

```go
type Adminer struct {
    ModelFields
    Username       string    `json:"username"`         // 登录用户名(唯一)
    PasswordHash   string    `json:"password_hash"`    // Bcrypt 哈希
}
```

### 2. OrangePi (设备) 模型

```go
type OrangePi struct {
    ModelFields

    Base                        string `json:"base"`          // 关联楼栋 base
    Name                        string `json:"name"`          // Orangepi 名称
    ICCTVAuthServiceRemotePort  int    `json:"icctv_auth_service_remote_port"` // 远程认证服务
    SSHRemotePort               int    `json:"ssh_remote_port"` // SSH 远程端口
    AdminPorts                  []int  `json:"admin_ports"`   // 可用管理端口列表(1~6)
    IsActive                    bool   `json:"is_active"`     // 是否在用
}
```

### 3. Building (建筑) 模型

```go
type Building struct {
    ModelFields

    Base     string `json:"base"`     // 物理园区/楼栋唯一标识
    ISmartID string `json:"ismartid"` // ismart 系统ID
    Name     string `json:"name"`     // 楼栋名称
    Remark   string `json:"remark"`   // 备注信息
    orangepis [] "json:orangepis[]"
}
```

### 4. PublicNetConfig (公网配置) 模型

```go
type PublicNetConfig struct {
    IP string `json:"external_ip"`  // 外部IP地址
}
```



