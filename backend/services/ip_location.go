package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// IPLocation IP归属地信息
type IPLocation struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	Region      string `json:"regionName"`
	City        string `json:"city"`
	ISP         string `json:"isp"`
	Query       string `json:"query"`
}

// ipLocationCache IP归属地缓存
var ipLocationCache = struct {
	sync.RWMutex
	data map[string]*IPLocation
}{data: make(map[string]*IPLocation)}

// httpClient 复用HTTP客户端
var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

// GetIPLocation 查询IP归属地
// 使用 ip-api.com 免费API (每分钟45次请求限制)
func GetIPLocation(ip string) *IPLocation {
	// 检查缓存
	ipLocationCache.RLock()
	if loc, ok := ipLocationCache.data[ip]; ok {
		ipLocationCache.RUnlock()
		return loc
	}
	ipLocationCache.RUnlock()

	// 过滤内网IP
	if isPrivateIP(ip) {
		loc := &IPLocation{
			Country: "内网",
			Region:  "本地",
			City:    "本地",
			Query:   ip,
		}
		cacheIPLocation(ip, loc)
		return loc
	}

	// 调用API查询
	url := fmt.Sprintf("http://ip-api.com/json/%s?lang=zh-CN", ip)
	resp, err := httpClient.Get(url)
	if err != nil {
		return &IPLocation{Country: "未知", Query: ip}
	}
	defer resp.Body.Close()

	var result struct {
		Status      string `json:"status"`
		Country     string `json:"country"`
		CountryCode string `json:"countryCode"`
		RegionName  string `json:"regionName"`
		City        string `json:"city"`
		ISP         string `json:"isp"`
		Query       string `json:"query"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return &IPLocation{Country: "未知", Query: ip}
	}

	if result.Status != "success" {
		return &IPLocation{Country: "未知", Query: ip}
	}

	loc := &IPLocation{
		Country:     result.Country,
		CountryCode: result.CountryCode,
		Region:      result.RegionName,
		City:        result.City,
		ISP:         result.ISP,
		Query:       result.Query,
	}

	// 缓存结果
	cacheIPLocation(ip, loc)

	return loc
}

// cacheIPLocation 缓存IP归属地
func cacheIPLocation(ip string, loc *IPLocation) {
	ipLocationCache.Lock()
	defer ipLocationCache.Unlock()
	ipLocationCache.data[ip] = loc
}

// isPrivateIP 判断是否为内网IP
func isPrivateIP(ip string) bool {
	// 简单判断内网IP
	if ip == "127.0.0.1" || ip == "::1" || ip == "localhost" {
		return true
	}
	// 10.x.x.x
	if len(ip) >= 3 && ip[:3] == "10." {
		return true
	}
	// 172.16-31.x.x
	if len(ip) >= 7 && ip[:4] == "172." {
		return true
	}
	// 192.168.x.x
	if len(ip) >= 8 && ip[:8] == "192.168." {
		return true
	}
	return false
}

// BatchGetIPLocations 批量查询IP归属地
func BatchGetIPLocations(ips []string) map[string]*IPLocation {
	result := make(map[string]*IPLocation)
	for _, ip := range ips {
		result[ip] = GetIPLocation(ip)
		// 简单限流，避免API限制
		time.Sleep(25 * time.Millisecond)
	}
	return result
}
