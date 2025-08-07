<template>
  <button class="back-to-top" @click="scrollToTop" :class="{ show: isVisible }">
    <i class="iconfont icon-up-btn"></i>
  </button>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'

const isVisible = ref(false)

const handleScroll = () => {
  isVisible.value = window.scrollY > 300
}

const scrollToTop = () => {
  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  })
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
.iconfont {
  font-size: 35px; 
}
.back-to-top {
  position: fixed;
  bottom: 50px;
  right: 10px;
  background-color: #99bfdf;
  color: white;
  border: none;
  border-radius: 50%;
  padding: 10px;
  cursor: pointer;
  transition: opacity 0.3s, transform 0.3s;
  opacity: 0;
  transform: translateY(100%);
  visibility: hidden;
  z-index: 1000;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.back-to-top:hover {
  background-color: #3894e0;
  opacity: 1;
  transform: translateY(0);
}

.back-to-top.show {
  opacity: 1;
  transform: translateY(0);
  visibility: visible;
}

</style>