<script lang="ts" setup>
import type { IDomEditor, IEditorConfig } from '@wangeditor/editor';
import type { ArticleRecord, SaveArticlePayload } from '#/api/system';

import { onBeforeUnmount, onMounted, reactive, ref, shallowRef } from 'vue';

import { Page } from '@vben/common-ui';
import {
  Button,
  Card,
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
import DOMPurify from 'dompurify';
import '@wangeditor/editor/dist/css/style.css';

import {
  createArticle,
  deleteArticle,
  getArticle,
  listArticles,
  uploadArticleImage,
  updateArticle,
} from '#/api/system';

import { normalizeOptionalStatus } from '../shared';
import ArticleContent from './content.vue';

defineOptions({ name: 'SystemArticle' });

const loading = ref(false);
const saving = ref(false);
const detailLoading = ref(false);
const viewing = ref(false);
const loadError = ref(false);
let requestVersion = 0;
const editor = shallowRef<IDomEditor>();
const editorVisible = ref(false);
const editingId = ref<number>();
const keyword = ref('');
const status = ref<number | undefined>();
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const rows = ref<ArticleRecord[]>([]);
const formRef = ref();
const form = reactive<SaveArticlePayload>({
  sort: 0,
  category: '',
  author: '',
  content: '',
  cover: '',
  status: 1,
  summary: '',
  title: '',
});

async function replaceEmbeddedImages(html: string): Promise<string> {
  if (!html.includes('data:image/')) return html;
  const documentHtml = new DOMParser().parseFromString(html, 'text/html');
  const images = documentHtml.querySelectorAll('img[src^="data:image/"]');
  for (const image of images) {
    const source = image.getAttribute('src');
    if (!source) continue;
    const blob = await fetch(source).then((response) => response.blob());
    const extension = blob.type.split('/')[1]?.replaceAll(/[^a-z0-9]/gi, '') || 'png';
    const file = new File([blob], `article-image-${Date.now()}.${extension}`, {
      type: blob.type,
    });
    image.setAttribute('src', await uploadArticleImage(file));
  }
  return documentHtml.body.innerHTML;
}

const editorConfig: Partial<IEditorConfig> = {
  placeholder: '请输入文章正文…',
  MENU_CONF: {
    uploadImage: {
      base64LimitSize: 0,
      allowedFileTypes: ['image/png', 'image/jpeg', 'image/gif', 'image/webp'],
      async customUpload(file, insertFn) {
        try {
          insertFn(await uploadArticleImage(file), file.name, '');
        } catch (error) {
          message.error(error instanceof Error ? error.message : '图片上传失败');
        }
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
  const version = ++requestVersion;
  loading.value = true;
  loadError.value = false;
  try {
    const result = await listArticles({
      keyword: keyword.value || undefined,
      page: page.value,
      pageSize: pageSize.value,
      status: normalizeOptionalStatus(status.value),
    });
    if (version !== requestVersion) return;
    rows.value = result.list;
    total.value = result.total;
  } catch {
    if (version === requestVersion) loadError.value = true;
  } finally {
    if (version === requestVersion) loading.value = false;
  }
}

function resetForm() {
  Object.assign(form, {
    sort: 0,
    category: '',
    author: '',
    content: '',
    cover: '',
    status: 1,
    summary: '',
    title: '',
  });
}

function openCreate() {
  viewing.value = false;
  editingId.value = undefined;
  resetForm();
  editor.value?.enable();
  editorVisible.value = true;
}

async function openEdit(row: ArticleRecord, readOnly = false) {
  if (detailLoading.value) return;
  detailLoading.value = true;
  try {
    const detail = await getArticle(row.id);
    editingId.value = row.id;
    viewing.value = readOnly;
    if (readOnly) editor.value?.disable();
    else editor.value?.enable();
    Object.assign(form, {
      sort: detail.sort ?? 0,
      category: detail.category ?? '',
      title: detail.title,
      author: detail.author,
      summary: detail.summary,
      content: DOMPurify.sanitize(detail.content, { USE_PROFILES: { html: true } }),
      cover: detail.cover,
      status: detail.status,
    });
    editorVisible.value = true;
  } catch {
    // The request interceptor displays the API error.
  } finally {
    detailLoading.value = false;
  }
}

function closeEditor() {
  if (saving.value) return;
  editorVisible.value = false;
}

function onEditorCreated(instance: IDomEditor) {
  editor.value = instance;
  if (viewing.value) instance.disable();
}

function onStatusChange(checked: boolean) {
  form.status = checked ? 1 : 0;
}

async function save() {
  if (saving.value || viewing.value) return;
  form.title = form.title.trim();
  form.category = form.category?.trim() ?? '';
  try {
    await formRef.value?.validate();
  } catch {
    return;
  }
  if (editor.value?.isEmpty()) {
    message.warning('请输入文章正文');
    return;
  }
  saving.value = true;
  try {
    form.content = await replaceEmbeddedImages(form.content);
  } catch (error) {
    message.error(error instanceof Error ? error.message : '正文图片上传失败');
    saving.value = false;
    return;
  }
  try {
    if (editingId.value) {
      await updateArticle(editingId.value, form);
      message.success('文章已更新');
    } else {
      await createArticle(form);
      message.success('文章已创建');
    }
    saving.value = false;
    closeEditor();
    await loadArticles();
  } catch {
    // Keep the user's draft when saving fails.
  } finally {
    saving.value = false;
  }
}

async function remove(id: number) {
  try {
    await deleteArticle(id);
    message.success('文章已删除');
    if (rows.value.length === 1 && page.value > 1) page.value -= 1;
    await loadArticles();
  } catch {
    // The request interceptor displays the API error.
  }
}

function search() {
  page.value = 1;
  loadArticles();
}

function changePage(nextPage: number, nextSize: number) {
  page.value = nextPage;
  pageSize.value = nextSize;
  loadArticles();
}

onMounted(loadArticles);
onBeforeUnmount(() => editor.value?.destroy());
</script>

<template>
  <Page auto-content-height>
    <Card :bordered="false" class="article-page">
      <template #title>
        <div class="page-heading">
          <div>
            <div class="eyebrow">CONTENT STUDIO</div>
            <h1>文章管理</h1>
            <p>文章按类别分组，cdn新闻序号越大越靠前；cdn帮助及其他类别序号越小越靠前。</p>
          </div>
          <Button
            v-access:code="['content:article:create']"
            type="primary"
            @click="openCreate"
            >新建文章</Button
          >
        </div>
      </template>
      <div class="filter-bar">
        <Input
          v-model:value="keyword"
          class="keyword-input"
          allow-clear
          placeholder="搜索文章标题"
          @press-enter="search"
        />
        <Select
          v-model:value="status"
          class="status-select"
          allow-clear
          placeholder="全部状态"
          :options="[
            { label: '已发布', value: 1 },
            { label: '草稿', value: 0 },
          ]"
          @change="search"
        />
        <Button @click="search">搜索</Button>
        <Button
          @click="
            keyword = '';
            status = undefined;
            search();
          "
          >重置</Button
        >
      </div>
      <Spin :spinning="loading">
        <Table
          :columns="columns"
          :data-source="rows"
          :pagination="false"
          row-key="id"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'status'"
              ><Tag :color="record.status === 1 ? 'success' : 'default'">{{
                record.status === 1 ? '已发布' : '草稿'
              }}</Tag></template
            >
            <template v-else-if="column.key === 'actions'">
              <Space>
                <Button
                  type="link"
                  size="small"
                  :disabled="detailLoading"
                  @click="openEdit(record, true)"
                  >查看</Button
                >
                <Button
                  v-access:code="['content:article:update']"
                  type="link"
                  size="small"
                  :disabled="detailLoading"
                  @click="openEdit(record)"
                  >编辑</Button
                >
                <Popconfirm
                  title="确定删除这篇文章吗？"
                  @confirm="remove(record.id)"
                  ><Button
                    v-access:code="['content:article:delete']"
                    type="link"
                    danger
                    size="small"
                    >删除</Button
                  ></Popconfirm
                >
              </Space>
            </template>
          </template>
        </Table>
        <div v-if="loadError" role="alert">
          文章加载失败，请重试。<Button @click="loadArticles">重试</Button>
        </div>
      </Spin>
      <div class="pagination-row">
        <span>共 {{ total }} 篇文章</span
        ><Pagination
          :current="page"
          :page-size="pageSize"
          :total="total"
          show-size-changer
          @change="changePage"
        />
      </div>
    </Card>

    <Modal
      :open="editorVisible"
      :confirm-loading="saving"
      :title="viewing ? '查看文章' : editingId ? '编辑文章' : '新建文章'"
      :footer="viewing ? null : undefined"
      :mask-closable="!saving"
      :closable="!saving"
      width="960px"
      ok-text="保存文章"
      cancel-text="取消"
      @cancel="closeEditor"
      @ok="save"
    >
      <Form
        ref="formRef"
        layout="vertical"
        :model="form"
        :disabled="viewing || saving"
      >
        <div class="form-grid">
          <FormItem
            label="标题"
            name="title"
            :rules="[{ required: true, message: '请输入文章标题' }]"
            ><Input
              v-model:value="form.title"
              :maxlength="120"
              show-count
              placeholder="输入文章标题"
          /></FormItem>
          <FormItem label="作者" name="author"
            ><Input
              v-model:value="form.author"
              :maxlength="40"
              placeholder="输入作者名称"
          /></FormItem>
        </div>
        <FormItem label="分类" name="category">
          <Input v-model:value="form.category" :maxlength="80" show-count placeholder="自行填写分类，例如：技术分享" />
        </FormItem>
        <FormItem
          label="排序"
          name="sort"
          extra="填写几就保存几，cdn新闻序号越大越靠前，cdn帮助及其他类别序号越小越靠前；同类别内排序号不可重复，不同类别可以使用相同排序号。"
          :rules="[{ required: true, type: 'integer', min: 0, max: 2147483647, message: '请输入 0 至 2147483647 之间的整数' }]"
        >
          <InputNumber v-model:value="form.sort" :min="0" :max="2147483647" :precision="0" />
        </FormItem>
        <FormItem label="摘要" name="summary"
          ><Input.TextArea
            v-model:value="form.summary"
            :rows="2"
            :maxlength="240"
            show-count
            placeholder="用于列表和分享卡片的简短摘要"
        /></FormItem>
        <FormItem label="封面地址" name="cover"
          ><Input v-model:value="form.cover" placeholder="可选，填写图片地址"
        /></FormItem>
        <FormItem label="发布状态" name="status"
          ><Switch
            :checked="form.status === 1"
            checked-children="发布"
            un-checked-children="草稿"
            @change="onStatusChange"
        /></FormItem>
        <FormItem
          label="正文"
          name="content"
          :rules="[{ required: true, message: '请输入文章正文' }]"
        >
          <ArticleContent v-if="viewing" :content="form.content" />
          <div v-show="!viewing" class="editor-shell">
            <Toolbar
              v-if="editor"
              :editor="editor"
              :default-config="{ excludeKeys: ['uploadVideo'] }"
              mode="default"
            /><Editor
              v-model="form.content"
              :default-config="editorConfig"
              style="height: 320px"
              mode="default"
              @on-created="onEditorCreated"
            />
          </div>
        </FormItem>
      </Form>
    </Modal>
  </Page>
</template>

<style scoped>
.article-page {
  min-height: 100%;
}
.page-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
}
.eyebrow {
  color: #0f766e;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.14em;
}
h1 {
  margin: 5px 0 0;
  color: inherit;
  font-size: 24px;
  line-height: 1.2;
}
.page-heading p {
  margin: 7px 0 0;
  color: #718096;
  font-size: 13px;
}
.filter-bar {
  display: flex;
  gap: 10px;
  margin: 4px 0 20px;
}
.keyword-input {
  width: 300px;
}
.status-select {
  width: 140px;
}
.pagination-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-top: 20px;
  color: #718096;
  font-size: 13px;
}
.form-grid {
  display: grid;
  grid-template-columns: 1fr 220px;
  gap: 16px;
}
.editor-shell :deep(.w-e-text-container) {
  min-height: 300px;
}
@media (max-width: 640px) {
  .page-heading,
  .filter-bar,
  .pagination-row {
    align-items: stretch;
    flex-direction: column;
  }
  .keyword-input,
  .status-select {
    width: 100%;
  }
  .form-grid {
    grid-template-columns: 1fr;
    gap: 0;
  }
}
</style>
