import { createRouter, createWebHistory } from 'vue-router'
import axios from 'axios'

import Login from './components/Login.vue'
import Regist from './components/Regist.vue'
import Home from './components/Home.vue'
import PasswordForget from './components/PasswordForget.vue'
import StudentHome from './components/StudentHome.vue'
import History from './components/History.vue'
import ManagerHistory from './components/ManagerHistory.vue'
import AssessmentSubmission from './components/AssessmentSubmission.vue'
import SystemBoard from './components/SystemBoard.vue'
import Ranking from './components/Ranking.vue'
import ManagerPage from './components/ManagerPage.vue'
import config from './config'
// 可以先导入一个空组件作为其他页面的占位

const routes = [
  {
    path: '/',
    redirect: '/login', // 修正重定向写法，不是restrict
    component: Login,
    meta: { public: true }
  },
  {
    path: '/login',
    component: Login, // 登录页路由
    meta: { public: true }
  },
  {
    path: '/regist',
    component: Regist, // 注册页路由
    meta: { public: true }
  },
  {
    path: '/passwordforget',
    component: PasswordForget,
    meta: { public: true }
  },
  {
    path: '/home',
    component: Home,
    children: [
      { path: 'student-home', component: StudentHome }, // 个人中心
      { path: 'history', component: History }, // 历史记录
      { path: 'assessment', component: AssessmentSubmission }, // 学生考核
      { path: 'ranking', component: Ranking }, // 排行榜
      { path: 'settings', component: SystemBoard }, // 系统设置
      {
        path: 'manager',
        component: ManagerPage,
        children: [
          {
            path: '',
            redirect: { name: 'StudentStatistics' }
          },
          {
            path: 'student-statistics',
            name: 'StudentStatistics',
            component: () => import('./components/StudentStatistics.vue')
          },
          {
            path: 'student-history/:studentId',
            name: 'ManagerStudentHistory',
            component: ManagerHistory,
            props: route => {
              const studentId = String(route.params.studentId || '').trim()
              return {
                readOnly: true,
                historyUrl: config.manager_history_url,
                historyTitle: `${studentId} 的提交历史`,
                backPath: '/home/manager/student-statistics',
                requestParams: { student_id: studentId }
              }
            }
          },
          {
            path: 'all-history',
            name: 'ManagerAllHistory',
            component: ManagerHistory,
            props: {
              readOnly: true,
              historyUrl: config.manager_all_history_url,
              historyTitle: '所有人的提交记录',
              backPath: '/home/manager/student-statistics',
              showStudentId: true,
              recordsPerPage: config.ManagerHistoryRecordPerPage,
              requestParams: {}
            }
          },
          {
            path: 'assessment-history',
            name: 'ManagerAssessmentHistory',
            component: ManagerHistory,
            props: {
              readOnly: true,
              historyUrl: config.manager_assessment_history_url,
              rerunAnomalRecordUrl: '',
              blockAllSubmitUrl: '',
              cancelSubmitBlockUrl: '',
              historyTitle: '所有考核提交记录',
              backPath: '/home/manager/student-statistics',
              showStudentId: true,
              recordsPerPage: config.ManagerHistoryRecordPerPage,
              requestParams: {}
            }
          },
          {
            path: 'feedback',
            name: 'ManagerFeedback',
            component: () => import('./components/FeedbackRecords.vue')
          }
        ]
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

/**
 * 校验当前 Cookie 是否仍然有效。
 *
 * Cookie 使用 HttpOnly 属性时，前端不能直接读取它，只能通过后端接口
 * 判断登录状态。后端约定 Cookie 失效时返回 401。
 */
async function checkLoginStatus() {
  try {
    const response = await fetch(`${config.base_url}/checkLoginStatus`, {
      credentials: 'include'
    })

    // 后端约定 Cookie 有效时返回 2xx，失效时返回 401。
    return response?.status !== 401
  } catch {
    // 网络异常时不允许进入需要登录的页面
    return false
  }
}

router.beforeEach(async (to) => {
  // 登录、注册和找回密码页面不需要登录校验
  if (to.meta.public) {
    return true
  }

  if (await checkLoginStatus()) {
    return true
  }

  return {
    path: '/login',
    query: {
      redirect: to.fullPath
    }
  }
})

// 处理接口请求过程中 Cookie 失效的情况。
// 项目中的 Axios 请求会共用这个全局响应拦截器。
let redirectingToLogin = false

function redirectToLogin() {
  const currentRoute = router.currentRoute.value

  // 登录页不再重复跳转；多个请求同时返回 401 时只处理一次。
  if (redirectingToLogin || currentRoute.meta.public) {
    return
  }

  redirectingToLogin = true
  router.replace({
    path: '/login',
    query: {
      redirect: currentRoute.fullPath
    }
  }).catch(() => {
    // 忽略重复导航或被其他导航取消的错误
  }).finally(() => {
    redirectingToLogin = false
  })
}

axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      redirectToLogin()
    }

    return Promise.reject(error)
  }
)

export default router
