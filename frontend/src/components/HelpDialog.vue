<script setup lang="ts">
const visible = defineModel({type:Boolean,required:true})
</script>
<template>
  <el-dialog v-model="visible" class="import-dialog" width="680" :title="$t('guide.help')">
    <div class="import-guide">
      <section class="guide-step"><h3>{{ $t('guide.helpImport') }}</h3><p>{{ $t('guide.helpImportBody') }}</p></section>
      <section class="guide-step"><h3>{{ $t('guide.helpStats') }}</h3>
        <dl class="glossary"><template v-for="term in ['pity','average','spent','outcome','outfit']" :key="term"><dt>{{ $t(`guide.terms.${term}`) }}</dt><dd>{{ $t(`guide.definitions.${term}`) }}</dd></template></dl>
      </section>
      <p>{{ $t('guide.partialHistory') }}</p>
    </div>
  </el-dialog>
</template>
