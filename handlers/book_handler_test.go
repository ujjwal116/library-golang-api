package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"library/mocks"
	"library/model"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Implement the necessary methods of *gorm.DB in the MockDB struct

type BookHandlerTestSuite struct {
	suite.Suite
	ctx               *gin.Context
	responseRecoreder *httptest.ResponseRecorder
	repoMock          *mocks.BookRepoMock
	bh                *BookHandler
}

func TestBookHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(BookHandlerTestSuite))
}
func (s *BookHandlerTestSuite) SetupTest() {
	s.repoMock = new(mocks.BookRepoMock)
	s.responseRecoreder = httptest.NewRecorder()
	s.ctx = PrepareGinContext(s.responseRecoreder)
	s.bh = &BookHandler{s.repoMock}
}

func (suite BookHandlerTestSuite) TestAddBook_success() {
	params := make(map[string]string)
	prepareContextForRequest(suite.ctx, "POST", &params, &model.Book{Id: 1, Name: "dummy1", Author: "dummy_author"})
	suite.repoMock.On("AddBook", mock.AnythingOfType("*model.Book")).Return(nil)
	suite.bh.AddBook(suite.ctx)
	actResp:=&model.Book{}
	json.Unmarshal(suite.responseRecoreder.Body.Bytes(), actResp)
	assert.Equal(suite.Suite.T(), &model.Book{Id: 1, Name: "dummy1", Author: "dummy_author"}, actResp)
}

func (suite BookHandlerTestSuite) TestAddBook_already_exists() {
	params := make(map[string]string)
	prepareContextForRequest(suite.ctx, "POST", &params, &model.Book{Id: 1, Name: "dummy1", Author: "dummy_author"})
	suite.repoMock.On("AddBook", mock.AnythingOfType("*model.Book")).Return(errors.New("Already exists"))
	suite.bh.AddBook(suite.ctx)
	actResp := &ErrorResponse{}
	json.Unmarshal(suite.responseRecoreder.Body.Bytes(), actResp)
	assert.Equal(suite.Suite.T(), "Already exists", actResp.Error)
}
func (suite *BookHandlerTestSuite) TestBookHandler_GetBooks_db_error_count() {
	params := make(map[string]string)
	prepareContextForRequest(suite.ctx, "GET", &params, nil)
	suite.repoMock.On("GetBooksCount", mock.AnythingOfType("*int64")).Return(errors.New("error occured"))
	suite.bh.GetBooks(suite.ctx)
	actResp := &ErrorResponse{}
	json.Unmarshal(suite.responseRecoreder.Body.Bytes(), actResp)
	assert.Equal(suite.Suite.T(), "error occured", actResp.Error)
}

func (suite *BookHandlerTestSuite) TestBookHandler_GetBooks_db_error_fetch_result() {

	params := make(map[string]string)
	prepareContextForRequest(suite.ctx, "GET", &params, nil)

	suite.repoMock.On("GetBooksCount", mock.AnythingOfType("*int64")).Return(nil).Run(func(args mock.Arguments) {
		v := args.Get(0).(*int64)
		*v = 9
	})
	suite.repoMock.On("GetBooksByPageAndPageSize",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("*[]model.Book")).Return(errors.New("error occured while fetching page"))

	suite.bh.GetBooks(suite.ctx)
	actResp := &ErrorResponse{}
	json.Unmarshal(suite.responseRecoreder.Body.Bytes(), actResp)
	assert.Equal(suite.Suite.T(), "error occured while fetching page", actResp.Error)

}
func (suite *BookHandlerTestSuite) TestBookHandler_GetBooks_success() {
	params := make(map[string]string)
	params["page"] = "1"
	prepareContextForRequest(suite.ctx, "GET", &params, nil)
	expactedResponse := []model.Book{{Id: 1, Name: "dummy1", Author: "dummy_author"}}
	fmt.Println(expactedResponse)
	suite.repoMock.On("GetBooksCount", mock.AnythingOfType("*int64")).Return(nil).Run(func(args mock.Arguments) {
		v := args.Get(0).(*int64)
		*v = int64(len(expactedResponse))
	})
	suite.repoMock.On("GetBooksByPageAndPageSize", mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("*[]model.Book")).Return(nil).Run(func(args mock.Arguments) {
		books := args.Get(2).(*[]model.Book)
		*books = []model.Book{{Id: 1, Name: "dummy1", Author: "dummy_author"}}
	})
	suite.bh.GetBooks(suite.ctx)
	actResp := &PagingResult{}
	json.Unmarshal(suite.responseRecoreder.Body.Bytes(), actResp)
	assert.Equal(suite.Suite.T(), expactedResponse, *actResp.Data)

}

func (suite *BookHandlerTestSuite) TestBookHandler_GetBooks_page_not_found() {
	params := make(map[string]string)
	params["page"] = "3"
	params["pageSize"] = "1"
	prepareContextForRequest(suite.ctx, "GET", &params, nil)
	expactedResponse := []model.Book{{Id: 1, Name: "dummy1", Author: "dummy_author"}, {Id: 2, Name: "dummy2", Author: "dummy_author2"}}
	fmt.Println(expactedResponse)
	suite.repoMock.On("GetBooksCount", mock.AnythingOfType("*int64")).Return(nil).Run(func(args mock.Arguments) {
		v := args.Get(0).(*int64)
		*v = 2
	})
	suite.bh.GetBooks(suite.ctx)
	actResp := &ErrorResponse{}
	json.Unmarshal(suite.responseRecoreder.Body.Bytes(), actResp)
	assert.Equal(suite.Suite.T(), "Page not found", actResp.Error)

}

func prepareContextForRequest(ctx *gin.Context, reqType string, params *map[string]string, body *model.Book) {
	ctx.Request.Method = reqType
	u := url.Values{}
	for k, v := range *params {
		u.Add(k, v)
	}
	ctx.Request.URL.RawQuery = u.Encode()
	if nil != body {
		jsonbytes, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
	}

}

func PrepareGinContext(rr *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = &http.Request{Header: make(http.Header), URL: &url.URL{}}
	ctx.Request.Header.Set("Content-Type", "application/json")
	return ctx
}
