A simple bank
* create and manage account
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
* execute the command ```make dropdb``` to drop the database
* execute the command ```make migrateup``` to run the migration
* execute the command ```make migrateup``` to run the migration