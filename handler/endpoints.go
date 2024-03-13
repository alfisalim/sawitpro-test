package handler

import (
	"encoding/json"
	"fmt"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"strings"
	"time"
)

func (s *Server) Login(ctx echo.Context) error {
	var req *generated.LoginRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&req)
	if err != nil {
		return err
	}

	err = s.Validator.Validate(req)
	if err != nil {
		return ctx.JSON(400, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	resGetProfile, err := s.Repository.GetProfile(ctx.Request().Context(), map[string]interface{}{
		"phone": req.PhoneNumber,
	})
	if err != nil {
		return ctx.JSON(500, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	if len(resGetProfile) == 0 {
		return ctx.JSON(400, generated.ErrorResponse{
			Message: fmt.Sprintf("User with phone number %s not found", req.PhoneNumber),
		})
	}

	if !utils.CheckPasswordHash(req.Password, resGetProfile[0].Password) {
		return ctx.JSON(400, generated.ErrorResponse{
			Message: fmt.Sprintf("Invalid password"),
		})
	}

	filterGetLoginData := map[string]interface{}{
		"user_id": resGetProfile[0].UserId,
	}

	resGetLogin, err := s.Repository.GetLogin(ctx.Request().Context(), filterGetLoginData)
	if err != nil {
		return ctx.JSON(500, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	jwtToken, expiresAt, _ := utils.GenerateToken(resGetProfile[0])

	if len(resGetLogin) == 0 {
		_, err = s.Repository.InsertIntoLogin(ctx.Request().Context(), repository.LoginModel{
			UserId:   resGetProfile[0].UserId,
			Ip:       ctx.Request().RemoteAddr,
			Token:    jwtToken,
			Expires:  expiresAt,
			Requests: 0,
		})
	} else {
		updatedData := map[string]interface{}{
			"ip":         ctx.Request().RemoteAddr,
			"requests":   resGetLogin[0].Requests + 1,
			"updated_at": time.Now(),
		}

		if resGetLogin[0].Expires < time.Now().Format("2006-01-02 15:04:05") {
			updatedData["token"] = jwtToken
			updatedData["expires"] = expiresAt
		} else {
			jwtToken = resGetLogin[0].Token
		}

		err = s.Repository.UpdateLogin(ctx.Request().Context(), filterGetLoginData, updatedData)
	}

	if err != nil {
		return ctx.JSON(500, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return ctx.JSON(200, generated.LoginResponse{
		Message: "success",
		Token:   jwtToken,
		UserID:  int(resGetProfile[0].UserId),
	})
}

func (s *Server) GetProfile(ctx echo.Context, params generated.GetProfileParams) error {
	claims, err := utils.ValidateToken(params.Authorization)
	if err != nil {
		return ctx.JSON(403, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	mapClaims := claims.(jwt.MapClaims)
	return ctx.JSON(200, generated.ProfileResponse{
		FullName:    mapClaims["FullName"].(string),
		Message:     "success",
		PhoneNumber: mapClaims["Phone"].(string),
	})
}

func (s *Server) UpdateProfile(ctx echo.Context, params generated.UpdateProfileParams) error {
	var req *generated.UpdateProfileRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&req)
	if err != nil {
		return err
	}

	err = s.Validator.Validate(req)
	if err != nil {
		return ctx.JSON(400, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	claims, err := utils.ValidateToken(params.Authorization)
	if err != nil {
		return ctx.JSON(403, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	mapClaims := claims.(jwt.MapClaims)
	var userID float64
	if _, ok := mapClaims["UserId"]; ok {
		userID = mapClaims["UserId"].(float64)
	}

	updatedBy := map[string]interface{}{
		"user_id": userID,
	}

	updatedData := map[string]interface{}{
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	}

	if req.FullName != nil {
		updatedData["full_name"] = req.FullName
	}

	if req.PhoneNumber != nil {
		updatedData["phone"] = req.PhoneNumber
	}

	err = s.Repository.UpdateProfile(ctx.Request().Context(), updatedBy, updatedData)
	if err != nil {
		code := 500
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			code = 409
		}
		return ctx.JSON(code, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return ctx.JSON(200, generated.ErrorResponse{
		Message: "success",
	})
}

func (s *Server) Register(ctx echo.Context) error {
	var req *generated.RegisterRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&req)
	if err != nil {
		return err
	}

	err = s.Validator.Validate(req)
	if err != nil {
		return ctx.JSON(400, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	resGetProfile, err := s.Repository.GetProfile(ctx.Request().Context(), map[string]interface{}{
		"phone": req.PhoneNumber,
	})
	if err != nil {
		return ctx.JSON(400, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	if len(resGetProfile) > 0 {
		return ctx.JSON(409, generated.ErrorResponse{
			Message: "phone number already exists",
		})
	}

	hashPassword, _ := utils.HashPassword(req.Password)

	now := time.Now().Format("2006-01-02 15:04:05")
	resCreateProfile, err := s.Repository.CreateProfile(ctx.Request().Context(), repository.Profile{
		FullName:  req.FullName,
		Password:  hashPassword,
		Phone:     req.PhoneNumber,
		Status:    1,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return ctx.JSON(500, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return ctx.JSON(200, generated.RegisterResponse{
		Message: "success",
		UserID:  int(resCreateProfile.UserId),
	})
}
