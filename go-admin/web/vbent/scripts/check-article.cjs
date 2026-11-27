// Run against the local frontend and API: node scripts/check-article.cjs
const assert = require('node:assert/strict');
const { chromium } = require('playwright');
const baseURL = process.env.ARTICLE_BASE_URL || 'http://localhost:5999';

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage({ permissions: ['clipboard-read', 'clipboard-write'] });
  if (process.env.ARTICLE_API_URL) {
    await page.route('**/api/**', async (route) => {
      const url = new URL(route.request().url());
      if (!url.pathname.startsWith('/api/')) {
        await route.continue();
        return;
      }
      const response = await route.fetch({ url: `${process.env.ARTICLE_API_URL}${url.pathname}${url.search}` });
      await route.fulfill({ response });
    });
  }
  const errors = [];
  const requests = [];
  let articleId;
  page.on('pageerror', (error) => errors.push(error.message));
  page.on('request', (request) => {
    if (new URL(request.url()).pathname.startsWith('/api/')) {
      requests.push(`${request.method()} ${new URL(request.url()).pathname}`);
    }
  });
  const title = `文章回归测试 ${Date.now()}`;
  try {
    await page.goto(`${baseURL}/system/article`);
    await page.getByPlaceholder('请输入用户名').waitFor();
    await page.evaluate(async () => {
      const { useAuthStore } = await import('/src/store/auth.ts');
      await useAuthStore().authLogin(
        { username: 'admin', password: 'admin123456' },
        () => {},
      );
    });
    // Simulate permissions cached before article management was added.
    await page.evaluate(() => {
      for (const key of Object.keys(localStorage)) {
        if (!key.endsWith('-access')) continue;
        const state = JSON.parse(localStorage.getItem(key));
        state.accessCodes = [];
        localStorage.setItem(key, JSON.stringify(state));
      }
    });
    await page.goto(`${baseURL}/system/article`);
    await page.getByRole('button', { name: '新建文章' }).waitFor();
    await page.waitForTimeout(2000);
    assert.equal(
      requests.filter((r) => r === 'GET /api/system/articles').length,
      1,
    );
    await page.getByRole('button', { name: '新建文章' }).click();
    const dialog = page.getByRole('dialog');
    await dialog.getByPlaceholder('输入文章标题').fill(title);
    await dialog.getByPlaceholder('自行填写分类，例如：技术分享').fill('技术分享');
    await dialog.getByRole('button', { name: '保存文章' }).click();
    await page.getByText('请输入文章正文', { exact: true }).first().waitFor();
    await dialog.locator('[contenteditable="true"]').fill('富文本正文测试');
    await dialog.locator('[contenteditable="true"]').press('End');
    await page.waitForTimeout(150);
    await page.evaluate(async () => {
      const canvas = document.createElement('canvas');
      canvas.width = canvas.height = 16;
      canvas.getContext('2d').fillRect(0, 0, 16, 16);
      const blob = await new Promise((resolve) => canvas.toBlob(resolve, 'image/png'));
      await navigator.clipboard.write([new ClipboardItem({ 'image/png': blob })]);
    });
    await dialog.locator('[contenteditable="true"]').press('Control+v');
    await dialog.locator('[contenteditable="true"] img[src^="data:image/png;base64,"]').waitFor({ state: 'attached' });
    await dialog.getByPlaceholder('输入作者名称').fill('回归测试');
    await dialog.getByRole('button', { name: '保存文章' }).click();
    await dialog.waitFor({ state: 'hidden' });
    const row = page.getByRole('row').filter({ hasText: title });
    await row.waitFor();
    assert.match(await row.innerText(), /技术分享/);
    articleId = await page.evaluate(async (keyword) => {
      const api = await import('/src/api/system/index.ts');
      return (await api.listArticles({ keyword })).list[0].id;
    }, title);
    await page.evaluate(async (id) => {
      const api = await import('/src/api/system/index.ts');
      const article = await api.getArticle(id);
      await api.updateArticle(id, {
        ...article,
        content: article.content + '<p><strong>格式验证</strong></p><script>window.articleUnsafe = true</script><p><a href="javascript:window.articleUnsafe=true">链接验证</a></p>',
      });
    }, articleId);
    await row.getByRole('button', { name: '查看', exact: true }).click();
    await dialog.getByPlaceholder('输入文章标题').waitFor();
    assert.equal(
      await dialog.getByPlaceholder('输入文章标题').isDisabled(),
      true,
    );
    assert.equal(await dialog.getByPlaceholder('自行填写分类，例如：技术分享').inputValue(), '技术分享');
    assert.match(
      await dialog.locator('.article-content').innerText(),
      /富文本正文测试/,
    );
    const pastedImage = dialog.locator('.article-content img[src^="data:image/png;base64,"]');
    await pastedImage.waitFor();
    assert.equal(await pastedImage.evaluate((img) => img.complete && img.naturalWidth === 16), true);
    assert.equal(await dialog.locator('.article-content strong').innerText(), '格式验证');
    assert.equal(await dialog.locator('.article-content script, .article-content a[href^="javascript:"]').count(), 0);
    assert.equal(await page.evaluate(() => window.articleUnsafe), undefined);
    await dialog.getByRole('button', { name: 'Close' }).click();
    await row.getByRole('button', { name: '编辑', exact: true }).click();
    await dialog.getByPlaceholder('输入文章标题').fill(`${title} 已修改`);
    await dialog.getByPlaceholder('自行填写分类，例如：技术分享').fill('产品更新');
    await dialog.locator('[contenteditable="true"]').fill('修改后的正文');
    await dialog.getByRole('switch').click();
    await dialog.getByRole('button', { name: '保存文章' }).click();
    await dialog.waitFor({ state: 'hidden' });
    const updated = page
      .getByRole('row')
      .filter({ hasText: `${title} 已修改` });
    await updated.waitFor();
    assert.match(await updated.innerText(), /草稿/);
    assert.match(await updated.innerText(), /产品更新/);
    await page.getByPlaceholder('搜索文章标题').fill(title);
    await page.getByPlaceholder('搜索文章标题').press('Enter');
    await updated.waitFor();
    await updated.getByRole('button', { name: '删除', exact: true }).click();
    await page.getByRole('button', { name: /^确\s*定$/ }).click();
    await updated.waitFor({ state: 'hidden' });
    articleId = undefined;
    await page.getByRole('button', { name: '新建文章' }).click();
    await dialog.locator('[contenteditable="true"]').waitFor();
    assert.equal(
      await dialog.getByPlaceholder('输入文章标题').inputValue(),
      '',
    );
    assert.equal(await dialog.getByPlaceholder('自行填写分类，例如：技术分享').inputValue(), '');
    await dialog.getByRole('button', { name: /^取\s*消$/ }).click();
    await dialog.waitFor({ state: 'hidden' });
    // wangEditor's selectionchange callback is debounced by 100 ms.
    await page.waitForTimeout(200);
    assert.deepEqual(errors, []);
    console.log(
      'PASS: no request loop; article create/read/update/search/delete; empty-content validation; editor reopen; no browser errors.',
    );
  } finally {
    if (articleId) {
      await page.screenshot({ path: 'article-test-failure.png' }).catch(() => {});
      await page
        .evaluate(async (id) => {
          const api = await import('/src/api/system/index.ts');
          await api.deleteArticle(id);
        }, articleId)
        .catch(() => {});
    }
    await browser.close();
  }
})().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
