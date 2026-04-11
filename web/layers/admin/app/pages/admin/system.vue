<template>
  <AdminPageContainer>
    <AdminPageHeader
      :title="$t('admin.system.title')"
      :subtitle="$t('admin.system.subtitle')">
      <template #actions>
        <UButton
          icon="i-tabler-refresh"
          color="neutral"
          variant="soft"
          :loading="refreshing"
          @click="handleRefresh">
          {{ $t('admin.system.refresh') }}
        </UButton>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <!-- Top stat cards -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 md:gap-4">
        <StatCard
          icon="i-tabler-brand-golang"
          icon-color="text-primary"
          icon-bg="bg-primary/10"
          :total="info?.go_version ?? '-'"
          :label="$t('admin.system.go_version')" />
        <StatCard
          icon="i-tabler-cpu"
          icon-color="text-success"
          icon-bg="bg-success/10"
          :total="info?.num_cpu ?? '-'"
          :label="$t('admin.system.cpu_cores')" />
        <StatCard
          icon="i-tabler-arrows-split"
          icon-color="text-warning"
          icon-bg="bg-warning/10"
          :total="info?.goroutines ?? '-'"
          :label="$t('admin.system.goroutines')" />
        <StatCard
          icon="i-tabler-clock"
          icon-color="text-info"
          icon-bg="bg-info/10"
          :total="info?.uptime ?? '-'"
          :label="$t('admin.system.uptime')" />
      </div>

      <!-- Detail cards -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 md:gap-6">
        <!-- Memory usage -->
        <UCard>
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-tabler-chart-bar" class="size-5 text-primary" />
              <span class="font-semibold text-highlighted">{{ $t('admin.system.memory') }}</span>
            </div>
          </template>
          <div class="space-y-4">
            <InfoRow :label="$t('admin.system.heap_alloc')" :value="info?.memory?.alloc_str ?? '-'" />
            <InfoRow :label="$t('admin.system.sys_memory')" :value="info?.memory?.sys_str ?? '-'" />
            <InfoRow :label="$t('admin.system.total_alloc')" :value="info?.memory?.total_alloc_str ?? '-'" />
            <InfoRow :label="$t('admin.system.gc_cycles')" :value="String(info?.memory?.num_gc ?? '-')" />
          </div>
        </UCard>

        <!-- Environment -->
        <UCard>
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-tabler-device-desktop" class="size-5 text-primary" />
              <span class="font-semibold text-highlighted">{{ $t('admin.system.environment') }}</span>
            </div>
          </template>
          <div class="space-y-4">
            <InfoRow :label="$t('admin.system.os')" :value="info?.os ?? '-'" />
            <InfoRow :label="$t('admin.system.arch')" :value="info?.arch ?? '-'" />
            <InfoRow :label="$t('admin.system.go_version')" :value="info?.go_version ?? '-'" />
            <InfoRow :label="$t('admin.system.db_size')" :value="info?.db_size || '-'" />
          </div>
        </UCard>
      </div>
    </AdminPageContent>
  </AdminPageContainer>
</template>

<script setup lang="ts">
definePageMeta({ layout: "admin" });

interface SystemInfo {
  os: string;
  arch: string;
  go_version: string;
  num_cpu: number;
  goroutines: number;
  memory: {
    alloc: number;
    total_alloc: number;
    sys: number;
    num_gc: number;
    alloc_str: string;
    total_alloc_str: string;
    sys_str: string;
  };
  uptime: string;
  uptime_sec: number;
  db_size_bytes: number;
  db_size: string;
}

const { apiFetch } = useApiFetch();

const { data: info, refresh } = useAsyncData("system-info", () =>
  apiFetch<SystemInfo>("/system/info")
);

const refreshing = ref(false);

async function handleRefresh() {
  refreshing.value = true;
  try {
    await refresh();
  } finally {
    refreshing.value = false;
  }
}
</script>

<script lang="ts">
const InfoRow = defineComponent({
  props: {
    label: { type: String, required: true },
    value: { type: String, required: true },
  },
  setup(props) {
    return () =>
      h("div", { class: "flex items-center justify-between" }, [
        h("span", { class: "text-sm text-muted" }, props.label),
        h("span", { class: "text-sm font-medium text-highlighted" }, props.value),
      ]);
  },
});
</script>
