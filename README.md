# Airways2

Airways2 is a REST API backend application for managing flights, users, airports, cities, and reviews. The project is written in Go using the `gin` framework and supports authentication, authorization, advanced filtering, and full CRUD operations for all key resources.

## Main Features

- **CRUD operations** for airplanes, flights, airports, cities, and users
- **Authentication and authorization** (registration, login, route protection)
- **Flight check-in** and flight reviews
- **Advanced search and filtering** (by price, city, date, airplane, etc.)
- **User-flight connection** (purchase/reservation, flight history)
- **RESTful API** structure with clear separation of concerns

## Project Structure

- `main.go` - Application entry point; initializes database, Redis, migrations, and routes
- `api/` - Logic for each domain (flights, users, airports, airplanes, cities)
- `routes/` - REST route definitions, grouping endpoints by entity
- `repository/` - Database access and migrations
- `redis/` - Redis cache system initialization and usage

## Technologies

- Go (Golang)
- Gin (web framework)
- SQL (database)
- Redis (cache/sessions)
- JWT (authentication)

## Installation and Running

1. Clone the repository:
   ```bash
   git clone https://github.com/AmnaGrcic/airways2.git
   cd airways2
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run locally:
   ```bash
   go run main.go
   ```
4. The API will be available at [http://localhost:8080](http://localhost:8080)

## Example API Endpoints

- `/users/register` - User registration
- `/users/login` - User login
- `/flights` - Manage flights (list, add, update, delete)
- `/airports` - Manage airports
- `/aeroplanes` - Manage airplanes
- `/cities` - Manage cities
- `/usersflights` - View user's flights, add reviews

## Author

Amna Grcic
