<script setup lang="ts">
import {WindowMinimise} from "../../wailsjs/runtime";
</script>

<template>
  <button class="window-control" :aria-label="$t('appearance.minimize')" :title="$t('appearance.minimize')" @click="WindowMinimise"><i-mdi-minus/></button>
</template>
