package typ

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"`
}

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Pwhash   string `json:"pwhash,omitempty"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}
