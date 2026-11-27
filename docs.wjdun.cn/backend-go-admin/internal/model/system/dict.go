package system

type CreateDictTypeRequest struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Status int8   `json:"status"`
	Remark string `json:"remark"`
}

type UpdateDictTypeRequest = CreateDictTypeRequest

type CreateDictDataRequest struct {
	DictTypeID int    `json:"dictTypeId" binding:"required"`
	Label      string `json:"label" binding:"required"`
	Value      string `json:"value" binding:"required"`
	TagType    string `json:"tagType"`
	CSSClass   string `json:"cssClass"`
	Sort       int    `json:"sort"`
	Status     int8   `json:"status"`
	IsDefault  bool   `json:"isDefault"`
	Remark     string `json:"remark"`
}

type UpdateDictDataRequest = CreateDictDataRequest

type DictTypeResponse struct {
	ID     int    `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status int8   `json:"status"`
	Remark string `json:"remark"`
}

type DictDataResponse struct {
	ID         int    `json:"id"`
	DictTypeID int    `json:"dictTypeId"`
	Label      string `json:"label"`
	Value      string `json:"value"`
	TagType    string `json:"tagType"`
	CSSClass   string `json:"cssClass"`
	Sort       int    `json:"sort"`
	Status     int8   `json:"status"`
	IsDefault  bool   `json:"isDefault"`
	Remark     string `json:"remark"`
}
