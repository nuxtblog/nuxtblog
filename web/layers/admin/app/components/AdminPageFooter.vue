<template>
  <div class="shrink-0 border-t border-default bg-default px-4 md:px-6 py-2.5 md:py-3">
    <div class="flex items-center justify-between gap-4 min-h-8">
      <div class="flex items-center gap-3 text-sm text-muted min-w-0">
        <slot name="left" />
      </div>
      <div class="flex items-center gap-2 shrink-0">
        <slot name="right" />
      </div>
    </div>
  </div>
</template>
