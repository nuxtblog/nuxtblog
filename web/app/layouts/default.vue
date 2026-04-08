<template>
  <div class="min-h-screen flex flex-col">
    <BlogHeader />
    <main class="grow max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-8">
      <slot />
    </main>
    <BlogFooter />
    <ClientOnly><BlogFloatingToolbar /></ClientOnly>
    <BlogBackToTop />
  </div>
</template>
