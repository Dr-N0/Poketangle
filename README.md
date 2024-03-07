# Poketangle

This is a basic Go web server with protected routes using a secret token for authentication.

## Prerequisites

- [Go](https://golang.org/dl/) installed on your machine.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/Dr-N0/Poketangle.git
   cd Poketangle
   ```

2. Set the environment variable for your secret token:

   ```bash
   export MY_SECRET_TOKEN=your-secret-token
   ```

   Replace "your-secret-token" with your actual secret token.

3. Run the Go program:

   ```bash
   go run main.go
   ```

   The server will start on port 8080.

## Usage

Visit the following routes in your browser or use a tool like `curl` to test the API:

- Home: [http://localhost:8080/](http://localhost:8080/)
- Hello: [http://localhost:8080/hello](http://localhost:8080/hello)
- Protected Guess Route: [http://localhost:8080/guess](http://localhost:8080/guess)

To access the protected "/guess" route, include the secret token in the `Authorization` header:

```bash
curl -X POST  -H "Content-Type: application/json"  -H "Authorization: Bearer your-secret-token"  -d '{"pokemon":"Pikachu","question":"is it a monotype?"}'  http://localhost:8080/guess
```

Replace "your-secret-token" with your actual secret token.