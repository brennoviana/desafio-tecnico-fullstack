import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { Provider } from 'react-redux';
import { store } from './store';
import { AuthInitializer } from './components/AuthInitializer';
import { Dashboard } from './pages/Dashboard';
import { VotingScreen } from './pages/VotingScreen';
import { ResultsScreen } from './pages/ResultsScreen';
import { SessionManager } from './pages/SessionManager';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';

function App() {
  return (
    <Provider store={store}>
      <AuthInitializer>
        <Router>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/topic/:topicId/vote" element={<VotingScreen />} />
            <Route path="/topic/:topicId/results" element={<ResultsScreen />} />
            <Route path="/topic/:topicId/manage" element={<SessionManager />} />
            
            <Route path="/" element={<Navigate to="/dashboard" replace />} />
            <Route path="/topics" element={<Navigate to="/dashboard" replace />} />
            
            <Route path="*" element={<Navigate to="/dashboard" replace />} />
          </Routes>
        </Router>
      </AuthInitializer>
    </Provider>
  );
}

export default App;
