package dto

type UserRegisterRequest struct {
	Name     string `json:"name" valid:"required~Name cannot be empty"`
	Email    string `json:"email" valid:"required~Email is required,email~Invalid email format"`
	Password string `json:"password" valid:"required~Password is required,minstringlength(8)~Password must be at least 8 characters long"`
}

type UserRegisterResponse struct {
	Result     string `json:"result"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type UserLoginRequest struct {
	Email    string `json:"email" valid:"required~Email is required,email~Invalid email format"`
	Password string `json:"password" valid:"required~Password is required"`
}

type UserLoginResponse struct {
	Result     string    `json:"result"`
	StatusCode int       `json:"statusCode"`
	Message    string    `json:"message"`
	Data       TokenData `json:"data"`
}

type TokenData struct {
	Token string `json:"token"`
}
