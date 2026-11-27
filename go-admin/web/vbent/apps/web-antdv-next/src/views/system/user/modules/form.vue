<script lang="ts" setup>
import type { UserRecord } from '#/api/system';

import { computed, nextTick, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';

import { message } from 'antdv-next';

import { useVbenForm } from '#/adapter/form';
import {
  createUser,
  listRoles,
  listUsers,
  replaceUserRoles,
  updateUser,
} from '#/api/system';

import {
  normalizeNumberArray,
  normalizeStatus,
  omitEmpty,
} from '../../shared';
import { useFormSchema } from '../data';

const emit = defineEmits<{ success: [] }>();

const formData = ref<UserRecord>();
const roles = ref<Array<{ label: string; value: number }>>([]);

const [Form, formApi] = useVbenForm({
  schema: useFormSchema(() => roles.value),
  showDefaultActions: false,
});

const [Drawer, drawerApi] = useVbenDrawer({
  async onConfirm() {
    const { valid } = await formApi.validate();
    if (!valid) return;

    const values = await formApi.getValues<Record<string, any>>();
    const roleIds = normalizeNumberArray(values.roleIds);

    if (!formData.value?.id && !values.password) {
      message.warning('新建用户必须填写初始密码');
      return;
    }

    drawerApi.lock();
    try {
      if (formData.value?.id) {
        await updateUser(
          formData.value.id,
          omitEmpty({
            avatar: values.avatar,
            email: values.email,
            mobile: values.mobile,
            nickname: values.nickname,
            organizationId: values.organizationId,
            realName: values.realName,
            remark: values.remark,
            status: normalizeStatus(values.status),
          }),
        );
        await replaceUserRoles(formData.value.id, roleIds);
      } else {
        await createUser(
          omitEmpty({
            avatar: values.avatar,
            email: values.email,
            mobile: values.mobile,
            nickname: values.nickname,
            organizationId: values.organizationId,
            password: values.password,
            realName: values.realName,
            remark: values.remark,
            status: normalizeStatus(values.status),
            username: values.username,
          }),
        );
        const result = await listUsers({
          keyword: values.username,
          page: 1,
          pageSize: 200,
        });
        const userId = result.list.find(
          (item) => item.username === values.username,
        )?.id;
        if (userId) {
          await replaceUserRoles(userId, roleIds);
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

    if (roles.value.length === 0) {
      const result = await listRoles({ page: 1, pageSize: 200 });
      roles.value = result.list.map((role) => ({
        label: `${role.name}（${role.code}）`,
        value: role.id,
      }));
    }

    const data = drawerApi.getData<Partial<UserRecord>>() as
      | Partial<UserRecord>
      | undefined;
    formData.value = data?.id ? (data as UserRecord) : undefined;
    formApi.resetForm();
    await nextTick();
    formApi.updateSchema([
      {
        componentProps: {
          disabled: !!formData.value?.id,
        },
        fieldName: 'username',
      },
    ]);

    formApi.setValues({
      ...data,
      organizationId: data?.organizationId || undefined,
      password: '',
      roleIds: normalizeNumberArray(data?.roleIds),
      status: normalizeStatus(data?.status),
    });
  },
});

const title = computed(() => (formData.value?.id ? '编辑用户' : '新增用户'));
</script>

<template>
  <Drawer class="w-full max-w-180" :title="title">
    <Form class="mx-4" />
  </Drawer>
</template>
