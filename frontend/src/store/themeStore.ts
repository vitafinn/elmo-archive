import {defineStore} from 'pinia'
import {ref} from 'vue'
export const useThemeStore = defineStore('theme', () => {
  const gradientFrom=ref(localStorage.getItem('themeFrom') || '#18212d')
  const gradientTo=ref(localStorage.getItem('themeTo') || '#18212d')
  const savedAppearance=localStorage.getItem('appearance')
  const appearance=ref(savedAppearance && ['jade','slate','contrast'].includes(savedAppearance) ? savedAppearance : 'jade')
  const textSize=ref(localStorage.getItem('textSize') === 'large' ? 'large' : 'comfortable')
  const apply=()=>{
    const root=document.documentElement
    root.dataset.appearance=appearance.value
    root.style.setProperty('--text-scale',textSize.value === 'large' ? '1.15' : '1')
    root.style.setProperty('--theme-from',gradientFrom.value)
    root.style.setProperty('--theme-to',gradientTo.value)
  }
  const updateAppearance=(mode:string,size:string)=>{
    appearance.value=['jade','slate','contrast'].includes(mode) ? mode : 'jade'
    textSize.value=size === 'large' ? 'large' : 'comfortable'
    localStorage.setItem('appearance',appearance.value)
    localStorage.setItem('textSize',textSize.value)
    apply()
  }
  const updateTheme=(from:string,to:string)=>{
    gradientFrom.value=from;gradientTo.value=to
    localStorage.setItem('themeFrom',from);localStorage.setItem('themeTo',to);apply()
  }
  apply()
  return {gradientFrom,gradientTo,appearance,textSize,updateTheme,updateAppearance}
})
