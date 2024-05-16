set schema 'cockaigne';



begin;



-------------------------------------------------------------------------------
-------------------------------------------------------------------------------
-- Table definitions
-------------------------------------------------------------------------------
-------------------------------------------------------------------------------
create table categories
(
    id        integer primary key,
    name      text    not null,
    is_active boolean not null default true
);



create table accounts
(
    id                      uuid                         not null default public.uuid_generate_v4() primary key,
    username                text                         not null,
    email                   text                         not null,
    password                text                         not null,
    active                  boolean                      not null default false,
    activation_code         integer                      null,
    change_password_code    text                         null,
    change_email_code       text                         null,
    new_email               text                         null,
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
    use_location_service    boolean                      not null default false,
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
    payment_state     text,
    "start"           timestamptz not null,
    "template"        boolean     not null default false,
    created           timestamptz not null default now()
);



create table dealer_ratings
(
    user_id   uuid        not null references accounts (id) on delete restrict on update cascade,
    dealer_id uuid        not null references accounts (id) on delete restrict on update cascade,
    stars     integer     not null,
    "text"    text        null,
    created   timestamptz not null default now(),
    constraint "dealer_ratings_pk" unique (user_id, dealer_id)
);



create table favorite_deals
(
    user_id uuid        not null references accounts (id) on delete restrict on update cascade,
    deal_id uuid        not null references deals (id) on delete restrict on update cascade,
    created timestamptz not null default now(),
    constraint "favorite_deals_pk" unique (user_id, deal_id)
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
    code                text        not null primary key,
    discount_in_percent int         not null,
    start               date        not null,
    "end"               date        not null,
    is_active           bool        not null,
    multi_use           bool        not null,
    comment             text        not null,
    created             timestamptz not null default now()
);



create table redeemed_vouchers
(
    account_id  uuid        not null references accounts (id) on delete restrict on update cascade,
    code        text        not null references vouchers (code) on delete restrict on update cascade,
    redeemed_at timestamptz not null default now(),
    constraint "redeemed_vouchers_pk" unique (account_id, code)
);



create table plans
(
    id                  serial primary key,
    name                text      not null,
    stripe_product_id   text      not null unique,
    stripe_price_id     text      not null unique,
    free_days_per_month integer   not null,
    active              boolean   not null default false,
    created             timestamp not null default now()
);



create table subscriptions
(
    id                     serial primary key,
    account_id             uuid        not null references accounts (id) on delete restrict on update cascade,
    plan_id                integer     not null references plans (id) on delete restrict on update cascade,
    stripe_subscription_id text        null,
    stripe_tracking_id     uuid        null,
    state                  text        not null,
    created                timestamptz not null default now(),
    activated              timestamptz null,
    canceled               timestamptz null,
    constraint "account_subscription_key" unique (account_id, stripe_subscription_id)
);



create table contact_messages
(
    account_id uuid        not null references accounts (id) on delete restrict on update cascade,
    message    text        not null,
    created    timestamptz not null default now()
);



create table contact_messages
(
    account_id uuid        not null references accounts (id) on delete restrict on update cascade,
    message    text        not null,
    created    timestamptz not null default now()
);



-------------------------------------------------------------------------------
-------------------------------------------------------------------------------
-- View definitions
-------------------------------------------------------------------------------
-------------------------------------------------------------------------------
create or replace view like_counts_view as
select deal_id,
       count(deal_id) as likecount
from likes
group by deal_id
order by likecount desc;



create or replace view favorite_counts_view as
select deal_id,
       count(deal_id) as favoritecount
from favorite_deals
group by deal_id
order by favoritecount desc;



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
  and d.payment_state = 'PAYED'
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
  and d.payment_state = 'PAYED'
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
  and d.payment_state = 'PAYED'
order by start_time;



create or replace view top_deals_view as
select id,
       title,
       username,
       dealer_id,
       category_id,
       location,
       coalesce(likecount, 0)     as likes,
       coalesce(favoritecount, 0) as favorites
from active_deals_view a
         left join like_counts_view l on l.deal_id = a.id
         left join favorite_counts_view f on f.deal_id = a.id
order by likes desc,
         favorites desc;


create or replace view dealer_ratings_view as
select r.user_id,
       r.dealer_id,
       r.stars,
       r.text,
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
select rv.account_id,
       v.code,
       v.start,
       v."end",
       v.discount_in_percent
from vouchers v
         join redeemed_vouchers rv on v.code = rv.code
where v.is_active
  and now() between v."start" and v."end";



-------------------------------------------------------------------------------
-------------------------------------------------------------------------------
-- Data
-------------------------------------------------------------------------------
-------------------------------------------------------------------------------
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



commit;
