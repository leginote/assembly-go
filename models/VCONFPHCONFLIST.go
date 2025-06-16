// 회의록 통합 API 객체
package models

// VCONFPHCONFLISTRequestParams represents request parameters for the VCONFPHCONFLIST API.
type VCONFPHCONFLISTRequestParams struct {
	Key    string `json:"KEY" validate:"required"`    // 인증키 (필수)
	Type   string `json:"Type" validate:"required"`   // 호출 문서 타입 (xml, json) (필수)
	Pindex string `json:"pIndex" validate:"required"` // 페이지 위치 (필수)
	Psize  string `json:"pSize" validate:"required"`  // 페이지 당 요청 숫자 (필수)
	ERACO  string `json:"ERACO" validate:"required"`  // 대수 (필수)
}

// VCONFPHCONFLISTResponse represents the response from the VCONFPHCONFLIST API.
type VCONFPHCONFLISTResponse struct {
	VCONFPHCONFLIST []VCONFPHCONFLIST `json:"VCONFPHCONFLIST"`
}

// VCONFPHCONFLIST represents a single entry in the VCONFPHCONFLIST API response.
type VCONFPHCONFLIST struct {
	Head []VCONFPHCONFLISTHead `json:"head"`
	Rows []VCONFPHCONFLISTRow  `json:"row"`
}

// VCONFPHCONFLISTHead represents the header information in the VCONFPHCONFLIST API response.
type VCONFPHCONFLISTHead struct {
	ListTotalCount int                   `json:"list_total_count"`
	Result         VCONFPHCONFLISTResult `json:"RESULT"`
}

// VCONFPHCONFLISTResult represents the result information in the VCONFPHCONFLIST API response.
type VCONFPHCONFLISTResult struct {
	Code    string `json:"CODE"`
	Message string `json:"MESSAGE"`
}

// VCONFPHCONFLISTRow represents each row of data in the VCONFPHCONFLIST API response.
type VCONFPHCONFLISTRow struct {
	CONF_ID  string `json:"CONF_ID"`  // 회의ID
	ERACO    string `json:"ERACO"`    // 대수
	SESS     string `json:"SESS"`     // 회기
	DGR      string `json:"DGR"`      // 차수
	CONF_DT  string `json:"CONF_DT"`  // 회의일자
	CONF_KND string `json:"CONF_KND"` // 회의종류
	CMIT_CD  string `json:"CMIT_CD"`  // 위원회코드
	CMIT_NM  string `json:"CMIT_NM"`  // 위원회명
	DOWN_URL string `json:"DOWN_URL"` // 다운로드 URL
}
