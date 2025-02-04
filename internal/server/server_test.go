package server

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dayterr/test-go-iq/internal/config"
	"github.com/dayterr/test-go-iq/internal/storage"
)

func TestIncreaseBalance(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	cb := []byte(`{
		"id": 1,
    "balance": 500
	}`)
	ib1 := []byte(`{
    "balance": 500
	}`)
	ib2 := []byte(`{
		"id": 1,
	}`)

	tests := []struct {
		url  string
		name string
		body []byte
		want want
	}{
		{
			url:  "/add",
			name: "test correct body",
			body: cb,
			want: want{
				code: http.StatusOK,
			},
		},
		{
			url:  "/add",
			name: "test incorrect body, without id",
			body: ib1,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			url:  "/add",
			name: "test incorrect body, without id",
			body: ib2,
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := config.GetConfig()
			if err != nil {
				log.Fatal("no config, can't start the program")
			}

			fmt.Println(config)

			h := NewHandler()
			s := storage.NewStorage(config.DatabaseURI)
			h.Storage = s
			r := CreateRouter(h)
			ts := httptest.NewServer(r)
			defer ts.Close()
			req, _ := testRequest(t, ts, http.MethodPost, tt.url, bytes.NewBuffer(tt.body))
			defer req.Body.Close()
			assert.Equal(t, tt.want.code, req.StatusCode, "Возвращаемый код не равен ожидаемому")
		})
	}

}

func TestTransferMoney(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	cb := []byte(`{
		"sender_id": 1,
    "amount": 100,
    "reciever_id": 3
	}`)
	ib1 := []byte(`{
		"amount": 100,
		"reciever_id": 3
	}`)
	ib2 := []byte(`{
		"sender_id": 1,
    "reciever_id": 3
	}`)
	ib3 := []byte(`{
		"sender_id": 1,
    "amount": 100,
	}`)

	tests := []struct {
		url  string
		name string
		body []byte
		want want
	}{
		{
			url:  "/transfer",
			name: "test correct body",
			body: cb,
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			url:  "/transfer",
			name: "test incorrect body, without sender_id",
			body: ib1,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			url:  "/transfer",
			name: "test incorrect body, without amount",
			body: ib2,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			url:  "/transfer",
			name: "test incorrect body, without reciever_id",
			body: ib3,
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := config.GetConfig()
			if err != nil {
				log.Fatal("no config, can't start the program")
			}

			fmt.Println(config)

			h := NewHandler()
			s := storage.NewStorage(config.DatabaseURI)
			h.Storage = s
			r := CreateRouter(h)
			ts := httptest.NewServer(r)
			defer ts.Close()
			req, _ := testRequest(t, ts, http.MethodPost, tt.url, bytes.NewBuffer(tt.body))
			defer req.Body.Close()
			assert.Equal(t, tt.want.code, req.StatusCode, "Возвращаемый код не равен ожидаемому")
		})
	}
}

func TestGetUserHistory(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		url  string
		name string
		body []byte
		want want
	}{
		{
			url:  "/history/1",
			name: "test correct user id",
			body: nil,
			want: want{
				code: http.StatusOK,
			},
		},
		{
			url:  "/history/-1",
			name: "test incorrect user id",
			body: nil,
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := config.GetConfig()
			if err != nil {
				log.Fatal("no config, can't start the program")
			}

			fmt.Println(config)

			h := NewHandler()
			s := storage.NewStorage(config.DatabaseURI)
			h.Storage = s
			r := CreateRouter(h)
			ts := httptest.NewServer(r)
			defer ts.Close()
			req, _ := testRequest(t, ts, http.MethodGet, tt.url, bytes.NewBuffer(tt.body))
			defer req.Body.Close()
			assert.Equal(t, tt.want.code, req.StatusCode, "Возвращаемый код не равен ожидаемому")
		})
	}

}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	fmt.Println(string(respBody))
	return resp, string(respBody)
}
