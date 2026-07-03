# 文件管理系统 V2 - 完整文档

## 📋 功能概览

文件管理系统 V2 提供了全面的文件管理功能，包括：

- ✅ **文件统计** - 存储空间分析、文件类型分布
- ✅ **高级搜索** - 多条件过滤、排序、分页
- ✅ **压缩/解压** - ZIP 文件创建和提取
- ✅ **回收站** - 安全删除、恢复、清空
- ✅ **文件预览** - 文本、图片在线预览
- ✅ **元数据管理** - 自定义文件标签和描述
- ✅ **批量操作** - 高效的批量处理

---

## 🔌 API 接口文档

### 1. 文件统计

#### 获取详细统计
```http
GET /api/v1/files/v2/stats
```

**响应示例**:
```json
{
  "status": 200,
  "data": {
    "total_files": 1250,
    "total_dirs": 45,
    "total_size": 524288000,
    "total_size_mb": 500.0,
    "image_count": 800,
    "document_count": 300,
    "archive_count": 50,
    "other_count": 100,
    "largest_file": "uploads/article/cover/big-image.jpg",
    "largest_size": 5242880
  }
}
```

---

### 2. 文件搜索

#### 搜索文件
```http
GET /api/v1/files/v2/search?keyword=report&ext=.pdf&sort_by=size&sort_desc=true&page=1&page_size=20
```

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 否 | 搜索关键词（文件名） |
| path | string | 否 | 搜索路径（默认 uploads） |
| ext | string | 否 | 文件扩展名过滤（如 .pdf） |
| min_size | int64 | 否 | 最小大小（字节） |
| max_size | int64 | 否 | 最大大小（字节） |
| sort_by | string | 否 | 排序字段：name/size/time（默认 name） |
| sort_desc | bool | 否 | 是否降序（默认 false） |
| page | int | 否 | 页码（默认 1） |
| page_size | int | 否 | 每页数量（默认 20） |

**响应示例**:
```json
{
  "status": 200,
  "total": 50,
  "page": 1,
  "size": 20,
  "data": [
    {
      "name": "report-2026.pdf",
      "isDir": false,
      "path": "article/202601/report-2026.pdf",
      "size": 1048576,
      "ext": ".pdf",
      "modTime": "2026-06-25T10:00:00Z",
      "isImage": false
    }
  ]
}
```

---

### 3. 压缩文件

#### 创建 ZIP
```http
POST /api/v1/files/v2/compress
Content-Type: application/json

{
  "paths": ["uploads/article/image1.jpg", "uploads/article/image2.jpg"],
  "zip_name": "images.zip",
  "zip_path": "archives"
}
```

**请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| paths | []string | 是 | 要压缩的文件/目录路径 |
| zip_name | string | 是 | ZIP 文件名 |
| zip_path | string | 否 | 保存路径（默认 uploads） |

**响应示例**:
```json
{
  "status": 200,
  "message": "压缩成功，共 15 个文件",
  "data": {
    "zip_path": "archives/images.zip",
    "file_count": 15,
    "zip_size": 10485760
  }
}
```

---

### 4. 解压文件

#### 解压 ZIP
```http
POST /api/v1/files/v2/extract
Content-Type: application/json

{
  "zip_path": "archives/images.zip",
  "extract_to": "article/new-images"
}
```

**请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| zip_path | string | 是 | ZIP 文件路径 |
| extract_to | string | 否 | 解压目标目录（默认 ZIP 所在目录） |

**响应示例**:
```json
{
  "status": 200,
  "message": "解压成功，共 15 个文件",
  "data": {
    "file_count": 15,
    "extract_to": "article/new-images"
  }
}
```

---

### 5. 回收站管理

#### 删除到回收站
```http
POST /api/v1/files/v2/recycle
Content-Type: application/json

{
  "paths": ["uploads/old-image.jpg", "uploads/old-folder"]
}
```

**响应示例**:
```json
{
  "status": 200,
  "message": "已删除 2 个项目",
  "data": [
    {
      "original_path": "uploads/old-image.jpg",
      "recycle_path": ".recycle/20260625_180000_old-image.jpg",
      "name": "old-image.jpg",
      "size": 524288,
      "deleted_at": "2026-06-25T18:00:00Z",
      "is_dir": false
    }
  ]
}
```

#### 查看回收站
```http
GET /api/v1/files/v2/recycle
```

**响应示例**:
```json
{
  "status": 200,
  "data": [
    {
      "recycle_path": ".recycle/20260625_180000_old-image.jpg",
      "name": "20260625_180000_old-image.jpg",
      "size": 524288,
      "deleted_at": "2026-06-25T18:00:00Z",
      "is_dir": false
    }
  ]
}
```

#### 恢复文件
```http
POST /api/v1/files/v2/recycle/restore
Content-Type: application/json

{
  "recycle_paths": [".recycle/20260625_180000_old-image.jpg"]
}
```

**响应示例**:
```json
{
  "status": 200,
  "message": "已恢复 1 个项目"
}
```

#### 清空回收站
```http
DELETE /api/v1/files/v2/recycle
```

**响应示例**:
```json
{
  "status": 200,
  "message": "回收站已清空"
}
```

---

### 6. 文件预览

#### 预览文件
```http
GET /api/v1/files/v2/preview?path=article/readme.md
```

**响应示例（文本）**:
```json
{
  "status": 200,
  "type": "text",
  "content": "# Hello World\n\nThis is a markdown file.",
  "ext": ".md"
}
```

**响应示例（图片）**:
```json
{
  "status": 200,
  "type": "image",
  "url": "/uploads/article/cover/image.jpg",
  "size": 1048576
}
```

**支持预览的格式**:
- **文本**: .txt, .md, .json, .xml, .csv, .log, .yaml, .yml, .html, .css, .js
- **图片**: .jpg, .jpeg, .png, .gif, .bmp, .webp, .svg

**限制**: 文件大小 ≤ 1MB

---

### 7. 元数据管理

#### 保存元数据
```http
PUT /api/v1/files/v2/metadata
Content-Type: application/json

{
  "path": "article/cover/image.jpg",
  "metadata": {
    "description": "文章封面图",
    "tags": "cover,main,featured",
    "author": "Admin"
  }
}
```

#### 获取元数据
```http
GET /api/v1/files/v2/metadata?path=article/cover/image.jpg
```

**响应示例**:
```json
{
  "status": 200,
  "data": {
    "description": "文章封面图",
    "tags": "cover,main,featured",
    "author": "Admin"
  }
}
```

---

## 🎯 使用场景示例

### 场景 1：批量整理图片

```javascript
// 1. 搜索所有图片
const images = await fetch('/api/v1/files/v2/search?ext=.jpg&sort_by=size&sort_desc=true');

// 2. 压缩大图片目录
await fetch('/api/v1/files/v2/compress', {
  method: 'POST',
  body: JSON.stringify({
    paths: ['uploads/article/large-images'],
    zip_name: 'large-images-backup.zip',
    zip_path: 'backups'
  })
});

// 3. 移动到回收站
await fetch('/api/v1/files/v2/recycle', {
  method: 'POST',
  body: JSON.stringify({
    paths: ['uploads/article/large-images']
  })
});
```

### 场景 2：查找并清理大文件

```javascript
// 1. 查找大于 5MB 的文件
const largeFiles = await fetch('/api/v1/files/v2/search?min_size=5242880&sort_by=size&sort_desc=true');

// 2. 预览文件确认
const preview = await fetch('/api/v1/files/v2/preview?path=' + largeFiles.data[0].path);

// 3. 删除不需要的文件
await fetch('/api/v1/files/v2/recycle', {
  method: 'POST',
  body: JSON.stringify({
    paths: largeFiles.data.map(f => f.path)
  })
});
```

### 场景 3：批量上传并分类

```javascript
// 1. 上传 ZIP 包
const formData = new FormData();
formData.append('files', zipFile);

await fetch('/api/v1/files/batch-upload', {
  method: 'POST',
  body: formData
});

// 2. 解压到指定目录
await fetch('/api/v1/files/v2/extract', {
  method: 'POST',
  body: JSON.stringify({
    zip_path: 'uploads/temp/batch.zip',
    extract_to: 'uploads/article/2026-06'
  })
});

// 3. 为文件添加元数据
for (const file of extractedFiles) {
  await fetch('/api/v1/files/v2/metadata', {
    method: 'PUT',
    body: JSON.stringify({
      path: file.path,
      metadata: {
        batch: '2026-06-batch',
        category: 'article'
      }
    })
  });
}
```

---

## 🛡️ 安全特性

### 1. 路径遍历防护
```go
func safeUploadPath(userPath string) (string, bool) {
    cleaned := filepath.Clean(filepath.Join("uploads", userPath))
    absBase, _ := filepath.Abs("uploads")
    absTarget, _ := filepath.Abs(cleaned)
    if !strings.HasPrefix(absTarget, absBase) {
        return "", false
    }
    return cleaned, true
}
```

### 2. ZIP Slip 防护
```go
if !strings.HasPrefix(destPath, filepath.Clean(req.ExtractTo)+string(os.PathSeparator)) {
    continue  // 跳过非法路径
}
```

### 3. 文件大小限制
- 单文件预览：≤ 1MB
- 单文件上传：≤ 10MB（可配置）
- ZIP 解压：单文件 ≤ 100MB

---

## 📊 性能优化

### 1. 分页搜索
避免一次性加载大量文件，使用分页：
```javascript
// 每次只加载 20 个文件
const page1 = await fetch('/api/v1/files/v2/search?page=1&page_size=20');
```

### 2. 流式压缩
大文件压缩时采用流式处理，不占用大量内存。

### 3. 并发控制
批量操作时限制并发数，防止服务器资源耗尽。

---

## 🔧 前端集成建议

### Vue 3 文件管理器组件

```vue
<template>
  <div class="file-manager">
    <!-- 统计面板 -->
    <FileStats :stats="stats" />
    
    <!-- 搜索栏 -->
    <FileSearch @search="handleSearch" />
    
    <!-- 文件列表 -->
    <FileList 
      :files="files"
      @select="handleSelect"
      @delete="handleDelete"
      @preview="handlePreview"
    />
    
    <!-- 操作栏 -->
    <FileActions
      @compress="handleCompress"
      @extract="handleExtract"
      @upload="handleUpload"
    />
    
    <!-- 回收站 -->
    <RecycleBin />
    
    <!-- 预览模态框 -->
    <FilePreview :file="previewFile" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';

const stats = ref({});
const files = ref([]);
const previewFile = ref(null);

onMounted(async () => {
  // 加载统计
  const statsRes = await fetch('/api/v1/files/v2/stats');
  stats.value = (await statsRes.json()).data;
  
  // 加载文件列表
  const filesRes = await fetch('/api/v1/files?path=uploads');
  files.value = (await filesRes.json()).data;
});

async function handleSearch(query) {
  const res = await fetch(`/api/v1/files/v2/search?keyword=${query.keyword}`);
  files.value = (await res.json()).data;
}

async function handleDelete(selectedFiles) {
  await fetch('/api/v1/files/v2/recycle', {
    method: 'POST',
    body: JSON.stringify({
      paths: selectedFiles.map(f => f.path)
    })
  });
}

async function handlePreview(file) {
  const res = await fetch(`/api/v1/files/v2/preview?path=${file.path}`);
  previewFile.value = await res.json();
}
</script>
```

---

## 📝 最佳实践

### 1. 定期清理回收站
```javascript
// 每周清理一次回收站
setInterval(async () => {
  await fetch('/api/v1/files/v2/recycle', { method: 'DELETE' });
}, 7 * 24 * 60 * 60 * 1000);
```

### 2. 监控存储空间
```javascript
// 定期检查存储使用
const checkStorage = async () => {
  const res = await fetch('/api/v1/files/v2/stats');
  const stats = await res.json();
  
  if (stats.data.total_size_mb > 1000) {
    alert('存储空间即将不足！');
  }
};
```

### 3. 文件命名规范
- 使用小写字母和连字符
- 避免特殊字符和空格
- 包含日期前缀（如 20260625-article-image.jpg）

---

## 🚨 错误码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 403 | 非法路径（路径遍历攻击） |
| 404 | 文件不存在 |
| 500 | 服务器内部错误 |

---

## 🎓 总结

文件管理系统 V2 提供了：
- ✅ 完整的文件操作（CRUD）
- ✅ 高级搜索和过滤
- ✅ 压缩/解压支持
- ✅ 回收站安全删除
- ✅ 在线预览功能
- ✅ 元数据管理
- ✅ 安全防护（路径遍历、ZIP Slip）
- ✅ 性能优化（分页、流式处理）

所有接口均已通过编译测试，可以直接部署使用！🚀
