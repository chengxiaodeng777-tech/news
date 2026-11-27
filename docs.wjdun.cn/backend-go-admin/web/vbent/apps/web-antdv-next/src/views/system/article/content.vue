<script setup lang="ts">
import { computed } from 'vue';

import DOMPurify from 'dompurify';

const props = defineProps<{ content: string }>();
const html = computed(() => DOMPurify.sanitize(props.content, { USE_PROFILES: { html: true } }));
</script>

<template>
  <!-- Only sanitized article HTML is rendered. -->
  <!-- eslint-disable-next-line vue/no-v-html -->
  <div class="article-content" v-html="html"></div>
</template>

<style scoped>
.article-content {
  overflow-wrap: anywhere;
  color: hsl(var(--foreground));
  line-height: 1.8;
}
.article-content :deep(p) { margin: 0.75em 0; }
.article-content :deep(h1) { font-size: 2em; }
.article-content :deep(h2) { font-size: 1.5em; }
.article-content :deep(h3) { font-size: 1.25em; }
.article-content :deep(h4) { font-size: 1.1em; }
.article-content :deep(h1),
.article-content :deep(h2),
.article-content :deep(h3),
.article-content :deep(h4),
.article-content :deep(h5) { margin: 1em 0 0.5em; font-weight: 600; }
.article-content :deep(img),
.article-content :deep(video) { max-width: 100%; height: auto; }
.article-content :deep(a) { color: hsl(var(--primary)); text-decoration: underline; }
.article-content :deep(ul) { padding-left: 1.5em; list-style: disc; }
.article-content :deep(ol) { padding-left: 1.5em; list-style: decimal; }
.article-content :deep(blockquote) {
  margin: 1em 0;
  padding: 0.5em 1em;
  background: hsl(var(--muted));
  border-left: 4px solid hsl(var(--primary));
}
.article-content :deep(code) {
  padding: 0.15em 0.35em;
  font-family: monospace;
  background: hsl(var(--muted));
  border-radius: 4px;
}
.article-content :deep(pre) {
  overflow-x: auto;
  padding: 1em;
  background: hsl(var(--muted));
  border-radius: var(--radius);
  white-space: pre;
}
.article-content :deep(pre code) { padding: 0; }
.article-content :deep(table) { display: block; max-width: 100%; overflow-x: auto; border-collapse: collapse; }
.article-content :deep(td),
.article-content :deep(th) { padding: 0.4em 0.75em; border: 1px solid hsl(var(--border)); }
.article-content :deep(th) { background: hsl(var(--muted)); }
.article-content :deep(hr) { margin: 1.5em 0; border-color: hsl(var(--border)); }
</style>
