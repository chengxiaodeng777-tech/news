package system

import (
	"context"
	"strings"
	"time"

	"base/ent"
	"base/ent/article"
	"entgo.io/ent/dialect/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"

	systemmodel "base/internal/model/system"
	"base/internal/pkg/apperror"

	"golang.org/x/net/html"
)

func articleResponse(v *ent.Article) systemmodel.ArticleResponse {
	return systemmodel.ArticleResponse{Sort: v.Sort, ID: v.ID, Title: v.Title, Category: v.Category, Summary: v.Summary, Content: v.Content, Cover: v.Cover, Status: v.Status, Author: v.Author, CreatedAt: v.CreatedAt.Format(time.RFC3339), UpdatedAt: v.UpdatedAt.Format(time.RFC3339)}
}
func (s *Service) ListArticles(ctx context.Context, q systemmodel.PageQuery) (*systemmodel.PageResult, error) {
	page, size := normalizePage(q.Page, q.PageSize)
	b := s.client.Article.Query()
	category := strings.TrimSpace(q.Category)
	if category != "" {
		b.Where(article.CategoryEQ(category))
	}
	if q.Keyword != "" {
		b.Where(article.TitleContainsFold(q.Keyword))
	}
	if q.Status != nil {
		b.Where(article.StatusEQ(*q.Status))
	}
	total, err := b.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := b.Order(orderArticlesForCategory(category)).Offset((page - 1) * size).Limit(size).All(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]systemmodel.ArticleResponse, 0, len(rows))
	for _, row := range rows {
		list = append(list, articleResponse(row))
	}
	return &systemmodel.PageResult{List: list, Total: total, Page: page, PageSize: size}, nil
}
func (s *Service) GetArticle(ctx context.Context, id int) (*systemmodel.ArticleResponse, error) {
	v, err := s.client.Article.Get(ctx, id)
	if ent.IsNotFound(err) {
		return nil, apperror.New(apperror.CodeNotFound, "文章不存在")
	}
	if err != nil {
		return nil, err
	}
	r := articleResponse(v)
	return &r, nil
}

// Public reads always filter publication status, regardless of client input.
func (s *Service) ListPublishedArticles(ctx context.Context, q systemmodel.PageQuery) (*systemmodel.PageResult, error) {
	page, size := normalizePage(q.Page, q.PageSize)
	b := s.client.Article.Query().Where(article.StatusEQ(1))
	category := strings.TrimSpace(q.Category)
	if category != "" {
		b.Where(article.CategoryEQ(category))
	}
	if q.Keyword != "" {
		b.Where(article.TitleContainsFold(q.Keyword))
	}
	total, err := b.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	// Explicit selection avoids loading the HTML body for list requests.
	rows, err := b.Order(orderArticlesForCategory(category)).
		Offset((page-1)*size).Limit(size).
		Select(article.FieldSort, article.FieldID, article.FieldTitle, article.FieldCategory, article.FieldSummary, article.FieldCover,
			article.FieldStatus, article.FieldAuthor, article.FieldCreatedAt, article.FieldUpdatedAt).All(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]systemmodel.ArticleListItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, systemmodel.ArticleListItem{
			Sort: row.Sort, ID: row.ID, Title: row.Title, Category: row.Category, Summary: row.Summary, Cover: row.Cover,
			Status: row.Status, Author: row.Author,
			CreatedAt: row.CreatedAt.Format(time.RFC3339), UpdatedAt: row.UpdatedAt.Format(time.RFC3339),
		})
	}
	return &systemmodel.PageResult{List: list, Total: total, Page: page, PageSize: size}, nil
}

// Sort CDN news by priority and CDN help by publication time within each category.
func orderArticlesForCategory(category string) func(*sql.Selector) {
	switch category {
	case "cdn新闻":
		return func(s *sql.Selector) {
			s.OrderBy(sql.Desc(s.C(article.FieldSort)), sql.Desc(s.C(article.FieldID)))
		}
	case "cdn帮助":
		return func(s *sql.Selector) {
			s.OrderBy(sql.Desc(s.C(article.FieldCreatedAt)), sql.Desc(s.C(article.FieldID)))
		}
	default:
		return orderArticles
	}
}

func orderArticles(s *sql.Selector) {
	s.OrderBy(s.C(article.FieldCategory))
	s.OrderExpr(sql.ExprFunc(func(b *sql.Builder) {
		b.WriteString("CASE WHEN ").Ident(s.C(article.FieldCategory)).WriteString(" = ").Arg("cdn新闻")
		b.WriteString(" THEN -").Ident(s.C(article.FieldSort))
		b.WriteString(" WHEN ").Ident(s.C(article.FieldCategory)).WriteString(" = ").Arg("cdn帮助")
		b.WriteString(" THEN NULL ELSE ").Ident(s.C(article.FieldSort)).WriteString(" END ASC")
	}))
	s.OrderExpr(sql.ExprFunc(func(b *sql.Builder) {
		b.WriteString("CASE WHEN ").Ident(s.C(article.FieldCategory)).WriteString(" = ").Arg("cdn帮助")
		b.WriteString(" THEN ").Ident(s.C(article.FieldCreatedAt)).WriteString(" END DESC")
	}))
	s.OrderBy(sql.Desc(s.C(article.FieldID)))
}

func (s *Service) GetPublishedArticle(ctx context.Context, id int) (*systemmodel.ArticleResponse, error) {
	v, err := s.client.Article.Query().Where(article.IDEQ(id), article.StatusEQ(1)).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.New(apperror.CodeNotFound, "文章不存在")
	}
	if err != nil {
		return nil, err
	}
	r := articleResponse(v)
	return &r, nil
}
func (s *Service) CreateArticle(ctx context.Context, req systemmodel.CreateArticleRequest) error {
	if err := validateArticle(&req); err != nil {
		return err
	}
	if err := s.validateArticleSort(ctx, 0, req.Category, req.Sort); err != nil {
		return err
	}
	err := s.client.Article.Create().SetTitle(req.Title).SetCategory(req.Category).SetSummary(req.Summary).SetContent(req.Content).SetCover(req.Cover).SetStatus(req.Status).SetAuthor(req.Author).SetNillableSort(req.Sort).Exec(ctx)
	return articleMutationError(err)
}
func (s *Service) UpdateArticle(ctx context.Context, id int, req systemmodel.UpdateArticleRequest) error {
	if err := validateArticle(&req); err != nil {
		return err
	}
	if err := s.validateArticleSort(ctx, id, req.Category, req.Sort); err != nil {
		return err
	}
	err := s.client.Article.UpdateOneID(id).SetTitle(req.Title).SetCategory(req.Category).SetSummary(req.Summary).SetContent(req.Content).SetCover(req.Cover).SetStatus(req.Status).SetAuthor(req.Author).SetNillableSort(req.Sort).Exec(ctx)
	return articleMutationError(err)
}
func (s *Service) DeleteArticle(ctx context.Context, id int) error {
	return articleMutationError(s.client.Article.DeleteOneID(id).Exec(ctx))
}

func (s *Service) validateArticleSort(ctx context.Context, id int, category string, value *int) error {
	if value == nil {
		if id == 0 {
			return apperror.New(apperror.CodeBadRequest, "请输入排序号")
		}
		current, err := s.client.Article.Get(ctx, id)
		if err != nil {
			return articleMutationError(err)
		}
		value = &current.Sort
	}
	if *value < 0 || *value > 2147483647 {
		return apperror.New(apperror.CodeBadRequest, "排序号必须为 0 至 2147483647 之间的整数")
	}
	query := s.client.Article.Query().Where(article.CategoryEQ(category), article.SortEQ(*value))
	if id != 0 {
		query.Where(article.IDNEQ(id))
	}
	exists, err := query.Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return apperror.New(apperror.CodeDuplicateResource, "排序号已被其他文章使用，请填写不同的排序号")
	}
	return nil
}

func articleMutationError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "article_sort_unique" {
		return apperror.New(apperror.CodeDuplicateResource, "排序号已被其他文章使用，请填写不同的排序号")
	}
	if ent.IsNotFound(err) {
		return apperror.New(apperror.CodeNotFound, "文章不存在")
	}
	return err
}

func validateArticle(req *systemmodel.CreateArticleRequest) error {
	req.Title = strings.TrimSpace(req.Title)
	req.Category = strings.TrimSpace(req.Category)
	req.Content = strings.TrimSpace(req.Content)
	if req.Title == "" || !hasArticleContent(req.Content) {
		return apperror.New(apperror.CodeBadRequest, "文章标题和正文不能为空")
	}
	return nil
}

// wangEditor represents an empty document as <p><br></p>.
func hasArticleContent(content string) bool {
	tokenizer := html.NewTokenizer(strings.NewReader(content))
	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			return false
		case html.TextToken:
			if strings.TrimSpace(string(tokenizer.Text())) != "" {
				return true
			}
		case html.StartTagToken, html.SelfClosingTagToken:
			name, _ := tokenizer.TagName()
			switch string(name) {
			case "img", "video", "audio", "iframe":
				return true
			}
		}
	}
}
