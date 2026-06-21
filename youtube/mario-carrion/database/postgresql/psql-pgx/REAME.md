# Create table
```
create table users(name VARCHAR not null, birth_year smallint null default 0);
INSERT INTO users(name, birth_year) VALUES('mario', 1900);
SELECT * FROM users;
go get github.com/jackc/pgx/v5
docker exec -it container_id psql -U user -d dbname
```

[Link](https://www.youtube.com/watch?v=2XCaKYH0Ydo&list=PL7yAAGMOat_EgwoQTvNUflrYL_4qzdB7f)
