package models

// 역대 국회의원 현황 요청 파라미터
type NprlapfmaufmqytetRequestParams struct {
	Key    string `json:"KEY" validate:"required"`    // 인증키 (필수)
	Type   string `json:"Type" validate:"required"`   // 호출 문서 타입 (xml, json) (필수)
	Pindex string `json:"pIndex" validate:"required"` // 페이지 위치 (필수)
	Psize  string `json:"pSize" validate:"required"`  // 페이지 당 요청 숫자 (필수)
	DAESU  string `json:"DAESU" validate:"required"`  // 대수
}

// 역대 국회의원 현황 요청 Optional 파라미터
type NprlapfmaufmqytetOptionalParams struct {
	DAE      *string `json:"DAE,omitempty"`      // 대별 및 소속정당(단체)
	DAE_NM   *string `json:"DAE_NM,omitempty"`   // 대별
	NAME     *string `json:"NAME,omitempty"`     // 이름
	NAME_HAN *string `json:"NAME_HAN,omitempty"` // 이름(한자)
	BIRTH    *string `json:"BIRTH,omitempty"`    // 생년월일
	POSI     *string `json:"POSI,omitempty"`     // 출생지
}

// 역대 국회의원 현황 Response
type NprlapfmaufmqytetResponse struct {
	Nprlapfmaufmqytet []Nprlapfmaufmqytet `json:"nprlapfmaufmqytet"`
}

// 역대 국회의원 현황 Data Structure
type Nprlapfmaufmqytet struct {
	Head []NprlapfmaufmqytetHead `json:"head"`
	Rows []NprlapfmaufmqytetRow  `json:"row"`
}

// Header 정보
type NprlapfmaufmqytetHead struct {
	ListTotalCount int                     `json:"list_total_count"`
	Result         NprlapfmaufmqytetResult `json:"RESULT"`
}

// Result 정보
type NprlapfmaufmqytetResult struct {
	Code    string `json:"CODE"`
	Message string `json:"MESSAGE"`
}

// Row 정보
type NprlapfmaufmqytetRow struct {
	DaeSu   string `json:"DAESU"`    // 대수
	Dae     string `json:"DAE"`      // 대별 및 소속정당(단체)
	DaeNm   string `json:"DAE_ NM"`  // 대별
	Name    string `json:"NAME"`     // 이름
	NameHan string `json:"NAME_HAN"` // 이름(한자)
	Ja      string `json:"JA"`       // 자
	Ho      string `json:"HO"`       // 호
	Birth   string `json:"BIRTH"`    // 생년월일
	Bon     string `json:"BON"`      // 본관
	Posi    string `json:"POSI"`     // 출생지
	Hak     string `json:"HAK"`      // 학력 및 경력
	Hobby   string `json:"HOBBY"`    // 종교 및 취미
	Book    string `json:"BOOK"`     // 저서
	Sang    string `json:"SANG"`     // 상훈
	Dead    string `json:"DEAD"`     // 기타정보(사망일)
	Url     string `json:"URL"`      // 회원정보 확인 헌정회 홈페이지 URL
}
