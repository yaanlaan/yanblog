# YANBLOG
## 项目简介

这是一个使用 Vue 3、Go 语言和现代 Web 技术构建的全栈博客系统。项目采用前后端分离架构，后端使用 Gin 框架提供 RESTful API，前端使用 Vue 3 构建响应式界面。

本项目支持 Docker 一键部署，配置灵活，开箱即用。

## 技术栈

### 后端技术

- **Go 语言**: 高性能的后端开发语言
- **Gin Web 框架**: 轻量级、高性能的 Web 框架
- **GORM**: 强大的 ORM 库，支持多种数据库
- **JWT**: JSON Web Token 用于身份认证
- **MySQL**: 关系型数据库存储

### 前端技术

- **Vue 3**: 使用 Composition API 构建现代化界面
- **TypeScript**: 提供类型安全的代码开发体验
- **Vite**: 极速的前端构建工具
- **Element Plus**: 优雅的 Vue 3 UI 组件库
- **Vue Router**: 官方路由管理器
- **Pinia**: 直观的状态管理库

### 部署与工具

- **Docker & Docker Compose**: 容器化部署，环境隔离，一键启动
- **Nginx**: 高性能反向代理服务器
- **Git**: 版本控制系统

## 功能特性

### 📝 文章管理

完整的文章发布、编辑、删除功能，支持 Markdown 格式编写，实时预览。

### 🏷️ 分类系统

灵活的文章分类管理，便于内容组织和检索。

### 🔐 用户权限

基于 JWT 的用户认证和权限控制系统，保障系统安全。超级管理员权限受到严格保护，防止误操作降权。

### ☁️ 文件上传

支持本地存储及云存储扩展，上传图片大小限制已优化至 50MB。

### ⚙️ 全配置化

前端页面展示信息（如头像、背景、社交链接、二维码等）完全通过配置文件管理，无需修改源码即可定制你的个人博客。

### 🌤️ 天气信息

集成第三方天气 API，实时显示天气状况，增添生活气息。

### 📊 系统监控

实时监控服务器状态，包括内存、CPU 使用率等，掌握系统运行情况。


## 视频示例
<iframe src="//player.bilibili.com/player.html?isOutside=true&aid=626357031&bvid=BV1yt4y1Q7SS&cid=210738676&p=1&autoplay=0" scrolling="no" border="0" frameborder="no" framespacing="0" allowfullscreen="true"></iframe>