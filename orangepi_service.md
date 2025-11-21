# 🎯 iCCTV OrangePi 认证服务

## 📋 项目简介

这是一个部署在 OrangePi 上的高性能认证网关服务，为 MediaMTX 视频流提供 Token 验证、权限控制和录像回放功能。

**核心功能**:
- ✅ Token 验证和权限控制（基于 buildings 的细粒度权限）
- ✅ 实时视频流反向代理（WebRTC/WHEP/WHIP）
- ✅ **实时录像下载**（NVR → MediaMTX → Progressive MP4）🔥
- ✅ 录像查询和管理
- ✅ 设备信息和状态监控

**技术特点**:
- 🚀 FastAPI 高性能异步框架
- 📦 模块化架构，易于扩展
- 🔒 多层安全验证
- 📊 完整的 API 文档
- 🐳 Docker 容器化部署

---

## 📋 完整API接口列表

| 编号 | 接口 | 方法 | 功能 | 使用场景 |
|------|------|------|------|----------|
| 1 | `/health` | GET | 健康检查 | 应用启动时检查服务状态 |
| 2 | `/api/device/info` | GET | 获取设备信息 | 设备状态监控页面 |
| 3 | `/api/auth/generate-token` | POST | 生成测试Token | 开发测试环境使用 |
| 4 | `/api/device/frpc/ports` | POST | **远程更新FRPC端口** | **OrangePi公网维护** |

**公网访问地址**: `http://39.108.49.167:29001` (通过FRP内网穿透)

---

## 📡 API 接口详解

> 下述编号与“完整API接口列表”一致，逐条展示接口含义、调试方式与响应示例。

### 1. GET /health

- **功能**：检查 Python 服务、MediaMTX API 以及 Docker 容器（mediamtx / frpc / ismart_auth_service）健康度。
- **测试命令**
  ```powershell
  Invoke-RestMethod -Uri "http://localhost:8889/health"
  ```
  ```bash
  curl http://39.108.49.167:29005/health
  ```
- **响应示例**
  ```json
  {
    "status": "healthy",
    "service": "mediamtx-auth",
    "docker_services": { "mediamtx": true, "frpc": true, "ismart_auth_service": true },
    "mediamtx_status": "connected",
    "frpc_status": "running"
  }
  ```

### 2. GET /api/device/info

- **功能**：返回 OrangePi 设备唯一标识（基于硬件指纹自动生成）、FRPC 服务器+端口、可用频道、录像开关及容器整体状态。
- **测试命令**
  ```powershell
  Invoke-RestMethod -Uri "http://localhost:8889/api/device/info"
  ```
- **响应示例**
  ```json
  {
    "device_id": "9d9f4da4-5f6b-5c4f-8551-3f6c7f73b052",
    "mediamtx_version": "v1.15.3",
    "frpc_server": "39.108.49.167",
    "frpc_auth_remote_port": 29005,
    "frpc_ssh_remote_port": 30005,
    "available_channels": ["channel1","channel2","channel3","channel4","channel5","channel6"],
    "status": "online"
  }
  ```

### 3. POST /api/auth/generate-token

- **功能**：在 DEBUG 模式下生成测试 Token（payload 与 `generate_token.py` 一致）。
- **请求体**
  ```json
  {
    "channels": [1, 2, 3],
    "building_id": "0314100"
  }
  ```
- **测试命令**
  ```powershell
  Invoke-RestMethod -Uri "http://localhost:8889/api/auth/generate-token" `
    -Method Post -ContentType "application/json" `
    -Body '{"channels":[1,2,3],"building_id":"0314100"}'
  ```

### 4. POST /api/device/frpc/ports

- **功能**：更新 `config/frpc.toml` 中认证/SSH 远程端口并调用 Docker API 重启 `frpc`，便于远程维护公网映射。
- **请求体**
  ```json
  {
    "icctv_orangepi_auth_remote_port": 29005,
    "orangepi_ssh_remote_port": 30005
  }
  ```
- **测试命令**
  ```powershell
  $body = @{ icctv_orangepi_auth_remote_port = 29005; orangepi_ssh_remote_port = 30005 } | ConvertTo-Json
  Invoke-RestMethod -Uri "http://localhost:8889/api/device/frpc/ports" -Method Post -Body $body -ContentType "application/json"
  ```

## 📚 完整使用示例

### 场景1: 校验权限并获取设备信息

```powershell
# 生成测试 token（与 generate_token.py 格式一致）
Invoke-RestMethod -Uri "http://localhost:8889/api/auth/generate-token" -Method Post -Body @{
  channels = @(1, 2, 3)  # 可访问的频道列表
  building_id = "0314100"  # 建筑ID
} | ConvertTo-Json -Depth 4

# 查询设备信息（包含可访问频道）
$token = "YOUR_TOKEN"
Invoke-RestMethod -Uri "http://localhost:8889/api/device/info?token=$token"
```

> 返回的 `available_channels` 即用户允许访问的回放/下载通道。

---

### 场景2: 远程更新FRPC端口

```powershell
$token = "YOUR_STAFF_TOKEN"
$body = @{
  icctv_orangepi_auth_remote_port = 29005  # 新的认证服务端口
  orangepi_ssh_remote_port = 30005         # 新的SSH远程端口
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8889/api/device/frpc/ports?token=$token" `
  -Method Post -ContentType "application/json" -Body $body
```

> 返回值会包含 `restarted=true/false`，用于判断 `frpc` 是否成功拉起。如重启失败，可手动执行 `docker-compose restart frpc`。

---

## 🔧 开发调试指南

### 本地开发环境搭建

1. **克隆项目**

```bash
git clone https://github.com/your-org/icctv_orangepi_service.git
cd icctv_orangepi_service
```

2. **安装依赖**

```bash
pip install -r auth_service/requirements.txt
```

3. **配置环境变量**

创建 `.env` 文件：

```bash
# 服务配置
DEBUG=true
TOKEN_MODE=simple
SECRET_KEY=your-secret-key-here
DEVICE_ID=orangepi-001

# MediaMTX配置
MEDIAMTX_HOST=localhost
MEDIAMTX_API_PORT=9997
MEDIAMTX_WEBRTC_PORT=8890
```

4. **启动服务**

```bash
# 启动MediaMTX（需要先安装）
./mediamtx

# 启动认证服务
cd auth_service
python app.py
```

5. **运行测试**

```bash
# 测试实时视频流接口
python test_channel_access.py

# 测试完整API
python test_complete_apis.py

# 测试设备信息
python test_device_info.py
```

### Docker部署

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f ismart_auth_service

# 重启服务
docker-compose restart ismart_auth_service

# 停止服务
docker-compose down
```

> 🛠️ **重要**：`ismart_auth_service` 需要同时挂载 `./config/frpc.toml:/app/config/frpc.toml` 与 `/var/run/docker.sock:/var/run/docker.sock`，并设置 `FRPC_CONFIG_PATH=/app/config/frpc.toml`、`DOCKER_SOCK_PATH=/var/run/docker.sock`、`FRPC_SERVICE_NAME=frpc`，这样远程端口更新接口才能落盘并重启 FRPC。

### 常见问题排查

#### 问题1: Token验证失败

**现象**: 访问channel时返回401

**排查步骤**:
1. 检查SECRET_KEY是否一致
2. 验证Token格式是否正确（base64.signature）
3. 检查Token是否过期（exp字段）
4. 查看服务日志确认错误详情

```bash
docker-compose logs ismart_auth_service | grep "Token"
```

#### 问题2: 权限验证失败

**现象**: 访问channel时返回403

**排查步骤**:
1. 确认用户的channels数组（如 [1, 2, 3] 表示可访问 channel1, channel2, channel3）
2. 检查channel命名是否正确（channel1, channel2 等）
3. 验证building_id字段值

**测试Token payload**:

```powershell
# 解码Token查看payload
$token = "YOUR_TOKEN"
$parts = $token -split '\.'
$payloadB64 = $parts[0]
$payloadJson = [System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String($payloadB64))
Write-Host $payloadJson
```

#### 问题3: 无法连接MediaMTX

**现象**: 服务启动但无法访问视频流

**排查步骤**:
1. 检查MediaMTX是否运行
2. 验证端口配置是否正确
3. 检查网络连接

```bash
# 检查MediaMTX API
curl http://localhost:9997/v3/config/global/get

# 检查MediaMTX paths
curl http://localhost:9997/v3/paths/list
```

---