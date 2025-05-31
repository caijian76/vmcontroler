<template>
  <div class="login-container">
    <h1>登录</h1>
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label for="username">用户名:</label>
        <input
          type="text"
          id="username"
          v-model="username"
          placeholder="请输入用户名"
          required
        />
      </div>
      <div class="form-group">
        <label for="password">密码:</label>
        <input
          type="password"
          id="password"
          v-model="password"
          placeholder="请输入密码"
          required
        />
      </div>
      <button type="submit" :disabled="isSubmitting">
        {{ isSubmitting ? '登录中...' : '登录' }}
      </button>
      <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useCookies } from 'vue3-cookies';
import { useRouter } from 'vue-router'; // 引入 useRouter

const { cookies } = useCookies();
const router = useRouter(); // 获取路由实例

// 定义响应式变量
const username = ref('');
const password = ref('');
const isSubmitting = ref(false);
const errorMessage = ref('');

// 处理登录逻辑
const handleLogin = async () => {
  isSubmitting.value = true;
  errorMessage.value = '';

  try {
    // 创建 FormData 对象
    const formData = new FormData();
    formData.append('username', username.value);
    formData.append('password', password.value);

    // 发送登录请求
    const response = await fetch('http://localhost:8080/login', {
      method: 'POST',
      // 注意：不要设置 Content-Type 头，浏览器会自动处理
      body: formData
    });

    if (!response.ok) {
      throw new Error('登录失败，请检查用户名和密码');
    }

    const data = await response.json();
    if (data.token) {
      // 保存 JWT token 到 cookie
      cookies.set('jwt_token', data.token, '1h'); // 设置 token 有效期为 1 小时
      // 登录成功后导航到 vmlist 页面
      router.push('/vmlist'); 
    } else {
      throw new Error('未收到有效的 token');
    }
  } catch (error) {
    errorMessage.value = error.message || '登录失败，请稍后重试';
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<style scoped>
.login-container {
  max-width: 400px;
  margin: 50px auto;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 5px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

/* 让 h1 元素内的文本居中 */
.login-container h1 {
  text-align: center;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
}

.form-group input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 3px;
}

button {
  width: 100%;
  padding: 10px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 3px;
  cursor: pointer;
}

button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.error-message {
  color: red;
  margin-top: 10px;
}
</style>
