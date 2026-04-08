import { defineNuxtPlugin } from "#app";
import type { Directive } from "vue";

interface HTMLElementWithHandler extends HTMLElement {
  __clickOutsideHandler__?: (event: MouseEvent) => void;
}

const vClickModalOutside: Directive<HTMLElementWithHandler> = {
  mounted(el, binding) {
    el.__clickOutsideHandler__ = (event: MouseEvent) => {
      // 排除触发按钮
      if ((event.target as HTMLElement).closest(".modal-control")) return;

      // 当前激活状态
      const isActive =
        (el instanceof HTMLDetailsElement && el.open) ||
        (el instanceof HTMLDialogElement && el.open) ||
        (!(el instanceof HTMLDetailsElement) &&
          !(el instanceof HTMLDialogElement));

      if (!isActive) return;

      // dialog 特殊逻辑
      if (el instanceof HTMLDialogElement) {
        const rect = el.getBoundingClientRect();
        const isInDialog =
          event.clientX >= rect.left &&
          event.clientX <= rect.right &&
          event.clientY >= rect.top &&
          event.clientY <= rect.bottom;

        // 点击 backdrop
        if (event.target === el) {
          el.close();
          return;
        }

        // 点击内容区则忽略
        if (isInDialog) return;
      } else {
        // 其他元素点击内部不处理
        if (el.contains(event.target as Node)) return;
      }

      // 关闭元素
      if (el instanceof HTMLDetailsElement) el.open = false;
      else if (el instanceof HTMLDialogElement) el.close();

      // 用户的回调
      const value = binding.value;
      if (value === undefined) return;

      if (typeof value === "function") {
        value(event);
      }
    };

    document.addEventListener("click", el.__clickOutsideHandler__);
  },

  unmounted(el) {
    if (el.__clickOutsideHandler__) {
      document.removeEventListener("click", el.__clickOutsideHandler__);
      delete el.__clickOutsideHandler__;
    }
  },
};

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.directive("click-modal-outside", vClickModalOutside);
});
