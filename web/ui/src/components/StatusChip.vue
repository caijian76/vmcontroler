<template>
  <v-chip
    size="x-small"
    :color="statusColor"
    variant="flat"
    rounded="lg"
  >
    <v-icon start :icon="statusIcon" size="12" />
    {{ statusLabel }}
  </v-chip>
</template>

<script setup>
import { computed } from 'vue';

const statusConfig = {
  Running:  { color: 'success', icon: 'mdi-play-circle', label: '运行中' },
  Stopped:  { color: 'warning', icon: 'mdi-stop-circle', label: '已停止' },
  Starting: { color: 'info',    icon: 'mdi-loading',     label: '启动中' },
  Stopping: { color: 'info',    icon: 'mdi-loading',     label: '停止中' },
  Pending:  { color: 'info',    icon: 'mdi-hourglass',   label: '等待中' },
};

const props = defineProps({ status: { type: String, default: '' } });

const statusColor = computed(() => statusConfig[props.status]?.color || 'grey');
const statusIcon = computed(() => statusConfig[props.status]?.icon || 'mdi-help-circle');
const statusLabel = computed(() => statusConfig[props.status]?.label || props.status || '未知');
</script>
