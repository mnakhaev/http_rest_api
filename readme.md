Database configuration (for MacOS)

Useful article - [click](https://postgrespro.ru/docs/postgrespro/9.5/tutorial-createdb) to read
1. Install postgres
2. Run `cd /Library/PostgreSQL/13/bin`
3. Run `sudo -u postgres ./createdb restapi_dev` to create new database
4. Install [migrate](https://www.youtube.com/redirect?redir_token=QUFFLUhqbGpFT05kbGNzLTZORnNseFpsbmFaNWFxUFg2d3xBQ3Jtc0tsRGxfYndfdTNWZjRPQ0p0LVFibmNnZ2JqaXRCUjI4bU04cEFCeExFaDlBb0s4U2FCZW9GNW44SU9Ld0lCczU0ZnJmUVdVT2lXZEwwYms3LUlZQ1BNVzFkZ19lbTBCSU5vblRycXh4QUlqVmZLd1YwSQ%3D%3D&event=comments&q=https%3A%2F%2Fgithub.com%2Fgolang-migrate%2Fmigrate%2Ftree%2Fmaster%2Fcmd%2Fmigrate&stzid=UgztW483m-Oh36Xr4nx4AaABAg) tool
5. Run ` migrate create -ext sql -dir migrations create_users`
6. Define up and down migrations
7. Run ` migrate -path migrations -database "postgres://localhost/restapi_dev?sslmode=disable&user=postgres&password=qwe123QWE" up`
8. Check DB: `sudo -u postgres ./psql -d restapi_dev`
9. Create new DB for tests: `sudo -u postgres ./createdb restapi_test`
10. Run migration for test DB: ` migrate -path migrations -database "postgres://localhost/restapi_test?sslmode=disable&user=postgres&password=qwe123QWE" up`

`model` - keeps all database models  

Store - kind of black-box instance, which provides public methods to work with the data.
It can contain multiple repositories: 
- User repository (create user / find in DB by parameters)
- To be updated

Store -> `config.go` - config for the store

Models - contains models of data representation.