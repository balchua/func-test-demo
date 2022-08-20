package main

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	fdb "github.com/balchua/func-test-demo/pkg/datastore"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/spf13/pflag"
)

type basket struct {
	crustSize  float32
	tomatoes   float32
	pineapples float32
	hams       float32
	dsConfig   datastoreConfig
}

func (b *basket) withCrustSize(crustSize float32) {
	b.crustSize = crustSize
}

func (b *basket) withIngredientsThresholds(ingredients *godog.Table) error {
	ds := fdb.NewIngredientsStore(b.dsConfig.host, b.dsConfig.port, b.dsConfig.user, b.dsConfig.password, b.dsConfig.dbName)
	header := ingredients.Rows[0].Cells
	for i := 1; i < len(ingredients.Rows); i++ {
		var min, max, crustSize float64
		var ingredientsType string
		var err error
		for n, cell := range ingredients.Rows[i].Cells {
			switch header[n].Value {
			case "min":
				if min, err = strconv.ParseFloat(cell.Value, 32); err != nil {
					return fmt.Errorf("invalid conversion: %v", err)
				}
			case "max":
				if max, err = strconv.ParseFloat(cell.Value, 32); err != nil {
					return fmt.Errorf("invalid conversion: %v", err)
				}
			case "crust size":
				if crustSize, err = strconv.ParseFloat(cell.Value, 32); err != nil {
					return fmt.Errorf("invalid conversion: %v", err)
				}
			case "ingredient_type":
				ingredientsType = cell.Value
			default:
				return fmt.Errorf("unexpected column name: %s", header[n].Value)
			}
		}
		ds.AddThreshold(float32(min), float32(max), ingredientsType, float32(crustSize))

	}
	ds.Close()
	return nil
}

func (b *basket) withIngredients(tomato, pineapples, hams float32) {
	b.hams = hams
	b.pineapples = pineapples
	b.tomatoes = tomato
}

func (b *basket) ofQuality(status string) error {
	p := NewPizza(b.dsConfig, b.crustSize, b.pineapples, b.hams, b.tomatoes)
	result := p.CookPizza()
	if status == result {
		return nil
	}
	return fmt.Errorf("expected result to be: %s, but actual is: %s", status, result)
}

func (b *basket) cleanup() {
	ds := fdb.NewIngredientsStore("localhost", 30432, "postgresadmin", "admin123", "postgres")
	ds.CleanTable()
	ds.Close()
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	d := datastoreConfig{
		host:     "localhost",
		port:     30432,
		user:     "postgresadmin",
		password: "admin123",
		dbName:   "postgres",
	}

	b := &basket{
		dsConfig: d,
	}

	ctx.BeforeScenario(func(sc *godog.Scenario) {
		fmt.Println("do anything before starting the scenario")
		b.cleanup()
	})
	//Given
	ctx.Step(`^the following thresholds$`, b.withIngredientsThresholds)

	//When
	ctx.Step(`^the crust size is (\d+) inches$`, b.withCrustSize)
	ctx.Step(`^the ingredients "([^"]*)", "([^"]*)", "([^"]*)"$`, b.withIngredients)
	// Then
	ctx.Step(`^it should be a "([^"]*)" pizza$`, b.ofQuality)
}

var opts = godog.Options{Output: colors.Colored(os.Stdout)}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	pflag.Parse()
	opts.Paths = pflag.Args()

	status := godog.TestSuite{
		Name:                "pizza",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	os.Exit(status)
}
