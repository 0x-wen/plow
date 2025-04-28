package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		IP      string `yaml:"ip"`
		Port    string `yaml:"port"`
		BaseUrl string `yaml:"base_url"`
	} `yaml:"server"`
	Endpoints []struct {
		Name    string `yaml:"name"`
		Path    string `yaml:"path"`
		Method  string `yaml:"method"`
		Payload string `yaml:"payload"`
	} `yaml:"endpoints"`
	Benchmark struct {
		Concurrency int    `yaml:"concurrency"`
		Duration    string `yaml:"duration"`
		Requests    int    `yaml:"requests"`
		Listen      string `yaml:"listen"`
		OutputDir   string `yaml:"output_dir"`
	} `yaml:"benchmark"`
}

func main() {
	// 读取配置文件
	configFile := "config.yaml"
	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		return
	}

	// 创建输出目录
	err = os.MkdirAll(config.Benchmark.OutputDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	// 遍历接口并压测
	for _, endpoint := range config.Endpoints {
		outputFile := filepath.Join(config.Benchmark.OutputDir, endpoint.Name+".json")
		// 构造 URL
		var url string
		if config.Server.BaseUrl != "" {
			url = fmt.Sprintf("%s%s", config.Server.BaseUrl, endpoint.Path)
		} else {
			url = fmt.Sprintf("http://%s:%s%s", config.Server.IP, config.Server.Port, endpoint.Path)
		}

		fmt.Printf("Testing %s...\n", endpoint.Path)

		// 构造 plow 命令
		cmd := exec.Command("sh", "-c", fmt.Sprintf(
			"plow %s -b '%s' -m %s -c %d -d %s -n %d --json --summary",
			url,
			endpoint.Payload,
			endpoint.Method,
			config.Benchmark.Concurrency,
			config.Benchmark.Duration,
			config.Benchmark.Requests,
		))

		fmt.Println("Command:", cmd.String())

		// 保存结果到文件
		output, err := os.Create(outputFile)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			continue
		}
		defer output.Close()

		cmd.Stdout = output
		cmd.Stderr = os.Stderr

		// 执行命令
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error running plow for %s: %v\n", endpoint.Path, err)
			continue
		}
		time.Sleep(time.Second * 5)
		fmt.Printf("Results saved to %s\n", outputFile)
	}
}
