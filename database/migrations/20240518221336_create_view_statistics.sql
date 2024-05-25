-- +goose Up
-- +goose StatementBegin
create or replace view
    statistics_view as
select d.id                         as deal_id,
       d.dealer_id,
       d.title,
       d.start,
       d.duration_in_hours,
       coalesce(c.clickcount, 0)    as clickcount,
       coalesce(f.favoritecount, 0) as favoritecount,
       coalesce(l.likecount, 0)     as likecount
from deals d
         left join click_counts_view c on d.id = c.deal_id
         left join favorite_counts_view f on d.id = f.deal_id
         left join like_counts_view l on d.id = l.deal_id;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
drop view statistics_view;
-- +goose StatementEnd
