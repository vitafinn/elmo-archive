<script setup lang="ts">
defineProps({
  title: {type: String, required: true},
  value: {type: Number, required: false},
  note: {type: String, required: false},
})
</script>

<template>
  <div class="stat-card">
    <div class="stat-label">{{ title }}</div>
    <div v-if="value!==undefined" class="stat-value">{{ (Math.round(value * 10) / 10).toLocaleString() }}</div>
    <div v-if="note" :class="value!==undefined?'stat-note':'stat-value'">{{ note }}</div>
  </div>
</template>
