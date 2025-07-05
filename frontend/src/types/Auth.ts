export interface User {
  id: number;
  cpf: string;
  name: string;
}

export interface LoginCredentials {
  cpf: string;
  password: string;
}

export interface RegisterCredentials {
  name: string;
  cpf: string;
  password: string;
}

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  loading: boolean;
}

export interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  loading: boolean;
  login: (credentials: LoginCredentials) => Promise<void>;
  register: (credentials: RegisterCredentials) => Promise<void>;
  logout: () => void;
} 