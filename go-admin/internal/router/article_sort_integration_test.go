package router

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"base/ent/article"
	"base/internal/config"
	"base/internal/db"
	systemmodel "base/internal/model/system"
	"base/internal/pkg/apperror"
	systemservice "base/internal/service/system"
)

func TestArticleSortUniquenessIntegration(t *testing.T) {
	if os.Getenv("GO_ADMIN_INTEGRATION") != "1" {
		t.Skip("set GO_ADMIN_INTEGRATION=1")
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
	if err := db.Migrate(ctx, database); err != nil {
		t.Fatal(err)
	}
	category := fmt.Sprintf("sort-test-%d", time.Now().UnixNano())
	otherCategory := category + "-other"
	t.Cleanup(func() {
		_, err := database.Client.Article.Delete().Where(article.CategoryIn(category, otherCategory)).Exec(ctx)
		if err != nil {
			t.Error(err)
		}
	})
	service := systemservice.New(database.Client)
	save := func(category, title string, sort int) error {
		return service.CreateArticle(ctx, systemmodel.CreateArticleRequest{Category: category, Title: title, Content: "<p>test</p>", Status: 1, Sort: &sort})
	}
	if err := save(category, "first", 7); err != nil {
		t.Fatal(err)
	}
	var appErr *apperror.Error
	if err := save(category, "duplicate", 7); !errors.As(err, &appErr) || appErr.Code != apperror.CodeDuplicateResource {
		t.Fatalf("expected duplicate error: %v", err)
	}
	if err := save(otherCategory, "same number in another category", 7); err != nil {
		t.Fatal(err)
	}
	first, err := database.Client.Article.Query().Where(article.CategoryEQ(category)).Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	number := 23
	payload := systemmodel.UpdateArticleRequest{Category: category, Title: "first", Content: "<p>updated</p>", Status: 1, Sort: &number}
	if err := service.UpdateArticle(ctx, first.ID, payload); err != nil {
		t.Fatal(err)
	}
	stored, err := database.Client.Article.Get(ctx, first.ID)
	if err != nil || stored.Sort != 23 {
		t.Fatalf("exact number not saved: %+v %v", stored, err)
	}
	if err := service.UpdateArticle(ctx, first.ID, payload); err != nil {
		t.Fatalf("own number rejected: %v", err)
	}
	if err := save(category, "second", 5); err != nil {
		t.Fatal(err)
	}
	number = 5
	if err := service.UpdateArticle(ctx, first.ID, payload); !errors.As(err, &appErr) || appErr.Code != apperror.CodeDuplicateResource {
		t.Fatalf("duplicate update accepted: %v", err)
	}
	stored, _ = database.Client.Article.Get(ctx, first.ID)
	if stored.Sort != 23 {
		t.Fatal("rejected update changed stored number")
	}
	number = 7
	payload.Category = otherCategory
	if err := service.UpdateArticle(ctx, first.ID, payload); !errors.As(err, &appErr) || appErr.Code != apperror.CodeDuplicateResource {
		t.Fatalf("category change bypassed uniqueness: %v", err)
	}
	for _, published := range []bool{false, true} {
		query := systemmodel.PageQuery{Category: category, Page: 1, PageSize: 1}
		if published {
			result, err := service.ListPublishedArticles(ctx, query)
			if err != nil || result.List.([]systemmodel.ArticleListItem)[0].Sort != 5 {
				t.Fatalf("public sort: %+v %v", result, err)
			}
		} else {
			result, err := service.ListArticles(ctx, query)
			if err != nil || result.List.([]systemmodel.ArticleResponse)[0].Sort != 5 {
				t.Fatalf("admin sort: %+v %v", result, err)
			}
		}
	}
	// Bypass the service pre-check to verify database protection under concurrent writes.
	start := make(chan struct{})
	results := make(chan error, 2)
	var workers sync.WaitGroup
	for i := 0; i < 2; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			results <- database.Client.Article.Create().SetCategory(category).SetTitle("race").SetContent("<p>test</p>").SetSort(99).Exec(ctx)
		}()
	}
	close(start)
	workers.Wait()
	close(results)
	successes := 0
	for err := range results {
		if err == nil {
			successes++
		}
	}
	if successes != 1 {
		t.Fatalf("concurrent writes: %d successes, want 1", successes)
	}
}
