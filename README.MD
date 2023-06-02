A simple bank
* Create and manage account
Owner, balance, currency
* Record all balance changes
Create an account entry for each change
* Money transfer transaction
Perform money transfer between two accounts consistently within a transaction

### How to create the migration
1. Install migration https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
2. Execute the command ```migrate create -ext sql -dir db/migration -seq init_schema```

### How to run the migration on windows: 
#### Download MAKE
* Execute the command ```make postgres``` to run postgres on docker
* Execute the command ```make createdb``` to create the database
* Execute the command ```make dropdb``` to drop the database
* Execute the command ```make migrateup``` to run the migration
* Execute the command ```make migratedown``` to undo the migration

### sqlc: 
You can install it: https://docs.sqlc.dev/en/latest/overview/install.html
Or run on docker: 
1. ```docker pull kjconroy/sqlc```
2. ```docker run --rm -v "${PWD}:/src" -w /src kjconroy/sqlc generate```
3. ```docker run --rm -v "$(Get-Location):/src" -w /src kjconroy/sqlc generate```