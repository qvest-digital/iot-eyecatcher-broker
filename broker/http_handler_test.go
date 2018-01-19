package broker

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTopicHandler(t *testing.T) {

	// given
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// and
	brokerMock := NewMockBroker(ctrl)
	brokerMock.EXPECT().LastMessage("testTopic").Return([]byte("theLastMessage"), nil)

	// and custom test router
	r := mux.NewRouter()
	r.HandleFunc("/{topic}", func(w http.ResponseWriter, r *http.Request) {
		GetLastTopicMessageHandler(brokerMock, w, r)
	}).Methods(http.MethodGet)
	req, _ := http.NewRequest("GET", "/testTopic", nil)
	rec := httptest.NewRecorder()

	// when
	r.ServeHTTP(rec, req)
	res, _ := ioutil.ReadAll(rec.Body)

	// then
	a.Equal([]byte("theLastMessage"), res)
	a.Equal(200, rec.Code)
}

func TestTopicHandlerTopicNotFound(t *testing.T) {

	// given
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// and
	brokerMock := NewMockBroker(ctrl)
	brokerMock.EXPECT().LastMessage("testTopic").Return(nil, errors.New("topic not found"))

	// and custom test router
	r := mux.NewRouter()
	r.HandleFunc("/{topic}", func(w http.ResponseWriter, r *http.Request) {
		GetLastTopicMessageHandler(brokerMock, w, r)
	}).Methods(http.MethodGet)
	req, _ := http.NewRequest("GET", "/testTopic", nil)
	rec := httptest.NewRecorder()

	// when
	r.ServeHTTP(rec, req)
	res, _ := ioutil.ReadAll(rec.Body)

	// then
	a.Equal([]byte{}, res)
	a.Equal(404, rec.Code)
}

func TestTopicHandlerTopicError(t *testing.T) {

	// given
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// and
	brokerMock := NewMockBroker(ctrl)
	brokerMock.EXPECT().LastMessage("testTopic").Return(nil, errors.New("some error"))

	// and custom test router
	r := mux.NewRouter()
	r.HandleFunc("/{topic}", func(w http.ResponseWriter, r *http.Request) {
		GetLastTopicMessageHandler(brokerMock, w, r)
	}).Methods(http.MethodGet)
	req, _ := http.NewRequest("GET", "/testTopic", nil)
	rec := httptest.NewRecorder()

	// when
	r.ServeHTTP(rec, req)
	res, _ := ioutil.ReadAll(rec.Body)

	// then
	a.Equal([]byte{}, res)
	a.Equal(500, rec.Code)
}

func TestTopicHandlerNoTopicError(t *testing.T) {

	// given
	a := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// and custom test router
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		GetLastTopicMessageHandler(nil, w, r)
	}).Methods(http.MethodGet)
	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	// when
	r.ServeHTTP(rec, req)
	res, _ := ioutil.ReadAll(rec.Body)

	// then
	a.Equal([]byte{}, res)
	a.Equal(400, rec.Code)
}

/*
func TestMessageHandler(t *testing.T) {

    // given
    a := assert.New(t)
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // and
    brokerMock := NewMockBroker(ctrl)
    brokerMock.EXPECT().Message("testTopic", []byte("testMessage"))

    // and custom test router
    r := mux.NewRouter()
    r.HandleFunc("/{topic}", func(w http.ResponseWriter, r *http.Request) {
        MessageHandler(brokerMock, w, r)
    }).Methods(http.MethodPost)
    req, _ := http.NewRequest("POST", "/testTopic", strings.NewReader("testMessage"))
    rec := httptest.NewRecorder()

    // when
    r.ServeHTTP(rec, req)
    res, _ := ioutil.ReadAll(rec.Body)

    // then
    a.Equal([]byte{}, res)
    a.Equal(201, rec.Code)
}

func TestMessageHandlerNoTopicError(t *testing.T) {

    // given
    a := assert.New(t)
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // and custom test router
    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        MessageHandler(nil, w, r)
    }).Methods(http.MethodPost)
    req, _ := http.NewRequest("POST", "/", strings.NewReader("testMessage"))
    rec := httptest.NewRecorder()

    // when
    r.ServeHTTP(rec, req)
    res, _ := ioutil.ReadAll(rec.Body)

    // then
    a.Equal([]byte{}, res)
    a.Equal(400, rec.Code)
}

func TestMessageHandlerReadError(t *testing.T) {

    // given
    a := assert.New(t)
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // and custom test router
    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        MessageHandler(nil, w, r)
    }).Methods(http.MethodPost)
    req, _ := http.NewRequest("POST", "/", ErroneousReader{})
    rec := httptest.NewRecorder()

    // when
    r.ServeHTTP(rec, req)
    res, _ := ioutil.ReadAll(rec.Body)

    // then
    a.Equal([]byte{}, res)
    a.Equal(500, rec.Code)
}
*/

type ErroneousReader struct{}

func (e ErroneousReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("some error")
}
