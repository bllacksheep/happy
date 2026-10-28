# happy

cnpg integrated backup-restore validation 

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
  db:
    name: "cluster-0001"
    uri: "postgres://postgres:test_pass@localhost:5432"
    backups: "s3://cloudnative-pg-bakups-010skdkv11/backups"
    stragegy: "checkpoint"
  healthchecks:
  - name: postgres_version
    statement: "select version()"
    return: "string"
  - name: database_count
    statement: "select count (*) from pg_database"
    return: "int"
```
