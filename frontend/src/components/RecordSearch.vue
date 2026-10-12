<script setup lang="ts">
const query = defineModel<string>({default:''})
defineProps<{count:number}>()
</script>
<template>
 <div class="record-search">
  <label><span>{{ $t('recordSearch.label') }}</span><input v-model="query" type="search" :placeholder="$t('recordSearch.placeholder')"></label>
  <button v-if="query" class="action-button" @click="query=''">{{ $t('recordSearch.clear') }}</button>
  <p v-if="query" role="status">{{ $t('recordSearch.count',{count}) }}</p>
  <small>{{ $t('recordSearch.scope') }}</small>
 </div>
</template>
<style scoped>
.record-search {display:flex;align-items:flex-end;flex-wrap:wrap;gap:10px 14px;color:var(--ui-muted)}
label {display:grid;gap:6px;flex:1;min-width:min(240px,100%);color:var(--ui-text)}
input {min-height:44px;width:100%;padding:10px 12px;border:1px solid var(--ui-border);border-radius:8px;background:var(--ui-input);color:var(--ui-text);font:inherit}
input:focus-visible {outline:2px solid var(--ui-focus);outline-offset:2px}
p {padding-bottom:10px} small {flex-basis:100%;font-size:.85em}
</style>
