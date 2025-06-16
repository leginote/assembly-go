// 국회의원 정보 통합 API
package models

// 국회의원 상세정보 조회 요청 파라미터
type AllNameMemberRequestParams struct {
	Key    string `json:"KEY" validate:"required"`    // 인증키 (필수)
	Type   string `json:"Type" validate:"required"`   // 호출 문서 타입 (xml, json) (필수)
	Pindex string `json:"pIndex" validate:"required"` // 페이지 위치 (필수)
	Psize  string `json:"pSize" validate:"required"`  // 페이지 당 요청 숫자 (필수)
}

type AllNameMemberOptionalParams struct {
	NAAS_CD      *string `json:"NAAS_CD,omitempty"`      // 국회의원코드
	NAAS_NM      *string `json:"NAAS_NM,omitempty"`      // 국회의원명
	PLPT_NM      *string `json:"PLPT_NM,omitempty"`      // 정당명
	BLNG_CMIT_NM *string `json:"BLNG_CMIT_NM,omitempty"` // 소속위원회명
}

// API 응답 구조체
type AllNameMemberResponse struct {
	AllNameMember []AllNameMember `json:"ALLNAMEMBER"`
}

type AllNameMember struct {
	Head []AllNameMemberHead `json:"head"`
	Rows []AllNameMemberRow  `json:"row"`
}

type AllNameMemberHead struct {
	ListTotalCount int                 `json:"list_total_count"`
	Result         AllNameMemberResult `json:"RESULT"`
}

type AllNameMemberResult struct {
	Code    string `json:"CODE"`
	Message string `json:"MESSAGE"`
}

// MemberVoteResponseRow represents each row of data in the API response
type AllNameMemberRow struct {
	NaasCd        string `json:"NAAS_CD"`         // 국회의원 코드
	NaasNm        string `json:"NAAS_NM"`         // 국회의원한글)
	NaasChNm      string `json:"NAAS_CH_NM"`      // 국회의원한자)
	NaasEnNm      string `json:"NAAS_EN_NM"`      // 국회의원영문)
	BirdyDivCd    string `json:"BIRDY_DIV_CD"`    // 생년월일 구분코드
	BirdyDt       string `json:"BIRDY_DT"`        // 생년월일
	DtyNm         string `json:"DTY_NM"`          // 직업명
	PlptNm        string `json:"PLPT_NM"`         // 배우자명
	ElecdNm       string `json:"ELECD_NM"`        // 선거구명
	ElecdDivNm    string `json:"ELECD_DIV_NM"`    // 선거구 구분명
	CmitNm        string `json:"CMIT_NM"`         // 위원회명
	BlngCmitNm    string `json:"BLNG_CMIT_NM"`    // 소속위원회명
	RlctDivNm     string `json:"RLCT_DIV_NM"`     // 역력구분명
	GteltEraco    string `json:"GTELT_ERACO"`     // 전화번호(국가번호 포함)
	NtrDiv        string `json:"NTR_DIV"`         // 국가 구분
	NaasTelNo     string `json:"NAAS_TEL_NO"`     // 국회의원 전화번호
	NaasEmailAddr string `json:"NAAS_EMAIL_ADDR"` // 국회의원 이메일 주소
	NaasHpUrl     string `json:"NAAS_HP_URL"`     // 국회의원 홈페이지 URL
	AideNm        string `json:"AIDE_NM"`         // 수석비서실장명
	ChfScrtNm     string `json:"CHF_SCRT_NM"`     // 비서실장명
	ScrtNm        string `json:"SCRT_NM"`         // 비서명
	BrfHst        string `json:"BRF_HST"`         // 경력사항
	OffmRnumNo    string `json:"OFFM_RNUM_NO"`    // 국회사무실 연락번호
	NaasPic       string `json:"NAAS_PIC"`        // 사진 정보
}
