# Demo functional test with Godog

This repository is an example of how to use cucumber/godog to write Behavior Driven Test with Gherkin Language.

## Pre-requisites

### Kubernetes cluster

You can install MicroK8s to quickly standup a kubernetes cluster.

``` shell
sudo snap install microk8s --channel 1.24/stable
microk8s status --wait-ready

microk8s enable dns rbac

```

### Postgres

Once the cluster is fully ready, install Postgres using the charts located [here](hack/manifests/postgres/)

#### Create the namespace

``` shell
kubectl create ns func-test
```

#### Install Postgres

If your system do not have `HugePages-2Mi` enabled.

``` shell
kubectl create ns func-test
helm upgrade --install --namespace func-test postgres -f hack/manifests/postgres/values-default.yaml hack/manifests/postgres/
```
But if your system has `HugePages-2Mi` enabled use the `values-hugepages.yaml`

``` shell
helm upgrade --install --namespace func-test postgres -f hack/manifests/postgres/values-default.yaml -f hack/manifests/postgres/values-hugepages.yaml hack/manifests/postgres/
```

#### Check Postgres

After the helm install, wait for Postgres to be available

``` shell
kubectl wait deployment -n func-test postgres --for condition=Available=True --timeout=90s
```

## Running the functional test

The feature is defined in the [file](features/make_pizza.feature)

In this particular test, we are going to detect whether the hawaiian pizza is created perfectly or not.

First, the perfect ingredients of a pizza are stored in a database table called `ingredients_threshold`

Currently defined as something like this 

``` sql
CREATE TABLE IF NOT EXISTS ingredients_thresholds (
   id SERIAL PRIMARY KEY,
   ingredient_type VARCHAR(10) NOT NULL,
   min_value NUMERIC(4,2) NOT NULL,
   max_value NUMERIC(4,2) NOT null,
   for_crust_size NUMERIC(4,2) not NULL
)
```

The columns :

`min_value` : defines the minimum parts of an ingredient to add.
`max_value` : defines the maximum parts of an ingredient to add.
`for_crust_size` : indicates the ingredients for a particular size of a crust.
`ingredient_type` : Defined as `H` for `Ham`, `T` for `tomato`and `P` for `pineapple`

To determine whether the hawaiian pizza is perfect, is that the chef must only add the type of ingredient within the range of the `min_value` and `max_value`

### Populating the table

Before starting the test, the data of the table `ingredients_thresholds` is truncated using the `BeforeScenario`.

Example:

``` go
	ctx.BeforeScenario(func(sc *godog.Scenario) {
		fmt.Println("do anything before starting the scenario")
		b.cleanup()
	})
```
The data is populated at every `Given` step.

Example in the feature file.

```
| min | max | ingredient_type | crust size |
| 10 | 30   | H | 12|
| 10 | 30   | P | 12|
| 0.5 | 1.0   | T | 12|
| 5 | 15   | H | 10|
| 10 | 15   | P | 10|
| 0.25 | 0.55   | T | 10|
```

The code inserts the data into the `ingredients_thresholds` table, via this piece of code. 

``` go
func (b *basket) withIngredientsThresholds(ingredients *godog.Table) error
```
### Parsing the functional Tests

``` go

	//Given
	ctx.Step(`^the following thresholds$`, b.withIngredientsThresholds)

	//When
	ctx.Step(`^the crust size is (\d+) inches$`, b.withCrustSize)
	ctx.Step(`^the ingredients "([^"]*)", "([^"]*)", "([^"]*)"$`, b.withIngredients)
	// Then
	ctx.Step(`^it should be a "([^"]*)" pizza$`, b.ofQuality)
```
Please refer to [make_pizza_test.go](make_pizza_test.go)
