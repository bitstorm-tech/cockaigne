-- +goose Up
-- +goose StatementBegin
create or replace view click_counts_view as
select deal_id,
       count(deal_id) as clickcount
from deal_clicks
group by deal_id
order by clickcount desc;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
drop view click_counts_view;
-- +goose StatementEnd
