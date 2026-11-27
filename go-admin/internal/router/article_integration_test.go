package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"base/ent"
	"base/internal/config"
	"base/internal/db"
	systemmodel "base/internal/model/system"
	systemservice "base/internal/service/system"

	"go.uber.org/zap"
)

// Opt in with GO_ADMIN_INTEGRATION=1; uses the configured database and removes its fixtures.
func TestPublicArticlesIntegration(t *testing.T) {
	if os.Getenv("GO_ADMIN_INTEGRATION") != "1" {
		t.Skip("set GO_ADMIN_INTEGRATION=1 to test against the configured database")
	}
	t.Chdir("../..")
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	database, err := db.Open(ctx, cfg.Database)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close(database) })
	keyword := fmt.Sprintf("public-article-test-%d", time.Now().UnixNano())
	sortBase := int(time.Now().UnixNano()%1000000000) + 100000000
	var rows []*ent.Article
	for i, status := range []int8{0, 1, 1} {
		row, err := database.Client.Article.Create().SetTitle(keyword).SetCategory("技术分享").SetContent("<p>正文</p>").SetStatus(status).SetSort(sortBase + i + 1).Save(ctx)
		if err != nil {
			t.Fatal(err)
		}
		rows = append(rows, row)
		t.Cleanup(func() {
			if err := database.Client.Article.DeleteOneID(row.ID).Exec(ctx); err != nil {
				t.Error(err)
			}
		})
	}
	r := New(ctx, cfg, zap.NewNop(), Dependencies{Database: database})
	// Sorting must be applied before pagination and remain stable after edits.
	service := systemservice.New(database.Client)
	sortValue := sortBase + 10
	if err := service.UpdateArticle(ctx, rows[2].ID, systemmodel.UpdateArticleRequest{
		Title: keyword, Category: "技术分享", Content: "<p>正文</p>", Status: 1, Sort: &sortValue,
	}); err != nil {
		t.Fatal(err)
	}
	assertOrder := func(wantFirst, wantSecond int) {
		t.Helper()
		for _, published := range []bool{false, true} {
			q := systemmodel.PageQuery{Keyword: keyword, Page: 1, PageSize: 1}
			status := int8(1)
			q.Status = &status
			for page, want := range []int{wantFirst, wantSecond} {
				q.Page = page + 1
				var result *systemmodel.PageResult
				var err error
				var got int
				if published {
					result, err = service.ListPublishedArticles(ctx, q)
					if err == nil && len(result.List.([]systemmodel.ArticleListItem)) == 1 {
						got = result.List.([]systemmodel.ArticleListItem)[0].ID
					}
				} else {
					result, err = service.ListArticles(ctx, q)
					if err == nil && len(result.List.([]systemmodel.ArticleResponse)) == 1 {
						got = result.List.([]systemmodel.ArticleResponse)[0].ID
					}
				}
				if err != nil || got != want {
					t.Fatalf("published=%v page=%d: id=%d want=%d err=%v", published, q.Page, got, want, err)
				}
			}
		}
	}
	assertOrder(rows[1].ID, rows[2].ID)
	if err := service.UpdateArticle(ctx, rows[2].ID, systemmodel.UpdateArticleRequest{
		Title: keyword, Category: "技术分享", Content: "<p>修改正文</p>", Status: 1,
	}); err != nil {
		t.Fatal(err)
	}
	updated, err := service.GetArticle(ctx, rows[2].ID)
	if err != nil || updated.Sort != sortValue {
		t.Fatalf("sort was not preserved: article=%+v err=%v", updated, err)
	}
	assertOrder(rows[1].ID, rows[2].ID)
	// An explicit unused smaller number moves the article ahead.
	sortValue = sortBase
	if err := service.UpdateArticle(ctx, rows[2].ID, systemmodel.UpdateArticleRequest{
		Title: keyword, Category: "技术分享", Content: rows[1].Content, Status: 1, Sort: &sortValue,
	}); err != nil {
		t.Fatal(err)
	}
	assertOrder(rows[2].ID, rows[1].ID)
	request := func(path string, want int) *httptest.ResponseRecorder {
		t.Helper()
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		if w.Code != want {
			t.Fatalf("%s: status %d, want %d: %s", path, w.Code, want, w.Body.String())
		}
		return w
	}
	for _, status := range []string{"", "&status=0", "&status=1", "&status=invalid"} {
		w := request("/api/articles?page=1&pageSize=1&keyword="+url.QueryEscape(keyword)+status, http.StatusOK)
		var result struct {
			Code int `json:"code"`
			Data struct {
				List     []systemmodel.ArticleResponse `json:"list"`
				Total    int                           `json:"total"`
				Page     int                           `json:"page"`
				PageSize int                           `json:"pageSize"`
			} `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		if result.Code != 0 || result.Data.Total != 2 || len(result.Data.List) != 1 || result.Data.Page != 1 || result.Data.PageSize != 1 {
			t.Fatalf("unexpected pagination: %s", w.Body.String())
		}
		if result.Data.List[0].Status != 1 {
			t.Fatal("draft leaked in public list")
		}
		if result.Data.List[0].Category != "技术分享" {
			t.Fatal("category missing from public list")
		}
		var raw struct {
			Data struct {
				List []map[string]json.RawMessage `json:"list"`
			} `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
			t.Fatal(err)
		}
		if _, exists := raw.Data.List[0]["content"]; exists {
			t.Fatal("public list must not include content, even as an empty field")
		}
	}
	for _, row := range rows {
		want := http.StatusOK
		if row.Status == 0 {
			want = http.StatusNotFound
		}
		w := request(fmt.Sprintf("/api/articles/%d", row.ID), want)
		if want == http.StatusOK {
			var result struct {
				Data systemmodel.ArticleResponse `json:"data"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
				t.Fatal(err)
			}
			if result.Data.ID != row.ID || result.Data.Content != row.Content || result.Data.Category != "技术分享" {
				t.Fatal("incorrect public detail")
			}
		}
	}
	request("/api/articles/2147483647", http.StatusNotFound)

	categoryKeyword := fmt.Sprintf("category-test-%d", time.Now().UnixNano())
	for i, fixture := range []struct {
		category string
		status   int8
	}{
		{"cdn帮助", 1}, {"cdn新闻", 1}, {"cdn帮助", 0}, {"技术分享", 1},
	} {
		row, err := database.Client.Article.Create().SetTitle(categoryKeyword).SetCategory(fixture.category).SetContent("<p>分类正文</p>").SetStatus(fixture.status).SetSort(sortBase + 20 + i).Save(ctx)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if err := database.Client.Article.DeleteOneID(row.ID).Exec(ctx); err != nil {
				t.Error(err)
			}
		})
	}
	for _, category := range []string{"cdn帮助", "cdn新闻", " cdn帮助 ", "不存在分类"} {
		w := request("/api/articles?page=1&pageSize=1&keyword="+url.QueryEscape(categoryKeyword)+"&category="+url.QueryEscape(category), http.StatusOK)
		var result struct {
			Data struct {
				List  []systemmodel.ArticleListItem `json:"list"`
				Total int                           `json:"total"`
			} `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		if category == "不存在分类" {
			if result.Data.Total != 0 || len(result.Data.List) != 0 {
				t.Fatal("unknown category must have an empty list")
			}
			continue
		}
		if result.Data.Total != 1 || len(result.Data.List) != 1 || result.Data.List[0].Category != strings.TrimSpace(category) {
			t.Fatalf("incorrect category filtering or total: %s", w.Body.String())
		}
	}
}
