package assembly_go

import (
	"net/http"
	"time"
)

// Downloader는 국회 문서 다운로드 기능을 정의하는 인터페이스입니다.
// SDK가 외부에 제공하는 공식적인 기능 목록입니다.
type Downloader interface {
	DownloadBill(innerBillId string) ([]byte, error)
	DownloadMeetingRecord(url string) ([]byte, error)
}

// Client는 Downloader 인터페이스의 구현체입니다.
// 실제 로직을 수행하는 주체이며, HTTP 클라이언트나 재시도 횟수 같은 내부 상태를 가집니다.
type Client struct {
	httpClient *http.Client
	baseURL    string
	retries    int
}

// NewClient는 새로운 SDK 클라이언트를 생성합니다.
// 사용자는 이 함수를 통해 SDK 사용을 시작합니다.
func NewClient(options ...Option) *Client {
	// 기본값 설정
	c := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second}, // 기본 타임아웃 30초
		baseURL:    "https://likms.assembly.go.kr",
		retries:    3,
	}

	// 사용자가 제공한 옵션으로 기본값 덮어쓰기
	for _, opt := range options {
		opt(c)
	}

	return c
}

// 이 아래에 인터페이스 메소드들을 구현해 나갈 것입니다.
// func (c *Client) DownloadBill(innerBillId string) ([]byte, error) { ... }
// func (c *Client) DownloadMeetingRecord(url string) ([]byte, error) { ... }
