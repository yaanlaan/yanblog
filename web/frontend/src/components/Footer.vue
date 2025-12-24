<template>
  <div class="footer-wrapper">
    <div class="copyright-container">
      <div class="container">
        <div class="ref-link">
          <span>
            <h2>言盐盐 - 无人在意的小角落</h2>
          </span>
          <span class="copyright"> Email: yanxia2425@foxmail.com </span>
          <span class="copyright"> Copyright ©2022-2025 言盐盐 All Rights Reserved </span>
          <span class="copyright">
            <!-- <a href="https://beian.miit.gov.cn" target="_blank" class="copyright-message">备案</a> -->
          </span>
        </div>

        <div class="ref-link">
          <h2><span class="iconfont icon-book" style="font-size:28px;"></span> 作品集</h2>

          <span>
            <span class="iconfont icon-lianjie"></span>
            <a href="https://github.com/yaanlaan/yanblog.git">yanblog</a>
          </span>

          <span>
            <span class="iconfont icon-lianjie"></span>
            <a href="https://github.com/yaanlaan/yanblog.git">yanblog</a>
          </span>

          <span>
            <span class="iconfont icon-lianjie"></span>
            <a href="https://github.com/yaanlaan/yanblog.git">yanblog</a>
          </span>
        </div>

        <div class="ref-link">
          <h2>相关链接</h2>
          <span>
            <span class="iconfont icon-icon-zhihu-circle-fill"></span>
            <a href="/">知乎@言盐盐</a>
          </span>

          <span>
            <span class="iconfont icon-bilibili-line"></span>
            <a href="/">Bilibili@土豆炸马铃薯0_0</a>
          </span>

          <span>
            <span class="iconfont icon-github-fill"></span>
            <a href="https://github.com/yaanlaan">Github@yaanlaan</a>
          </span>

          <span>
            <span class="iconfont icon-lianjie"></span>
            <a href="https://kirigaya.cn">本博客参考博客链接（知乎：恢锦）</a>
          </span>

          <!-- <span>
            <span class="iconfont icon-steam"></span>
            <a href="https://steamcommunity.com">Steam@yaanlaan</a>
          </span> -->
          
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 页脚组件
</script>

<style scoped>
.footer-wrapper {
  width: 100%;
}

.copyright-container {
  width: 100%;
  background-color: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(10px);
  border-top: 1px solid rgba(66, 184, 131, 0.2);
  padding: 30px 0;
  margin-top: 50px;
}

.copyright-container a {
  color: inherit;
  text-decoration: none;
  transition: color 0.3s;
}

.copyright-container a:hover {
  color: #42b883;
}

.container {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  justify-content: space-evenly;
}

.copyright-message {
  color: rgb(127, 136, 143);
}

.copyright-container .title {
  display: flex;
  font-size: 15px;
  align-items: center;
  margin-bottom: 20px;
}

.copyright-container .copyright {
  font-size: 18px;
}

.copyright-container .ref-link {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  font-size: 0.9rem;
}

.copyright-container .ref-link h2 {
  font-size: 32px;
}

.copyright-container .ref-link>span {
  font-size: 18px;
  margin: 7px;
  display: flex;
  align-items: center;
}

.copyright-container .ref-link .iconfont {
  margin-right: 8px;
  font-size: 18px;
}

.copyright-container .ref-link a {
  color: rgb(127, 136, 143);
  text-decoration-line: underline;
  text-decoration-style: solid;
  text-decoration-color: currentColor;
  text-decoration-thickness: 1px;
  text-underline-offset: 0.5em;
  background-repeat: unset;
  background-image: unset;
  background-size: unset;
  background-position: unset;
}

.copyright-container .ref-link a:hover {
  color: #2cb497ff;
  text-decoration-thickness: 3px;
}

@media (max-width: 992px) {
  .container {
    flex-direction: column;
    align-items: center;
  }

  .copyright-container .ref-link {
    align-items: center;
    margin-bottom: 30px;
  }

  .copyright-container .ref-link h2 {
    text-align: center;
  }

  .copyright-container .ref-link>span {
    align-items: center;
    text-align: center;
  }
}
</style>