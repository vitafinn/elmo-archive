<script setup lang="ts">
import {Quit} from "../../wailsjs/runtime";
</script>

<template>
  <button class="window-control close-control" :aria-label="$t('appearance.close')" :title="$t('appearance.close')" @click="Quit"><i-mdi-close/></button>
</template>
