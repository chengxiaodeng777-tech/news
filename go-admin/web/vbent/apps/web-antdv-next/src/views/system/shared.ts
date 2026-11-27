export interface TreeNode {
  children?: TreeNode[];
  id: number;
  name?: string;
  title?: string;
}

export interface AntTreeNode {
  children?: AntTreeNode[];
  disabled?: boolean;
  key: number;
  title: string;
  value?: number;
}

export const statusOptions = [
  { label: '启用', value: 1 },
  { label: '停用', value: 0 },
];

export const statusTagOptions = [
  { color: 'success', label: '启用', value: 1 },
  { color: 'error', label: '停用', value: 0 },
];

export type StatusValue = boolean | number | string | null | undefined;

export function normalizeStatus(status: StatusValue, fallback = 1) {
  if (status === true) return 1;
  if (status === false) return 0;

  const value = Number(status);
  if (value === 1 || value === 0) {
    return value;
  }
  return fallback;
}

export function normalizeNumberArray(value?: Array<number | string>) {
  return (value ?? [])
    .map((item) => Number(item))
    .filter((item) => Number.isFinite(item));
}

export function toGridResult<T>(list: T[], total?: number) {
  return {
    items: list,
    total: total ?? list.length,
  };
}

export function normalizeOptionalStatus(status: StatusValue) {
  if (status === null || status === undefined || status === '') {
    return undefined;
  }
  return normalizeStatus(status);
}

export function statusText(status: StatusValue) {
  const value = normalizeOptionalStatus(status);
  if (value === 1) return '启用';
  if (value === 0) return '停用';
  return '-';
}

export function statusColor(status: StatusValue) {
  const value = normalizeOptionalStatus(status);
  if (value === 1) return 'success';
  if (value === 0) return 'error';
  return 'default';
}

export function toTreeOptions<T extends TreeNode>(
  nodes: T[] = [],
  disabledId?: number,
): AntTreeNode[] {
  return nodes.map((node) => ({
    children: toTreeOptions(node.children as T[] | undefined, disabledId),
    disabled: node.id === disabledId,
    key: node.id,
    title: node.title || node.name || String(node.id),
    value: node.id,
  }));
}

export function toTreeData<T extends TreeNode>(nodes: T[] = []): AntTreeNode[] {
  return nodes.map((node) => ({
    children: toTreeData(node.children as T[] | undefined),
    key: node.id,
    title: node.title || node.name || String(node.id),
  }));
}

export function flattenTree<T extends TreeNode>(nodes: T[] = []): T[] {
  return nodes.flatMap((node) => [
    node,
    ...flattenTree(node.children as T[] | undefined),
  ]);
}

export function omitEmpty<T extends Record<string, any>>(data: T): T {
  return Object.fromEntries(
    Object.entries(data).filter(([, value]) => value !== undefined),
  ) as T;
}
