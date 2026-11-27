<script lang="ts" setup>
import type { IDomEditor, IEditorConfig } from '@wangeditor/editor';

import { onBeforeUnmount, reactive, ref, shallowRef } from 'vue';

import { Page } from '@vben/common-ui';
import {
  Button,
  Card,
  Empty,
  Form,
  FormItem,
  Input,
  InputNumber,
  Modal,
  Pagination,
  Popconfirm,
  Select,
  Space,
  Spin,
  Switch,
  Table,
  Tag,
  message,
} from 'antdv-next';
import { Editor, Toolbar } from '@wangeditor/editor-for-vue';
import '@wangeditor/editor/dist/css/style.css';

import {
  createArticle,
  deleteArticle,
  getArticleList,
  updateArticle,
  type ArticleApi,
} from '#/api';

const loading = ref(false);
const saving = ref(false);
const editor = shallowRef<IDomEditor>();
const editorVisible = ref(false);
const editingId = ref<number>();
const keyword = ref('');
const status = ref<number | undefined>();
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const rows = ref<ArticleApi.Article[]>([]);
const formRef = ref();
const form = reactive({
  sort: 0,
  category: '',
  title: '',
  summary: '',
  content: '',
  cover: '',
  status: 1 as 0 | 1,
  author: '',
});

const toolbarConfig = {};
const editorConfig: Partial<IEditorConfig> = {
  placeholder: '在这里写下文章正文...',
  MENU_CONF: {
    uploadImage: {
      base64LimitSize: 5 * 1024 * 1024,
      allowedFileTypes: ['image/*'],
      customUpload() {
        message.warning('图片大小不能超过 5 MB');
      },
    },
  },
};

const columns = [
  { title: '排序', dataIndex: 'sort', key: 'sort', width: 90 },
  { title: '标题', dataIndex: 'title', key: 'title', ellipsis: true },
  { title: '分类', dataIndex: 'category', key: 'category', width: 140, ellipsis: true },
  { title: '作者', dataIndex: 'author', key: 'author', width: 140 },
  { title: '状态', key: 'status', width: 100 },
  { title: '更新时间', dataIndex: 'updatedAt', key: 'updatedAt', width: 190 },
  { title: '操作', key: 'actions', width: 180 },
];

async function loadArticles() {
  loading.value = true;
  try {
    const result = await getArticleList({
      keyword: keyword.value || undefined,
      status: status.value,
      page: page.value,
      pageSize: pageSize.value,
    });
    rows.value = result.list ?? [];
    total.value = result.total ?? 0;
  } finally {
    loading.value = false;
  }
}

function resetForm() {
  Object.assign(form, { sort: 0, title: '', category: '', summary: '', content: '', cover: '', status: 1, author: '' });
}

function openCreate() {
  editingId.value = undefined;
  resetForm();
  editorVisible.value = true;
}

function openEdit(row: ArticleApi.Article) {
  editingId.value = row.id;
  Object.assign(form, row);
  editorVisible.value = true;
}

function closeEditor() {
  editorVisible.value = false;
}

function handleEditorCreated(instance: IDomEditor) {
  editor.value = instance;
}

async function saveArticle() {
  await formRef.value?.validate();
  saving.value = true;
  try {
    const payload = { ...form };
    if (editingId.value) {
      await updateArticle(editingId.value, payload);
      message.success('文章已更新');
    } else {
      await createArticle(payload);
      message.success('文章已创建');
    }
    closeEditor();
    await loadArticles();
  } finally {
    saving.value = false;
  }
}

async function removeArticle(id: number) {
  await deleteArticle(id);
  message.success('文章已删除');
  if (rows.value.length === 1 && page.value > 1) page.value -= 1;
  await loadArticles();
}

function onSearch() {
  page.value = 1;
  loadArticles();
}

function onPageChange(nextPage: number, nextPageSize: number) {
  page.value = nextPage;
  pageSize.value = nextPageSize;
  loadArticles();
}

function onStatusChange(checked: boolean) {
  form.status = checked ? 1 : 0;
}

onBeforeUnmount(() => editor.value?.destroy());
loadArticles();
</script>

<template>
  <Page auto-content-height>
    <Card class="article-page" :bordered="false">
      <template #title>
        <div class="page-heading">
          <div>
            <div class="eyebrow">CONTENT STUDIO</div>
            <h1>文章管理</h1>
            <p>维护文章内容、发布状态和排序，排序值越小越靠前。</p>
          </div>
          <Button type="primary" @click="openCreate">新建文章</Button>
        </div>
      </template>

      <div class="filter-bar">
        <Input v-model:value="keyword" allow-clear class="keyword-input" placeholder="搜索文章标题" @press-enter="onSearch" />
        <Select v-model:value="status" allow-clear class="status-select" placeholder="全部状态" :options="[{ label: '已发布', value: 1 }, { label: '草稿', value: 0 }]" @change="onSearch" />
        <Button @click="onSearch">搜索</Button>
      </div>

      <Spin :spinning="loading">
        <Table :columns="columns" :data-source="rows" :pagination="false" row-key="id">
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'status'">
              <Tag :color="record.status === 1 ? 'success' : 'default'">{{ record.status === 1 ? '已发布' : '草稿' }}</Tag>
            </template>
            <template v-else-if="column.key === 'actions'">
              <Space>
                <Button type="link" size="small" @click="openEdit(record)">编辑</Button>
                <Popconfirm title="确定删除这篇文章吗？" @confirm="removeArticle(record.id)">
                  <Button type="link" danger size="small">删除</Button>
                </Popconfirm>
              </Space>
            </template>
          </template>
        </Table>
        <Empty v-if="!loading && rows.length === 0" description="暂无文章" />
      </Spin>

      <div class="pagination-row">
        <span>共 {{ total }} 篇文章</span>
        <Pagination :current="page" :page-size="pageSize" :total="total" show-size-changer @change="onPageChange" />
      </div>
    </Card>

    <Modal :open="editorVisible" :confirm-loading="saving" :title="editingId ? '编辑文章' : '新建文章'" width="960px" ok-text="保存文章" cancel-text="取消" @cancel="closeEditor" @ok="saveArticle">
      <Form ref="formRef" layout="vertical" :model="form">
        <div class="form-grid">
          <FormItem label="标题" name="title" :rules="[{ required: true, message: '请输入文章标题' }]">
            <Input v-model:value="form.title" :maxlength="120" show-count placeholder="输入文章标题" />
          </FormItem>
          <FormItem label="作者" name="author"><Input v-model:value="form.author" :maxlength="40" placeholder="输入作者名称" /></FormItem>
        </div>
        <FormItem label="摘要" name="summary"><Input.TextArea v-model:value="form.summary" :rows="2" :maxlength="240" show-count placeholder="用于列表和分享卡片的简短摘要" /></FormItem>
        <FormItem label="分类" name="category"><Input v-model:value="form.category" :maxlength="80" show-count placeholder="自行填写分类，例如：技术分享" /></FormItem>
        <FormItem label="封面地址" name="cover"><Input v-model:value="form.cover" placeholder="可选，填写图片地址" /></FormItem>
        <FormItem label="发布状态" name="status"><Switch :checked="form.status === 1" checked-children="发布" un-checked-children="草稿" @change="onStatusChange" /></FormItem>
        <FormItem label="排序" name="sort" extra="数值越小越靠前；相同数值按文章 ID 倒序，修改内容不会改变顺序。" :rules="[{ required: true, type: 'integer', min: 0, max: 2147483647, message: '请输入 0 至 2147483647 之间的整数' }]">
          <InputNumber v-model:value="form.sort" :min="0" :max="2147483647" :precision="0" />
        </FormItem>
        <FormItem label="正文" name="content" :rules="[{ required: true, message: '请输入文章正文' }]">
          <div class="editor-shell">
            <Toolbar v-if="editor" :editor="editor" :default-config="toolbarConfig" mode="default" />
            <Editor v-model="form.content" :default-config="editorConfig" mode="default" @on-created="handleEditorCreated" />
          </div>
        </FormItem>
      </Form>
    </Modal>
  </Page>
</template>

<style scoped>
.article-page { min-height: 100%; }
.page-heading { display: flex; align-items: center; justify-content: space-between; gap: 24px; }
.eyebrow { color: #0f766e; font-size: 11px; font-weight: 700; letter-spacing: 0.14em; }
h1 { margin: 5px 0 0; color: #17212b; font-size: 24px; line-height: 1.2; }
.page-heading p { margin: 7px 0 0; color: #718096; font-size: 13px; }
.filter-bar { display: flex; gap: 10px; margin: 4px 0 20px; }
.keyword-input { width: 300px; }
.status-select { width: 140px; }
.pagination-row { display: flex; justify-content: space-between; align-items: center; gap: 16px; margin-top: 20px; color: #718096; font-size: 13px; }
.form-grid { display: grid; grid-template-columns: 1fr 220px; gap: 16px; }
.editor-shell :deep(.w-e-text-container) { min-height: 300px; }
@media (max-width: 640px) {
  .page-heading, .filter-bar, .pagination-row { align-items: stretch; flex-direction: column; }
  .keyword-input, .status-select { width: 100%; }
  .form-grid { grid-template-columns: 1fr; gap: 0; }
}
</style>
