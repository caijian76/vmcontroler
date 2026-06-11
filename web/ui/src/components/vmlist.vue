<template>
  <div class="page-shell">
    <v-snackbar
      v-model="snackbarVisible"
      :timeout="30000"
      color="error"
      location="top"
      variant="elevated"
      multi-line
    >
      {{ snackbarMessage }}
      <template v-slot:actions>
        <v-btn color="blue" variant="text" @click="snackbarVisible = false"> x </v-btn>
      </template>
    </v-snackbar>

    <section class="panel-card">
      <header class="panel-header">
        <div>
          <p class="eyebrow">KubeVirt</p>
          <h2>虚拟机列表</h2>
          <p class="subtitle">管理实例运行状态、查看就绪情况与启动/停止操作。</p>
        </div>
        <button class="logout-btn" @click="logout">退出登录</button>
      </header>

      <div v-if="isLoading" class="empty-state">加载中...</div>

      <div v-else-if="vmList.length" class="table-wrap">
        <table>
      <thead>
        <tr>
          <th>名称</th>
          <th>基本配置</th>
          <th>就绪状态</th>
          <th>当前状态</th>
          <th>启动时间</th>
          <th>开关</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="vm in vmList" :key="vm.Name">
          <!-- 第一列根据虚拟机状态决定是否显示超链接 -->
          <td>
            <template v-if="vm.Status === 'Running'">
              <a 
                :href="`/novnc/vnc.html?path=/vm/${vm.Name}/vnc&autoconnect=true`" 
                target="_blank"
                rel="noopener noreferrer"
              >
                {{ vm.Name }}
              </a>
            </template>
            <template v-else>
              {{ vm.Name }}
            </template>
          </td>
          
          <td>{{ vm.CPU || '0 vCPU' }} / {{ vm.Memory || '未配置' }}</td>
          <td>{{ vm.Ready ? '就绪' : '未就绪' }}</td>
          <td>{{ vm.Status }}</td>
          <td>{{ vm.Status === 'Running' ? (vm.StartTime || '') : '' }}</td>
          <td>
            <button
              class="power-btn"
              :class="vm.Status === 'Running' ? 'stop' : 'start'"
              type="button"
              :disabled="togglingVmNames.includes(vm.Name)"
              @click="toggleVmStatus(vm)"
            >
              {{ togglingVmNames.includes(vm.Name) ? '处理中...' : (vm.Status === 'Running' ? '停止' : '启动') }}
            </button>
          </td>
        </tr>
      </tbody>
        </table>
      </div>

      <div v-else class="empty-state">暂无虚拟机数据</div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { useCookies } from 'vue3-cookies';
import { useRouter } from 'vue-router';

// 引入 vue3-cookies
const { cookies } = useCookies();
const router = useRouter();

// 定义响应式变量
const vmList = ref([]);
const isLoading = ref(true);
const errorMessage = ref('');
const togglingVmNames = ref([]);
const snackbarVisible = ref(false);
const snackbarMessage = ref('');
// 定义定时器变量
let refreshInterval;

const showToast = (message) => {
  if (!message) {
    return;
  }

  snackbarMessage.value = message;
  snackbarVisible.value = true;
};

const getErrorMessage = async (response, fallbackMessage) => {
  try {
    const data = await response.json();
    if (data && typeof data === 'object') {
      if (data.error) {
        return `${fallbackMessage}：${data.error}`;
      }
      if (data.message) {
        return `${fallbackMessage}：${data.message}`;
      }
    }
  } catch (error) {
    // 忽略 JSON 解析失败，回退到默认提示
  }

  return fallbackMessage;
};

const getJwtToken = () => {
  const jwtToken = cookies.get('jwt_token');
  if (!jwtToken) {
    throw new Error('未找到有效的 JWT Token，请重新登录');
  }
  return jwtToken;
};

// 获取虚拟机列表的函数
const fetchVmList = async () => {
  try {
    const jwtToken = getJwtToken();

    // 发送请求获取虚拟机列表，添加 Authorization 请求头（同源请求，不写死 IP）
    const response = await fetch('/api/vm', {
      headers: {
        Authorization: `${jwtToken}`
      }
    });
    if (!response.ok) {
      const message = await getErrorMessage(response, '请求失败，请稍后重试');
      throw new Error(message);
    }
    // 解析响应数据
    vmList.value = await response.json();
  } catch (error) {
    errorMessage.value = error.message;
    showToast(error.message);
  } finally {
    // 无论请求成功或失败，都结束加载状态
    isLoading.value = false;
  }
};

const toggleVmStatus = async (vm) => {
  if (togglingVmNames.value.includes(vm.Name)) {
    return;
  }

  togglingVmNames.value = [...togglingVmNames.value, vm.Name];
  errorMessage.value = '';

  try {
    const jwtToken = getJwtToken();
    const action = vm.Status === 'Running' ? 'stop' : 'start';
    const response = await fetch(`/api/vm/${encodeURIComponent(vm.Name)}/${action}`, {
      headers: {
        Authorization: `${jwtToken}`
      }
    });

    if (!response.ok) {
      const message = await getErrorMessage(response, '操作失败，请稍后重试');
      console.log('操作失败:', message);
      throw new Error(message);
    }

    await fetchVmList();
  } catch (error) {
    errorMessage.value = error.message;
    showToast(error.message);
  } finally {
    togglingVmNames.value = togglingVmNames.value.filter((name) => name !== vm.Name);
  }
};

// 退出登录函数
const logout = () => {
  // 清除 jwt_token Cookie
  cookies.remove('jwt_token');
  // 跳转回登录页面
  router.push('/'); 
};

// 组件挂载后调用获取虚拟机列表的函数，并设置定时器
onMounted(() => {
  fetchVmList();
  refreshInterval = setInterval(fetchVmList, 15*1000);
});

// 组件卸载时清除定时器
onUnmounted(() => {
  clearInterval(refreshInterval);
});
</script>

<style scoped>
.page-shell {
  min-height: 100vh;
  padding: 24px;
  background:
    radial-gradient(circle at top, rgba(59, 130, 246, 0.10), transparent 30%),
    linear-gradient(135deg, #0f172a 0%, #111827 45%, #172554 100%);
  color: #e5eefb;
}

.panel-card {
  max-width: 1200px;
  margin: 0 auto;
  background: rgba(15, 23, 42, 0.86);
  border: 1px solid rgba(148, 163, 184, 0.18);
  border-radius: 18px;
  box-shadow: 0 18px 30px rgba(15, 23, 42, 0.35);
  padding: 20px;
  backdrop-filter: blur(10px);
}

.panel-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 18px;
}

.eyebrow {
  text-transform: uppercase;
  letter-spacing: 0.18em;
  color: #93c5fd;
  font-size: 12px;
  margin-bottom: 6px;
}

h2 {
  margin: 0 0 6px;
  font-size: 26px;
  color: #fff;
}

.subtitle {
  margin: 0;
  color: #cbd5e1;
  font-size: 14px;
}

.logout-btn,
.power-btn {
  border: none;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
  transition: transform 0.12s ease, opacity 0.2s ease, background 0.2s ease;
}

.logout-btn:hover,
.power-btn:hover:not(:disabled) {
  transform: translateY(-1px);
}

.logout-btn {
  background: linear-gradient(135deg, #2563eb, #1d4ed8);
  color: #fff;
  padding: 10px 14px;
}

.table-wrap {
  overflow-x: auto;
}

.table-wrap table {
  width: 100%;
  border-collapse: collapse;
  background: rgba(15, 23, 42, 0.92);
  border-radius: 12px;
  overflow: hidden;
}

th, td {
  border-bottom: 1px solid rgba(148, 163, 184, 0.16);
  padding: 12px 10px;
  text-align: left;
  font-size: 14px;
}

th {
  background: rgba(30, 41, 59, 0.95);
  color: #bfdbfe;
  font-weight: 700;
}

tr:nth-child(even) {
  background: rgba(30, 41, 59, 0.45);
}

tr:hover {
  background: rgba(51, 65, 85, 0.55);
}

.power-btn {
  min-width: 72px;
  padding: 8px 10px;
  color: #fff;
}

.power-btn.start {
  background: linear-gradient(135deg, #22c55e, #16a34a);
}

.power-btn.stop {
  background: linear-gradient(135deg, #f59e0b, #d97706);
}

.power-btn:disabled {
  opacity: 0.65;
  cursor: not-allowed;
}

.empty-state {
  padding: 28px;
  text-align: center;
  color: #cbd5e1;
  background: rgba(30, 41, 59, 0.5);
  border-radius: 12px;
}

a {
  color: #93c5fd;
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}
</style>
