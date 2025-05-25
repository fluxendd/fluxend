# Fluxton Frontend

Frontend awesomeness that makes Fluxton whole - built with React and React Router.

## Prerequisites

- Node.js 20.16.0 (recommended to use [Volta](https://volta.sh/) for Node.js version management)
- Yarn package manager

## Getting Started

Follow these steps to set up the project on your local machine:

### 1. Clone the repository

```bash
git clone <repository-url>
cd fluxend-frontend
```

### 2. Install dependencies

This project uses Yarn as the package manager:

```bash
yarn install
```

### 3. Set up environment variables

Copy the sample environment file and modify as needed:

```bash
cp .env.sample .env
```

The following environment variables are required:

```
VITE_FLX_API_BASE_URL=http://fluxton.io/api
VITE_FLX_DEFAULT_ACCEPT_HEADER=application/json
VITE_FLX_DEFAULT_CONTENT_TYPE=application/json
```

Adjust the `VITE_FLX_API_BASE_URL` to point to your local API if needed.

### 4. Start the development server

Run the following command to start the development server:

```bash
yarn dev
```

This will start the React Router development server, and your application will be available at http://localhost:3000 (or the port specified in your environment variables).

## Available Scripts

- `yarn dev` - Starts the development server
- `yarn build` - Builds the application for production
- `yarn start` - Serves the built application
- `yarn typecheck` - Generates React Router types and runs TypeScript type checking

## Project Structure

- `/app` - Main application code
  - `/components` - Reusable UI components
  - `/hooks` - Custom React hooks
  - `/lib` - Utility functions and shared logic
  - `/routes` - Application routes and page components
  - `/services` - API services and data fetching
  - `/tools` - Helper tools and utilities
- `/public` - Static assets

## Technologies

This project uses:

- React 19
- React Router 7
- TypeScript
- Tailwind CSS
- Radix UI Components
- React Hook Form
- Zod for validation
- Tanstack React Query
- Tanstack Table
- Zustand for state management
- Sonner for toast notifications

## Development Notes

- This application uses React Router's SSR capabilities by default (configurable in `react-router.config.ts`)
- Tailwind CSS is used for styling
- Path aliases are configured (`~/*` maps to `./app/*`)
- TypeScript strict mode is enabled

## Deployment

### Using Docker

The project includes a multi-stage Dockerfile optimized for production:

```bash
# Build the Docker image
docker build -t fluxend-frontend .

# Run the container
docker run -p 3000:3000 fluxend-frontend
```

### Manual Deployment

Build the application for production:

```bash
yarn build
```

Then start the production server:

```bash
yarn start
```

You can deploy the built application to any hosting provider that supports Node.js applications.

## License

See the [LICENSE](LICENSE) file for details.