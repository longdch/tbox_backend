package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tbox_backend/internal/constants"
	"tbox_backend/internal/dto"
	"tbox_backend/internal/helpers"
	"tbox_backend/internal/services"
	"tbox_backend/internal/validator"
)

const OtpRequestKey = "OtpRequest"

type Router struct {
	userService        services.IUserService
	userValidator      validator.IUserValidator
	phoneNumberLimiter *helpers.PhoneNumberRateLimiters
}

func NewRouter(userService services.IUserService, userValidator validator.IUserValidator, phoneNumberLimiter *helpers.PhoneNumberRateLimiters) *Router {
	return &Router{
		userService:        userService,
		userValidator:      userValidator,
		phoneNumberLimiter: phoneNumberLimiter,
	}
}

func (r *Router) IndexRouter(rg *gin.Engine) {
	gr := rg.Group("/api")
	{
		gr.POST("/generate_otp", r.rateLimit, r.generateOtpHandler)
		gr.POST("/resend_otp", r.rateLimit, r.resendOtpHandler)
		gr.POST("/login", r.loginHandler)
	}
}

// @Summary Generate otp
// @Description Generate otp and send otp to phone number. OTP will be printed in console log.
// @Accept json
// @Produce json
// @Param Body body dto.GenerateOtpRequest true "Body"
// @Success 200 {object} dto.GenerateOtpResponse
// @Router /generate_otp [post]
func (r *Router) generateOtpHandler(ctx *gin.Context) {
	generateOtpRequest := ctx.MustGet(OtpRequestKey)
	err := r.userService.GenerateOtp(generateOtpRequest.(dto.GenerateOtpRequest).PhoneNumber)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, dto.NewGenerateOtpResponse(constants.SomethingWentWrongStatus, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewGenerateOtpResponse(constants.SuccessStatus, "Success"))
	return
}

// @Summary Resend otp
// @Description Generate new otp and send otp to phone number. OTP will be printed in console log.
// @Accept json
// @Produce json
// @Param Body body dto.GenerateOtpRequest true "Body"
// @Success 200 {object} dto.GenerateOtpResponse
// @Router /resend_otp [post]
func (r *Router) resendOtpHandler(ctx *gin.Context) {
	generateOtpRequest := ctx.MustGet(OtpRequestKey)
	err := r.userService.ResendOtp(generateOtpRequest.(dto.GenerateOtpRequest).PhoneNumber)
	if err != nil {
		ctx.JSON(http.StatusOK, dto.NewGenerateOtpResponse(constants.SomethingWentWrongStatus, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewGenerateOtpResponse(constants.SuccessStatus, "Success"))
	return
}

// @Summary Login
// @Description Verify otp and return access_token. OTP is not necessary after the fist successful login
// @Accept  json
// @Produce  json
// @Param Body body dto.LoginRequest true "Body"
// @Success 200 {object} dto.LoginResponse
// @Router /login [post]
func (r *Router) loginHandler(ctx *gin.Context) {
	var loginRequest dto.LoginRequest
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusOK, dto.NewLoginResponse(constants.InvalidRequestStatus, err.Error(), ""))
		return
	}

	token, err := r.userService.Login(loginRequest.PhoneNumber, loginRequest.Otp)
	if err != nil {
		ctx.JSON(http.StatusOK, dto.NewLoginResponse(constants.SomethingWentWrongStatus, err.Error(), token))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewLoginResponse(constants.SuccessStatus, "Success", token))
	return
}

func (r *Router) rateLimit(ctx *gin.Context) {
	var generateOtpRequest dto.GenerateOtpRequest
	if err := ctx.ShouldBindJSON(&generateOtpRequest); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, dto.NewGenerateOtpResponse(constants.InvalidRequestStatus, err.Error()))
		return
	}

	if valid := r.userValidator.IsPhoneNumberValid(generateOtpRequest.PhoneNumber); !valid {
		ctx.AbortWithStatusJSON(http.StatusOK, dto.NewGenerateOtpResponse(constants.InvalidRequestStatus, "Phone number invalid "))
		return
	}

	limiter := r.phoneNumberLimiter.GetLimiter(generateOtpRequest.PhoneNumber)
	if !limiter.Allow() {
		ctx.AbortWithStatusJSON(http.StatusOK, dto.NewGenerateOtpResponse(constants.TooManyRequestStatus, "Too many requests "))
		return
	}

	ctx.Set(OtpRequestKey, generateOtpRequest)
	return
}

