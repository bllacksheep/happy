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
# config
happy:
  cluster:
    name: "cluster-0001"
    namespace: "db"
    uri: "postgres://postgres:test_pass@localhost:5432"
    backups: "s3://cloudnative-pg-bakups-010skdkv11/backups"
    strategy: "checkpoint"
  healthchecks:
    global:
      every: 30s
      rise: 3
      fall: 3
    # explicitly on by default, use this to turn off
    happy_checks:
      check_cluster_endpoint_health: true
      check_wal_replication_health: true
      check_last_backup_health: true
      check_db_instance_health: true
      check_rpo_health: true
    custom_checks:
      - name: postgres_version
        statement: "select version()"
        every: 10s
        rise: 3
        fall: 3
      - name: database_count
        statement: "select count (*) from pg_database"
        every: 10s
        rise: 3
        fall: 3
```
checking cluster status manually from pod
```
k exec happy -n happy-system -it -- sh
curl --cacert $CACERT -H "Authorization: Bearer $TOKEN" https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_SERVICE_PORT/apis/postgresql.cnpg.io/v1/namespaces/db/clusters/cluster-0001

```
