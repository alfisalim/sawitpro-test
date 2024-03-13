package handler

import (
	"errors"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/middlewares"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"strings"
	"testing"
)

type endpointsTestSuite struct {
	suite.Suite
	service generated.ServerInterface
}

func (e *endpointsTestSuite) TestLogin() {
	// Expectations
	ctrl := gomock.NewController(e.T())
	mockRepository := repository.NewMockRepositoryInterface(ctrl)
	mockValidator := middlewares.NewMockCustomValidatorInterface(ctrl)
	e.service = NewServer(NewServerOptions{
		Repository: mockRepository,
		Validator:  mockValidator,
	})

	e.Run("Negative Scenario, Failed Decode Body Req", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "12124", "Password": "test123", "Channel": "M"?}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		err := e.service.Login(newContext)
		e.Error(err)
	})

	e.Run("Negative Scenario, Invalid Body Req", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "123", "Password": "test123"}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("some error")).Times(1)

		err := e.service.Login(newContext)
		e.NoError(err)
	})

	e.Run("Negative Scenario, Failed Get Profile from DB", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "123", "Password": "test123"}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		mockRepository.EXPECT().GetProfile(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error")).Times(1)

		err := e.service.Login(newContext)
		e.NoError(err)
	})

	e.Run("Negative Scenario, User Not Found", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "123", "Password": "test123"}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		mockRepository.EXPECT().GetProfile(gomock.Any(), gomock.Any()).Return([]repository.Profile{}, nil).Times(1)

		err := e.service.Login(newContext)
		e.NoError(err)
	})

	resGetProfile := []repository.Profile{
		{
			UserId:    123,
			FullName:  "TEST",
			Password:  "TEST",
			Phone:     "TEST",
			Status:    1,
			CreatedAt: "TEST",
			UpdatedAt: "TEST",
		},
	}

	e.Run("Negative Scenario, Invalid Password", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "123", "Password": "test123"}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		mockRepository.EXPECT().GetProfile(gomock.Any(), gomock.Any()).Return(resGetProfile, nil).Times(1)

		err := e.service.Login(newContext)
		e.NoError(err)
	})

	e.Run("Negative Scenario, Failed get data login", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "123", "Password": "test123"}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		resGetProfileDum := resGetProfile
		resGetProfileDum[0].Password, _ = utils.HashPassword("test123")
		mockRepository.EXPECT().GetProfile(gomock.Any(), gomock.Any()).Return(resGetProfile, nil).Times(1)

		mockRepository.EXPECT().GetLogin(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error")).Times(1)

		err := e.service.Login(newContext)
		e.NoError(err)
	})

	e.Run("Negative Scenario, Failed insert login", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "123", "Password": "test123"}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		resGetProfileDum := resGetProfile
		resGetProfileDum[0].Password, _ = utils.HashPassword("test123")
		mockRepository.EXPECT().GetProfile(gomock.Any(), gomock.Any()).Return(resGetProfile, nil).Times(1)

		mockRepository.EXPECT().GetLogin(gomock.Any(), gomock.Any()).Return([]repository.LoginModel{}, nil).Times(1)

		mockRepository.EXPECT().InsertIntoLogin(gomock.Any(), gomock.Any()).Return(repository.LoginModel{}, errors.New("some error")).Times(1)

		err := e.service.Login(newContext)
		e.NoError(err)
	})

	resGetLogin := []repository.LoginModel{
		{
			LoginId:  123,
			UserId:   123,
			Ip:       "TEST",
			Token:    "TEST",
			Expires:  "TEST",
			Requests: 123,
		},
	}

	e.Run("Negative Scenario, Failed update login", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "123", "Password": "test123"}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		resGetProfileDum := resGetProfile
		resGetProfileDum[0].Password, _ = utils.HashPassword("test123")
		mockRepository.EXPECT().GetProfile(gomock.Any(), gomock.Any()).Return(resGetProfile, nil).Times(1)

		mockRepository.EXPECT().GetLogin(gomock.Any(), gomock.Any()).Return(resGetLogin, nil).Times(1)

		mockRepository.EXPECT().UpdateLogin(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("some error")).Times(1)

		err := e.service.Login(newContext)
		e.NoError(err)
	})

	resGetLoginTmp := []repository.LoginModel{
		{
			LoginId:  123,
			UserId:   123,
			Ip:       "TEST",
			Token:    "TEST",
			Expires:  "2020-12-12 00:00:00",
			Requests: 123,
		},
	}

	e.Run("Positive Scenario", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "123", "Password": "test123"}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		resGetProfileDum := resGetProfile
		resGetProfileDum[0].Password, _ = utils.HashPassword("test123")
		mockRepository.EXPECT().GetProfile(gomock.Any(), gomock.Any()).Return(resGetProfile, nil).Times(1)

		mockRepository.EXPECT().GetLogin(gomock.Any(), gomock.Any()).Return(resGetLoginTmp, nil).Times(1)

		mockRepository.EXPECT().UpdateLogin(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := e.service.Login(newContext)
		e.NoError(err)
	})
}

func (e *endpointsTestSuite) TestGetProfile() {
	// Expectations
	ctrl := gomock.NewController(e.T())
	mockRepository := repository.NewMockRepositoryInterface(ctrl)
	e.service = NewServer(NewServerOptions{
		Repository: mockRepository,
	})

	reqDum := httptest.NewRequest(echo.GET, "http://localhost:1323/profile", nil)
	reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Act
	c := echo.New()
	newContext := c.NewContext(reqDum, rec)

	e.Run("Negative Scenario, Failed Validate Token", func() {
		req := &generated.GetProfileParams{
			Authorization: "Bearer abcdefg",
		}
		err := e.service.GetProfile(newContext, *req)
		e.NoError(err)
	})

	e.Run("Postitive Scenario, Success", func() {
		req := &generated.GetProfileParams{
			Authorization: "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjMsIkZ1bGxOYW1lIjoiQWxmaSBTYWxpbSIsIlBhc3N3b3JkIjoiIiwiUGhvbmUiOiIrNjI4MTIzMTEyNjIiLCJTdGF0dXMiOjAsIkNyZWF0ZWRBdCI6IiIsIlVwZGF0ZWRBdCI6IiIsImV4cCI6MTcxMDUxNTAzMn0.FTvOH-mmLb86KmMwTL-2bg2HcqOW3DEkUiwNF0eurg7ZFserq-MRgG9Bmx6SOBy8jrbqfDu1xUjTcQmgAE_w1IloznWbbrlVRNZVjHOemYMvucz93LtfFDr16Yl6smz8q5aK6NilA7-bLh-aKs16Dl0TGBr6TpPNE41pl7mE_bY",
		}
		err := e.service.GetProfile(newContext, *req)
		e.NoError(err)
	})
}

func (e *endpointsTestSuite) TestUpdateProfile() {
	// Expectations
	ctrl := gomock.NewController(e.T())
	mockRepository := repository.NewMockRepositoryInterface(ctrl)
	mockValidator := middlewares.NewMockCustomValidatorInterface(ctrl)
	e.service = NewServer(NewServerOptions{
		Repository: mockRepository,
		Validator:  mockValidator,
	})

	req := generated.UpdateProfileParams{
		Authorization: "Bearer abcdefg",
	}

	e.Run("Negative Scenario, Failed Decode Body Req", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "12124", "fullName": "test123"?}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		err := e.service.UpdateProfile(newContext, req)
		e.Error(err)
	})

	e.Run("Negative Scenario, Invalid Body Req", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "12124"}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("some error")).Times(1)

		err := e.service.UpdateProfile(newContext, req)
		e.NoError(err)
	})

	e.Run("Negative Scenario, Invalid token", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "12124", "fullName": "test123"}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		err := e.service.UpdateProfile(newContext, req)
		e.NoError(err)
	})

	req = generated.UpdateProfileParams{
		Authorization: "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjMsIkZ1bGxOYW1lIjoiQWxmaSBTYWxpbSIsIlBhc3N3b3JkIjoiIiwiUGhvbmUiOiIrNjI4MTIzMTEyNjIiLCJTdGF0dXMiOjAsIkNyZWF0ZWRBdCI6IiIsIlVwZGF0ZWRBdCI6IiIsImV4cCI6MTcxMDUxNTAzMn0.FTvOH-mmLb86KmMwTL-2bg2HcqOW3DEkUiwNF0eurg7ZFserq-MRgG9Bmx6SOBy8jrbqfDu1xUjTcQmgAE_w1IloznWbbrlVRNZVjHOemYMvucz93LtfFDr16Yl6smz8q5aK6NilA7-bLh-aKs16Dl0TGBr6TpPNE41pl7mE_bY",
	}

	e.Run("Negative Scenario, Failed update profile", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "12124", "fullName": "test123"}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		mockRepository.EXPECT().UpdateProfile(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("some error SQLSTATE 23505"))

		err := e.service.UpdateProfile(newContext, req)
		e.NoError(err)
	})

	e.Run("Positive Scenario, Return success", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "12124", "fullName": "test123"}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		mockRepository.EXPECT().UpdateProfile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		err := e.service.UpdateProfile(newContext, req)
		e.NoError(err)
	})
}

func (e *endpointsTestSuite) TestRegister() {
	// Expectations
	ctrl := gomock.NewController(e.T())
	mockRepository := repository.NewMockRepositoryInterface(ctrl)
	mockValidator := middlewares.NewMockCustomValidatorInterface(ctrl)
	e.service = NewServer(NewServerOptions{
		Repository: mockRepository,
		Validator:  mockValidator,
	})

	e.Run("Negative Scenario, Failed Decode Body Req", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "12124", "fullName": "test123"?}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		err := e.service.Register(newContext)
		e.Error(err)
	})

	e.Run("Negative Scenario, Invalid Body Req", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "12124"}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("some error")).Times(1)

		err := e.service.Register(newContext)
		e.NoError(err)
	})

	e.Run("Negative Scenario, Faield Get Profile", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "12124", "fullName": "test123", "password": "<PASSWORD>"}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		mockRepository.EXPECT().GetProfile(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))

		err := e.service.Register(newContext)
		e.NoError(err)
	})

	resGetProfile := []repository.Profile{
		{
			UserId:    123,
			FullName:  "TEST",
			Password:  "TEST",
			Phone:     "TEST",
			Status:    123,
			CreatedAt: "TEST",
			UpdatedAt: "TEST",
		},
	}
	e.Run("Negative Scenario, Phone Number Already Exists", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "12124", "fullName": "test123", "password": "<PASSWORD>"}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		mockRepository.EXPECT().GetProfile(gomock.Any(), gomock.Any()).Return(resGetProfile, nil)

		err := e.service.Register(newContext)
		e.NoError(err)
	})

	e.Run("Negative Scenario, Failed Insert New User Into DB", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "12124", "fullName": "test123", "password": "<PASSWORD>"}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		mockRepository.EXPECT().GetProfile(gomock.Any(), gomock.Any()).Return(nil, nil)

		mockRepository.EXPECT().CreateProfile(gomock.Any(), gomock.Any()).Return(repository.Profile{}, errors.New("some error"))

		err := e.service.Register(newContext)
		e.NoError(err)
	})

	e.Run("Positive Scenario, Should return data", func() {
		bodyReader := strings.NewReader(`{"phoneNumber": "12124", "fullName": "test123", "password": "<PASSWORD>"}`)
		reqDum := httptest.NewRequest(echo.POST, "http://localhost:1323/login", bodyReader)
		reqDum.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New()
		newContext := c.NewContext(reqDum, rec)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		mockRepository.EXPECT().GetProfile(gomock.Any(), gomock.Any()).Return(nil, nil)

		mockRepository.EXPECT().CreateProfile(gomock.Any(), gomock.Any()).Return(repository.Profile{}, nil)

		err := e.service.Register(newContext)
		e.NoError(err)
	})
}

func TestEndpoints(t *testing.T) {
	suite.Run(t, new(endpointsTestSuite))
}
