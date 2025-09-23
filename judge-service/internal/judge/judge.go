package judge

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// JudgeRequest 判题请求结构
type JudgeRequest struct {
	SubmissionID string      `json:"submission_id"`
	Code         string      `json:"code"`
	Language     string      `json:"language"`
	TimeLimit    int         `json:"time_limit"`    // 毫秒
	MemoryLimit  int         `json:"memory_limit"`  // MB
	TestCases    []TestCase  `json:"test_cases"`
}

// TestCase Verilog测试用例结构
type TestCase struct {
	Testbench    string `json:"testbench"`     // testbench代码
	ExpectedVCD  string `json:"expected_vcd"`  // 期望的VCD文件内容或关键信号值
	Description  string `json:"description"`   // 测试用例描述
	SimTime      int    `json:"sim_time"`      // 仿真时间（时间单位）
}

// JudgeResult 判题结果结构
type JudgeResult struct {
	SubmissionID  string    `json:"submission_id"`
	Status        string    `json:"status"`
	Score         int       `json:"score"`
	RunTime       int       `json:"run_time"`       // 毫秒
	Memory        int       `json:"memory"`         // KB
	ErrorMessage  string    `json:"error_message"`
	PassedTests   int       `json:"passed_tests"`
	TotalTests    int       `json:"total_tests"`
	JudgedAt      time.Time `json:"judged_at"`
}

// Judge 判题器结构
type Judge struct {
	workDir string
}

// NewJudge 创建新的判题器
func NewJudge(workDir string) *Judge {
	return &Judge{
		workDir: workDir,
	}
}

// Judge 执行判题
func (j *Judge) Judge(ctx context.Context, req *JudgeRequest) (*JudgeResult, error) {
	result := &JudgeResult{
		SubmissionID: req.SubmissionID,
		TotalTests:   len(req.TestCases),
		JudgedAt:     time.Now(),
	}

	// 创建临时工作目录
	tempDir, err := j.createTempDir(req.SubmissionID)
	if err != nil {
		result.Status = "system_error"
		result.ErrorMessage = fmt.Sprintf("Failed to create temp directory: %v", err)
		return result, nil
	}
	defer os.RemoveAll(tempDir)

	// 编译代码 - 注意：这里需要从测试用例中获取testbench
	// 暂时使用第一个测试用例的testbench进行编译检查
	if len(req.TestCases) == 0 {
		result.Status = "system_error"
		result.ErrorMessage = "No test cases provided"
		return result, nil
	}
	if err := j.compileVerilog(tempDir, req.Code, req.TestCases[0].Testbench); err != nil {
		result.Status = "compile_error"
		result.ErrorMessage = err.Error()
		return result, nil
	}

	// 运行测试用例
	passed := 0
	totalRunTime := 0
	maxMemory := 0

	for i, testCase := range req.TestCases {
		// 检查上下文是否取消
		select {
		case <-ctx.Done():
			result.Status = "system_error"
			result.ErrorMessage = "Judge timeout"
			return result, nil
		default:
		}

		// 运行单个测试用例
		testResult, err := j.runSingleTest(ctx, tempDir, testCase, req.Code, req.TimeLimit, req.MemoryLimit)
		if err != nil {
			result.Status = "system_error"
			result.ErrorMessage = fmt.Sprintf("Test case %d failed: %v", i+1, err)
			return result, nil
		}

		totalRunTime += testResult.RunTime
		if testResult.Memory > maxMemory {
			maxMemory = testResult.Memory
		}

		if testResult.Status == "accepted" {
			passed++
		} else if result.Status == "" {
			// 设置第一个失败的状态
			result.Status = testResult.Status
			result.ErrorMessage = testResult.ErrorMessage
		}
	}

	result.PassedTests = passed
	result.RunTime = totalRunTime
	result.Memory = maxMemory
	result.Score = (passed * 100) / len(req.TestCases)

	if passed == len(req.TestCases) {
		result.Status = "accepted"
		result.ErrorMessage = ""
	}

	return result, nil
}

// createTempDir 创建临时工作目录
func (j *Judge) createTempDir(submissionID string) (string, error) {
	tempDir := filepath.Join(j.workDir, fmt.Sprintf("judge_%s_%d", submissionID, time.Now().UnixNano()))
	return tempDir, os.MkdirAll(tempDir, 0755)
}

// compileVerilog 编译Verilog代码和testbench
func (j *Judge) compileVerilog(tempDir, designCode, testbenchCode string) error {
	// 写入设计文件
	designFile := filepath.Join(tempDir, "design.v")
	if err := os.WriteFile(designFile, []byte(designCode), 0644); err != nil {
		return fmt.Errorf("failed to write design file: %v", err)
	}

	// 写入testbench文件
	testbenchFile := filepath.Join(tempDir, "testbench.v")
	if err := os.WriteFile(testbenchFile, []byte(testbenchCode), 0644); err != nil {
		return fmt.Errorf("failed to write testbench file: %v", err)
	}

	// 使用iverilog编译设计和testbench
	cmd := exec.Command("iverilog", "-o", filepath.Join(tempDir, "simulation"), designFile, testbenchFile)
	cmd.Dir = tempDir
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output))
	}

	return nil
}

// runSingleTest 运行单个Verilog测试用例
func (j *Judge) runSingleTest(ctx context.Context, tempDir string, testCase TestCase, designCode string, timeLimit, memoryLimit int) (*JudgeResult, error) {
	result := &JudgeResult{}
	
	// 为每个测试用例重新编译（因为testbench可能不同）
	if err := j.compileVerilog(tempDir, designCode, testCase.Testbench); err != nil {
		result.Status = "compile_error"
		result.ErrorMessage = err.Error()
		return result, nil
	}
	
	// 执行仿真
	executable := filepath.Join(tempDir, "simulation")
	vcdFile := filepath.Join(tempDir, "output.vcd")
	
	cmd := exec.Command("vvp", executable)
	cmd.Dir = tempDir
	
	// 设置超时
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeLimit)*time.Millisecond)
	defer cancel()
	
	startTime := time.Now()
	output, err := cmd.CombinedOutput()
	runTime := int(time.Since(startTime).Milliseconds())
	
	result.RunTime = runTime
	result.Memory = 1024 // 简化处理，实际应该获取真实内存使用
	
	// 检查超时
	if timeoutCtx.Err() == context.DeadlineExceeded {
		result.Status = "time_limit_exceeded"
		return result, nil
	}
	
	if err != nil {
		result.Status = "runtime_error"
		result.ErrorMessage = fmt.Sprintf("Simulation failed: %s", string(output))
		return result, nil
	}

	// 检查VCD文件是否生成
	if _, err := os.Stat(vcdFile); os.IsNotExist(err) {
		result.Status = "runtime_error"
		result.ErrorMessage = "VCD file not generated"
		return result, nil
	}

	// 比较VCD输出
	if j.compareVCD(vcdFile, testCase.ExpectedVCD) {
		result.Status = "accepted"
	} else {
		result.Status = "wrong_answer"
		result.ErrorMessage = "VCD output does not match expected results"
	}
	
	return result, nil
}

// compareVCD 比较VCD文件输出
func (j *Judge) compareVCD(actualVCDFile, expectedPattern string) bool {
	// 读取生成的VCD文件
	vcdContent, err := os.ReadFile(actualVCDFile)
	if err != nil {
		return false
	}

	// 解析期望的模式（可以是具体的信号值序列或正则表达式）
	return j.matchVCDPattern(string(vcdContent), expectedPattern)
}

// matchVCDPattern 匹配VCD模式
func (j *Judge) matchVCDPattern(vcdContent, pattern string) bool {
	// 如果pattern是JSON格式的信号值检查
	if strings.HasPrefix(pattern, "{") {
		return j.checkSignalValues(vcdContent, pattern)
	}
	
	// 否则作为正则表达式处理
	matched, err := regexp.MatchString(pattern, vcdContent)
	return err == nil && matched
}

// checkSignalValues 检查特定信号的值
func (j *Judge) checkSignalValues(vcdContent, signalPattern string) bool {
	// 简化实现：检查关键信号在特定时间点的值
	// 实际实现中可以解析JSON格式的期望值并与VCD中的信号值比较
	
	// 示例：检查是否包含特定的信号变化模式
	lines := strings.Split(vcdContent, "\n")
	for _, line := range lines {
		// 解析VCD格式的信号变化
		if strings.HasPrefix(line, "#") {
			// 时间戳行
			continue
		}
		if len(line) > 0 && (line[0] == '0' || line[0] == '1' || line[0] == 'x' || line[0] == 'z') {
			// 信号值变化行
			// 这里可以实现具体的信号值检查逻辑
		}
	}
	
	// 简化处理：如果VCD文件包含期望的模式字符串则认为匹配
	return strings.Contains(vcdContent, signalPattern)
}

// parseVCDSignals 解析VCD文件中的信号值（辅助函数）
func (j *Judge) parseVCDSignals(vcdContent string) map[string][]string {
	signals := make(map[string][]string)
	lines := strings.Split(vcdContent, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		
		// 解析信号值变化
		if len(line) > 1 && (line[0] == '0' || line[0] == '1' || line[0] == 'x' || line[0] == 'z') {
			value := string(line[0])
			signalID := line[1:]
			signals[signalID] = append(signals[signalID], value)
		}
	}
	
	return signals
}