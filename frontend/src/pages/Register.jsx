import { useState, useEffect } from 'react';
import { useNavigate, Link, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import api from '../utils/api';
import { Wallet, Mail, Lock, User, CreditCard, Eye, EyeOff, Key, CheckCircle, AlertCircle } from 'lucide-react';
import toast from 'react-hot-toast';

const Register = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { register } = useAuth();
  const [formData, setFormData] = useState({
    fullName: '',
    email: location.state?.email || '',
    cnic: '',
    password: '',
    confirmPassword: '',
  });
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [showPrivateKey, setShowPrivateKey] = useState(false);
  const [privateKey, setPrivateKey] = useState('');
  const [emailVerified, setEmailVerified] = useState(location.state?.verified || false);
  const [showVerificationSuccess, setShowVerificationSuccess] = useState(false);

  useEffect(() => {
    // If email is provided from verification, mark as verified
    if (location.state?.verified && location.state?.email) {
      setEmailVerified(true);
      setFormData((prev) => ({ ...prev, email: location.state.email }));
      setShowVerificationSuccess(true);
      // Hide success message after 5 seconds
      setTimeout(() => setShowVerificationSuccess(false), 5000);
    }
  }, [location.state]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
    
    // If email changes, reset verification status
    if (name === 'email' && value !== location.state?.email) {
      setEmailVerified(false);
    }
  };

  const handleVerifyEmail = async () => {
    if (!formData.email) {
      toast.error('Please enter your email address');
      return;
    }

    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      toast.error('Please enter a valid email address');
      return;
    }

    // Send OTP before navigating
    try {
      const res = await api.post('/otp/generate', { email: formData.email });
      toast.success('OTP sent to your email!');
      // Show dev OTP if in development mode
      if (res.data.otp) {
        toast.success(`Dev OTP: ${res.data.otp}`, { duration: 10000 });
      }
      // Navigate to verification page
      navigate('/verify-email', { state: { email: formData.email } });
    } catch (error) {
      const message = error.response?.data?.error || 'Failed to send OTP';
      toast.error(message);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!emailVerified) {
      toast.error('Please verify your email first');
      return;
    }

    // Validate CNIC format
    const cnicRegex = /^\d{5}-\d{7}-\d{1}$/;
    if (!cnicRegex.test(formData.cnic)) {
      toast.error('Invalid CNIC format. Expected: 12345-1234567-1');
      return;
    }

    if (formData.password !== formData.confirmPassword) {
      toast.error('Passwords do not match');
      return;
    }

    if (formData.password.length < 6) {
      toast.error('Password must be at least 6 characters');
      return;
    }

    // Validate password strength
    const hasLetter = /[a-zA-Z]/.test(formData.password);
    const hasNumber = /[0-9]/.test(formData.password);
    if (!hasLetter || !hasNumber) {
      toast.error('Password must contain at least one letter and one number');
      return;
    }

    setLoading(true);
    try {
      const result = await register(
        formData.email,
        formData.password,
        formData.fullName,
        formData.cnic
      );
      
      if (result.success) {
        setPrivateKey(result.privateKey);
        // Save private key to localStorage
        localStorage.setItem('privateKey', result.privateKey);
        setShowPrivateKey(true);
      }
    } catch (error) {
      console.error('Registration error:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCopyPrivateKey = () => {
    navigator.clipboard.writeText(privateKey);
    toast.success('Private key copied to clipboard!');
  };

  const handleContinueToDashboard = () => {
    navigate('/dashboard');
  };

  if (showPrivateKey) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-indigo-50 via-blue-50 to-purple-50 flex items-center justify-center p-4 relative overflow-hidden">
        {/* Animated background elements */}
        <div className="absolute inset-0 overflow-hidden pointer-events-none">
          <div className="absolute -top-40 -right-40 w-80 h-80 bg-green-300 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-blob"></div>
          <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-blue-300 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-blob animation-delay-2000"></div>
        </div>

        <div className="bg-white/90 backdrop-blur-lg rounded-3xl shadow-3xl p-8 max-w-2xl w-full animate-scale-in relative z-10 border border-white/20">
          <div className="text-center mb-6">
            <div className="inline-flex items-center justify-center w-20 h-20 bg-gradient-to-br from-green-500 to-emerald-600 rounded-2xl mb-4 shadow-lg animate-bounce">
              <CheckCircle className="w-10 h-10 text-white" />
            </div>
            <h2 className="text-4xl font-bold bg-gradient-to-r from-gray-900 to-green-700 bg-clip-text text-transparent mb-2">Registration Successful!</h2>
            <p className="text-green-600 font-bold mb-2 text-lg">✓ User has been successfully registered</p>
            <p className="text-gray-600">Your wallet has been created. Please save your private key securely before continuing.</p>
          </div>

          <div className="bg-gradient-to-r from-red-50 to-orange-50 border-l-4 border-red-500 rounded-r-xl p-5 mb-6 shadow-md transform transition-all hover:scale-[1.02]">
            <div className="flex">
              <div className="flex-shrink-0">
                <AlertCircle className="h-6 w-6 text-red-600 animate-pulse" />
              </div>
              <div className="ml-3">
                <h3 className="text-sm font-bold text-red-900">⚠️ Important Warning</h3>
                <p className="text-sm text-red-800 mt-1 font-medium">
                  Your private key is the ONLY way to access your wallet. Save it securely and never share it with anyone!
                </p>
              </div>
            </div>
          </div>

          <div className="bg-gradient-to-br from-gray-50 to-blue-50 rounded-2xl p-5 mb-6 border-2 border-blue-200 shadow-inner">
            <label className="block text-sm font-bold text-gray-700 mb-3 flex items-center">
              <Key className="w-4 h-4 mr-2 text-blue-600" />
              Your Private Key
            </label>
            <div className="bg-white border-2 border-blue-300 rounded-xl p-4 font-mono text-sm break-all shadow-lg hover:shadow-xl transition-shadow">
              {privateKey}
            </div>
          </div>

          <div className="flex gap-3">
            <button
              onClick={handleCopyPrivateKey}
              className="flex-1 bg-gradient-to-r from-blue-600 to-indigo-600 text-white py-3.5 rounded-xl font-bold hover:from-blue-700 hover:to-indigo-700 transition-all duration-300 shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 active:translate-y-0"
            >
              📋 Copy Private Key
            </button>
            <button
              onClick={handleContinueToDashboard}
              className="flex-1 bg-gradient-to-r from-green-600 to-emerald-600 text-white py-3.5 rounded-xl font-bold hover:from-green-700 hover:to-emerald-700 transition-all duration-300 shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 active:translate-y-0"
            >
              Login to System →
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-50 via-blue-50 to-purple-50 flex items-center justify-center p-4 relative overflow-hidden">
      {/* Animated background */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-purple-300 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-blob"></div>
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-blue-300 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-blob animation-delay-2000"></div>
        <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-80 h-80 bg-indigo-300 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-blob animation-delay-4000"></div>
      </div>

      <div className="bg-white/80 backdrop-blur-lg rounded-3xl shadow-3xl p-8 max-w-md w-full relative z-10 border border-white/20">
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-20 h-20 bg-gradient-to-br from-blue-600 to-indigo-700 rounded-2xl mb-4 shadow-lg transform transition-all duration-300 hover:scale-110 hover:rotate-6">
            <Wallet className="w-10 h-10 text-white" />
          </div>
          <h2 className="text-4xl font-bold bg-gradient-to-r from-gray-900 via-blue-900 to-indigo-900 bg-clip-text text-transparent mb-2">Create Account</h2>
          <p className="text-gray-600">Join the blockchain revolution</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Email Verification Success Banner */}
          {showVerificationSuccess && (
            <div className="bg-green-50 border-l-4 border-green-500 p-4 mb-4 rounded">
              <div className="flex items-center">
                <CheckCircle className="h-5 w-5 text-green-600 mr-2" />
                <div>
                  <p className="text-sm font-medium text-green-800">
                    Email Verified Successfully!
                  </p>
                  <p className="text-xs text-green-700 mt-1">
                    You can now complete your registration below.
                  </p>
                </div>
              </div>
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Full Name
            </label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <User className="h-5 w-5 text-gray-400" />
              </div>
              <input
                type="text"
                name="fullName"
                value={formData.fullName}
                onChange={handleChange}
                required
                className="pl-10 w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="John Doe"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Email Address
            </label>
            <div className="flex gap-2">
              <div className="relative flex-1">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <Mail className="h-5 w-5 text-gray-400" />
                </div>
                <input
                  type="email"
                  name="email"
                  value={formData.email}
                  onChange={handleChange}
                  required
                  disabled={emailVerified}
                  className={`pl-10 w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                    emailVerified
                      ? 'bg-green-50 border-green-500 cursor-not-allowed'
                      : 'border-gray-300'
                  }`}
                  placeholder="john@example.com"
                />
                {emailVerified && (
                  <div className="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
                    <CheckCircle className="h-5 w-5 text-green-600" />
                  </div>
                )}
              </div>
              <button
                type="button"
                onClick={handleVerifyEmail}
                disabled={emailVerified}
                className={`px-4 py-3 rounded-lg font-medium transition whitespace-nowrap ${
                  emailVerified
                    ? 'bg-green-100 text-green-700 cursor-not-allowed'
                    : 'bg-blue-600 text-white hover:bg-blue-700'
                }`}
              >
                {emailVerified ? 'Verified' : 'Verify'}
              </button>
            </div>
            {!emailVerified && (
              <p className="text-xs text-orange-600 mt-1 flex items-center gap-1">
                <AlertCircle className="w-3 h-3" />
                Please verify your email before registration
              </p>
            )}
            {emailVerified && (
              <p className="text-xs text-green-600 mt-1 flex items-center gap-1">
                <CheckCircle className="w-3 h-3" />
                Email verified successfully
              </p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              CNIC Number
            </label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <CreditCard className="h-5 w-5 text-gray-400" />
              </div>
              <input
                type="text"
                name="cnic"
                value={formData.cnic}
                onChange={handleChange}
                required
                className="pl-10 w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="42101-1234567-8"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Password
            </label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <Lock className="h-5 w-5 text-gray-400" />
              </div>
              <input
                type={showPassword ? 'text' : 'password'}
                name="password"
                value={formData.password}
                onChange={handleChange}
                required
                className="pl-10 pr-10 w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="••••••••"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute inset-y-0 right-0 pr-3 flex items-center"
              >
                {showPassword ? (
                  <EyeOff className="h-5 w-5 text-gray-400" />
                ) : (
                  <Eye className="h-5 w-5 text-gray-400" />
                )}
              </button>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Confirm Password
            </label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <Lock className="h-5 w-5 text-gray-400" />
              </div>
              <input
                type="password"
                name="confirmPassword"
                value={formData.confirmPassword}
                onChange={handleChange}
                required
                className="pl-10 w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="••••••••"
              />
            </div>
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-blue-600 text-white py-3 rounded-lg font-semibold hover:bg-blue-700 transition disabled:bg-gray-400 disabled:cursor-not-allowed"
          >
            {loading ? 'Creating Account...' : 'Create Account'}
          </button>
        </form>

        <p className="text-center text-gray-600 mt-6">
          Already have an account?{' '}
          <Link to="/login" className="text-blue-600 hover:text-blue-700 font-semibold">
            Sign In
          </Link>
        </p>
      </div>
    </div>
  );
};

export default Register;
