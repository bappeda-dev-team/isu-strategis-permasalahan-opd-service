# API E-KAK 


### Kebutuhan

- Go versi 1.23.0
- PostgreSQL 16
- java 17
- Flyway


### Cara install

- download flyway dan pastikan java versi lebih dari 17
```sh
brew install flyway
```

- buat database bernama db_opd

- add echo framework
```sh
go get github.com/labstack/echo/v4
```

```sh
go get github.com/labstack/echo/v4/middleware
```


### Run server

ketikkan perintah:

```sh
go run main.go
```

untuk hot reload(only macOS/linux)

install
```sh
go install github.com/cespare/reflex@latest
```

```sh
reflex -s -r '\.go$$' go run .
```


untuk menghentikan server, tekan Ctrl + c



### Documentasi API With Swagger Echo
- Untuk menampilkan Ui Swagger
```sh
http://localhost:8001/swagger/index.html
```

- Untuk install swagger echo
```sh
go install github.com/swaggo/swag/cmd/swag@latest
```
```sh
go get -d github.com/swaggo/swag/cmd/swag
```
```sh
go get -u github.com/swaggo/echo-swagger
```

###Maintenance Database
- Untuk menambahkan table,buat file di folder db/migrations dan gunakan format yg ditentukan flyway
```sh
V1__create_table_a
V2__create_table_b
```