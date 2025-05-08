package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
	"zhaoxin2025/config"
)

// WxLogin 发送请求到微信API并返回openid和session_key和error
func WxLogin(code string) (string, string, error) {
	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 基础 URL
	baseURL := "https://api.weixin.qq.com/sns/jscode2session"

	// 创建查询参数
	params := url.Values{}
	params.Add("appid", config.Config.AppID)
	params.Add("secret", config.Config.AppSecret)
	params.Add("js_code", code)
	params.Add("grant_type", "authorization_code")

	// 将参数添加到 URL
	reqURL := baseURL + "?" + params.Encode()

	// 创建请求
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("HTTP错误: %d", resp.StatusCode)
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// 解析 JSON 响应
	var wxResp struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		UnionID    string `json:"unionid,omitempty"`
		ErrCode    int    `json:"errcode,omitempty"`
		ErrMsg     string `json:"errmsg,omitempty"`
	}
	if err := json.Unmarshal(body, &wxResp); err != nil {
		return "", "", fmt.Errorf("解析微信响应失败: %w", err)
	}

	// 检查响应是否包含错误
	if wxResp.ErrCode != 0 {
		return "", "", fmt.Errorf("微信API错误: %d - %s", wxResp.ErrCode, wxResp.ErrMsg)
	}

	// 返回 openid 和 session_key
	return wxResp.OpenID, wxResp.SessionKey, nil
}
