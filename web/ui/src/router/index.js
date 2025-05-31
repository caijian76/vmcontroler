import { createRouter, createWebHistory } from 'vue-router';
import Login from '../components/login.vue';
import VmList from '../components/vmlist.vue';

const routes = [
  {
    path: '/',
    name: 'Login',
    component: Login
  },
  {
    path: '/vmlist',
    name: 'VmList',
    component: VmList
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;
