<template>
  <v-container class="py-8" fluid>
    <v-snackbar v-model="snackbarVisible" :timeout="10000" color="error" location="top" variant="elevated" multi-line>
      {{ snackbarMessage }}
      <template v-slot:actions>
        <v-btn color="blue" variant="text" @click="snackbarVisible = false"> x </v-btn>
      </template>
    </v-snackbar>

    <v-card elevation="4" class="mx-auto" max-width="1400">
      <v-card-title>
        <div>
          
          <h2 class="text-h4 font-weight-bold mt-1">虚拟机列表</h2>
          <p class="text-body-1 text-grey mt-1">管理实例运行状态、新建虚拟机、查看就绪情况与启动/停止操作。(删除，修改等操作暂不支持，请在Kubernetes集群中手动操作)</p>
        </div>
      </v-card-title>

      <v-card-actions class="px-6 position-absolute top-0 right-0 mt-1">
        <v-btn color="error" variant="elevated" @click="logout" class="white--text">
          <v-icon left>mdi-logout</v-icon>
          退出登录
        </v-btn>
      </v-card-actions>

      <v-card-actions class="px-6">
        <v-btn color="primary" variant="elevated" @click="openCreateDialog" class="white--text">
          <v-icon left>mdi-plus</v-icon>
          新建虚拟机
        </v-btn>
      </v-card-actions>

      <v-card-text>
        <v-progress-linear indeterminate v-if="isLoading" class="mb-4"></v-progress-linear>

        <v-data-table v-else :items="vmList" :headers="headers" :loading="isLoading" class="elevation-1"
          hide-default-footer>
          <template v-slot:item.name="{ item }">
            <a v-if="item.Status === 'Running'" :href="`/novnc/vnc.html?path=/vm/${item.Name}/vnc&autoconnect=true`"
              target="_blank" rel="noopener noreferrer" class="text-blue-500 hover:text-blue-400">
              {{ item.Name }}
            </a>
            <span v-else>{{ item.Name }}</span>
          </template>

          <template v-slot:item.config="{ item }">
            {{ item.CPU || '0 vCPU' }} / {{ item.Memory || '未配置' }}
          </template>

          <template v-slot:item.node="{ item }">
            {{ item.Status === 'Running' ? (item.NodeName || '未调度') : '未调度' }}
          </template>

          <template v-slot:item.ready="{ item }">
            <v-chip :color="item.Ready ? 'success' : 'error'" size="small" class="white--text">
              {{ item.Ready ? '就绪' : '未就绪' }}
            </v-chip>
          </template>

          <template v-slot:item.status="{ item }">
            <v-chip :color="item.Status === 'Running' ? 'success' : 'error'" size="small" class="white--text">
              {{ item.Status }}
            </v-chip>
          </template>

          <template v-slot:item.startTime="{ item }">
            {{ item.Status === 'Running' ? (item.StartTime || '') : '' }}
          </template>

          <template v-slot:item.actions="{ item }">
            <v-btn color="orange" :disabled="togglingVmNames.includes(item.Name)"
              :loading="togglingVmNames.includes(item.Name)" size="small" class="white--text"
              @click="toggleVmStatus(item)">
              {{ item.Status != 'Stopped' ? '停止' : '启动' }}
            </v-btn>
          </template>
        </v-data-table>

        <v-empty-state v-if="!isLoading && !vmList.length" class="py-12">
          <v-icon size="48" color="grey">mdi-server-off</v-icon>
          <template v-slot:heading>暂无虚拟机数据</template>
          <template v-slot:subtitle>点击上方按钮创建新的虚拟机</template>
        </v-empty-state>
      </v-card-text>
    </v-card>

    <v-dialog v-model="dialog" max-width="800" persistent>
      <v-card title="配置新建虚拟机参数">
        <v-card-text>
          <v-form ref="createVmForm" v-model="formValid" lazy-validation>
            <v-row gap="12">
              <v-col cols="12" sm="6">
                <v-text-field v-model="createForm.name" label="虚拟机名称" placeholder="仅小写字母、数字、横杠" required
                  :rules="nameRules" outlined dense></v-text-field>
              </v-col>
              <v-col cols="12" sm="6">
                <v-select v-model="createForm.node" label="在哪个节点创建" :items="nodeList" required outlined
                  dense></v-select>
              </v-col>
            </v-row>

            <v-row gap="12">
              <v-col cols="12" sm="6">
                <v-text-field v-model.number="createForm.cpu" label="CPU(vCPU)" type="number" min="1" required
                  :rules="numRules" outlined dense></v-text-field>
              </v-col>
              <v-col cols="12" sm="6">
                <v-text-field v-model.number="createForm.memory" label="内存大小(GiB)" type="number" min="1" required
                  :rules="numRules" outlined dense></v-text-field>
              </v-col>
            </v-row>

            <v-row gap="12">
              <v-col cols="12" sm="6">
                <v-select v-model="createForm.networkType" label="网络模式" :items="networkOptions" required outlined
                  dense></v-select>
              </v-col>
              <v-col cols="12" sm="6">
                <v-text-field v-model.number="createForm.diskSize" label="磁盘存储大小(GiB)" type="number" min="10" required
                  :rules="numRules" outlined dense></v-text-field>
              </v-col>
            </v-row>

            <v-row gap="12">
              <v-col cols="12" sm="4">
                <v-checkbox label="挂载ISO盘" v-model="createForm.mountISO"></v-checkbox>
              </v-col>
              <v-col v-if="createForm.mountISO" cols="12" sm="8">
                <v-text-field v-model="createForm.isoPath" label="ISO文件路径" placeholder="例如: /root/ISO/ubuntu-22.04.iso"
                  required :rules="pathRules" outlined dense></v-text-field>
              </v-col>
            </v-row>

            <v-switch v-model="createForm.autoStart" label="创建后自动开机" color="primary" class="mt-4">
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
  </v-container>
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
  { title: 'masquerade网络(缺省)', value: 'masquerade' },
  { title: 'bridge桥接网络', value: 'bridge' }
]);

const headers = ref([
  { title: '名称', align: 'start', key: 'name' },
  { title: '基本配置', align: 'start', key: 'config' },
  { title: '运行节点', align: 'start', key: 'node' },
  { title: '就绪状态', align: 'start', key: 'ready' },
  { title: '当前状态', align: 'start', key: 'status' },
  { title: '启动时间', align: 'start', key: 'startTime' },
  { title: '操作', align: 'center', key: 'actions' }
]);

const nameRules = [
  v => !!v || '虚拟机名称不能为空',
  v => /^[a-z0-9-]+$/.test(v) || '只能小写字母、数字、横杠'
];
const pathRules = [v => !!v || 'ISO盘路径不能为空', v => /^[a-zA-Z0-9-/_.]+$/.test(v) || '只能包含字母、数字、斜杠、横杆、下划线、点号'];
const numRules = [v => v >= 1 || '数值必须大于等于1'];

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
  } catch (err) { }
  return fallbackMessage;
};

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
    const nodes = await response.json();
    nodeList.value = Array.isArray(nodes) ? nodes : [];
  } catch (error) {
    showToast(error.message);
    nodeList.value = [];
  }
};

const getJwtToken = () => {
  const jwtToken = cookies.get('jwt_token');
  if (!jwtToken) throw new Error('未找到有效的 JWT Token，请重新登录');
  return jwtToken;
};

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
      headers: { Authorization: `${jwtToken}` }
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
 //   console.log ('Submitting create VM form with data:', createForm.value);
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