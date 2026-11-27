<script lang="ts" setup>
import type { AnalysisOverviewItem } from '@vben/common-ui';
import type { MonitorRecord } from '#/api/system';

import { computed, onMounted, ref } from 'vue';

import { AnalysisChartCard, AnalysisOverview, Page } from '@vben/common-ui';
import {
  SvgBellIcon,
  SvgCakeIcon,
  SvgCardIcon,
  SvgDownloadIcon,
} from '@vben/icons';

import { getMonitor } from '#/api/system';

const loading = ref(false);
const monitor = ref<MonitorRecord>();

const overviewItems = computed<AnalysisOverviewItem[]>(() => {
  const runtime = monitor.value?.runtime;
  const database = monitor.value?.database;
  const redis = monitor.value?.redis;
  const host = monitor.value?.host;

  return [
    {
      icon: SvgCardIcon,
      title: '协程数',
      totalTitle: '当前协程数',
      totalValue: runtime?.goroutines ?? 0,
      value: runtime?.cpus ?? 0,
    },
    {
      icon: SvgCakeIcon,
      title: '数据库连接',
      totalTitle: '数据库打开连接',
      totalValue: database?.openConnections ?? 0,
      value: database?.inUseConnections ?? 0,
    },
    {
      icon: SvgDownloadIcon,
      title: 'Redis连接',
      totalTitle: 'Redis总连接数',
      totalValue: redis?.totalConns ?? 0,
      value: redis?.idleConns ?? 0,
    },
    {
      icon: SvgBellIcon,
      title: '内存(MB)',
      totalTitle: '系统总内存',
      totalValue: host?.totalMemoryMB ?? 0,
      value: runtime?.memoryAllocMB ?? 0,
    },
  ];
});

async function loadMonitor() {
  loading.value = true;
  try {
    monitor.value = await getMonitor();
  } finally {
    loading.value = false;
  }
}

onMounted(loadMonitor);
</script>

<template>
  <Page auto-content-height>
    <div class="p-5">
      <AnalysisOverview :items="overviewItems" />

      <div class="mt-5 grid grid-cols-1 gap-4 xl:grid-cols-2">
        <AnalysisChartCard :loading="loading" title="应用概览">
          <div class="space-y-3 text-sm">
            <div class="flex justify-between">
              <span>应用名称</span>
              <span>{{ monitor?.overview.appName }}</span>
            </div>
            <div class="flex justify-between">
              <span>运行环境</span>
              <span>{{ monitor?.overview.environment }}</span>
            </div>
            <div class="flex justify-between">
              <span>服务地址</span>
              <span>{{ monitor?.overview.address }}</span>
            </div>
            <div class="flex justify-between">
              <span>启动时间</span>
              <span>{{ monitor?.overview.startTime }}</span>
            </div>
            <div class="flex justify-between">
              <span>运行时长</span>
              <span>{{ monitor?.overview.uptime }}</span>
            </div>
            <div class="flex justify-between">
              <span>当前时间</span>
              <span>{{ monitor?.overview.currentTime }}</span>
            </div>
          </div>
        </AnalysisChartCard>

        <AnalysisChartCard :loading="loading" title="数据库状态">
          <div class="space-y-3 text-sm">
            <div class="flex justify-between">
              <span>状态</span>
              <span>{{ monitor?.database.status }}</span>
            </div>
            <div class="flex justify-between">
              <span>驱动</span>
              <span>{{ monitor?.database.driver }}</span>
            </div>
            <div class="flex justify-between">
              <span>打开连接</span>
              <span>{{ monitor?.database.openConnections }}</span>
            </div>
            <div class="flex justify-between">
              <span>使用中连接</span>
              <span>{{ monitor?.database.inUseConnections }}</span>
            </div>
            <div class="flex justify-between">
              <span>空闲连接</span>
              <span>{{ monitor?.database.idleConnections }}</span>
            </div>
            <div class="flex justify-between">
              <span>等待次数</span>
              <span>{{ monitor?.database.waitCount }}</span>
            </div>
          </div>
        </AnalysisChartCard>

        <AnalysisChartCard :loading="loading" title="Redis状态">
          <div class="space-y-3 text-sm">
            <div class="flex justify-between">
              <span>状态</span>
              <span>{{ monitor?.redis.status }}</span>
            </div>
            <div class="flex justify-between">
              <span>地址</span>
              <span>{{ monitor?.redis.address }}</span>
            </div>
            <div class="flex justify-between">
              <span>DB</span>
              <span>{{ monitor?.redis.db }}</span>
            </div>
            <div class="flex justify-between">
              <span>连接池大小</span>
              <span>{{ monitor?.redis.poolSize }}</span>
            </div>
            <div class="flex justify-between">
              <span>总连接数</span>
              <span>{{ monitor?.redis.totalConns }}</span>
            </div>
            <div class="flex justify-between">
              <span>Ping耗时</span>
              <span>{{ monitor?.redis.pingDuration }}</span>
            </div>
          </div>
        </AnalysisChartCard>

        <AnalysisChartCard :loading="loading" title="运行时与主机状态">
          <div class="space-y-3 text-sm">
            <div class="flex justify-between">
              <span>主机名</span>
              <span>{{ monitor?.host.hostname }}</span>
            </div>
            <div class="flex justify-between">
              <span>系统/架构</span>
              <span>{{ monitor?.host.os }} / {{ monitor?.host.arch }}</span>
            </div>
            <div class="flex justify-between">
              <span>Go版本</span>
              <span>{{ monitor?.runtime.goVersion }}</span>
            </div>
            <div class="flex justify-between">
              <span>协程数</span>
              <span>{{ monitor?.runtime.goroutines }}</span>
            </div>
            <div class="flex justify-between">
              <span>内存分配(MB)</span>
              <span>{{ monitor?.runtime.memoryAllocMB }}</span>
            </div>
            <div class="flex justify-between">
              <span>磁盘剩余(GB)</span>
              <span>{{ monitor?.host.freeDiskGB }}</span>
            </div>
          </div>
        </AnalysisChartCard>
      </div>
    </div>
  </Page>
</template>
