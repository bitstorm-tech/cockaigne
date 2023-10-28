set schema 'cockaigne';

begin;

create table categories
(
    id        integer primary key,
    name      text    not null,
    is_active boolean not null default true
);



insert into categories (id, name)
values (1, 'Elektronik & Technik'),
       (2, 'Unterhaltung & Gaming'),
       (3, 'Lebensmittel & Haushalt'),
       (4, 'Fashion, Schmuck & Lifestyle'),
       (5, 'Beauty, Wellness & Gesundheit'),
       (6, 'Family & Kids'),
       (7, 'Home & Living'),
       (8, 'Baumarkt & Garten'),
       (9, 'Auto, Fahhrad & Motorrad'),
       (10, 'Gastronomie, Bars & Cafes'),
       (11, 'Kultur & Freizeit'),
       (12, 'Sport & Outdoor'),
       (13, 'Reisen, Hotels & Übernachtungen'),
       (14, 'Dienstleistungen & Finanzen'),
       (15, 'Floristik'),
       (16, 'Sonstiges');

create table accounts
(
    id                      uuid                         not null default public.uuid_generate_v4() primary key,
    username                text                         not null,
    email                   text                         not null,
    password                text                         not null,
    is_dealer               boolean                      not null default false,
    street                  text                         null,
    age                     integer                      null,
    gender                  text                         null,
    default_category        integer                      null references categories (id) on delete restrict on update cascade,
    house_number            text                         null,
    city                    text                         null,
    zip                     integer                      null,
    phone                   text                         null,
    tax_id                  text                         null,
    use_current_location    boolean                      not null default false,
    search_radius_in_meters integer                      not null default 500,
    "location"              public.geometry(point, 4326) null,
    created                 timestamptz                  not null default now()
);

create unique index accounts_username_idx on accounts (lower(username));
create index accounts_location_idx on accounts using GIST (location);



create table deals
(
    id                uuid        not null default public.uuid_generate_v4() primary key,
    dealer_id         uuid        not null references accounts (id) on delete restrict on update cascade,
    title             text        not null,
    description       text        not null,
    category_id       integer     not null references categories (id) on delete restrict on update cascade,
    duration_in_hours integer     not null,
    "start"           timestamptz not null,
    "template"        boolean     not null default false,
    created           timestamptz not null default now()
);



create table dealer_ratings
(
    user_id     uuid        not null references accounts (id) on delete restrict on update cascade,
    dealer_id   uuid        not null references accounts (id) on delete restrict on update cascade,
    stars       integer     not null,
    rating_text text        null,
    created     timestamptz not null default now(),
    constraint "dealer_ratings_pk" unique (user_id, dealer_id)
);



create table hot_deals
(
    user_id uuid        not null references accounts (id) on delete restrict on update cascade,
    deal_id uuid        not null references deals (id) on delete restrict on update cascade,
    created timestamptz not null default now(),
    constraint "hot_deals_pk" unique (user_id, deal_id)
);



create table reported_deals
(
    reporter_id uuid        not null references accounts (id) on delete restrict on update cascade,
    deal_id     uuid        not null references deals (id) on delete restrict on update cascade,
    reason      text        not null,
    created     timestamptz not null default now(),
    constraint "reported_deals_pk" unique (reporter_id, deal_id)
);



create table favorite_dealers
(
    user_id   uuid        not null references accounts (id) on delete restrict on update cascade,
    dealer_id uuid        not null references accounts (id) on delete restrict on update cascade,
    created   timestamptz not null default now(),
    constraint "favorite_dealer_pk" unique (user_id, dealer_id)
);



create table likes
(
    user_id uuid        not null references accounts (id) on delete restrict on update cascade,
    deal_id uuid        not null references deals (id) on delete restrict on update cascade,
    created timestamptz not null default now(),
    constraint "likes_pk" unique (user_id, deal_id)
);



create table selected_categories
(
    user_id     uuid        not null references accounts (id) on delete restrict on update cascade,
    category_id integer     not null references categories (id) on delete restrict on update cascade,
    created     timestamptz not null default now(),
    constraint "selected_categories_pk" unique (user_id, category_id)
);



create table vouchers
(
    id               integer primary key,
    code             text        not null unique,
    start            date        not null,
    "end"            date        null,
    duration_in_days integer     null,
    is_active        bool        not null,
    multi_use        bool        not null,
    comment          text        not null,
    created          timestamptz not null default now()
);



create table activated_vouchers
(
    user_id      uuid        not null references accounts (id) on delete restrict on update cascade,
    voucher_code text        not null references vouchers (code) on delete restrict on update cascade,
    activated    timestamptz not null default now(),
    constraint "activated_vouchers_pk" unique (user_id, voucher_code)
);



create table plans
(
    id                  serial primary key,
    stripe_product_id   text      not null,
    free_days_per_month integer   not null,
    active              boolean   not null default false,
    created             timestamp not null default now()
);



create table subscriptions
(
    user_id                uuid        not null references accounts (id) on delete restrict on update cascade,
    plan_id                integer     not null references plans (id) on delete restrict on update cascade,
    stripe_subscription_id text        not null,
    active                 boolean     not null default false,
    created                timestamptz not null default now(),
    constraint "subscriptions_pk" unique (user_id, stripe_subscription_id)
);



create or replace view like_counts_view as
select deal_id,
       count(deal_id) as likecount
from likes
group by deal_id
order by likecount desc;



create or replace view dealer_view as
select a.id,
       a.username,
       a.street,
       a.house_number,
       a.phone,
       a.zip,
       a.city,
       (select c.name
        from categories c
        where c.id = a.default_category) as category
from accounts a
where is_dealer is true;



create or replace view active_deals_view as
select d.id,
       d.dealer_id,
       d.title,
       d.description,
       d.category_id,
       d.duration_in_hours,
       d.start,
       d.start::time            as start_time,
       a.username,
       a.location,
       coalesce(c.likecount, 0) as likes
from deals d
         join accounts a on d.dealer_id = a.id
         left join like_counts_view c on c.deal_id = d.id
where d.template = false
  and now() between d."start" and d."start" + (d.duration_in_hours || ' hours')::interval
order by start_time;



create or replace view future_deals_view as
select d.id,
       d.dealer_id,
       d.title,
       d.description,
       d.category_id,
       d.duration_in_hours,
       d.start,
       d.start::time            as start_time,
       a.username,
       a.location,
       coalesce(c.likecount, 0) as likes
from deals d
         join accounts a on d.dealer_id = a.id
         left join like_counts_view c on c.deal_id = d.id
where d.template = false
  and d."start" > now()
order by start_time;



create or replace view past_deals_view as
select d.id,
       d.dealer_id,
       d.title,
       d.description,
       d.category_id,
       d.duration_in_hours,
       d.start,
       d.start::time            as start_time,
       a.username,
       a.location,
       coalesce(c.likecount, 0) as likes
from deals d
         join accounts a on d.dealer_id = a.id
         left join like_counts_view c on c.deal_id = d.id
where d.template = false
  and (d."start" + (d.duration_in_hours || ' hours')::interval) < now()
order by start_time;



create or replace view dealer_ratings_view as
select r.user_id,
       r.dealer_id,
       r.stars,
       r.rating_text,
       r.created,
       a.username
from dealer_ratings r
         join accounts a on r.user_id = a.id;



create or replace view favorite_dealers_view as
select f.user_id,
       f.dealer_id,
       a.username
from favorite_dealers f
         join accounts a on f.dealer_id = a.id;



create or replace view invoice_metadata_view as
select d.dealer_id,
       extract(
               year
               from
               d.start
       )                        as year,
       extract(
               month
               from
               d.start
       )                        as month,
       count(d.id)              as deals,
       sum(d.duration_in_hours) as total_duration_in_min
from accounts a
         join deals d on d.dealer_id = a.id
where d.template = false
group by d.dealer_id,
         year,
         month;



create or replace view active_vouchers_view as
select av.user_id,
       av.activated,
       v.code,
       v.start,
       v."end",
       v.duration_in_days
from vouchers v
         join activated_vouchers av on v.code = av.voucher_code
where v.is_active
    and (now() between v."start" and v."end")
   or (now() between v."start" and v."start" + (v.duration_in_days || ' days')::interval);



commit;