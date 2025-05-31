<template>
  <div>
    <!-- 右上角添加退出登录按钮 -->
    <div style="text-align: right; margin-bottom: 16px;">
      <button @click="logout">退出登录</button>
    </div>
    <h2>虚拟机列表</h2>
    <!-- 若加载中，显示加载提示 -->
    <div v-if="isLoading">加载中...</div>
    <!-- 若有错误，显示错误信息 -->
    <div v-else-if="errorMessage">{{ errorMessage }}</div>
    <!-- 若虚拟机列表不为空，显示表格 -->
    <table v-else-if="vmList.length">
      <thead>
        <tr>
          <th>名称</th>
          <th>运行策略</th>
          <th>就绪状态</th>
          <th>当前状态</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="vm in vmList" :key="vm.Name">
          <!-- 第一列根据虚拟机状态决定是否显示超链接 -->
          <td>
            <template v-if="vm.Status === 'Running'">
              <a 
                :href="`http://localhost:8080/novnc/vnc.html?path=/vm/${vm.Name}/vnc&autoconnect=true`" 
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
          <td>{{ vm.Run }}</td>
          <td>{{ vm.Ready ? '就绪' : '未就绪' }}</td>
          <td>{{ vm.Status }}</td>
        </tr>
      </tbody>
    </table>
    <!-- 若虚拟机列表为空，显示无数据提示 -->
    <div v-else>暂无虚拟机数据</div>
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
// 定义定时器变量
let refreshInterval;

// 获取虚拟机列表的函数
const fetchVmList = async () => {
  try {
    // 读取 jwt_token Cookie
    const jwtToken = cookies.get('jwt_token');
    if (!jwtToken) {
      throw new Error('未找到有效的 JWT Token，请重新登录');
    }

    // 发送请求获取虚拟机列表，添加 Authorization 请求头
    const response = await fetch('http://localhost:8080/api/vm', {
      headers: {
        Authorization: `${jwtToken}`
      }
    });
    if (!response.ok) {
      throw new Error('请求失败，请稍后重试');
    }
    // 解析响应数据
    vmList.value = await response.json();
  } catch (error) {
    // 捕获错误并显示错误信息
    errorMessage.value = error.message;
  } finally {
    // 无论请求成功或失败，都结束加载状态
    isLoading.value = false;
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
  refreshInterval = setInterval(fetchVmList, 5000);
});

// 组件卸载时清除定时器
onUnmounted(() => {
  clearInterval(refreshInterval);
});
</script>

<style scoped>
table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 16px;
}

th, td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}

th {
  background-color: #f2f2f2;
}

tr:nth-child(even) {
  background-color: #f9f9f9;
}

tr:hover {
  background-color: #f5f5f5;
}

a {
  color: #007bff;
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}
</style>
