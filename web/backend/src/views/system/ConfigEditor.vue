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

              <el-form-item label="允许访问域名">
                  <el-select
                    v-model="configForm.allowed_hosts"
                    multiple
                    filterable
                    allow-create
                    default-first-option
                    :reserve-keyword="false"
                    placeholder="输入域名并回车添加"
                  >
                  </el-select>
                  <div class="form-tip">用于 Vite 开发环境的 allowedHosts 设置，修改后可能需要重启服务</div>
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

          <!-- 评论配置 -->
          <el-tab-pane label="评论系统" name="comment">
            <el-form label-width="120px" class="config-form" v-if="configForm.comment">
              <el-form-item label="启用评论">
                <el-switch v-model="configForm.comment.enable" />
              </el-form-item>
              
              <template v-if="configForm.comment.enable">
                <el-form-item label="类型">
                  <el-select v-model="configForm.comment.type">
                    <el-option label="Giscus" value="giscus" />
                  </el-select>
                </el-form-item>

                <div v-if="configForm.comment.type === 'giscus' && configForm.comment.giscus" class="sub-config">
                  <el-divider content-position="left">Giscus 配置</el-divider>
                  <el-form-item label="GitHub 仓库">
                    <el-input v-model="configForm.comment.giscus.repo" placeholder="username/repo" />
                  </el-form-item>
                  <el-form-item label="仓库 ID">
                    <el-input v-model="configForm.comment.giscus.repo_id" />
                  </el-form-item>
                  <el-form-item label="分类 (Category)">
                    <el-input v-model="configForm.comment.giscus.category" />
                  </el-form-item>
                  <el-form-item label="分类 ID">
                    <el-input v-model="configForm.comment.giscus.category_id" />
                  </el-form-item>
                </div>
              </template>
            </el-form>
            <div v-else class="empty-tip">
              配置文件中缺少该字段，请切换到源码模式手动添加。
            </div>
          </el-tab-pane>

          <!-- 页脚设置 -->
          <el-tab-pane label="页脚设置" name="footer">
            <el-form label-width="120px" class="config-form" v-if="configForm.footer">
              <el-form-item label="副标题">
                <el-input v-model="configForm.footer.tagline" placeholder="博客名下方的副标题" />
              </el-form-item>
              <el-form-item label="版权信息">
                <el-input v-model="configForm.footer.copyright" />
              </el-form-item>
              <el-form-item label="联系邮箱">
                <el-input v-model="configForm.footer.email" />
              </el-form-item>
              <el-form-item label="Powered By">
                <el-input v-model="configForm.footer.powered_by" />
              </el-form-item>
              <el-form-item label="Powered By 链接">
                <el-input v-model="configForm.footer.powered_by_link" />
              </el-form-item>

              <el-divider content-position="left">ICP 备案</el-divider>
              <el-form-item label="显示 ICP">
                <el-switch v-model="configForm.footer.icp.show" />
              </el-form-item>
              <el-form-item label="ICP 文本" v-if="configForm.footer.icp.show">
                <el-input v-model="configForm.footer.icp.text" />
              </el-form-item>
              <el-form-item label="ICP 链接" v-if="configForm.footer.icp.show">
                <el-input v-model="configForm.footer.icp.link" />
              </el-form-item>

              <el-divider content-position="left">作品集</el-divider>
              <el-form-item label="显示作品集">
                <el-switch v-model="configForm.footer.portfolio.show" />
              </el-form-item>
              <template v-if="configForm.footer.portfolio.show">
                <el-form-item label="标题">
                  <el-input v-model="configForm.footer.portfolio.title" />
                </el-form-item>
                <el-form-item label="图标">
                  <el-input v-model="configForm.footer.portfolio.icon" />
                </el-form-item>
                <el-form-item label="项目列表">
                  <el-button type="primary" plain size="small" @click="addPortfolioItem">添加项目</el-button>
                  <el-table :data="configForm.footer.portfolio.items" border class="mt-8">
                    <el-table-column label="名称">
                      <template #default="{ row }"><el-input v-model="row.name" size="small" /></template>
                    </el-table-column>
                    <el-table-column label="链接">
                      <template #default="{ row }"><el-input v-model="row.url" size="small" /></template>
                    </el-table-column>
                    <el-table-column label="图标">
                      <template #default="{ row }"><el-input v-model="row.icon" size="small" /></template>
                    </el-table-column>
                    <el-table-column label="操作" width="80" align="center">
                      <template #default="{ $index }">
                        <el-button type="danger" :icon="Delete" circle size="small" @click="removePortfolioItem($index)" />
                      </template>
                    </el-table-column>
                  </el-table>
                </el-form-item>
              </template>

              <el-divider content-position="left">相关链接</el-divider>
              <el-form-item label="显示相关链接">
                <el-switch v-model="configForm.footer.related_links.show" />
              </el-form-item>
              <template v-if="configForm.footer.related_links.show">
                <el-form-item label="标题">
                  <el-input v-model="configForm.footer.related_links.title" />
                </el-form-item>
                <el-form-item label="链接列表">
                  <el-button type="primary" plain size="small" @click="addRelatedLink">添加链接</el-button>
                  <el-table :data="configForm.footer.related_links.items" border class="mt-8">
                    <el-table-column label="名称">
                      <template #default="{ row }"><el-input v-model="row.name" size="small" /></template>
                    </el-table-column>
                    <el-table-column label="链接">
                      <template #default="{ row }"><el-input v-model="row.url" size="small" /></template>
                    </el-table-column>
                    <el-table-column label="图标">
                      <template #default="{ row }"><el-input v-model="row.icon" size="small" /></template>
                    </el-table-column>
                    <el-table-column label="操作" width="80" align="center">
                      <template #default="{ $index }">
                        <el-button type="danger" :icon="Delete" circle size="small" @click="removeRelatedLink($index)" />
                      </template>
                    </el-table-column>
                  </el-table>
                </el-form-item>
              </template>
            </el-form>
            <div v-else class="empty-tip">配置文件中缺少 footer 字段，请先在源码模式添加。</div>
          </el-tab-pane>

          <!-- 社交链接 -->
          <el-tab-pane label="社交链接" name="socials">
            <div class="list-editor">
              <el-button type="primary" plain size="small" @click="addSocial" class="mb-4">添加社交链接</el-button>
              <el-table :data="configForm.socials" border>
                <el-table-column label="名称">
                  <template #default="{ row }"><el-input v-model="row.name" size="small" /></template>
                </el-table-column>
                <el-table-column label="链接">
                  <template #default="{ row }"><el-input v-model="row.url" size="small" /></template>
                </el-table-column>
                <el-table-column label="图标" width="160">
                  <template #default="{ row }"><el-input v-model="row.icon" size="small" /></template>
                </el-table-column>
                <el-table-column label="颜色" width="120">
                  <template #default="{ row }"><el-input v-model="row.color" size="small" /></template>
                </el-table-column>
                <el-table-column label="圆形" width="70" align="center">
                  <template #default="{ row }"><el-switch v-model="row.is_circle" size="small" /></template>
                </el-table-column>
                <el-table-column label="操作" width="80" align="center">
                  <template #default="{ $index }">
                    <el-button type="danger" :icon="Delete" circle size="small" @click="removeSocial($index)" />
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-tab-pane>

          <!-- 联系方式 -->
          <el-tab-pane label="联系方式" name="contacts">
            <el-form label-width="120px" class="config-form" v-if="configForm.contacts">
              <el-form-item label="显示联系方式">
                <el-switch v-model="configForm.contacts.show" />
              </el-form-item>
              <el-form-item label="微信二维码">
                <image-uploader v-model="configForm.contacts.wechat_qr" />
              </el-form-item>
              <el-divider content-position="left">联系项目</el-divider>
              <el-button type="primary" plain size="small" @click="addContactItem" class="mb-4">添加联系方式</el-button>
              <el-table :data="configForm.contacts.items" border>
                <el-table-column label="名称">
                  <template #default="{ row }"><el-input v-model="row.name" size="small" /></template>
                </el-table-column>
                <el-table-column label="链接">
                  <template #default="{ row }"><el-input v-model="row.url" size="small" /></template>
                </el-table-column>
                <el-table-column label="图标" width="160">
                  <template #default="{ row }"><el-input v-model="row.icon" size="small" /></template>
                </el-table-column>
                <el-table-column label="颜色" width="120">
                  <template #default="{ row }"><el-input v-model="row.color" size="small" /></template>
                </el-table-column>
                <el-table-column label="圆形" width="70" align="center">
                  <template #default="{ row }"><el-switch v-model="row.is_circle" size="small" /></template>
                </el-table-column>
                <el-table-column label="操作" width="80" align="center">
                  <template #default="{ $index }">
                    <el-button type="danger" :icon="Delete" circle size="small" @click="removeContactItem($index)" />
                  </template>
                </el-table-column>
              </el-table>
            </el-form>
            <div v-else class="empty-tip">配置文件中缺少 contacts 字段，请先在源码模式添加。</div>
          </el-tab-pane>

          <!-- 音乐播放器 -->
          <el-tab-pane label="音乐播放器" name="music">
            <el-form label-width="120px" class="config-form" v-if="configForm.music_player">
              <el-form-item label="显示播放器">
                <el-switch v-model="configForm.music_player.show" />
              </el-form-item>
              <el-form-item label="播放器 URL" v-if="configForm.music_player.show">
                <el-input v-model="configForm.music_player.url" placeholder="如：//music.163.com/outchain/player..." />
              </el-form-item>
            </el-form>
            <div v-else class="empty-tip">配置文件中缺少 music_player 字段，请先在源码模式添加。</div>
          </el-tab-pane>

          <!-- 默认图片 -->
          <el-tab-pane label="默认图片" name="defaultImages">
            <el-form label-width="120px" class="config-form">
              <div class="form-tip" style="margin-bottom:16px">文章无封面、头像加载失败、二维码缺失时的回退图片</div>
              <el-form-item label="默认封面">
                <image-uploader v-model="configForm.default_images.cover" />
              </el-form-item>
              <el-form-item label="默认头像">
                <image-uploader v-model="configForm.default_images.avatar" />
              </el-form-item>
              <el-form-item label="默认二维码">
                <image-uploader v-model="configForm.default_images.qr_code" />
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <!-- 开发设置 -->
          <el-tab-pane label="开发设置" name="dev">
            <el-form label-width="140px" class="config-form">
              <el-form-item label="开发环境 Admin 端口">
                <el-input-number v-model="configForm.dev_admin_port" :min="1024" :max="65535" />
                <div class="form-tip">本地开发时 admin 面板的端口号，Docker 部署会自动覆盖</div>
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
  ElUpload,
  ElButton,
  ElIcon
} from 'element-plus'
import type { UploadProps } from 'element-plus'
import {
  Check,
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
        // 确保数组存在
        if (!Array.isArray(configForm.value.allowed_hosts)) {
          configForm.value.allowed_hosts = []
        if (!configForm.value.default_images) {
          configForm.value.default_images = { cover: '', avatar: '', qr_code: '' }
        }
        }
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

// 页脚 — 作品集 & 相关链接
const addPortfolioItem = () => {
  if (!configForm.value.footer?.portfolio) return
  if (!configForm.value.footer.portfolio.items) configForm.value.footer.portfolio.items = []
  configForm.value.footer.portfolio.items.push({ name: '', url: '', icon: '' })
}
const removePortfolioItem = (i: number) => configForm.value.footer?.portfolio?.items?.splice(i, 1)

const addRelatedLink = () => {
  if (!configForm.value.footer?.related_links) return
  if (!configForm.value.footer.related_links.items) configForm.value.footer.related_links.items = []
  configForm.value.footer.related_links.items.push({ name: '', url: '', icon: '' })
}
const removeRelatedLink = (i: number) => configForm.value.footer?.related_links?.items?.splice(i, 1)

// 社交链接
const addSocial = () => {
  if (!Array.isArray(configForm.value.socials)) configForm.value.socials = []
  configForm.value.socials.push({ name: '', url: '', icon: '', color: '', is_circle: true })
}
const removeSocial = (i: number) => configForm.value.socials?.splice(i, 1)

// 联系方式
const addContactItem = () => {
  if (!configForm.value.contacts) configForm.value.contacts = { show: false, wechat_qr: '', items: [] }
  if (!configForm.value.contacts.items) configForm.value.contacts.items = []
  configForm.value.contacts.items.push({ name: '', url: '', icon: '', color: '', is_circle: true })
}
const removeContactItem = (i: number) => configForm.value.contacts?.items?.splice(i, 1)

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

.sub-config {
  margin-top: 15px;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 4px;
  border: 1px solid #eee;
}

.empty-tip {
  color: #909399;
  font-size: 14px;
  text-align: center;
  padding: 20px;
}
</style>
