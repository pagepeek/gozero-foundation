package twitter

type AccessTokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Error        string `json:"error"`
	ErrDesc      string `json:"error_description"`
}

type UserData struct {
	Data struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
	}
	Status int    `json:"status"`
	Title  string `json:"title"`
	Error  string `json:"detail"`
}
