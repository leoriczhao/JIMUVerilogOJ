<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getProblems, type Problem } from '@/api/problems'
import { getNews, type News } from '@/api/news'
import { getForumPosts, type ForumPost } from '@/api/forum'
import { ElMessage } from 'element-plus'

const problems = ref<Problem[]>([])
const newsList = ref<News[]>([])
const forumPosts = ref<ForumPost[]>([])

// 获取难度类型
const getDifficultyType = (difficulty: string) => {
  switch (difficulty) {
    case 'Easy': return 'success'
    case 'Medium': return 'warning'
    case 'Hard': return 'danger'
    default: return 'info'
  }
}

// 格式化日期
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('zh-CN')
}

// 加载数据
const loadData = async () => {
  try {
    // 分别加载数据，避免一个失败影响其他
    const loadProblems = async () => {
      try {
        const res = await getProblems({ page: 1, limit: 6 })
        problems.value = res?.problems || []
      } catch (error) {
        console.error('加载题目失败:', error)
        problems.value = []
      }
    }

    const loadNews = async () => {
      try {
        const res = await getNews({ page: 1, limit: 4 })
        newsList.value = res?.news || []
      } catch (error) {
        console.error('加载新闻失败:', error)
        newsList.value = []
      }
    }

    const loadForum = async () => {
      try {
        const res = await getForumPosts()
        forumPosts.value = res?.posts?.slice(0, 4) || []
      } catch (error) {
        console.error('加载论坛失败:', error)
        forumPosts.value = []
      }
    }

    // 并行加载，但各自处理错误
    await Promise.all([
      loadProblems(),
      loadNews(),
      loadForum()
    ])
  } catch (error) {
    console.error('加载数据失败:', error)
    ElMessage.error('部分数据加载失败，请稍后重试')
  }
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="home">
    <!-- 欢迎横幅 -->
    <div class="hero-section">
      <div class="hero-content">
        <h1>欢迎来到 Verilog OJ</h1>
        <p>专业的Verilog在线判题系统，提升您的硬件设计技能</p>
        <div class="hero-buttons">
          <el-button type="primary" size="large" @click="$router.push('/problems')">
            开始刷题
          </el-button>
          <el-button size="large" @click="$router.push('/register')">
            立即注册
          </el-button>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <div class="content-container">
        <!-- 题库部分 -->
        <div class="section">
          <div class="section-header">
            <h2>最新题目</h2>
            <el-button text @click="$router.push('/problems')">
              查看更多 <el-icon><ArrowRight /></el-icon>
            </el-button>
          </div>
          <div class="problems-grid">
            <el-card 
              v-for="problem in problems" 
              :key="problem.id" 
              class="problem-card"
              @click="$router.push(`/problems/${problem.id}`)"
            >
              <div class="problem-info">
                <h3>{{ problem.title }}</h3>
                <p class="problem-desc">{{ problem.description }}</p>
                <div class="problem-meta">
                  <el-tag :type="getDifficultyType(problem.difficulty)">
                    {{ problem.difficulty }}
                  </el-tag>
                  <span class="stats">
                    <el-icon><View /></el-icon> {{ problem.submit_count }}
                    <el-icon><Check /></el-icon> {{ problem.accepted_count }}
                  </span>
                </div>
              </div>
            </el-card>
          </div>
        </div>

        <!-- 新闻部分 -->
        <div class="section">
          <div class="section-header">
            <h2>最新新闻</h2>
            <el-button text @click="$router.push('/news')">
              查看更多 <el-icon><ArrowRight /></el-icon>
            </el-button>
          </div>
          <div class="news-list">
            <el-card 
              v-for="news in newsList" 
              :key="news.id" 
              class="news-card"
              @click="$router.push(`/news/${news.id}`)"
            >
              <div class="news-content">
                <h3>{{ news.title }}</h3>
                <p>{{ news.summary || news.content.substring(0, 100) }}...</p>
                <div class="news-meta">
                  <span class="date">{{ formatDate(news.created_at) }}</span>
                  <el-tag v-if="news.category" size="small">{{ news.category }}</el-tag>
                </div>
              </div>
            </el-card>
          </div>
        </div>

        <!-- 论坛部分 -->
        <div class="section">
          <div class="section-header">
            <h2>热门讨论</h2>
            <el-button text @click="$router.push('/forum')">
              查看更多 <el-icon><ArrowRight /></el-icon>
            </el-button>
          </div>
          <div class="forum-list">
            <el-card 
              v-for="post in forumPosts" 
              :key="post.id" 
              class="forum-card"
              @click="$router.push(`/forum/posts/${post.id}`)"
            >
              <div class="forum-content">
                <h3>{{ post.title }}</h3>
                <p>{{ post.content.substring(0, 80) }}...</p>
                <div class="forum-meta">
                  <div class="author">
                    <el-avatar :size="24" :src="post.user?.avatar">
                      {{ post.user?.nickname?.charAt(0) || post.user?.username?.charAt(0) }}
                    </el-avatar>
                    <span>{{ post.user?.nickname || post.user?.username }}</span>
                  </div>
                  <div class="stats">
                    <span><el-icon><ChatDotRound /></el-icon> {{ post.reply_count }}</span>
                    <span><el-icon><View /></el-icon> {{ post.view_count }}</span>
                  </div>
                </div>
              </div>
            </el-card>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.home {
  min-height: 100vh;
}

.hero-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 120px 0 80px;
  text-align: center;
}

.hero-content {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 20px;
}

.hero-content h1 {
  font-size: 64px;
  font-weight: bold;
  margin-bottom: 24px;
  line-height: 1.2;
}

.hero-content p {
  font-size: 24px;
  margin-bottom: 48px;
  opacity: 0.9;
  line-height: 1.5;
}

.hero-buttons {
  display: flex;
  gap: 24px;
  justify-content: center;
  flex-wrap: wrap;
}

.hero-buttons .el-button {
  padding: 12px 32px;
  font-size: 18px;
  border-radius: 8px;
}

.main-content {
  padding: 60px 0;
  background: #f8f9fa;
}

.content-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 40px;
}

.section {
  margin-bottom: 60px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.section-header h2 {
  font-size: 36px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0;
}

.problems-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 24px;
}

.problem-card {
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;
}

.problem-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.problem-info h3 {
  margin: 0 0 12px 0;
  color: #2c3e50;
  font-size: 22px;
  font-weight: 600;
}

.problem-desc {
  color: #7f8c8d;
  margin-bottom: 15px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.problem-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stats {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #95a5a6;
  font-size: 14px;
}

.news-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 24px;
}

.news-card {
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;
}

.news-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
}

.news-content h3 {
  margin: 0 0 12px 0;
  color: #2c3e50;
  font-size: 20px;
  font-weight: 600;
}

.news-content p {
  color: #7f8c8d;
  margin-bottom: 15px;
  line-height: 1.5;
}

.news-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.date {
  color: #95a5a6;
  font-size: 14px;
}

.forum-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 24px;
}

.forum-card {
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;
}

.forum-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
}

.forum-content h3 {
  margin: 0 0 12px 0;
  color: #2c3e50;
  font-size: 20px;
  font-weight: 600;
}

.forum-content p {
  color: #7f8c8d;
  margin-bottom: 15px;
  line-height: 1.5;
}

.forum-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.author {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #7f8c8d;
}

.stats {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #95a5a6;
  font-size: 14px;
}

/* 桌面端优化 */
@media (min-width: 1200px) {
  .hero-content h1 {
    font-size: 72px;
  }
  
  .hero-content p {
    font-size: 28px;
  }
  
  .problems-grid {
    grid-template-columns: repeat(3, 1fr);
  }
  
  .news-list {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .forum-list {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .content-container {
    padding: 0 60px;
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  .hero-content h1 {
    font-size: 32px;
  }
  
  .hero-content p {
    font-size: 16px;
  }
  
  .problems-grid,
  .news-list,
  .forum-list {
    grid-template-columns: 1fr;
  }
  
  .section-header {
    flex-direction: column;
    gap: 15px;
    align-items: flex-start;
  }
  
  .content-container {
    padding: 0 20px;
  }
}
</style>
