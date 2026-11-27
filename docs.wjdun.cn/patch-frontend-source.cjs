const fs = require('node:fs');
const file = 'backend-go-admin/web/vbent/apps/web-antdv-next/src/views/system/article/index.vue';
let source = fs.readFileSync(file, 'utf8');
const original = /      \/\/ Store pasted\/local images in the article HTML until an upload service is available\.\r?\n      base64LimitSize: 5 \* 1024 \* 1024,\r?\n      allowedFileTypes: \['image\/\*'\],\r?\n      customUpload\(\) \{\r?\n        message\.warning\([^\r\n]*\);\r?\n      \},/;
if (!original.test(source)) throw new Error('Expected original image upload configuration');
source = source.replace(original, `      base64LimitSize: 0,
      allowedFileTypes: ['image/png', 'image/jpeg', 'image/gif', 'image/webp'],
      async customUpload(file, insertFn) {
        try {
          insertFn(await uploadArticleImage(file), file.name, '');
        } catch (error) {
          message.error(error instanceof Error ? error.message : '图片上传失败');
          throw error;
        }
      },`);
fs.writeFileSync(file, source, 'utf8');
