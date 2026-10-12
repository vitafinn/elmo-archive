<script setup lang="ts">
import {computed} from 'vue'
import {useUserStore} from '../store/userStore'
const userStore = useUserStore()
const selected = computed(() => userStore.userList.find(user => user.id === userStore.userId))
</script>
<template>
  <el-dropdown v-if="selected" style="--wails-draggable:no-drag">
    <button class="account-button" :aria-label="$t('guide.switchAccount')">{{ $t(`server.${selected.server}`) }} · {{ selected.uid }} <span aria-hidden="true">⌄</span></button>
    <template #dropdown><el-dropdown-menu>
      <el-dropdown-item v-for="user in userStore.userList" :key="user.id" @click="userStore.updateUserId(user.id)">
        <span>{{ user.id === userStore.userId ? '✓ ' : '' }}{{ $t(`server.${user.server}`) }} · {{ user.uid }}</span>
      </el-dropdown-item>
    </el-dropdown-menu></template>
  </el-dropdown>
</template>
