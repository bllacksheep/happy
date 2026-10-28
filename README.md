# happy

cnpg integrated backup-restore validation 

temp instructions for dev
```
source .awscreds
nerdctl run -d --name yolo -p 5432:5432 -e POSTGRES_PASSWORD="test_pass" postgres:15.19-alpine3.24
telnet localhost 5432
export DATABASE_URL="postgres://postgres:test_pass@localhost:5432"
```
