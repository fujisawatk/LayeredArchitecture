package handler

import (
	"LayeredArchitecture/domain"
	"LayeredArchitecture/interfaces/middleware"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/takashabe/go-ddd-sample/domain/repository"
)

func Routes() *httprouter.Router {

	router := httprouter.New()

	// Post Route
	router.GET("/post/:id", middleware.Authenticate(postHandler.HandlePostGet))
	router.GET("/posts/index", middleware.Authenticate(postHandler.HandlePostsGet))
	router.POST("/post/create", middleware.Authenticate(postHandler.HandlePostCreate))
	router.PUT("/post/:id", middleware.Authenticate(postHandler.HandlePostUpdate))
	router.DELETE("/post/:id", middleware.Authenticate(postHandler.HandlePostDelete))

	return router
}

func prepareServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(Routes())
}

func sendRequest(t *testing.T, method, url string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	return res
}

func TestHandlePostGet(t *testing.T) {
	cases := []struct {
		input    int
		wantBody []byte
		wantCode int
		mock     func(*repository.MockUserRepository)
	}{
		{
			input:    1,
			wantBody: []byte(`{"id":1,"name":"foo"}`),
			wantCode: http.StatusOK,
			mock: func(r *repository.MockUserRepository) {
				r.EXPECT().Get(gomock.Any(), 1).Return(&domain.User{ID: 1, Name: "foo"}, nil)
			},
		},
		{
			input:    0,
			wantBody: nil,
			wantCode: http.StatusNotFound,
			mock: func(r *repository.MockUserRepository) {
				r.EXPECT().Get(gomock.Any(), 0).Return(nil, sql.ErrNoRows)
			},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := repository.NewMockUserRepository(ctrl)
			tt.mock(userRepo)

			h := &Handler{
				Repository: userRepo,
			}
			ts := prepareServer(t, h)
			defer ts.Close()

			res := sendRequest(t, "GET", fmt.Sprintf("%s/user/%d", ts.URL, tt.input), nil)
			defer res.Body.Close()

			if tt.wantCode != res.StatusCode {
				t.Errorf("want %d, got %d", tt.wantCode, res.StatusCode)
			}
			if res.StatusCode != http.StatusOK {
				return
			}

			payload, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("want non error, got %#v", err)
			}
			if diff := cmp.Diff(tt.wantBody, payload); diff != "" {
				t.Errorf("body mismatch %s", string(diff))
			}
		})
	}
}

func TestHandlePostsGet(t *testing.T) {

}

func TestHandlePostCreate(t *testing.T) {

}

func TestHandlePostUpdate(t *testing.T) {

}

func TestHandlePostDelete(t *testing.T) {

}
