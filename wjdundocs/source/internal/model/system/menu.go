package system

type CreateMenuRequest struct {
	ParentID       int    `json:"parentId"`
	Name           string `json:"name" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Path           string `json:"path"`
	Component      string `json:"component"`
	Icon           string `json:"icon"`
	Type           string `json:"type" binding:"required"`
	PermissionCode string `json:"permissionCode"`
	Redirect       string `json:"redirect"`
	Sort           int    `json:"sort"`
	Status         int8   `json:"status"`
	Hidden         bool   `json:"hidden"`
	KeepAlive      bool   `json:"keepAlive"`
	ExternalLink   string `json:"externalLink"`
	Remark         string `json:"remark"`
}

type UpdateMenuRequest = CreateMenuRequest

type MenuResponse struct {
	ID             int            `json:"id"`
	ParentID       int            `json:"parentId"`
	Name           string         `json:"name"`
	Title          string         `json:"title"`
	Path           string         `json:"path"`
	Component      string         `json:"component"`
	Icon           string         `json:"icon"`
	Type           string         `json:"type"`
	PermissionCode string         `json:"permissionCode"`
	Redirect       string         `json:"redirect"`
	Sort           int            `json:"sort"`
	Status         int8           `json:"status"`
	Hidden         bool           `json:"hidden"`
	KeepAlive      bool           `json:"keepAlive"`
	ExternalLink   string         `json:"externalLink"`
	Remark         string         `json:"remark"`
	Children       []MenuResponse `json:"children,omitempty"`
}
