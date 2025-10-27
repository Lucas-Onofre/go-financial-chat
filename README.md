💬 Financial Chat - Multi-Room

This project is a multi-room chat system with a Go backend and a static HTML/JS frontend. It allows users to:

Log in and register.

Join multiple chat rooms simultaneously.

Send messages and commands to interact with a financial bot.

Visually differentiate messages from the bot, users, and errors.

The frontend is a static HTML file (index.html), so we need a simple HTTP server to serve it. Opening the file directly with file:// causes WebSocket connections to close automatically.

🚀 How to Run

Make sure you have Docker and Python 3 installed.

Create a .env file in the project root with all required environment variables for the backend services (e.g., database credentials, JWT secret, broker URLs). Example:

In the terminal, from the project root, run the following command to start backend and frontend:

docker compose up --build -d && \
(cd ./chat-frontend && python3 -m http.server 8080)


Open the frontend in your browser:

http://localhost:8080


Log in or register a new user.

Create or join chat rooms and start sending messages.

🔹 Notes

Frontend: served by the Python HTTP server only to bypass browser security restrictions for WebSockets.

Backend: all services (chat, bot, etc.) are defined in the docker-compose.yml.