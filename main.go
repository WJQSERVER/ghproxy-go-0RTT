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

func main() {
	//加载配置
	config, err := config.LoadConfig("/data/ghproxy/config/config.yaml")
	if err != nil {
		fmt.Printf("Failed to load config: %v", err)
		return
	}
	fmt.Printf("Loaded config: %v", config)

	//初始化日志
	logFile, err := os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Log Initialization Failed: > %s", err)
		fmt.Printf("Failed to open log file: %s", config.LogFilePath)
		fmt.Printf("Please check the log file path and permissions.")
	} else {
		defer logFile.Close()
		log.SetOutput(logFile)
		log.Println("Log Initialization Complete")
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "https://ghproxy0rtt.1888866.xyz/")
	})

	router.NoRoute(noRouteHandler(config))

	err = router.Run(fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
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

		proxy(c, rawPath, config)
	}
}

func proxy(c *gin.Context, u string, config *config.Config) {
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
	//req.Header.Del("Host")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	//resp, err := http.DefaultClient.Do(req)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	//resp, err = client.Do(req)
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
			c.Redirect(http.StatusFound, finalURL)
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

	c.Status(resp.StatusCode)
	if _, err := io.Copy(c.Writer, resp.Body); err != nil {
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
