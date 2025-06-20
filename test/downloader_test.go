package assembly_go_test

import (
	"assembly_go"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 테스트용 모의 서버와 클라이언트를 설정하는 헬퍼 함수
func setupTestClient(handler http.HandlerFunc) (*httptest.Server, *assembly_go.Client) {

	server := httptest.NewServer(handler)
	// WithBaseURL을 사용하여 클라이언트가 모의 서버를 가리키도록 설정
	// downloader.go의 downloadBillPdf는 "/filegate/sender30" 경로를 사용하므로
	// BaseURL에 해당 경로를 포함하지 않도록 순수 서버 주소만 전달합니다.
	client, _ := assembly_go.NewClient("TEST_API_KEY", assembly_go.WithBaseURL(server.URL))
	return server, client
}

func TestDownloadBillPdf(t *testing.T) {
	t.Run("성공적인 PDF 다운로드", func(t *testing.T) {
		expectedPdfData := []byte("%PDF-1.4 sample content")
		innerBillId := "test-bill-id"

		handler := func(w http.ResponseWriter, r *http.Request) {
			// 요청 경로와 파라미터가 올바른지 확인
			if r.URL.Path != "/filegate/sender30" {
				t.Errorf("예상 경로 /filegate/sender30, 실제 경로 %s", r.URL.Path)
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}
			if r.URL.Query().Get("bookId") != innerBillId {
				t.Errorf("예상 bookId %s, 실제 bookId %s", innerBillId, r.URL.Query().Get("bookId"))
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			w.Write(expectedPdfData)
		}
		server, client := setupTestClient(handler)
		defer server.Close()

		pdfData, err := client.DownloadBill(innerBillId)
		if err != nil {
			t.Fatalf("다운로드 중 에러 발생: %v", err)
		}

		if !bytes.Equal(pdfData, expectedPdfData) {
			t.Errorf("다운로드된 데이터가 예상과 다릅니다.")
		}
	})

	t.Run("빈 ID로 인한 에러", func(t *testing.T) {
		// 모의 서버가 필요 없는 테스트
		client, _ := assembly_go.NewClient("TEST_API_KEY")
		_, err := client.DownloadBill("")
		if !errors.Is(err, assembly_go.ErrInvalidID) {
			t.Errorf("예상 에러: %v, 실제 에러: %v", assembly_go.ErrInvalidID, err)
		}
	})

	t.Run("서버 에러로 인한 다운로드 실패", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		server, client := setupTestClient(handler)
		defer server.Close()

		_, err := client.DownloadBill("any-id")
		if err == nil {
			t.Fatal("에러가 발생해야 했지만, 발생하지 않았습니다.")
		}
	})

	t.Run("HTML 컨텐츠 수신 시 에러", func(t *testing.T) {
		htmlContent := []byte("<!DOCTYPE html><html><body>Error</body></html>")
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Write(htmlContent)
		}
		server, client := setupTestClient(handler)
		defer server.Close()

		_, err := client.DownloadBill("any-id")
		if !errors.Is(err, assembly_go.ErrHTMLContent) {
			t.Errorf("예상 에러: %v, 실제 에러: %v", assembly_go.ErrHTMLContent, err)
		}
	})
}

func TestDownloadPdfWithUrl(t *testing.T) {
	t.Run("성공적인 URL 기반 PDF 다운로드", func(t *testing.T) {
		expectedPdfData := []byte("%PDF-1.4 from url")
		handler := func(w http.ResponseWriter, r *http.Request) {
			// 특정 경로 요청에만 응답
			if r.URL.Path == "/test.pdf" {
				w.Write(expectedPdfData)
			} else {
				http.NotFound(w, r)
			}
		}
		server, client := setupTestClient(handler)
		defer server.Close()

		pdfUrl := fmt.Sprintf("%s/test.pdf", server.URL)
		pdfData, err := client.DownloadBill(pdfUrl)
		if err != nil {
			t.Fatalf("다운로드 중 에러 발생: %v", err)
		}

		if !bytes.Equal(pdfData, expectedPdfData) {
			t.Errorf("다운로드된 데이터가 예상과 다릅니다.")
		}
	})

	t.Run("빈 URL로 인한 에러", func(t *testing.T) {
		client, _ := assembly_go.NewClient("TEST_API_KEY")
		_, err := client.DownloadBill("")
		if !errors.Is(err, assembly_go.ErrInvalidID) {
			t.Errorf("예상 에러: %v, 실제 에러: %v", assembly_go.ErrInvalidID, err)
		}
	})

	t.Run("서버 에러로 인한 다운로드 실패", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		server, client := setupTestClient(handler)
		defer server.Close()

		_, err := client.DownloadBill(server.URL + "/any.pdf")
		if err == nil {
			t.Fatal("에러가 발생해야 했지만, 발생하지 않았습니다.")
		}
	})
}

func TestIsHTMLContent(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{"HTML Doctype", []byte("<!DOCTYPE html>..."), true},
		{"Lowercase html tag", []byte("<html>..."), true},
		{"Uppercase HTML tag", []byte("<HTML>..."), true},
		{"Not HTML", []byte("%PDF-1.4..."), false},
		{"Empty slice", []byte(""), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if result := assembly_go.IsHTMLContent(tc.input); result != tc.expected {
				t.Errorf("입력값: %s, 예상: %v, 결과: %v", tc.input, tc.expected, result)
			}
		})
	}
}
