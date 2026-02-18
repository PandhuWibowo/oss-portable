import { ref, onMounted } from 'vue'

export function useTheme() {
    const isLight = ref(true)

    function applyTheme(light) {
        if (light) {
            document.documentElement.setAttribute('data-theme', 'light')
        } else {
            document.documentElement.removeAttribute('data-theme')
        }
    }

    function toggleTheme() {
        isLight.value = !isLight.value
        applyTheme(isLight.value)
        try { localStorage.setItem('theme', isLight.value ? 'light' : 'dark') } catch { }
    }

    onMounted(() => {
        try {
            const saved = localStorage.getItem('theme')
            isLight.value = saved !== 'dark'
        } catch {
            isLight.value = true
        }
        applyTheme(isLight.value)
    })

    return { isLight, toggleTheme }
}
