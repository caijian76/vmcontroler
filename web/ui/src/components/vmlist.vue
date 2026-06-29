<template>
  <div>
    <v-app-bar color="primary" elevation="4" scroll-behavior="elevate">
      <template v-slot:prepend>
        <v-avatar
          size="40"
          rounded="lg"
          color="white"
          class="mr-2"
        >
          <v-icon color="primary" size="24">mdi-server</v-icon>
        </v-avatar>
      </template>
      <v-app-bar-title class="font-weight-bold">虚拟机管理平台</v-app-bar-title>
      <v-spacer />
      <v-btn
        variant="text"
        color="white"
        @click="logout"
        prepend-icon="mdi-logout"
      >
        退出登录
      </v-btn>
    </v-app-bar>

    <v-main>
      <v-container fluid class="pa-6">
        <v-snackbar
          v-model="snackbarVisible"
          :timeout="6000"
          color="error"
          location="top"
          variant="elevated"
          multi-line
        >
          {{ snackbarMessage }}
          <template v-slot:actions>
            <v-btn color="white" variant="text" @click="snackbarVisible = false">
              <v-icon>mdi-close</v-icon>
            </v-btn>
          </template>
        </v-snackbar>

        <v-row>
          <v-col cols="12">
            <v-card elevation="4" class="rounded-xl overflow-hidden">
              <v-card-item class="bg-surface-light pt-5 pb-3">
                <template v-slot:title>
                  <div class="d-flex align-center">
                    <v-icon color="primary" class="mr-2" size="28">mdi-server-network</v-icon>
                    <span class="text-h5 font-weight-bold">虚拟机列表</span>
                  </div>
                </template>
                <template v-slot:subtitle>
                  <span class="text-body-2 text-medium-emphasis">
                    管理实例运行状态、新建虚拟机、查看就绪情况与启动/停止操作
                  </span>
                </template>
              </v-card-item>

              <v-divider />

              <v-card-actions class="px-4 py-3">
                <v-btn
                  color="primary"
                  variant="elevated"
                  prepend-icon="mdi-plus-circle-outline"
                  @click="openCreateDialog"
                  rounded="lg"
                >
                  新建虚拟机
                </v-btn>
                <v-spacer />
                <v-chip
                  color="info"
                  variant="tonal"
                  size="small"
                  label
                >
                  <v-icon start size="18">mdi-refresh</v-icon>
                  每 15 秒自动刷新
                </v-chip>
              </v-card-actions>

              <v-divider />

              <v-card-text class="pa-0">
                <v-progress-linear
                  v-if="isLoading"
                  color="primary"
                  indeterminate
                  height="3"
                />

                <v-data-table
                  v-else
                  :items="vmList"
                  :headers="headers"
                  :loading="isLoading"
                  hover
                  items-per-page="-1"
                  class="elevation-0"
                >
                  <template v-slot:item.name="{ item }">
                    <a
                      v-if="item.Status === 'Running'"
                      :href="`/novnc/vnc.html?path=/vm/${item.Name}/vnc&autoconnect=true`"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="text-decoration-none font-weight-medium"
                      style="color: #1a73e8"
                    >
                      <v-icon size="18" class="mr-1">mdi-open-in-new</v-icon>
                      {{ item.Name }}
                    </a>
                    <span v-else class="text-medium-emphasis">{{ item.Name }}</span>
                  </template>

                  <template v-slot:item.config="{ item }">
                    <v-chip size="small" variant="tonal" color="purple" class="mr-1">
                      <v-icon start size="16">mdi-cpu-64-bit</v-icon>
                      {{ item.CPU || '0' }}
                    </v-chip>
                    <v-chip size="small" variant="tonal" color="teal">
                      <v-icon start size="16">mdi-memory</v-icon>
                      {{ item.Memory || '未配置' }}
                    </v-chip>
                  </template>

                  <template v-slot:item.Node="{ item }">
                    <v-chip
                      v-if="item.Status === 'Running' && item.NodeName"
                      size="small"
                      variant="tonal"
                      color="blue-grey"
                    >
                      <v-icon start size="16">mdi-server</v-icon>
                      {{ item.NodeName }}
                    </v-chip>
                    <span v-else class="text-medium-emphasis text-caption">未调度</span>
                  </template>

                  <template v-slot:item.Ready="{ item }">
                    <v-chip
                      :color="item.Ready ? 'success' : 'warning'"
                      size="small"
                      variant="tonal"
                      label
                    >
                      <v-icon start size="16">
                        {{ item.Ready ? 'mdi-check-circle' : 'mdi-alert-circle' }}
                      </v-icon>
                      {{ item.Ready ? '就绪' : '未就绪' }}
                    </v-chip>
                  </template>

                  <template v-slot:item.status="{ item }">
                    <v-chip
                      :color="item.Status === 'Running' ? 'success' : 'grey-darken-1'"
                      size="small"
                      variant="tonal"
                      label
                    >
                      <v-icon start size="16">
                        {{ item.Status === 'Running' ? 'mdi-play-circle' : 'mdi-stop-circle' }}
                      </v-icon>
                      {{ item.Status }}
                    </v-chip>
                  </template>

                  <template v-slot:item.startTime="{ item }">
                    <span v-if="item.Status === 'Running' && item.StartTime" class="text-body-2">
                      <v-icon size="16" class="mr-1">mdi-clock-outline</v-icon>
                      {{ item.StartTime }}
                    </span>
                    <span v-else class="text-medium-emphasis text-caption">—</span>
                  </template>

                  <template v-slot:item.actions="{ item }">
                    <v-btn
                      :color="item.Status !== 'Stopped' ? 'error' : 'success'"
                      variant="tonal"
                      size="small"
                      rounded="lg"
                      :disabled="togglingVmNames.includes(item.Name)"
                      :loading="togglingVmNames.includes(item.Name)"
                      @click="toggleVmStatus(item)"
                    >
                      <v-icon start size="18">
                        {{ item.Status !== 'Stopped' ? 'mdi-stop' : 'mdi-play' }}
                      </v-icon>
                      {{ item.Status !== 'Stopped' ? '停止' : '启动' }}
                    </v-btn>
                  </template>
                </v-data-table>

                <v-empty-state
                  v-if="!isLoading && !vmList.length"
                  class="py-16"
                  icon="mdi-server-off"
                  title="暂无虚拟机数据"
                  text="点击上方「新建虚拟机」按钮创建您的第一个虚拟机"
                  color="grey"
                >
                  <template v-slot:media>
                    <v-icon size="80" color="grey-lighten-1">mdi-server-off</v-icon>
                  </template>
                </v-empty-state>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <v-dialog v-model="dialog" max-width="750" persistent>
          <v-card rounded="xl">
            <v-card-item class="bg-primary pb-4">
              <template v-slot:prepend>
                <v-avatar rounded="lg" color="white" size="44">
                  <v-icon color="primary" size="28">mdi-plus-circle</v-icon>
                </v-avatar>
              </template>
              <template v-slot:title>
                <span class="text-h6 font-weight-bold text-white">新建虚拟机</span>
              </template>
              <template v-slot:subtitle>
                <span class="text-body-2" style="color: rgba(255,255,255,0.7)">
                  配置虚拟机参数并提交创建任务
                </span>
              </template>
            </v-card-item>

            <v-card-text class="pt-6">
              <v-form ref="createVmForm" v-model="formValid" lazy-validation>
                <v-row dense>
                  <v-col cols="12" sm="6">
                    <v-text-field
                      v-model="createForm.name"
                      label="虚拟机名称"
                      placeholder="仅小写字母、数字、横杠"
                      variant="outlined"
                      density="comfortable"
                      :rules="nameRules"
                      prepend-inner-icon="mdi-tag-outline"
                    />
                  </v-col>
                  <v-col cols="12" sm="6">
                    <v-select
                      v-model="createForm.node"
                      label="目标节点"
                      :items="nodeList"
                      variant="outlined"
                      density="comfortable"
                      prepend-inner-icon="mdi-server"
                    />
                  </v-col>
                </v-row>

                <v-row dense>
                  <v-col cols="12" sm="6">
                    <v-text-field
                      v-model.number="createForm.cpu"
                      label="CPU (vCPU)"
                      type="number"
                      min="1"
                      variant="outlined"
                      density="comfortable"
                      :rules="numRules"
                      prepend-inner-icon="mdi-cpu-64-bit"
                      suffix="核"
                    />
                  </v-col>
                  <v-col cols="12" sm="6">
                    <v-text-field
                      v-model.number="createForm.memory"
                      label="内存大小"
                      type="number"
                      min="1"
                      variant="outlined"
                      density="comfortable"
                      :rules="numRules"
                      prepend-inner-icon="mdi-memory"
                      suffix="GiB"
                    />
                  </v-col>
                </v-row>

                <v-row dense>
                  <v-col cols="12" sm="6">
                    <v-select
                      v-model="createForm.networkType"
                      label="网络模式"
                      :items="networkOptions"
                      variant="outlined"
                      density="comfortable"
                      prepend-inner-icon="mdi-lan"
                    />
                  </v-col>
                  <v-col cols="12" sm="6">
                    <v-text-field
                      v-model.number="createForm.diskSize"
                      label="磁盘存储大小"
                      type="number"
                      min="10"
                      variant="outlined"
                      density="comfortable"
                      :rules="numRules"
                      prepend-inner-icon="mdi-harddisk"
                      suffix="GiB"
                    />
                  </v-col>
                </v-row>

                <v-row dense>
                  <v-col cols="12" sm="4">
                    <v-checkbox
                      label="挂载 ISO 镜像"
                      v-model="createForm.mountISO"
                      color="primary"
                      density="compact"
                    />
                  </v-col>
                  <v-col v-if="createForm.mountISO" cols="12" sm="8">
                    <v-text-field
                      v-model="createForm.isoPath"
                      label="ISO 文件路径"
                      placeholder="例如: /root/ISO/ubuntu-22.04.iso"
                      variant="outlined"
                      density="comfortable"
                      :rules="pathRules"
                      prepend-inner-icon="mdi-disc"
                    />
                  </v-col>
                </v-row>

                <v-switch
                  v-model="createForm.autoStart"
                  label="创建后自动开机"
                  color="primary"
                  hide-details
                  class="mt-2"
                />
              </v-form>
            </v-card-text>

            <v-divider />

            <v-card-actions class="pa-4">
              <v-spacer />
              <v-btn
                variant="tonal"
                color="grey"
                rounded="lg"
                @click="closeCreateDialog"
              >
                取消
              </v-btn>
              <v-btn
                color="primary"
                variant="elevated"
                rounded="lg"
                :loading="submitting"
                prepend-icon="mdi-check-circle"
                @click="submitCreateVm"
              >
                确认创建
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-container>
    </v-main>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { useCookies } from 'vue3-cookies';
import { useRouter } from 'vue-router';

const { cookies } = useCookies();
const router = useRouter();

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
  { title: 'masquerade 网络（缺省）', value: 'masquerade' },
  { title: 'bridge 桥接网络', value: 'bridge' },
]);

const headers = ref([
  { title: '名称', align: 'start', key: 'Name', sortable: true },
  { title: '基本配置', align: 'start', key: 'config', sortable: false },
  { title: '运行节点', align: 'start', key: 'Node', sortable: true },
  { title: '就绪状态', align: 'start', key: 'Ready', sortable: true },
  { title: '当前状态', align: 'start', key: 'status', sortable: false },
  { title: '启动时间', align: 'start', key: 'startTime', sortable: false },
  { title: '操作', align: 'center', key: 'actions', sortable: false },
]);

const nameRules = [
  v => !!v || '虚拟机名称不能为空',
  v => /^[a-z0-9-]+$/.test(v) || '只能小写字母、数字、横杠',
];
const pathRules = [
  v => !!v || 'ISO 路径不能为空',
  v => /^[a-zA-Z0-9-/_.]+$/.test(v) || '只能包含字母、数字、斜杠、横杆、下划线、点号',
];
const numRules = [v => v >= 1 || '数值必须大于等于 1'];

let refreshInterval;

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
  } catch (err) {
    // ignore
  }
  return fallbackMessage;
};

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
    const nodes = await response.json();
    nodeList.value = Array.isArray(nodes) ? nodes : [];
  } catch (error) {
    showToast(error.message);
    nodeList.value = [];
  }
};

const getJwtToken = () => {
  const jwtToken = cookies.get('jwt_token');
  if (!jwtToken) {
    // showToast('登录已过期，请重新登录');
    setTimeout(() => {
      router.push('/');
    }, 2000);
    throw new Error('未找到有效的 JWT Token, 请重新登录');
  }
  return jwtToken;
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
    togglingVmNames.value = togglingVmNames.value.filter(name => name !== vm.Name);
  }
};

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

onMounted(() => {
  fetchVmList();
  fetchNodeList();
  refreshInterval = setInterval(fetchVmList, 15 * 1000);
});

onUnmounted(() => {
  clearInterval(refreshInterval);
});
</script>