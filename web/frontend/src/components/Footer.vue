<template>
  <div class="footer-wrapper">
    <div class="copyright-container">
      <div class="container">
        <div class="ref-link">
          <span>
            <h2>{{ siteInfo.blog_name }} - 无人在意的小角落</h2>
          </span>
          <span class="copyright"> Email: {{ siteInfo.footer.email }} </span>
          <span class="copyright"> {{ siteInfo.footer.copyright }} </span>
          <span class="copyright" v-if="siteInfo.footer.icp?.show">
            <a :href="siteInfo.footer.icp.link" target="_blank" class="copyright-message">{{ siteInfo.footer.icp.text }}</a>
          </span>
        </div>

        <div class="ref-link" v-if="siteInfo.footer.portfolio?.show">
          <h2><span class="iconfont" :class="siteInfo.footer.portfolio.icon" style="font-size:28px;"></span> {{ siteInfo.footer.portfolio.title }}</h2>

          <span v-for="(item, index) in siteInfo.footer.portfolio.items" :key="index">
            <span class="iconfont" :class="item.icon"></span>
            <a :href="item.url" target="_blank">{{ item.name }}</a>
          </span>
        </div>

        <div class="ref-link" v-if="siteInfo.footer.related_links?.show">
          <h2>{{ siteInfo.footer.related_links.title }}</h2>
          <span v-for="(item, index) in siteInfo.footer.related_links.items" :key="index">
            <span class="iconfont" :class="item.icon"></span>
            <a :href="item.url" target="_blank">{{ item.name }}</a>
          </span>
        </div>

          <!-- <span>
            <span class="iconfont icon-steam"></span>
            <a href="https://steamcommunity.com">Steam@yaanlaan</a>
          </span> -->
      </div>
    </div>
  </div>
</template>

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

<script setup lang="ts">
import { useSiteInfoStore } from '@/stores/siteInfo'
import { storeToRefs } from 'pinia'

const siteInfoStore = useSiteInfoStore()
const { siteInfo } = storeToRefs(siteInfoStore)
</script>