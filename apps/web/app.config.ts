import {
  crmButtonUi,
  crmInputUi,
  crmModalUi,
  crmPaginationUi,
  crmSelectMenuUi,
  crmSelectUi,
  crmTableUi,
} from '@crm/ui-kit'

export default defineAppConfig({
  ui: {
    button: crmButtonUi,
    input: crmInputUi,
    select: crmSelectUi,
    selectMenu: crmSelectMenuUi,
    table: crmTableUi,
    pagination: crmPaginationUi,
    modal: crmModalUi,
  },
})
