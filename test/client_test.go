package assembly_go_test

import (
	"assembly_go"
	"assembly_go/models"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

// mockServerAndClient는 테스트를 위한 모의 서버와 해당 서버를 가리키는 Client를 생성합니다.
func mockServerAndClient(handler http.HandlerFunc) (*httptest.Server, *assembly_go.Client) {
	server := httptest.NewServer(handler)
	client, _ := assembly_go.NewClient("TEST_API_KEY", assembly_go.WithBaseURL(server.URL))
	return server, client
}

func TestNewClient(t *testing.T) {
	t.Run("API 키가 필수임을 확인", func(t *testing.T) {
		_, err := assembly_go.NewClient("")
		if err == nil {
			t.Error("API 키가 없을 때 에러를 반환해야 하지만, nil을 반환했습니다.")
		}
	})

	t.Run("기본 옵션으로 클라이언트 생성", func(t *testing.T) {
		client, err := assembly_go.NewClient("TEST_API_KEY")
		if err != nil {
			t.Fatalf("클라이언트 생성 시 에러 발생: %v", err)
		}
		if client == nil {
			t.Fatal("클라이언트가 nil입니다.")
		}
	})

	t.Run("WithTimeout 옵션이 실제 타임아웃을 유발하는지 테스트", func(t *testing.T) {
		// 일부러 응답을 지연시키는 모의 서버
		slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(50 * time.Millisecond) // 클라이언트 타임아웃보다 길게 설정
			w.WriteHeader(http.StatusOK)
		}))
		defer slowServer.Close()

		// 매우 짧은 타임아웃을 가진 클라이언트 생성
		client, err := assembly_go.NewClient("TEST_API_KEY",
			assembly_go.WithBaseURL(slowServer.URL),
			assembly_go.WithTimeout(20*time.Millisecond), // 서버 지연시간보다 짧게 설정
		)
		if err != nil {
			t.Fatalf("클라이언트 생성 시 에러 발생: %v", err)
		}

		// API 호출 시 타임아웃 에러가 발생하는지 확인
		_, err = client.FetchApiData("any-endpoint", http.MethodGet, nil)
		if err == nil {
			t.Fatal("타임아웃 에러가 발생해야 하지만, 에러가 발생하지 않았습니다.")
		}

		// 에러가 url.Error 타입이며, Timeout() 메서드가 true를 반환하는지 확인
		var urlErr *url.Error
		if errors.As(err, &urlErr) && urlErr.Timeout() {
			// 기대했던 타임아웃 에러이므로 테스트 성공
		} else {
			t.Errorf("기대했던 타임아웃 에러가 아닙니다. 실제 에러: %v", err)
		}
	})
}

func TestStructToMapString(t *testing.T) {
	type TestStruct struct {
		Name    string `json:"name"`
		Age     int    `json:"age"`
		IsMajor bool   `json:"is_major"`
		Field   string `json:"-"` // 무시되어야 할 필드
	}

	t.Run("구조체를 map[string]string으로 변환", func(t *testing.T) {
		s := TestStruct{Name: "John", Age: 30, IsMajor: true, Field: "ignore"}
		expected := map[string]string{
			"name":     "John",
			"age":      "30",
			"is_major": "true",
		}

		result, err := assembly_go.StructToMapString(s)
		if err != nil {
			t.Fatalf("StructToMapString 변환 중 에러 발생: %v", err)
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("기대값: %v, 결과값: %v", expected, result)
		}
	})

	t.Run("비-구조체 입력에 대한 에러 처리", func(t *testing.T) {
		_, err := assembly_go.StructToMapString("not a struct")
		if err == nil {
			t.Error("비-구조체 입력에 대해 에러를 반환해야 하지만, nil을 반환했습니다.")
		}
	})
}

func TestFetchApiData(t *testing.T) {
	t.Run("API 데이터 가져오기 성공", func(t *testing.T) {
		expectedBody := `{"status": "ok"}`
		handler := func(w http.ResponseWriter, r *http.Request) {
			// URL 파라미터 검증
			if r.URL.Query().Get("KEY") != "TEST_API_KEY" {
				t.Error("API 키가 쿼리 파라미터에 없습니다.")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if r.URL.Query().Get("Type") != "json" {
				t.Error("Type=json 파라미터가 쿼리에 없습니다.")
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			fmt.Fprint(w, expectedBody)
		}
		server, client := mockServerAndClient(handler)
		defer server.Close()

		body, err := client.FetchApiData("test-endpoint", http.MethodGet, nil)
		if err != nil {
			t.Fatalf("FetchApiData 실행 중 에러 발생: %v", err)
		}

		if string(body) != expectedBody {
			t.Errorf("기대값: %s, 결과값: %s", expectedBody, string(body))
		}
	})

	t.Run("API 서버 에러 처리", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		server, client := mockServerAndClient(handler)
		defer server.Close()

		_, err := client.FetchApiData("error-endpoint", http.MethodGet, nil)
		if err == nil {
			t.Fatal("서버 에러 시 에러를 반환해야 하지만, nil을 반환했습니다.")
		}

		if !errors.Is(err, assembly_go.ErrDownloadFailed) {
			t.Errorf("기대 에러: %v, 실제 에러: %v", assembly_go.ErrDownloadFailed, err)
		}
	})
}

// 각 Fetch 메서드에 대한 테스트
func TestFetchMethods(t *testing.T) {
	// 성공 시 반환될 모의 JSON 응답
	// 각 API의 실제 응답 구조를 기반으로 작성되었습니다.
	mockResponses := map[string]string{
		"TVBPMBILL11":       `{"TVBPMBILL11":[{"head":[{"list_total_count":1, "RESULT":{"CODE":"INFO-000","MESSAGE":"정상 처리되었습니다."}}],"row":[{"BILL_ID":"PRC_E2O3A0C7A0E9A1E0E064E0A11A0A045B"}]}]}`,
		"ALLNAMEMBER":       `{"ALLNAMEMBER":[{"head":[{"list_total_count":1, "RESULT":{"CODE":"INFO-000","MESSAGE":"정상 처리되었습니다."}}],"row":[{"NAAS_CD":"1234"}]}]}`,
		"VCONFBILLCONFLIST": `{"VCONFBILLCONFLIST":[{"head":[{"list_total_count":1, "RESULT":{"CODE":"INFO-000","MESSAGE":"정상 처리되었습니다."}}],"row":[{"BILL_ID":"PRC_E2O3A0C7A0E9A1E0E064E0A11A0A045B"}]}]}`,
		"VCONFPHCONFLIST":   `{"VCONFPHCONFLIST":[{"head":[{"list_total_count":1, "RESULT":{"CODE":"INFO-000","MESSAGE":"정상 처리되었습니다."}}],"row":[{"CONF_ID":"1"}]}]}`,
		"nojepdqqaweusdfbi": `{"nojepdqqaweusdfbi":[{"head":[{"list_total_count":1, "RESULT":{"CODE":"INFO-000","MESSAGE":"정상 처리되었습니다."}}],"row":[{"BILL_ID":"PRC_E2O3A0C7A0E9A1E0E064E0A11A0A045B"}]}]}`,
		"nprlapfmaufmqytet": `{"nprlapfmaufmqytet":[{"head":[{"list_total_count":1, "RESULT":{"CODE":"INFO-000","MESSAGE":"정상 처리되었습니다."}}],"row":[{"DAESU":"21"}]}]}`,
		"nwvrqwxyaytdsfvhu": `{"nwvrqwxyaytdsfvhu":[{"head":[{"list_total_count":1, "RESULT":{"CODE":"INFO-000","MESSAGE":"정상 처리되었습니다."}}],"row":[{"HG_NM":"홍길동"}]}]}`,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		endpoint := r.URL.Path[1:] // URL 경로에서 맨 앞의 '/' 제거
		if resp, ok := mockResponses[endpoint]; ok {
			fmt.Fprint(w, resp)
		} else {
			http.NotFound(w, r)
		}
	}
	server, client := mockServerAndClient(http.HandlerFunc(handler))
	defer server.Close()

	// 테스트 케이스 정의
	testCases := []struct {
		name      string
		fetchFunc func() (interface{}, error)
		// 반환된 결과가 유효한지 검증하는 함수
		validator func(t *testing.T, resp interface{})
	}{
		{
			name: "FetchBills",
			fetchFunc: func() (interface{}, error) {
				return client.FetchBills(models.TVBPMBILL11RequestParams{})
			},
			validator: func(t *testing.T, resp interface{}) {
				r := resp.(*models.TVBPMBILL11Response)
				if len(r.TVBPMBILL11) == 0 || len(r.TVBPMBILL11[0].Rows) == 0 {
					t.Error("응답 데이터가 비어있습니다.")
				} else if r.TVBPMBILL11[0].Rows[0].BillId != "PRC_E2O3A0C7A0E9A1E0E064E0A11A0A045B" {
					t.Errorf("잘못된 데이터가 파싱되었습니다: %s", r.TVBPMBILL11[0].Rows[0].BillId)
				}
			},
		},
		{
			name: "FetchAllMembers",
			fetchFunc: func() (interface{}, error) {
				return client.FetchAllMembers(models.AllNameMemberRequestParams{})
			},
			validator: func(t *testing.T, resp interface{}) {
				r := resp.(*models.AllNameMemberResponse)
				if len(r.AllNameMember) == 0 || len(r.AllNameMember[0].Rows) == 0 {
					t.Error("응답 데이터가 비어있습니다.")
				} else if r.AllNameMember[0].Rows[0].NaasCd != "1234" {
					t.Errorf("잘못된 데이터가 파싱되었습니다: %s", r.AllNameMember[0].Rows[0].NaasCd)
				}
			},
		},
		{
			name: "FetchBillConferenceList",
			fetchFunc: func() (interface{}, error) {
				return client.FetchBillConferenceList(models.VCONFBILLCONFLISTRequestParams{})
			},
			validator: func(t *testing.T, resp interface{}) {
				r := resp.(*models.VCONFBILLCONFLISTResponse)
				if len(r.VCONFBILLCONFLIST) == 0 || len(r.VCONFBILLCONFLIST[0].Rows) == 0 {
					t.Error("응답 데이터가 비어있습니다.")
				} else if r.VCONFBILLCONFLIST[0].Rows[0].BillId != "PRC_E2O3A0C7A0E9A1E0E064E0A11A0A045B" {
					t.Errorf("잘못된 데이터가 파싱되었습니다: %s", r.VCONFBILLCONFLIST[0].Rows[0].BillId)
				}
			},
		},
		{
			name: "FetchMeetingConferenceList",
			fetchFunc: func() (interface{}, error) {
				return client.FetchMeetingConferenceList(models.VCONFPHCONFLISTRequestParams{})
			},
			validator: func(t *testing.T, resp interface{}) {
				r := resp.(*models.VCONFPHCONFLISTResponse)
				if len(r.VCONFPHCONFLIST) == 0 || len(r.VCONFPHCONFLIST[0].Rows) == 0 {
					t.Error("응답 데이터가 비어있습니다.")
				} else if r.VCONFPHCONFLIST[0].Rows[0].CONF_ID != "1" {
					t.Errorf("잘못된 데이터가 파싱되었습니다: %s", r.VCONFPHCONFLIST[0].Rows[0].CONF_ID)
				}
			},
		},
		{
			name: "FetchMemberVoteResult",
			fetchFunc: func() (interface{}, error) {
				return client.FetchMemberVoteResult(models.NojepdqqaweusdfbiRequestParams{})
			},
			validator: func(t *testing.T, resp interface{}) {
				r := resp.(*models.NojepdqqaweusdfbiResponse)
				if len(r.Nojepdqqaweusdfbi) == 0 || len(r.Nojepdqqaweusdfbi[0].Rows) == 0 {
					t.Error("응답 데이터가 비어있습니다.")
				} else if r.Nojepdqqaweusdfbi[0].Rows[0].BillID != "PRC_E2O3A0C7A0E9A1E0E064E0A11A0A045B" {
					t.Errorf("잘못된 데이터가 파싱되었습니다: %s", r.Nojepdqqaweusdfbi[0].Rows[0].BillID)
				}
			},
		},
		{
			name: "FetchHistoricalMembers",
			fetchFunc: func() (interface{}, error) {
				return client.FetchHistoricalMembers(models.NprlapfmaufmqytetRequestParams{})
			},
			validator: func(t *testing.T, resp interface{}) {
				r := resp.(*models.NprlapfmaufmqytetResponse)
				if len(r.Nprlapfmaufmqytet) == 0 || len(r.Nprlapfmaufmqytet[0].Rows) == 0 {
					t.Error("응답 데이터가 비어있습니다.")
				} else if r.Nprlapfmaufmqytet[0].Rows[0].DaeSu != "21" {
					t.Errorf("잘못된 데이터가 파싱되었습니다: %s", r.Nprlapfmaufmqytet[0].Rows[0].DaeSu)
				}
			},
		},
		{
			name: "FetchMemberDetails",
			fetchFunc: func() (interface{}, error) {
				return client.FetchMemberDetails(models.NwvrqwxyaytdsfvhuRequestParams{})
			},
			validator: func(t *testing.T, resp interface{}) {
				r := resp.(*models.NwvrqwxyaytdsfvhuResponse)
				if len(r.Nwvrqwxyaytdsfvhu) == 0 || len(r.Nwvrqwxyaytdsfvhu[0].Rows) == 0 {
					t.Error("응답 데이터가 비어있습니다.")
				} else if r.Nwvrqwxyaytdsfvhu[0].Rows[0].HgNm != "홍길동" {
					t.Errorf("잘못된 데이터가 파싱되었습니다: %s", r.Nwvrqwxyaytdsfvhu[0].Rows[0].HgNm)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name+" 성공", func(t *testing.T) {
			resp, err := tc.fetchFunc()
			if err != nil {
				t.Fatalf("데이터를 가져오는 중 에러 발생: %v", err)
			}
			if resp == nil {
				t.Fatal("응답이 nil입니다.")
			}
			tc.validator(t, resp)
		})

		t.Run(tc.name+" JSON 파싱 에러", func(t *testing.T) {
			// 고의로 잘못된 JSON을 반환하는 서버 설정
			errServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"invalid_json":}`)
			}))
			defer errServer.Close()

			errClient, _ := assembly_go.NewClient("TEST_API_KEY", assembly_go.WithBaseURL(errServer.URL))

			// 임시 클라이언트로 테스트 함수 실행
			var tempFetchFunc func() (interface{}, error)
			switch tc.name {
			case "FetchBills":
				tempFetchFunc = func() (interface{}, error) { return errClient.FetchBills(models.TVBPMBILL11RequestParams{}) }
			case "FetchAllMembers":
				tempFetchFunc = func() (interface{}, error) { return errClient.FetchAllMembers(models.AllNameMemberRequestParams{}) }
			case "FetchBillConferenceList":
				tempFetchFunc = func() (interface{}, error) {
					return errClient.FetchBillConferenceList(models.VCONFBILLCONFLISTRequestParams{})
				}
			case "FetchMeetingConferenceList":
				tempFetchFunc = func() (interface{}, error) {
					return errClient.FetchMeetingConferenceList(models.VCONFPHCONFLISTRequestParams{})
				}
			case "FetchMemberVoteResult":
				tempFetchFunc = func() (interface{}, error) {
					return errClient.FetchMemberVoteResult(models.NojepdqqaweusdfbiRequestParams{})
				}
			case "FetchHistoricalMembers":
				tempFetchFunc = func() (interface{}, error) {
					return errClient.FetchHistoricalMembers(models.NprlapfmaufmqytetRequestParams{})
				}
			case "FetchMemberDetails":
				tempFetchFunc = func() (interface{}, error) {
					return errClient.FetchMemberDetails(models.NwvrqwxyaytdsfvhuRequestParams{})
				}
			}

			_, err := tempFetchFunc()
			if err == nil {
				t.Error("잘못된 JSON에 대해 에러를 반환해야 하지만, nil을 반환했습니다.")
			}
		})
	}
}

// Client 구조체의 비공개 필드(httpClient)에 접근하기 위한 헬퍼 메서드 추가
// 실제 assembly_go 패키지에 이 코드를 추가해야 합니다.
// func (c *Client) GetHttpClient() *http.Client {
// 	return c.httpClient
// }
// 위 헬퍼가 실제 코드에 없다는 가정 하에, 테스트 코드 내에서만 사용할 수 있는 방법을 사용합니다.
// 하지만 Go에서는 패키지 외부에서 비공개 필드에 직접 접근할 수 없으므로,
// 테스트를 위해 공개 헬퍼를 추가하는 것이 일반적입니다.
// 여기서는 `WithTimeout` 옵션이 잘 작동하는지 검증하기 위해 GetHttpClient가 있다고 가정합니다.

// assembly_go 패키지 내부에 다음 함수를 추가했다고 가정하고 테스트를 진행합니다.
// (실제 프로젝트에서는 이 함수를 client.go에 추가해야 합니다)
/*
package assembly_go

import "net/http"

func (c *Client) GetHttpClient() *http.Client {
    return c.httpClient
}
*/
