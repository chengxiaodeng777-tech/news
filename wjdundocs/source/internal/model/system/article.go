package system

type CreateArticleRequest struct {
	Sort     *int   `json:"sort" binding:"omitempty,gte=0,lte=2147483647"`
	Category string `json:"category" binding:"max=80"`
	Title    string `json:"title" binding:"required,max=120"`
	Summary  string `json:"summary" binding:"max=240"`
	Content  string `json:"content" binding:"required"`
	Cover    string `json:"cover"`
	Status   int8   `json:"status" binding:"oneof=0 1"`
	Author   string `json:"author" binding:"max=40"`
}
type UpdateArticleRequest = CreateArticleRequest
type ArticleListItem struct {
	Sort      int    `json:"sort"`
	Category  string `json:"category"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Cover     string `json:"cover"`
	Status    int8   `json:"status"`
	Author    string `json:"author"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ArticleResponse struct {
	Sort      int    `json:"sort"`
	Category  string `json:"category"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Content   string `json:"content"`
	Cover     string `json:"cover"`
	Status    int8   `json:"status"`
	Author    string `json:"author"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
