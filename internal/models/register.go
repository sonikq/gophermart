package models

type RegisterUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (m RegisterUserRequest) Validate() error {
	if m.Login == "" {
		return ErrEmptyLogin
	}

	if m.Password == "" {
		return ErrEmptyPassword
	}

	return nil
}
