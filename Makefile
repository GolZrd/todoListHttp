build:
run:

migrate-up:
	migrate -path migrations -database 'postgres://postgres:qwerty@localhost:5432/mainpet?sslmode=disable' up

migrate-down:
	migrate -path migrations -database 'postgres://postgres:qwerty@localhost:5432/mainpet?sslmode=disable' down