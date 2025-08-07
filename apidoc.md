# 博客系统 API 文档

## 1. 概述

这是一个使用 Go 语言和 Gin 框架构建的博客系统后端 API。API 遵循 RESTful 设计原则，提供用户管理、文章管理、分类管理等功能。

### 1.1 基础 URL

所有 API 接口的基础 URL 为: `http://localhost:3000/api/v1`

### 1.2 响应格式

所有 API 接口返回统一的 JSON 格式:

```json
{
  "status": 200,        // 状态码
  "data": {},           // 返回数据
  "message": "成功",     // 状态消息
  "total": 100          // 总数（分页接口）
}
```

### 1.3 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 操作成功 |
| 400 | 请求参数错误 |
| 401 | 未授权访问 |
| 500 | 服务器内部错误 |

### 1.4 认证机制

部分接口需要通过 JWT Token 进行认证。用户登录成功后，系统会返回一个 token，访问需要认证的接口时需要在请求头中添加：

```
Authorization: Bearer <token>
```

---

## 2. 用户管理接口

### 2.1 用户注册

**接口地址**: `POST /user/add`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |
| role | int | 否 | 用户角色(1: 管理员, 2: 普通用户) |

**请求示例**:

```json
{
  "username": "testuser",
  "password": "123456",
  "role": 2
}
```

**响应示例**:

```json
{
  "status": 200,
  "message": "用户创建成功"
}
```

### 2.2 用户登录

**接口地址**: `POST /login`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**请求示例**:

```json
{
  "username": "testuser",
  "password": "123456"
}
```

**响应示例**:

```json
{
  "status": 200,
  "message": "登录成功",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 2.3 获取用户列表

**接口地址**: `GET /users`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| pagesize | int | 否 | 每页条数(-1表示获取所有) |
| pagenum | int | 否 | 页码(-1表示获取所有) |

**响应示例**:

```json
{
  "status": 200,
  "data": [
    {
      "id": 1,
      "username": "admin",
      "role": 1
    },
    {
      "id": 2,
      "username": "testuser",
      "role": 2
    }
  ],
  "total": 2,
  "message": "OK"
}
```

### 2.4 搜索用户

**接口地址**: `GET /users/search`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| pagesize | int | 否 | 每页条数(-1表示获取所有) |
| pagenum | int | 否 | 页码(-1表示获取所有) |
| keyword | string | 否 | 搜索关键词 |
| role | int | 否 | 用户角色 |

**响应示例**:

```json
{
  "status": 200,
  "data": [
    {
      "id": 2,
      "username": "testuser",
      "role": 2
    }
  ],
  "total": 1,
  "message": "OK"
}
```

### 2.5 编辑用户（需要认证）

**接口地址**: `PUT /user/:id`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |
| role | int | 否 | 用户角色 |

**请求示例**:

```json
{
  "username": "newusername",
  "password": "newpassword",
  "role": 2
}
```

**响应示例**:

```json
{
  "status": 200,
  "message": "用户修改成功"
}
```

### 2.6 删除用户（需要认证）

**接口地址**: `DELETE /user/:id`

**响应示例**:

```json
{
  "status": 200,
  "message": "用户删除成功"
}
```

---

## 3. 分类管理接口

### 3.1 添加分类（需要认证）

**接口地址**: `POST /category/add`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 分类名称 |
| img | string | 否 | 分类图片URL |
| top | int | 否 | 是否置顶(0: 不置顶, 1: 置顶) |

**请求示例**:

```json
{
  "name": "技术分享",
  "img": "https://example.com/image.jpg",
  "top": 0
}
```

**响应示例**:

```json
{
  "status": 200,
  "data": {
    "id": 1,
    "name": "技术分享",
    "img": "https://example.com/image.jpg",
    "top": 0
  },
  "message": "分类创建成功"
}
```

### 3.2 获取分类列表

**接口地址**: `GET /category`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| pagesize | int | 否 | 每页条数(-1表示获取所有) |
| pagenum | int | 否 | 页码(-1表示获取所有) |

**响应示例**:

```json
{
  "status": 200,
  "data": [
    {
      "id": 1,
      "name": "技术分享",
      "img": "https://example.com/image.jpg",
      "top": 0
    }
  ],
  "total": 1,
  "message": "OK"
}
```

### 3.3 搜索分类

**接口地址**: `GET /category/search`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| pagesize | int | 否 | 每页条数(-1表示获取所有) |
| pagenum | int | 否 | 页码(-1表示获取所有) |
| keyword | string | 否 | 搜索关键词 |

**响应示例**:

```json
{
  "status": 200,
  "data": [
    {
      "id": 1,
      "name": "技术分享",
      "img": "https://example.com/image.jpg",
      "top": 0
    }
  ],
  "total": 1,
  "message": "OK"
}
```

### 3.4 获取分类信息

**接口地址**: `GET /category/info/:id`

**响应示例**:

```json
{
  "status": 200,
  "data": {
    "id": 1,
    "name": "技术分享",
    "img": "https://example.com/image.jpg",
    "top": 0
  },
  "message": "OK"
}
```

### 3.5 编辑分类（需要认证）

**接口地址**: `PUT /category/:id`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 分类名称 |
| img | string | 否 | 分类图片URL |
| top | int | 否 | 是否置顶(0: 不置顶, 1: 置顶) |

**请求示例**:

```json
{
  "name": "技术文章",
  "img": "https://example.com/new-image.jpg",
  "top": 1
}
```

**响应示例**:

```json
{
  "status": 200,
  "message": "分类修改成功"
}
```

### 3.6 删除分类（需要认证）

**接口地址**: `DELETE /category/:id`

**响应示例**:

```json
{
  "status": 200,
  "message": "分类删除成功"
}
```

---

## 4. 文章管理接口

### 4.1 添加文章（需要认证）

**接口地址**: `POST /article/add`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| title | string | 是 | 文章标题 |
| content | string | 是 | 文章内容 |
| img | string | 否 | 文章图片URL |
| desc | string | 否 | 文章摘要 |
| cid | int | 是 | 分类ID |
| top | int | 否 | 是否置顶(0: 不置顶, 1: 置顶) |

**请求示例**:

```json
{
  "title": "Go语言学习笔记",
  "content": "Go语言是一门现代编程语言...",
  "img": "https://example.com/article-image.jpg",
  "desc": "Go语言学习心得分享",
  "cid": 1,
  "top": 0
}
```

**响应示例**:

```json
{
  "status": 200,
  "data": {
    "id": 1,
    "title": "Go语言学习笔记",
    "content": "Go语言是一门现代编程语言...",
    "img": "https://example.com/article-image.jpg",
    "desc": "Go语言学习心得分享",
    "cid": 1,
    "top": 0
  },
  "message": "文章创建成功"
}
```

### 4.2 获取文章列表

**接口地址**: `GET /article`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| pagesize | int | 否 | 每页条数(-1表示获取所有) |
| pagenum | int | 否 | 页码(-1表示获取所有) |

**响应示例**:

```json
{
  "status": 200,
  "data": [
    {
      "id": 1,
      "title": "Go语言学习笔记",
      "content": "Go语言是一门现代编程语言...",
      "img": "https://example.com/article-image.jpg",
      "desc": "Go语言学习心得分享",
      "cid": 1,
      "top": 0
    }
  ],
  "total": 1,
  "message": "OK"
}
```

### 4.3 搜索文章

**接口地址**: `GET /article/search`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| pagesize | int | 否 | 每页条数(-1表示获取所有) |
| pagenum | int | 否 | 页码(-1表示获取所有) |
| keyword | string | 否 | 搜索关键词 |
| cid | int | 否 | 分类ID |

**响应示例**:

```json
{
  "status": 200,
  "data": [
    {
      "id": 1,
      "title": "Go语言学习笔记",
      "content": "Go语言是一门现代编程语言...",
      "img": "https://example.com/article-image.jpg",
      "desc": "Go语言学习心得分享",
      "cid": 1,
      "top": 0
    }
  ],
  "total": 1,
  "message": "OK"
}
```

### 4.4 获取置顶文章

**接口地址**: `GET /article/top`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| num | int | 否 | 获取文章数量(默认6篇) |

**响应示例**:

```json
{
  "status": 200,
  "data": [
    {
      "id": 1,
      "title": "Go语言学习笔记",
      "content": "Go语言是一门现代编程语言...",
      "img": "https://example.com/article-image.jpg",
      "desc": "Go语言学习心得分享",
      "cid": 1,
      "top": 1
    }
  ],
  "message": "OK"
}
```

### 4.5 获取分类下的文章

**接口地址**: `GET /article/list/:id`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| pagesize | int | 否 | 每页条数(-1表示获取所有) |
| pagenum | int | 否 | 页码(-1表示获取所有) |

**响应示例**:

```json
{
  "status": 200,
  "data": [
    {
      "id": 1,
      "title": "Go语言学习笔记",
      "content": "Go语言是一门现代编程语言...",
      "img": "https://example.com/article-image.jpg",
      "desc": "Go语言学习心得分享",
      "cid": 1,
      "top": 0
    }
  ],
  "total": 1,
  "message": "OK"
}
```

### 4.6 获取文章详情

**接口地址**: `GET /article/info/:id`

**响应示例**:

```json
{
  "status": 200,
  "data": {
    "id": 1,
    "title": "Go语言学习笔记",
    "content": "Go语言是一门现代编程语言...",
    "img": "https://example.com/article-image.jpg",
    "desc": "Go语言学习心得分享",
    "cid": 1,
    "top": 0
  },
  "message": "OK"
}
```

### 4.7 编辑文章（需要认证）

**接口地址**: `PUT /article/:id`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| title | string | 是 | 文章标题 |
| content | string | 是 | 文章内容 |
| img | string | 否 | 文章图片URL |
| desc | string | 否 | 文章摘要 |
| cid | int | 是 | 分类ID |
| top | int | 否 | 是否置顶(0: 不置顶, 1: 置顶) |

**请求示例**:

```json
{
  "title": "Go语言学习笔记（更新版）",
  "content": "Go语言是一门现代编程语言，具有并发性强等特点...",
  "img": "https://example.com/article-image.jpg",
  "desc": "Go语言学习心得分享（更新版）",
  "cid": 1,
  "top": 1
}
```

**响应示例**:

```json
{
  "status": 200,
  "message": "文章修改成功"
}
```

### 4.8 删除文章（需要认证）

**接口地址**: `DELETE /article/:id`

**响应示例**:

```json
{
  "status": 200,
  "message": "文章删除成功"
}
```

---

## 5. 系统功能接口

### 5.1 文件上传（需要认证）

**接口地址**: `POST /upload`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| file | file | 是 | 上传的文件 |

**响应示例**:

```json
{
  "status": 200,
  "message": "文件上传成功",
  "url": "https://example.com/uploaded-file.jpg"
}
```

### 5.2 获取天气信息

**接口地址**: `GET /weather`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| city | string | 否 | 城市名称(默认为配置中的城市) |

**响应示例**:

```json
{
  "status": 200,
  "data": {
    "city": "北京",
    "temperature": 25,
    "description": "晴",
    "humidity": 60,
    "wind_speed": 3.5
  },
  "message": "OK"
}
```

### 5.3 获取系统状态

**接口地址**: `GET /system/status`

**响应示例**:

```json
{
  "status": 200,
  "message": "OK",
  "data": {
    "status": "online",
    "uptime": "1天2小时30分钟15秒",
    "memory_usage": 45.5,
    "cpu_usage": 12.3,
    "goroutines": 25,
    "start_time": 1623456789000
  }
}
```