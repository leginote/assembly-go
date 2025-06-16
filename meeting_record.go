package assembly_go

import (
	"fmt"
	"net/http"
)

// DownloadMeetingRecord는 Downloader 인터페이스를 구현합니다.
// API가 제공하는 최종 회의록 URL을 받아 PDF를 다운로드합니다.
// 이 메소드는 내부적으로 크롤링을 수행하지 않습니다.
func (c *Client) DownloadMeetingRecord(pdfURL string) ([]byte, error) {
	if pdfURL == "" {
		return nil, ErrInvalidID // 빈 URL에 대한 에러 처리
	}

	req, err := http.NewRequest("GET", pdfURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}
	// 필요시 User-Agent 등 헤더 설정
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	// 공통 다운로드 헬퍼 함수를 호출하여 바로 다운로드 실행
	return c.download(req)
}
