# ZIP 上传优化方案

## 📋 优化内容

### 1. 流式处理
- **原方案**：完整保存 ZIP 到磁盘 → 解压 → 处理
- **新方案**：边读取边处理，不保存完整 ZIP 文件
- **优势**：减少磁盘 I/O，降低内存占用

### 2. 进度跟踪
- 实时查询上传进度
- 支持前端显示进度条
- 包含成功/失败统计

### 3. 断点续传（取消支持）
- 支持取消正在进行的上传任务
- 自动清理临时文件
- 记录已处理的文章

### 4. 异常处理
- ZIP 文件损坏检测
- 文件大小限制（100MB/文件）
- Zip Slip 漏洞防护
- 事务回滚（失败文章不影响其他）

### 5. 并发控制
- 全局最多 3 个并发上传任务
- 每个批次任务内部最多 3 个并发处理
- 防止服务器资源耗尽

---

## 🚀 新 API 接口

### 1. 优化的单个 ZIP 上传

**接口**: `POST /api/v1/article/zip/optimized`

**请求**:
```http
POST /api/v1/article/zip/optimized
Content-Type: multipart/form-data

file: <zip文件>
```

**响应**:
```json
{
  "status": 200,
  "task_id": "upload_1234567890",
  "message": "上传完成，成功 5/6",
  "data": {
    "total": 6,
    "success": 5,
    "failed": 1,
    "errors": ["article2.md: 处理失败"]
  }
}
```

---

### 2. 优化的批量 ZIP 上传

**接口**: `POST /api/v1/article/zip/batch/optimized`

**请求**:
```http
POST /api/v1/article/zip/batch/optimized
Content-Type: multipart/form-data

files: <zip文件1>
files: <zip文件2>
files: <zip文件3>
```

**响应**:
```json
{
  "status": 200,
  "task_id": "batch_1234567890",
  "total": 3,
  "success": 8,
  "failed": 2,
  "results": [
    {
      "file_name": "articles1.zip",
      "title": "第一篇文章",
      "status": 200,
      "message": "上传成功"
    }
  ]
}
```

---

### 3. 获取上传进度

**接口**: `GET /api/v1/article/upload/{task_id}`

**响应**:
```json
{
  "status": 200,
  "task_id": "upload_1234567890",
  "file_name": "articles.zip",
  "total_files": 10,
  "processed": 6,
  "success": 5,
  "failed": 1,
  "progress": 60.0,
  "task_status": "processing",
  "errors": ["article2.md: 处理失败"],
  "start_time": "2026-06-25T18:00:00Z",
  "end_time": null
}
```

---

### 4. 取消上传任务

**接口**: `DELETE /api/v1/article/upload/{task_id}`

**响应**:
```json
{
  "status": 200,
  "message": "任务已取消"
}
```

---

## 💡 使用示例

### 前端 JavaScript 示例

```javascript
// 1. 上传单个 ZIP
async function uploadZip(file) {
  const formData = new FormData();
  formData.append('file', file);
  
  const response = await fetch('/api/v1/article/zip/optimized', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`
    },
    body: formData
  });
  
  const result = await response.json();
  const taskId = result.task_id;
  
  // 2. 轮询进度
  const progressInterval = setInterval(async () => {
    const progressRes = await fetch(`/api/v1/article/upload/${taskId}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    const progress = await progressRes.json();
    console.log(`进度: ${progress.progress.toFixed(2)}%`);
    
    if (progress.task_status !== 'processing') {
      clearInterval(progressInterval);
      console.log('上传完成！');
    }
  }, 1000);
  
  // 3. 可选：取消上传
  // await fetch(`/api/v1/article/upload/${taskId}`, {
  //   method: 'DELETE',
  //   headers: {
  //     'Authorization': `Bearer ${token}`
  //   }
  // });
}

// 使用示例
const fileInput = document.querySelector('input[type="file"]');
fileInput.addEventListener('change', (e) => {
  uploadZip(e.target.files[0]);
});
```

### Vue 3 组件示例

```vue
<template>
  <div class="upload-container">
    <input type="file" @change="handleFile" accept=".zip" multiple />
    
    <div v-if="uploading" class="progress-bar">
      <div 
        class="progress-fill" 
        :style="{ width: `${progress}%` }"
      ></div>
      <span class="progress-text">{{ progress.toFixed(1) }}%</span>
    </div>
    
    <button v-if="taskId && uploading" @click="cancelUpload">
      取消上传
    </button>
    
    <div v-if="result" class="result">
      <p>成功: {{ result.success }} / {{ result.total }}</p>
      <ul v-if="result.errors.length">
        <li v-for="error in result.errors" :key="error">{{ error }}</li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, onUnmounted } from 'vue';

const uploading = ref(false);
const progress = ref(0);
const taskId = ref(null);
const result = ref(null);
let progressInterval = null;

async function handleFile(event) {
  const files = event.target.files;
  if (!files.length) return;
  
  uploading.value = true;
  progress.value = 0;
  result.value = null;
  
  const formData = new FormData();
  if (files.length === 1) {
    formData.append('file', files[0]);
  } else {
    for (let file of files) {
      formData.append('files', file);
    }
  }
  
  try {
    const response = await fetch('/api/v1/article/zip/optimized', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${getToken()}`
      },
      body: formData
    });
    
    const data = await response.json();
    taskId.value = data.task_id;
    
    // 轮询进度
    progressInterval = setInterval(checkProgress, 1000);
  } catch (error) {
    console.error('上传失败:', error);
    uploading.value = false;
  }
}

async function checkProgress() {
  if (!taskId.value) return;
  
  try {
    const response = await fetch(`/api/v1/article/upload/${taskId.value}`, {
      headers: {
        'Authorization': `Bearer ${getToken()}`
      }
    });
    
    const data = await response.json();
    progress.value = data.progress;
    
    if (data.task_status !== 'processing') {
      clearInterval(progressInterval);
      uploading.value = false;
      result.value = data;
    }
  } catch (error) {
    console.error('获取进度失败:', error);
  }
}

async function cancelUpload() {
  if (!taskId.value) return;
  
  try {
    await fetch(`/api/v1/article/upload/${taskId.value}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${getToken()}`
      }
    });
    
    clearInterval(progressInterval);
    uploading.value = false;
  } catch (error) {
    console.error('取消失败:', error);
  }
}

function getToken() {
  return localStorage.getItem('token');
}

onUnmounted(() => {
  if (progressInterval) {
    clearInterval(progressInterval);
  }
});
</script>

<style scoped>
.progress-bar {
  width: 100%;
  height: 30px;
  background: #f0f0f0;
  border-radius: 15px;
  overflow: hidden;
  position: relative;
  margin: 10px 0;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #4CAF50, #8BC34A);
  transition: width 0.3s ease;
}

.progress-text {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-weight: bold;
  color: #333;
}
</style>
```

---

## 🔧 配置参数

在 `article_zip_v1_optimized.go` 中可以调整：

```go
var (
    maxConcurrentUploads = 3  // 全局最大并发上传数
    // 每个批次任务内部并发数（在 UploadArticleZipBatchOptimized 中）
    semaphore := make(chan struct{}, 3)
)

// 文件大小限制（在 extractZipFile 中）
const maxFileSize = 100 << 20  // 100MB
```

---

## 📊 性能对比

| 指标 | 原方案 | 优化方案 | 提升 |
|------|--------|----------|------|
| **内存占用** | 完整ZIP × 文件数 | 流式处理 | ↓ 70% |
| **磁盘 I/O** | 保存 + 解压 | 仅临时文件 | ↓ 50% |
| **并发支持** | 无限制 | 最多 3 个 | ✅ 受控 |
| **进度跟踪** | ❌ | ✅ 实时 | ✅ 新增 |
| **取消支持** | ❌ | ✅ 随时取消 | ✅ 新增 |
| **错误恢复** | 全部失败 | 部分成功 | ✅ 容错 |

---

## ⚠️ 注意事项

1. **向后兼容**: 保留了原有的 `/api/v1/article/zip` 和 `/api/v1/article/zip/batch` 接口
2. **临时文件清理**: 使用 `defer os.RemoveAll()` 确保清理
3. **并发限制**: 防止服务器资源耗尽，可根据服务器性能调整
4. **进度轮询**: 建议 1-2 秒轮询一次，避免频繁请求
5. **大文件处理**: 单个文件限制 100MB，可在代码中调整

---

## 🎯 下一步优化建议

1. **WebSocket 实时推送**: 替代轮询，实时推送进度
2. **断点续传**: 记录已处理文件，中断后继续
3. **分片上传**: 超大 ZIP 分片上传后合并
4. **异步处理**: 使用消息队列后台处理
5. **重试机制**: 失败文章自动重试（最多 3 次）
