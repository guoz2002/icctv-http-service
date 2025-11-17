## 📁 项目目录结构说明（icctv-http-service）
icctv-http-service
├─ main.go                # 程序入口，启动 HTTP 服务
├─ go.mod                 # Go 模块配置
├─ readme.md              # API & 模型文档
├─ struct.md              # 项目结构说明
│
├─ controllers/           # 控制器层：接 HTTP 请求
│  ├─ management_*.go     # 管理相关接口
│  ├─ production_*.go     # 生产相关接口
│  └─ ...                 # 未来的 admin/orangepi/building 控制器
│
├─ services/              # 服务层：业务逻辑
│  ├─ admin_service.go    # 管理员相关逻辑
│  ├─ orangepi_service.go # 设备管理逻辑
│  ├─ building_service.go # 建筑管理逻辑
│  └─ ...                 # 其他业务
│
├─ models/                # 模型层：数据结构 & ORM
│  ├─ model.go            # ModelFields 等通用字段
│  ├─ adminer.go          # Adminer 模型
│  ├─ orangepi.go         # OrangePi 模型
│  ├─ building.go         # Building 模型
│  └─ ...                 # 其他表模型
│
├─ databases/             # 数据库初始化 & 连接
│  └─ db.go
│
├─ routes/                # 路由注册（URL → controller）
│  └─ routes.go
│
├─ middlewares/           # 中间件（JWT、日志、CORS 等）
│  └─ jwt.go
│
└─ utils/                 # 工具函数
   ├─ gorm_pagination_order.go
   └─ struct_convert.go

项目整体是一个基于 Go 的 HTTP 服务，按职责拆分为控制层、服务层、模型层等，方便后续扩展 OrangPi 管理后台功能。

## 🧱 顶层文件

- **`main.go`**: 服务入口，启动 HTTP Server（默认端口 `8080`），挂载路由。
- **`go.mod`**: Go Module 配置，定义模块名与 Go 版本。
- **`readme.md`**: 项目 API 文档与数据模型说明。
- **`struct.md`**: 当前这个文件，用于描述项目目录结构和职责划分。

## 📂 目录结构

- **`controllers/`**: 控制器层
  - 负责具体 HTTP 接口的请求处理（解析参数、调用 service、返回 JSON）。
  - 按业务拆分为管理、生产等控制器文件，后续可新增 `admin_controller.go`、`orangepi_controller.go` 等。

- **`databases/`**: 数据库访问与初始化
  - 封装数据库连接（例如 MySQL、PostgreSQL、SQLite 等）的初始化。
  - 提供统一的 `DB` 实例给 `models` 和 `services` 使用。

- **`middlewares/`**: 中间件
  - 放置 JWT 鉴权、日志、跨域等 HTTP 中间件。
  - 可在这里实现接口级别的权限控制（如 Adminer 验证）。

- **`models/`**: 数据模型层
  - 存放所有业务实体定义，例如：
    - `ModelFields` 通用字段（ID、CreatedAt、UpdatedAt、DeletedAt）。
    - `Adminer` 管理员账号模型。
    - `OrangePi` 设备模型。
    - `Building` 楼栋/建筑模型。
    - `PublicNet` 公网配置等。
  - 同时可包含数据库迁移、校验逻辑等辅助代码。

- **`routes/`**: 路由注册
  - 负责将 URL 路径映射到对应的控制器 handler。
  - 示例：将 `/api/admin` 映射到 `Adminer` 相关控制器函数。

- **`services/`**: 业务服务层
  - 承载具体业务逻辑（不直接处理 HTTP），由 `controllers` 调用。
  - 典型职责：
    - 处理 Adminer 登录、Token 生成。
    - OrangePi/Building/PublicNet 的 CRUD 业务规则。
    - 调用外部服务或本地进程（如 FRPC 管理）。

- **`utils/`**: 工具与通用辅助函数
  - 放一些与业务无关的通用方法，例如：
    - 结构体转换、分页、排序工具。
    - 时间处理、字符串工具等。

## 🎯 使用建议

- 新增接口时：**优先在 `controllers/` + `services/` + `models/` 中分层实现**，避免所有逻辑都堆在 `main.go`。
- 修改数据结构时：优先修改 `models/` 中的模型定义，再在 `services/` 中调整对应业务逻辑。


