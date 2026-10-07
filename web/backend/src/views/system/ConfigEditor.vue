<template>
  <div class="config-editor">
    <el-card class="editor-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span class="card-title">个性化与前端配置</span>
            <el-tag :type="mode === 'form' ? 'success' : 'info'" class="mode-tag" effect="light">
              {{ mode === 'form' ? '可视化配置模式' : 'YAML源码模式' }}
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

            <el-button type="primary" @click="saveConfig" :loading="saving" :icon="Check">
              保存并热应用
            </el-button>
          </div>
        </div>
      </template>

      <!-- 源码模式 -->
      <div v-show="mode === 'yaml'" class="yaml-mode">
        <el-alert title="直接编辑 config.yaml 源码，修改保存后即时生效，无需重新打包编译。" type="info" show-icon :closable="false" class="mb-4" />
        <el-input 
          v-model="configContent" 
          type="textarea" 
          :rows="26" 
          placeholder="正在加载配置文件..." 
          class="yaml-editor"
          spellcheck="false" 
        />
      </div>

      <!-- 可视化模式 -->
      <div v-if="mode === 'form'" class="form-mode">
        <el-tabs v-model="activeTab" class="config-tabs">
          <!-- 基本设置 -->
          <el-tab-pane label="基本设置与图像" name="basic">
            <el-form label-width="130px" class="config-form">
              <el-form-item label="博客名称">
                <el-input v-model="configForm.blog_name" placeholder="博客全站名称" />
              </el-form-item>

              <el-form-item label="Logo文字">
                <el-input v-model="configForm.logo_text" placeholder="导航栏 Logo 文字" />
              </el-form-item>

              <el-form-item label="Logo图标">
                <ResourceImagePicker 
                  v-model="configForm.logo_image" 
                  label="Logo 图标"
                  upload-type="system"
                  :preset-options="logoPresets"
                  tip="显示在导航栏左上角的博客标志。支持直接上传、从媒体库选择或输入外链。"
                />
              </el-form-item>

              <el-form-item label="网站 Favicon">
                <ResourceImagePicker 
                  v-model="configForm.favicon" 
                  label="Favicon 图标"
                  upload-type="system"
                  :preset-options="faviconPresets"
                  tip="浏览器标签页小图标。支持 svg、ico、png 格式，保存后热生效。"
                />
              </el-form-item>

              <el-form-item label="全局网页背景">
                <ResourceImagePicker 
                  v-model="configForm.background_image" 
                  label="全站背景图"
                  upload-type="system"
                  aspect-ratio="wide"
                  :preset-options="globalBgPresets"
                  tip="整个博客页面的底层壁纸。上传新图即可全站热修改生效；留空则自动使用内置经典小猫背景。"
                />
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
                />
                <div class="form-tip">用于 Vite 开发环境的 allowedHosts 设置</div>
              </el-form-item>

              <el-form-item label="Admin管理路径">
                <el-input v-model="configForm.admin_url" placeholder="/admin" />
              </el-form-item>

              <el-form-item label="Iconfont URL">
                <el-input v-model="configForm.iconfont_url" placeholder="自定义阿里矢量图标库 CSS 链接" />
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <!-- 作者信息与签名 -->
          <el-tab-pane label="博主信息与简言" name="author">
            <el-form label-width="130px" class="config-form">
              <el-form-item label="作者昵称">
                <el-input v-model="configForm.author_name" placeholder="博主显示昵称" />
              </el-form-item>

              <el-form-item label="个性签名 (简言)">
                <el-input 
                  v-model="configForm.author_bio" 
                  type="textarea" 
                  :rows="3" 
                  placeholder="一句话介绍自己（如：欢迎光临 (￣▽￣)~*）" 
                />
                <div class="form-tip">显示在侧边栏名片个人头像下方，支持换行或表情符号。可在【名言语录】页一键将喜欢的名言设为个性签名。</div>
              </el-form-item>

              <el-form-item label="站长头像">
                <ResourceImagePicker 
                  v-model="configForm.author_avatar" 
                  label="站长头像"
                  upload-type="avatar"
                  :preset-options="avatarPresets"
                  tip="博主个人名片及文章底部的头像，支持正方形 JPG/PNG/WebP。"
                />
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <!-- 首页 Hero -->
          <el-tab-pane label="首页设置" name="hero">
            <div class="hero-config-layout" v-if="configForm.hero">
              <el-form label-width="130px" class="config-form hero-form">
                <el-form-item label="主标题">
                  <el-input v-model="configForm.hero.title" placeholder="如：草木山石<br>日月星辰" />
                  <div class="form-tip">支持 HTML 换行标签 &lt;br&gt;</div>
                </el-form-item>
                
                <el-form-item label="副标题">
                  <el-input v-model="configForm.hero.subtitle" placeholder="如：yaan's blog" />
                </el-form-item>

                <el-form-item label="欢迎语">
                  <el-input v-model="configForm.hero.welcome" type="textarea" :rows="2" placeholder="如：Welcome to<br>Yaan's Blog" />
                </el-form-item>

                <el-form-item label="欢迎卡片背景图">
                  <ResourceImagePicker 
                    v-model="configForm.hero.welcome_image" 
                    label="欢迎大图"
                    upload-type="system"
                    aspect-ratio="wide"
                    :preset-options="heroBgPresets"
                    tip="首页右侧欢迎卡片的大背景图，建议 16:9 比例高质感壁纸。"
                  />
                </el-form-item>

                <el-form-item label="装饰技能标签">
                  <div class="skills-editor">
                    <div class="skills-tags">
                      <el-tag
                        v-for="(skill, idx) in heroSkillsList"
                        :key="idx"
                        closable
                        class="skill-tag"
                        @close="removeHeroSkill(idx)"
                      >
                        {{ skill }}
                      </el-tag>
                    </div>

                    <div class="skill-input-row">
                      <el-input
                        v-model="newSkillInput"
                        placeholder="输入技术/标签按回车添加"
                        size="small"
                        class="skill-input"
                        @keyup.enter="addHeroSkill"
                      />
                      <el-button size="small" type="primary" plain @click="addHeroSkill">
                        添加标签
                      </el-button>
                    </div>

                    <div class="quick-skills">
                      <span class="quick-title">推荐常用标签：</span>
                      <el-tag
                        v-for="s in recommendedSkills"
                        :key="s"
                        size="small"
                        class="quick-skill-tag"
                        effect="plain"
                        @click="addRecommendedSkill(s)"
                      >
                        + {{ s }}
                      </el-tag>
                    </div>
                    <div class="form-tip">显示在首页左侧卡片右侧的动态漂浮装饰徽章</div>
                  </div>
                </el-form-item>
              </el-form>

              <!-- 首页 Hero 实时预览区 -->
              <div class="hero-preview-wrap">
                <div class="preview-header">
                  <span class="preview-title">首页卡片实时视觉预览</span>
                  <el-tag size="small" type="success">所见即所得</el-tag>
                </div>

                <div class="hero-preview-container">
                  <div class="hero-grid-sim">
                    <!-- 左侧介绍卡片 -->
                    <div class="sim-intro-card">
                      <div class="sim-intro-content">
                        <div class="sim-title" v-html="configForm.hero.title || '草木山石<br>日月星辰'"></div>
                        <div class="sim-subtitle">{{ configForm.hero.subtitle || "yaan's blog" }}</div>
                      </div>
                      <div class="sim-stack-visual">
                        <div 
                          v-for="(skill, idx) in heroSkillsList.slice(0, 6)" 
                          :key="idx" 
                          class="sim-tech-badge"
                          :class="`badge-${idx % 6}`"
                        >
                          {{ skill }}
                        </div>
                      </div>
                    </div>

                    <!-- 右侧欢迎卡片 -->
                    <div 
                      class="sim-welcome-card"
                      :style="{ backgroundImage: `url(${configForm.hero.welcome_image || '/uploads/defaults/hero.jpg'})` }"
                    >
                      <div class="sim-welcome-overlay"></div>
                      <div class="sim-welcome-content">
                        <div class="sim-welcome-title" v-html="configForm.hero.welcome || 'Welcome to<br>Yaan\'s Blog'"></div>
                        <div class="sim-btn">随便逛逛</div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </el-tab-pane>

          <!-- 名言语录 / 简言库 -->
          <el-tab-pane label="名言语录与简言库" name="quotes">
            <div class="quotes-page-layout">
              <!-- 左侧：语录管理列表 -->
              <div class="quotes-manage-col">
                <div class="quotes-actions-bar">
                  <div class="left-btns">
                    <el-button type="primary" size="small" :icon="Plus" @click="addQuote">
                      添加语录
                    </el-button>
                    <el-button type="success" plain size="small" :icon="Document" @click="openBatchQuotes">
                      批量导入 / 编辑 ({{ configForm.quotes?.length || 0 }}条)
                    </el-button>
                    <el-button type="warning" plain size="small" :icon="MagicStick" @click="quotesInspirationVisible = true">
                      灵感语录推荐
                    </el-button>
                  </div>
                </div>

                <div class="quotes-list-container" v-if="configForm.quotes && configForm.quotes.length > 0">
                  <div 
                    v-for="(quote, index) in configForm.quotes" 
                    :key="index" 
                    class="quote-card-item"
                  >
                    <div class="quote-index-badge">{{ index + 1 }}</div>
                    
                    <div class="quote-input-box">
                      <el-input 
                        v-model="configForm.quotes[index]" 
                        placeholder="请输入名言或感悟句子..." 
                        type="textarea"
                        :rows="2"
                        autosize
                      />
                    </div>

                    <div class="quote-item-ops">
                      <el-tooltip content="上移调整次序" placement="top">
                        <el-button 
                          :icon="Top" 
                          circle 
                          size="small" 
                          :disabled="index === 0" 
                          @click="moveQuote(index, -1)" 
                        />
                      </el-tooltip>
                      <el-tooltip content="下移调整次序" placement="top">
                        <el-button 
                          :icon="Bottom" 
                          circle 
                          size="small" 
                          :disabled="index === configForm.quotes.length - 1" 
                          @click="moveQuote(index, 1)" 
                        />
                      </el-tooltip>
                      <el-tooltip content="一键设为博主个性签名" placement="top">
                        <el-button 
                          :icon="User" 
                          type="success"
                          circle 
                          size="small" 
                          @click="setQuoteAsBio(quote)" 
                        />
                      </el-tooltip>
                      <el-tooltip content="删除此条语录" placement="top">
                        <el-button 
                          type="danger" 
                          :icon="Delete" 
                          circle 
                          size="small" 
                          @click="removeQuote(index)" 
                        />
                      </el-tooltip>
                    </div>
                  </div>
                </div>

                <el-empty v-else description="暂无语录，点击上方【添加语录】或从【灵感推荐】快速填充" />
              </div>

              <!-- 右侧：前台 ProfileCard 拟真名片与气泡实时预览 -->
              <div class="quotes-preview-col">
                <div class="preview-header">
                  <span class="preview-title">博主名片气泡与轮播实时预览</span>
                  <el-tag size="small" type="primary">悬停头像看气泡</el-tag>
                </div>

                <div class="mock-profile-card">
                  <div class="mock-header">
                    <div 
                      class="mock-avatar-wrap"
                      @mouseenter="hoverBubbleVisible = true"
                      @mouseleave="hoverBubbleVisible = false"
                    >
                      <img 
                        :src="configForm.author_avatar || '/uploads/defaults/avatar.jpg'" 
                        class="mock-avatar" 
                        alt="Avatar" 
                      />
                      <div class="mock-status-dot"></div>

                      <!-- 悬停气泡对话框 -->
                      <transition name="pop">
                        <div class="mock-speech-bubble" v-if="hoverBubbleVisible">
                          {{ currentHoverQuote }}
                        </div>
                      </transition>
                    </div>
                  </div>

                  <div class="mock-body">
                    <div class="mock-name">{{ configForm.author_name || 'Yaan' }}</div>
                    <div class="mock-bio">{{ configForm.author_bio || '欢迎光临 (￣▽￣)~*' }}</div>

                    <div class="mock-stats">
                      <div class="stat"><span class="val">12</span><span class="lbl">文章</span></div>
                      <div class="stat"><span class="val">5</span><span class="lbl">分类</span></div>
                    </div>

                    <!-- 名言轮播区域模拟 -->
                    <div class="mock-quote-area">
                      <transition name="fade" mode="out-in">
                        <p class="mock-quote-text" :key="currentCarouselQuote">
                          "{{ currentCarouselQuote }}"
                        </p>
                      </transition>
                    </div>
                    <div class="mock-hint">（此卡片即时呈现前台名片动画效果）</div>
                  </div>
                </div>
              </div>
            </div>
          </el-tab-pane>

          <!-- 快捷链接 -->
          <el-tab-pane label="快捷链接" name="shortcuts">
            <div class="shortcuts-editor">
              <el-button type="primary" plain size="small" @click="addShortcut" class="mb-4">添加链接</el-button>
              <el-table :data="configForm.shortcuts" style="width: 100%" border>
                <el-table-column label="名称" width="160">
                  <template #default="{ row }">
                    <el-input v-model="row.name" size="small" />
                  </template>
                </el-table-column>
                <el-table-column label="链接">
                  <template #default="{ row }">
                    <el-input v-model="row.url" size="small" />
                  </template>
                </el-table-column>
                <el-table-column label="图标 (Iconfont 类名)" width="180">
                  <template #default="{ row }">
                    <el-input v-model="row.icon" size="small" placeholder="如 icon-document" />
                  </template>
                </el-table-column>
                <el-table-column label="渐变或颜色代码">
                  <template #default="{ row }">
                    <el-input v-model="row.color" size="small" placeholder="linear-gradient(...)" />
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
            <el-form label-width="130px" class="config-form" v-if="configForm.page_title">
              <el-form-item label="默认标题">
                <el-input v-model="configForm.page_title.default" placeholder="浏览器标签栏常驻标题" />
              </el-form-item>
              <el-form-item label="失焦标题">
                <el-input v-model="configForm.page_title.blur" placeholder="离开网页时的趣味标题，如 (- - )" />
                <div class="form-tip">用户切换至其他标签页时触发的搞怪标题</div>
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <!-- 评论系统 -->
          <el-tab-pane label="评论系统" name="comment">
            <el-form label-width="130px" class="config-form" v-if="configForm.comment">
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
                  <el-divider content-position="left">Giscus 参数设置</el-divider>
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
          </el-tab-pane>

          <!-- 页脚设置 -->
          <el-tab-pane label="页脚设置" name="footer">
            <el-form label-width="130px" class="config-form" v-if="configForm.footer">
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
                <el-form-item label="项目列表">
                  <el-button type="primary" plain size="small" @click="addPortfolioItem">添加项目</el-button>
                  <el-table :data="configForm.footer.portfolio.items" border style="margin-top: 10px;">
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
                  <el-table :data="configForm.footer.related_links.items" border style="margin-top: 10px;">
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
          </el-tab-pane>

          <!-- 社交链接 -->
          <el-tab-pane label="社交媒体" name="socials">
            <div class="list-editor">
              <el-button type="primary" plain size="small" @click="addSocial" class="mb-4">添加社交链接</el-button>
              <el-table :data="configForm.socials" border>
                <el-table-column label="平台名称" width="140">
                  <template #default="{ row }"><el-input v-model="row.name" size="small" /></template>
                </el-table-column>
                <el-table-column label="链接 URL">
                  <template #default="{ row }"><el-input v-model="row.url" size="small" /></template>
                </el-table-column>
                <el-table-column label="图标 (Iconfont)" width="180">
                  <template #default="{ row }"><el-input v-model="row.icon" size="small" /></template>
                </el-table-column>
                <el-table-column label="主题颜色" width="130">
                  <template #default="{ row }"><el-input v-model="row.color" size="small" /></template>
                </el-table-column>
                <el-table-column label="圆底" width="70" align="center">
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

          <!-- 联系方式卡片 -->
          <el-tab-pane label="侧边栏联系方式" name="contacts">
            <el-form label-width="130px" class="config-form" v-if="configForm.contacts">
              <el-form-item label="显示联系方式">
                <el-switch v-model="configForm.contacts.show" />
              </el-form-item>
              <el-divider content-position="left">联系项目</el-divider>
              <el-button type="primary" plain size="small" @click="addContactItem" class="mb-4">添加联系项</el-button>
              <el-table :data="configForm.contacts.items" border>
                <el-table-column label="项目名称" width="140">
                  <template #default="{ row }"><el-input v-model="row.name" size="small" /></template>
                </el-table-column>
                <el-table-column label="目标链接">
                  <template #default="{ row }"><el-input v-model="row.url" size="small" /></template>
                </el-table-column>
                <el-table-column label="图标" width="180">
                  <template #default="{ row }"><el-input v-model="row.icon" size="small" /></template>
                </el-table-column>
                <el-table-column label="颜色" width="130">
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
          </el-tab-pane>

          <!-- 音乐播放器 -->
          <el-tab-pane label="音乐播放器" name="music">
            <el-form label-width="130px" class="config-form" v-if="configForm.music_player">
              <el-form-item label="显示播放器">
                <el-switch v-model="configForm.music_player.show" />
              </el-form-item>
              <el-form-item label="播放器 URL" v-if="configForm.music_player.show">
                <el-input v-model="configForm.music_player.url" placeholder="//music.163.com/outchain/player..." />
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <!-- 默认回退图片 -->
          <el-tab-pane label="默认回退图片" name="defaultImages">
            <el-form label-width="130px" class="config-form">
              <div class="form-tip mb-4">文章未设封面或用户头像读取失败时的兜底显示图片</div>
              
              <el-form-item label="默认文章封面">
                <ResourceImagePicker 
                  v-model="configForm.default_images.cover" 
                  label="默认封面"
                  upload-type="system"
                  aspect-ratio="wide"
                  :preset-options="coverPresets"
                  tip="文章无封面时的备选图，留空将使用纯色优雅背景。"
                />
              </el-form-item>

              <el-form-item label="默认作者头像">
                <ResourceImagePicker 
                  v-model="configForm.default_images.avatar" 
                  label="默认头像"
                  upload-type="avatar"
                  :preset-options="avatarPresets"
                  tip="作者头像未设置或加载失败时的全局回退头像。"
                />
              </el-form-item>
            </el-form>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-card>

    <!-- 批量导入/编辑语录弹窗 -->
    <el-dialog
      v-model="batchQuotesVisible"
      title="批量导入 / 编辑名言语录"
      width="680px"
      append-to-body
    >
      <div class="batch-quotes-body">
        <p class="batch-tip">在下方文本框中直接编辑，每行一条语录。可以直接从笔记中全选复制粘贴过来：</p>
        <el-input
          v-model="batchQuotesText"
          type="textarea"
          :rows="12"
          placeholder="月亮想着我的心事，一只猫吃了我的奶酪&#10;草木山石，日月星辰&#10;雾霭山岚，风光雨霁"
        />
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="batchQuotesVisible = false">取消</el-button>
          <el-button type="success" plain @click="appendBatchQuotes">追加到现有列表</el-button>
          <el-button type="primary" @click="replaceBatchQuotes">全量替换</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 灵感推荐语录弹窗 -->
    <el-dialog
      v-model="quotesInspirationVisible"
      title="精选语录灵感库"
      width="720px"
      append-to-body
    >
      <div class="inspiration-container">
        <el-tabs tab-position="left">
          <el-tab-pane 
            v-for="cat in quoteCategories" 
            :key="cat.category" 
            :label="cat.category"
          >
            <div class="inspiration-list">
              <div 
                v-for="(item, idx) in cat.items" 
                :key="idx" 
                class="inspiration-item"
              >
                <span class="insp-text">“{{ item }}”</span>
                <el-button 
                  size="small" 
                  type="primary" 
                  plain 
                  :icon="Plus" 
                  @click="addInspirationQuote(item)"
                >
                  添加
                </el-button>
              </div>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
      <template #footer>
        <el-button @click="quotesInspirationVisible = false">完成</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import {
  ElMessage
} from 'element-plus'
import {
  Check,
  Operation,
  Document,
  Plus,
  Delete,
  Top,
  Bottom,
  User,
  MagicStick
} from '@element-plus/icons-vue'
import { systemApi } from '@/services/api'
import yaml from 'js-yaml'
import ResourceImagePicker from '@/components/common/ResourceImagePicker.vue'

defineOptions({
  name: 'ConfigEditor'
})

// --- 预设快捷选项 ---
const logoPresets = [
  { label: '经典 Logo', value: '/static/logo/yaan.png' },
  { label: '系统默认 Logo', value: '/uploads/defaults/logo.png' }
]

const faviconPresets = [
  { label: 'SVG 格式', value: '/favicon.svg' },
  { label: '默认 Favicon', value: '/uploads/defaults/favicon.svg' },
  { label: 'ICO 格式', value: '/uploads/defaults/favicon.ico' }
]

const avatarPresets = [
  { label: '默认头像', value: '/uploads/defaults/avatar.jpg' },
  { label: '备用头像', value: '/static/avatar/avatar.jpg' }
]

const heroBgPresets = [
  { label: '默认星空大图', value: '/uploads/defaults/hero.jpg' },
  { label: '1412 动漫壁纸', value: '/static/img/1412.jpg' }
]

const globalBgPresets = [
  { label: '默认猫咪壁纸 (留空)', value: '' },
  { label: '星空壁纸', value: '/uploads/defaults/hero.jpg' },
  { label: '1412 壁纸', value: '/static/img/1412.jpg' }
]

const coverPresets = [
  { label: '纯色/内置占位 (留空)', value: '' },
  { label: '默认封面', value: '/uploads/defaults/hero.jpg' }
]

// --- 主状态 ---
const mode = ref<'yaml' | 'form'>('form')
const activeTab = ref('basic')
const configContent = ref('')
const saving = ref(false)

const configForm = ref<any>({
  blog_name: '',
  logo_text: '',
  logo_image: '',
  favicon: '',
  background_image: '',
  allowed_hosts: [],
  author_name: '',
  author_bio: '',
  author_avatar: '',
  default_images: { cover: '', avatar: '' },
  hero: {
    title: '',
    subtitle: '',
    welcome: '',
    welcome_image: '',
    skills: []
  },
  quotes: [],
  shortcuts: [],
  page_title: { default: '', blur: '' },
  comment: { enable: false, type: 'giscus', giscus: {} },
  footer: { icp: {}, portfolio: { items: [] }, related_links: { items: [] } },
  socials: [],
  contacts: { show: false, items: [] },
  music_player: { show: false, url: '' }
})

// --- 首页技能标签 ---
const newSkillInput = ref('')
const recommendedSkills = ['Go', 'Vue', 'React', 'TS', 'Docker', 'Gin', 'Linux', 'Python', 'Rust', 'K8s']

const heroSkillsList = computed(() => {
  if (!configForm.value.hero) return []
  if (!Array.isArray(configForm.value.hero.skills)) {
    configForm.value.hero.skills = ['JS', 'Vue', 'React', 'HTML5', 'CSS3', 'TS']
  }
  return configForm.value.hero.skills
})

const addHeroSkill = () => {
  const val = newSkillInput.value.trim()
  if (!val) return
  if (!configForm.value.hero.skills.includes(val)) {
    configForm.value.hero.skills.push(val)
  }
  newSkillInput.value = ''
}

const addRecommendedSkill = (skill: string) => {
  if (!configForm.value.hero.skills.includes(skill)) {
    configForm.value.hero.skills.push(skill)
  }
}

const removeHeroSkill = (idx: number) => {
  configForm.value.hero.skills.splice(idx, 1)
}

// --- 名言语录管理与实时模拟 ---
const batchQuotesVisible = ref(false)
const batchQuotesText = ref('')
const quotesInspirationVisible = ref(false)
const hoverBubbleVisible = ref(false)
const carouselIndex = ref(0)
let carouselTimer: any = null

const quoteCategories = [
  {
    category: '诗词雅韵',
    items: [
      '月亮想着我的心事，一只猫吃了我的奶酪',
      '草木山石，日月星辰',
      '雾霭山岚，风光雨霁',
      '山高水长，江湖未远',
      '醉后不知天在水，满船清梦压星河',
      '星垂平野阔，月涌大江流',
      '行到水穷处，坐看云起时'
    ]
  },
  {
    category: '极客思考',
    items: [
      'Talk is cheap. Show me the code.',
      'Stay hungry, stay foolish.',
      '代码即诗篇，逻辑即宇宙。',
      'Simple is better than complex.',
      'Hello World, Hello Universe.',
      '凡是过往，皆为序章。'
    ]
  },
  {
    category: '生活随想',
    items: [
      '慢品人间烟火色，闲观万事岁月长。',
      '心有山海，静而不争。',
      '生活明朗，万物可爱。',
      '慢慢走，沿途皆有风景。',
      '保持热爱，奔赴山海。'
    ]
  }
]

const currentCarouselQuote = computed(() => {
  const list = configForm.value.quotes
  if (!Array.isArray(list) || list.length === 0) return '草木山石，日月星辰'
  return list[carouselIndex.value % list.length] || list[0]
})

const currentHoverQuote = computed(() => {
  const list = configForm.value.quotes
  if (!Array.isArray(list) || list.length === 0) return '欢迎光临我的博客！'
  return list[Math.floor(Math.random() * list.length)] || list[0]
})

const addQuote = () => {
  if (!Array.isArray(configForm.value.quotes)) configForm.value.quotes = []
  configForm.value.quotes.push('')
}

const removeQuote = (index: number) => {
  configForm.value.quotes.splice(index, 1)
}

const moveQuote = (index: number, step: number) => {
  const target = index + step
  if (target < 0 || target >= configForm.value.quotes.length) return
  const temp = configForm.value.quotes[index]
  configForm.value.quotes[index] = configForm.value.quotes[target]
  configForm.value.quotes[target] = temp
}

const setQuoteAsBio = (text: string) => {
  if (!text) return
  configForm.value.author_bio = text
  ElMessage.success('已将此语录设为博主个性签名！')
}

const openBatchQuotes = () => {
  batchQuotesText.value = (configForm.value.quotes || []).join('\n')
  batchQuotesVisible.value = true
}

const replaceBatchQuotes = () => {
  const lines = batchQuotesText.value
    .split('\n')
    .map(s => s.trim())
    .filter(s => s.length > 0)
  configForm.value.quotes = lines
  batchQuotesVisible.value = false
  ElMessage.success(`已全量更新为 ${lines.length} 条名言语录`)
}

const appendBatchQuotes = () => {
  const lines = batchQuotesText.value
    .split('\n')
    .map(s => s.trim())
    .filter(s => s.length > 0)
  if (!Array.isArray(configForm.value.quotes)) configForm.value.quotes = []
  configForm.value.quotes.push(...lines)
  batchQuotesVisible.value = false
  ElMessage.success(`已追加 ${lines.length} 条名言语录`)
}

const addInspirationQuote = (item: string) => {
  if (!Array.isArray(configForm.value.quotes)) configForm.value.quotes = []
  if (!configForm.value.quotes.includes(item)) {
    configForm.value.quotes.push(item)
    ElMessage.success('已加入语录列表')
  } else {
    ElMessage.info('语录列表中已存在该条内容')
  }
}

// --- 快捷链接与页脚操作 ---
const addShortcut = () => {
  if (!configForm.value.shortcuts) configForm.value.shortcuts = []
  configForm.value.shortcuts.push({ name: '新链接', url: '#', icon: '', color: '' })
}
const removeShortcut = (index: number) => configForm.value.shortcuts.splice(index, 1)

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

const addSocial = () => {
  if (!Array.isArray(configForm.value.socials)) configForm.value.socials = []
  configForm.value.socials.push({ name: '', url: '', icon: '', color: '', is_circle: true })
}
const removeSocial = (i: number) => configForm.value.socials?.splice(i, 1)

const addContactItem = () => {
  if (!configForm.value.contacts) configForm.value.contacts = { show: false, items: [] }
  if (!configForm.value.contacts.items) configForm.value.contacts.items = []
  configForm.value.contacts.items.push({ name: '', url: '', icon: '', color: '', is_circle: true })
}
const removeContactItem = (i: number) => configForm.value.contacts?.items?.splice(i, 1)

// --- 加载与保存 ---
const loadConfig = async () => {
  try {
    const res = await systemApi.getFrontEndConfig()
    if (res.data.status === 200) {
      configContent.value = res.data.data
      try {
        const parsed = yaml.load(configContent.value) as any || {}
        if (!Array.isArray(parsed.allowed_hosts)) parsed.allowed_hosts = []
        if (!parsed.default_images) parsed.default_images = { cover: '', avatar: '' }
        if (!parsed.hero) parsed.hero = { title: '', subtitle: '', welcome: '', welcome_image: '', skills: [] }
        if (!Array.isArray(parsed.hero.skills)) parsed.hero.skills = ['JS', 'Vue', 'React', 'HTML5', 'CSS3', 'TS']
        if (!Array.isArray(parsed.quotes)) parsed.quotes = []
        if (!parsed.background_image) parsed.background_image = ''
        configForm.value = parsed
      } catch (e) {
        console.error('Initial YAML parse error', e)
        mode.value = 'yaml'
        ElMessage.warning('配置文件格式复杂或有误，已自动切至源码模式')
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
    try {
      const parsed = yaml.load(configContent.value) as any
      if (typeof parsed !== 'object' || parsed === null) {
        throw new Error('YAML 格式有误')
      }
      configForm.value = parsed
      mode.value = 'form'
    } catch (e: any) {
      ElMessage.error('YAML 语法错误，无法切换到可视化模式: ' + e.message)
    }
  } else {
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
      ElMessage.success('配置保存成功！前台页面已即时热更新生效')
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

onMounted(() => {
  loadConfig()
  carouselTimer = setInterval(() => {
    carouselIndex.value++
  }, 3500)
})

onUnmounted(() => {
  if (carouselTimer) clearInterval(carouselTimer)
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
  gap: 12px;
}

.card-title {
  font-size: 17px;
  font-weight: 600;
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
  max-width: 820px;
  margin-top: 15px;
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

/* 首页 Hero 配置布局与预览 */
.hero-config-layout {
  display: flex;
  gap: 30px;
  align-items: flex-start;
}

.hero-form {
  flex: 1;
  max-width: 580px;
}

.skills-editor {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.skills-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.skill-tag {
  user-select: none;
}

.skill-input-row {
  display: flex;
  gap: 8px;
  max-width: 320px;
}

.quick-skills {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.quick-title {
  font-size: 12px;
  color: #909399;
}

.quick-skill-tag {
  cursor: pointer;
  transition: all 0.2s;
}

.quick-skill-tag:hover {
  color: var(--el-color-primary);
  border-color: var(--el-color-primary);
}

.hero-preview-wrap {
  width: 420px;
  background: #fbfbfb;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
}

.preview-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.hero-preview-container {
  width: 100%;
}

.hero-grid-sim {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.sim-intro-card {
  height: 160px;
  background: #1a1a1a;
  border-radius: 12px;
  padding: 16px;
  color: #fff;
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
}

.sim-intro-content {
  z-index: 2;
}

.sim-title {
  font-size: 20px;
  font-weight: 800;
  line-height: 1.3;
  margin-bottom: 6px;
}

.sim-subtitle {
  font-size: 12px;
  opacity: 0.6;
  font-family: monospace;
}

.sim-stack-visual {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 55%;
  overflow: hidden;
  opacity: 0.85;
}

.sim-tech-badge {
  position: absolute;
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: bold;
  color: #fff;
  transform: rotate(-10deg);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
}

.badge-0 { background: #f7df1e; color: #323330; top: 15%; right: 40%; }
.badge-1 { background: #409eff; top: 50%; right: 20%; }
.badge-2 { background: #61dafb; color: #20232a; top: 10%; right: 10%; }
.badge-3 { background: #e34f26; bottom: 15%; right: 45%; }
.badge-4 { background: #1572b6; bottom: 20%; right: 8%; }
.badge-5 { background: #3178c6; top: 40%; right: 55%; }

.sim-welcome-card {
  height: 140px;
  border-radius: 12px;
  background-size: cover;
  background-position: center;
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: flex-end;
  padding: 16px;
  color: #fff;
}

.sim-welcome-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(to top, rgba(0,0,0,0.85) 0%, rgba(0,0,0,0.1) 70%);
}

.sim-welcome-content {
  position: relative;
  z-index: 2;
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  width: 100%;
}

.sim-welcome-title {
  font-size: 15px;
  font-weight: 700;
  line-height: 1.3;
}

.sim-btn {
  background: rgba(255, 255, 255, 0.25);
  backdrop-filter: blur(4px);
  border: 1px solid rgba(255, 255, 255, 0.4);
  padding: 4px 10px;
  border-radius: 14px;
  font-size: 11px;
}

/* 名言语录管理布局 */
.quotes-page-layout {
  display: flex;
  gap: 30px;
  align-items: flex-start;
  margin-top: 15px;
}

.quotes-manage-col {
  flex: 1;
}

.quotes-actions-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.left-btns {
  display: flex;
  gap: 10px;
}

.quotes-list-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-height: 650px;
  overflow-y: auto;
  padding-right: 6px;
}

.quote-card-item {
  display: flex;
  align-items: center;
  gap: 12px;
  background: #fbfbfb;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  padding: 10px 14px;
  transition: all 0.2s ease;
}

.quote-card-item:hover {
  background: #fff;
  border-color: var(--el-color-primary-light-5);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.quote-index-badge {
  font-size: 13px;
  font-weight: bold;
  color: #909399;
  width: 24px;
  text-align: center;
}

.quote-input-box {
  flex: 1;
}

.quote-item-ops {
  display: flex;
  align-items: center;
  gap: 6px;
}

/* 博客名片拟真预览区 */
.quotes-preview-col {
  width: 330px;
  background: #fbfbfb;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
}

.mock-profile-card {
  background: #fff;
  border-radius: 12px;
  border: 1px solid var(--el-border-color-light);
  overflow: visible;
  text-align: center;
  margin-top: 10px;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.05);
}

.mock-header {
  height: 80px;
  background: linear-gradient(135deg, #409eff 0%, #205bb5 100%);
  border-radius: 12px 12px 0 0;
  position: relative;
  margin-bottom: 35px;
}

.mock-avatar-wrap {
  width: 70px;
  height: 70px;
  position: absolute;
  bottom: -35px;
  left: 50%;
  transform: translateX(-50%);
  cursor: pointer;
}

.mock-avatar {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  border: 3px solid #fff;
  object-fit: cover;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  transition: transform 0.3s;
}

.mock-avatar-wrap:hover .mock-avatar {
  transform: scale(0.92);
}

.mock-status-dot {
  position: absolute;
  bottom: 3px;
  right: 3px;
  width: 14px;
  height: 14px;
  background-color: #2ecc71;
  border: 2px solid #fff;
  border-radius: 50%;
}

.mock-speech-bubble {
  position: absolute;
  bottom: 115%;
  left: 50%;
  transform: translateX(-50%);
  background: #fff;
  color: #303133;
  padding: 6px 10px;
  border-radius: 6px;
  font-size: 12px;
  width: max-content;
  max-width: 200px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  border: 1px solid #ebeef5;
  z-index: 20;
  line-height: 1.4;
}

.mock-speech-bubble::before {
  content: '';
  position: absolute;
  bottom: -5px;
  left: 50%;
  transform: translateX(-50%) rotate(45deg);
  width: 8px;
  height: 8px;
  background: #fff;
  border-right: 1px solid #ebeef5;
  border-bottom: 1px solid #ebeef5;
}

.mock-body {
  padding: 0 16px 16px;
}

.mock-name {
  font-size: 17px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 4px;
}

.mock-bio {
  font-size: 12px;
  color: #909399;
  margin-bottom: 12px;
}

.mock-stats {
  display: flex;
  justify-content: space-around;
  padding: 8px 0;
  border-top: 1px solid #f0f0f0;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 12px;
}

.mock-stats .stat {
  display: flex;
  flex-direction: column;
}

.mock-stats .val {
  font-size: 15px;
  font-weight: bold;
  color: #303133;
}

.mock-stats .lbl {
  font-size: 11px;
  color: #909399;
}

.mock-quote-area {
  background: #f8f9fa;
  border-left: 3px solid #409eff;
  border-radius: 6px;
  padding: 8px 10px;
  min-height: 52px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.mock-quote-text {
  font-size: 12px;
  color: #606266;
  font-style: italic;
  margin: 0;
  line-height: 1.4;
}

.mock-hint {
  font-size: 11px;
  color: #c0c4cc;
  margin-top: 8px;
}

/* 灵感弹窗 */
.inspiration-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-height: 380px;
  overflow-y: auto;
  padding-left: 10px;
}

.inspiration-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #f8f9fa;
  border-radius: 6px;
  transition: all 0.2s;
}

.inspiration-item:hover {
  background: #f0f2f5;
}

.insp-text {
  font-size: 13px;
  color: #303133;
  flex: 1;
  margin-right: 12px;
}

.batch-tip {
  font-size: 13px;
  color: #606266;
  margin-bottom: 10px;
}

.sub-config {
  margin-top: 15px;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 6px;
}

/* 动画 */
.pop-enter-active,
.pop-leave-active {
  transition: all 0.25s cubic-bezier(0.68, -0.55, 0.265, 1.55);
}

.pop-enter-from,
.pop-leave-to {
  opacity: 0;
  transform: translate(-50%, 6px) scale(0.85);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.35s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
