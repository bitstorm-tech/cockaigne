package main

import (
	"github.com/bitstorm-tech/cockaigne/internal/deal"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	persistence.ConnectToDb()
	persistence.DB.Create(&categories)
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
