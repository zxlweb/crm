import { crmInputUi } from './nuxt-ui-input'

/** Nuxt UI USelect — 与 Input 同视觉语言 */
export const crmSelectUi = {
  ...crmInputUi,
  form: 'form-select',
  trailingIcon: 'i-heroicons-chevron-down-20-solid',
  default: {
    ...crmInputUi.default,
    trailingIcon: 'i-heroicons-chevron-down-20-solid',
  },
} as const
