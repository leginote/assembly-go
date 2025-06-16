package assembly_go

import (
	"assembly_go/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect" // StructToMapString에 필요
	"strconv" // StructToMapString에 필요
	"time"
)

// Downloader는 국회 문서 다운로드 기능을 정의하는 인터페이스입니다.
// SDK가 외부에 제공하는 공식적인 기능 목록입니다.

type Downloader interface {
	DownloadBill(innerBillId string) ([]byte, error)
	DownloadMeetingRecord(pdfURL string) ([]byte, error) // 기존 이름 유지 (외부 노출 인터페이스)
	// 새로운 OpenAPI 래핑 메서드들을 여기에 추가할 수 있습니다.

	FetchBills(params models.TVBPMBILL11RequestParams) (*models.TVBPMBILL11Response, error)
	FetchAllMembers(params models.AllNameMemberRequestParams) (*models.AllNameMemberResponse, error)
	FetchBillConferenceList(params models.VCONFBILLCONFLISTRequestParams) (*models.VCONFBILLCONFLISTResponse, error)
	FetchMeetingConferenceList(params models.VCONFPHCONFLISTRequestParams) (*models.VCONFPHCONFLISTResponse, error)
	FetchMemberVoteResult(params models.NojepdqqaweusdfbiRequestParams) (*models.NojepdqqaweusdfbiResponse, error)
	FetchHistoricalMembers(params models.NprlapfmaufmqytetRequestParams) (*models.NprlapfmaufmqytetResponse, error)
	FetchMemberDetails(params models.NwvrqwxyaytdsfvhuRequestParams) (*models.NwvrqwxyaytdsfvhuResponse, error)
}

// Client는 Downloader 인터페이스의 구현체입니다.
// 실제 로직을 수행하는 주체이며, HTTP 클라이언트나 재시도 횟수 같은 내부 상태를 가집니다.
type Client struct {
	httpClient *http.Client
	baseURL    string
	retries    int
	apiKey     string // API Key를 Client 구조체에 포함시킵니다.
}

// NewClient는 새로운 SDK 클라이언트를 생성합니다.
// 사용자는 이 함수를 통해 SDK 사용을 시작합니다.
func NewClient(apiKey string, options ...Option) *Client {
	// 기본값 설정
	c := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},      // 기본 타임아웃 30초
		baseURL:    "https://open.assembly.go.kr/portal/openapi", // OpenAPI 기본 URL
		retries:    3,
		apiKey:     apiKey, // API Key 설정
	}

	// 사용자가 제공한 옵션으로 기본값 덮어쓰기
	for _, opt := range options {
		opt(c)
	}

	return c
}

// 이 아래에 인터페이스 메소드들을 구현해 나갈 것입니다.
// 기존 DownloadBill 및 DownloadMeetingRecord 메서드는 그대로 유지하거나,
// Client 구조체의 메서드로 옮겨와 Client의 httpClient를 사용하도록 변경할 수 있습니다.
// 현재 leginote-assembly-go/bill.go 와 leginote-assembly-go/meeting_record.go 에 구현되어 있습니다.
// 이들은 Client의 download 헬퍼 함수를 사용하도록 이미 구성되어 있으므로 적절합니다.

// download는 실제 HTTP GET 요청 및 재시도를 처리하는 내부 헬퍼 함수입니다.
// 기존 leginote/assembly-go/common.go 에서 가져와 Client의 메서드로 변경합니다.
func (c *Client) download(req *http.Request) ([]byte, error) {
	var body []byte
	var err error
	var resp *http.Response

	for i := 0; i < c.retries; i++ {
		if i > 0 {
			time.Sleep(2 * time.Second) // 재시도 전 2초 대기
		}

		resp, err = c.httpClient.Do(req)
		if err != nil {
			continue // 요청 자체에 실패하면 재시도
		}

		if resp.StatusCode == http.StatusOK {
			body, err = io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
			}
			return body, nil // 성공 시 즉시 반환
		}
		resp.Body.Close()
	}

	// 모든 재시도 실패 시 마지막 응답과 에러를 기반으로 에러 반환
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}
	return nil, fmt.Errorf("%w: received status code %d", ErrDownloadFailed, resp.StatusCode)
}

// FetchApiData 함수는 URL, 헤더, 메서드, 요청 인자를 받아 HTTP 요청을 보내고 결과를 반환합니다.
// 이 함수는 leginote-worker-bill/api/assembly/common.go 에서 가져와 SDK 내부 헬퍼 함수로 사용합니다.
func (c *Client) FetchApiData(endpoint string, method string, params map[string]string) ([]byte, error) {
	fullURL := fmt.Sprintf("%s/%s", c.baseURL, endpoint)
	parsedURL, err := url.ParseRequestURI(fullURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	query := parsedURL.Query()
	query.Add("KEY", c.apiKey)
	query.Add("Type", "json") // OpenAPI 응답 타입을 JSON으로 고정
	for k, v := range params {
		query.Add(k, v)
	}
	parsedURL.RawQuery = query.Encode()

	req, err := http.NewRequest(method, parsedURL.String(), nil) // GET 요청이므로 body는 nil
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestFailed, err) // SDK 에러 사용
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:115.0) Gecko/20100101 Firefox/115.0")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err) // SDK 에러 사용
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err) // SDK 에러 사용
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("HTTP request failed with status code %d", resp.StatusCode)
		if len(respBody) > 100 {
			errMsg += fmt.Sprintf(": %s...", respBody[:100])
		} else if len(respBody) > 0 {
			errMsg += fmt.Sprintf(": %s", respBody)
		}
		return nil, fmt.Errorf("%w: %s", ErrDownloadFailed, errMsg) // SDK 에러 사용
	}
	return respBody, nil
}

// StructToMapString는 구조체를 map[string]string으로 변환합니다.
// leginote-worker-bill/worker/worker.go 에서 가져와 SDK 내부 헬퍼 함수로 사용합니다.
func StructToMapString(obj interface{}) (map[string]string, error) {
	result := make(map[string]string)
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct")
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		jsonTag := fieldType.Tag.Get("json")

		key := fieldType.Name
		if jsonTag != "" && jsonTag != "-" {
			key = jsonTag
		}

		switch field.Kind() {
		case reflect.String:
			result[key] = field.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result[key] = strconv.FormatInt(field.Int(), 10)
		case reflect.Float32, reflect.Float64:
			result[key] = strconv.FormatFloat(field.Float(), 'f', -1, 64)
		case reflect.Bool:
			result[key] = strconv.FormatBool(field.Bool())
		default:
			result[key] = fmt.Sprintf("%v", field.Interface())
		}
	}
	return result, nil
}

// FetchBills는 TVBPMBILL11 OpenAPI를 호출하여 법률안 심사 및 처리 정보를 가져옵니다.
func (c *Client) FetchBills(params models.TVBPMBILL11RequestParams) (*models.TVBPMBILL11Response, error) {
	reqParamsMap, err := StructToMapString(params)
	if err != nil {
		return nil, fmt.Errorf("parameter conversion failed: %w", err)
	}

	data, err := c.FetchApiData("TVBPMBILL11", http.MethodGet, reqParamsMap) // TVBPMBILL11은 GET 메소드 사용
	if err != nil {
		return nil, err
	}

	var resp models.TVBPMBILL11Response
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err) // 응답 파싱 실패도 SDK 에러로 처리
	}
	return &resp, nil
}

// FetchAllMembers는 ALLNAMEMBER OpenAPI를 호출하여 국회의원 정보 통합 데이터를 가져옵니다.
func (c *Client) FetchAllMembers(params models.AllNameMemberRequestParams) (*models.AllNameMemberResponse, error) {
	reqParamsMap, err := StructToMapString(params)
	if err != nil {
		return nil, fmt.Errorf("parameter conversion failed: %w", err)
	}

	data, err := c.FetchApiData("ALLNAMEMBER", http.MethodGet, reqParamsMap)
	if err != nil {
		return nil, err
	}

	var resp models.AllNameMemberResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}
	return &resp, nil
}

// FetchBillConferenceList는 VCONFBILLCONFLIST OpenAPI를 호출하여 의안별 회의록 목록을 가져옵니다.
func (c *Client) FetchBillConferenceList(params models.VCONFBILLCONFLISTRequestParams) (*models.VCONFBILLCONFLISTResponse, error) {
	// int 필드가 있는 요청 파라미터를 StructToMapString이 처리할 수 있도록 수정하거나,
	// 여기서 직접 map[string]string을 구성해야 합니다.
	// 현재 StructToMapString은 int를 string으로 변환하는 로직을 가지고 있으므로 활용 가능합니다.
	reqParamsMap, err := StructToMapString(params)
	if err != nil {
		return nil, fmt.Errorf("parameter conversion failed: %w", err)
	}

	data, err := c.FetchApiData("VCONFBILLCONFLIST", http.MethodGet, reqParamsMap)
	if err != nil {
		return nil, err
	}

	var resp models.VCONFBILLCONFLISTResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}
	return &resp, nil
}

// FetchMeetingConferenceList는 VCONFPHCONFLIST OpenAPI를 호출하여 회의록 통합 데이터를 가져옵니다.
func (c *Client) FetchMeetingConferenceList(params models.VCONFPHCONFLISTRequestParams) (*models.VCONFPHCONFLISTResponse, error) {
	reqParamsMap, err := StructToMapString(params)
	if err != nil {
		return nil, fmt.Errorf("parameter conversion failed: %w", err)
	}

	data, err := c.FetchApiData("VCONFPHCONFLIST", http.MethodGet, reqParamsMap)
	if err != nil {
		return nil, err
	}

	var resp models.VCONFPHCONFLISTResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}
	return &resp, nil
}

// FetchMemberVoteResult는 nojepdqqaweusdfbi OpenAPI를 호출하여 국회의원 본회의 표결정보를 가져옵니다.
func (c *Client) FetchMemberVoteResult(params models.NojepdqqaweusdfbiRequestParams) (*models.NojepdqqaweusdfbiResponse, error) {
	reqParamsMap, err := StructToMapString(params)
	if err != nil {
		return nil, fmt.Errorf("parameter conversion failed: %w", err)
	}

	data, err := c.FetchApiData("nojepdqqaweusdfbi", http.MethodGet, reqParamsMap)
	if err != nil {
		return nil, err
	}

	var resp models.NojepdqqaweusdfbiResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}
	return &resp, nil
}

// FetchHistoricalMembers는 nprlapfmaufmqytet OpenAPI를 호출하여 역대 국회의원 현황을 가져옵니다.
func (c *Client) FetchHistoricalMembers(params models.NprlapfmaufmqytetRequestParams) (*models.NprlapfmaufmqytetResponse, error) {
	reqParamsMap, err := StructToMapString(params)
	if err != nil {
		return nil, fmt.Errorf("parameter conversion failed: %w", err)
	}

	data, err := c.FetchApiData("nprlapfmaufmqytet", http.MethodGet, reqParamsMap)
	if err != nil {
		return nil, err
	}

	var resp models.NprlapfmaufmqytetResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}
	return &resp, nil
}

// FetchMemberDetails는 nwvrqwxyaytdsfvhu OpenAPI를 호출하여 국회의원 인적사항을 가져옵니다.
func (c *Client) FetchMemberDetails(params models.NwvrqwxyaytdsfvhuRequestParams) (*models.NwvrqwxyaytdsfvhuResponse, error) {
	reqParamsMap, err := StructToMapString(params)
	if err != nil {
		return nil, fmt.Errorf("parameter conversion failed: %w", err)
	}

	data, err := c.FetchApiData("nwvrqwxyaytdsfvhu", http.MethodGet, reqParamsMap)
	if err != nil {
		return nil, err
	}

	var resp models.NwvrqwxyaytdsfvhuResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}
	return &resp, nil
}
