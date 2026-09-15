<template>
  <template v-if="isLoading">
    <slot name="loading"></slot>
  </template>

  <template v-else-if="hasAccess">
    <slot :user-info="resolvedUserInfo" :vip="vipLevel"></slot>
  </template>

  <template v-else>
    <slot name="fallback" :error="error"></slot>
  </template>
</template>

<script>
import { computed, onMounted, ref } from 'vue';
import config from '../config';

export default {
  name: 'VipOnly',

  props: {
    // 推荐传入已有的用户信息，避免重复请求学生信息接口。
    userInfo: {
      type: Object,
      default: null
    },
    // 也可以只传 vip 等级：<VipOnly :vip="userInfo.vip">...</VipOnly>
    vip: {
      type: [Number, String],
      default: null
    },
    // 未传 userInfo/vip 时，是否自动请求当前用户信息。
    autoLoad: {
      type: Boolean,
      default: true
    }
  },

  setup(props) {
    const fetchedUserInfo = ref(null);
    const isLoading = ref(false);
    const error = ref(null);

    const resolvedUserInfo = computed(() => props.userInfo || fetchedUserInfo.value);

    const vipLevel = computed(() => {
      const value = props.userInfo?.vip
        ?? props.vip
        ?? fetchedUserInfo.value?.vip
        ?? config.VIP_NONE;
      const numericValue = Number(value);
      return Number.isFinite(numericValue) ? numericValue : config.VIP_NONE;
    });

    const hasAccess = computed(() => vipLevel.value >= config.VIP_SUPER);

    const loadUserInfo = async () => {
      isLoading.value = true;
      error.value = null;
      try {
        const response = await fetch(`${config.base_url}/home/studentInfo`, {
          method: 'GET',
          credentials: 'include'
        });
        const data = await response.json();
        if (!response.ok || !data?.status) {
          throw new Error(data?.msg || '获取用户信息失败');
        }
        fetchedUserInfo.value = data.data || null;
      } catch (requestError) {
        error.value = requestError;
        fetchedUserInfo.value = null;
      } finally {
        isLoading.value = false;
      }
    };

    onMounted(() => {
      if (!props.userInfo && props.vip === null && props.autoLoad) {
        loadUserInfo();
      }
    });

    return {
      error,
      hasAccess,
      isLoading,
      resolvedUserInfo,
      vipLevel
    };
  }
};
</script>
