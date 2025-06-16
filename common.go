package assembly_go

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// download는 실제 HTTP GET 요청 및 재시도를 처리하는 내부 헬퍼 함수입니다.
func (c *Client) download(req *http.Request) ([]byte, error) {
	var body []byte
	var err error
	var resp *http.Response

	for i := 0; i < c.retries; i++ {
		if i > 0 {
			time.Sleep(2 * time.Second) // 재시도 전 2초 대기
		}

		resp, err = c.httpClient.Do(req)
		if err != nil {
			continue // 요청 자체에 실패하면 재시도
		}

		if resp.StatusCode == http.StatusOK {
			body, err = io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
			}
			return body, nil // 성공 시 즉시 반환
		}
		resp.Body.Close()
	}

	// 모든 재시도 실패 시 마지막 응답과 에러를 기반으로 에러 반환
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}
	return nil, fmt.Errorf("%w: received status code %d", ErrDownloadFailed, resp.StatusCode)
}
