import { crmInputUi } from './nuxt-ui-input'

/**
 * Nuxt UI USelectMenu — 用 Headless UI Listbox 渲染的可自定义下拉。
 * 触发器复用 input 风格；菜单弹层与 design-system 对齐。
 */
export const crmSelectMenuUi = {
  // 触发器（外层）样式直接复用 input 的视觉
  ...crmInputUi,
  // 菜单容器
  container: 'z-30 group',
  trigger: 'flex items-center w-full',
  width: 'w-full',
  height: 'max-h-72',
  base: 'relative focus:outline-none overflow-y-auto scroll-py-1',
  background: 'bg-ds-bg-elevated',
  shadow: 'shadow-ds-lg',
  rounded: 'rounded-xl',
  padding: 'p-1',
  ring: 'ring-1 ring-ds-border',
  empty: 'text-sm text-ds-fg-muted px-3 py-2',
  // 单条 option
  option: {
    base: 'cursor-pointer select-none relative flex items-center justify-between gap-2 transition-colors duration-150',
    rounded: 'rounded-lg',
    padding: 'px-2.5 py-1.5',
    size: 'text-sm',
    color: 'text-ds-fg',
    container: 'flex items-center gap-2 min-w-0',
    active: 'bg-ds-brand-subtle/60 text-ds-fg-brand',
    inactive: 'hover:bg-ds-bg-muted',
    selected: 'pe-7 font-semibold text-ds-fg-brand',
    disabled: 'cursor-not-allowed opacity-50',
    empty: 'text-sm text-ds-fg-muted px-3 py-2',
    icon: {
      base: 'flex-shrink-0 h-4 w-4',
      active: 'text-ds-fg-brand',
      inactive: 'text-ds-fg-subtle',
    },
    selectedIcon: {
      wrapper: 'absolute inset-y-0 end-0 flex items-center',
      padding: 'pe-2',
      base: 'h-4 w-4 text-ds-fg-brand flex-shrink-0',
    },
    avatar: {
      base: 'flex-shrink-0',
      size: '2xs' as const,
    },
    chip: {
      base: 'flex-shrink-0 w-2 h-2 mx-1 rounded-full',
    },
  },
  transition: {
    leaveActiveClass: 'transition ease-in duration-150',
    leaveFromClass: 'opacity-100',
    leaveToClass: 'opacity-0',
  },
  popper: {
    placement: 'bottom-start' as const,
    strategy: 'absolute' as const,
  },
  default: {
    selectedIcon: 'i-heroicons-check-20-solid',
    trailingIcon: 'i-heroicons-chevron-down-20-solid',
    clearSearchOnClose: false,
    showCreateOptionWhen: 'empty' as const,
    searchablePlaceholder: { label: 'Search...' },
    empty: { label: 'No options.' },
    optionEmpty: { label: 'No results for "{query}".' },
  },
  arrow: {
    base: 'invisible before:visible before:block before:rotate-45 before:z-[-1] before:w-2 before:h-2',
    rounded: 'before:rounded-sm',
    shadow: 'before:shadow',
    ring: 'before:ring-1 before:ring-ds-border',
    background: 'before:bg-ds-bg-elevated',
    placement: "group-data-[popper-placement*='right']:-left-1 group-data-[popper-placement*='left']:-right-1 group-data-[popper-placement*='top']:-bottom-1 group-data-[popper-placement*='bottom']:-top-1",
  },
} as const
