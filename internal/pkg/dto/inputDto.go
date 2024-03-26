package dto

type (
	PermissionsInputDTO struct {
		UserModule    *bool `json:"user_module" example:"true"`
		ProfileModule *bool `json:"profile_module" example:"true"`
	}

	ProfileInputDTO struct {
		Name        *string             `json:"name" example:"ADMIN"`
		Permissions PermissionsInputDTO `json:"permissions"`
	}

	UserInputDTO struct {
		Name      *string `json:"name" example:"John Cena"`
		Email     *string `json:"email" example:"john.cena@email.com"`
		Status    *bool   `json:"status" example:"true"`
		ProfileID *uint   `json:"profile_id" example:"1"`
	}

	PasswordInputDTO struct {
		Password        *string `json:"password" example:"secret"`
		PasswordConfirm *string `json:"password_confirm" example:"secret"`
	}

	LoginInputDTO struct {
		Email    string `json:"email" example:"admin@admin.com"`
		Password string `json:"password" example:"12345678"`
		Expire   bool   `json:"expire" example:"false"`
	}

	// Accounts example
	OperatorInputDTO struct {
		Name *string `json:"name" example:"TIM"`
	}

	PhoneInputDTO struct {
		Number     *string `json:"number" example:"TIM"`
		OperatorID *uint   `json:"operator_id" example:"1"`
	}

	SiteInputDTO struct {
		Name *string `json:"name" example:"Google"`
		URL  *string `json:"url" example:"https://www.google.com/"`
	}

	AccountInputDTO struct {
		Username *string `json:"username" example:"username"`
		Password *string `json:"password" example:"secret"`
		SiteID   *uint   `json:"site_id" example:"1"`
		PhoneID  *uint   `json:"phone_id" example:"1"`
		MailID   *uint   `json:"mail_id" example:"1"`
	}
)

func (p PasswordInputDTO) IsValid() bool {
	if p.Password == nil || p.PasswordConfirm == nil {
		return false
	}

	if len(*p.Password) < 5 || len(*p.PasswordConfirm) < 5 {
		return false
	}

	return *p.Password == *p.PasswordConfirm
}
