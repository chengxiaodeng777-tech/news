package db

import (
	"context"
	"database/sql"
)

// A trigger enforces uniqueness without silently renumbering legacy duplicates.
// Serialize writers before checking so concurrent saves cannot claim the same number.
func migrateArticleSort(ctx context.Context, database *sql.DB) error {
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `
CREATE OR REPLACE FUNCTION enforce_article_sort_unique() RETURNS trigger
LANGUAGE plpgsql VOLATILE AS $$
BEGIN
 PERFORM pg_advisory_xact_lock(812904103);
 IF EXISTS (SELECT 1 FROM article WHERE category = NEW.category AND sort = NEW.sort AND id IS DISTINCT FROM NEW.id) THEN
  RAISE EXCEPTION 'article sort already exists'
   USING ERRCODE = '23505', CONSTRAINT = 'article_sort_unique';
 END IF;
 RETURN NEW;
END;
$$;
DROP TRIGGER IF EXISTS article_sort_unique ON article;
CREATE TRIGGER article_sort_unique BEFORE INSERT OR UPDATE OF sort, category ON article
FOR EACH ROW EXECUTE FUNCTION enforce_article_sort_unique();
`)
	if err != nil {
		return err
	}
	return tx.Commit()
}
