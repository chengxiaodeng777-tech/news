<script lang="ts" setup>
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();

onMounted(() => {
  router.replace('/dashboard/monitor');
});
</script>

<template>
  <div class="p-5 text-sm text-gray-500">正在跳转到系统监控首页...</div>
</template>
