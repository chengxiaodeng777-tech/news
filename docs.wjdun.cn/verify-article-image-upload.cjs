const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');

const bundle = fs.readFileSync('js/article-BCbH4Lzw.js', 'utf8');
assert.match(fs.readFileSync('index.html', 'utf8'), /\/jse\/index-index-DjwYEQvf\.js/);
assert.match(fs.readFileSync('jse/index-index-DjwYEQvf.js', 'utf8'), /bootstrap-DLLMD7gc\.js/);
assert.match(fs.readFileSync('js/bootstrap-DLLMD7gc.js', 'utf8'), /article-BCbH4Lzw\.js/);
assert.match(bundle, /base64LimitSize:0,allowedFileTypes/);
assert.match(bundle, /O\.content=yield replaceEmbeddedImages\(O\.content\)/);
assert.doesNotMatch(bundle, /图片大小不能超过 5 MB/);
assert.doesNotMatch(bundle, /Nt\.error\(error instanceof Error\?error\.message:`图片上传失败`\);throw error/);

const start = bundle.indexOf('function uploadArticleImage(file){');
const end = bundle.indexOf('async function replaceEmbeddedImages(', start);
assert.ok(start >= 0 && end > start);
const source = bundle.slice(start, end) + '\nthis.uploadArticleImage=uploadArticleImage;';

async function verify() {
  const success = { uploadArticleFile: ({ onSuccess }) => onSuccess({ data: { url: '/uploads/example.png' } }) };
  vm.runInNewContext(source, success);
  assert.equal(await success.uploadArticleImage({}), '/uploads/example.png');

  const missingUrl = { uploadArticleFile: ({ onSuccess }) => onSuccess({ data: {} }) };
  vm.runInNewContext(source, missingUrl);
  await assert.rejects(missingUrl.uploadArticleImage({}), /未返回图片地址/);

  const failure = { uploadArticleFile: ({ onError }) => onError(new Error('upload failed')) };
  vm.runInNewContext(source, failure);
  await assert.rejects(failure.uploadArticleImage({}), /upload failed/);

  console.log('Article image upload checks passed');
}

verify().catch(error => {
  console.error(error);
  process.exitCode = 1;
});
