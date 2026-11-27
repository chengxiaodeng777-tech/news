import type { Recordable } from '@vben/types';

import { requestClient } from '#/api/request';

export namespace ArticleApi {
  export interface Article {
    sort: number;
    category: string;
    id: number;
    title: string;
    summary: string;
    content: string;
    cover: string;
    status: 0 | 1;
    author: string;
    createdAt: string;
    updatedAt: string;
  }

  export interface ArticlePage {
    list: Article[];
    total: number;
    page: number;
    pageSize: number;
  }
}

export function getArticleList(params: Recordable<any>) {
  return requestClient.get<ArticleApi.ArticlePage>('/system/articles', { params });
}

export function getArticle(id: number) {
  return requestClient.get<ArticleApi.Article>(`/system/articles/${id}`);
}

export function createArticle(data: Omit<ArticleApi.Article, 'id' | 'createdAt' | 'updatedAt'>) {
  return requestClient.post('/system/articles', data);
}

export function updateArticle(id: number, data: Omit<ArticleApi.Article, 'id' | 'createdAt' | 'updatedAt'>) {
  return requestClient.put(`/system/articles/${id}`, data);
}

export function deleteArticle(id: number) {
  return requestClient.delete(`/system/articles/${id}`);
}
