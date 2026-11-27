import type { VbenFormSchema } from '#/adapter/form';
import type { VxeTableGridColumns } from '#/adapter/vxe-table';
import type { DictDataRecord, DictTypeRecord } from '#/api/system';

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

export function useTypeFormSchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'code',
      label: '字典编码',
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'name',
      label: '字典名称',
      rules: 'required',
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

export function useDataFormSchema(): VbenFormSchema[] {
  return [
    {
      component: 'Input',
      fieldName: 'label',
      label: '标签',
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'value',
      label: '值',
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'tagType',
      label: '标签类型',
    },
    {
      component: 'Input',
      fieldName: 'cssClass',
      label: 'CSS 类名',
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
      component: 'Switch',
      fieldName: 'isDefault',
      label: '默认',
    },
    {
      component: 'Textarea',
      fieldName: 'remark',
      label: '备注',
    },
  ];
}

export function useTypeColumns(): VxeTableGridColumns<DictTypeRecord> {
  return [
    { field: 'code', minWidth: 150, title: '字典编码' },
    { field: 'name', minWidth: 150, title: '字典名称' },
    {
      cellRender: { name: 'CellTag', options: statusTagOptions },
      field: 'status',
      title: '状态',
      width: 90,
    },
    {
      align: 'center',
      field: 'operation',
      fixed: 'right',
      slots: { default: 'action' },
      title: '操作',
      width: 190,
    },
  ];
}

export function useDataColumns(): VxeTableGridColumns<DictDataRecord> {
  return [
    { field: 'label', minWidth: 120, title: '标签' },
    { field: 'value', minWidth: 120, title: '值' },
    { field: 'tagType', minWidth: 120, title: '标签类型' },
    { field: 'sort', title: '排序', width: 90 },
    {
      cellRender: { name: 'CellTag', options: statusTagOptions },
      field: 'status',
      title: '状态',
      width: 90,
    },
    {
      cellRender: {
        name: 'CellTag',
        options: [
          { color: 'success', label: '是', value: true },
          { color: 'default', label: '否', value: false },
        ],
      },
      field: 'isDefault',
      title: '默认',
      width: 90,
    },
    {
      align: 'center',
      field: 'operation',
      fixed: 'right',
      slots: { default: 'action' },
      title: '操作',
      width: 150,
    },
  ];
}
