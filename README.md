This is gator, a  little learning project. It seemed interesting to me as it let me practice using SQL in a practical way, with other languages.

Gator requires Postgres and Go to run.

You can use go install at the root of the repo to get it set up.

As for Postgres, you'll need to provide the program with your Postgres instance, it should look something like "postgres://USER:PASSWORD@HOST:PORT/DBNAME?sslmode=disable"

You should set db_url by doing "export DB_URL=postgres://postgres:postgres@localhost:5432/gator?sslmode=disable
"

