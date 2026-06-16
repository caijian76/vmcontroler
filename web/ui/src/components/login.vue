<template>
  <v-container class="justify-center w-50"  >
    <div class="min-h-screen bg-gradient-to-br from-blue-500 via-purple-500 to-indigo-600 flex items-center justify-center" >


      <v-card  elevation="12" class="mx-auto max-w-md w-full mx-4 bg-white/95 backdrop-blur-sm rounded-2xl">
        <v-card-title class="text-center pb-10">
          <div class="flex flex-col items-center">
            <div class="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center mb-4 shadow-lg">
              <v-icon size="32" class="white--text">mdi-server</v-icon>
            </div>
            <span class="text-h5 text-gray-500 mt-1">虚拟机管理平台</span>
          </div>
        </v-card-title>

        <v-card-text class="mx-auto max-w-md w-full mx-4 bg-white/95 backdrop-blur-sm rounded-2xl w-66">
          <v-form @submit.prevent="handleLogin" ref="loginForm">
            <v-text-field 
              v-model="username"
              label="用户名"
              placeholder="请输入用户名"
              required
              autofocus
              outlined
              rounded-lg
              color="blue"
              prepend-inner-icon="mdi-account"
              class="transition-all duration-600 px-10 "
            />
            <v-text-field
              v-model="password"
              label="密码"
              type="password"
              placeholder="请输入密码"
              required
              outlined
              rounded-lg
              color="blue"
              prepend-inner-icon="mdi-lock"
              class="mt-4 transition-all duration-300 px-10"
              @keydown.enter="handleLogin"
            />
          </v-form>
          <v-alert
            v-if="errorMessage"
            type="error"
            border="left"
            class="mt-4"
            dense
            transition="slide-y-reverse-transition"
          >
            {{ errorMessage }}
          </v-alert>
        </v-card-text>

        <v-card-actions class="pb-6 px-6 justify-center">
          <v-btn
            type="submit"
            color="blue"
            variant="elevated"
            class="w-full rounded-lg text-white font-medium"
            :disabled="isSubmitting"
            :loading="isSubmitting"
            @click="handleLogin"
          >
            {{ isSubmitting ? '登录中...' : '登录' }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </div>
  </v-container>
</template>

<script setup>
import { ref } from 'vue';
import { useCookies } from 'vue3-cookies';
import { useRouter } from 'vue-router';

const { cookies } = useCookies();
const router = useRouter();

const username = ref('');
const password = ref('');
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
