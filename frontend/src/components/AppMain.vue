<script setup lang="ts">
import PoolMain from "./PoolMain.vue";
import PoolRadio from "./PoolRadio.vue";

import {usePoolStore} from "../store/poolStore";
const poolStore = usePoolStore()
</script>

<template>
  <div class="dashboard-shell">
    <nav class="pool-navigation" :aria-label="$t('dashboard.procurement')">
      <p class="eyebrow">{{ $t('dashboard.procurement') }}</p>
      <PoolRadio v-for="poolType in poolStore.poolTypes" :key="poolType" :poolType="poolType"/>
    </nav>
    <main class="dashboard-content"><PoolMain :key="poolStore.poolType"/></main>
  </div>
</template>
