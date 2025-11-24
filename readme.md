# OrangePi 后台管理系统 - HTTP Service API 文档

## � 项目简介

本项目是一个基于 Go 语言开发的 OrangePi 设备后台管理系统的 HTTP 服务端,提供设备管理、建筑信息管理、录像文件管理、公网配置等功能。

## �📋 完整API接口列表

### 核心管理接口

| 编号 | 接口 | 方法 | 功能 | 使用场景 | 权限要求 |
|------|------|------|------|----------|----------|
| 1 | `/health` | GET | 健康检查 | 应用启动时检查服务状态 | 无 |
| 2 | `/api/auth/public` | GET | 获取公开访问Token | 第三方系统集成 | 无(有速率限制) |(TODO)
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

### nvr(网络硬盘录像机) (nvr)

| 编号 | 接口 | 方法 | 功能 | 权限要求 |
|------|------|------|------|----------|
| 12 | `/api/nvr` | GET | 查询nvr列表/详情 | 管理员 |
| 13 | `/api/nvr` | POST | 创建nvr信息 | 管理员 |
| 14 | `/api/nvr` | PUT | 更新nvr信息 | 管理员 |
| 15 | `/api/nvr` | DELETE | 删除nvr信息 | 管理员 |



### 设备与网络配置(TODO)

| 编号 | 接口 | 方法 | 功能 | 权限要求 |
|------|------|------|------|----------|
| 16 | `/api/device/info` | GET | 获取设备信息 | 操作员 |（TODO）
| 17 | `/api/device/ports` | POST | 远程更新FRPC端口并重启 | 管理员 |（TODO）
| 19 | `/api/publicnet/config` | PUT | 修改公网配置 | 管理员 |（TODO）



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
    
    // 关联关系 (一对多)
    OrangePis []OrangePi `json:"orangepis,omitempty"` // 关联的OrangePi设备列表
    NVRs      []NVR      `json:"nvrs,omitempty"`      // 关联的NVR设备列表
}
```

### 4. PublicNetConfig (公网配置) 模型

```go
type PublicNetConfig struct {
    IP string `json:"external_ip"`  // 外部IP地址
}
```

### 5. NVR (网络硬盘录像机) 模型

#### 用户认证信息子模型

```go
// AdminUser 管理员账户信息
type AdminUser struct {
    Name     string `json:"name"`     // 管理员用户名
    Password string `json:"password"` // 管理员密码
}

// User 普通用户账户信息
type User struct {
    Name     string `json:"name"`     // 用户名
    Password string `json:"password"` // 密码
}

// ChannelURL RTSP频道地址
type ChannelURL struct {
    Channel int    `json:"channel"` // 通道号
    URL     string `json:"url"`     // RTSP 地址
}
```

#### NVR 完整模型

```go
type NVR struct {
    ModelFields
    
    Name       string       `json:"name"`        // NVR 名称
    URL        string       `json:"url"`         // NVR 访问地址 (IP:Port)
    BuildingID int64        `json:"building_id"` // 关联建筑ID (外键)
    AdminUser  AdminUser    `json:"admin_user"`  // 管理员账户
    Users      []User       `json:"users"`       // 普通用户列表
    RTSPUrls   []ChannelURL `json:"rtsp_urls"`   // RTSP 地址列表
    
    // 关联关系 (多对一)
    Building   Building     `json:"building,omitempty"` // 所属建筑
}
```

**字段说明:**
- `name`: NVR 设备的显示名称
- `url`: NVR 的管理访问地址，格式为 `IP:Port` (如: `192.168.1.100:80`)
- `building_id`: 关联的建筑ID
- `admin_user`: 管理员账户信息
- `users`: 普通用户列表
- `rtsp_urls`: RTSP流地址列表，每个包含通道号和对应的RTSP URL

**完整数据示例:**

```json
{
    "id": 1,
    "name": "A楼监控系统",
    "url": "192.168.1.100:8080",
    "building_id": 5,
    "admin_user": {
        "name": "admin",
        "password": "admin123"
    },
    "users": [
        {
            "name": "operator1",
            "password": "pass123"
        },
        {
            "name": "operator2",
            "password": "pass456"
        }
    ],
    "rtsp_urls": [
        {
            "channel": 1,
            "url": "rtsp://192.168.1.100:554/stream1"
        },
        {
            "channel": 2,
            "url": "rtsp://192.168.1.100:554/stream2"
        },
        {
            "channel": 3,
            "url": "rtsp://192.168.1.100:554/stream3"
        }
    ],
    "building": {
        "id": 5,
        "name": "办公楼A栋"
    },
    "createdAt": "2025-11-24T10:00:00Z",
    "updatedAt": "2025-11-24T10:00:00Z"
}
```