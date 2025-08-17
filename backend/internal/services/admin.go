package services

// AdminRepository 定义管理端所需的统计仓储接口
type AdminRepository interface {
	CountUsers() (int64, error)
	CountProblems() (int64, error)
	CountSubmissions() (int64, error)
}

// AdminService 管理端服务
type AdminService struct {
	repo AdminRepository
}

// NewAdminService 创建管理端服务
func NewAdminService(repo AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

// SystemStats 系统统计数据
type SystemStats struct {
	Users        int64 `json:"users"`
	Problems     int64 `json:"problems"`
	Submissions  int64 `json:"submissions"`
	JudgesOnline int64 `json:"judges_online"`
}

// GetSystemStats 返回系统统计数据
func (s *AdminService) GetSystemStats() (*SystemStats, error) {
	users, err := s.repo.CountUsers()
	if err != nil {
		return nil, err
	}
	problems, err := s.repo.CountProblems()
	if err != nil {
		return nil, err
	}
	submissions, err := s.repo.CountSubmissions()
	if err != nil {
		return nil, err
	}

	// TODO: 与判题服务联动统计在线评测机数量
	judgesOnline := int64(1)

	return &SystemStats{
		Users:        users,
		Problems:     problems,
		Submissions:  submissions,
		JudgesOnline: judgesOnline,
	}, nil
}
