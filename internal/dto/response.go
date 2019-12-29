package dto

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type GenerateOtpResponse struct {
	Response
}

func NewGenerateOtpResponse(status int, message string) *GenerateOtpResponse {
	return &GenerateOtpResponse{Response: Response{
		Status:  status,
		Message: message,
	}}
}

type LoginResponse struct {
	Response
	Token string `json:"token"`
}

func NewLoginResponse(status int, message string, token string) *LoginResponse {
	return &LoginResponse{
		Response: Response{
			Status:  status,
			Message: message,
		},
		Token: token,
	}
}
