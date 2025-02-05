package models

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (m LoginRequest) Validate() error {
	if m.Login == "" {
		return ErrEmptyLogin
	}

	if m.Password == "" {
		return ErrEmptyPassword
	}

	return nil
}
