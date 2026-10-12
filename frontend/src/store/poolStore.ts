import {defineStore} from "pinia";
import {ref} from "vue";
import {handler} from "../../wailsjs/go/models.ts";
import {GetPoolInfo, GetPoolTypes, GetOutfitInfo} from "../../wailsjs/go/handler/App";
import {useUserStore} from "./userStore.ts";
import {NotifyError} from "../utils/notify.ts";
import PoolInfo = handler.PoolInfo;

export const usePoolStore = defineStore('pool', () => {
    const poolType = ref(1)
    const poolTypes = ref<number[]>([])
    const poolInfo = ref<PoolInfo>()
    const outfitInfo = ref<handler.OutfitInfo>()
    const loading = ref(false)
    const loadError = ref(false)
    const userStore = useUserStore()
    let generation = 0

    const updatePoolInfo = async (silent = false) => {
        const request = ++generation
        if (!silent) loading.value = true
        loadError.value = false
        const userId = userStore.userId || 0
        if (!silent) { poolInfo.value = undefined; outfitInfo.value = undefined }
        try {
            const types = await GetPoolTypes(userId)
            if (request !== generation) return
            poolTypes.value = types || []
            if (!poolTypes.value.includes(poolType.value)) {
                poolType.value = poolTypes.value[0] || 1
            }
            if (!userId || !poolTypes.value.length) return
            if ([9,10].includes(poolType.value)) {
                const result = await GetOutfitInfo(userId, poolType.value)
                if (request === generation) outfitInfo.value = result
                return
            }
            const result = await GetPoolInfo(userId, poolType.value)
            if (request === generation) poolInfo.value = result
        } catch (err) {
            if (request === generation) { loadError.value = true; NotifyError('Error', err) }
        } finally {
            if (request === generation) loading.value = false
        }
    }

    return {poolType, poolTypes, poolInfo, outfitInfo, loading, loadError, updatePoolInfo, init: updatePoolInfo}
})
