# 关于本站

基于 Go + Vue 3 构建的个人博客系统，支持 Markdown 写作、暗黑模式、Docker 一键部署。

## 技术栈

- **后端**: Go, Gin, GORM, JWT
- **前端**: Vue 3, TypeScript, Vite, Pinia
- **数据库**: SQLite / MySQL
- **部署**: Docker, Nginx

## 功能

- Markdown 文章编辑，支持数学公式和流程图
- 文章分类、标签、置顶、归档
- 全站暗黑模式，响应式设计
- 后台可视化配置，无需修改代码
- 文件管理器，支持批量上传

## 代码示例

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    r.GET("/api/v1/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    r.Run(":8080")
}
```

```vue
<script setup lang="ts">
import { ref } from 'vue'
const message = ref('Hello YanBlog!')
</script>

<template>
  <h1>{{ message }}</h1>
</template>
```
