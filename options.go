package assembly_go

import (
	"net/http"
	"time"
)

// Option은 Client 설정을 위한 함수 타입입니다.
type Option func(*Client)

// WithHTTPClient는 사용자 정의 http.Client를 설정하는 옵션입니다.
// 타임아웃 등을 직접 제어하고 싶을 때 유용합니다.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL은 API 요청의 기본 URL을 설정하는 옵션입니다.
// 기본값은 "https://likms.assembly.go.kr" 입니다.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithRetries는 요청 실패 시 재시도 횟수를 설정하는 옵션입니다.
// 기본값은 3입니다.
func WithRetries(count int) Option {
	return func(c *Client) {
		c.retries = count
	}
}

// WithTimeout은 요청 타임아웃을 설정하는 옵션입니다.
// 이 옵션은 새로운 http.Client를 생성하여 설정합니다.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}
