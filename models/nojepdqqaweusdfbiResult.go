// 국회의원 본회의 표결정보
package models

import "encoding/json"

type NojepdqqaweusdfbiRequestParams struct {
	Key     string `json:"KEY" validate:"required"`    // 인증키 (필수)
	Type    string `json:"Type" validate:"required"`   // 호출 문서 타입 (xml, json) (필수)
	Pindex  string `json:"pIndex" validate:"required"` // 페이지 위치 (필수)
	Psize   string `json:"pSize" validate:"required"`  // 페이지 당 요청 숫자 (필수)
	AGE     string `json:"AGE"  validate:"required"`   // 대
	BILL_ID string `json:BILL_ID validate:"required"`  // 의안ID

}

type NojepdqqaweusdfbiOptionalParams struct {
	HG_NM             *string `json:"HG_NM,omitempty"`             // 의원
	POLY_NM           *string `json:"POLY_NM,omitempty"`           // 정당
	MEMBER_NO         *string `json:"MEMBER_NO,omitempty"`         // 의원번호
	VOTE_DATE         *string `json:"VOTE_DATE,omitempty"`         // 의결일자
	BILL_NO           *string `json:"BILL_NO,omitempty"`           // 의안번호
	BILL_NAME         *string `json:"BILL_NAME,omitempty"`         // 의안명
	CURR_COMMITTEE    *string `json:"CURR_COMMITTEE,omitempty"`    // 소관위원회
	RESULT_VOTE_MOD   *string `json:"RESULT_VOTE_MOD,omitempty"`   // 표결결과
	CURR_COMMITTEE_ID *string `json:"CURR_COMMITTEE_ID,omitempty"` // 소관위코드
	MONA_CD           *string `json:"MONA_CD,omitempty"`           // 국회의원코드
}

// 국회의원표결정보
// NojepdqqaweusdfbiResponse represents the top-level response structure from the API
type NojepdqqaweusdfbiResponse struct {
	Nojepdqqaweusdfbi []Nojepdqqaweusdfbi `json:"nojepdqqaweusdfbi"`
}

type Nojepdqqaweusdfbi struct {
	Head []NojepdqqaweusdfbiHead `json:"head"`
	Rows []NojepdqqaweusdfbiRow  `json:"row"`
}

type NojepdqqaweusdfbiHead struct {
	ListTotalCount int                     `json:"list_total_count"`
	Result         NojepdqqaweusdfbiResult `json:"RESULT"`
}

type NojepdqqaweusdfbiResult struct {
	Code    string `json:"CODE"`
	Message string `json:"MESSAGE"`
}

// MemberVoteResponseRow represents each row of data in the API response
type NojepdqqaweusdfbiRow struct {
	MemberName              string      `json:"HG_NM"`             // 의원
	MemberNameHanja         *string     `json:"HJ_NM"`             // 한자명
	PartyName               string      `json:"POLY_NM"`           // 정당
	Constituency            string      `json:"ORIG_NM"`           // 선거구
	MemberNumber            string      `json:"MEMBER_NO"`         // 의원번호
	PartyCode               string      `json:"POLY_CD"`           // 소속정당코드
	ConstituencyCode        string      `json:"ORIG_CD"`           // 선거구코드
	VoteDate                string      `json:"VOTE_DATE"`         // 의결일자
	BillNumber              string      `json:"BILL_NO"`           // 의안번호
	BillName                string      `json:"BILL_NAME"`         // 의안명
	BillID                  string      `json:"BILL_ID"`           // 의안ID
	LawTitle                string      `json:"LAW_TITLE"`         // 법률명
	JurisdictionCommittee   string      `json:"CURR_COMMITTEE"`    // 소관위원회
	VoteResult              string      `json:"RESULT_VOTE_MOD"`   // 표결결과
	DepartmentCode          string      `json:"DEPT_CD"`           // 부서코드(사용안함)
	JurisdictionCommitteeID string      `json:"CURR_COMMITTEE_ID"` // 소관위코드
	DisplayOrder            json.Number `json:"DISP_ORDER"`        // 표시정렬순서 (nullable)
	BillURL                 string      `json:"BILL_URL"`          // 의안URL
	BillNameURL             string      `json:"BILL_NAME_URL"`     // 의안링크
	SessionCode             json.Number `json:"SESSION_CD"`        // 회기
	CurrentsCode            json.Number `json:"CURRENTS_CD"`       // 차수
	Age                     json.Number `json:"AGE"`               // 대
	MemberCode              string      `json:"MONA_CD"`           // 국회의원코드
}
