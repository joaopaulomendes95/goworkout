# AGENTS.md - GoWorkout Project Guidelines

This file provides guidelines for AI agents working on this codebase.

## Project Overview

GoWorkout is a full-stack web application for tracking workouts.

- **Backend**: Go with Chi router
- **Frontend**: Svelte 5 + Tailwind CSS (SvelteKit)
- **Database**: PostgreSQL with Goose migrations
- **Authentication**: Custom token-based auth (HttpOnly cookies)

---

## Build Commands

### Go Backend

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Run integration tests
make itest

# Live reload (requires air)
make watch

# Run a single test
go test ./path/to/package -run TestName -v

# Docker development
make docker-run
make docker-down
```

### SvelteKit Frontend

```bash
cd frontend

# Install dependencies
npm install

# Development server
npm run dev

# Build for production
npm run build

# Type checking
npm run check

# Type checking with watch
npm run check:watch

# Formatting (write)
npm run format

# Linting (check only)
npm run lint

# Database (Drizzle)
npm run db:push
npm run db:migrate
npm run db:studio
```

---

## Code Style Guidelines

### Go Backend

**Formatting**
- Use `gofmt` or goimports automatically
- Run `go vet` before committing

**Naming Conventions**
- Use PascalCase for exported functions/types (e.g., `HandleRegisterUser`)
- Use camelCase for private functions/variables
- Use meaningful, descriptive names

**Error Handling**
- Return errors to the caller, don't log and ignore unless truly unrecoverable
- Use `utils.WriteJSON` for consistent JSON responses
- Always handle `nil` pointers (check before use)

**Imports**
- Group imports: standard library → external packages → internal packages
- Use alias if package name conflicts

**Authentication**
- All protected routes use middleware: `s.Middleware.RequireUser(handler)`
- Tokens validated via `Authorization: Bearer <token>` header

### SvelteKit Frontend (Svelte 5)

**CRITICAL: Use Svelte 5 Runes**

```svelte
<!-- WRONG (Svelte 4) -->
<script>
  import { writable } from 'svelte/store';
  let count = writable(0);
  $: doubled = count * 2;
  export let title;
</script>

<!-- CORRECT (Svelte 5) -->
<script>
  let count = $state(0);
  let doubled = $derived(count * 2);
  let { title } = $props();
</script>
```

**Svelte 5 Migration Rules**

| Old (Svelte 4)           | New (Svelte 5)         |
|--------------------------|-----------------------|
| `import { writable }`   | `$state()`            |
| `import { derived }`    | `$derived()`          |
| `$: `                    | `$effect()`           |
| `export let prop`        | `let { prop } = $props()` |
| `bind:value` (prop)     | `$bindable()`         |
| `on:click`               | `onclick`             |
| `on:input`               | `oninput`             |
| `$page` store            | `$props()` from load  |
| `<slot />`               | `@render snippet()`   |

**Never use:**
- `svelte/store` (writable, readable, derived, get)
- `$:` reactive declarations
- `export let` for props
- `on:event` handlers (use onclick, oninput, etc.)

**Data Fetching**
- Use `+page.server.ts` load functions for SSR data
- Use form actions (`+page.server.ts` actions) for mutations
- Pass data via `$props()` in Svelte components

**Example: Server Load Function**
```typescript
// +page.server.ts
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals }) => {
  // Use locals.token for authenticated API calls
  return { data: ... };
};
```

**Example: Server Form Action**
```typescript
// +page.server.ts
export const actions = {
  default: async ({ request, locals }) => {
    const formData = await request.formData();
    // ... process form
    return { success: true };
  }
};
```

**Example: Svelte Component**
```svelte
<script lang="ts">
  let { data, form } = $props();
  let count = $state(0);
  let doubled = $derived(count * 2);
</script>

{#if form?.success}
  <p>Success!</p>
{/if}
```

**Formatting (Prettier)**
- Tabs for indentation
- Single quotes
- Print width: 100
- Use `npm run format` before committing

**TypeScript**
- Never use `any` - use proper types or `unknown`
- Define interfaces for API responses
- Use `$props()` with typed destructuring

**Event Handlers (Svelte 5)**
```svelte
<!-- WRONG -->
<button on:click={handleClick}>Click</button>

<!-- CORRECT -->
<button onclick={handleClick}>Click</button>
```

---

## Architecture

### Authentication Flow

1. User logs in → Go returns `auth_token`
2. SvelteKit stores as HttpOnly cookie
3. `hooks.server.ts` reads cookie → sets `locals.token`
4. All API calls include `Authorization: Bearer <token>` header
5. Go middleware validates token

**Never make API calls from client-side** - always use server load/actions.

### Project Structure

```
cmd/app/main.go       # Entry point
internal/
  api/                # HTTP handlers
  database/           # DB connection
  middleware/         # Auth middleware
  server/             # Server config, routes
  store/              # Data access layer
  tokens/             # Token handling
  utils/              # Helpers
frontend/src/
  routes/             # SvelteKit pages
    (protected)/      # Authenticated routes
  lib/                # Shared code
```

---

## Running Tests

### Go
```bash
# Run all tests
make test

# Run specific test
go test ./internal/api -run TestUserAPI -v

# Integration tests
make itest
```

### Frontend
```bash
cd frontend
npm run check        # Type checking
npm run lint        # Linting
```

---

## Common Issues

- **Svelte 5 warnings about `state_referenced_locally`**: These are informational - the code is correct when initializing form state with `$state(form?.field || '')`.
- **Missing user in request**: Use `middleware.GetUser(r)` after `Authenticate` middleware runs.
- **Cookie not found**: Ensure HttpOnly cookies are set with `httpOnly: true` and `secure` in production.
