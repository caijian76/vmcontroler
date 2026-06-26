<template>
  <v-sheet
    class="d-flex align-center justify-center"
    height="100vh"
    :style="{
      background: 'linear-gradient(135deg, #1a73e8 0%, #6c3bbf 50%, #1a237e 100%)',
    }"
  >
    <v-container class="d-flex align-center justify-center">
      <v-row justify="center">
        <v-col cols="12" sm="8" md="5" lg="4">
          <v-card
            elevation="16"
            class="rounded-xl"
            :style="{
              background: 'rgba(255, 255, 255, 0.92)',
              backdropFilter: 'blur(20px)',
            }"
          >
            <v-card-item class="pt-10 pb-4">
              <template v-slot:prepend>
                <v-avatar
                  size="64"
                  rounded="lg"
                  :style="{
                    background: 'linear-gradient(135deg, #1a73e8, #6c3bbf)',
                  }"
                  class="elevation-4"
                >
                  <v-icon size="36" color="white">mdi-server</v-icon>
                </v-avatar>
              </template>
              <template v-slot:title>
                <span class="text-h5 font-weight-bold">虚拟机管理平台</span>
              </template>
              <template v-slot:subtitle>
                <span class="text-body-2 text-grey-darken-1">请使用您的凭据登录</span>
              </template>
            </v-card-item>

            <v-card-text class="px-8 pb-2">
              <v-form @submit.prevent="handleLogin" ref="loginForm" validate-on="submit">
                <v-text-field
                  v-model="username"
                  label="用户名"
                  placeholder="请输入用户名"
                  variant="outlined"
                  density="comfortable"
                  color="primary"
                  prepend-inner-icon="mdi-account-outline"
                  :rules="[v => !!v || '请输入用户名']"
                  class="mb-3"
                  autofocus
                />

                <v-text-field
                  v-model="password"
                  label="密码"
                  placeholder="请输入密码"
                  type="password"
                  variant="outlined"
                  density="comfortable"
                  color="primary"
                  prepend-inner-icon="mdi-lock-outline"
                  :rules="[v => !!v || '请输入密码']"
                  class="mb-2"
                  @keydown.enter="handleLogin"
                />

                <v-alert
                  v-if="errorMessage"
                  type="error"
                  variant="tonal"
                  density="compact"
                  border="start"
                  class="mb-4 mt-2"
                  closable
                  @click:close="errorMessage = ''"
                >
                  {{ errorMessage }}
                </v-alert>
              </v-form>
            </v-card-text>

            <v-card-actions class="px-8 pb-10">
              <v-btn
                color="primary"
                size="large"
                block
                rounded="lg"
                variant="flat"
                :disabled="isSubmitting"
                :loading="isSubmitting"
                @click="handleLogin"
                :style="{
                  background: 'linear-gradient(135deg, #1a73e8, #6c3bbf)',
                }"
              >
                <v-icon start v-if="!isSubmitting">mdi-login</v-icon>
                {{ isSubmitting ? '登录中...' : '登 录' }}
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-sheet>
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
  const { valid } = await loginForm.value.validate();
  if (!valid) return;

  isSubmitting.value = true;
  errorMessage.value = '';

  try {
    const formData = new FormData();
    formData.append('username', username.value);
    formData.append('password', password.value);

    const response = await fetch('/login', {
      method: 'POST',
      body: formData,
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