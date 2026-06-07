# 🚀 快速开始指南

> 本指南专为没有编程基础的用户设计，按照步骤操作即可部署你的博客系统。

## 📋 目录

- [方式一：Docker 部署（推荐，最简单）](#方式一 docker 部署推荐最简单)
- [方式二：本地开发部署](#方式二本地开发部署)
- [配置说明](#配置说明)
- [常见问题](#常见问题)

---

## 方式一：Docker 部署（推荐，最简单）

### 前置要求

- 安装 [Docker Desktop](https://www.docker.com/products/docker-desktop/)（Windows/Mac）或 Docker（Linux）
- 确保 Docker 正在运行

### 步骤

#### 1️⃣ 克隆项目

```bash
git clone <你的项目地址>
cd yanblog
```

#### 2️⃣ 修改配置文件

打开 `config/config.yaml` 文件，修改以下内容：

```yaml
server:
  AppMode: production  # 生产环境设为 production
  HttpPort: :8080      # 服务端口，一般不需要改
  SiteUrl: https://your-domain.com  # 你的网站域名（可选）

database:
  Db: MYSQL
  DbHost: db           # Docker 环境下固定为 db
  DbPort: 3306
  DbUser: root
  DbPassWord: your_password  # 设置你的数据库密码（重要！）
  DbName: yanblog

# JWT 密钥（重要！）
# 在命令行运行：openssl rand -hex 32
# 将生成的字符串填到下面
JwtKey: 将这里替换成随机字符串

weather:
  Provider: openweathermap
  ApiKey: 你的天气 API（可选，不需要可以留空）
  DefaultCity: Beijing  # 默认显示的城市
```

> ⚠️ **必须修改的项**：
> - `DbPassWord`：设置一个强密码
> - `JwtKey`：生成一个随机字符串

#### 3️ 一键启动

**Windows 用户**：
```bash
.\docker.ps1
```

**Mac/Linux 用户**：
```bash
./docker.sh
```

#### 4️⃣ 访问博客

打开浏览器访问：
- **前台**：http://localhost:3002
- **后台**：http://localhost:3001

#### 5️⃣ 首次使用

1. 访问后台管理页面
2. 首次访问会自动跳转到注册页面
3. 注册管理员账号
4. 登录后台开始管理博客

---

## 方式二：本地开发部署

### 前置要求

- **后端**：
  - Go 1.21+
  - MySQL 8.0+
  
- **前端**：
  - Node.js 18+
  - npm 或 pnpm

### 步骤

#### 1️⃣ 准备数据库

1. 安装 MySQL
2. 创建数据库：

```sql
CREATE DATABASE yanblog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### 2️⃣ 配置后端

编辑 `config/config.yaml`：

```yaml
server:
  AppMode: debug  # 开发环境设为 debug
  HttpPort: :8080
  SiteUrl: 

database:
  Db: MYSQL
  DbHost: localhost  # 本地数据库
  DbPort: 3306
  DbUser: root
  DbPassWord: 你的数据库密码
  DbName: yanblog

JwtKey: 随机字符串

weather:
  Provider: openweathermap
  ApiKey: 
  DefaultCity: Beijing
```

#### 3️⃣ 启动后端

```bash
go run main.go
```

看到以下日志表示成功：
```
Starting Gin server on :8080
```

#### 4️⃣ 配置前端

编辑 `web/frontend/public/config.yaml`：

```yaml
# 后端 API 地址
apiUrl: http://localhost:8080/api/v1

# 后台管理地址
adminUrl: http://localhost:3001
```

#### 5️ 启动前台

```bash
cd web/frontend
npm install
npm run dev
```

#### 6️⃣ 启动后台

```bash
cd web/backend
npm install
npm run dev
```

#### 7️⃣ 访问

- **前台**：http://localhost:5173
- **后台**：http://localhost:5174

---

## 📖 配置说明

### 后端配置 (`config/config.yaml`)

#### 服务器配置

| 配置项 | 说明 | 示例 | 是否必须 |
|--------|------|------|----------|
| `AppMode` | 运行模式 | `debug` / `production` | 是 |
| `HttpPort` | 服务端口 | `:8080` | 是 |
| `SiteUrl` | 网站域名 | `https://blog.example.com` | 否 |

#### 数据库配置

| 配置项 | 说明 | 示例 | 是否必须 |
|--------|------|------|----------|
| `Db` | 数据库类型 | `MYSQL` | 是 |
| `DbHost` | 数据库地址 | `localhost` 或 `db` | 是 |
| `DbPort` | 数据库端口 | `3306` | 是 |
| `DbUser` | 数据库用户名 | `root` | 是 |
| `DbPassWord` | 数据库密码 | 你的密码 | **必须** |
| `DbName` | 数据库名称 | `yanblog` | 是 |

#### JWT 配置

| 配置项 | 说明 | 生成方式 | 是否必须 |
|--------|------|----------|----------|
| `JwtKey` | JWT 密钥 | `openssl rand -hex 32` | **必须** |

>  **如何生成 JwtKey**：
> - Windows（需要安装 Git Bash）：`openssl rand -hex 32`
> - Mac/Linux：`openssl rand -hex 32`
> - 在线生成：访问 https://randomkeygen.com/ 复制一个密钥

#### 天气配置（可选）

| 配置项 | 说明 | 示例 | 是否必须 |
|--------|------|------|----------|
| `Provider` | 天气提供商 | `openweathermap` | 否 |
| `ApiKey` | API 密钥 | 你的密钥 | 否 |
| `DefaultCity` | 默认城市 | `Beijing` | 否 |

> 🌤️ **获取天气 API 密钥**：
> 1. 访问 https://openweathermap.org/
> 2. 注册账号
> 3. 在 API keys 页面生成密钥
> 4. 复制到配置文件

### 前端配置 (`web/frontend/public/config.yaml`)

```yaml
# 后端 API 地址
# Docker 部署：http://localhost:8080/api/v1
# 本地部署：http://localhost:8080/api/v1
apiUrl: http://localhost:8080/api/v1

# 后台管理地址
# 一般不需要修改
adminUrl: http://localhost:3001
```

---

##  常见问题

### Q1: 启动时提示 "database connection failed"

**原因**：数据库配置错误或数据库未启动

**解决方法**：
1. 检查 `config/config.yaml` 中的数据库配置
2. 确保 MySQL 正在运行
3. Docker 用户确保 `docker-compose.yaml` 中的数据库服务正常

### Q2: 前台无法访问后端 API

**原因**：跨域问题或 API 地址配置错误

**解决方法**：
1. 检查 `web/frontend/public/config.yaml` 中的 `apiUrl`
2. 确保后端服务正在运行
3. 检查端口是否正确

### Q3: 登录时提示 "JWT key not set"

**原因**：未配置 JWT 密钥

**解决方法**：
1. 生成 JWT 密钥：`openssl rand -hex 32`
2. 将生成的字符串填入 `config/config.yaml` 的 `JwtKey` 字段
3. 重启服务

### Q4: 文件上传失败

**原因**：上传目录不存在或权限不足

**解决方法**：
```bash
# 创建上传目录
mkdir -p uploads/images
mkdir -p uploads/files

# 设置权限（Linux/Mac）
chmod -R 755 uploads
```

### Q5: Docker 启动失败

**原因**：端口被占用或 Docker 配置问题

**解决方法**：
1. 检查端口是否被占用：
   ```bash
   # Windows PowerShell
   netstat -ano | findstr :8080
   
   # Mac/Linux
   lsof -i :8080
   ```
2. 修改 `docker-compose.yaml` 中的端口映射
3. 重启 Docker Desktop

### Q6: 如何修改网站标题和 Logo？

**方法**：
1. 访问后台管理系统
2. 进入"配置管理"页面
3. 修改以下配置：
   - 网站标题
   - Logo URL
   - 网站图标（favicon）
   - 头像
   - 背景图片
4. 保存配置

### Q7: 如何自定义域名？

**步骤**：
1. 购买域名
2. 配置 DNS 解析到你的服务器 IP
3. 修改 `config/config.yaml` 中的 `SiteUrl`
4. 配置 Nginx 反向代理（可选）
5. 申请 SSL 证书（推荐）

---

## 🎯 下一步

配置完成后，你可以：

1. ️ **发布文章**：在后台管理创建你的第一篇文章
2. 🎨 **自定义外观**：在配置管理修改网站样式
3. 📊 **查看数据**：在仪表盘查看访问统计
4. 🏷️ **管理分类标签**：组织你的内容

---

## 📞 需要帮助？

如果遇到问题：

1. 查看项目 [README.md](README.md)
2. 检查 [常见问题](#常见问题)
3. 查看项目 [Issues](https://github.com/your-repo/yanblog/issues)
4. 提交 Issue 描述你的问题

---

**祝你使用愉快！** 🎉
