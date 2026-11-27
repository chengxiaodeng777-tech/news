// Patch the deployed article bundle; this checkout does not include its source project.
const fs = require('node:fs');
const path = require('node:path');

const file = path.join(__dirname, 'js/article-BCbH4Lzw.js');
let bundle = fs.readFileSync(file, 'utf8');

function replaceOnce(before, after) {
  const first = bundle.indexOf(before);
  if (first < 0 || bundle.indexOf(before, first + before.length) >= 0) {
    throw new Error(`Expected one occurrence of: ${before.slice(0, 80)}`);
  }
  bundle = bundle.replace(before, after);
}

replaceOnce(
  'import{D as re,_ as ie,c as ae,h as oe,t as se}from"./system-D05vO9Nm.js";',
  'import{D as re,I as uploadArticleFile,_ as ie,c as ae,h as oe,t as se}from"./system-D05vO9Nm.js";',
);

replaceOnce(
  'qOe={class:`editor-shell`},JOe=',
  `qOe={class:\`editor-shell\`};
function uploadArticleImage(file){
  return new Promise((resolve,reject)=>{
    uploadArticleFile({file,onSuccess(response){
      const data=response?.data??response;
      const url=typeof data===\`string\`?data:data?.url??data?.fileUrl??data?.path??data?.file?.url;
      if(typeof url!==\`string\`||!url.trim()){
        reject(new Error(\`图片上传接口未返回图片地址\`));
        return;
      }
      resolve(url.trim());
    },onError:reject});
  });
}
async function replaceEmbeddedImages(html){
  if(!html.includes(\`data:image/\`))return html;
  const documentHtml=new DOMParser().parseFromString(html,\`text/html\`);
  for(const image of documentHtml.querySelectorAll(\`img[src^="data:image/"]\`)){
    const blob=await fetch(image.getAttribute(\`src\`)).then(response=>response.blob());
    const extension=(blob.type.split(\`/\`)[1]||\`png\`).replace(/[^a-z0-9]/gi,\`\`)||\`png\`;
    const file=new File([blob],\`article-image-\${Date.now()}.\${extension}\`,{type:blob.type});
    image.setAttribute(\`src\`,await uploadArticleImage(file));
  }
  return documentHtml.body.innerHTML;
}
var JOe=`,
);

replaceOnce(
  'uploadImage:{base64LimitSize:5*1024*1024,allowedFileTypes:[`image/*`],customUpload(){Nt.warning(`图片大小不能超过 5 MB`)}}',
  'uploadImage:{base64LimitSize:0,allowedFileTypes:[`image/png`,`image/jpeg`,`image/gif`,`image/webp`],async customUpload(file,insertFn){try{insertFn(await uploadArticleImage(file),file.name,``)}catch(error){Nt.error(error instanceof Error?error.message:`图片上传失败`)}}}',
);

replaceOnce(
  'r.value=!0;try{p.value?(yield re(p.value,O)',
  'r.value=!0;try{O.content=yield replaceEmbeddedImages(O.content);p.value?(yield re(p.value,O)',
);

replaceOnce(
  'yield I()}catch(e){}finally{r.value=!1}}}),me.apply(this,arguments)}',
  'yield I()}catch(e){Nt.error(e instanceof Error?e.message:`文章保存失败`)}finally{r.value=!1}}}),me.apply(this,arguments)}',
);

fs.writeFileSync(file, bundle, 'utf8');
console.log(`Patched ${file}`);
