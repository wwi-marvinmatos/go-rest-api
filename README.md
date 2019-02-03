# apiServer

https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql

PostgreSQL getting started on osx
https://www.codementor.io/engineerapart/getting-started-with-postgresql-on-mac-osx-are8jcopb

SQL to create the table

start psql 
create the user
set permissions
create teh db
set db privs
create teh relation

Queries:
SELECT * FROM "products";

clean up table
DELETE FROM products;
ALTER SEQUENCE products_id_seq RESTART WITH 1;

Use the below to test the endpoints
CREATE
curl -i -H 'Content-Type: application/json' -d '{"name":"test product","price":11.22}' http://127.0.0.1:8000/product
curl -i -H 'Content-Type: application/json' -d '{"name":"Air Pods","price":159.99}' http://127.0.0.1:8000/product
curl -i -H 'Content-Type: application/json' -d '{"name":"iPhone 8 Plus","price":759.99}' http://127.0.0.1:8000/product
curl -i -H 'Content-Type: application/json' -d '{"name":"iPhone 8","price":559.99}' http://127.0.0.1:8000/product
curl -i -H 'Content-Type: application/json' -d '{"name":"Macbook Pro 15","price":2559.99}' http://127.0.0.1:8000/product

GET
curl -i http://127.0.0.1:8080/products
curl -i http://127.0.0.1:8080/product/1

DELETE
curl -i -X DELETE http://127.0.0.1:8080/product/5
curl -i http://127.0.0.1:8080/product/5

curl -i -X DELETE http://127.0.0.1:8080/product/4
curl -i http://127.0.0.1:8080/product/4

curl -i http://127.0.0.1:8080/product


CREATE TABLE products (
    id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT products_pkey PRIMARY KEY (id)
);

APP_DB_USERNAME=marvin
APP_DB_PASSWORD=test
APP_DB_NAME=api_server
APP_DB_SSLMODE=disabled


these should be set via the cmd line 
export TEST_DB_USERNAME=marvin
export TEST_DB_PASSWORD=test
export TEST_DB_NAME=api_server
export TEST_DB_SSLMODE=disabled
