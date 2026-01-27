<template>
  <div class="config-editor">
    <el-card class="editor-card">
      <template #header>
        <div class="card-header">
          <span>前端配置管理 (config.yaml)</span>
          <div>
            <el-button @click="resetConfig" :icon="Refresh">重置/刷新</el-button>
            <el-button type="primary" @click="saveConfig" :loading="saving" :icon="Check">保存配置</el-button>
          </div>
        </div>
      </template>

      <div class="tip-box">
        <el-alert
          title="警告：请遵循 YAML 语法格式进行修改。错误的格式可能导致前台页面崩溃。修改后刷新前台页面即可生效。"
          type="warning"
          show-icon
          :closable="false"
        />
      </div>

      <el-input
        v-model="configContent"
        type="textarea"
        :rows="25"
        placeholder="正在加载配置文件..."
        class="yaml-editor"
        spellcheck="false"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Check, Refresh } from '@element-plus/icons-vue'
import { systemApi } from '@/services/api' // We need to add this

const configContent = ref('')
const saving = ref(false)

const loadConfig = async () => {
  try {
    const res = await systemApi.getFrontEndConfig()
    if (res.data.status === 200) {
      configContent.value = res.data.data
    } else {
      ElMessage.error(res.data.message || '加载配置失败')
    }
  } catch (error) {
    console.error(error)
    ElMessage.error('网络错误，无法加载配置')
  }
}

const saveConfig = async () => {
  if (!configContent.value) return

  try {
    saving.value = true
    const res = await systemApi.updateFrontEndConfig({ content: configContent.value })
    if (res.data.status === 200) {
      ElMessage.success('配置保存成功')
    } else {
      ElMessage.error(res.data.message || '保存配置失败')
    }
  } catch (error) {
    console.error(error)
    ElMessage.error('保存失败，请检查网络或配置格式')
  } finally {
    saving.value = false
  }
}

const resetConfig = () => {
  ElMessageBox.confirm(
    '确定要放弃当前的修改并重新加载吗？',
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  )
    .then(() => {
      loadConfig()
    })
    .catch(() => {})
}

onMounted(() => {
  loadConfig()
})
</script>

<style scoped>
.config-editor {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.tip-box {
  margin-bottom: 20px;
}

.yaml-editor :deep(.el-textarea__inner) {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
  background-color: #f8f9fa;
  color: #333;
}
</style>
