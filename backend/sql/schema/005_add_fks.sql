-- +goose Up

ALTER TABLE "recurring_themes"
ADD FOREIGN KEY("user_id") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "episode_references"
ADD FOREIGN KEY("recurring_theme_id") REFERENCES "recurring_themes"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "bingo_card_cells"
ADD FOREIGN KEY("recurring_theme_id") REFERENCES "recurring_themes"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "bingo_card_cells"
ADD FOREIGN KEY("user_id") REFERENCES "users"("id")
ON UPDATE NO ACTION ON DELETE CASCADE;

-- +goose Down

