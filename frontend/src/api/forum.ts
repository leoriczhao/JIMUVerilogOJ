import api from './index'

export interface ForumPost {
  id: number
  title: string
  content: string
  category: string
  tags: string
  user: any
  reply_count: number
  view_count: number
  is_locked: boolean
  created_at: string
  updated_at: string
}

export interface ForumReply {
  id: number
  content: string
  author: any
  parent_id: number
  post_id: number
  created_at: string
  updated_at: string
}

// 获取帖子列表
export const getForumPosts = () => {
  return api.get<{
    posts: ForumPost[]
    total: number
    page: number
    limit: number
  }>('/forum/posts') as unknown as Promise<{
    posts: ForumPost[]
    total: number
    page: number
    limit: number
  }>
}

// 获取帖子详情
export const getForumPost = (id: number) => {
  return api.get<ForumPost>(`/forum/posts/${id}`)
}

// 获取帖子回复列表
export const getForumReplies = (id: number, params?: {
  page?: number
  limit?: number
}) => {
  return api.get<{
    replies: ForumReply[]
    total: number
    page: number
    limit: number
  }>(`/forum/posts/${id}/replies`, { params })
} 