import { createRouter, createWebHistory } from 'vue-router'

import Login from './components/Login.vue'
import Regist from './components/Regist.vue'
import Home from './components/Home.vue'
import PasswordForget from './components/PasswordForget.vue'
import StudentHome from './components/StudentHome.vue'
import History from './components/History.vue'
import AssessmentSubmission from './components/AssessmentSubmission.vue'
import SystemBoard from './components/SystemBoard.vue'
import Ranking from './components/Ranking.vue'
// 可以先导入一个空组件作为其他页面的占位

const routes = [
  {
    path: '/',
    redirect: '/login', // 修正重定向写法，不是restrict
    component: Login
  },
  {
    path: '/login',
    component: Login // 登录页路由
  },
  {
    path: '/regist',
    component: Regist // 注册页路由
  },
  {
    path: '/passwordforget',
    component: PasswordForget
  },
  {
    path: '/home',
    component: Home,
    children: [
      { path: 'student-home', component: StudentHome }, // 个人中心
      { path: 'history', component: History }, // 历史记录
      { path: 'assessment', component: AssessmentSubmission }, // 学生考核
      { path: 'ranking', component: Ranking }, // 排行榜
      { path: 'settings', component: SystemBoard } // 系统设置
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
