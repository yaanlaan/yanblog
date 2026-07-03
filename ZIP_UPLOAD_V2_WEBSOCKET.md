# V2 增强版 ZIP 上传 - WebSocket 实时推送

## 🚀 新特性

### 1. WebSocket 实时推送进度
- 无需轮询，服务器主动推送
- 实时更新进度、速度、ETA
- 支持多客户端同时监听

### 2. 断点续传支持
- 自动保存 ZIP 文件（如果有失败）
- 一键重试失败文件
- 记录已处理文件列表

### 3. 上传历史记录
- 保留最近 100 条记录
- 持久化到磁盘（JSON）
- 支持查询和清空

### 4. 详细的错误信息
- 每个文件的错误详情
- 是否已重试标记
- 失败文件列表

---

## 📡 API 接口

### 1. V2 上传
```http
POST /api/v1/article/zip/v2
Content-Type: multipart/form-data

file: <zip文件>
```

### 2. WebSocket 实时进度
```javascript
const ws = new WebSocket('ws://192.168.0.221:3002/api/v1/article/upload/v2/{task_id}/ws');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log(`进度: ${data.progress}%`);
  console.log(`速度: ${data.speed} 文件/秒`);
  console.log(`预计剩余: ${data.eta}`);
};
```

### 3. 重试失败文件
```http
POST /api/v1/article/upload/v2/{task_id}/retry
```

### 4. 上传历史
```http
GET /api/v1/article/upload/history
```

---

## 💻 完整前端示例

### Vue 3 + WebSocket 组件

```vue
<template>
  <div class="upload-v2">
    <!-- 文件选择 -->
    <input 
      type="file" 
      @change="handleFile" 
      accept=".zip" 
      :disabled="uploading"
    />
    
    <!-- 上传按钮 -->
    <button 
      @click="startUpload" 
      :disabled="!selectedFile || uploading"
    >
      开始上传
    </button>
    
    <!-- 进度条 -->
    <div v-if="uploading" class="progress-container">
      <div class="progress-bar">
        <div 
          class="progress-fill" 
          :style="{ width: `${progress}%` }"
        ></div>
      </div>
      
      <div class="progress-info">
        <span>{{ processed }}/{{ totalFiles }}</span>
        <span>{{ progress.toFixed(1) }}%</span>
        <span v-if="speed > 0">{{ speed.toFixed(1) }} 文件/秒</span>
        <span v-if="eta">剩余 {{ eta }}</span>
      </div>
      
      <!-- 取消按钮 -->
      <button @click="cancelUpload" class="btn-cancel">
        取消上传
      </button>
    </div>
    
    <!-- 结果 -->
    <div v-if="result" class="result">
      <h3>上传结果</h3>
      <p>成功: {{ result.success }} / {{ result.total }}</p>
      <p>失败: {{ result.failed }}</p>
      
      <!-- 错误列表 -->
      <div v-if="result.errors.length" class="errors">
        <h4>错误详情</h4>
        <ul>
          <li v-for="(error, index) in result.errors" :key="index">
            <strong>{{ error.file_name }}</strong>: {{ error.error }}
            <span v-if="error.retried" class="retried-badge">已重试</span>
          </li>
        </ul>
        
        <!-- 重试按钮 -->
        <button 
          v-if="result.failed > 0" 
          @click="retryFailed"
          class="btn-retry"
        >
          重试失败文件 ({{ result.failed }})
        </button>
      </div>
    </div>
    
    <!-- 历史记录 -->
    <div class="history">
      <h3>上传历史</h3>
      <button @click="loadHistory" class="btn-load">
        加载历史
      </button>
      
      <table v-if="history.length">
        <thead>
          <tr>
            <th>文件名</th>
            <th>总数</th>
            <th>成功</th>
            <th>失败</th>
            <th>耗时</th>
            <th>状态</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in history" :key="item.task_id">
            <td>{{ item.file_name }}</td>
            <td>{{ item.total_files }}</td>
            <td>{{ item.success }}</td>
            <td>{{ item.failed }}</td>
            <td>{{ item.duration }}</td>
            <td>
              <span :class="`status-${item.status}`">
                {{ item.status }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onUnmounted } from 'vue';

const selectedFile = ref(null);
const uploading = ref(false);
const progress = ref(0);
const processed = ref(0);
const totalFiles = ref(0);
const speed = ref(0);
const eta = ref('');
const result = ref(null);
const taskId = ref(null);
const history = ref([]);

let ws = null;

function handleFile(event) {
  selectedFile.value = event.target.files[0];
}

async function startUpload() {
  if (!selectedFile.value) return;
  
  uploading.value = true;
  progress.value = 0;
  result.value = null;
  
  const formData = new FormData();
  formData.append('file', selectedFile.value);
  
  try {
    // 1. 开始上传
    const response = await fetch('/api/v1/article/zip/v2', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${getToken()}`
      },
      body: formData
    });
    
    const data = await response.json();
    taskId.value = data.task_id;
    
    // 2. 建立 WebSocket 连接
    const wsUrl = `ws://${window.location.host}/api/v1/article/upload/v2/${data.task_id}/ws`;
    ws = new WebSocket(wsUrl);
    
    ws.onmessage = (event) => {
      const progressData = JSON.parse(event.data);
      
      progress.value = progressData.progress;
      processed.value = progressData.processed;
      totalFiles.value = progressData.total_files;
      speed.value = progressData.speed;
      eta.value = progressData.eta;
      
      // 上传完成
      if (progressData.status !== 'processing') {
        ws.close();
        uploading.value = false;
        result.value = progressData;
      }
    };
    
    ws.onerror = (error) => {
      console.error('WebSocket 错误:', error);
      // 降级到轮询
      startPolling();
    };
    
  } catch (error) {
    console.error('上传失败:', error);
    uploading.value = false;
  }
}

// 轮询降级方案
function startPolling() {
  const interval = setInterval(async () => {
    const response = await fetch(`/api/v1/article/upload/v2/${taskId.value}`, {
      headers: {
        'Authorization': `Bearer ${getToken()}`
      }
    });
    
    const data = await response.json();
    progress.value = data.progress;
    processed.value = data.processed;
    totalFiles.value = data.total_files;
    speed.value = data.speed;
    eta.value = data.eta;
    
    if (data.task_status !== 'processing') {
      clearInterval(interval);
      uploading.value = false;
      result.value = data;
    }
  }, 1000);
}

async function cancelUpload() {
  if (!taskId.value) return;
  
  try {
    await fetch(`/api/v1/article/upload/v2/${taskId.value}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${getToken()}`
      }
    });
    
    if (ws) {
      ws.close();
    }
    
    uploading.value = false;
  } catch (error) {
    console.error('取消失败:', error);
  }
}

async function retryFailed() {
  if (!taskId.value) return;
  
  try {
    const response = await fetch(`/api/v1/article/upload/v2/${taskId.value}/retry`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${getToken()}`
      }
    });
    
    const data = await response.json();
    console.log('开始重试:', data);
    
    // 重新监听进度
    startUpload();
  } catch (error) {
    console.error('重试失败:', error);
  }
}

async function loadHistory() {
  try {
    const response = await fetch('/api/v1/article/upload/history', {
      headers: {
        'Authorization': `Bearer ${getToken()}`
      }
    });
    
    const data = await response.json();
    history.value = data.data;
  } catch (error) {
    console.error('加载历史失败:', error);
  }
}

function getToken() {
  return localStorage.getItem('token');
}

onUnmounted(() => {
  if (ws) {
    ws.close();
  }
});
</script>

<style scoped>
.upload-v2 {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.progress-container {
  margin: 20px 0;
}

.progress-bar {
  width: 100%;
  height: 30px;
  background: #f0f0f0;
  border-radius: 15px;
  overflow: hidden;
  margin-bottom: 10px;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #4CAF50, #8BC34A);
  transition: width 0.3s ease;
}

.progress-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
  font-size: 14px;
  color: #666;
}

.errors {
  margin-top: 20px;
  padding: 15px;
  background: #fff3cd;
  border-radius: 8px;
}

.errors ul {
  list-style: none;
  padding: 0;
}

.errors li {
  padding: 8px;
  margin: 5px 0;
  background: white;
  border-radius: 4px;
}

.retried-badge {
  background: #ffc107;
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  margin-left: 8px;
}

.btn-cancel {
  background: #dc3545;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 5px;
  cursor: pointer;
}

.btn-retry {
  background: #ffc107;
  color: #333;
  border: none;
  padding: 10px 20px;
  border-radius: 5px;
  cursor: pointer;
  margin-top: 10px;
}

.btn-load {
  background: #007bff;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 5px;
  cursor: pointer;
  margin-bottom: 10px;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 10px;
}

th, td {
  padding: 10px;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

th {
  background: #f8f9fa;
  font-weight: bold;
}

.status-completed {
  color: #28a745;
  font-weight: bold;
}

.status-failed {
  color: #dc3545;
  font-weight: bold;
}

.status-cancelled {
  color: #ffc107;
  font-weight: bold;
}
</style>
```

---

## 🔧 配置参数

```go
var (
    maxConcurrentV2 = 3              // 最大并发数
    historyFile = "./data/upload_history.json"  // 历史文件路径
)

type UploadTaskV2 struct {
    MaxRetries: 3,  // 最大重试次数
}
```

---

## 📊 WebSocket 消息格式

### 进度消息
```json
{
  "type": "progress",
  "task_id": "upload_v2_1234567890",
  "file_name": "articles.zip",
  "total_files": 10,
  "processed": 6,
  "success": 5,
  "failed": 1,
  "progress": 60.0,
  "status": "processing",
  "speed": 2.5,
  "eta": "2秒",
  "retry_count": 0,
  "errors": [
    {
      "file_name": "article2.md",
      "error": "处理失败",
      "retried": false
    }
  ],
  "start_time": "2026-06-25T18:00:00Z",
  "end_time": null
}
```

---

## ⚠️ 注意事项

1. **WebSocket 降级**: 如果 WebSocket 不可用，自动降级到轮询
2. **ZIP 保留**: 如果有失败文件，保留 ZIP 用于重试
3. **历史限制**: 最多保留 100 条记录
4. **并发限制**: 全局最多 3 个并发任务
5. **生产环境**: WebSocket 应该检查 Origin 防止 CSRF

---

## 🎯 性能对比

| 特性 | V1 | V2 |
|------|----|----|
| **进度更新** | 轮询 | WebSocket 实时推送 |
| **延迟** | 1-2秒 | <100ms |
| **服务器压力** | 高（频繁请求） | 低（长连接） |
| **断点续传** | ❌ | ✅ |
| **重试机制** | ❌ | ✅ |
| **历史记录** | ❌ | ✅ |
| **速度显示** | ❌ | ✅ |
| **ETA** | ❌ | ✅ |
