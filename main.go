package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"ghproxy/config"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req/v3"
)

var (
	exps = []*regexp.Regexp{
		regexp.MustCompile(`^(?:https?://)?github\.com/([^/]+)/([^/]+)/(?:releases|archive)/.*$`),
		regexp.MustCompile(`^(?:https?://)?github\.com/([^/]+)/([^/]+)/(?:blob|raw)/.*$`),
		regexp.MustCompile(`^(?:https?://)?github\.com/([^/]+)/([^/]+)/(?:info|git-).*$`),
		regexp.MustCompile(`^(?:https?://)?raw\.github(?:usercontent|)\.com/([^/]+)/([^/]+)/.+?/.+$`),
		regexp.MustCompile(`^(?:https?://)?gist\.github\.com/([^/]+)/.+?/.+$`),
	}
)

var (
	router *gin.Engine
	cfg    *config.Config
)

func init() {
	// 加载配置
	var err error
	cfg, err = config.LoadConfig("/data/ghproxy/config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("Loaded config: %v\n", cfg)

	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化路由
	router = gin.Default()

	// 定义路由
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://ghproxy0rtt.1888866.xyz/")
	})

	// 健康检查
	router.GET("/api/healthcheck", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// 未匹配路由处理
	router.NoRoute(noRouteHandler(cfg))
}

func main() {
	//初始化日志
	logFile, err := os.OpenFile(cfg.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Log Initialization Failed: > %s", err)
		fmt.Printf("Failed to open log file: %s", cfg.LogFilePath)
		fmt.Printf("Please check the log file path and permissions.")
	} else {
		defer logFile.Close()
		log.SetOutput(logFile)
		log.Println("Log Initialization Complete")
	}
	// 启动服务器
	err = router.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}

	fmt.Println("Program finished")
}

//reserved for future use
/*func handler(c *gin.Context, config *config.Config) {
	rawPath := strings.TrimPrefix(c.Request.URL.RequestURI(), "/")
	re := regexp.MustCompile(`^(http:|https:)?/?/?(.*)`)
	matches := re.FindStringSubmatch(rawPath)

	rawPath = "https://" + matches[2]

	matches = checkURL(rawPath)
	if matches == nil {
		c.String(http.StatusForbidden, "Invalid input.")
		return
	}

	if exps[1].MatchString(rawPath) {
		rawPath = strings.Replace(rawPath, "/blob/", "/raw/", 1)
	}

	proxy(c, rawPath, config)
}*/

func noRouteHandler(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawPath := strings.TrimPrefix(c.Request.URL.RequestURI(), "/")
		re := regexp.MustCompile(`^(http:|https:)?/?/?(.*)`)
		matches := re.FindStringSubmatch(rawPath)

		rawPath = "https://" + matches[2]

		matches = checkURL(rawPath)
		if matches == nil {
			c.String(http.StatusForbidden, "Invalid input.")
			return
		}

		if exps[1].MatchString(rawPath) {
			rawPath = strings.Replace(rawPath, "/blob/", "/raw/", 1)
		}

		//日志记录
		log.Printf("Request: %s %s", c.Request.Method, rawPath)
		log.Printf("Matches: %v", matches)

		//proxyGit(c, rawPath, config)
		switch {
		case exps[0].MatchString(rawPath):
			log.Printf("%s Matched EXPS[0] - USE proxy-chrome", rawPath)
			proxychrome(c, rawPath, config)
		case exps[1].MatchString(rawPath):
			log.Printf("%s Matched EXPS[1] - USE proxy-chrome", rawPath)
			proxychrome(c, rawPath, config)
		case exps[2].MatchString(rawPath):
			log.Printf("%s Matched EXPS[2] - USE proxy-git", rawPath)
			proxyGit(c, rawPath, config)
		case exps[3].MatchString(rawPath):
			log.Printf("%s Matched EXPS[3] - USE proxy-chrome", rawPath)
			proxychrome(c, rawPath, config)
		case exps[4].MatchString(rawPath):
			log.Printf("%s Matched EXPS[4] - USE proxy-chrome", rawPath)
			proxychrome(c, rawPath, config)
		default:
			c.String(http.StatusForbidden, "Invalid input.")
			return
		}
	}
}

/*func proxy(c *gin.Context, u string, config *config.Config) {
	req, err := http.NewRequest(c.Request.Method, u, c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("server error %v", err))
		log.Printf("Failed to create request: %v", err)
		return
	}
	defer c.Request.Body.Close()

	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	req.Header.Del("Host")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")

	//resp, err := http.DefaultClient.Do(req)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("server error %v", err))
		log.Printf("Failed to send request: %v", err)
		return
	}
	defer resp.Body.Close()

	contentLength := resp.Header.Get("Content-Length")
	if contentLength != "" {
		size, err := strconv.Atoi(contentLength)
		if err == nil && size > config.SizeLimit {
			finalURL := resp.Request.URL.String()
			c.Redirect(http.StatusMovedPermanently, finalURL)
			log.Printf("%s - Redirecting to %s due to size limit (%d bytes)", time.Now().Format("2006-01-02 15:04:05"), finalURL, size)
			return
		}
	}

	resp.Header.Del("Content-Security-Policy")
	resp.Header.Del("Referrer-Policy")
	resp.Header.Del("Strict-Transport-Security")

	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	if config.CORSOrigin {
		resp.Header.Set("Access-Control-Allow-Origin", "*")
	} else {
		resp.Header.Del("Access-Control-Allow-Origin")
	}

	c.Status(resp.StatusCode)
	if _, err := io.Copy(c.Writer, resp.Body); err != nil {
		return
	}
}*/

// 使用req库伪装git
func proxyGit(c *gin.Context, u string, config *config.Config) {
	method := c.Request.Method
	log.Printf("%s Method: %s", u, method)
	client := req.C().SetUserAgent("git/2.33.1")

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("server error %v", err))
		log.Printf("Failed to read request body: %v", err)
		return
	}
	defer c.Request.Body.Close()

	// 创建新的请求
	req := client.R().SetBody(body)

	// 复制请求头
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.SetHeader(key, value)
		}
	}

	// 发送请求并处理响应
	resp, err := sendRequest(req, method, u)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("server error %v", err))
		log.Printf("Failed to send request: %v", err)
		return
	}
	defer resp.Body.Close()

	// 检查响应内容长度并处理重定向
	if err := handleResponseSize(resp, config, c); err != nil {
		log.Printf("Error handling response size: %v", err)
		return
	}

	// 删除不必要的响应头
	resp.Header.Del("Content-Security-Policy")
	resp.Header.Del("Referrer-Policy")
	resp.Header.Del("Strict-Transport-Security")

	// 复制响应头
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	if config.CORSOrigin {
		c.Header("Access-Control-Allow-Origin", "*")
	} else {
		c.Header("Access-Control-Allow-Origin", "")
	}

	c.Status(resp.StatusCode)
	if _, err := io.Copy(c.Writer, resp.Body); err != nil {
		log.Printf("Failed to copy response body: %v", err)
		return
	}
}

// 发送请求并返回响应
func sendRequest(req *req.Request, method, url string) (*req.Response, error) {
	switch method {
	case "GET":
		return req.Get(url)
	case "POST":
		return req.Post(url)
	case "PUT":
		return req.Put(url)
	case "DELETE":
		return req.Delete(url)
	default:
		return nil, fmt.Errorf("unsupported method: %s", method)
	}
}

// 处理响应内容长度
func handleResponseSize(resp *req.Response, config *config.Config, c *gin.Context) error {
	contentLength := resp.Header.Get("Content-Length")
	if contentLength != "" {
		size, err := strconv.Atoi(contentLength)
		if err == nil && size > config.SizeLimit {
			finalURL := resp.Request.URL.String()
			c.Redirect(http.StatusMovedPermanently, finalURL)
			log.Printf("%s - Redirecting to %s due to size limit (%d bytes)", time.Now().Format("2006-01-02 15:04:05"), finalURL, size)
			return fmt.Errorf("response size exceeds limit")
		}
	}
	return nil
}

// 使用req库伪装chrome浏览器
func proxychrome(c *gin.Context, u string, config *config.Config) {
	method := c.Request.Method
	log.Printf("%s Method: %s", u, method)
	client := req.C().
		SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36").
		SetTLSFingerprintChrome()

	client.ImpersonateChrome()

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("server error %v", err))
		log.Printf("Failed to read request body: %v", err)
		return
	}
	defer c.Request.Body.Close()

	// 创建新的请求
	req := client.R().
		SetBody(body).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")

	// 复制请求头
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.SetHeader(key, value)
		}
	}

	// 发送请求并处理响应
	resp, err := sendRequest(req, method, u)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("server error %v", err))
		log.Printf("Failed to send request: %v", err)
		return
	}
	defer resp.Body.Close()

	// 检查响应内容长度并处理重定向
	if err := handleResponseSize(resp, config, c); err != nil {
		log.Printf("Error handling response size: %v", err)
		return
	}

	// 删除不必要的响应头
	resp.Header.Del("Content-Security-Policy")
	resp.Header.Del("Referrer-Policy")
	resp.Header.Del("Strict-Transport-Security")

	// 复制响应头
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	if config.CORSOrigin {
		c.Header("Access-Control-Allow-Origin", "*")
	} else {
		c.Header("Access-Control-Allow-Origin", "")
	}

	c.Status(resp.StatusCode)
	if _, err := io.Copy(c.Writer, resp.Body); err != nil {
		log.Printf("Failed to copy response body: %v", err)
		return
	}
}

func checkURL(u string) []string {
	for _, exp := range exps {
		if matches := exp.FindStringSubmatch(u); matches != nil {
			log.Printf("URL matched: %s, Matches: %v", u, matches[1:])
			return matches[1:]
		}
	}
	log.Printf("Invalid URL: %s", u)
	return nil
}
