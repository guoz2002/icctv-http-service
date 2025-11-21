# 🧹 MySQL 日志文件清理方案

## 问题

MySQL 在 `/data/mysql/` 目录中生成了大量日志文件，占用磁盘空间：

```
binlog.000001 - 000004    # 二进制日志（2.8MB + 6.8KB + 180B + 25KB）
general_log.*              # 通用查询日志
slow_log.*                 # 慢查询日志
#ib_redo*, #ib_16384*      # InnoDB 日志
```

## ✅ 解决方案

已修改 `docker-compose.yml`，添加以下 MySQL 参数来禁用日志生成：

### 修改内容

```yaml
command:
  - --character-set-server=utf8mb4
  - --collation-server=utf8mb4_unicode_ci
  - --skip-log-bin              # ❌ 禁用二进制日志
  - --general-log=0             # ❌ 禁用通用查询日志
  - --slow-query-log=0          # ❌ 禁用慢查询日志
```

## 📊 参数说明

| 参数 | 作用 | 影响 |
|------|------|------|
| `--skip-log-bin` | 禁用二进制日志 | 减少磁盘占用，但无法进行主从复制 |
| `--general-log=0` | 禁用通用查询日志 | 不记录所有 SQL 查询（开发环境可接受） |
| `--slow-query-log=0` | 禁用慢查询日志 | 不记录执行时间长的查询 |

## 🗑️ 清理现有日志

由于你已经有了这些日志文件在 `/data/mysql/` 目录中，建议清理它们：

### 方案 1: 手动删除（推荐）

```bash
# 停止 MySQL 容器
docker stop icctv-mysql

# 删除旧日志文件
cd data/mysql
rm -f binlog.*
rm -f general_log.*
rm -f slow_log.*
rm -f #ib_redo*
rm -f #ib_16384*

# 重启 MySQL
docker start icctv-mysql
```

### 方案 2: 完全清除（谨慎）

```bash
# 停止容器
docker stop icctv-mysql

# 清空整个数据目录（会丢失数据库数据！）
rm -rf data/mysql/*

# 重启容器（会重新初始化数据库）
docker-compose up -d
```

### 方案 3: PowerShell 清理

```powershell
# 停止容器
docker stop icctv-mysql

# 删除日志文件
cd "C:\Users\20216\Documents\GitHub\icctv-http-service\data\mysql"
Remove-Item "binlog.*"
Remove-Item "general_log.*"
Remove-Item "slow_log.*"
Remove-Item "#ib_redo*"
Remove-Item "#ib_16384*"

# 重启
docker start icctv-mysql
```

## 🔄 重启 MySQL 应用修改

修改 `docker-compose.yml` 后，需要重启 MySQL 容器才能生效：

```bash
# 方案 A: 只重启 MySQL
docker-compose restart mysql

# 方案 B: 重新部署
docker-compose down
docker-compose up -d
```

## 📈 效果对比

### 修改前（有日志）
```
data/mysql/
  ├── binlog.000001 (2.8MB)
  ├── binlog.000002 (6.8KB)
  ├── binlog.000003 (180B)
  ├── binlog.000004 (25KB)
  ├── general_log.CSV (大)
  ├── slow_log.CSV (大)
  ├── #ib_redo10_tmp
  ├── #ib_redo11_tmp
  ├── ... (更多临时文件)
  └── mysql/ (实际数据库文件)
```

### 修改后（仅数据）
```
data/mysql/
  ├── icctv_http_service/ (你的数据库)
  ├── mysql/              (系统数据库)
  ├── ibdata1             (表空间数据)
  └── mysql.ibd           (系统表数据)
```

## ⚠️ 注意事项

### ❌ 禁用日志的缺点

1. **无法进行主从复制** - 如果你需要 MySQL 主从架构，需要保留二进制日志
2. **无法恢复到指定时间点** - 没有二进制日志无法通过 binlog 进行时间点恢复
3. **无法调试性能问题** - 没有慢查询日志，难以找到性能瓶颈

### ✅ 开发环境通常可以接受

对于开发/测试环境，这些限制通常可以接受，因为：
- 不需要高可用性
- 不需要复杂的数据恢复
- 优先考虑节省磁盘空间

### 生产环境建议

如果用于生产环境，建议：
1. 保留二进制日志（用于备份和恢复）
2. 定期归档和清理日志文件
3. 使用单独的日志卷/分区
4. 配置日志轮换策略

## 🔍 验证修改

修改后运行测试，确认没有新的日志文件生成：

```powershell
# 查看 data/mysql 目录大小变化
Get-ChildItem -Path "data/mysql" -Recurse | Measure-Object -Sum -Property Length

# 运行测试
.\test_api.ps1

# 再次检查目录大小（应该没有增加）
```

## 📚 相关资源

- [MySQL 日志文件](https://dev.mysql.com/doc/refman/8.0/en/binary-log.html)
- [二进制日志选项](https://dev.mysql.com/doc/refman/8.0/en/replication-options-binary-log.html)
- [通用查询日志](https://dev.mysql.com/doc/refman/8.0/en/query-log.html)
- [慢查询日志](https://dev.mysql.com/doc/refman/8.0/en/slow-query-log.html)

---

**总结**: 已禁用 MySQL 日志生成。建议清理现有日志文件以节省空间。修改后需重启 MySQL 容器。

