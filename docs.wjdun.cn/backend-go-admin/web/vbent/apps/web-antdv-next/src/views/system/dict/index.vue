<script lang="ts" setup>
import type { VxeTableGridOptions } from '#/adapter/vxe-table';
import type { DictDataRecord, DictTypeRecord } from '#/api/system';

import { ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';
import { Plus } from '@vben/icons';

import { Button, message } from 'antdv-next';

import { useVbenVxeGrid, VbenTableAction } from '#/adapter/vxe-table';
import {
  deleteDictData,
  deleteDictType,
  listDictData,
  listDictTypes,
} from '#/api/system';

import { normalizeOptionalStatus, toGridResult } from '../shared';
import {
  useDataColumns,
  useGridFormSchema,
  useTypeColumns,
} from './data';
import DataForm from './modules/data-form.vue';
import TypeForm from './modules/type-form.vue';

defineOptions({ name: 'SystemDict' });

const selectedType = ref<DictTypeRecord>();

const [TypeFormModal, typeFormModalApi] = useVbenModal({
  connectedComponent: TypeForm,
  destroyOnClose: true,
});

const [DataFormModal, dataFormModalApi] = useVbenModal({
  connectedComponent: DataForm,
  destroyOnClose: true,
});

const [TypeGrid, typeGridApi] = useVbenVxeGrid({
  formOptions: {
    schema: useGridFormSchema(),
    submitOnChange: true,
  },
  gridOptions: {
    columns: useTypeColumns(),
    height: 'auto',
    keepSource: true,
    proxyConfig: {
      ajax: {
        query: async ({ page }, formValues) => {
          const result = await listDictTypes({
            keyword: formValues.keyword,
            page: page.currentPage,
            pageSize: page.pageSize,
            status: normalizeOptionalStatus(formValues.status),
          });
          return toGridResult(result.list, result.total);
        },
      },
    },
    rowConfig: {
      keyField: 'id',
    },
    toolbarConfig: {
      custom: true,
      export: false,
      refresh: true,
      search: true,
      zoom: true,
    },
  } as VxeTableGridOptions<DictTypeRecord>,
});

const [DataGrid, dataGridApi] = useVbenVxeGrid({
  formOptions: {
    schema: useGridFormSchema(),
    submitOnChange: true,
  },
  gridOptions: {
    columns: useDataColumns(),
    height: 'auto',
    keepSource: true,
    proxyConfig: {
      ajax: {
        query: async ({ page }, formValues) => {
          if (!selectedType.value?.id) {
            return toGridResult([], 0);
          }
          const result = await listDictData({
            dictTypeId: selectedType.value.id,
            keyword: formValues.keyword,
            page: page.currentPage,
            pageSize: page.pageSize,
            status: normalizeOptionalStatus(formValues.status),
          });
          return toGridResult(result.list, result.total);
        },
      },
    },
    rowConfig: {
      keyField: 'id',
    },
    toolbarConfig: {
      custom: true,
      export: false,
      refresh: true,
      search: true,
      zoom: true,
    },
  } as VxeTableGridOptions<DictDataRecord>,
});

function onCreateType() {
  typeFormModalApi.setData({}).open();
}

function onEditType(row: DictTypeRecord) {
  typeFormModalApi.setData(row).open();
}

function onSelectType(row: DictTypeRecord) {
  selectedType.value = row;
  dataGridApi.query();
}

async function onDeleteType(row: DictTypeRecord) {
  await deleteDictType(row.id);
  if (selectedType.value?.id === row.id) {
    selectedType.value = undefined;
    dataGridApi.query();
  }
  message.success('删除成功');
  typeGridApi.query();
}

function onCreateData() {
  if (!selectedType.value?.id) {
    message.warning('请先选择字典类型');
    return;
  }
  dataFormModalApi.setData({ dictTypeId: selectedType.value.id }).open();
}

function onEditData(row: DictDataRecord) {
  dataFormModalApi.setData(row).open();
}

async function onDeleteData(row: DictDataRecord) {
  await deleteDictData(row.id);
  message.success('删除成功');
  dataGridApi.query();
}
</script>

<template>
  <Page auto-content-height>
    <TypeFormModal @success="typeGridApi.query()" />
    <DataFormModal @success="dataGridApi.query()" />
    <div class="flex size-full gap-4">
      <div class="w-2/5">
        <TypeGrid table-title="字典类型">
          <template #toolbar-tools>
            <Button type="primary" @click="onCreateType">
              <Plus class="size-5" />
              新增类型
            </Button>
          </template>
          <template #action="{ row }">
            <VbenTableAction
              :actions="[
                {
                  text: '选择',
                  icon: 'lucide:check',
                  onClick: () => onSelectType(row),
                },
                {
                  text: '编辑',
                  icon: 'lucide:edit',
                  onClick: () => onEditType(row),
                },
              ]"
              :dropdown-actions="[
                {
                  text: '删除',
                  icon: 'lucide:trash-2',
                  danger: true,
                  popConfirm: {
                    title: `确认删除字典类型 ${row.name}？`,
                    confirm: () => onDeleteType(row),
                  },
                },
              ]"
              align="center"
            />
          </template>
        </TypeGrid>
      </div>
      <div class="w-3/5">
        <DataGrid :table-title="selectedType ? `字典数据：${selectedType.name}` : '字典数据'">
          <template #toolbar-tools>
            <Button type="primary" @click="onCreateData">
              <Plus class="size-5" />
              新增数据
            </Button>
          </template>
          <template #action="{ row }">
            <VbenTableAction
              :actions="[
                {
                  text: '编辑',
                  icon: 'lucide:edit',
                  onClick: () => onEditData(row),
                },
              ]"
              :dropdown-actions="[
                {
                  text: '删除',
                  icon: 'lucide:trash-2',
                  danger: true,
                  popConfirm: {
                    title: `确认删除字典数据 ${row.label}？`,
                    confirm: () => onDeleteData(row),
                  },
                },
              ]"
              align="center"
            />
          </template>
        </DataGrid>
      </div>
    </div>
  </Page>
</template>
