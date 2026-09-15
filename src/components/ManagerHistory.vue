<script>
import History from './History.vue';
import config from '../config';

// 管理平台历史页继承普通历史页的展示、分页、下载和结果查看能力。
// 通过默认属性切换为只读模式，并由路由传入要查看的学生 ID。
export default {
  name: 'ManagerHistory',
  extends: History,
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
