<script lang="ts" setup>
import type { DataNode } from 'antdv-next/dist/tree';

import type { RoleRecord } from '#/api/system';

import { computed, nextTick, ref } from 'vue';

import { Tree, useVbenDrawer } from '@vben/common-ui';

import { Spin } from 'antdv-next';

import { useVbenForm } from '#/adapter/form';
import {
  createRole,
  listMenuTree,
  listRoles,
  replaceRoleMenus,
  updateRole,
} from '#/api/system';

import { normalizeNumberArray, normalizeStatus } from '../../shared';
import { useFormSchema } from '../data';

const emit = defineEmits<{ success: [] }>();

const formData = ref<RoleRecord>();
const menus = ref<DataNode[]>([]);
const loadingMenus = ref(false);

const [Form, formApi] = useVbenForm({
  schema: useFormSchema(),
  showDefaultActions: false,
});

const [Drawer, drawerApi] = useVbenDrawer({
  async onConfirm() {
    const { valid } = await formApi.validate();
    if (!valid) return;

    const values = await formApi.getValues<Record<string, any>>();
    const menuIds = normalizeNumberArray(values.menuIds);

    drawerApi.lock();
    try {
      if (formData.value?.id) {
        await updateRole(formData.value.id, {
          name: values.name,
          remark: values.remark,
          sort: values.sort ?? 0,
          status: normalizeStatus(values.status),
        });
        await replaceRoleMenus(formData.value.id, menuIds);
      } else {
        await createRole({
          code: values.code,
          name: values.name,
          remark: values.remark,
          sort: values.sort ?? 0,
          status: normalizeStatus(values.status),
        });
        const result = await listRoles({
          keyword: values.code,
          page: 1,
          pageSize: 200,
        });
        const roleId = result.list.find((item) => item.code === values.code)?.id;
        if (roleId) {
          await replaceRoleMenus(roleId, menuIds);
        }
      }
      emit('success');
      drawerApi.close();
    } finally {
      drawerApi.unlock();
    }
  },
  async onOpenChange(isOpen) {
    if (!isOpen) return;

    if (menus.value.length === 0) {
      await loadMenus();
    }

    const data = drawerApi.getData<Partial<RoleRecord>>() as
      | Partial<RoleRecord>
      | undefined;
    formData.value = data?.id ? (data as RoleRecord) : undefined;
    formApi.resetForm();
    await nextTick();
    formApi.updateSchema([
      {
        componentProps: {
          disabled: !!formData.value?.id,
        },
        fieldName: 'code',
      },
    ]);
    formApi.setValues({
      ...data,
      code: data?.code ?? '',
      menuIds: normalizeNumberArray(data?.menuIds),
      sort: data?.sort ?? 0,
      status: normalizeStatus(data?.status),
    });
  },
});

async function loadMenus() {
  loadingMenus.value = true;
  try {
    menus.value = (await listMenuTree()) as unknown as DataNode[];
  } finally {
    loadingMenus.value = false;
  }
}

const title = computed(() => (formData.value?.id ? '编辑角色' : '新增角色'));
</script>

<template>
  <Drawer class="w-full max-w-180" :title="title">
    <Form class="mx-4">
      <template #menuIds="slotProps">
        <Spin :spinning="loadingMenus" :classes="{ root: 'w-full' }">
          <Tree
            bordered
            :default-expanded-level="2"
            label-field="title"
            multiple
            :tree-data="menus"
            value-field="id"
            v-bind="slotProps"
          />
        </Spin>
      </template>
    </Form>
  </Drawer>
</template>
