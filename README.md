
## 💬 Financial Chat - Multi-Room

This project is a **multi-room chat system** with a Go backend and a static HTML/JS frontend. It allows users to:

- Log in and register.
- Join multiple chat rooms simultaneously.
- Send messages and commands to interact with a financial bot.

> The frontend is a static HTML file (`index.html`), so we need a **simple HTTP server** to serve it. Opening the file directly with `file://` causes WebSocket connections to close automatically.


## 🚀 How to Run

#### 1. **Create a `.env` file** in the project root with all required environment variables for the services.
You can use `env.example` as a reference or for testing purposes. Example:

`
JWT_SECRET=your_secret_here
DATABASE_URL=postgres://user:password@db:5432/chat
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/`




#### 2. In the terminal, from the project root, run the following command to start backend and frontend:

Make sure you have Docker and Python 3 installed.

``
docker compose up --build -d && (cd ./chat-frontend && python3 -m http.server 3000)``

Open the frontend in your browser:
`http://localhost:3000`

Log in or register a new user.
Create or join chat rooms and start sending messages.

## 🔹 Technologies Used
Backend: Go (Golang)

WebSocket handling: Gorilla WebSocket package

Frontend: HTML, CSS, JavaScript (vanilla)

Containerization: Docker, Docker Compose

Bot communication: RabbitMQ (message broker)

Authentication: JWT (JSON Web Tokens)