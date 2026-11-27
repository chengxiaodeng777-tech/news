<script lang="ts" setup>
import type { MenuRecord } from '#/api/system';

import { computed, nextTick, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';

import { useVbenForm } from '#/adapter/form';
import { createMenu, updateMenu } from '#/api/system';

import { normalizeStatus } from '../../shared';
import { useFormSchema } from '../data';

const emit = defineEmits<{ success: [] }>();

const formData = ref<MenuRecord>();

const [Form, formApi] = useVbenForm({
  commonConfig: {
    formItemClass: 'col-span-2 md:col-span-1',
  },
  schema: useFormSchema(),
  showDefaultActions: false,
  wrapperClass: 'grid-cols-2 gap-x-4',
});

const [Drawer, drawerApi] = useVbenDrawer({
  async onConfirm() {
    const { valid } = await formApi.validate();
    if (!valid) return;

    const values = await formApi.getValues<Record<string, any>>();
    const payload = {
      component: values.component,
      externalLink: values.externalLink,
      hidden: !!values.hidden,
      icon: values.icon,
      keepAlive: !!values.keepAlive,
      name: values.name,
      parentId: values.parentId ?? 0,
      path: values.path,
      permissionCode: values.permissionCode,
      redirect: values.redirect,
      remark: values.remark,
      sort: values.sort ?? 0,
      status: normalizeStatus(values.status),
      title: values.title,
      type: values.type,
    };

    drawerApi.lock();
    try {
      await (formData.value?.id
        ? updateMenu(formData.value.id, payload)
        : createMenu(payload));
      emit('success');
      drawerApi.close();
    } finally {
      drawerApi.unlock();
    }
  },
  async onOpenChange(isOpen) {
    if (!isOpen) return;
    const data = drawerApi.getData<Partial<MenuRecord>>() as
      | Partial<MenuRecord>
      | undefined;
    formData.value = data?.id ? (data as MenuRecord) : undefined;
    formApi.resetForm();
    await nextTick();
    formApi.setValues({
      ...data,
      parentId: data?.parentId || undefined,
      sort: data?.sort ?? 0,
      status: normalizeStatus(data?.status),
      type: data?.type ?? 'menu',
    });
  },
});

const title = computed(() => (formData.value?.id ? '编辑菜单' : '新增菜单'));
</script>

<template>
  <Drawer class="w-full max-w-200" :title="title">
    <Form class="mx-4" />
  </Drawer>
</template>
