import { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import api from '../utils/api';
import { Mail, ArrowLeft, Check, RefreshCw } from 'lucide-react';
import toast from 'react-hot-toast';

const VerifyEmail = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const email = location.state?.email || '';
  
  const [otp, setOtp] = useState(['', '', '', '', '', '']);
  const [loading, setLoading] = useState(false);
  const [resending, setResending] = useState(false);
  const [displayedOtp, setDisplayedOtp] = useState('');
  const [timer, setTimer] = useState(300); // 5 minutes in seconds

  useEffect(() => {
    if (!email) {
      navigate('/register');
      return;
    }
    // OTP is already sent when user clicked "Verify" button on Register page
    // Don't send again on mount to avoid duplicate emails
  }, [email, navigate]);

  useEffect(() => {
    if (timer > 0) {
      const interval = setInterval(() => {
        setTimer((prev) => prev - 1);
      }, 1000);
      return () => clearInterval(interval);
    }
  }, [timer]);

  const sendOTP = async () => {
    setResending(true);
    try {
      const res = await api.post('/otp/generate', { email });
      toast.success('\u2713 New OTP sent to your email');      // In development, display the OTP (remove in production)
      if (res.data.otp) {
        setDisplayedOtp(res.data.otp);
        toast.success(`Development OTP: ${res.data.otp}`, { duration: 10000 });
      }
      setTimer(300); // Reset timer
    } catch (error) {
      const message = error.response?.data?.error || 'Failed to send OTP';
      toast.error(message);
    } finally {
      setResending(false);
    }
  };

  const handleOtpChange = (index, value) => {
    if (value.length > 1) {
      value = value[0];
    }

    if (!/^\d*$/.test(value)) {
      return;
    }

    const newOtp = [...otp];
    newOtp[index] = value;
    setOtp(newOtp);

    // Auto-focus next input
    if (value && index < 5) {
      document.getElementById(`otp-${index + 1}`)?.focus();
    }
  };

  const handleKeyDown = (index, e) => {
    if (e.key === 'Backspace' && !otp[index] && index > 0) {
      document.getElementById(`otp-${index - 1}`)?.focus();
    }
  };

  const handlePaste = (e) => {
    e.preventDefault();
    const pastedData = e.clipboardData.getData('text').slice(0, 6);
    if (!/^\d+$/.test(pastedData)) {
      return;
    }

    const newOtp = [...otp];
    for (let i = 0; i < pastedData.length && i < 6; i++) {
      newOtp[i] = pastedData[i];
    }
    setOtp(newOtp);

    // Focus last filled input or first empty
    const lastIndex = Math.min(pastedData.length, 5);
    document.getElementById(`otp-${lastIndex}`)?.focus();
  };

  const handleVerify = async () => {
    const otpCode = otp.join('');
    if (otpCode.length !== 6) {
      toast.error('Please enter complete OTP');
      return;
    }

    setLoading(true);
    try {
      await api.post('/otp/verify', { email, otp: otpCode });
      toast.success('✓ Email verified successfully! Complete your registration now.', { duration: 3000 });
      setTimeout(() => {
        navigate('/register', { state: { email, verified: true } });
      }, 1500);
    } catch (error) {
      const message = error.response?.data?.error || 'Invalid OTP';
      toast.error(message);
      setOtp(['', '', '', '', '', '']);
      document.getElementById('otp-0')?.focus();
    } finally {
      setLoading(false);
    }
  };

  const formatTime = (seconds) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
      <div className="max-w-md w-full">
        {/* Back Button */}
        <button
          onClick={() => navigate('/register')}
          className="flex items-center gap-2 text-gray-600 hover:text-gray-800 mb-6"
        >
          <ArrowLeft className="w-5 h-5" />
          Back to Registration
        </button>

        {/* Verification Card */}
        <div className="bg-white rounded-2xl shadow-xl p-8">
          {/* Icon */}
          <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <Mail className="w-8 h-8 text-blue-600" />
          </div>

          {/* Title */}
          <h1 className="text-3xl font-bold text-gray-800 text-center mb-2">
            Verify Your Email
          </h1>
          <p className="text-gray-600 text-center mb-8">
            We've sent a 6-digit code to<br />
            <span className="font-semibold text-gray-800">{email}</span>
          </p>

          {/* Development OTP Display */}
          {displayedOtp && (
            <div className="bg-yellow-50 border-l-4 border-yellow-500 p-4 mb-6">
              <p className="text-sm text-yellow-800">
                <strong>Development Mode:</strong> OTP is {displayedOtp}
              </p>
            </div>
          )}

          {/* OTP Input */}
          <div className="flex gap-2 justify-center mb-6">
            {otp.map((digit, index) => (
              <input
                key={index}
                id={`otp-${index}`}
                type="text"
                inputMode="numeric"
                maxLength={1}
                value={digit}
                onChange={(e) => handleOtpChange(index, e.target.value)}
                onKeyDown={(e) => handleKeyDown(index, e)}
                onPaste={handlePaste}
                className="w-12 h-14 text-center text-2xl font-bold border-2 border-gray-300 rounded-lg focus:border-blue-500 focus:ring-2 focus:ring-blue-200 outline-none transition"
              />
            ))}
          </div>

          {/* Timer */}
          <div className="text-center mb-6">
            {timer > 0 ? (
              <p className="text-sm text-gray-600">
                Code expires in{' '}
                <span className="font-semibold text-blue-600">{formatTime(timer)}</span>
              </p>
            ) : (
              <p className="text-sm text-red-600">OTP has expired</p>
            )}
          </div>

          {/* Verify Button */}
          <button
            onClick={handleVerify}
            disabled={loading || otp.join('').length !== 6}
            className="w-full bg-blue-600 text-white py-3 rounded-lg font-semibold hover:bg-blue-700 transition disabled:bg-gray-400 disabled:cursor-not-allowed flex items-center justify-center gap-2 mb-4"
          >
            {loading ? (
              <>
                <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
                Verifying...
              </>
            ) : (
              <>
                <Check className="w-5 h-5" />
                Verify Email
              </>
            )}
          </button>

          {/* Resend Button */}
          <button
            onClick={sendOTP}
            disabled={resending || timer > 240} // Can resend after 1 minute
            className="w-full text-blue-600 py-2 rounded-lg font-medium hover:bg-blue-50 transition disabled:text-gray-400 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            {resending ? (
              <>
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
                Sending...
              </>
            ) : (
              <>
                <RefreshCw className="w-4 h-4" />
                {timer > 240 ? `Resend in ${formatTime(300 - timer)}` : 'Resend OTP'}
              </>
            )}
          </button>

          {/* Help Text */}
          <div className="mt-6 pt-6 border-t border-gray-200">
            <p className="text-sm text-gray-600 text-center">
              Didn't receive the code? Check your spam folder or click resend.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default VerifyEmail;
