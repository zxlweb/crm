<template>
  <div class="flex min-h-screen bg-ds-bg font-sans text-ds-fg">
    <AdminSidebar />
    <div class="flex min-w-0 flex-1 flex-col">
      <AdminTopbar />
      <main class="flex-1 overflow-auto p-6 lg:p-8">
        <slot />
      </main>
    </div>
  </div>
</template>
