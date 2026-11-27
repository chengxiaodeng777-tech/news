import type { VbenFormSchema } from '#/adapter/form';
import type { VxeTableGridColumns } from '#/adapter/vxe-table';
import type { OrganizationRecord } from '#/api/system';

import { statusOptions, statusTagOptions } from '../shared';

export const organizationTypeOptions = [
  { label: '平台', value: 'platform' },
  { label: '公司', value: 'company' },
  { label: '部门', value: 'department' },
  { label: '小组', value: 'group' },
];

export function useFormSchema(): VbenFormSchema[] {
  return [
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
      fieldName: 'parentId',
      label: '上级组织',
    },
    {
      component: 'Input',
      fieldName: 'name',
      label: '组织名称',
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'code',
      label: '组织编码',
      rules: 'required',
    },
    {
      component: 'Select',
      componentProps: {
        options: organizationTypeOptions,
      },
      fieldName: 'type',
      label: '组织类型',
    },
    {
      component: 'Input',
      fieldName: 'leader',
      label: '负责人',
    },
    {
      component: 'Input',
      fieldName: 'phone',
      label: '电话',
    },
    {
      component: 'Input',
      fieldName: 'email',
      label: '邮箱',
    },
    {
      component: 'InputNumber',
      fieldName: 'sort',
      label: '排序',
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

export function useColumns(): VxeTableGridColumns<OrganizationRecord> {
  return [
    {
      align: 'left',
      field: 'name',
      fixed: 'left',
      title: '组织名称',
      treeNode: true,
      width: 220,
    },
    { field: 'code', minWidth: 140, title: '组织编码' },
    {
      cellRender: { name: 'CellTag', options: organizationTypeOptions },
      field: 'type',
      title: '类型',
      width: 100,
    },
    { field: 'leader', minWidth: 120, title: '负责人' },
    { field: 'phone', minWidth: 120, title: '电话' },
    {
      cellRender: { name: 'CellTag', options: statusTagOptions },
      field: 'status',
      title: '状态',
      width: 90,
    },
    { field: 'sort', title: '排序', width: 90 },
    {
      align: 'center',
      field: 'operation',
      fixed: 'right',
      slots: { default: 'action' },
      title: '操作',
      width: 210,
    },
  ];
}
