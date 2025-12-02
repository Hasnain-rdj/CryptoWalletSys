import { createContext, useContext, useState, useEffect } from 'react';
import api from '../utils/api';
import toast from 'react-hot-toast';

const AuthContext = createContext();

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};

export const AuthProvider = ({ children }) => {
  const [currentUser, setCurrentUser] = useState(null);
  const [userProfile, setUserProfile] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Check if user has token in localStorage
    const token = localStorage.getItem('token');
    const savedUser = localStorage.getItem('user');

    if (token && savedUser) {
      setCurrentUser(JSON.parse(savedUser));
      setUserProfile(JSON.parse(savedUser));
      
      // Optionally refresh profile from backend
      refreshProfile().catch(() => {
        // If token is invalid, clear auth
        logout();
      });
    }
    
    setLoading(false);
  }, []);

  const register = async (email, password, fullName, cnic) => {
    try {
      const response = await api.post('/register', {
        email,
        password,
        fullName,
        cnic,
      });

      // Save token and user data
      const { token, user, privateKey } = response.data;
      localStorage.setItem('token', token);
      localStorage.setItem('user', JSON.stringify(user));
      
      setCurrentUser(user);
      setUserProfile(user);

      toast.success('Registration successful! Please save your private key.');
      
      return {
        success: true,
        privateKey,
        user,
      };
    } catch (error) {
      const message = error.response?.data?.error || 'Registration failed';
      toast.error(message);
      throw new Error(message);
    }
  };

  const login = async (email, password) => {
    try {
      const response = await api.post('/login', { email, password });
      
      const { token, user } = response.data;
      localStorage.setItem('token', token);
      localStorage.setItem('user', JSON.stringify(user));
      
      setCurrentUser(user);
      setUserProfile(user);
      
      toast.success('Login successful!');
      return user;
    } catch (error) {
      const message = error.response?.data?.error || 'Login failed';
      toast.error(message);
      throw error;
    }
  };

  const logout = async () => {
    try {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      localStorage.removeItem('privateKey'); // Also remove private key on logout
      
      setCurrentUser(null);
      setUserProfile(null);
      
      toast.success('Logged out successfully');
    } catch (error) {
      toast.error('Logout failed');
      throw error;
    }
  };

  const refreshProfile = async () => {
    try {
      const response = await api.get('/profile');
      const user = response.data.user;
      
      setUserProfile(user);
      setCurrentUser(user);
      localStorage.setItem('user', JSON.stringify(user));
      
      return user;
    } catch (error) {
      console.error('Failed to refresh profile:', error);
      throw error;
    }
  };

  const value = {
    currentUser,
    userProfile,
    loading,
    register,
    login,
    logout,
    refreshProfile,
  };

  return (
    <AuthContext.Provider value={value}>
      {!loading && children}
    </AuthContext.Provider>
  );
};

export default AuthContext;
