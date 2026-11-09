# happy

cnpg integrated backup-restore validation, checks for semantic integraity of data based on your checks

spin up a quick cnpg cluster on k3s https://github.com/bllacksheep/cloudnativepg

temp instructions for dev (always changing)

```bash
# bash
source .awscreds
nerdctl run -d --name yolo -p 5432:5432 -e POSTGRES_PASSWORD="test_pass" postgres:15.19-alpine3.24
telnet localhost 5432
export DB_CONNECTION_URI="postgres://postgres:test_pass@localhost:5432"
```
basic health checks by cluster
will eventually be custom resource definitions
crd for Happy, RPOTestScheudle and Check or Test
```yaml
# minimal example
happy:
  clusters:
    - name: "cluster-0001"
      namespace: "db"
      backup_test:
        enabled: true
        schedule: "* * * * *"
        target_namespace: "db"
      tests:
        - name: database_count
          query: "select count (*) from pg_database"
          expect: 5
```
mess with defaults
```yaml
# more complete example
happy:
  checks:
    - name: postgres_version
      query: "select version()"
      expect: "some version"
      enabled: true
  check_cluster_endpoint: true
  check_wal_replication: true
  check_last_backup: true
  check_db_instance: true
  check_rpo_health: true
  threshold:
    every: "10s"
    rise: 3
    fall: 3
  clusters:
    - name: "cluster-0001"
      namespace: "db-1"
      use_defaults: true
      backup_test:
        enabled: true
        target_namespace: "db"
      use_pooler: true
      check_rpo_health: false
      exempt_from:
        - postgres_version
      checks:
        - name: database_count
          query: "select count (*) from pg_database"
          expect: 5
    - name: "cluster-0002"
      namespace: "db-2"
      use_pooler: false
      test_backups: false
      use_defaults: false
      checks: []

# defaults, true values not required
# use_pooler true     when false, falls back to use_secondary
# use_primary false
# use_secondary false
# use_instance string by name
# enabled true
```
checking cluster status manually from pod
```
k exec happy -n happy-system -it -- sh
curl --cacert $CACERT -H "Authorization: Bearer $TOKEN" https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_SERVICE_PORT/apis/postgresql.cnpg.io/v1/namespaces/db/clusters/cluster-0001
```

```
nerdctl build --tag maclighiche/happy:0.0.x .
nerdctl push maclighiche/happy:0.0.x

test
kubectl port-forward svc/pooler-rw 5432:5432 -n db
```

