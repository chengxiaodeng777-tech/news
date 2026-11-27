package system

type CreateRoleRequest struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Status int8   `json:"status"`
	Sort   int    `json:"sort"`
	Remark string `json:"remark"`
}

type UpdateRoleRequest struct {
	Name   string `json:"name" binding:"required"`
	Status int8   `json:"status"`
	Sort   int    `json:"sort"`
	Remark string `json:"remark"`
}

type AssignRoleMenusRequest struct {
	MenuIDs []int `json:"menuIds" binding:"required"`
}

type RoleResponse struct {
	ID      int    `json:"id"`
	Code    string `json:"code"`
	Name    string `json:"name"`
	Status  int8   `json:"status"`
	Sort    int    `json:"sort"`
	MenuIDs []int  `json:"menuIds,omitempty"`
	Remark  string `json:"remark"`
}
