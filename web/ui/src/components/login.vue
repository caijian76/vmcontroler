<template>
  <div class="login-page">
    <!-- Animated background -->
    <div class="bg-animation">
      <div class="blob blob-1"></div>
      <div class="blob blob-2"></div>
      <div class="blob blob-3"></div>
    </div>

    <div class="login-container">
      <v-card class="login-card elevation-8" rounded="xl">
        <v-card-text class="pa-8 pt-10">
          <!-- Logo & Title -->
          <div class="text-center mb-8">
            <div class="logo-wrapper mb-4">
              <v-icon size="48" color="primary">mdi-server-network</v-icon>
            </div>
            <h1 class="text-h4 font-weight-bold text-primary mb-1">VM Control Center</h1>
            <p class="text-body-2 text-medium-emphasis">KubeVirt 虚拟机管理平台</p>
          </div>

          <v-form @submit.prevent="handleLogin" ref="loginForm">
            <!-- Username -->
            <v-text-field
              v-model="username"
              label="用户名"
              placeholder="输入用户名"
              required
              autofocus
              variant="outlined"
              rounded="lg"
              density="comfortable"
              color="primary"
              prepend-inner-icon="mdi-account-outline"
              :rules="[v => !!v || '请输入用户名']"
              class="mb-2"
              hide-details="auto"
            />

            <!-- Password -->
            <v-text-field
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              label="密码"
              placeholder="输入密码"
              required
              variant="outlined"
              rounded="lg"
              density="comfortable"
              color="primary"
              prepend-inner-icon="mdi-lock-outline"
              :rules="[v => !!v || '请输入密码']"
              class="mb-2"
              hide-details="auto"
              @keydown.enter="handleLogin"
            >
              <template v-slot:append-inner>
                <v-icon
                  :icon="showPassword ? 'mdi-eye-off-outline' : 'mdi-eye-outline'"
                 @click="showPassword = !showPassword"
                  class="cursor-pointer"
                  size="20"
                />
              </template>
            </v-text-field>

            <!-- Error Alert -->
            <v-alert
              v-if="errorMessage"
              type="error"
              variant="tonal"
              rounded="lg"
              class="mt-4 mb-4"
              density="compact"
            >
              {{ errorMessage }}
            </v-alert>

            <!-- Submit Button -->
            <v-btn
              type="submit"
              color="primary"
              size="large"
              block
              rounded="lg"
              :loading="isSubmitting"
              :disabled="isSubmitting"
              class="mt-2 login-btn"
              @click="handleLogin"
            >
              {{ isSubmitting ? '登录中...' : '登 录' }}
            </v-btn>
          </v-form>
        </v-card-text>
      </v-card>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useCookies } from 'vue3-cookies';
import { useRouter } from 'vue-router';

const { cookies } = useCookies();
const router = useRouter();

const username = ref('');
const password = ref('');
const showPassword = ref(false);
const isSubmitting = ref(false);
const errorMessage = ref('');
const loginForm = ref(null);

const handleLogin = async () => {
  if (!username.value.trim() || !password.value.trim()) {
    errorMessage.value = '请输入用户名和密码';
    return;
  }

  isSubmitting.value = true;
  errorMessage.value = '';

  try {
    const formData = new FormData();
    formData.append('username', username.value);
    formData.append('password', password.value);

    const response = await fetch('/login', {
      method: 'POST',
      body: formData
    });

    if (!response.ok) {
      throw new Error('登录失败，请检查用户名和密码');
    }

    const data = await response.json();
    if (data.token) {
      cookies.set('jwt_token', data.token, '1h');
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
.login-page {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #0f172a 100%);
  overflow: hidden;
  z-index: 1;
}

.bg-animation {
  position: absolute;
  inset: 0;
  overflow: hidden;
}

.blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
  animation: float 20s ease-in-out infinite;
}

.blob-1 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, #3b82f6, transparent);
  top: -100px;
  left: -100px;
  animation-delay: 0s;
}

.blob-2 {
  width: 350px;
  height: 350px;
  background: radial-gradient(circle, #8b5cf6, transparent);
  bottom: -80px;
  right: -80px;
  animation-delay: -7s;
}

.blob-3 {
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, #06b6d4, transparent);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: -14s;
}

@keyframes float {
  0%, 100% { transform: translate(0, 0) scale(1); }
  25% { transform: translate(30px, -50px) scale(1.1); }
  50% { transform: translate(-20px, 20px) scale(0.95); }
  75% { transform: translate(50px, 30px) scale(1.05); }
}

.login-container {
  position: relative;
  z-index: 2;
  width: 100%;
  max-width: 440px;
  padding: 16px;
}

.login-card {
  background: rgba(255, 255, 255, 0.95) !important;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.logo-wrapper {
  width: 72px;
  height: 72px;
  margin: 0 auto;
  background: linear-gradient(135deg, #3b82f6, #8b5cf6);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 32px rgba(59, 130, 246, 0.3);
}

.login-btn {
  font-weight: 600;
  letter-spacing: 0.5px;
  height: 48px;
}
</style>
