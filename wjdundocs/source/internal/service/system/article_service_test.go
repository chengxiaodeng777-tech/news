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
	want := `SELECT * FROM "article" ORDER BY "article"."category", CASE WHEN "article"."category" = $1 THEN -"article"."sort" ELSE "article"."sort" END ASC LIMIT 1 OFFSET 1`
	if query != want || !reflect.DeepEqual(args, []any{"cdn新闻"}) {
		t.Fatalf("unexpected article ordering: %s %v", query, args)
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
