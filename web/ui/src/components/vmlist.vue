<template>
  <v-app>
    <!-- Top Navigation Bar -->
    <v-app-bar color="surface" elevation="2" height="64" class="px-4 navbar">
      <div class="d-flex align-center ga-3">
        <div class="nav-logo">
          <v-icon color="white" size="24">mdi-server-network</v-icon>
        </div>
        <h1 class="text-subtitle-1 font-weight-bold text-primary d-none d-sm-block">VM Control Center</h1>
        <v-divider vertical class="mx-2 d-none d-sm-block"></v-divider>
        <v-chip size="x-small" color="primary" variant="tonal">
          <v-icon start size="14">mdi-cloud</v-icon>
          KubeVirt
        </v-chip>
      </div>

      <v-spacer />

      <div class="d-flex align-center ga-2">
        <v-chip size="x-small" variant="tonal" color="success">
          <v-icon start size="14">mdi-server</v-icon>
          {{ vmList.length }}
        </v-chip>
        <v-chip size="x-small" variant="tonal" color="info">
          <v-icon start size="14">mdi-network-node</v-icon>
          {{ nodeList.length }}
        </v-chip>
        <v-btn
          icon="mdi-refresh"
          @click="refreshData"
          :loading="isLoading"
          size="small"
          rounded="lg"
          variant="text"
        />
        <v-btn color="error" variant="tonal" size="small" rounded="lg" @click="logout" prepend-icon="mdi-logout">
          <span class="d-none d-sm-inline">退出</span>
        </v-btn>
      </div>
    </v-app-bar>

    <!-- Main Content -->
    <v-main class="bg-dashboard">
      <v-container fluid class="pa-4 pa-md-6">
        <!-- Stats Cards -->
        <v-row class="mb-4">
          <v-col cols="6" md="3">
            <v-card rounded="xl" class="stat-card" variant="flat" elevation="1">
              <v-card-text class="pa-4">
                <div class="d-flex align-center justify-space-between">
                  <div>
                    <p class="text-caption text-medium-emphasis mb-1">总虚拟机</p>
                    <p class="text-h4 font-weight-bold mb-0 text-primary">{{ vmList.length }}</p>
                  </div>
                  <div class="stat-icon stat-icon--blue">
                    <v-icon color="primary" size="24">mdi-server</v-icon>
                  </div>
                </div>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="6" md="3">
            <v-card rounded="xl" class="stat-card" variant="flat" elevation="1">
              <v-card-text class="pa-4">
                <div class="d-flex align-center justify-space-between">
                  <div>
                    <p class="text-caption text-medium-emphasis mb-1">运行中</p>
                    <p class="text-h4 font-weight-bold mb-0 text-success">{{ runningCount }}</p>
                  </div>
                  <div class="stat-icon stat-icon--green">
                    <v-icon color="success" size="24">mdi-play-circle</v-icon>
                  </div>
                </div>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="6" md="3">
            <v-card rounded="xl" class="stat-card" variant="flat" elevation="1">
              <v-card-text class="pa-4">
                <div class="d-flex align-center justify-space-between">
                  <div>
                    <p class="text-caption text-medium-emphasis mb-1">已停止</p>
                    <p class="text-h4 font-weight-bold mb-0 text-warning">{{ stoppedCount }}</p>
                  </div>
                  <div class="stat-icon stat-icon--orange">
                    <v-icon color="warning" size="24">mdi-stop-circle</v-icon>
                  </div>
                </div>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="6" md="3">
            <v-card rounded="xl" class="stat-card" variant="flat" elevation="1">
              <v-card-text class="pa-4">
                <div class="d-flex align-center justify-space-between">
                  <div>
                    <p class="text-caption text-medium-emphasis mb-1">就绪</p>
                    <p class="text-h4 font-weight-bold mb-0 text-purple">{{ readyCount }}</p>
                  </div>
                  <div class="stat-icon stat-icon--purple">
                    <v-icon color="purple" size="24">mdi-check-circle</v-icon>
                  </div>
                </div>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <!-- VM Table Card -->
        <v-card rounded="xl" elevation="1" class="table-card">
          <v-card-item class="pb-2">
            <template v-slot:title>
              <div class="d-flex align-center ga-2">
                <v-icon color="primary" size="22">mdi-vm</v-icon>
                <span class="text-h6 font-weight-bold">虚拟机列表</span>
              </div>
            </template>
            <template v-slot:subtitle>
              <span class="text-caption text-medium-emphasis">
                自动刷新 15s · 点击名称打开 VNC 控制台
              </span>
            </template>
            <template v-slot:append>
              <v-btn color="primary" variant="flat" rounded="lg" size="small" @click="openCreateDialog" prepend-icon="mdi-plus">
                新建虚拟机
              </v-btn>
            </template>
          </v-card-item>

          <v-divider />

          <v-card-text class="pa-0">
            <v-progress-linear indeterminate color="primary" v-if="isLoading" rounded height="3" />

            <v-data-table-virtual
              v-else
              :items="vmList"
              :headers="headers"
              class="vm-table"
              hide-default-footer
              :items-per-page="-1"
              :search="search"
            >
              <template v-slot:no-data>
                <v-empty-state
                  v-if="!isLoading && !vmList.length"
                  headline="暂无虚拟机数据"
                  text="点击下方按钮创建新的虚拟机"
                  class="text-center py-12"
                >
                  <template v-slot:icon>
                    <v-icon size="64" color="medium-emphasis">mdi-server-network-off</v-icon>
                  </template>
                </v-empty-state>
              </template>

              <!-- Name -->
              <template v-slot:item.name="{ item }">
                <div class="d-flex align-center ga-2">
                  <v-icon
                    v-if="item.Status === 'Running'"
                    color="success"
                    size="14"
                  >mdi-circle-small</v-icon>
                  <v-icon v-else color="disabled" size="14">mdi-circle-small-outline</v-icon>
                  <a
                    v-if="item.Status === 'Running' && item.Name"
                    :href="`/novnc/vnc.html?path=/vm/${encodeURIComponent(item.Name)}/vnc&autoconnect=true`"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="text-primary font-weight-medium"
                  >
                    {{ item.Name }}
                  </a>
                  <span v-else class="font-weight-medium">{{ item.Name || '-' }}</span>
                </div>
              </template>

              <!-- Config -->
              <template v-slot:item.config="{ item }">
                <div class="d-flex flex-column ga-0.5">
                  <div class="d-flex align-center ga-1">
                    <v-icon size="14" color="medium-emphasis">mdi-cpu-64-bit</v-icon>
                    <span class="text-body-2">{{ item.CPU || '0 vCPU' }}</span>
                  </div>
                  <div class="d-flex align-center ga-1">
                    <v-icon size="14" color="medium-emphasis">mdi-memory</v-icon>
                    <span class="text-body-2">{{ item.Memory || '-' }}</span>
                  </div>
                </div>
              </template>

              <!-- Node -->
              <template v-slot:item.node="{ item }">
                <v-chip size="x-small" variant="tonal" color="secondary" rounded="lg">
                  <v-icon start size="14">mdi-network-node</v-icon>
                  {{ item.Status === 'Running' ? (item.NodeName || '未调度') : '未调度' }}
                </v-chip>
              </template>

              <!-- Ready -->
              <template v-slot:item.ready="{ item }">
                <v-chip
                  size="x-small"
                  :color="item.Ready ? 'success' : 'warning'"
                  variant="flat"
                  rounded="lg"
                >
                  <v-icon start size="12">{{ item.Ready ? 'mdi-check' : 'mdi-alert' }}</v-icon>
                  {{ item.Ready ? '就绪' : '未就绪' }}
                </v-chip>
              </template>

              <!-- Status -->
              <template v-slot:item.status="{ item }">
                <StatusChip :status="item.Status" />
              </template>

              <!-- StartTime -->
              <template v-slot:item.startTime="{ item }">
                <span v-if="item.Status === 'Running' && item.StartTime" class="text-body-2 text-medium-emphasis">
                  <v-icon size="14" start>mdi-clock-outline</v-icon>
                  {{ item.StartTime }}
                </span>
                <span v-else class="text-body-2 text-disabled">—</span>
              </template>

              <!-- Actions -->
              <template v-slot:item.actions="{ item }">
                <v-btn-group density="compact" rounded="lg" variant="tonal">
                  <v-btn
                    v-if="item.Status === 'Stopped'"
                    size="x-small"
                    color="success"
                    variant="tonal"
                    rounded="lg"
                    :loading="togglingVmNames.includes(item.Name)"
                    :disabled="togglingVmNames.includes(item.Name)"
                    @click="toggleVmStatus(item)"
                  >
                    <v-icon start size="14">mdi-play</v-icon>
                    启动
                  </v-btn>
                  <v-btn
                    v-if="item.Status === 'Running'"
                    size="x-small"
                    color="warning"
                    variant="tonal"
                    rounded="lg"
                    :loading="togglingVmNames.includes(item.Name)"
                    :disabled="togglingVmNames.includes(item.Name)"
                    @click="toggleVmStatus(item)"
                  >
                    <v-icon start size="14">mdi-pause</v-icon>
                    停止
                  </v-btn>
                  <v-btn
                    v-if="['Starting','Stopping','Pending'].includes(item.Status)"
                    size="x-small"
                    rounded="lg"
                    disabled
                  >
                    <v-icon size="14" color="info">mdi-loading</v-icon>
                    处理中
                  </v-btn>
                </v-btn-group>
              </template>
            </v-data-table-virtual>
          </v-card-text>
        </v-card>
      </v-container>
    </v-main>

    <!-- Create VM Dialog -->
    <v-dialog v-model="dialog" max-width="680" persistent>
      <v-card rounded="xl" elevation="8">
        <v-card-item>
          <template v-slot:title>
            <div class="d-flex align-center ga-2">
              <v-icon color="primary">mdi-plus-box-multiple</v-icon>
              新建虚拟机
            </div>
          </template>
          <template v-slot:subtitle>
            配置虚拟机的基本参数和资源规格
          </template>
          <template v-slot:append>
            <v-btn icon="mdi-close" @click="closeCreateDialog" size="small" rounded="lg" variant="text" />
          </template>
        </v-card-item>

        <v-divider />

        <v-card-text class="pt-6">
          <v-form ref="createVmForm" v-model="formValid" lazy-validation>
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="createForm.name"
                  label="虚拟机名称"
                  placeholder="仅小写字母、数字、横杠"
                  required
                  :rules="nameRules"
                  variant="outlined"
                  rounded="lg"
                  density="comfortable"
                  color="primary"
                  prepend-inner-icon="mdi-name"
                  hide-details="auto"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="createForm.node"
                  label="目标节点"
                  :items="nodeList"
                  required
                  variant="outlined"
                  rounded="lg"
                  density="comfortable"
                  color="primary"
                  prepend-inner-icon="mdi-server"
                  hide-details="auto"
                />
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model.number="createForm.cpu"
                  label="CPU (vCPU)"
                  type="number"
                  min="1"
                  required
                  :rules="numRules"
                  variant="outlined"
                  rounded="lg"
                  density="comfortable"
                  color="primary"
                  prepend-inner-icon="mdi-cpu-64-bit"
                  hide-details="auto"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model.number="createForm.memory"
                  label="内存 (GiB)"
                  type="number"
                  min="1"
                  required
                  :rules="numRules"
                  variant="outlined"
                  rounded="lg"
                  density="comfortable"
                  color="primary"
                  prepend-inner-icon="mdi-memory"
                  hide-details="auto"
                />
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="12" md="6">
                <v-select
                  v-model="createForm.networkType"
                  label="网络模式"
                  :items="networkOptions"
                  required
                  variant="outlined"
                  rounded="lg"
                  density="comfortable"
                  color="primary"
                  prepend-inner-icon="mdi-network"
                  hide-details="auto"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model.number="createForm.diskSize"
                  label="磁盘存储 (GiB)"
                  type="number"
                  min="10"
                  required
                  :rules="numRules"
                  variant="outlined"
                  rounded="lg"
                  density="comfortable"
                  color="primary"
                  prepend-inner-icon="mdi-harddisk"
                  hide-details="auto"
                />
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="12">
                <v-switch
                  v-model="createForm.mountISO"
                  label="挂载 ISO 盘"
                  color="primary"
                  hide-details
                  class="mb-2"
                />
              </v-col>
              <v-col v-if="createForm.mountISO" cols="12">
                <v-text-field
                  v-model="createForm.isoPath"
                  label="ISO 文件路径"
                  placeholder="/root/ISO/ubuntu-22.04.iso"
                  :rules="pathRules"
                  variant="outlined"
                  rounded="lg"
                  density="comfortable"
                  color="primary"
                  prepend-inner-icon="mdi-cd-video-frame"
                  hide-details="auto"
                />
              </v-col>
            </v-row>

            <v-switch v-model="createForm.autoStart" label="创建后自动启动" color="primary" hide-details />
          </v-form>
        </v-card-text>

        <v-divider />

        <v-card-actions class="pa-4">
          <v-spacer />
          <v-btn variant="text" color="medium-emphasis" rounded="lg" @click="closeCreateDialog">
            取消
          </v-btn>
          <v-btn
            color="primary"
            variant="flat"
            rounded="lg"
            :loading="submitting"
            @click="submitCreateVm"
          >
            确认创建
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Snackbar -->
    <v-snackbar
      v-model="snackbarVisible"
      :timeout="5000"
      location="top"
      variant="tonal"
      color="error"
      rounded="lg"
    >
      {{ snackbarMessage }}
      <template v-slot:actions>
        <v-btn color="white" variant="text" size="small" @click="snackbarVisible = false">关闭</v-btn>
      </template>
    </v-snackbar>
  </v-app>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useCookies } from 'vue3-cookies';
import { useRouter } from 'vue-router';
import StatusChip from './StatusChip.vue';

const { cookies } = useCookies();
const router = useRouter();

// --- Data ---
const vmList = ref([]);
const nodeList = ref([]);
const isLoading = ref(true);
const togglingVmNames = ref([]);
const snackbarVisible = ref(false);
const snackbarMessage = ref('');
const dialog = ref(false);
const submitting = ref(false);
const formValid = ref(false);
const createVmForm = ref(null);
const search = ref('');
let refreshInterval;

const createForm = ref({
  name: '',
  cpu: 2,
  memory: 4,
  node: '',
  diskSize: 20,
  networkType: 'masquerade',
  autoStart: true,
  mountISO: false,
  isoPath: '',
});

const networkOptions = ref([
  { title: 'Masquerade 网络（默认）', value: 'masquerade' },
  { title: 'Bridge 桥接网络', value: 'bridge' },
]);

const headers = ref([
  { title: '名称', align: 'start', key: 'name', sortable: true },
  { title: '资源配置', align: 'start', key: 'config', sortable: false },
  { title: '运行节点', align: 'start', key: 'node', sortable: true },
  { title: '就绪', align: 'start', key: 'ready', sortable: false },
  { title: '状态', align: 'start', key: 'status', sortable: true },
  { title: '启动时间', align: 'start', key: 'startTime', sortable: true },
  { title: '操作', align: 'center', key: 'actions', sortable: false },
]);

const nameRules = [
  v => !!v || '虚拟机名称不能为空',
  v => /^[a-z0-9-]+$/.test(v) || '只能小写字母、数字、横杠',
];
const pathRules = [
  v => !createForm.value.mountISO || !!v || 'ISO 盘路径不能为空',
  v => !createForm.value.mountISO || /^[a-zA-Z0-9-/_.]+$/.test(v) || '只能包含字母、数字、斜杠、横杆、下划线、点号',
];
const numRules = [v => v >= 1 || '数值必须大于等于 1'];

// --- Computed ---
const runningCount = computed(() => vmList.value.filter(v => v.Status === 'Running').length);
const stoppedCount = computed(() => vmList.value.filter(v => v.Status === 'Stopped').length);
const readyCount = computed(() => vmList.value.filter(v => v.Ready).length);

// --- Helpers ---
const showToast = (message) => {
  if (!message) return;
  snackbarMessage.value = message;
  snackbarVisible.value = true;
};

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

const getJwtToken = () => {
  const jwtToken = cookies.get('jwt_token');
  if (!jwtToken) throw new Error('未找到有效的 JWT Token，请重新登录');
  return jwtToken;
};

// --- API Calls ---
const fetchNodeList = async () => {
  try {
    const jwtToken = getJwtToken();
    const response = await fetch('/api/nodes', {
      headers: { Authorization: `${jwtToken}` },
    });
    if (!response.ok) {
      const msg = await getErrorMessage(response, '获取节点列表失败');
      throw new Error(msg);
    }
    nodeList.value = Array.isArray(await response.json()) ? (await response.json()).slice() : [];
  } catch (error) {
    showToast(error.message);
    nodeList.value = [];
  }
};

const fetchVmList = async () => {
  try {
    const jwtToken = getJwtToken();
    const response = await fetch('/api/vm', {
      headers: { Authorization: `${jwtToken}` },
    });
    if (!response.ok) {
      const msg = await getErrorMessage(response, '请求失败，请稍后重试');
      throw new Error(msg);
    }
    const vms = await response.json();
    vmList.value = vms != null ? vms : [];
  } catch (error) {
    showToast(error.message);
  } finally {
    isLoading.value = false;
  }
};

const toggleVmStatus = async (vm) => {
  if (togglingVmNames.value.includes(vm.Name)) return;
  togglingVmNames.value = [...togglingVmNames.value, vm.Name];

  try {
    const jwtToken = getJwtToken();
    const action = vm.Status === 'Stopped' ? 'start' : 'stop';
    const response = await fetch(`/api/vm/${encodeURIComponent(vm.Name)}/${action}`, {
      headers: { Authorization: `${jwtToken}` },
    });
    if (!response.ok) {
      const msg = await getErrorMessage(response, '操作失败，请稍后重试');
      throw new Error(msg);
    }
    await fetchVmList();
  } catch (error) {
    showToast(error.message);
  } finally {
    togglingVmNames.value = togglingVmNames.value.filter((name) => name !== vm.Name);
  }
};

const refreshData = () => {
  isLoading.value = true;
  fetchVmList();
};

// --- Dialog ---
const openCreateDialog = () => {
  createForm.value = {
    name: '',
    cpu: 2,
    memory: 4,
    node: nodeList.value.length ? nodeList.value[0] : '',
    diskSize: 20,
    networkType: 'masquerade',
    autoStart: true,
    mountISO: false,
    isoPath: '',
  };
  dialog.value = true;
};

const closeCreateDialog = () => {
  dialog.value = false;
  submitting.value = false;
};

const submitCreateVm = async () => {
  await createVmForm.value.validate();
  if (!formValid.value) return;

  submitting.value = true;
  try {
    const jwtToken = getJwtToken();
    const res = await fetch('/api/vm/create', {
      method: 'POST',
      headers: {
        Authorization: `${jwtToken}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(createForm.value),
    });

    if (!res.ok) {
      const errMsg = await getErrorMessage(res, '创建虚拟机失败');
      throw new Error(errMsg);
    }

    showToast('虚拟机创建任务提交成功！');
    closeCreateDialog();
    await fetchVmList();
  } catch (err) {
    showToast(err.message);
  } finally {
    submitting.value = false;
  }
};

const logout = () => {
  cookies.remove('jwt_token');
  router.push('/');
};

// --- Lifecycle ---
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
.bg-dashboard {
  background: #f1f5f9;
}

.stat-card {
  border: 1px solid rgba(0, 0, 0, 0.06);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.08) !important;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon--blue   { background: rgba(59, 130, 246, 0.1); }
.stat-icon--green  { background: rgba(34, 197, 94, 0.1); }
.stat-icon--orange { background: rgba(249, 115, 22, 0.1); }
.stat-icon--purple { background: rgba(168, 85, 247, 0.1); }

.table-card {
  border: 1px solid rgba(0, 0, 0, 0.06);
}

.vm-table {
  font-size: 0.875rem;
}

.navbar {
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}

.nav-logo {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #3b82f6, #8b5cf6);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
