<script lang="ts" setup>
import type { VxeTableGridOptions } from '#/adapter/vxe-table';
import type { UserRecord } from '#/api/system';

import { onMounted, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { Plus } from '@vben/icons';

import { Button, message } from 'antdv-next';

import { useVbenVxeGrid, VbenTableAction } from '#/adapter/vxe-table';
import { deleteUser, listRoles, listUsers } from '#/api/system';

import { normalizeOptionalStatus, toGridResult } from '../shared';
import { useColumns, useGridFormSchema } from './data';
import Form from './modules/form.vue';

defineOptions({ name: 'SystemUser' });

const roles = ref(new Map<number, string>());

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  connectedComponent: Form,
  destroyOnClose: true,
});

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions: {
    schema: useGridFormSchema(),
    submitOnChange: true,
  },
  gridOptions: {
    columns: useColumns((roleIds) =>
      (roleIds ?? [])
        .map((roleId) => roles.value.get(roleId) || String(roleId))
        .join('、'),
    ),
    height: 'auto',
    keepSource: true,
    proxyConfig: {
      ajax: {
        query: async ({ page }, formValues) => {
          const result = await listUsers({
            keyword: formValues.keyword,
            page: page.currentPage,
            pageSize: page.pageSize,
            status: normalizeOptionalStatus(formValues.status),
          });
          return toGridResult(result.list, result.total);
        },
      },
    },
    rowConfig: {
      keyField: 'id',
    },
    toolbarConfig: {
      custom: true,
      export: false,
      refresh: true,
      search: true,
      zoom: true,
    },
  } as VxeTableGridOptions<UserRecord>,
});

function onCreate() {
  formDrawerApi.setData({}).open();
}

function onEdit(row: UserRecord) {
  formDrawerApi.setData(row).open();
}

async function onDelete(row: UserRecord) {
  await deleteUser(row.id);
  message.success('删除成功');
  gridApi.query();
}

async function loadRoles() {
  const result = await listRoles({ page: 1, pageSize: 200 });
  roles.value = new Map(result.list.map((role) => [role.id, role.name]));
}

onMounted(loadRoles);
</script>

<template>
  <Page auto-content-height>
    <FormDrawer @success="gridApi.query()" />
    <Grid table-title="用户管理">
      <template #toolbar-tools>
        <Button type="primary" @click="onCreate">
          <Plus class="size-5" />
          新增用户
        </Button>
      </template>
      <template #action="{ row }">
        <VbenTableAction
          :actions="[
            {
              text: '编辑',
              icon: 'lucide:edit',
              onClick: () => onEdit(row),
            },
          ]"
          :dropdown-actions="[
            {
              text: '删除',
              icon: 'lucide:trash-2',
              danger: true,
              popConfirm: {
                title: `确认删除用户 ${row.username}？`,
                confirm: () => onDelete(row),
              },
            },
          ]"
          align="center"
        />
      </template>
    </Grid>
  </Page>
</template>
