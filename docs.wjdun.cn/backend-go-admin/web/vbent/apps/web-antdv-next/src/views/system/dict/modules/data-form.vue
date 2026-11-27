<script lang="ts" setup>
import type { DictDataRecord } from '#/api/system';

import { computed, nextTick, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import { useVbenForm } from '#/adapter/form';
import { createDictData, updateDictData } from '#/api/system';

import { normalizeStatus } from '../../shared';
import { useDataFormSchema } from '../data';

const emit = defineEmits<{ success: [] }>();

const formData = ref<DictDataRecord>();
const dictTypeId = ref<number>();

const [Form, formApi] = useVbenForm({
  layout: 'vertical',
  schema: useDataFormSchema(),
  showDefaultActions: false,
});

const [Modal, modalApi] = useVbenModal({
  async onConfirm() {
    const { valid } = await formApi.validate();
    if (!valid || !dictTypeId.value) return;

    const values = await formApi.getValues<Record<string, any>>();
    const payload = {
      cssClass: values.cssClass,
      dictTypeId: dictTypeId.value,
      isDefault: !!values.isDefault,
      label: values.label,
      remark: values.remark,
      sort: values.sort ?? 0,
      status: normalizeStatus(values.status),
      tagType: values.tagType,
      value: values.value,
    };

    modalApi.lock();
    try {
      await (formData.value?.id
        ? updateDictData(formData.value.id, payload)
        : createDictData(payload));
      emit('success');
      modalApi.close();
    } finally {
      modalApi.lock(false);
    }
  },
  async onOpenChange(isOpen) {
    if (!isOpen) return;
    const data = modalApi.getData<
      Partial<DictDataRecord> & { dictTypeId: number }
    >() as (Partial<DictDataRecord> & { dictTypeId: number }) | undefined;
    formData.value = data?.id ? (data as DictDataRecord) : undefined;
    dictTypeId.value = data?.dictTypeId;
    formApi.resetForm();
    await nextTick();
    formApi.setValues({
      ...data,
      isDefault: !!data?.isDefault,
      sort: data?.sort ?? 0,
      status: normalizeStatus(data?.status),
    });
  },
});

const title = computed(() =>
  formData.value?.id ? '编辑字典数据' : '新增字典数据',
);
</script>

<template>
  <Modal class="w-140" :title="title">
    <Form class="mx-4" />
  </Modal>
</template>
