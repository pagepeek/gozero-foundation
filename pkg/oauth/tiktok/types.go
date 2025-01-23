package tiktok

type AccessTokenResp struct {
	OpenId           string `json:"open_id"`
	Scope            string `json:"scope"`
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
	Error            string `json:"error"`
	ErrDesc          string `json:"error_description"`
}

type UserData struct {
	Data struct {
		User struct {
			AvatarUrl string `json:"avatar_url"`
			OpenId    string `json:"open_id"`
			UnionId   string `json:"union_id"`
			Name      string `json:"display_name"`
		}
	}
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		LogId   string `json:"log_id"`
	}
}
