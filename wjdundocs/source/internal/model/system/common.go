package system

type IDPath struct {
	ID int `uri:"id" binding:"required"`
}

type PageQuery struct {
	Category string `form:"category"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Keyword  string `form:"keyword"`
	Status   *int8  `form:"status"`
}

type PageResult struct {
	List     any `json:"list"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type IDListResponse struct {
	IDs []int `json:"ids"`
}

type DictDataQuery struct {
	DictTypeID int `form:"dictTypeId"`
}
