import { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';
import api from '../utils/api';
import Navbar from '../components/Navbar';
import { Send, Wallet, ArrowLeft, Check } from 'lucide-react';
import toast from 'react-hot-toast';

const SendMoney = () => {
  const { userProfile } = useAuth();
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    receiverWalletId: '',
    amount: '',
    note: '',
  });
  const [balance, setBalance] = useState(0);
  const [loading, setLoading] = useState(false);
  const [validating, setValidating] = useState(false);
  const [isValidWallet, setIsValidWallet] = useState(null);

  useEffect(() => {
    fetchBalance();
  }, []);

  const fetchBalance = async () => {
    try {
      const res = await api.get('/balance');
      setBalance(res.data.balance || 0);
    } catch (error) {
      console.error('Error fetching balance:', error);
    }
  };

  const validateWallet = async (walletId) => {
    if (!walletId || walletId.length < 10) {
      setIsValidWallet(null);
      return;
    }

    if (walletId === userProfile?.walletId) {
      setIsValidWallet(false);
      toast.error('Cannot send to your own wallet');
      return;
    }

    setValidating(true);
    try {
      const res = await api.get(`/wallet/validate/${walletId}`);
      setIsValidWallet(res.data.valid);
      if (!res.data.valid) {
        toast.error('Invalid wallet ID');
      }
    } catch (error) {
      setIsValidWallet(false);
    } finally {
      setValidating(false);
    }
  };

  const handleWalletIdChange = (e) => {
    const value = e.target.value;
    setFormData({ ...formData, receiverWalletId: value });
    validateWallet(value);
  };

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!isValidWallet) {
      toast.error('Please enter a valid wallet ID');
      return;
    }

    const amount = parseFloat(formData.amount);
    if (amount <= 0) {
      toast.error('Amount must be greater than 0');
      return;
    }

    if (amount < 0.01) {
      toast.error('Minimum transaction amount is 0.01 BC');
      return;
    }

    if (amount > 1000000) {
      toast.error('Maximum transaction amount is 1,000,000 BC');
      return;
    }

    if (amount > balance) {
      toast.error('Insufficient balance');
      return;
    }

    // Get private key from localStorage
    const privateKey = localStorage.getItem('privateKey');
    if (!privateKey) {
      toast.error('Private key not found. Please login again.');
      return;
    }

    setLoading(true);
    try {
      const res = await api.post('/transaction', {
        receiverWalletId: formData.receiverWalletId,
        amount: amount,
        note: formData.note || '',
        privateKey: privateKey,
      });

      toast.success('Transaction created successfully!');
      // Reset form
      setFormData({ receiverWalletId: '', amount: '', note: '' });
      setIsValidWallet(null);
      // Refresh balance
      await fetchBalance();
      // Navigate to transactions after a short delay
      setTimeout(() => navigate('/transactions'), 1500);
    } catch (error) {
      const message = error.response?.data?.error || 'Failed to create transaction';
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <Navbar />
      <div className="min-h-screen bg-gray-50 py-8">
        <div className="max-w-2xl mx-auto px-4">
        {/* Header */}
        <div className="mb-8">
          <button
            onClick={() => navigate('/dashboard')}
            className="flex items-center gap-2 text-gray-600 hover:text-gray-800 mb-4"
          >
            <ArrowLeft className="w-5 h-5" />
            Back to Dashboard
          </button>
          <h1 className="text-3xl font-bold text-gray-800">Send Money</h1>
          <p className="text-gray-600 mt-2">Transfer cryptocurrency to another wallet</p>
        </div>

        {/* Balance Card */}
        <div className="bg-gradient-to-r from-blue-600 to-indigo-700 text-white rounded-xl p-6 mb-8">
          <p className="text-blue-100 text-sm mb-1">Available Balance</p>
          <h2 className="text-4xl font-bold">
            {balance.toFixed(2)} <span className="text-xl text-blue-100">BC</span>
          </h2>
        </div>

        {/* Send Form */}
        <div className="bg-white rounded-xl shadow-lg p-8">
          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Receiver Wallet ID */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Recipient Wallet ID
              </label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <Wallet className="h-5 w-5 text-gray-400" />
                </div>
                <input
                  type="text"
                  name="receiverWalletId"
                  value={formData.receiverWalletId}
                  onChange={handleWalletIdChange}
                  required
                  className={`pl-10 pr-10 w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                    isValidWallet === true
                      ? 'border-green-500'
                      : isValidWallet === false
                      ? 'border-red-500'
                      : 'border-gray-300'
                  }`}
                  placeholder="Enter recipient's wallet ID"
                />
                {validating && (
                  <div className="absolute inset-y-0 right-0 pr-3 flex items-center">
                    <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-blue-600"></div>
                  </div>
                )}
                {isValidWallet === true && !validating && (
                  <div className="absolute inset-y-0 right-0 pr-3 flex items-center">
                    <Check className="h-5 w-5 text-green-500" />
                  </div>
                )}
              </div>
              {isValidWallet === false && (
                <p className="text-red-500 text-sm mt-1">Invalid or non-existent wallet ID</p>
              )}
            </div>

            {/* Amount */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Amount (BC)
              </label>
              <input
                type="number"
                name="amount"
                value={formData.amount}
                onChange={handleChange}
                required
                step="0.01"
                min="0.01"
                max={balance}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="0.00"
              />
              <p className="text-sm text-gray-500 mt-1">
                Maximum: {balance.toFixed(2)} BC
              </p>
            </div>

            {/* Note */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Note (Optional)
              </label>
              <textarea
                name="note"
                value={formData.note}
                onChange={handleChange}
                rows={3}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Add a note for this transaction..."
              />
            </div>

            {/* Transaction Summary */}
            {formData.amount && parseFloat(formData.amount) > 0 && (
              <div className="bg-gray-50 rounded-lg p-4">
                <h3 className="font-semibold text-gray-800 mb-3">Transaction Summary</h3>
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-gray-600">Amount to send:</span>
                    <span className="font-semibold">{parseFloat(formData.amount).toFixed(2)} BC</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Remaining balance:</span>
                    <span className="font-semibold">
                      {(balance - parseFloat(formData.amount)).toFixed(2)} BC
                    </span>
                  </div>
                </div>
              </div>
            )}

            {/* Submit Button */}
            <button
              type="submit"
              disabled={loading || !isValidWallet || !formData.amount}
              className="w-full bg-blue-600 text-white py-3 rounded-lg font-semibold hover:bg-blue-700 transition disabled:bg-gray-400 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              {loading ? (
                <>
                  <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
                  Processing...
                </>
              ) : (
                <>
                  <Send className="w-5 h-5" />
                  Send Money
                </>
              )}
            </button>
          </form>
        </div>

        {/* Info Card */}
        <div className="bg-blue-50 border-l-4 border-blue-500 p-4 mt-6">
          <div className="flex">
            <div className="ml-3">
              <h3 className="text-sm font-medium text-blue-800">Transaction Information</h3>
              <div className="mt-2 text-sm text-blue-700">
                <ul className="list-disc list-inside space-y-1">
                  <li>Transactions are irreversible once confirmed</li>
                  <li>Your transaction will be added to the pending pool</li>
                  <li>It will be confirmed when the next block is mined</li>
                  <li>Double-check the recipient wallet ID before sending</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
      </div>
    </>
  );
};

export default SendMoney;
