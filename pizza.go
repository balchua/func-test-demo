package main

import fdb "github.com/balchua/func-test-demo/pkg/datastore"

type datastoreConfig struct {
	host     string
	port     int
	user     string
	password string
	dbName   string
}

type pizza struct {
	crustSize  float32
	tomatoes   float32
	pineapples float32
	hams       float32
	ds         datastoreConfig
}

func NewPizza(ds datastoreConfig, crustSize, pineapples, hams, tomatoes float32) pizza {
	return pizza{
		crustSize:  crustSize,
		tomatoes:   tomatoes,
		pineapples: pineapples,
		hams:       hams,
		ds:         ds,
	}
}

func ingredientWithin[T int32 | float32](min, max, count T) bool {
	if count >= min && count <= max {
		return true
	}
	return false
}

func (p *pizza) CookPizza() string {
	ds := fdb.NewIngredientsStore(p.ds.host, p.ds.port, p.ds.user, p.ds.password, p.ds.dbName)
	items := ds.GetThresholdsByCrustSize(p.crustSize)
	for _, item := range items {
		if item.IngredientType == "H" && !ingredientWithin(item.MinValue, item.MaxValue, p.hams) {
			return "not perfect"
		}
		if item.IngredientType == "P" && !ingredientWithin(item.MinValue, item.MaxValue, p.pineapples) {
			return "not perfect"
		}
		if item.IngredientType == "T" && !ingredientWithin(item.MinValue, item.MaxValue, p.tomatoes) {
			return "not perfect"
		}
	}
	return "perfect"
}
