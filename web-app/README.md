# crypto-balance-bot

`crypto-balance-bot` is a subscription system designed to notify users about movements in their subscribed blockchain accounts. The application is built using Next.js for the frontend and Nest.js for the backend, providing a seamless experience for managing blockchain account subscriptions.

## Features

- User authentication with sign-up and login functionality.
- User management page for handling blockchain account subscriptions.
- Notifications for balance movements in subscribed accounts.
- Support for multiple blockchains, including Bitcoin and Ethereum.

## Project Structure

- **backend/**: Contains the Nest.js application.
  - **src/**: Source code for the backend.
    - **auth/**: Authentication module for user login and registration.
    - **users/**: Module for managing user data.
    - **blockchain/**: Module for handling blockchain subscriptions.
  - **package.json**: Backend dependencies and scripts.
  - **nest-cli.json**: Configuration for the Nest CLI.
  - **tsconfig.json**: TypeScript configuration for the backend.

- **frontend/**: Contains the Next.js application.
  - **pages/**: Application pages including landing, login, signup, and user management.
  - **components/**: Reusable components such as Layout and Navbar.
  - **styles/**: Global styles for the application.
  - **package.json**: Frontend dependencies and scripts.
  - **tsconfig.json**: TypeScript configuration for the frontend.

- **docker-compose.yml**: Configuration for orchestrating backend and frontend services.

## Getting Started

To get started with the project, follow these steps:

1. Clone the repository:
   ```
   git clone <repository-url>
   cd crypto-balance-bot
   ```

2. Set up the backend:
   - Navigate to the `backend` directory and install dependencies:
     ```
     cd backend
     npm install
     ```

3. Set up the frontend:
   - Navigate to the `frontend` directory and install dependencies:
     ```
     cd frontend
     npm install
     ```

4. Run the application using Docker:
   ```
   docker-compose up
   ```

5. Access the application at `http://localhost:3000`.

## TODO

- [ ] Implement additional blockchain support.
- [ ] Enhance user interface for better user experience.
- [ ] Add more notification options for users.