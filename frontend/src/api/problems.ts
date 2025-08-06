import api from './index'

export interface Problem {
  id: number
  title: string
  description: string
  input_desc: string
  output_desc: string
  difficulty: 'Easy' | 'Medium' | 'Hard'
  category: string
  tags: string[]
  time_limit: number
  memory_limit: number
  submit_count: number
  accepted_count: number
  is_public: boolean
  author: any
  created_at: string
}

export interface ProblemList {
  problems: Problem[]
  total: number
  page: number
  limit: number
}

// 获取题目列表
export const getProblems = (params?: {
  page?: number
  limit?: number
  difficulty?: string
  category?: string
}) => {
  return api.get<ProblemList>('/problems', { params }) as unknown as Promise<ProblemList>
}

// 获取题目详情
export const getProblem = (id: number) => {
  return api.get<Problem>(`/problems/${id}`)
}

// 获取题目测试用例
export const getProblemTestCases = (id: number) => {
  return api.get(`/problems/${id}/testcases`)
} 