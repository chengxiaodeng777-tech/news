import { requestClient } from '#/api/request';

export interface PageQuery {
  keyword?: string;
  page?: number;
  pageSize?: number;
  status?: number;
}

export interface PageResult<T> {
  list: T[];
  page: number;
  pageSize: number;
  total: number;
}

export interface ArticleRecord {
  sort: number;
  category: string;
  author: string;
  content: string;
  cover: string;
  createdAt: string;
  id: number;
  status: number;
  summary: string;
  title: string;
  updatedAt: string;
}

export interface SaveArticlePayload {
  sort?: number;
  category?: string;
  author?: string;
  content: string;
  cover?: string;
  status: number;
  summary?: string;
  title: string;
}

export interface UserRecord {
  avatar: string;
  email: string;
  id: number;
  mobile: string;
  nickname: string;
  organizationId: number;
  organizationName: string;
  realName: string;
  remark: string;
  roleIds?: number[];
  status: number;
  username: string;
}

export interface SaveUserPayload {
  avatar?: string;
  email?: string;
  mobile?: string;
  nickname?: string;
  organizationId?: number;
  password?: string;
  realName?: string;
  remark?: string;
  status: number;
  username?: string;
}

export interface RoleRecord {
  code: string;
  id: number;
  menuIds?: number[];
  name: string;
  remark: string;
  sort: number;
  status: number;
}

export interface SaveRolePayload {
  code?: string;
  name: string;
  remark?: string;
  sort: number;
  status: number;
}

export interface MenuRecord {
  children?: MenuRecord[];
  component: string;
  externalLink: string;
  hidden: boolean;
  icon: string;
  id: number;
  keepAlive: boolean;
  name: string;
  parentId: number;
  path: string;
  permissionCode: string;
  redirect: string;
  remark: string;
  sort: number;
  status: number;
  title: string;
  type: string;
}

export interface SaveMenuPayload {
  component?: string;
  externalLink?: string;
  hidden: boolean;
  icon?: string;
  keepAlive: boolean;
  name: string;
  parentId: number;
  path?: string;
  permissionCode?: string;
  redirect?: string;
  remark?: string;
  sort: number;
  status: number;
  title: string;
  type: string;
}

export interface OrganizationRecord {
  children?: OrganizationRecord[];
  code: string;
  email: string;
  id: number;
  leader: string;
  name: string;
  parentId: number;
  phone: string;
  remark: string;
  sort: number;
  status: number;
  type: string;
}

export interface SaveOrganizationPayload {
  code: string;
  email?: string;
  leader?: string;
  name: string;
  parentId: number;
  phone?: string;
  remark?: string;
  sort: number;
  status: number;
  type?: string;
}

export interface DictTypeRecord {
  code: string;
  id: number;
  name: string;
  remark: string;
  status: number;
}

export interface SaveDictTypePayload {
  code: string;
  name: string;
  remark?: string;
  status: number;
}

export interface DictDataRecord {
  cssClass: string;
  dictTypeId: number;
  id: number;
  isDefault: boolean;
  label: string;
  remark: string;
  sort: number;
  status: number;
  tagType: string;
  value: string;
}

export interface SaveDictDataPayload {
  cssClass?: string;
  dictTypeId: number;
  isDefault: boolean;
  label: string;
  remark?: string;
  sort: number;
  status: number;
  tagType?: string;
  value: string;
}

export function listUsers(params: PageQuery) {
  return requestClient.get<PageResult<UserRecord>>('/system/users', { params });
}

export function listArticles(params: PageQuery) {
  return requestClient.get<PageResult<ArticleRecord>>('/system/articles', {
    params,
  });
}

export function getArticle(id: number) {
  return requestClient.get<ArticleRecord>(`/system/articles/${id}`);
}

export async function uploadArticleImage(file: File): Promise<string> {
  const result = await requestClient.upload<{ url: string }>('/upload', { file });
  if (!result?.url) throw new Error('图片上传接口未返回图片地址');
  return result.url;
}

export function createArticle(data: SaveArticlePayload) {
  return requestClient.post<boolean>('/system/articles', data);
}

export function updateArticle(id: number, data: SaveArticlePayload) {
  return requestClient.put<boolean>(`/system/articles/${id}`, data);
}

export function deleteArticle(id: number) {
  return requestClient.delete<boolean>(`/system/articles/${id}`);
}

export function createUser(data: SaveUserPayload) {
  return requestClient.post<boolean>('/system/users', data);
}

export function updateUser(id: number, data: SaveUserPayload) {
  return requestClient.put<boolean>(`/system/users/${id}`, data);
}

export function deleteUser(id: number) {
  return requestClient.delete<boolean>(`/system/users/${id}`);
}

export function replaceUserRoles(id: number, roleIds: number[]) {
  return requestClient.put<boolean>(`/system/users/${id}/roles`, { roleIds });
}

export function resetUserPassword(id: number, password: string) {
  return requestClient.put<boolean>(`/system/users/${id}/password`, {
    password,
  });
}

export function listRoles(params: PageQuery) {
  return requestClient.get<PageResult<RoleRecord>>('/system/roles', { params });
}

export function createRole(data: SaveRolePayload) {
  return requestClient.post<boolean>('/system/roles', data);
}

export function updateRole(id: number, data: SaveRolePayload) {
  return requestClient.put<boolean>(`/system/roles/${id}`, data);
}

export function deleteRole(id: number) {
  return requestClient.delete<boolean>(`/system/roles/${id}`);
}

export function replaceRoleMenus(id: number, menuIds: number[]) {
  return requestClient.put<boolean>(`/system/roles/${id}/menus`, { menuIds });
}

export function listMenus() {
  return requestClient.get<MenuRecord[]>('/system/menus');
}

export function listMenuTree() {
  return requestClient.get<MenuRecord[]>('/system/menus/tree');
}

export function createMenu(data: SaveMenuPayload) {
  return requestClient.post<boolean>('/system/menus', data);
}

export function updateMenu(id: number, data: SaveMenuPayload) {
  return requestClient.put<boolean>(`/system/menus/${id}`, data);
}

export function deleteMenu(id: number) {
  return requestClient.delete<boolean>(`/system/menus/${id}`);
}

export function listOrganizations() {
  return requestClient.get<OrganizationRecord[]>('/system/organizations');
}

export function listOrganizationTree() {
  return requestClient.get<OrganizationRecord[]>('/system/organizations/tree');
}

export function createOrganization(data: SaveOrganizationPayload) {
  return requestClient.post<boolean>('/system/organizations', data);
}

export function updateOrganization(id: number, data: SaveOrganizationPayload) {
  return requestClient.put<boolean>(`/system/organizations/${id}`, data);
}

export function deleteOrganization(id: number) {
  return requestClient.delete<boolean>(`/system/organizations/${id}`);
}

export function listDictTypes(params: PageQuery) {
  return requestClient.get<PageResult<DictTypeRecord>>('/system/dict-types', {
    params,
  });
}

export function createDictType(data: SaveDictTypePayload) {
  return requestClient.post<boolean>('/system/dict-types', data);
}

export function updateDictType(id: number, data: SaveDictTypePayload) {
  return requestClient.put<boolean>(`/system/dict-types/${id}`, data);
}

export function deleteDictType(id: number) {
  return requestClient.delete<boolean>(`/system/dict-types/${id}`);
}

export function listDictData(params: PageQuery & { dictTypeId?: number }) {
  return requestClient.get<PageResult<DictDataRecord>>('/system/dict-data', {
    params,
  });
}

export function createDictData(data: SaveDictDataPayload) {
  return requestClient.post<boolean>('/system/dict-data', data);
}

export function updateDictData(id: number, data: SaveDictDataPayload) {
  return requestClient.put<boolean>(`/system/dict-data/${id}`, data);
}

export function deleteDictData(id: number) {
  return requestClient.delete<boolean>(`/system/dict-data/${id}`);
}

export interface MonitorOverview {
  address: string;
  appName: string;
  currentTime: string;
  environment: string;
  requestIdEnabled: boolean;
  startTime: string;
  uptime: string;
}

export interface MonitorDatabase {
  driver: string;
  idleConnections: number;
  inUseConnections: number;
  maxOpenConnections: number;
  openConnections: number;
  status: string;
  waitCount: number;
  waitDuration: string;
}

export interface MonitorRedis {
  address: string;
  db: number;
  idleConns: number;
  pingDuration: string;
  poolHits: number;
  poolMisses: number;
  poolSize: number;
  staleConns: number;
  status: string;
  timeouts: number;
  totalConns: number;
}

export interface MonitorRuntime {
  cpus: number;
  goVersion: string;
  goroutines: number;
  heapAllocMB: number;
  heapSysMB: number;
  lastGC: string;
  memoryAllocMB: number;
  memorySysMB: number;
  numGC: number;
}

export interface MonitorHost {
  arch: string;
  cpuLoad: number;
  freeDiskGB: number;
  freeMemoryMB: number;
  hostname: string;
  os: string;
  totalDiskGB: number;
  totalMemoryMB: number;
}

export interface MonitorRecord {
  database: MonitorDatabase;
  host: MonitorHost;
  overview: MonitorOverview;
  redis: MonitorRedis;
  runtime: MonitorRuntime;
}

export function getMonitor() {
  return requestClient.get<MonitorRecord>('/system/monitor');
}
