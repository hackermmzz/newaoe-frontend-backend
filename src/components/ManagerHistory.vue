<template>
  <History v-bind="$props" />
</template>

<script>
import History from './History.vue';
import config from '../config';

// 管理平台历史页复用普通历史页的完整模板和逻辑，
// 由路由传入只读状态、接口地址和学生筛选参数。
export default {
  name: 'ManagerHistory',
  components: {
    History
  },
  props: {
    readOnly: {
      type: Boolean,
      default: true
    },
    historyUrl: {
      type: String,
      default: config.manager_history_url
    },
    historyTitle: {
      type: String,
      default: '学生提交历史'
    },
    backPath: {
      type: String,
      default: '/home/manager/student-statistics'
    },
    requestParams: {
      type: Object,
      default: () => ({})
    }
  }
};
</script>
