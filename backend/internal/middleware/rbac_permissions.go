package middleware

// 权限常量定义
// 采用 resource.action[.condition] 格式命名

// ========== 用户相关权限 ==========
const (
	// 用户资料权限
	PermUserProfileRead    = "user.profile.read"
	PermUserProfileUpdate  = "user.profile.update"
	PermUserPasswordChange = "user.password.change"
	PermUserAvatarUpload   = "user.avatar.upload"

	// 用户管理权限
	PermUserCreate = "user.create"
	PermUserRead   = "user.read"
	PermUserUpdate = "user.update"
	PermUserDelete = "user.delete"
	PermUserBan    = "user.ban"
	PermUserUnban  = "user.unban"
)

// ========== 题目相关权限 ==========
const (
	// 题目基础权限
	PermProblemCreate = "problem.create"
	PermProblemRead   = "problem.read"
	PermProblemList   = "problem.list"

	// 题目更新权限
	PermProblemUpdateOwn = "problem.update.own"
	PermProblemUpdateAll = "problem.update.all"

	// 题目删除权限
	PermProblemDeleteOwn = "problem.delete.own"
	PermProblemDeleteAll = "problem.delete.all"

	// 题目发布权限
	PermProblemPublish = "problem.publish"
	PermProblemArchive = "problem.archive"

	// 题目管理权限
	PermProblemManageAll = "problem.manage.all"
)

// ========== 测试用例相关权限 ==========
const (
	// 测试用例查看权限
	PermTestcaseReadSample = "testcase.read.sample"
	PermTestcaseReadOwn    = "testcase.read.own"
	PermTestcaseReadAll    = "testcase.read.all"

	// 测试用例管理权限
	PermTestcaseCreate = "testcase.create"
	PermTestcaseUpdate = "testcase.update"
	PermTestcaseDelete = "testcase.delete"
)

// ========== 提交相关权限 ==========
const (
	// 提交基础权限
	PermSubmissionCreate = "submission.create"
	PermSubmissionRead   = "submission.read"
	PermSubmissionList   = "submission.list"

	// 提交管理权限
	PermSubmissionManage  = "submission.manage"
	PermSubmissionDelete  = "submission.delete"
	PermSubmissionRejudge = "submission.rejudge"
)

// ========== 论坛相关权限 ==========
const (
	// 发帖权限
	PermForumPostCreate = "forum.post.create"
	PermForumPostRead   = "forum.post.read"

	// 回复权限
	PermForumReplyCreate = "forum.reply.create"
	PermForumReplyRead   = "forum.reply.read"

	// 编辑权限
	PermForumEditOwn = "forum.edit.own"
	PermForumEditAll = "forum.edit.all"

	// 论坛管理权限
	PermForumModerate   = "forum.moderate"
	PermForumPostLock   = "forum.post.lock"
	PermForumPostSticky = "forum.post.sticky"
	PermForumDelete     = "forum.delete"
)

// ========== 新闻相关权限 ==========
const (
	// 新闻查看权限
	PermNewsRead = "news.read"
	PermNewsList = "news.list"

	// 新闻管理权限
	PermNewsCreate  = "news.create"
	PermNewsUpdate  = "news.update"
	PermNewsDelete  = "news.delete"
	PermNewsPublish = "news.publish"
	PermNewsArchive = "news.archive"
)

// ========== 系统管理权限 ==========
const (
	// 用户管理
	PermManageUsers       = "manage.users"
	PermManageRoles       = "manage.roles"
	PermManagePermissions = "manage.permissions"

	// 系统配置
	PermManageSystem = "manage.system"
	PermManageConfig = "manage.config"
	PermManageLogs   = "manage.logs"

	// 内容管理
	PermManageContent = "manage.content"
	PermManageNews    = "manage.news"
	PermManageForum   = "manage.forum"

	// 统计信息
	PermStatsRead  = "stats.read"
	PermStatsAdmin = "stats.admin"
)

// ========== 通配符权限定义 ==========
const (
	// 全权限
	PermAll = "*"

	// 资源级通配符
	PermUserAll       = "user.*"
	PermProblemAll    = "problem.*"
	PermTestcaseAll   = "testcase.*"
	PermSubmissionAll = "submission.*"
	PermForumAll      = "forum.*"
	PermNewsAll       = "news.*"
	PermManageAll     = "manage.*"
	PermStatsAll      = "stats.*"

	// 用户资料通配符
	PermUserProfileAll = "user.profile.*"

	// 题目操作通配符
	PermProblemUpdate = "problem.update.*"
	PermProblemDelete = "problem.delete.*"

	// 论坛操作通配符
	PermForumPostAll  = "forum.post.*"
	PermForumReplyAll = "forum.reply.*"
)

// ========== 权限分组 ==========
var (
	// 基础用户权限
	UserBasicPermissions = []string{
		PermUserProfileRead,
		PermUserProfileUpdate,
		PermUserPasswordChange,
		PermUserAvatarUpload,
	}

	// 学生权限
	StudentPermissions = append(UserBasicPermissions,
		PermProblemRead,
		PermProblemList,
		PermTestcaseReadSample,
		PermSubmissionCreate,
		PermSubmissionRead,
		PermSubmissionList,
		PermForumPostCreate,
		PermForumPostRead,
		PermForumReplyCreate,
		PermForumReplyRead,
		PermNewsRead,
		PermNewsList,
		PermStatsRead,
	)

	// 教师权限（继承学生权限）
	TeacherPermissions = append(StudentPermissions,
		PermProblemCreate,
		PermProblemUpdateOwn,
		PermProblemDeleteOwn,
		PermTestcaseReadOwn,
		PermTestcaseCreate,
		PermTestcaseUpdate,
		PermTestcaseDelete,
		PermNewsCreate,
		PermNewsUpdate,
		PermNewsDelete,
	)

	// 管理员权限
	AdminPermissions = append(TeacherPermissions,
		PermUserCreate,
		PermUserRead,
		PermUserUpdate,
		PermUserDelete,
		PermUserBan,
		PermUserUnban,
		PermProblemUpdateAll,
		PermProblemDeleteAll,
		PermProblemPublish,
		PermProblemArchive,
		PermTestcaseReadAll,
		PermSubmissionManage,
		PermSubmissionDelete,
		PermSubmissionRejudge,
		PermForumEditAll,
		PermForumModerate,
		PermForumPostLock,
		PermForumPostSticky,
		PermForumDelete,
		PermNewsPublish,
		PermNewsArchive,
		PermManageUsers,
		PermManageSystem,
		PermManageConfig,
		PermManageContent,
		PermStatsAdmin,
	)

	// 超级管理员权限
	SuperAdminPermissions = []string{PermAll}
)

// ========== 权限描述 ==========
var PermissionDescriptions = map[string]string{
	// 用户权限
	PermUserProfileRead:    "查看个人信息",
	PermUserProfileUpdate:  "更新个人信息",
	PermUserPasswordChange: "修改密码",
	PermUserAvatarUpload:   "上传头像",

	// 题目权限
	PermProblemCreate:    "创建题目",
	PermProblemRead:      "查看题目",
	PermProblemUpdateOwn: "更新自己的题目",
	PermProblemUpdateAll: "更新任意题目",
	PermProblemDeleteOwn: "删除自己的题目",
	PermProblemDeleteAll: "删除任意题目",

	// 测试用例权限
	PermTestcaseReadSample: "查看样例测试用例",
	PermTestcaseReadOwn:    "查看自己题目的测试用例",
	PermTestcaseReadAll:    "查看所有测试用例",

	// 提交权限
	PermSubmissionCreate: "提交代码",
	PermSubmissionRead:   "查看提交记录",
	PermSubmissionManage: "管理提交记录",

	// 论坛权限
	PermForumPostCreate:  "发帖",
	PermForumReplyCreate: "回复",
	PermForumEditOwn:     "编辑自己的帖子",
	PermForumModerate:    "管理论坛",

	// 管理权限
	PermManageUsers:  "管理用户",
	PermManageSystem: "管理系统",
	PermStatsRead:    "查看统计信息",

	// 通配符权限
	PermAll: "所有权限",
}
