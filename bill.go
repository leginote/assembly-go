package assembly_go

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

// isHTMLContent는 다운로드된 데이터가 HTML 문서인지 확인합니다.
// 기존 leginote-worker-bill/util/validation.go 의 로직을 가져왔습니다.
func isHTMLContent(data []byte) bool {
	htmlIndicators := [][]byte{
		[]byte("<!DOCTYPE html>"),
		[]byte("<html>"),
		[]byte("<HTML>"),
	}
	for _, indicator := range htmlIndicators {
		if bytes.HasPrefix(data, indicator) {
			return true
		}
	}
	return false
}

// DownloadBill은 Downloader 인터페이스를 구현합니다.
// innerBillId를 사용하여 법안 PDF 원문을 다운로드합니다.
func (c *Client) DownloadBill(innerBillId string) ([]byte, error) {
	if innerBillId == "" {
		return nil, ErrInvalidID
	}

	// URL 정의 및 파라미터 설정 (기존 billpdf.go 로직)
	params := url.Values{}
	params.Add("dummy", "dummy")
	params.Add("bookId", innerBillId)
	params.Add("type", "1")
	fullURL := fmt.Sprintf("%s/filegate/sender30?%s", c.baseURL, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}

	// User-Agent 헤더 추가
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	// 공통 다운로드 함수 호출
	pdfBinary, err := c.download(req)
	if err != nil {
		return nil, err
	}

	// 다운로드된 파일이 PDF가 아닌 HTML인지 확인
	if isHTMLContent(pdfBinary) {
		return nil, fmt.Errorf("%w (innerBillId: %s)", ErrHTMLContent, innerBillId)
	}

	return pdfBinary, nil
}
