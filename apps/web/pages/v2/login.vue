<script setup lang="ts">
/** @deprecated 已合并至 /login，保留重定向 */
definePageMeta({ layout: false })
const theme = useTheme()
theme.setTheme('v2')
await navigateTo('/login', { replace: true })
</script>
