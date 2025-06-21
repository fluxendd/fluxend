# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

- `yarn dev` - Start development server
- `yarn build` - Build for production
- `yarn start` - Serve built application  
- `yarn typecheck` - Generate React Router types and run TypeScript checking

## Architecture Overview

This is a React Router 7 SSR application with the following key architectural patterns:

### Routing Structure
- Uses React Router's file-based routing with programmatic route configuration in `app/routes.ts`
- Two main layout levels:
  1. `app-layout.tsx` - Main authenticated app wrapper
  2. `project-layout.tsx` - Project-specific context wrapper
- Custom `routeFolder` utility (`app/lib/router.ts`) for dynamic route discovery with optional sidebar components

### Authentication & Services
- Token-based authentication using secure HTTP-only cookies (`sessionCookie`)
- Service layer pattern with `initializeServices()` factory that creates API service instances
- Auth token passed down through layout loaders and used to initialize services
- Server/client auth token utilities in `app/lib/auth.ts`

### State Management
- React Query (TanStack Query) for server state management
- Zustand for client-side state management
- Services layer handles all API interactions with consistent error handling

### UI Architecture
- Shadcn/ui components with Radix UI primitives
- Tailwind CSS with custom theme system
- Sidebar-based navigation using Shadcn sidebar components
- Path aliases: `~/*` maps to `./app/*`

### Key Patterns
- Loader functions for SSR data fetching with error boundaries
- Outlet context pattern for passing data between layouts and routes
- Custom hooks for common functionality (mobile detection, theme, data filtering)
- Service injection pattern through layout loaders

## Project Structure Notes

- `/routes/` - Page components with nested routing
- `/components/shared/` - Layout and shared components  
- `/components/ui/` - Shadcn UI components
- `/services/` - API service layer
- `/lib/` - Utilities and shared logic
- `/hooks/` - Custom React hooks

## Environment Setup

Required environment variables:
```
VITE_FLX_API_BASE_URL=http://fluxend.app/api
VITE_FLX_DEFAULT_ACCEPT_HEADER=application/json
VITE_FLX_DEFAULT_CONTENT_TYPE=application/json
```

## Key Technologies

- React 19 with React Router 7 (SSR enabled)
- TypeScript with strict mode
- Tailwind CSS v4 with Shadcn/ui components
- TanStack Query + TanStack Table + TanStack Virtual
- React Hook Form with Zod validation
- Zustand for client state
- Node.js 20.16.0 (Volta managed)