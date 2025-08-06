import api from './index'

export interface News {
  id: number
  title: string
  content: string
  summary: string
  category: string
  tags: string[]
  is_published: boolean
  author: any
  created_at: string
  updated_at: string
}

// 获取新闻列表
export const getNews = (params?: {
  page?: number
  limit?: number
}) => {
  return api.get<{
    news: News[]
    total: number
    page: number
    limit: number
  }>('/news', { params }) as unknown as Promise<{
    news: News[]
    total: number
    page: number
    limit: number
  }>
}

// 获取新闻详情
export const getNewsDetail = (id: number) => {
  return api.get<{ news: News }>(`/news/${id}`)
} 