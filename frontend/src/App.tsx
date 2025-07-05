import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { ProtectedRoute } from './components/ProtectedRoute';
import { TopicsPage } from './pages/TopicsPage';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          {/* Public routes */}
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          
          {/* Protected routes */}
          <Route 
            path="/topics" 
            element={
              <ProtectedRoute>
                <TopicsPage />
              </ProtectedRoute>
            } 
          />
          
          {/* Redirect root to topics */}
          <Route path="/" element={<Navigate to="/topics" replace />} />
          
          {/* Fallback for unknown routes */}
          <Route path="*" element={<Navigate to="/topics" replace />} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
