package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
	"zhaoxin2025/config"
	"zhaoxin2025/model"
)

type AccessTokenInfo struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	FetchTime   time.Time
}

func (token *AccessTokenInfo) IsExpired() bool {
	return time.Since(token.FetchTime) > time.Duration(token.ExpiresIn)*time.Second
}

type WxMessage struct {
	AccessToken      string `json:"access_token"`
	TemplateID       string `json:"template_id"`
	Page             string `json:"page"`
	OpenID           string `json:"touser"`
	MiniProgramState string `json:"miniprogram_state"`
	Lang             string `json:"lang"`
	Data             any
}

type RegisterMessage struct {
	Name   Field `json:"thing8"`
	ReType Field `json:"thing18"`
	Interv Field `json:"phrase5"`
	Time   Field `json:"time4"`
	Note   Field `json:"thing22"`
}

type IntervTimeMessage struct {
	Time Field `json:"time4"`
	Site Field `json:"thing7"`
	Note Field `json:"thing8"`
}

type IntervResMessage struct {
	Phrase Field `json:"phrase01"`
	Thing  Field `json:"thing01"`
}

type Field struct {
	Value string `json:"value"`
}

type WxMessageResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

var AccessToken AccessTokenInfo

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

func SendMessage(openid string, template_id string, data any) error {
	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	if AccessToken.IsExpired() {
		if err := GetAccessToken(); err != nil {
			return fmt.Errorf("获取AccessToken失败: %w", err)
		}
	}
	
	// 基础 URL
	baseURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%s", AccessToken.AccessToken)

	// 创建请求体
	message := WxMessage{
		TemplateID:       template_id,
		Page:             "",
		OpenID:           openid,
		MiniProgramState: "developer",
		Lang:             "zh_CN",
		Data:             data,
	}

	// 将请求体转换为 JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("JSON编码失败: %w", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP错误: %d", resp.StatusCode)
	}

	var wxResp WxMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&wxResp); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if wxResp.ErrCode != 0 {
		return fmt.Errorf("微信API错误: %d - %s", wxResp.ErrCode, wxResp.ErrMsg)
	}

	return nil
}

func SendRegister(record model.Stu) error {
	if err := SendMessage(record.OpenID, config.Config.TemplateID[0], RegisterMessage{
		Name:   Field{record.Name},
		ReType: Field{string(record.Depart)},
		Interv: Field{Pass(record.Interv.Pass)},
		Time:   Field{Time(record.Interv.Time)},
		Note:   Field{"欢迎报名挑战网招新面试，期待你的如约而至！"},
	}); err != nil {
		return fmt.Errorf("发送订阅消息失败: %w", err)
	}
	return nil
}

func SendResult(stu model.Stu) error {
	if err := SendMessage(stu.OpenID, config.Config.TemplateID[1], IntervResMessage{
		Phrase: Field{Value: Pass(stu.Interv.Pass)},
		Thing:  Field{Value: "点击即可查看详细信息，挑战，无处不在！"},
	}); err != nil {
		return fmt.Errorf("发送订阅消息失败: %w", err)
	}
	return nil
}

func SendTime(stu model.Stu) error {
	if err := SendMessage(stu.OpenID, config.Config.TemplateID[2], IntervTimeMessage{
		Time: Field{Value: stu.Interv.Time.Format("2006年01月02日 15:04:05")},
		Site: Field{Value: "挑战阁楼"},
		Note: Field{Value: "请您按时参加面试，特殊情况请联系管理员"},
	}); err != nil {
		return fmt.Errorf("发送订阅消息失败: %w", err)
	}
	return nil
}

func GetAccessToken() error {
	// 设置请求体
	URL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		config.Config.AppID,
		config.Config.AppSecret)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}
	// 发请求
	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println("GetAccessToken:", err)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	// 解析数据
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var info AccessTokenInfo
	if err := json.Unmarshal(respBody, &info); err != nil {
		return err
	}
	AccessToken = info
	AccessToken.FetchTime = time.Now()
	return nil
}
