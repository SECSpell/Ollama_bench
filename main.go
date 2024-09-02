package main

import (
	"bytes"      // 提供了操作字节切片的函数 // Provides functions for manipulating byte slices
	"context"    // 提供了跨 API 边界的取消、超时和截止日期功能 // Provides cancellation, timeout, and deadline functionality across API boundaries
	"encoding/json"  // 实现了 JSON 的编解码 // Implements encoding and decoding of JSON
	"fmt"        // 实现格式化 I/O // Implements formatted I/O
	"io"         // 提供了 I/O 原语的基本接口 // Provides basic interfaces to I/O primitives
	"log"        // 实现了简单的日志功能 // Implements a simple logging package
	"math/rand"  // 实现伪随机数生成器 // Implements pseudo-random number generators
	"net/http"   // 提供 HTTP 客户端和服务器实现 // Provides HTTP client and server implementations
	"os"         // 提供了操作系统函数的平台无关接口 // Provides a platform-independent interface to operating system functionality
	"path/filepath"  // 实现了兼容各操作系统的文件路径操作 // Implements utility routines for manipulating filename paths in a way compatible with all supported operating systems
	"strconv"    // 实现了基本数据类型和其字符串表示的相互转换 // Implements conversions to and from string representations of basic data types
	"sync"       // 提供了基本的同步原语 // Provides basic synchronization primitives
	"time"       // 提供了时间的测量和显示功能 // Provides functionality for measuring and displaying time

	"golang.org/x/sync/semaphore"  // 提供了带权重的信号量实现 // Provides a weighted semaphore implementation
)

var questions = []string{
	"Why is the sky blue?", "Why do we dream?", "Why is the ocean salty?", "Why do leaves change color?",
	"Why do birds sing?", "Why do we have seasons?", "Why do stars twinkle?", "Why do we yawn?",
	"Why is the sun hot?", "Why do cats purr?", "Why do dogs bark?", "Why do fish swim?",
	"Why do we have fingerprints?", "Why do we sneeze?", "Why do we have eyebrows?", "Why do we have hair?",
	"Why do we have nails?", "Why do we have teeth?", "Why do we have bones?", "Why do we have muscles?",
	"Why do we have blood?", "Why do we have a heart?", "Why do we have lungs?", "Why do we have a brain?",
	"Why do we have skin?", "Why do we have ears?", "Why do we have eyes?", "Why do we have a nose?",
	"Why do we have a mouth?", "Why do we have a tongue?", "Why do we have a stomach?", "Why do we have intestines?",
	"Why do we have a liver?", "Why do we have kidneys?", "Why do we have a bladder?", "Why do we have a pancreas?",
	"Why do we have a spleen?", "Why do we have a gallbladder?", "Why do we have a thyroid?", "Why do we have adrenal glands?",
	"Why do we have a pituitary gland?", "Why do we have a hypothalamus?", "Why do we have a thymus?", "Why do we have lymph nodes?",
	"Why do we have a spinal cord?", "Why do we have nerves?", "Why do we have a circulatory system?", "Why do we have a respiratory system?",
	"Why do we have a digestive system?", "Why do we have an immune system?",
}

// Config 结构体定义了配置文件的结构
// Config struct defines the structure of the configuration file
type Config struct {
	BaseURL   string `json:"base_url"`
	APIKey    string `json:"apikey"`
	ModelName string `json:"model_name"`
}

type RequestPayload struct {
	Model     string `json:"model"`
	Messages  []Message `json:"messages"`
	Stream    bool   `json:"stream"`
	Temperature float64 `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
	Error string `json:"error"`
}

// loadConfig 加载配置文件
// loadConfig loads the configuration file
func loadConfig() (*Config, error) {
	// 获取可执行文件的路径
	// Get the path of the executable file
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("unable to get executable path: %v", err)
	}
	
	// 获取可执行文件所在的目录
	// Get the directory of the executable file
	exeDir := filepath.Dir(exePath)
	
	// 构建配置文件的完整路径
	// Build the full path of the configuration file
	configPath := filepath.Join(exeDir, "config.json")
	
	// 读取配置文件
	// Read the configuration file
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果文件不存在，则创建一个新的配置文件
			// If the file doesn't exist, create a new configuration file
			config := &Config{
				BaseURL:   "http://localhost:11434/v1/chat/completions",
				APIKey:    "your_api_key_here",
				ModelName: "llama3.1:latest",
			}
			data, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return nil, fmt.Errorf("unable to create default configuration: %v", err)
			}
			if err := os.WriteFile(configPath, data, 0644); err != nil {
				return nil, fmt.Errorf("unable to write configuration file %s: %v", configPath, err)
			}
			fmt.Printf("New configuration file created at %s. Please edit and fill in the correct API key.\n", configPath)
			return config, nil
		}
		return nil, fmt.Errorf("unable to read configuration file %s: %v", configPath, err)
	}

	// 解析配置文件
	// Parse the configuration file
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unable to parse configuration file %s: %v", configPath, err)
	}

	return &config, nil
}

// fetch 发送 HTTP 请求并获取响应
// fetch sends an HTTP request and retrieves the response
func fetch(ctx context.Context, client *http.Client, config *Config, question string) (int, float64, error) {
	startTime := time.Now()

	payload := RequestPayload{
		Model: config.ModelName,
		Messages: []Message{
			{Role: "user", Content: question},
		},
		Stream: false,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", config.BaseURL, io.NopCloser(bytes.NewReader(jsonData)))
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, 0, err
	}

	if response.Error != "" {
		return 0, 0, fmt.Errorf("API error: %s", response.Error)
	}

	totalTokens := response.Usage.TotalTokens
	elapsedTime := time.Since(startTime).Seconds()

	return totalTokens, elapsedTime, nil
}

func main() {
	// 打印当前工作目录和可执行文件路径，用于调试
	// Print current working directory and executable path for debugging
	// currentDir, _ := os.Getwd()
	// exePath, _ := os.Executable()
	// fmt.Printf("Current working directory: %s\n", currentDir)
	// fmt.Printf("Executable path: %s\n", exePath)

	var maxConcurrentRequests, totalRequests int
	var err error

	switch len(os.Args) {
	case 1:
		maxConcurrentRequests = 1
		totalRequests = 4
	case 3:
		maxConcurrentRequests, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalf("Invalid max concurrent requests: %v", err)
		}
		totalRequests, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid total requests: %v", err)
		}
	default:
		log.Fatalf("Usage: %s <C> <N> (or no arguments for default values)", os.Args[0])
	}

	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	client := &http.Client{}
	sem := semaphore.NewWeighted(int64(maxConcurrentRequests))
	var wg sync.WaitGroup

	totalTokens := 0
	totalTime := 0.0

	// 添加英文提示信息
	// Add English prompt message
	fmt.Println("Starting request test, please wait...")

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := sem.Acquire(context.Background(), 1); err != nil {
				log.Printf("Failed to acquire semaphore: %v", err)
				return
			}
			defer sem.Release(1)

			question := questions[rand.Intn(len(questions))]
			tokens, elapsedTime, err := fetch(context.Background(), client, config, question)
			if err != nil {
				log.Printf("Request failed: %v", err)
				return
			}

			totalTokens += tokens
			totalTime += elapsedTime
		}()
	}

	wg.Wait()

	avgTimePerRequest := totalTime / float64(totalRequests)
	tokensPerSecond := float64(totalTokens) / totalTime

	fmt.Printf("Performance Results:\n")
	fmt.Printf("  Total requests            : %d\n", totalRequests)
	fmt.Printf("  Max concurrent requests   : %d\n", maxConcurrentRequests)
	fmt.Printf("  Total tokens              : %d\n", totalTokens)
	fmt.Printf("  Total time                : %.2f seconds\n", totalTime)
	fmt.Printf("  Average time per request  : %.2f seconds\n", avgTimePerRequest)
	fmt.Printf("  Tokens per second         : %.2f\n", tokensPerSecond)

	fmt.Println("\nPress Enter to exit...")
	fmt.Scanln()
}
