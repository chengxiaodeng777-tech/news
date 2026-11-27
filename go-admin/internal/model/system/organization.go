package system

type CreateOrganizationRequest struct {
	ParentID int    `json:"parentId"`
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Type     string `json:"type"`
	Leader   string `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Sort     int    `json:"sort"`
	Status   int8   `json:"status"`
	Remark   string `json:"remark"`
}

type UpdateOrganizationRequest = CreateOrganizationRequest

type OrganizationResponse struct {
	ID       int                    `json:"id"`
	ParentID int                    `json:"parentId"`
	Name     string                 `json:"name"`
	Code     string                 `json:"code"`
	Type     string                 `json:"type"`
	Leader   string                 `json:"leader"`
	Phone    string                 `json:"phone"`
	Email    string                 `json:"email"`
	Sort     int                    `json:"sort"`
	Status   int8                   `json:"status"`
	Remark   string                 `json:"remark"`
	Children []OrganizationResponse `json:"children,omitempty"`
}
