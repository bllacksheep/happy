# happy

cnpg integrated backup-restore validation 

spin up a quick cnpg cluster on k3s https://github.com/bllacksheep/cloudnativepg

temp instructions for dev

```bash
# bash
source .awscreds
nerdctl run -d --name yolo -p 5432:5432 -e POSTGRES_PASSWORD="test_pass" postgres:15.19-alpine3.24
telnet localhost 5432
export DB_CONNECTION_URI="postgres://postgres:test_pass@localhost:5432"
```
basic health checks by cluster
```yaml
happy:
  default:
    healthchecks:
      check_cluster_endpoint_health: true
      check_wal_replication_health: true
      check_last_backup_health: true
      check_db_instance_health: true
      check_rpo_health: true
    threshold:
      every: "10s"
      rise: 3
      fall: 3
  clusters:
    - name: "cluster-0001"
      namespace: "db"
      checks:
        - name: postgres_version
          statement: "select version()"
          enabled: false
        - name: database_count
          statement: "select count (*) from pg_database"
          every: "30s"
          rise: 2
          fall: 5
    - name: "cluster-0002"
      namespace: "db-1"
      use_pooler: false
      checks: []
    - name: "cluster-0003"
      namespace: "db-2"
      use_primary: true
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

