<script setup lang="ts">
/** @deprecated 使用 /settings?tab=roles */
definePageMeta({ layout: 'app', middleware: 'auth' })

await navigateTo({ path: '/settings', query: { tab: 'roles' } }, { replace: true })
</script>
