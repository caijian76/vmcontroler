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
          <v-btn class="new-vm-btn" @click="openCreateDialog">新建虚拟机</v-btn>
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
              <th>运行节点 </th>
              <th>就绪状态</th>
              <th>当前状态</th>
              <th>启动时间</th>
              <th>开关</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="vm in vmList" :key="vm.Name">
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
              <td>{{ vm.Status === 'Running' ? (vm.NodeName || '未调度') : '未调度' }}</td>
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

    <!-- 新建虚拟机弹窗 -->
    <v-dialog v-model="dialog" max-width="800" persistent>
      <v-card title="配置新建虚拟机参数">
        <v-card-text>
          <v-form ref="createVmForm" v-model="formValid" lazy-validation>
            <!-- 虚拟机名称 -->
            <v-text-field
              v-model="createForm.name"
              label="虚拟机名称"
              placeholder="仅小写字母、数字、横杠"
              required
              :rules="nameRules"
              outlined
              dense
            ></v-text-field>

            <!-- CPU + 内存 同一行两栏 -->
            <v-row no-gutters>
              <v-col cols="5" pr="4">
                <v-text-field
                  v-model.number="createForm.cpu"
                  label="CPU(vCPU)"
                  type="number"
                  min="1"
                  required
                  :rules="numRules"
                  outlined
                  dense
                ></v-text-field>
              </v-col>
            <v-col cols="2"></v-col>  
  
              <v-col cols="5">
                <v-text-field
                  v-model.number="createForm.memory"
                  label="内存大小(GiB)"
                  type="number"
                  min="1"
                  required
                  :rules="numRules"
                  outlined
                  dense
                ></v-text-field>
              </v-col>
            </v-row>
            <v-row no-gutters>
              <v-col cols="5">
            <!-- 节点选择 -->
            <v-select
              v-model="createForm.node"
              label="在哪个节点创建"
              :items="nodeList"
              required
              outlined
              dense
            ></v-select>
            </v-col>
            <v-col cols="2"></v-col>  
            <v-col cols="5">    
            <!-- 存储磁盘大小 -->
            <v-text-field
              v-model.number="createForm.diskSize"
              label="磁盘存储大小(GiB)"
              type="number"
              min="10"
              required
              :rules="numRules"
              outlined
              dense
            ></v-text-field>
            </v-col>
            </v-row>

            <!-- 网络模式 -->
            <v-select
              v-model="createForm.networkType"
              label="网络模式"
              :items="networkOptions"
              required
              outlined
              dense
            ></v-select>

            <!-- 开机自动启动 -->

            <v-switch
             v-model="createForm.autoStart"
             label="创建后自动开机"
             color="primary"
             thumb-color="orange-darken-2">
            </v-switch>
          </v-form>
        </v-card-text>

        <template v-slot:actions>
          <v-spacer></v-spacer>
          <v-btn @click="closeCreateDialog" color="grey">取消</v-btn>
          <v-btn color="primary" :loading="submitting" @click="submitCreateVm">确认创建</v-btn>
        </template>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { useCookies } from 'vue3-cookies';
import { useRouter } from 'vue-router';

const { cookies } = useCookies();
const router = useRouter();

// 页面基础状态
const vmList = ref([]);
const nodeList = ref([]);
const isLoading = ref(true);
const errorMessage = ref('');
const togglingVmNames = ref([]);
const snackbarVisible = ref(false);
const snackbarMessage = ref('');
const dialog = ref(false);
const submitting = ref(false);
const formValid = ref(false);
const createVmForm = ref(null);


 
// 新建虚拟机表单数据
const createForm = ref({
  name: '',
  cpu: 2,
  memory: 4,
  node: '',
  diskSize: 20,
  networkType: 'bridge',
  autoStart: true
});



// 网络模式选项
const networkOptions = ref([
  { title: '桥接网络 bridge', value: 'bridge' },
  { title: '集群默认Pod网络 pod', value: 'pod' },
  { title: '隔离本地网络 isolated', value: 'isolated' }
]);

// 表单校验规则
const nameRules = [
  v => !!v || '虚拟机名称不能为空',
  v => /^[a-z0-9-]+$/.test(v) || '只能小写字母、数字、横杠'
];
const numRules = [v => v >= 1 || '数值必须大于等于1'];

let refreshInterval;

// 消息提示
const showToast = (message) => {
  if (!message) return;
  snackbarMessage.value = message;
  snackbarVisible.value = true;
};

// 错误解析
const getErrorMessage = async (response, fallbackMessage) => {
  try {
    const data = await response.json();
    if (data && typeof data === 'object') {
      if (data.error) return `${fallbackMessage}：${data.error}`;
      if (data.message) return `${fallbackMessage}：${data.message}`;
    }
  } catch (err) {}
  return fallbackMessage;
};

// 拉取节点列表
const fetchNodeList = async () => {
  try {
    const jwtToken = getJwtToken();
    const response = await fetch('/api/nodes', {
      headers: { Authorization: `${jwtToken}` }
    });
    if (!response.ok) {
      const msg = await getErrorMessage(response, '获取节点列表失败');
      throw new Error(msg);
    }
    // 后端返回数组如 ["node1","node2"]
    const nodes = await response.json();
    nodeList.value = Array.isArray(nodes) ? nodes : [];
  } catch (error) {
    showToast(error.message);
    nodeList.value = [];
  }
};

// 获取Token
const getJwtToken = () => {
  const jwtToken = cookies.get('jwt_token');
  if (!jwtToken) throw new Error('未找到有效的 JWT Token，请重新登录');
  return jwtToken;
};

// 拉取虚拟机列表
const fetchVmList = async () => {
  try {
    const jwtToken = getJwtToken();
    const response = await fetch('/api/vm', {
      headers: { Authorization: `${jwtToken}` }
    });
    if (!response.ok) {
      const msg = await getErrorMessage(response, '请求失败，请稍后重试');
      throw new Error(msg);
    }
    vmList.value = await response.json();
  } catch (error) {
    errorMessage.value = error.message;
    showToast(error.message);
  } finally {
    isLoading.value = false;
  }
};

// 启停虚拟机
const toggleVmStatus = async (vm) => {
  if (togglingVmNames.value.includes(vm.Name)) return;
  togglingVmNames.value = [...togglingVmNames.value, vm.Name];
  errorMessage.value = '';

  try {
    const jwtToken = getJwtToken();
    const action = vm.Status === 'Running' ? 'stop' : 'start';
    const response = await fetch(`/api/vm/${encodeURIComponent(vm.Name)}/${action}`, {
      headers: { Authorization: `${jwtToken}` }
    });
    if (!response.ok) {
      const msg = await getErrorMessage(response, '操作失败，请稍后重试');
      throw new Error(msg);
    }
    await fetchVmList();
    showToast(`${vm.Name} ${action === 'start' ? '启动' : '停止'}指令下发成功`);
  } catch (error) {
    errorMessage.value = error.message;
    showToast(error.message);
  } finally {
    togglingVmNames.value = togglingVmNames.value.filter(name => name !== vm.Name);
  }
};

// 打开弹窗并重置表单
const openCreateDialog = () => {
  // 重置为默认值
  createForm.value = {
    name: '',
    cpu: 2,
    memory: 4,
    node: nodeList.value.length ? nodeList.value[0] : '',
    diskSize: 20,
    networkType: 'bridge',
    autoStart: true
  };
  dialog.value = true;
};

// 关闭弹窗
const closeCreateDialog = () => {
  dialog.value = false;
  submitting.value = false;
};

// 提交创建虚拟机接口
const submitCreateVm = async () => {
  // 触发表单校验
  await createVmForm.value.validate();
  if (!formValid.value) return;

  submitting.value = true;
  try {
    const jwtToken = getJwtToken();
    const res = await fetch('/api/vm/create', {
      method: 'POST',
      headers: {
        Authorization: `${jwtToken}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(createForm.value)
    });

    if (!res.ok) {
      const errMsg = await getErrorMessage(res, '创建虚拟机失败');
      throw new Error(errMsg);
    }

    showToast('虚拟机创建任务提交成功！');
    closeCreateDialog();
    // 刷新列表
    await fetchVmList();
  } catch (err) {
    showToast(err.message);
  } finally {
    submitting.value = false;
  }
};

// 退出登录
const logout = () => {
  cookies.remove('jwt_token');
  router.push('/');
};

onMounted(() => {
  fetchVmList();
  fetchNodeList();
  refreshInterval = setInterval(fetchVmList, 15 * 1000);
});

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
  padding-bottom: 10px;
}

.logout-btn,
.power-btn,
.new-vm-btn {
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

.new-vm-btn {
  min-width: 72px;
  padding: 8px 10px;
  color: #fff;
  background: linear-gradient(135deg, #e26b3b, #ca1e1e);
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