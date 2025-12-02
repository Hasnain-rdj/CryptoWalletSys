import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import {
  LayoutDashboard,
  Send,
  History,
  Box,
  User,
  FileText,
  LogOut,
  Menu,
  X,
} from 'lucide-react';
import { useState } from 'react';
import toast from 'react-hot-toast';

const Navbar = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { logout, userProfile } = useAuth();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  const handleLogout = () => {
    logout();
    toast.success('Logged out successfully');
    navigate('/login');
  };

  const navItems = [
    { path: '/dashboard', icon: LayoutDashboard, label: 'Dashboard' },
    { path: '/send', icon: Send, label: 'Send Money' },
    { path: '/transactions', icon: History, label: 'Transactions' },
    { path: '/blockchain', icon: Box, label: 'Blockchain' },
    { path: '/profile', icon: User, label: 'Profile' },
    { path: '/reports', icon: FileText, label: 'Reports' },
  ];

  const isActive = (path) => location.pathname === path;

  return (
    <nav className="bg-white/80 backdrop-blur-lg shadow-lg sticky top-0 z-50 border-b border-gray-200/50">
      <div className="max-w-7xl mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          {/* Logo with Gradient */}
          <div
            className="flex items-center gap-2 cursor-pointer group"
            onClick={() => navigate('/dashboard')}
          >
            <div className="w-11 h-11 bg-gradient-to-br from-blue-600 via-indigo-600 to-purple-600 rounded-xl flex items-center justify-center shadow-lg transform transition-all duration-300 group-hover:scale-110 group-hover:rotate-6">
              <span className="text-white font-black text-xl">BC</span>
            </div>
            <span className="text-xl font-black bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent hidden sm:block">
              BlockChain Wallet
            </span>
          </div>

          {/* Desktop Navigation with Enhanced Styling */}
          <div className="hidden md:flex items-center gap-1">
            {navItems.map((item) => (
              <button
                key={item.path}
                onClick={() => navigate(item.path)}
                className={`group relative flex items-center gap-2 px-4 py-2 rounded-xl transition-all duration-300 ${
                  isActive(item.path)
                    ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white font-bold shadow-lg'
                    : 'text-gray-600 hover:bg-gradient-to-r hover:from-gray-100 hover:to-blue-50 hover:text-blue-600 font-semibold'
                }`}
              >
                <item.icon className={`w-4 h-4 transition-transform duration-300 ${isActive(item.path) ? '' : 'group-hover:scale-110'}`} />
                <span className="text-sm">{item.label}</span>
                {isActive(item.path) && (
                  <div className="absolute bottom-0 left-0 right-0 h-1 bg-white rounded-full"></div>
                )}
              </button>
            ))}
          </div>

          {/* User Menu with Enhanced Design */}
          <div className="hidden md:flex items-center gap-4">
            <div className="text-right bg-gradient-to-r from-blue-50 to-indigo-50 px-4 py-2 rounded-xl border border-blue-100">
              <p className="text-sm font-bold text-gray-800">
                {userProfile?.fullName || 'User'}
              </p>
              <p className="text-xs text-gray-600">{userProfile?.email || ''}</p>
            </div>
            <button
              onClick={handleLogout}
              className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-red-600 to-red-700 text-white rounded-xl hover:from-red-700 hover:to-red-800 transition-all duration-300 shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 font-semibold"
            >
              <LogOut className="w-4 h-4" />
              Logout
            </button>
          </div>

          {/* Mobile Menu Button with Animation */}
          <button
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            className="md:hidden p-2 rounded-xl hover:bg-gradient-to-r hover:from-blue-50 hover:to-indigo-50 transition-all duration-300 transform hover:scale-110"
          >
            {isMobileMenuOpen ? (
              <X className="w-6 h-6 text-gray-600" />
            ) : (
              <Menu className="w-6 h-6 text-gray-600" />
            )}
          </button>
        </div>

        {/* Mobile Menu with Slide Animation */}
        {isMobileMenuOpen && (
          <div className="md:hidden py-4 border-t border-gray-200 animate-slide-in-left">
            <div className="space-y-2">
              {navItems.map((item, index) => (
                <button
                  key={item.path}
                  onClick={() => {
                    navigate(item.path);
                    setIsMobileMenuOpen(false);
                  }}
                  style={{ animationDelay: `${index * 50}ms` }}
                  className={`w-full flex items-center gap-2 px-4 py-3 rounded-xl transition-all duration-300 animate-fade-in-up ${
                    isActive(item.path)
                      ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white font-bold shadow-lg'
                      : 'text-gray-600 hover:bg-gradient-to-r hover:from-gray-100 hover:to-blue-50'
                  }`}
                >
                  <item.icon className="w-5 h-5" />
                  <span className="font-semibold">{item.label}</span>
                </button>
              ))}
              <div className="pt-4 border-t border-gray-200 mt-4">
                <div className="px-4 pb-3 bg-gradient-to-r from-blue-50 to-indigo-50 rounded-xl p-3 mb-3">
                  <p className="text-sm font-bold text-gray-800">
                    {userProfile?.fullName || 'User'}
                  </p>
                  <p className="text-xs text-gray-600">{userProfile?.email || ''}</p>
                </div>
                <button
                  onClick={() => {
                    handleLogout();
                    setIsMobileMenuOpen(false);
                  }}
                  className="w-full flex items-center gap-2 px-4 py-3 bg-gradient-to-r from-red-600 to-red-700 text-white rounded-xl hover:from-red-700 hover:to-red-800 transition-all duration-300 shadow-lg font-semibold"
                >
                  <LogOut className="w-5 h-5" />
                  Logout
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
