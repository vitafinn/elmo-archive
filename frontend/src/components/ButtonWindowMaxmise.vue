<script setup lang="ts">
import {WindowIsMaximised, WindowMaximise, WindowUnmaximise} from "../../wailsjs/runtime";

const handleWindowMaximise = async () => {
  await WindowIsMaximised().then(res => res ? WindowUnmaximise() : WindowMaximise())
}
</script>

<template>
  <button class="window-control" :aria-label="$t('appearance.maximize')" :title="$t('appearance.maximize')" @click="handleWindowMaximise"><i-mdi-crop-square/></button>
</template>
