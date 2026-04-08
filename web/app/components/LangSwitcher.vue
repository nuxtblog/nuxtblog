<script setup lang="ts">
const { locale, switchLocale } = useLocaleSwitch()
const { getLanguages } = useSiteApi()

const { data: langsData } = await useAsyncData('site-languages', getLanguages, {
  default: () => ({ list: [] as Array<{ code: string; name: string; label: string }> }),
})

const items = computed(() => [
  langsData.value.list.map(lang => ({
    label: lang.name,
    trailingIcon: locale.value === lang.code ? 'i-tabler-check' : undefined,
    onClick: () => switchLocale(lang.code),
  })),
])

</script>

<template>
  <UDropdownMenu :items="items">
    <UButton color="neutral" variant="ghost" icon="i-tabler-language" />
  </UDropdownMenu>
</template>
