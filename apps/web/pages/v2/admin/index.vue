<script setup lang="ts">
/** @deprecated 已合并至 /admin，保留重定向 */
definePageMeta({ layout: false })
const theme = useTheme()
theme.setTheme('v2')
await navigateTo('/admin', { replace: true })
</script>
