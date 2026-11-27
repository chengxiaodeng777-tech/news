<script lang="ts" setup>
import type { DictTypeRecord } from '#/api/system';

import { computed, nextTick, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { useVbenForm } from '#/adapter/form';
import { createDictType, updateDictType } from '#/api/system';

import { normalizeStatus } from '../../shared';
import { useTypeFormSchema } from '../data';

const emit = defineEmits<{ success: [] }>();

const formData = ref<DictTypeRecord>();

const [Form, formApi] = useVbenForm({
  layout: 'vertical',
  schema: useTypeFormSchema(),
  showDefaultActions: false,
});

const [Modal, modalApi] = useVbenModal({
  async onConfirm() {
    const { valid } = await formApi.validate();
    if (!valid) return;

    const values = await formApi.getValues<Record<string, any>>();
    modalApi.lock();
    try {
      await (formData.value?.id
        ? updateDictType(formData.value.id, {
            code: values.code,
            name: values.name,
            remark: values.remark,
            status: normalizeStatus(values.status),
          })
        : createDictType({
            code: values.code,
            name: values.name,
            remark: values.remark,
            status: normalizeStatus(values.status),
          }));
      emit('success');
      modalApi.close();
    } finally {
      modalApi.lock(false);
    }
  },
  async onOpenChange(isOpen) {
    if (!isOpen) return;
    const data = modalApi.getData<Partial<DictTypeRecord>>() as
      | Partial<DictTypeRecord>
      | undefined;
    formData.value = data?.id ? (data as DictTypeRecord) : undefined;
    formApi.resetForm();
    await nextTick();
    formApi.setValues({
      ...data,
      status: normalizeStatus(data?.status),
    });
  },
});

const title = computed(() =>
  formData.value?.id ? '编辑字典类型' : '新增字典类型',
);
</script>

<template>
  <Modal class="w-140" :title="title">
    <Form class="mx-4" />
  </Modal>
</template>
