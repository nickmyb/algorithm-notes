package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mozillazg/request"
	"github.com/nickmyb/algorithm-notes/ctl/util"
)

// 所有接口地址都从 util.Site 拼出来，换站点只改那一个常量。
var (
	// AllProblemURL 返回全部题目的题号、标题、slug、难度、通过率，不含题目描述。
	AllProblemURL = util.Site + "/api/problems/all/"
	// QraphqlURL 是 GraphQL 入口，题目描述从这里取。
	QraphqlURL = util.Site + "/graphql/"
	// LoginPageURL 只用来当请求头里的 Referer。
	LoginPageURL = util.Site + "/accounts/login/"
)

var req *request.Request

// sentCookie 记录本次进程有没有真的带着 Cookie 去请求。
//
// 它和「接口返回的 user_name 非空」是两个不同的信号：前者是本地事实（我们发没发凭据），
// 后者是服务器的判断（认不认这个凭据）。两个合起来才能分辨出「Cookie 过期」——
// 那种情况下接口不报错，只是返回未登录数据，光看任何一个信号都会误判。
var sentCookie bool

func newReq() *request.Request {
	if req == nil {
		req = signin()
	}
	return req
}

func signin() *request.Request {
	cfg := getConfig()
	sentCookie = cfg.Cookie != ""
	req := request.NewRequest(&http.Client{Timeout: 30 * time.Second})
	req.Headers = map[string]string{
		"Content-Type":    "application/json",
		"Accept-Encoding": "",
		"cookie":          cfg.Cookie,
		"x-csrftoken":     cfg.CSRFtoken,
		"Referer":         LoginPageURL,
		"origin":          util.Site,
	}
	return req
}

func getRaw(URL string) []byte {
	req := newReq()
	resp, err := req.Get(URL)
	if err != nil || resp == nil {
		fmt.Printf("getRaw: Get Error: %v\n", err)
		return []byte{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("getRaw: Read Error: %s\n", err.Error())
		return []byte{}
	}
	// 非 200（如被限流的 429）直接返回空，让调用方做返回值校验后跳过，避免拿错误页继续解析
	if resp.StatusCode != 200 {
		fmt.Printf("getRaw: non-200 status %d for %s\n", resp.StatusCode, URL)
		return []byte{}
	}
	return body
}

func getProblemAllList() []byte {
	return getRaw(AllProblemURL)
}

func getQraphql(payload string) []byte {
	req := newReq()
	resp, err := req.PostForm(QraphqlURL, bytes.NewBuffer([]byte(payload)))
	if err != nil || resp == nil {
		fmt.Printf("getQraphql: Post Error: %v\n", err)
		return []byte{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("getQraphql: Read Error: %s\n", err.Error())
		return []byte{}
	}
	// 非 200（如被限流的 429）直接返回空，让调用方做返回值校验后跳过
	if resp.StatusCode != 200 {
		fmt.Printf("getQraphql: non-200 status %d\n", resp.StatusCode)
		return []byte{}
	}
	return body
}

// questionDetail 是一道题的题目描述。
//
// leetcode.cn 的 GraphQL 同时返回英文原文和官方中文翻译，而且不需要登录：
// Content 是英文原文，TranslatedContent 是官方中文翻译，两者都是 HTML。
// leetcode.com 只有 Content，TranslatedContent 会是空串，代码里按空处理即可。
type questionDetail struct {
	QuestionFrontendID string `json:"questionFrontendId"`
	Title              string `json:"title"`
	TranslatedTitle    string `json:"translatedTitle"`
	Difficulty         string `json:"difficulty"`
	Content            string `json:"content"`
	TranslatedContent  string `json:"translatedContent"`
}

const questionDetailQuery = `query questionData($titleSlug: String!) {
  question(titleSlug: $titleSlug) {
    questionFrontendId
    title
    translatedTitle
    difficulty
    content
    translatedContent
  }
}`

// getQuestionDetail 按 slug 取题目描述。取不到时返回 nil，调用方把 README 留空即可，不算致命错误。
func getQuestionDetail(slug string) *questionDetail {
	payload, err := json.Marshal(map[string]any{
		"operationName": "questionData",
		"query":         questionDetailQuery,
		"variables":     map[string]string{"titleSlug": slug},
	})
	if err != nil {
		fmt.Printf("拼 GraphQL 请求失败: %v\n", err)
		return nil
	}

	body := getQraphql(string(payload))
	if len(body) == 0 {
		return nil
	}

	var resp struct {
		Data struct {
			Question *questionDetail `json:"question"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Printf("解析题目描述失败: %v\n", err)
		return nil
	}
	return resp.Data.Question
}
