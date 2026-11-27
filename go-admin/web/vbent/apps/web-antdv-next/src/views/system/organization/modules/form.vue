<script lang="ts" setup>
import type { OrganizationRecord } from '#/api/system';

import { computed, nextTick, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { useVbenForm } from '#/adapter/form';
import { createOrganization, updateOrganization } from '#/api/system';

import { normalizeStatus } from '../../shared';
import { useFormSchema } from '../data';

const emit = defineEmits<{ success: [] }>();

const formData = ref<OrganizationRecord>();

const [Form, formApi] = useVbenForm({
  layout: 'vertical',
  schema: useFormSchema(),
  showDefaultActions: false,
});

const [Modal, modalApi] = useVbenModal({
  async onConfirm() {
    const { valid } = await formApi.validate();
    if (!valid) return;

    const values = await formApi.getValues<Record<string, any>>();
    const payload = {
      code: values.code,
      email: values.email,
      leader: values.leader,
      name: values.name,
      parentId: values.parentId ?? 0,
      phone: values.phone,
      remark: values.remark,
      sort: values.sort ?? 0,
      status: normalizeStatus(values.status),
      type: values.type,
    };

    modalApi.lock();
    try {
      await (formData.value?.id
        ? updateOrganization(formData.value.id, payload)
        : createOrganization(payload));
      emit('success');
      modalApi.close();
    } finally {
      modalApi.lock(false);
    }
  },
  async onOpenChange(isOpen) {
    if (!isOpen) return;
    const data = modalApi.getData<Partial<OrganizationRecord>>() as
      | Partial<OrganizationRecord>
      | undefined;
    formData.value = data?.id ? (data as OrganizationRecord) : undefined;
    formApi.resetForm();
    await nextTick();
    formApi.setValues({
      ...data,
      parentId: data?.parentId || undefined,
      sort: data?.sort ?? 0,
      status: normalizeStatus(data?.status),
      type: data?.type ?? 'department',
    });
  },
});

const title = computed(() => (formData.value?.id ? '编辑组织' : '新增组织'));
</script>

<template>
  <Modal class="w-160" :title="title">
    <Form class="mx-4" />
  </Modal>
</template>
