<template>
  <div class="config-editor">
    <el-card class="editor-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span>前端配置管理</span>
            <el-tag :type="mode === 'form' ? 'success' : 'info'" class="mode-tag">
              {{ mode === 'form' ? '可视化模式' : '源码模式' }}
            </el-tag>
          </div>
          <div class="header-actions">
            <el-button-group class="mode-switch">
              <el-button :type="mode === 'form' ? 'primary' : ''" @click="switchMode('form')" :icon="Operation">
                可视化编辑
              </el-button>
              <el-button :type="mode === 'yaml' ? 'primary' : ''" @click="switchMode('yaml')" :icon="Document">
                YAML源码
              </el-button>
            </el-button-group>

            <el-divider direction="vertical" />

            <el-button @click="resetConfig" :icon="Refresh">重置</el-button>
            <el-button type="primary" @click="saveConfig" :loading="saving" :icon="Check">保存更改</el-button>
          </div>
        </div>
      </template>

      <!-- 源码模式 -->
      <div v-show="mode === 'yaml'" class="yaml-mode">
        <el-alert title="直接编辑 YAML 源码，请注意缩进格式。" type="info" show-icon :closable="false" class="mb-4" />
        <el-input v-model="configContent" type="textarea" :rows="25" placeholder="正在加载配置文件..." class="yaml-editor"
          spellcheck="false" />
      </div>

      <!-- 可视化模式 -->
      <div v-if="mode === 'form'" class="form-mode">
        <el-tabs v-model="activeTab" class="config-tabs">
          <!-- 基本设置 -->
          <el-tab-pane label="基本信息" name="basic">
            <el-form label-width="120px" class="config-form">
              <el-form-item label="博客名称">
                <el-input v-model="configForm.blog_name" />
              </el-form-item>

              <el-form-item label="Logo文字">
                <el-input v-model="configForm.logo_text" />
              </el-form-item>

              <el-form-item label="Logo图标">
                <image-uploader v-model="configForm.logo_image" />
              </el-form-item>

              <el-form-item label="Favicon">
                <image-uploader v-model="configForm.favicon" />
              </el-form-item>

              <el-form-item label="Admin地址">
                <el-input v-model="configForm.admin_url" />
              </el-form-item>

              <el-form-item label="Iconfont URL">
                <el-input v-model="configForm.iconfont_url" />
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <!-- 作者信息 -->
          <el-tab-pane label="作者信息" name="author">
            <el-form label-width="120px" class="config-form">
              <el-form-item label="作者昵称">
                <el-input v-model="configForm.author_name" />
              </el-form-item>
              <el-form-item label="个性签名">
                <el-input v-model="configForm.author_bio" type="textarea" :rows="2" />
              </el-form-item>
              <el-form-item label="头像">
                <image-uploader v-model="configForm.author_avatar" />
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <!-- 首页 Hero -->
          <el-tab-pane label="首页设置" name="hero">
            <el-form label-width="120px" class="config-form" v-if="configForm.hero">
              <el-form-item label="主标题">
                <el-input v-model="configForm.hero.title" />
                <div class="form-tip">支持 HTML 标签，如 &lt;br&gt;</div>
              </el-form-item>
              <el-form-item label="副标题">
                <el-input v-model="configForm.hero.subtitle" />
              </el-form-item>
              <el-form-item label="欢迎语">
                <el-input v-model="configForm.hero.welcome" type="textarea" :rows="2" />
              </el-form-item>
              <el-form-item label="欢迎图片">
                <image-uploader v-model="configForm.hero.welcome_image" />
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <!-- 名言/Quotes -->
          <el-tab-pane label="名言语录" name="quotes">
            <div class="list-editor">
              <el-button type="primary" plain size="small" @click="addQuote" class="mb-4">添加语录</el-button>
              <div v-for="(quote, index) in configForm.quotes" :key="index" class="list-item">
                <el-input v-model="configForm.quotes[index]" placeholder="请输入内容" class="mr-2">
                  <template #prepend>{{ index + 1 }}</template>
                </el-input>
                <el-button type="danger" :icon="Delete" circle @click="removeQuote(index)" />
              </div>
            </div>
          </el-tab-pane>

          <!-- 快捷链接 -->
          <el-tab-pane label="快捷链接" name="shortcuts">
            <div class="shortcuts-editor">
              <el-button type="primary" plain size="small" @click="addShortcut" class="mb-4">添加链接</el-button>
              <el-table :data="configForm.shortcuts" style="width: 100%" border>
                <el-table-column label="名称" width="150">
                  <template #default="{ row }">
                    <el-input v-model="row.name" size="small" />
                  </template>
                </el-table-column>
                <el-table-column label="链接">
                  <template #default="{ row }">
                    <el-input v-model="row.url" size="small" />
                  </template>
                </el-table-column>
                <el-table-column label="图标" width="150">
                  <template #default="{ row }">
                    <el-input v-model="row.icon" size="small" />
                  </template>
                </el-table-column>
                <el-table-column label="颜色代码">
                  <template #default="{ row }">
                    <el-input v-model="row.color" size="small" />
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="80" align="center">
                  <template #default="{ $index }">
                    <el-button type="danger" :icon="Delete" circle size="small" @click="removeShortcut($index)" />
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-tab-pane>

          <!-- 页面标题 -->
          <el-tab-pane label="页面标题" name="pagetitle">
            <el-form label-width="120px" class="config-form" v-if="configForm.page_title">
              <el-form-item label="默认标题">
                <el-input v-model="configForm.page_title.default" />
              </el-form-item>
              <el-form-item label="失焦标题">
                <el-input v-model="configForm.page_title.blur" />
                <div class="form-tip">页面失去焦点时（切换Tab）显示的搞怪标题</div>
              </el-form-item>
            </el-form>
          </el-tab-pane>

        </el-tabs>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, defineComponent, h } from 'vue'
import {
  ElMessage,
  ElMessageBox,
  ElUpload,
  ElButton,
  ElIcon
} from 'element-plus'
import type { UploadProps } from 'element-plus'
import {
  Check,
  Refresh,
  Operation,
  Document,
  Plus,
  Delete
} from '@element-plus/icons-vue'
import { systemApi } from '@/services/api'
import yaml from 'js-yaml'

defineOptions({
  name: 'ConfigEditor'
})

// --- 简单的图片上传内联组件 ---
const ImageUploader = defineComponent({
  props: ['modelValue'],
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const handleSuccess: UploadProps['onSuccess'] = (response) => {
      if (response.status === 200) {
        emit('update:modelValue', response.url)
        ElMessage.success('上传成功')
      } else {
        ElMessage.error(response.message || '上传失败')
      }
    }

    const handleError = () => {
      ElMessage.error('上传出错')
    }

    const beforeUpload: UploadProps['beforeUpload'] = (rawFile) => {
      if (rawFile.size / 1024 / 1024 > 5) {
        ElMessage.error('图片大小不能超过 5MB!')
        return false
      }
      return true
    }

    // 预览图片
    const previewImage = () => {
       if (props.modelValue) {
         window.open(props.modelValue, '_blank')
       }
    }

    return () => h('div', { class: 'custom-image-uploader' }, [
      h(ElUpload, {
        class: 'avatar-uploader',
        action: '/api/v1/upload',
        data: { type: 'common' }, // 使用 common 类型
        showFileList: false,
        onSuccess: handleSuccess,
        onError: handleError,
        beforeUpload: beforeUpload,
        accept: 'image/*'
      }, {
        default: () => [
          props.modelValue ?
            h('img', { src: props.modelValue, class: 'avatar' }) :
            h(ElIcon, { class: 'avatar-uploader-icon' }, () => h(Plus))
        ]
      }),
      // 如果有图片，显示预览URL和操作
      props.modelValue && h('div', { class: 'image-actions' }, [
        h('span', { class: 'image-url', title: props.modelValue }, props.modelValue),
        h(ElButton, {
           size: 'small', 
           link: true, 
           type: 'primary',
           onClick: previewImage
        }, () => '查看原图')
      ])
    ])
  }
})

// --- 主逻辑 ---
const mode = ref<'yaml' | 'form'>('form')
const activeTab = ref('basic')
const configContent = ref('')
const configForm = ref<any>({})
const saving = ref(false)

const loadConfig = async () => {
  try {
    const res = await systemApi.getFrontEndConfig()
    if (res.data.status === 200) {
      configContent.value = res.data.data
      try {
        configForm.value = yaml.load(configContent.value) || {}
      } catch (e) {
        console.error('Initial YAML parse error', e)
        mode.value = 'yaml'
        ElMessage.warning('配置文件格式复杂或有误，已切换到源码模式')
      }
    } else {
      ElMessage.error(res.data.message || '加载配置失败')
    }
  } catch (error) {
    console.error(error)
    ElMessage.error('网络错误，无法加载配置')
  }
}

const switchMode = (targetMode: 'yaml' | 'form') => {
  if (targetMode === mode.value) return

  if (targetMode === 'form') {
    // YAML -> Form
    try {
      const parsed = yaml.load(configContent.value)
      if (typeof parsed !== 'object' || parsed === null) {
        throw new Error('YAML must evaluate to an object')
      }
      configForm.value = parsed
      mode.value = 'form'
    } catch (e: any) {
      ElMessage.error('YAML 语法错误，无法切换到可视化模式: ' + e.message)
    }
  } else {
    // Form -> YAML
    try {
      configContent.value = yaml.dump(configForm.value)
      mode.value = 'yaml'
    } catch (e: any) {
      ElMessage.error('转换失败: ' + e.message)
    }
  }
}

const saveConfig = async () => {
  let contentToSend = configContent.value

  if (mode.value === 'form') {
    try {
      contentToSend = yaml.dump(configForm.value)
      configContent.value = contentToSend
    } catch (e: any) {
      ElMessage.error('生成配置失败: ' + e.message)
      return
    }
  }

  try {
    saving.value = true
    const res = await systemApi.updateFrontEndConfig({ content: contentToSend })
    if (res.data.status === 200) {
      ElMessage.success('配置保存成功，刷新前台页面生效')
    } else {
      ElMessage.error(res.data.message || '保存配置失败')
    }
  } catch (error) {
    console.error(error)
    ElMessage.error('保存失败，请检查网络')
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

// 数组操作辅助函数
const addQuote = () => {
  if (!configForm.value.quotes) configForm.value.quotes = []
  configForm.value.quotes.push('')
}

const removeQuote = (index: number) => {
  configForm.value.quotes.splice(index, 1)
}

const addShortcut = () => {
  if (!configForm.value.shortcuts) configForm.value.shortcuts = []
  configForm.value.shortcuts.push({ name: 'New Link', url: '#', icon: '', color: '' })
}

const removeShortcut = (index: number) => {
  configForm.value.shortcuts.splice(index, 1)
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

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.yaml-editor :deep(.el-textarea__inner) {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
}

.config-form {
  max-width: 800px;
  margin-top: 20px;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
  margin-top: 4px;
}

.mb-4 {
  margin-bottom: 16px;
}

.mr-2 {
  margin-right: 8px;
}

.list-item {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

/* 动态组件样式 */
:deep(.avatar-uploader .el-upload) {
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
  width: 140px; 
  height: 140px;
  background-color: #fafafa;
}

:deep(.avatar-uploader .el-upload:hover) {
  border-color: var(--el-color-primary);
}

:deep(.avatar-uploader-icon) {
  font-size: 28px;
  color: #8c939d;
  width: 140px;
  height: 140px;
  text-align: center;
  line-height: 140px;
  display: flex;
  align-items: center;
  justify-content: center;
}

:deep(.avatar) {
  width: 140px;
  height: 140px;
  display: block;
  object-fit: contain; /* 改为 contain 以免裁切重要logo */
  background-color: #f0f0f0;
}

:deep(.custom-image-uploader) {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 10px;
}

:deep(.image-actions) {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 12px;
  color: #666;
  max-width: 100%;
}

:deep(.image-url) {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 300px;
  background: #f4f4f5;
  padding: 2px 6px;
  border-radius: 4px;
}
</style>
