import { crmButtonUi, crmInputUi, crmPaginationUi, crmSelectUi, crmTableUi } from '@crm/ui-kit'

export default defineAppConfig({
  ui: {
    button: crmButtonUi,
    input: crmInputUi,
    select: crmSelectUi,
    table: crmTableUi,
    pagination: crmPaginationUi,
  },
})
