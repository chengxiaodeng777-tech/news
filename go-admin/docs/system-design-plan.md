# 系统管理模块实现计划

本文档用于规划当前项目的基础系统管理能力，采用经典后端设计，不围绕现有前端 mock 结构做兼容，后续由前端适配后端接口与数据结构。

权限部分直接采用 `Casbin` 作为授权引擎，不实现多租户。

## 一、目标范围

本阶段实现以下基础能力：

- 用户管理
- 角色管理
- 菜单管理
- 权限控制
- 组织管理
- 字典管理
- 登录认证
- 当前用户菜单与权限查询

设计原则：

- 采用经典 `RBAC` 权限模型
- 授权引擎使用 `Casbin`
- 不启用多租户 `domain`
- 后端数据结构优先，前端后续适配
- 菜单不仅表示前端路由，也承载按钮和接口权限点
- 组织与菜单均采用树形结构
- 字典采用“字典类型 + 字典数据”的标准模型

## 二、领域模型设计

计划建设以下核心实体：

### 1. 用户

`User`

建议字段：

- `id`
- `username`
- `password_hash`
- `real_name`
- `nickname`
- `mobile`
- `email`
- `avatar`
- `status`
- `is_super_admin`
- `organization_id`
- `remark`
- `created_at`
- `updated_at`

说明：

- `username` 全局唯一
- `password_hash` 保存密码哈希，不保存明文密码
- `is_super_admin` 用于超级管理员放行
- 用户归属一个组织

### 2. 角色

`Role`

建议字段：

- `id`
- `code`
- `name`
- `status`
- `sort`
- `remark`
- `created_at`
- `updated_at`

说明：

- `code` 用于角色唯一标识
- 一个用户可绑定多个角色

### 3. 菜单

`Menu`

建议字段：

- `id`
- `parent_id`
- `name`
- `title`
- `path`
- `component`
- `icon`
- `type`
- `permission_code`
- `redirect`
- `sort`
- `status`
- `hidden`
- `keep_alive`
- `external_link`
- `remark`
- `created_at`
- `updated_at`

说明：

- `type` 用于区分目录、菜单、按钮、接口
- `permission_code` 用于后端接口权限校验，也可供前端按钮权限判断
- 菜单采用树结构
- 后续前端动态菜单主要依赖本表中的目录和菜单类型数据

### 4. 组织

`Organization`

建议字段：

- `id`
- `parent_id`
- `name`
- `code`
- `type`
- `leader`
- `phone`
- `email`
- `sort`
- `status`
- `remark`
- `created_at`
- `updated_at`

说明：

- 组织采用树结构
- `type` 可表示集团、公司、部门、小组等

### 5. 字典类型

`DictType`

建议字段：

- `id`
- `code`
- `name`
- `status`
- `remark`
- `created_at`
- `updated_at`

### 6. 字典数据

`DictData`

建议字段：

- `id`
- `dict_type_id`
- `label`
- `value`
- `tag_type`
- `css_class`
- `sort`
- `status`
- `is_default`
- `remark`
- `created_at`
- `updated_at`

说明：

- 通过 `dict_type_id` 关联到字典类型
- 适合做状态枚举、菜单类型、组织类型等系统基础数据

### 7. 关联关系

#### 用户角色关联

`UserRole`

字段建议：

- `id`
- `user_id`
- `role_id`

#### 角色菜单关联

`RoleMenu`

字段建议：

- `id`
- `role_id`
- `menu_id`

## 三、Casbin 权限模型设计

本项目采用 `Casbin + RBAC`。

说明：

- 用户、角色、菜单、组织、字典仍然保留为业务实体
- 权限判定统一交给 Casbin
- 不使用多租户，不引入 `domain`
- 不做过度复杂的 `ABAC`

### 1. 模型选择

采用标准 RBAC 模型：

- 用户与角色关系交给 `g` 规则
- 角色与权限关系交给 `p` 规则

建议模型：

- `sub`：主体，通常是角色编码，也可以是用户标识
- `obj`：资源，一般使用权限码
- `act`：动作，例如 `read`、`create`、`update`、`delete`

示例：

- `p, admin, system:user, read`
- `p, admin, system:user, create`
- `g, alice, admin`

### 2. 权限码设计

权限码建议拆成“资源 + 动作”，便于 Casbin 判定：

- `obj = system:user`
- `act = list`
- `act = create`
- `act = update`
- `act = delete`

也可以保留传统单字符串权限码，例如：

- `system:user:list`
- `system:user:create`

如果采用单字符串权限码，则可约定：

- `obj = system:user:list`
- `act = use`

为降低第一版复杂度，建议统一采用以下方式：

- `obj`：权限码，例如 `system:user:list`
- `act`：固定值 `access`

示例：

- `p, admin, system:user:list, access`
- `p, admin, system:user:create, access`
- `g, admin, super_admin`

这样更容易与后台管理系统现有“权限码”思维保持一致。

### 3. 菜单与 Casbin 的关系

菜单仍然保存在业务表中，不直接用 Casbin 替代菜单管理。

职责划分如下：

- `Menu` 表：维护目录、菜单、按钮、接口等结构化数据
- `Casbin policy`：负责判断当前用户是否拥有某个权限码

推荐做法：

1. 菜单表中的按钮/接口类型记录 `permission_code`
2. 角色绑定菜单后，从菜单中提取 `permission_code`
3. 同步写入 Casbin `p` 策略
4. 用户绑定角色后，写入 Casbin `g` 策略

### 4. 菜单树返回策略

前端动态菜单仍由后端业务表生成，不直接从 Casbin policy 生成。

建议流程：

1. 查询用户角色
2. 查询角色关联菜单
3. 从菜单表中筛出目录、菜单类型
4. 构建树形菜单返回前端

Casbin 主要负责：

- 接口权限校验
- 按钮权限判断
- 权限码聚合校验

### 5. Casbin 落库策略

第一版建议直接把 Casbin 策略存数据库。

建议增加策略表：

`CasbinRule`

建议字段：

- `id`
- `ptype`
- `v0`
- `v1`
- `v2`
- `v3`
- `v4`
- `v5`

说明：

- `ptype = p` 表示权限策略
- `ptype = g` 表示角色继承/用户角色关系

### 6. 同步策略原则

以下场景需要同步 Casbin 策略：

- 创建/修改/删除角色
- 给角色分配菜单
- 给用户分配角色
- 修改菜单权限码
- 删除菜单

第一版建议：

- 在角色分配菜单、用户分配角色时同步更新 Casbin
- 不做复杂的增量补偿机制
- 必要时提供“重建策略”功能，按业务表全量重建 Casbin policy

## 四、权限模型设计

采用 `Casbin RBAC`：

- 用户绑定多个角色
- 角色绑定多个菜单
- 权限从角色聚合到用户

权限规则：

1. 用户登录后，先根据用户查询角色
2. 再根据角色查询菜单与权限点
3. 同时基于 Casbin 生成或加载授权关系
4. 聚合得到当前用户菜单树和权限码列表
5. 超级管理员跳过常规权限校验

权限承载方式：

- 页面级权限：菜单类型为“目录/菜单”
- 按钮级权限：菜单类型为“按钮”
- 接口级权限：菜单类型为“接口”

建议菜单类型枚举：

- `catalog`
- `menu`
- `button`
- `api`

## 五、后端接口规划

### 1. 认证接口

- `POST /api/auth/login`
- `POST /api/auth/logout`
- `POST /api/auth/refresh`
- `GET /api/user/info`

### 2. 当前用户权限接口

- `GET /api/menu/all`
- `GET /api/auth/codes`

说明：

- `/api/menu/all` 返回当前用户可见菜单树
- `/api/auth/codes` 返回当前用户所有权限码列表

### 3. 用户管理

- `GET /api/system/users`
- `GET /api/system/users/{id}`
- `POST /api/system/users`
- `PUT /api/system/users/{id}`
- `DELETE /api/system/users/{id}`
- `PUT /api/system/users/{id}/roles`
- `PUT /api/system/users/{id}/password`

### 4. 角色管理

- `GET /api/system/roles`
- `GET /api/system/roles/{id}`
- `POST /api/system/roles`
- `PUT /api/system/roles/{id}`
- `DELETE /api/system/roles/{id}`
- `PUT /api/system/roles/{id}/menus`

### 5. 菜单管理

- `GET /api/system/menus`
- `GET /api/system/menus/tree`
- `GET /api/system/menus/{id}`
- `POST /api/system/menus`
- `PUT /api/system/menus/{id}`
- `DELETE /api/system/menus/{id}`

### 6. 组织管理

- `GET /api/system/organizations`
- `GET /api/system/organizations/tree`
- `GET /api/system/organizations/{id}`
- `POST /api/system/organizations`
- `PUT /api/system/organizations/{id}`
- `DELETE /api/system/organizations/{id}`

### 7. 字典管理

- `GET /api/system/dict/types`
- `GET /api/system/dict/types/{id}`
- `POST /api/system/dict/types`
- `PUT /api/system/dict/types/{id}`
- `DELETE /api/system/dict/types/{id}`
- `GET /api/system/dict/data`
- `GET /api/system/dict/data/{id}`
- `POST /api/system/dict/data`
- `PUT /api/system/dict/data/{id}`
- `DELETE /api/system/dict/data/{id}`
- `GET /api/system/dict/data/by-type/{code}`

## 六、返回数据规划

### 1. 当前用户信息

建议包含：

- 用户基础信息
- 角色编码列表
- 默认首页路径
- 所属组织信息

### 2. 当前用户菜单树

后端返回树形结构即可，建议包含：

- `id`
- `parentId`
- `name`
- `title`
- `path`
- `component`
- `icon`
- `redirect`
- `type`
- `hidden`
- `keepAlive`
- `children`

说明：

- 返回给前端动态菜单时，只返回目录与菜单类型
- 按钮与接口类型用于权限聚合，不作为菜单树直接返回

### 3. 当前用户权限码列表

示例：

- `system:user:create`
- `system:user:update`
- `system:role:assign`
- `system:menu:list`

## 七、目录分层建议

建议增加如下目录：

- `internal/model`
- `internal/repository`
- `internal/service`
- `internal/handler/system`
- `internal/handler/auth`
- `internal/casbin`
- `internal/pkg/response`
- `internal/pkg/page`
- `internal/pkg/password`

职责说明：

- `repository`：数据库读写
- `service`：业务逻辑聚合
- `handler`：HTTP 参数解析与响应输出
- `model`：请求体、响应体、分页对象、视图对象
- `casbin`：Casbin 初始化、模型、策略同步、权限判断封装

## 八、实现顺序

### 第一阶段：模型与基础设施

1. 编写 `ent/schema`
2. 生成 Ent 代码
3. 接入自动迁移或初始化迁移能力
4. 接入 Casbin 模型与策略存储
5. 增加基础种子数据能力

### 第二阶段：认证与权限基础

1. 用户登录
2. JWT 签发与刷新
3. 当前用户信息查询
4. 当前用户权限码聚合
5. 当前用户菜单树聚合
6. Casbin 权限校验中间件

### 第三阶段：系统管理 CRUD

1. 用户管理
2. 角色管理
3. 菜单管理
4. 组织管理
5. 字典管理

### 第四阶段：中间件与通用能力

1. JWT 鉴权中间件
2. Casbin 权限校验中间件
3. 分页查询封装
4. 统一响应结构
5. 错误码与状态码约定

## 九、初始化基础数据

计划初始化以下内容：

### 1. 超级管理员

- 用户名：`admin`
- 密码：初始化默认密码，首次可手动修改
- 标记为超级管理员

### 2. 管理员角色

- 角色编码：`admin`
- 角色名称：`管理员`

### 3. Casbin 初始策略

- 超级管理员角色拥有全部核心系统权限
- `admin` 用户默认绑定管理员角色

### 4. 基础组织

- 根组织节点
- 示例公司/部门节点

### 5. 基础菜单

- 系统管理
- 用户管理
- 角色管理
- 菜单管理
- 组织管理
- 字典管理

### 6. 基础字典

- 用户状态
- 角色状态
- 菜单类型
- 组织类型
- 是否默认

## 十、注意事项

- 密码必须使用哈希保存，不可明文落库
- Casbin 只负责授权，不替代业务表
- 删除组织前要校验是否存在子组织或挂载用户
- 删除角色前要校验是否仍被用户引用
- 删除菜单前要校验是否存在子菜单
- 字典类型删除前要校验是否仍存在字典数据
- 菜单树与组织树接口建议统一按树形输出
- 权限码命名风格全项目统一，例如：`system:user:list`
- 修改角色菜单、用户角色、菜单权限码后，要同步更新 Casbin 策略
- 建议保留“重建 Casbin 策略”能力，便于数据修复

## 十一、下一步执行建议

下一步建议直接进入以下工作：

1. 先创建 `ent/schema`，同时补 `CasbinRule` 表
2. 再补 `repository/service/handler/casbin` 基础目录
3. 然后实现认证、菜单树、权限码、Casbin 校验四个最核心能力

这三个完成后，系统管理模块的骨架就基本成立，后续再继续补用户、角色、组织、字典 CRUD。
