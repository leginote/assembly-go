package models

//국회의원 인적사항

type NwvrqwxyaytdsfvhuRequestParams struct {
	Key    string `json:"KEY" validate:"required"`    // 인증키 (필수)
	Type   string `json:"Type" validate:"required"`   // 호출 문서 타입 (xml, json) (필수)
	Pindex string `json:"pIndex" validate:"required"` // 페이지 위치 (필수)
	Psize  string `json:"pSize" validate:"required"`  // 페이지 당 요청 숫자 (필수)
}

type NwvrqwxyaytdsfvhuOptionalParams struct {
	HG_NM      *string `json:"HG_NM,omitempty"`      // 이름
	POLY_NM    *string `json:"POLY_NM,omitempty"`    // 정당명
	ORIG_NM    *string `json:"ORIG_NM,omitempty"`    // 선거구
	CMITS      *string `json:"CMITS,omitempty"`      // 소속 위원회 목록
	SEX_GBN_NM *string `json:"SEX_GBN_NM,omitempty"` // 성별
	MONA_CD    *string `json:"MONA_CD,omitempty"`    // 국회의원코드
}

//국회의원 인적사항

type NwvrqwxyaytdsfvhuResponse struct {
	Nwvrqwxyaytdsfvhu []Nwvrqwxyaytdsfvhu `json:"nwvrqwxyaytdsfvhu"`
}

type Nwvrqwxyaytdsfvhu struct {
	Head []NwvrqwxyaytdsfvhuHead `json:"head"`
	Rows []NwvrqwxyaytdsfvhuRow  `json:"row"`
}

type NwvrqwxyaytdsfvhuHead struct {
	ListTotalCount int                     `json:"list_total_count"`
	Result         NwvrqwxyaytdsfvhuResult `json:"RESULT"`
}

type NwvrqwxyaytdsfvhuResult struct {
	Code    string `json:"CODE"`
	Message string `json:"MESSAGE"`
}

// MemberVoteResponseRow represents each row of data in the API response
type NwvrqwxyaytdsfvhuRow struct {
	HgNm       string  `json:"HG_NM"`        // 이름
	HjNm       *string `json:"HJ_NM"`        // 한자명
	EngNm      *string `json:"ENG_NM"`       // 영문명칭
	BthGbnNm   *string `json:"BTH_GBN_NM"`   // 음/양력
	BthDate    *string `json:"BTH_DATE"`     // 생년월일
	JobResNm   *string `json:"JOB_RES_NM"`   // 직책명
	PolyNm     *string `json:"POLY_NM"`      // 정당명
	OrigNm     *string `json:"ORIG_NM"`      // 선거구
	ElectGbnNm *string `json:"ELECT_GBN_NM"` // 선거구구분
	CmitNm     *string `json:"CMIT_NM"`      // 대표 위원회
	Cmits      *string `json:"CMITS"`        // 소속 위원회 목록
	ReeleGbnNm *string `json:"REELE_GBN_NM"` // 재선
	Units      *string `json:"UNITS"`        // 당선
	SexGbnNm   *string `json:"SEX_GBN_NM"`   // 성별
	TelNo      *string `json:"TEL_NO"`       // 전화번호
	EMail      *string `json:"E_MAIL"`       // 이메일
	Homepage   *string `json:"HOMEPAGE"`     // 홈페이지
	Staff      *string `json:"STAFF"`        // 보좌관
	Secretary  *string `json:"SECRETARY"`    // 선임비서관
	Secretary2 *string `json:"SECRETARY2"`   // 비서관
	MonaCd     string  `json:"MONA_CD"`      // 국회의원코드
	MemTitle   *string `json:"MEM_TITLE"`    // 약력
	AssemAddr  *string `json:"ASSEM_ADDR"`   // 사무실 호실
}
