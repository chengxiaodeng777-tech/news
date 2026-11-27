import type { VbenFormSchema } from '#/adapter/form';
import type { VxeTableGridColumns } from '#/adapter/vxe-table';
import type { MenuRecord } from '#/api/system';

import { statusOptions, statusTagOptions } from '../shared';

export const menuTypeOptions = [
  { color: 'processing', label: '目录', value: 'catalog' },
  { color: 'success', label: '菜单', value: 'menu' },
  { color: 'warning', label: '按钮', value: 'button' },
  { color: 'error', label: '接口', value: 'api' },
];

export function useFormSchema(): VbenFormSchema[] {
  return [
    {
      component: 'ApiTreeSelect',
      componentProps: {
        allowClear: true,
        api: async () => {
          const { listMenuTree } = await import('#/api/system');
          return await listMenuTree();
        },
        childrenField: 'children',
        class: 'w-full',
        labelField: 'title',
        treeDefaultExpandAll: true,
        valueField: 'id',
      },
      fieldName: 'parentId',
      label: '上级菜单',
    },
    {
      component: 'RadioGroup',
      componentProps: {
        buttonStyle: 'solid',
        optionType: 'button',
        options: menuTypeOptions,
      },
      defaultValue: 'menu',
      fieldName: 'type',
      label: '菜单类型',
    },
    {
      component: 'Input',
      fieldName: 'name',
      label: '路由名称',
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'title',
      label: '显示标题',
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'path',
      label: '路由路径',
    },
    {
      component: 'Input',
      fieldName: 'component',
      label: '组件路径',
    },
    {
      component: 'Input',
      fieldName: 'redirect',
      label: '重定向',
    },
    {
      component: 'IconPicker',
      componentProps: {
        prefix: 'carbon',
      },
      fieldName: 'icon',
      label: '图标',
    },
    {
      component: 'Input',
      fieldName: 'permissionCode',
      label: '权限码',
    },
    {
      component: 'Input',
      fieldName: 'externalLink',
      label: '外链',
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
      component: 'Checkbox',
      fieldName: 'hidden',
      label: '隐藏菜单',
      renderComponentContent: () => ({ default: () => '隐藏菜单' }),
    },
    {
      component: 'Checkbox',
      fieldName: 'keepAlive',
      label: '缓存页面',
      renderComponentContent: () => ({ default: () => '缓存页面' }),
    },
    {
      component: 'Textarea',
      fieldName: 'remark',
      label: '备注',
    },
  ];
}

export function useColumns(): VxeTableGridColumns<MenuRecord> {
  return [
    {
      align: 'left',
      field: 'title',
      fixed: 'left',
      title: '菜单名称',
      treeNode: true,
      width: 220,
    },
    {
      cellRender: { name: 'CellTag', options: menuTypeOptions },
      field: 'type',
      title: '类型',
      width: 90,
    },
    { field: 'name', minWidth: 140, title: '路由名称' },
    { field: 'path', minWidth: 160, title: '路由路径' },
    { field: 'component', minWidth: 180, title: '组件' },
    { field: 'permissionCode', minWidth: 160, title: '权限码' },
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
