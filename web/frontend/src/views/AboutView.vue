<template>
  <div class="about-page">
    <div class="container">
      <div class="content-wrapper">
        <div class="main-content">
          <div class="about-content">
            <h1>关于我们</h1>
            <p>这是一个使用 Vue 3 和 Go 语言构建的博客系统。</p>
            <p>该系统提供了完整的博客功能，包括文章发布、分类管理、用户管理等。</p>
            <h2>技术栈</h2>
            <ul>
              <li>前端：Vue 3, TypeScript, Element Plus</li>
              <li>后端：Go, Gin 框架</li>
              <li>数据库：MySQL</li>
              <li>部署：Docker</li>
            </ul>
            <h2>功能特性</h2>
            <ul>
              <li>响应式设计，支持各种设备</li>
              <li>Markdown 文章编辑</li>
              <li>文章分类管理</li>
              <li>用户权限管理</li>
              <li>文件上传（七牛云）</li>
            </ul>
          </div>
        </div>
        
        <div class="sidebar">
          <div class="sidebar-card">
            <div class="card-header">
              <h3>相关信息</h3>
            </div>
            <div class="card-content">
              <p>这是关于页面的侧边栏。</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.about-page {
  width: 100%;
  min-height: calc(100vh - 200px);
}

.content-wrapper {
  display: flex;
  gap: 30px;
  min-height: calc(100vh - 280px);
}

.main-content {
  flex: 1;
  min-width: 0;
}

.about-content {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 40px;
}

.about-content h1 {
  font-size: 32px;
  margin-bottom: 20px;
  color: #333;
}

.about-content h2 {
  font-size: 24px;
  margin: 30px 0 15px;
  color: #333;
}

.about-content p {
  font-size: 16px;
  line-height: 1.6;
  color: #666;
  margin-bottom: 15px;
}

.about-content ul {
  padding-left: 20px;
  margin-bottom: 20px;
}

.about-content li {
  margin-bottom: 10px;
  color: #666;
}

.sidebar {
  width: 300px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.sidebar-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.card-header {
  padding: 15px 20px;
  border-bottom: 1px solid #eee;
  background: #f8f9fa;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.card-content {
  padding: 20px;
}

@media (max-width: 992px) {
  .content-wrapper {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
  }
  
  .about-content {
    padding: 20px;
  }
  
  .about-content h1 {
    font-size: 28px;
  }
  
  .about-content h2 {
    font-size: 22px;
  }
}
</style>