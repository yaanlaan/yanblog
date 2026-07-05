import type { Directive, DirectiveBinding } from 'vue'

interface LazyImageElement extends HTMLImageElement {
  _lazySrc?: string
  _lazyObserver?: IntersectionObserver
}

const lazy: Directive = {
  mounted(el: LazyImageElement, binding: DirectiveBinding) {
    const src = binding.value
    const defaultSrc = el.getAttribute('data-default-src') || '/assets/img/无封面.jpg'
    
    if (!src || src.trim() === '') {
      el.src = defaultSrc
      el.style.opacity = '1'
      return
    }
    
    el._lazySrc = src

    const placeholder = document.createElement('div')
    placeholder.className = 'lazy-placeholder'
    placeholder.style.cssText = `
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: linear-gradient(90deg, var(--color-border) 25%, var(--color-border-hover) 50%, var(--color-border) 75%);
      background-size: 200% 100%;
      animation: lazy-shimmer 1.5s ease-in-out infinite;
      border-radius: inherit;
    `
    el.parentElement?.style.position === 'static' && (el.parentElement.style.position = 'relative')
    el.style.opacity = '0'
    el.style.transition = 'opacity 0.4s ease'

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting && el._lazySrc) {
            const img = new Image()
            img.onload = () => {
              el.src = el._lazySrc!
              placeholder.remove()
              el.style.opacity = '1'
              observer.unobserve(el)
            }
            img.onerror = () => {
              placeholder.remove()
              el.src = defaultSrc
              el.style.opacity = '1'
              observer.unobserve(el)
            }
            img.src = el._lazySrc
          }
        })
      },
      {
        rootMargin: '100px',
        threshold: 0.1
      }
    )

    el._lazyObserver = observer
    observer.observe(el)
  },
  unmounted(el: LazyImageElement) {
    el._lazyObserver?.disconnect()
  }
}

export default lazy