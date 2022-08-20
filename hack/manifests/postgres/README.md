# Install

```shell
helm upgrade --install postgres --namespace postgres -f values-default.yaml -f values-local.yaml .
```

This will deploy postgres

##  Connect

Exec into the pod and then execute the following

```shell
psql -d postgres -U postgresadmin
```

##  Prepare for temporal

Port forward the servive port of postgres

```shell
kubectl -n postgres port-forward svc/postgres-service 5432:5432
```

```shell
export SQL_PLUGIN=postgres
export SQL_HOST=localhost
export SQL_PORT=5432
export SQL_USER=postgresadmin
export SQL_PASSWORD=admin123

temporal-sql-tool create-database -database temporal
SQL_DATABASE=temporal temporal-sql-tool setup-schema -v 0.0
SQL_DATABASE=temporal temporal-sql-tool update -schema-dir schema/postgresql/v96/temporal/versioned

temporal-sql-tool create-database -database temporal_visibility
SQL_DATABASE=temporal_visibility temporal-sql-tool setup-schema -v 0.0
SQL_DATABASE=temporal_visibility temporal-sql-tool update -schema-dir schema/postgresql/v96/visibility/versioned
```

## Install temporal

```shell
kubectl create ns temporal

helm dependency build

helm install -n temporal -f values/values.postgresql.yaml temporal \
  --set elasticsearch.enabled=false \
  --set server.config.persistence.default.sql.user=postgresadmin \
  --set server.config.persistence.default.sql.password=admin123 \
  --set server.config.persistence.visibility.sql.user=postgresadmin \
  --set server.config.persistence.visibility.sql.password=admin123 \
  --set server.config.persistence.default.sql.host=postgres-service.postgres \
  --set server.config.persistence.visibility.sql.host=postgres-service.postgres . --timeout 900s
```  