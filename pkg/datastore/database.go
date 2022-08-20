package datastore

import (
	"database/sql"
	"fmt"

	"log"

	_ "github.com/lib/pq"
)

type IngredientsThresholds struct {
	MinValue       float32
	MaxValue       float32
	IngredientType string
	CrustSize      float32
}

type ingredientsStore struct {
	db *sql.DB
}

func NewIngredientsStore(host string, port int, user string, password string, dbname string) *ingredientsStore {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	database, err := sql.Open("postgres", psqlconn)
	checkError(err)

	// check db
	err = database.Ping()
	checkError(err)
	return &ingredientsStore{
		db: database,
	}
}

func (i *ingredientsStore) AddThreshold(min, max float32, ingredientsType string, crustSize float32) {
	// insert
	// hardcoded
	addThreshold := `insert into "ingredients_thresholds"("min_value", "max_value", "ingredient_type","for_crust_size") values($1, $2, $3, $4)`
	_, e := i.db.Exec(addThreshold, min, max, ingredientsType, crustSize)
	checkError(e)
}

func (i *ingredientsStore) CleanTable() {
	// insert
	// hardcoded
	truncate := `truncate table ingredients_thresholds`
	_, e := i.db.Exec(truncate)
	checkError(e)
}

func (i *ingredientsStore) GetThresholdsByCrustSize(crustSize float32) []IngredientsThresholds {

	findThresholdByCrustSize := `select min_value, max_value, ingredient_type, for_crust_size from ingredients_thresholds where for_crust_size = $1`
	items := []IngredientsThresholds{}
	rows, err := i.db.Query(findThresholdByCrustSize, crustSize)
	defer rows.Close()
	if err != nil {
		log.Fatalf("%v", err)
	}
	for rows.Next() {
		it := IngredientsThresholds{}
		rows.Scan(&it.MinValue, &it.MaxValue, &it.IngredientType, &it.CrustSize)
		if err != nil {
			log.Fatalf("%v", err)
		}
		items = append(items, it)
	}
	return items

}

func (i *ingredientsStore) Close() {
	i.db.Close()
}
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
