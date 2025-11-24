# 测试脚本使用说明

## test_nvr_api.ps1 - 完整 API 测试脚本

### 📝 功能说明

该脚本会按顺序测试以下所有 API 接口:

#### 1. 核心管理接口
- ✅ `GET /health` - 健康检查
- ✅ `POST /api/auth/login` - 管理员登录

#### 2. 管理员账户管理 (Adminer)
- ✅ `GET /api/admin` - 查询管理员列表
- ✅ `POST /api/admin` - 创建管理员账户
- ✅ `PUT /api/admin` - 更新管理员信息
- ✅ `DELETE /api/admin` - 删除管理员账户

#### 3. 建筑信息管理 (Building)
- ✅ `GET /api/building` - 查询建筑列表
- ✅ `POST /api/building` - 创建建筑信息
- ✅ `PUT /api/building` - 更新建筑信息
- ✅ `DELETE /api/building` - 删除建筑信息

#### 4. OrangePi 设备管理
- ✅ `GET /api/device` - 查询OrangePi设备列表
- ✅ `POST /api/device` - 创建OrangePi设备记录
- ✅ `PUT /api/device` - 更新OrangePi设备信息
- ✅ `DELETE /api/device` - 删除OrangePi设备记录

#### 5. NVR 管理 (包含 RTSPUrls 字段)
- ✅ `GET /api/nvr` - 查询NVR列表
- ✅ `GET /api/nvr?id=X` - 查询NVR详情
- ✅ `POST /api/nvr` - 创建NVR (含 RTSPUrls)
- ✅ `PUT /api/nvr?id=X` - 更新NVR (含 RTSPUrls)
- ✅ `DELETE /api/nvr` - 删除NVR

#### 6. Bind 关联关系管理
- ✅ `POST /api/bind/building-orangepi` - 绑定Building和OrangePi
- ✅ `GET /api/bind/building-orangepi/{building_id}` - 查询Building关联的OrangePi
- ✅ `DELETE /api/bind/building-orangepi` - 解绑Building和OrangePi
- ✅ `POST /api/bind/building-nvr` - 绑定Building和NVR
- ✅ `GET /api/bind/building-nvr/{building_id}` - 查询Building关联的NVR
- ✅ `DELETE /api/bind/building-nvr` - 解绑Building和NVR

---

## 🚀 使用方法

### 1. 基本用法 (使用默认参数)

```powershell
.\test_nvr_api.ps1
```

默认参数:
- BaseUrl: `http://127.0.0.1:8080`
- Username: `admin`
- Password: `123456`

### 2. 自定义参数

```powershell
# 指定不同的服务地址
.\test_nvr_api.ps1 -BaseUrl "http://192.168.1.100:8080"

# 指定不同的登录凭据
.\test_nvr_api.ps1 -Username "admin" -Password "mypassword"

# 组合使用
.\test_nvr_api.ps1 -BaseUrl "http://192.168.1.100:8080" -Username "admin" -Password "mypassword"
```

---

## 📊 输出说明

### 成功输出示例

```
==========================================
ICCTV HTTP Service - Complete API Test
==========================================

========== 核心管理接口 ==========

[1] Health Check
    ✓ SUCCESS

[2] Admin Login
    ✓ SUCCESS

========== 管理员账户管理 (Adminer) ==========

[3] Query Admin List
    ✓ SUCCESS

[4] Create Admin
    ✓ SUCCESS
    Created Admin ID: 123

...

==========================================
Test Summary
==========================================

Total:  25
Passed: 25
Failed: 0

✅ All tests PASSED!
```

### 失败输出示例

```
[5] Create Building
    ✗ FAILED: 401 Unauthorized

==========================================
Test Summary
==========================================

Total:  10
Passed: 9
Failed: 1

❌ Some tests FAILED!
```

---

## 🧪 测试内容详解

### NVR 测试 (重点)

脚本会创建包含完整 RTSPUrls 数据的 NVR:

```json
{
    "name": "测试NVR设备",
    "url": "192.168.1.100:8080",
    "building_id": 1,
    "admin_user": {
        "name": "admin",
        "password": "admin123"
    },
    "users": [
        {"name": "operator1", "password": "pass123"},
        {"name": "operator2", "password": "pass456"}
    ],
    "rtsp_urls": [
        {"channel": 1, "url": "rtsp://192.168.1.100:554/stream1"},
        {"channel": 2, "url": "rtsp://192.168.1.100:554/stream2"},
        {"channel": 3, "url": "rtsp://192.168.1.100:554/stream3"},
        {"channel": 4, "url": "rtsp://192.168.1.100:554/stream4"}
    ]
}
```

然后测试更新 RTSPUrls:

```json
{
    "name": "更新后的NVR设备",
    "rtsp_urls": [
        {"channel": 1, "url": "rtsp://192.168.1.100:554/updated_stream1"},
        {"channel": 2, "url": "rtsp://192.168.1.100:554/updated_stream2"},
        {"channel": 5, "url": "rtsp://192.168.1.100:554/new_stream5"}
    ]
}
```

---

## 🗑️ 清理测试数据

脚本执行完成后会询问是否删除测试数据:

```
是否删除测试数据? (y/n):
```

- 输入 `y` 或 `Y`: 删除所有测试数据
- 输入 `n` 或 `N`: 保留测试数据,并显示创建的资源ID

保留数据时的输出:
```
测试数据保留:
  Admin ID: 1234
  Building ID: 5678
  OrangePi ID: 9012
  NVR ID: 3456
```

---

## ⚠️ 注意事项

### 1. 前置条件
- 确保服务已启动: `go run main.go`
- 确保数据库连接正常
- 确保管理员账户存在 (username: admin, password: 123456)

### 2. 网络要求
- 脚本需要能访问 `BaseUrl` 指定的服务地址
- 默认超时时间为 100 秒

### 3. 权限要求
- 登录账户需要管理员权限
- 所有测试接口都需要认证

### 4. 数据安全
- 测试会创建真实数据
- 建议在测试环境运行
- 生产环境请谨慎使用

---

## 🔧 故障排查

### 问题1: 连接失败
```
Login FAILED: 无法连接到远程服务器
```
**解决方案:**
- 检查服务是否启动
- 确认 BaseUrl 是否正确
- 检查防火墙设置

### 问题2: 认证失败
```
Login FAILED: 401 Unauthorized
```
**解决方案:**
- 确认用户名密码是否正确
- 检查数据库中是否存在该管理员账户

### 问题3: 创建失败
```
Create Building FAILED: building_id already exists
```
**解决方案:**
- 数据已存在,可忽略
- 或修改脚本中的随机ID生成范围

---

## 📈 扩展使用

### 自动化集成

```powershell
# 在 CI/CD 中使用
$result = .\test_nvr_api.ps1
if ($LASTEXITCODE -ne 0) {
    Write-Error "API tests failed"
    exit 1
}
```

### 性能测试

```powershell
# 多次运行测试
for ($i = 1; $i -le 10; $i++) {
    Write-Host "Run $i"
    .\test_nvr_api.ps1
}
```

---

## 📞 技术支持

如有问题,请检查:
1. 服务日志: `logs/*.log`
2. 数据库状态: MySQL 连接和表结构
3. API 文档: `readme.md`

---

**最后更新:** 2025-11-24
