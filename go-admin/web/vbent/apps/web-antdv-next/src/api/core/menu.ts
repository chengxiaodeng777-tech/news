import type { RouteRecordStringComponent } from '@vben/types';

import { requestClient } from '#/api/request';

interface BackendMenuMeta {
  hidden?: boolean;
  icon?: string;
  title: string;
}

interface BackendMenu {
  children?: BackendMenu[];
  component: string;
  meta: BackendMenuMeta;
  name: string;
  path: string;
  redirect?: string;
}

function transformBackendMenu(menu: BackendMenu): RouteRecordStringComponent {
  const { children, component, meta, redirect, ...route } = menu;

  return {
    ...route,
    ...(redirect ? { redirect } : {}),
    children: children?.map(transformBackendMenu),
    component: component === 'LAYOUT' ? 'BasicLayout' : component,
    meta: {
      hideInMenu: meta.hidden,
      icon: meta.icon,
      title: meta.title,
    },
  };
}

export async function getAllMenusApi() {
  const menus = await requestClient.get<BackendMenu[]>('/menu/all');
  return menus.map(transformBackendMenu);
}
