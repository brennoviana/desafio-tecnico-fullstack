# Redux Store Structure

## Overview
This application uses Redux Toolkit for state management, specifically managing authentication state including:
- User information (name, CPF)
- Authentication token
- Authentication status
- Loading states
- Error handling

## Store Structure

### Authentication Slice (`authSlice.ts`)
Manages all authentication-related state and actions:

#### State
```typescript
interface AuthState {
  user: User | null;        // Current logged-in user
  token: string | null;     // JWT authentication token
  isAuthenticated: boolean; // Authentication status
  loading: boolean;         // Loading state for async operations
  error: string | null;     // Error messages
}
```

#### Actions
- `loginUser` - Async thunk for user login
- `registerUser` - Async thunk for user registration
- `logoutUser` - Async thunk for user logout
- `clearError` - Clear error messages
- `initializeAuth` - Initialize auth state from localStorage

### Store Configuration (`index.ts`)
- Configures the Redux store with the auth slice
- Exports typed `RootState` and `AppDispatch` types

### Typed Hooks (`../hooks/redux.ts`)
- `useAppDispatch()` - Typed dispatch hook
- `useAppSelector()` - Typed selector hook

## Usage Examples

### Login
```typescript
const dispatch = useAppDispatch();
const { loading, error } = useAppSelector((state) => state.auth);

// Login user
dispatch(loginUser({ cpf: '12345678901', password: 'password' }));
```

### Check Authentication
```typescript
const { isAuthenticated, user } = useAppSelector((state) => state.auth);

if (isAuthenticated) {
  console.log(`Welcome ${user?.name}`);
}
```

### Logout
```typescript
const dispatch = useAppDispatch();
dispatch(logoutUser());
```

## Components Integration

### AuthInitializer
Initializes authentication state on app startup by checking localStorage for existing tokens.

### Pages
- `LoginPage` - Uses `loginUser` action
- `RegisterPage` - Uses `registerUser` action  
- `TopicsPage` - Uses `logoutUser` action and auth selectors

### Protected Routes
- `ProtectedRoute` - Uses auth selectors to check authentication status

## Token Management
- JWT tokens are stored in localStorage
- Tokens are automatically included in API requests
- Auth state is initialized from localStorage on app startup
- Tokens are cleared on logout 