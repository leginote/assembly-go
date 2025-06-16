// 의안별 회의록 목록
package models

type VCONFBILLCONFLISTRequestParams struct {
	Key     string `json:"KEY" validate:"required"`    // 인증키 (필수)
	Type    string `json:"Type" validate:"required"`   // 호출 문서 타입 (xml, json) (필수)
	Pindex  int    `json:"pIndex" validate:"required"` // 페이지 위치 (필수)
	Psize   int    `json:"pSize" validate:"required"`  // 페이지 당 요청 숫자 (필수)
	BILL_ID string `json:BILL_ID validate:"required"`  // 의안ID
}

type VCONFBILLCONFLISTResponse struct {
	VCONFBILLCONFLIST []VCONFBILLCONFLIST `json:"VCONFBILLCONFLIST"`
}

type VCONFBILLCONFLIST struct {
	Head []VCONFBILLCONFLISTHead `json:"head"`
	Rows []VCONFBILLCONFLISTRow  `json:"row"`
}

type VCONFBILLCONFLISTHead struct {
	ListTotalCount int                     `json:"list_total_count"` // 총 리스트 수
	Result         VCONFBILLCONFLISTResult `json:"RESULT"`           // 처리 결과
}

type VCONFBILLCONFLISTResult struct {
	Code    string `json:"CODE"`    // 처리 코드
	Message string `json:"MESSAGE"` // 처리 메시지
}

type VCONFBILLCONFLISTRow struct {
	BillId         string `json:"BILL_ID"`  // 의안 ID
	BillName       string `json:"BILL_NM"`  // 의안명
	ConferenceKind string `json:"CONF_KND"` // 회의록 종류
	ConferenceId   string `json:"CONF_ID"`  // 회의록 ID
	EraCode        string `json:"ERACO"`    // 대수
	Session        string `json:"SESS"`     // 회기
	Degree         string `json:"DGR"`      // 차수
	ConferenceDate string `json:"CONF_DT"`  // 회의 날짜
	DownloadUrl    string `json:"DOWN_URL"` // 다운로드 URL
}
