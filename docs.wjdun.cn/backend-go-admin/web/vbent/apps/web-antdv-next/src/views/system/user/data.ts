import type { VbenFormSchema } from '#/adapter/form';
import type { VxeTableGridColumns } from '#/adapter/vxe-table';
import type { UserRecord } from '#/api/system';

import { statusOptions, statusTagOptions } from '../shared';

export function useGridFormSchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'keyword',
      label: '关键字',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: statusOptions,
      },
      fieldName: 'status',
      label: '状态',
    },
  ];
}

export function useFormSchema(roleOptions: () => Array<{ label: string; value: number }>): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'username',
      label: '用户名',
      rules: 'required',
    },
    {
      component: 'InputPassword',
      fieldName: 'password',
      label: '初始密码',
    },
    {
      component: 'Input',
      fieldName: 'realName',
      label: '姓名',
    },
    {
      component: 'Input',
      fieldName: 'nickname',
      label: '昵称',
    },
    {
      component: 'ApiTreeSelect',
      componentProps: {
        allowClear: true,
        api: async () => {
          const { listOrganizationTree } = await import('#/api/system');
          return await listOrganizationTree();
        },
        childrenField: 'children',
        class: 'w-full',
        labelField: 'name',
        treeDefaultExpandAll: true,
        valueField: 'id',
      },
      fieldName: 'organizationId',
      label: '组织',
    },
    {
      component: 'Select',
      componentProps: () => ({
        allowClear: true,
        mode: 'multiple',
        optionFilterProp: 'label',
        options: roleOptions(),
      }),
      fieldName: 'roleIds',
      label: '角色',
    },
    {
      component: 'Input',
      fieldName: 'mobile',
      label: '手机号',
    },
    {
      component: 'Input',
      fieldName: 'email',
      label: '邮箱',
    },
    {
      component: 'Input',
      fieldName: 'avatar',
      label: '头像',
    },
    {
      component: 'RadioGroup',
      componentProps: {
        buttonStyle: 'solid',
        optionType: 'button',
        options: statusOptions,
      },
      defaultValue: 1,
      fieldName: 'status',
      label: '状态',
    },
    {
      component: 'Textarea',
      fieldName: 'remark',
      label: '备注',
    },
  ];
}

export function useColumns(
  getRoleNames: (roleIds?: number[]) => string,
): VxeTableGridColumns<UserRecord> {
  return [
    { field: 'id', title: 'ID', width: 80 },
    { field: 'username', title: '用户名', minWidth: 120 },
    { field: 'realName', title: '姓名', minWidth: 120 },
    { field: 'organizationName', title: '组织', minWidth: 140 },
    {
      field: 'roleIds',
      formatter: ({ row }) => getRoleNames(row.roleIds),
      minWidth: 180,
      title: '角色',
    },
    {
      cellRender: { name: 'CellTag', options: statusTagOptions },
      field: 'status',
      title: '状态',
      width: 90,
    },
    { field: 'remark', minWidth: 140, title: '备注' },
    {
      align: 'center',
      field: 'operation',
      fixed: 'right',
      slots: { default: 'action' },
      title: '操作',
      width: 160,
    },
  ];
}
