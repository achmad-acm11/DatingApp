
# Dating App

This is a simple Dating App API using Go, Gin, and GORM.

## Project Structure
```bash
dating-app/
├── main.go
├── configs/
│   └── database.go
├── controllers/
│   └── orderController.go
│   └── packageController.go
│   └── userController.go
├── entities/
│   └── matchQueue.go
│   └── order.go
│   └── productPackage.go
│   └── user.go
├── middlewares/
│   └── errorHandler.go
│   └── jwtTokenHandler.go
├── migrations/
│   └── migration.go
├── repositories/
│   └── matchQueueRepository.go
│   └── orderRepository.go
│   └── packageRepository.go
│   └── userRepository.go
├── requests/
│   └── orderRequest.go
│   └── packageRequest.go
│   └── userRequest.go
├── routes/
│   └── orderRoute.go
│   └── packageRoute.go
│   └── userRoute.go
├── seeds/
│   └── packageSeed.go
│   └── seeds.go
│   └── userSeed.go
├── services/
│   └── orderService.go
│   └── packageService.go
│   └── userService.go
├── go.mod
└── README.md
```

## Getting Started

### Clone Repository
```bash
git clone https://github.com/achmad-acm11/DatingApp.git
```

### Environment Variables
To run this project, you will need to add the following environment variables to your .env file.
```bash
DB_USER="username-your-database"
DB_PASSWORD="password-your-database"
DB_NAME="name-your-database"
DB_HOST="hostname-your-database"
DB_PORT="port-your-database"
JWT_SECRET="key-for-jwt-token"
JWT_KEY_ISS="issuer-for-jwt-token"
JWT_EXPIRED="time-expired-for-jwt-token" # in hour units
```

, or you can duplicate `.env.example` and change to `.env` then fill all environment variables.


### Run Application
```bash
go mod tidy
go mod download
go run .
```
By default the application will run on port 8080, you can change whatever you want in the `.env` file of the `APP_PORT` environment variable.

## Tech Stack
* Go 
* Gin Framework
* GORM
* MySQL
* JWT Token
* API Doc _(Postman File are included in the doc folder)_ 

