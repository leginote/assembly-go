// assembly_go/downloader.go
package assembly_go

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

// downloadBillPdf는 innerBillId를 사용하여 법안 PDF 원문을 다운로드합니다.
// 이는 기존 leginote-worker-bill/downloader/billpdf.go의 로직을 가져옵니다.
func (c *Client) downloadBillPdf(innerBillId string) ([]byte, error) {
	if innerBillId == "" {
		return nil, ErrInvalidID
	}

	// URL 정의 및 파라미터 설정
	params := url.Values{}
	params.Add("dummy", "dummy")
	params.Add("bookId", innerBillId)
	params.Add("type", "1")
	fullURL := fmt.Sprintf("%s/filegate/sender30?%s", c.baseURL, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	pdfBinary, err := c.download(req) // Client의 공통 download 헬퍼 함수 사용
	if err != nil {
		return nil, err
	}

	// isHTMLContent는 SDK 내부에서만 사용되므로 여기에 정의하거나, 별도 util 파일로 분리 가능
	// leginote-worker-bill/util/validation.go 로직 사용
	if isHTMLContent(pdfBinary) {
		return nil, fmt.Errorf("%w (innerBillId: %s)", ErrHTMLContent, innerBillId)
	}

	return pdfBinary, nil
}

// downloadMeetingRecordPdf는 API가 제공하는 최종 회의록 URL을 받아 PDF를 다운로드합니다.
// 기존 leginote-assembly-go/meeting_record.go 의 로직을 가져옵니다.
func (c *Client) downloadPdfWithUrl(pdfURL string) ([]byte, error) {
	if pdfURL == "" {
		return nil, ErrInvalidID
	}

	req, err := http.NewRequest("GET", pdfURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	return c.download(req) // Client의 공통 download 헬퍼 함수 사용
}

// isHTMLContent는 다운로드된 데이터가 HTML 문서인지 확인합니다. (이 함수는 이 파일 내에서만 사용될 수 있습니다.)
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
