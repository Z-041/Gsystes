import type { Directive, DirectiveBinding } from 'vue'
import { useAuthStore } from '@/stores/auth'

export const permission: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding<string>) {
    const authStore = useAuthStore()
    const code = binding.value
    if (code && !authStore.hasPermission(code)) {
      el.style.display = 'none'
    }
  },
}
