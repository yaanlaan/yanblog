<template>
  <div class="backend-config">
    <el-card class="config-card">
      <template #header>
        <div class="card-header">
          <span>后端配置管理</span>
          <el-button type="primary" @click="saveConfig" :loading="saving" :icon="Check">
            保存配置
          </el-button>
        </div>
      </template>

      <el-alert
        title="修改后端配置后需要重启服务才能生效"
        type="warning"
        show-icon
        :closable="false"
        class="mb-4"
      />

      <div v-loading="loading" class="config-body">
        <!-- 服务器配置 -->
        <div class="config-section">
          <h3 class="section-title">
            <el-icon><Monitor /></el-icon>
            服务器配置
          </h3>
          <el-form label-width="120px" class="config-form">
            <el-form-item label="运行模式">
              <el-select v-model="config.server.AppMode">
                <el-option label="Debug" value="debug" />
                <el-option label="Release" value="release" />
              </el-select>
            </el-form-item>
            <el-form-item label="监听端口">
              <el-input v-model="config.server.HttpPort" placeholder="如 :8080" />
            </el-form-item>
            <el-form-item label="站点 URL">
              <el-input v-model="config.server.SiteUrl" placeholder="如 https://example.com" />
              <div class="form-tip">用于生成站点地图、SEO 等功能的完整链接</div>
            </el-form-item>
          </el-form>
        </div>

        <!-- 数据库配置 -->
        <div class="config-section">
          <h3 class="section-title">
            <el-icon><Coin /></el-icon>
            数据库配置
          </h3>
          <el-form label-width="120px" class="config-form">
            <el-form-item label="数据库类型">
              <el-select v-model="config.database.Db">
                <el-option label="SQLite" value="SQLite" />
                <el-option label="MySQL" value="MYSQL" />
              </el-select>
            </el-form-item>
            <template v-if="config.database.Db !== 'SQLite'">
              <el-form-item label="主机地址">
                <el-input v-model="config.database.DbHost" placeholder="localhost" />
              </el-form-item>
              <el-form-item label="端口">
                <el-input-number v-model="config.database.DbPort" :min="1" :max="65535" />
              </el-form-item>
              <el-form-item label="用户名">
                <el-input v-model="config.database.DbUser" placeholder="root" />
              </el-form-item>
              <el-form-item label="密码">
                <el-input
                  v-model="dbPassword"
                  type="password"
                  show-password
                  placeholder="留空则不修改密码"
                />
                <div class="form-tip">出于安全考虑，当前密码不会回显。输入新密码可修改。</div>
              </el-form-item>
            </template>
            <el-form-item label="数据库名称">
              <el-input v-model="config.database.DbName" placeholder="如 data/yanblog.db" />
              <div class="form-tip">SQLite 时为文件路径；MySQL 时为数据库名</div>
            </el-form-item>
          </el-form>
        </div>

        <!-- 天气设置 -->
        <div class="config-section">
          <h3 class="section-title">
            <el-icon><Sunny /></el-icon>
            天气设置
          </h3>
          <el-form label-width="120px" class="config-form">
            <el-form-item label="天气服务">
              <div class="form-tip">
                使用 <a href="https://open-meteo.com" target="_blank">Open-Meteo</a> 免费天气 API，无需注册和密钥。
              </div>
            </el-form-item>
            <el-form-item label="默认城市">
              <el-input v-model="config.weather.DefaultCity" placeholder="如 Hefei, Shanghai" />
              <div class="form-tip">前端未指定城市时使用的默认城市（支持中英文）</div>
            </el-form-item>
          </el-form>
        </div>

        <!-- 其他设置 -->
        <div class="config-section">
          <h3 class="section-title">
            <el-icon><Setting /></el-icon>
            其他设置
          </h3>
          <el-form label-width="120px" class="config-form">
            <el-form-item label="前端配置路径">
              <el-input v-model="config.FrontEndConfigPath" placeholder="config/frontend/config.yaml" />
            </el-form-item>
          </el-form>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Check, Monitor, Coin, Sunny, Setting } from '@element-plus/icons-vue'
import { systemApi } from '@/services/api'

defineOptions({
  name: 'BackendConfig'
})

const loading = ref(false)
const saving = ref(false)
const dbPassword = ref('')

const config = reactive({
  server: {
    AppMode: 'debug',
    HttpPort: ':8080',
    SiteUrl: ''
  },
  database: {
    Db: 'SQLite',
    DbHost: 'localhost',
    DbPort: 3306,
    DbUser: 'root',
    DbPassWord: '',
    DbName: 'data/yanblog.db'
  },
  weather: {
    DefaultCity: 'Hefei'
  },
  FrontEndConfigPath: 'config/frontend/config.yaml'
})

const loadConfig = async () => {
  loading.value = true
  try {
    const res = await systemApi.getBackendConfig()
    if (res.data.status === 200) {
      const data = res.data.data
      if (data.server) Object.assign(config.server, data.server)
      if (data.database) {
        Object.assign(config.database, data.database)
        // 密码不回显
        config.database.DbPassWord = ''
      }
      if (data.weather) Object.assign(config.weather, data.weather)
      if (data.FrontEndConfigPath) config.FrontEndConfigPath = data.FrontEndConfigPath
    }
  } catch (e) {
    ElMessage.error('加载后端配置失败')
    console.error(e)
  } finally {
    loading.value = false
  }
}

const saveConfig = async () => {
  saving.value = true
  try {
    const payload: any = { ...config }
    // 只在用户输入了新密码时才发送密码
    if (dbPassword.value) {
      payload.database = { ...config.database, DbPassWord: dbPassword.value }
    } else {
      // 不发送密码字段，避免覆盖
      payload.database = { ...config.database }
      delete payload.database.DbPassWord
    }

    const res = await systemApi.updateBackendConfig(payload)
    if (res.data.status === 200) {
      ElMessage.success('后端配置保存成功，请重启服务使其生效')
      dbPassword.value = '' // 清空密码输入
    } else {
      ElMessage.error(res.data.message || '保存失败')
    }
  } catch (e) {
    ElMessage.error('保存失败')
    console.error(e)
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadConfig()
})
</script>

<style scoped>
.backend-config {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.mb-4 {
  margin-bottom: 16px;
}

.config-body {
  max-width: 800px;
}

.config-section {
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid #ebeef5;
}

.config-section:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  margin: 0 0 20px 0;
}

.section-title .el-icon {
  color: #409eff;
}

.config-form {
  margin-top: 0;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
  margin-top: 4px;
}
</style>
