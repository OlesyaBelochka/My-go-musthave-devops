package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleUpdateMetrics(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{
		{"count", "/update/counter/PollCount/1",
			want{http.StatusOK, "text/plain"},
		},
		{"gauge", "/update/gauge/HeapInuse/126197760",
			want{http.StatusOK, "text/plain"},
		},
		{"bad", "/update/wewfsf/HeapInuse",
			want{http.StatusNotImplemented, "text/plain"},
		},
		{"noneGauge", "/update/gauge/HeapInuse/none",
			want{http.StatusBadRequest, "text/plain"},
		},
		{"testCounter100", "/update/counter/testCounter/100",
			want{http.StatusOK, "text/plain"},
		},
		{"noneCounter", "/update/counter/testCounter/non",
			want{http.StatusBadRequest, "text/plain"},
		},
	}
	fmt.Println("start..")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(HandleUpdateMetrics)
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, result.StatusCode, tt.want.code)
			assert.Contains(t, tt.want.contentType, result.Header.Get("Content-Type"))

		})
	}
}
