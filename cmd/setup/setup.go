package main

import (
	"os"

	"github.com/bitstorm-tech/cockaigne/internal/account"
	"github.com/bitstorm-tech/cockaigne/internal/deal"
	"github.com/bitstorm-tech/cockaigne/internal/like"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/gofiber/fiber/v2/log"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const createLikeCountsViewQuery = `
create or replace view
  cockaigne.like_counts as
select
  deal_id,
  count(deal_id) as like_count
from
  likes
group by
  deal_id
order by
  like_count desc;
`

const createActiveDealsViewQuery = `
create or replace view
  cockaigne.active_deals as
select
  d.id,
  d.dealer_id,
  d.title,
  d.description,
  d.category_id,
  d.duration_in_hours,
  d.start,
  d.start::time as start_time,
  a.username,
  a.location,
  coalesce(c.like_count, 0) as likes
from
  deals d
  join accounts a on d.dealer_id = a.id
  left join like_counts c on c.deal_id = d.id
where
  d.is_template = false
  and now() between d."start" and d."start"  + (d.duration_in_hours || ' hours')::interval
order by
  start_time;
`

func main() {
	createSchemaAndInstallPostGis()
	persistence.ConnectToDb()
	persistence.DB.Exec("drop view if exists cockaigne.like_counts cascade;")
	persistence.DB.Exec("drop view if exists cockaigne.active_deals cascade;")
	err := persistence.DB.AutoMigrate(&account.Account{}, &deal.Deal{}, &deal.Category{}, &like.Like{}, &account.FavoriteCategory{})
	if err != nil {
		panic("can't migrate database: " + err.Error())
	}
	persistence.DB.Updates(&categories)
	persistence.DB.Exec(createLikeCountsViewQuery)
	persistence.DB.Exec(createActiveDealsViewQuery)
}

func createSchemaAndInstallPostGis() {
	db, err := gorm.Open(postgres.Open(persistence.ConnectionString))
	if err != nil {
		panic("can't create cockaigen schema: " + err.Error())
	}

	user := os.Getenv("PGUSER")
	log.Info("Create postgis extension ...")
	db.Exec("create extension if not exists postgis;")
	log.Info("Create schema cockaigne ...")
	db.Exec("create schema if not exists cockaigne authorization " + user + ";")

	_db, _ := db.DB()
	_db.Close()
}

var categories = []deal.Category{
	{Name: "Elektronik & Technik"},
	{Name: "Unterhaltung & Gaming"},
	{Name: "Lebensmittel & Haushalt"},
	{Name: "Fashion, Schmuck & Lifestyle"},
	{Name: "Beauty, Wellness & Gesundheit"},
	{Name: "Family & Kids"},
	{Name: "Home & Living"},
	{Name: "Baumarkt & Garten"},
	{Name: "Auto, Fahhrad & Motorrad"},
	{Name: "Gastronomie, Bars & Cafes"},
	{Name: "Kultur & Freizeit"},
	{Name: "Sport & Outdoor"},
	{Name: "Reisen, Hotels & Übernachtungen"},
	{Name: "Dienstleistungen & Finanzen"},
	{Name: "Floristik"},
	{Name: "Sonstiges"},
}
