# Chirpy

Chirpy is a simple microblogging platform where users can post, view, and manage short messages called chirps. The application supports user authentication, token-based authorization, and an optional premium membership with exclusive features.

## Features

- **User Authentication:** Secure login and account management using JWT and refresh tokens.
- **Chirp Management:** Create, view, edit (premium users), and delete chirps.
- **Premium Membership:** "Chirpy Red" members get special privileges, such as editing chirps.
- **Webhooks:** Integration with Polka for managing premium memberships.
- **Sorting and Filtering:** Retrieve chirps by author with sorting by creation date.

## Why Use Chirpy?

Chirpy is a minimalistic yet functional microblogging app designed for learning and experimentation. It demonstrates key backend development concepts, including RESTful API design, database migrations, secure authentication, and third-party integrations.

## Technologies Used

- **Go**
- **PostgreSQL**
- **Goose** (Database migrations)
- **SQLC** (Database queries)
- **JWT** (JSON Web Tokens for authentication)
- **HTML** (Frontend structure)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Sacarianos/chirpy.git
   cd chirpy
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up your environment variables by creating a `.env` file in the root directory:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=chirpy
   JWT_SECRET=your_jwt_secret
   POLKA_KEY= your_key
   ```

4. Run database migrations:
   ```bash
   go run cmd/migrate/main.go
   ```

5. Start the server:
   ```bash
   go run cmd/server/main.go
   ```

## Usage

### API Endpoints

#### Authentication
- `POST /api/login`: Login and receive an access token and refresh token.
- `POST /api/refresh`: Refresh your access token using a refresh token.
- `POST /api/revoke`: Revoke a refresh token.

#### Users
- `PUT /api/users`: Update your email and password.

#### Chirps
- `GET /api/chirps`: Retrieve all chirps or filter by author.
- `POST /api/chirps`: Create a new chirp.
- `PUT /api/chirps/{chirpID}`: Edit a chirp (premium users only).
- `DELETE /api/chirps/{chirpID}`: Delete a chirp.

#### Webhooks
- `POST /api/polka/webhooks`: Handle subscription updates from Polka.

## Contributing

Feel free to fork the project and submit pull requests. Contributions are welcome and appreciated!

## License

This project is licensed under the MIT License. See the LICENSE file for details.

---

Thank you for checking out Chirpy!
