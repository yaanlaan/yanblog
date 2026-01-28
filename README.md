# Go+Vue 前后端分离博客系统

[原项目链接](https://github.com/wejectchen/Ginblog.git)

后端主要参考了 ginblog，由于原教程是 22 年的有一些东西不一样了，同时我也希望设计一些更丰富的功能。

## 项目简介

这是一个使用 Go 语言和 Vue.js 构建的前后端分离博客系统。后端采用 Gin 框架，前端使用 Vue 3 + TypeScript + Vite 构建。支持 Docker 一键部署。

## 功能特性

### 后端功能

- 用户管理（注册、登录、权限控制）
- 文章管理（创建、编辑、删除、分类）
- 分类管理
- 文件上传
- 天气信息获取
- 系统状态监控
- JWT Token 认证
- RESTful API 设计

### 前端功能

- **全站暗黑模式适配**：支持跟随系统自动切换或手动切换，精心调教的配色，解决了加载闪烁问题。
- **文章归档与标签云**：新增时间轴样式的归档页面和侧边栏标签云，方便按时间和标签检索文章。
- **精致的 UI/UX**：
  - 加载动画与骨架屏优化，提升视觉体验。
  - 自定义滚动条样式。
  - 侧边栏卡片（天气、热门文章、服务器状态）视觉升级。
- 响应式设计，适配不同设备
- 文章列表展示与分类筛选
- 文章详情页面（支持目录、评论）
- 管理后台（用户、文章、分类管理）
- 搜索功能
- **全配置化管理**：头像、背景图、Logo、二维码等均可通过配置文件修改，无需修改代码。

## 快速开始 (Docker 部署)

这是最推荐的部署方式，简单快捷。

### 1. 准备配置文件

在项目根目录下，你需要准备好配置文件。

将docker_field文件夹下的docker_config.yaml复制重命名为config.yaml，记得填写好相关的信息（如果想自己本地运行，把docker_field文件夹下的docker_config.yaml复制重命名为config.yaml放到/config下，还有一个放到/web/frontend/public/config.yaml）

### 2. 启动服务

确保你已经安装了 Docker 和 Docker Compose，然后在项目根目录下运行：

```bash
docker-compose up --build -d
```

### 3. 访问服务

启动成功后，你可以通过以下地址访问：

- **博客前台**: [http://localhost:3002](http://localhost:3002)
- **后台管理**: [http://localhost:3001](http://localhost:3001)
- **后端 API**: [http://localhost:8080](http://localhost:8080)

### 4. 后续维护

- **修改前端配置**：直接编辑 `docker_field/frontend/config.yaml`，保存后刷新浏览器即可生效。
- **修改关于页面**：直接编辑 `docker_field/frontend/static/about.md`，支持 Markdown 语法。保存后刷新浏览器即可看到更新，无需重启容器。
- **更换图片资源**：将图片放入 `docker_field/frontend/static/` 目录（例如 `avatar.jpg`），然后在配置文件中引用 `/static/avatar.jpg`。
- **修改后端配置**：编辑 `docker_field/backend/config.yaml`，保存后需要重启后端容器：`docker-compose restart backend`。
- **数据备份**：
    - 数据库数据位于 `docker_field/mysql/data`
    - 上传的文件位于 `docker_field/uploads`


## API 文档

详细的 API 接口文档请查看 [apidoc.md](apidoc.md)

## 配置说明

项目使用 YAML 格式进行配置

## 预览

### 前端
===由于是整页捕获的，所以有些下半是白色，但其实正常不会==
![首页](./readme_src/前端预览/1.png)
![文章](./readme_src/前端预览/2.png)
![分类](./readme_src/前端预览/3.jpeg)
![归档](./readme_src/前端预览/4.jpeg)
![关于](./readme_src/前端预览/5.png)
![测试文章](./readme_src/前端预览/6.jpeg)

### 后端

![仪表盘](./readme_src/后端预览/1.jpeg)
![用户列表](./readme_src/后端预览/2.jpeg)
![分类列表](./readme_src/后端预览/3.jpeg)
![文章列表](./readme_src/后端预览/4.jpeg)
![媒体库](./readme_src/后端预览/5.jpeg)
![文章详情](./readme_src/后端预览/6.jpeg)
![文章详情](./readme_src/后端预览/7.jpeg)