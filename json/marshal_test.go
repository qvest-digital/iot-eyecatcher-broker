package json

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"time"
)

func TestNewMarshaller(t *testing.T) {

	// given
	a := assert.New(t)

	// when
	testSubject := NewMarshaller()

	// then
	a.NotNil(testSubject.(*marshaller).jsonTemplate)
}

func TestMarshal(t *testing.T) {

	// given
	a := assert.New(t)
	testSubject := NewMarshaller()

	// when
	res, _ := testSubject.Marshal("myTopic", time.Unix(50, 0), []byte("\"hello\""))

	// then
	a.Equal([]byte(expectedJson), res)
}

func TestMarshalBadJson(t *testing.T) {

	// given
	a := assert.New(t)
	testSubject := NewMarshaller()
	testSubject.(*marshaller).jsonTemplate = ErroneousTemplate{}

	// when
	res, err := testSubject.Marshal("myTopic", time.Unix(50, 0), []byte("\"hello\""))

	// then
	a.Nil(res)
	a.Error(err, "some error")
}

const expectedJson = `{
    "topic": "myTopic",
    "timestamp": 50,
    "message": "hello"
}`

type ErroneousTemplate struct{}

func (e ErroneousTemplate) Execute(io.Writer, interface{}) error {
	return errors.New("some error")
}

// Just to please the coverage of the mocks
func TestMockPleaserMarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	x := NewMockMarshaller(ctrl)
	x.EXPECT().Marshal(gomock.Any(), gomock.Any(), gomock.Any())
	x.Marshal("", time.Now(), nil)
}
