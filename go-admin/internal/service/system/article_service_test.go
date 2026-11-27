package system

import (
	"reflect"
	"testing"

	"base/ent/article"
	systemmodel "base/internal/model/system"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
)

func TestOrderArticles(t *testing.T) {
	selector := sql.Dialect(dialect.Postgres).Select("*").From(sql.Table(article.Table))
	orderArticles(selector)
	selector.Limit(1).Offset(1)
	query, args := selector.Query()
	want := `SELECT * FROM "article" ORDER BY "article"."category", CASE WHEN "article"."category" = $1 THEN -"article"."sort" WHEN "article"."category" = $2 THEN NULL ELSE "article"."sort" END ASC, CASE WHEN "article"."category" = $3 THEN "article"."created_at" END DESC, "article"."id" DESC LIMIT 1 OFFSET 1`
	if query != want || !reflect.DeepEqual(args, []any{"cdn新闻", "cdn帮助", "cdn帮助"}) {
		t.Fatalf("unexpected article ordering: %s %v", query, args)
	}
}

func TestOrderArticlesForCategory(t *testing.T) {
	for _, tc := range []struct {
		category string
		want     string
	}{
		{"cdn新闻", `SELECT * FROM "article" ORDER BY "article"."sort" DESC, "article"."id" DESC LIMIT 1 OFFSET 1`},
		{"cdn帮助", `SELECT * FROM "article" ORDER BY "article"."created_at" DESC, "article"."id" DESC LIMIT 1 OFFSET 1`},
	} {
		t.Run(tc.category, func(t *testing.T) {
			selector := sql.Dialect(dialect.Postgres).Select("*").From(sql.Table(article.Table))
			orderArticlesForCategory(tc.category)(selector)
			selector.Limit(1).Offset(1)
			query, args := selector.Query()
			if query != tc.want || len(args) != 0 {
				t.Fatalf("unexpected article ordering: %s %v", query, args)
			}
		})
	}
}

func TestValidateArticleContent(t *testing.T) {
	for _, content := range []string{"", "  ", "<p><br></p>", "<p>&nbsp; </p>", "<p><strong></strong></p>"} {
		req := systemmodel.CreateArticleRequest{Title: "文章", Content: content}
		if err := validateArticle(&req); err == nil {
			t.Errorf("accepted empty document %q", content)
		}
	}
	for _, content := range []string{"<p>正文</p>", `<p><img src="https://example.com/image.png"></p>`, "plain text"} {
		req := systemmodel.CreateArticleRequest{Title: "  文章  ", Content: content}
		if err := validateArticle(&req); err != nil {
			t.Errorf("rejected document %q: %v", content, err)
		}
		if req.Title != "文章" {
			t.Errorf("title was not trimmed: %q", req.Title)
		}
	}
	req := systemmodel.CreateArticleRequest{Title: " \t ", Content: "<p>正文</p>"}
	if err := validateArticle(&req); err == nil {
		t.Error("accepted a blank title")
	}
}
