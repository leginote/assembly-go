package models

type TVBPMBILL11RequestParams struct {
	Key    string `json:"KEY" validate:"required"`    // 인증키 (필수)
	Type   string `json:"Type" validate:"required"`   // 호출 문서 타입 (xml, json) (필수)
	Pindex string `json:"pIndex" validate:"required"` // 페이지 위치 (필수)
	Psize  string `json:"pSize" validate:"required"`  // 페이지 당 요청 숫자 (필수)
}

type TVBPMBILL11OptionalParams struct {
	BILL_ID           string `json:"BILL_ID"`           // 의안ID
	BILL_NO           string `json:"BILL_NO"`           // 의안번호
	AGE               string `json:"AGE"`               // 대
	BILL_NAME         string `json:"BILL_NAME"`         // 의안명(한글)
	PROPOSER          string `json:"PROPOSER"`          // 제안자
	PROPOSER_KIND     string `json:"PROPOSER_KIND"`     // 제안자구분
	CURR_COMMITTEE_ID string `json:"CURR_COMMITTEE_ID"` // 소관위코드
	CURR_COMMITTEE    string `json:"CURR_COMMITTEE"`    // 소관위
	PROC_RESULT_CD    string `json:"PROC_RESULT_CD"`    // 본회의심의결과
	PROC_DT           string `json:"PROC_DT"`           // 의결일
}

// 법률안 심사 및 처리(의안검색)
type TVBPMBILL11Response struct {
	TVBPMBILL11 []TVBPMBILL11 `json:"TVBPMBILL11"`
}

type TVBPMBILL11 struct {
	Head []TVBPMBILL11Head `json:"head"`
	Rows []TVBPMBILL11Row  `json:"row"`
}

type TVBPMBILL11Head struct {
	ListTotalCount int               `json:"list_total_count"`
	Result         TVBPMBILL11Result `json:"RESULT"`
}

type TVBPMBILL11Result struct {
	Code    string `json:"CODE"`
	Message string `json:"MESSAGE"`
}

// BillResponseRow represents each row of data in the API response
type TVBPMBILL11Row struct {
	BillId                               string  `json:"BILL_ID"`            // 의안 ID
	BillNumber                           string  `json:"BILL_NO"`            // 의안 번호
	Age                                  string  `json:"AGE"`                // 대수
	BillName                             string  `json:"BILL_NAME"`          // 의안명(한글)
	Proposer                             string  `json:"PROPOSER"`           // 제안자
	ProposerDivision                     string  `json:"PROPOSER_KIND"`      // 제안자 구분
	ProposeDate                          string  `json:"PROPOSE_DT"`         // 제안일
	JurisdictionCommitteeCode            string  `json:"CURR_COMMITTEE_ID"`  // 소관위 코드
	JurisdictionCommittee                string  `json:"CURR_COMMITTEE"`     // 소관위
	SubmitDate                           string  `json:"COMMITTEE_DT"`       // 소관위 회부일
	CommitteeReviewProcessDate           *string `json:"COMMITTEE_PROC_DT"`  // 위원회 심사 처리일 (nullable)
	DetailUrl                            string  `json:"LINK_URL"`           // 의안 상세정보 URL
	LeadProposer                         string  `json:"RST_PROPOSER"`       // 대표발의자
	LegislationAndJudiciaryProcessResult *string `json:"LAW_PROC_RESULT_CD"` // 법사위 처리 결과 코드 (nullable)
	LegislationAndJudiciaryProcessDate   *string `json:"LAW_PROC_DT"`        // 법사위 처리일 (nullable)
	LegislationAndJudiciaryPresentDate   *string `json:"LAW_PRESENT_DT"`     // 법사위 상정일 (nullable)
	LegislationAndJudiciarySubmitDate    *string `json:"LAW_SUBMIT_DT"`      // 법사위 회부일 (nullable)
	CommitteeProcessResult               *string `json:"CMT_PROC_RESULT_CD"` // 소관위 처리 결과 코드 (nullable)
	CommitteeProcessDate                 *string `json:"CMT_PROC_DT"`        // 소관위 처리일 (nullable)
	CommitteePresentDate                 *string `json:"CMT_PRESENT_DT"`     // 소관위 상정일 (nullable)
	LeadProposerCode                     string  `json:"RST_MONA_CD"`        // 대표발의자 코드
	PlenarySessionReviewResult           string  `json:"PROC_RESULT_CD"`     // 본회의 심의 결과
	ResolutionDate                       string  `json:"PROC_DT"`            // 의결일
}
