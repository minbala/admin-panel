## Assessment Test

## API Documentation (Swagger)

To generate swagger documentation

```sh
$ cd back_end
$ swag fmt
$ swag init --parseDependency --parseInternal
```

Read documentation at [doc](http://localhost:8080/docs/index.html):

To do database migration

First change Database Setting in conf/app.env accordingly and create necessary database

```sh
$ cd back_end
$ go run main.go - migrate
```
default account email => minbala33@gmail.com
password => minbala33

To run the backend application

```sh
$ cd back_end
$ go run main.go - serve
```

to run test cases

```sh
$ cd back_end
$ go test ./...
```
To run front end
```sh
$ cd front_end
$ npm install
$ npm run dev

```# admin-panel
