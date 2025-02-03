package handlers

type newUserForm struct {
	Name            string `json:"name" binding:"required"`
	Surname         string `json:"surname" binding:"required"`
	Email           string `json:"email" binding:"required"`
	PhoneNumber     string `json:"phone_number" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}
type userLoginForm struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userInfoUpdateFrom struct {
	Name    string `json:"name" binding:"required"`    
	Surname    string `json:"surname" binding:"required"`    
	PhoneNumber string `json:"phone_number" binding:"required"` 
}
