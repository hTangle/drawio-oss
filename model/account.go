package model

type Account struct {
	Email    string `json:"email"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (a *Account) CheckAccount() bool {
	return GetLocalEditorConf("").IsValidUser(a.Email, a.Password)
}
