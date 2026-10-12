<script setup lang="ts">
import {usePoolStore} from "../store/poolStore.ts";

const props = defineProps({poolType: {type: Number, required: true}})
const poolStore = usePoolStore()

const handleClick = () => {
  poolStore.poolType = props.poolType
}
</script>

<template>
  <button type="button" class="pool-nav-item" :class="{'is-selected': poolStore.poolType===poolType}" :aria-current="poolStore.poolType===poolType ? 'page' : undefined" @click="handleClick">
    {{ $te(`gacha.type.${poolType}`) ? $t(`gacha.type.${poolType}`) : $t('gacha.unknownType', {id: poolType}) }}
  </button>
</template>
