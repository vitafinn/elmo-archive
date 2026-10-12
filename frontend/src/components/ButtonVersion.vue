<script setup lang="ts">
import {onMounted, ref} from 'vue'
import {GetVersion} from '../../wailsjs/go/handler/App'
import {BrowserOpenURL} from '../../wailsjs/runtime/runtime'
const version=ref('dev')
// Set by the release workflow from its actual owner/repository, not a personal handle.
const repository=import.meta.env.VITE_REPOSITORY_URL || ''
const releases=/^https:\/\/github\.com\/[A-Za-z0-9-]+\/[A-Za-z0-9_.-]+\/?$/.test(repository) ? repository.replace(/\/$/,'')+'/releases' : ''
onMounted(async()=>{version.value=await GetVersion() || 'dev'})
</script>
<template>
 <button v-if="releases" class="text-link ml-2" @click="BrowserOpenURL(releases)">{{ version }} · GitHub Releases ↗</button>
 <span v-else class="detail-note ml-2">{{ version }}</span>
</template>
