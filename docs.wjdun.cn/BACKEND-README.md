# 文章发布后端修复

此包仅包含 Go 后端程序和源码；前端文件在独立的 `frontend-article-publish.zip` 中。

## 部署

1. 备份并停止现有 Go 服务。
2. 用 `go-admin-linux-amd64` 替换原 Linux amd64 服务程序，赋予执行权限。保留服务器上的 `configs` 和数据库。
3. 设置一个持久化、可写的图片目录，例如 `ARTICLE_UPLOAD_DIR=/var/lib/go-admin/uploads/articles`；未设置时使用服务工作目录下的 `uploads/articles`。
4. 重启服务，确认反向代理将 `/api/upload` 与 `/api/uploads/articles/*` 转发到该服务。代理的请求体大小限制应按实际图片大小配置。
5. 再部署独立的前端包，并清理旧 JS 缓存。

部署后执行 `curl -i -X POST https://docs.wjdun.cn/api/upload`。未登录时新后端会返回 401；若仍返回 404，说明请求没有到达新后端，请检查正在运行的程序和反向代理转发规则。

后端使用数据库中的当前角色菜单权限校验文章接口，避免权限变更后仍使用旧缓存；新增文章图片上传与读取接口。普通用户仍需拥有文章创建或编辑权限。

`source/` 是修复后的 Go 源码，不含服务器配置文件和依赖缓存。
