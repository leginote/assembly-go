package assembly_go

import "errors"

// SDK에서 반환될 수 있는 주요 에러들입니다.
var (
	// ErrInvalidID는 innerBillId가 비어있는 등 ID가 유효하지 않을 때 발생합니다.
	ErrInvalidID = errors.New("provided ID is invalid or empty")

	// ErrRequestFailed는 HTTP 요청 생성에 실패했을 때 발생합니다.
	ErrRequestFailed = errors.New("failed to create HTTP request")

	// ErrDownloadFailed는 HTTP 상태 코드가 200이 아니거나 응답을 읽는 데 실패했을 때 발생합니다.
	ErrDownloadFailed = errors.New("failed to download content")

	// ErrHTMLContent는 다운로드한 내용이 기대했던 파일이 아닌 HTML 문서일 때 발생합니다.
	ErrHTMLContent = errors.New("downloaded content is an HTML document, not the expected file")
)
