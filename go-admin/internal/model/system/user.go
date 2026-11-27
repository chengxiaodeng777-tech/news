package system

type CreateUserRequest struct {
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	RealName       string `json:"realName"`
	Nickname       string `json:"nickname"`
	Mobile         string `json:"mobile"`
	Email          string `json:"email"`
	Avatar         string `json:"avatar"`
	Status         int8   `json:"status"`
	OrganizationID int    `json:"organizationId"`
	Remark         string `json:"remark"`
}

type UpdateUserRequest struct {
	RealName       string `json:"realName"`
	Nickname       string `json:"nickname"`
	Mobile         string `json:"mobile"`
	Email          string `json:"email"`
	Avatar         string `json:"avatar"`
	Status         int8   `json:"status"`
	OrganizationID int    `json:"organizationId"`
	Remark         string `json:"remark"`
}

type AssignUserRolesRequest struct {
	RoleIDs []int `json:"roleIds" binding:"required"`
}

type ResetUserPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID               int    `json:"id"`
	Username         string `json:"username"`
	RealName         string `json:"realName"`
	Nickname         string `json:"nickname"`
	Mobile           string `json:"mobile"`
	Email            string `json:"email"`
	Avatar           string `json:"avatar"`
	Status           int8   `json:"status"`
	OrganizationID   int    `json:"organizationId"`
	OrganizationName string `json:"organizationName"`
	RoleIDs          []int  `json:"roleIds,omitempty"`
	Remark           string `json:"remark"`
}
