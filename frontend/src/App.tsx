import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { Provider } from 'react-redux';
import { store } from './store';
import { AuthInitializer } from './components/AuthInitializer';
import { TopicsPage } from './pages/TopicsPage';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';

function App() {
  return (
    <Provider store={store}>
      <AuthInitializer>
        <Router>
          <Routes>
            {/* Public routes */}
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            
            {/* Public routes */}
            <Route path="/topics" element={<TopicsPage />} />
            
            {/* Redirect root to topics */}
            <Route path="/" element={<Navigate to="/topics" replace />} />
            
            {/* Fallback for unknown routes */}
            <Route path="*" element={<Navigate to="/topics" replace />} />
          </Routes>
        </Router>
      </AuthInitializer>
    </Provider>
  );
}

export default App;
