# API 文档

## 基础信息

- API 前缀: `/api/v1`
- 数据格式: JSON
- 字符编码: UTF-8

## 认证方式

部分接口需要认证，认证方式为 JWT Token：
- 请求头添加 `Authorization: Bearer <token>`

## 用户相关接口

### 1. 用户注册
- 请求URL: `POST /api/v1/user/add`
- 请求参数:
  ```json
  {
    "username": "用户名",
    "password": "密码"
  }
  ```
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "创建成功",
    "data": null
  }
  ```

### 2. 用户登录
- 请求URL: `POST /api/v1/login`
- 请求参数:
  ```json
  {
    "username": "用户名",
    "password": "密码"
  }
  ```
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "登录成功",
    "data": {
      "token": "JWT_TOKEN"
    }
  }
  ```

### 3. 获取用户列表
- 请求URL: `GET /api/v1/users`
- 请求参数:
  - `pagesize`: 每页条数
  - `pagenum`: 页码
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "OK",
    "data": {
      "total": 100,
      "data": [
        {
          "ID": 1,
          "CreatedAt": "2022-01-01T00:00:00Z",
          "UpdatedAt": "2022-01-01T00:00:00Z",
          "username": "用户名",
          "role": 1
        }
      ]
    }
  }
  ```

### 4. 搜索用户
- 请求URL: `GET /api/v1/users/search`
- 请求参数:
  - `pagesize`: 每页条数
  - `pagenum`: 页码
  - `keyword`: 搜索关键词
- 响应格式同获取用户列表

### 5. 编辑用户
- 请求URL: `PUT /api/v1/user/:id`
- 请求参数:
  ```json
  {
    "username": "用户名",
    "password": "密码",
    "role": 1
  }
  ```
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "修改成功",
    "data": null
  }
  ```

### 6. 删除用户
- 请求URL: `DELETE /api/v1/user/:id`
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "删除成功",
    "data": null
  }
  ```

## 分类相关接口

### 1. 创建分类
- 请求URL: `POST /api/v1/category/add`
- 请求参数:
  ```json
  {
    "name": "分类名称"
  }
  ```
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "创建成功",
    "data": null
  }
  ```

### 2. 获取分类列表
- 请求URL: `GET /api/v1/category`
- 请求参数:
  - `pagesize`: 每页条数
  - `pagenum`: 页码
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "OK",
    "data": {
      "total": 100,
      "data": [
        {
          "ID": 1,
          "CreatedAt": "2022-01-01T00:00:00Z",
          "UpdatedAt": "2022-01-01T00:00:00Z",
          "name": "分类名称",
          "article_count": 10
        }
      ]
    }
  }
  ```

### 3. 搜索分类
- 请求URL: `GET /api/v1/category/search`
- 请求参数:
  - `pagesize`: 每页条数
  - `pagenum`: 页码
  - `keyword`: 搜索关键词
- 响应格式同获取分类列表

### 4. 获取分类信息
- 请求URL: `GET /api/v1/category/info/:id`
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "OK",
    "data": {
      "ID": 1,
      "CreatedAt": "2022-01-01T00:00:00Z",
      "UpdatedAt": "2022-01-01T00:00:00Z",
      "name": "分类名称"
    }
  }
  ```

### 5. 编辑分类
- 请求URL: `PUT /api/v1/category/:id`
- 请求参数:
  ```json
  {
    "name": "分类名称"
  }
  ```
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "修改成功",
    "data": null
  }
  ```

### 6. 删除分类
- 请求URL: `DELETE /api/v1/category/:id`
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "删除成功",
    "data": null
  }
  ```

## 文章相关接口

### 1. 创建文章
- 请求URL: `POST /api/v1/article/add`
- 请求参数:
  ```json
  {
    "title": "文章标题",
    "cid": 1,
    "desc": "文章描述",
    "content": "文章内容",
    "img": "图片URL"
  }
  ```
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "创建成功",
    "data": null
  }
  ```

### 2. 获取文章列表
- 请求URL: `GET /api/v1/article`
- 请求参数:
  - `pagesize`: 每页条数
  - `pagenum`: 页码
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "OK",
    "data": {
      "total": 100,
      "data": [
        {
          "ID": 1,
          "CreatedAt": "2022-01-01T00:00:00Z",
          "UpdatedAt": "2022-01-01T00:00:00Z",
          "title": "文章标题",
          "cid": 1,
          "desc": "文章描述",
          "content": "文章内容",
          "img": "图片URL",
          "Category": {
            "name": "分类名称"
          }
        }
      ]
    }
  }
  ```

### 3. 搜索文章
- 请求URL: `GET /api/v1/article/search`
- 请求参数:
  - `pagesize`: 每页条数
  - `pagenum`: 页码
  - `keyword`: 搜索关键词 (可选)
  - `cid`: 分类ID (可选)
- 响应格式同获取文章列表

### 4. 获取分类下文章
- 请求URL: `GET /api/v1/article/list/:id`
- 请求参数:
  - `pagesize`: 每页条数
  - `pagenum`: 页码
- 响应格式同获取文章列表

### 5. 获取文章详情
- 请求URL: `GET /api/v1/article/info/:id`
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "OK",
    "data": {
      "ID": 1,
      "CreatedAt": "2022-01-01T00:00:00Z",
      "UpdatedAt": "2022-01-01T00:00:00Z",
      "title": "文章标题",
      "cid": 1,
      "desc": "文章描述",
      "content": "文章内容",
      "img": "图片URL"
    }
  }
  ```

### 6. 获取置顶文章
- 请求URL: `GET /api/v1/article/top`
- 请求参数:
  - `num`: 获取文章数量 (可选，默认3)
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "OK",
    "data": [
      {
        "ID": 1,
        "CreatedAt": "2022-01-01T00:00:00Z",
        "UpdatedAt": "2022-01-01T00:00:00Z",
        "title": "文章标题",
        "cid": 1,
        "desc": "文章描述",
        "content": "文章内容",
        "img": "图片URL"
      }
    ]
  }
  ```

### 7. 编辑文章
- 请求URL: `PUT /api/v1/article/:id`
- 请求参数:
  ```json
  {
    "title": "文章标题",
    "cid": 1,
    "desc": "文章描述",
    "content": "文章内容",
    "img": "图片URL"
  }
  ```
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "修改成功",
    "data": null
  }
  ```

### 8. 删除文章
- 请求URL: `DELETE /api/v1/article/:id`
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "删除成功",
    "data": null
  }
  ```

## 其他接口

### 1. 上传文件
- 请求URL: `POST /api/v1/upload`
- 请求方式: multipart/form-data
- 请求参数:
  - `file`: 文件
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "上传成功",
    "data": "文件URL"
  }
  ```

### 2. 获取天气信息
- 请求URL: `GET /api/v1/weather`
- 请求参数:
  - `city`: 城市名称 (可选，默认合肥)
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "OK",
    "data": {
      "city": "Hefei",
      "temperature": 25.5,
      "description": "多云",
      "humidity": 65,
      "wind_speed": 2.1
    }
  }
  ```

### 3. 获取系统状态
- 请求URL: `GET /api/v1/system/status`
- 响应示例:
  ```json
  {
    "status": 200,
    "message": "OK",
    "data": {
      "status": "online",
      "uptime": "1天2小时30分钟",
      "memory_usage": 35.5,
      "cpu_usage": 12.3,
      "goroutines": 15,
      "start_time": 1609459200000
    }
  }
  ```