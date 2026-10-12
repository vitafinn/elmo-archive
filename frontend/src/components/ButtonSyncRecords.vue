<script setup lang="ts">
import {PrepareImport} from '../../wailsjs/go/handler/App'
import {useImportStore} from '../store/importStore'
const flow=useImportStore()
const syncRecords=async()=>{
 if(flow.busy || flow.preparing || flow.confirming)return
 flow.preparing=true;flow.error='';flow.preview=null;flow.acknowledged=false
 try{flow.preview=await PrepareImport('os');flow.confirming=true}
 catch(err){flow.error=String(err);flow.visible=true}
 finally{flow.preparing=false}
}
</script>
<template>
 <el-button type="primary" :loading="flow.busy || flow.preparing" :disabled="flow.busy || flow.preparing || flow.confirming" @click="syncRecords">{{ flow.busy ? $t('sync.loading') : $t('guide.sync') }}</el-button>
</template>
