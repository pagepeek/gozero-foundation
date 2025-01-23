package instagram

type AccessTokenResp struct {
	UserID      uint64   `json:"user_id"`
	AccessToken string   `json:"access_token"`
	Permissions []string `json:"permissions"`
	Code        int      `json:"code"`
	ErrType     string   `json:"error_type"`
	ErrMsg      string   `json:"error_message"`
}

type UserTokenResp struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Code        int    `json:"code"`
	ErrType     string `json:"error_type"`
	ErrMsg      string `json:"error_message"`
}

type UserData struct {
	ID       string `json:"id"`
	UserId   string `json:"user_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Code     int    `json:"code"`
	ErrType  string `json:"error_type"`
	ErrMsg   string `json:"error_message"`
}

type PublishResp struct {
	ID      string `json:"id"`
	Uri     string `json:"uri"`
	Code    int    `json:"code"`
	ErrType string `json:"error_type"`
	ErrMsg  string `json:"error_message"`
}

type UploadResp struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	DebugInfo struct {
		Retriable bool   `json:"retriable"`
		Type      string `json:"type"`
		Message   string `json:"message"`
	}
}

type UploadOption struct {
	MediaType      string
	AudioName      string
	Caption        string
	Collaborators  []string
	Children       []string
	CoverUrl       string
	IsCarouselItem bool
	LocationID     string
	ProductTags    []map[string]string
	ShareToFeed    bool
	ThumbOffset    int
	UploadType     string
	UserTags       []map[string]string
}
