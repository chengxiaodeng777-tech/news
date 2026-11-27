package auth

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type UserInfoResponse struct {
	UserID       int      `json:"userId"`
	Username     string   `json:"username"`
	RealName     string   `json:"realName"`
	Nickname     string   `json:"nickname"`
	Avatar       string   `json:"avatar"`
	HomePath     string   `json:"homePath"`
	Roles        []string `json:"roles"`
	Organization string   `json:"organization"`
}

type MenuMeta struct {
	Title  string `json:"title"`
	Icon   string `json:"icon"`
	Hidden bool   `json:"hidden"`
}

type MenuResponse struct {
	Name      string         `json:"name"`
	Path      string         `json:"path"`
	Component string         `json:"component"`
	Redirect  string         `json:"redirect"`
	Meta      MenuMeta       `json:"meta"`
	Children  []MenuResponse `json:"children,omitempty"`
}
