package tool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// SearchTool is a tool for searching the web.
type SearchTool struct{}

func (t *SearchTool) Name() string {
	return "search"
}

func (t *SearchTool) Description() string {
	return "A tool for searching the web using DuckDuckGo."
}

func (t *SearchTool) Execute(args json.RawMessage) (string, error) {
	var params struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	return t.search(params.Query)
}

func (t *SearchTool) search(query string) (string, error) {
	encodedQuery := url.QueryEscape(query)

	// 使用DuckDuckGo HTML搜索端点
	url := "https://html.duckduckgo.com/html/?q=" + encodedQuery

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// 设置User-Agent避免被阻止
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("search failed with status code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	os.WriteFile("search.html", body, 0644)
	return t.parseHTML(body)
}

func (t *SearchTool) parseHTML(body []byte) (string, error) {
	content := string(body)

	// 简单的HTML解析，提取搜索结果
	var results []string

	// 查找搜索结果标题和描述
	// DuckDuckGo HTML结构中的搜索结果
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 查找标题
		if strings.Contains(line, "result__title") || strings.Contains(line, "result__a") {
			// 提取标题文本
			start := strings.Index(line, ">")
			if start > 0 {
				end := strings.Index(line[start+1:], "<")
				if end > 0 {
					title := strings.TrimSpace(line[start+1 : start+1+end])
					if title != "" && !strings.Contains(title, "DuckDuckGo") {
						results = append(results, title)
					}
				}
			}
		}

		// 查找描述
		if strings.Contains(line, "result__snippet") {
			start := strings.Index(line, ">")
			if start > 0 {
				end := strings.Index(line[start+1:], "<")
				if end > 0 {
					desc := strings.TrimSpace(line[start+1 : start+1+end])
					if desc != "" {
						results = append(results, "  "+desc)
					}
				}
			}
		}
	}

	// 如果没有找到结构化数据，尝试提取纯文本
	if len(results) == 0 {
		// 移除HTML标签
		re := regexp.MustCompile(`<[^>]*>`)
		cleanText := re.ReplaceAllString(content, "")
		cleanText = strings.TrimSpace(cleanText)

		// 分割成句子并取前几行有意义的内容
		lines := strings.Split(cleanText, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && len(line) > 10 && !strings.Contains(line, "DuckDuckGo") {
				results = append(results, line)
				if len(results) >= 5 {
					break
				}
			}
		}
	}

	if len(results) == 0 {
		return "No search results found. This might be due to the search query or temporary API limitations.", nil
	}

	return strings.Join(results, "\n"), nil
}
