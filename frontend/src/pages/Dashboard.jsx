import { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';
import api from '../utils/api';
import Navbar from '../components/Navbar';
import { 
  Wallet, 
  Send, 
  ArrowDownToLine, 
  TrendingUp, 
  Clock, 
  Eye,
  Copy,
  RefreshCw
} from 'lucide-react';
import toast from 'react-hot-toast';

const Dashboard = () => {
  const { userProfile, refreshProfile } = useAuth();
  const navigate = useNavigate();
  const [balance, setBalance] = useState(0);
  const [transactions, setTransactions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  useEffect(() => {
    fetchDashboardData();
  }, []);

  const fetchDashboardData = async () => {
    try {
      setLoading(true);
      
      // Fetch balance
      const balanceRes = await api.get('/balance');
      setBalance(balanceRes.data.balance || 0);

      // Fetch recent transactions
      const txRes = await api.get('/transactions');
      setTransactions(txRes.data.transactions?.slice(0, 5) || []);

      await refreshProfile();
    } catch (error) {
      console.error('Error fetching dashboard data:', error);
      toast.error('Failed to load dashboard data');
    } finally {
      setLoading(false);
    }
  };

  const handleRefresh = async () => {
    setRefreshing(true);
    await fetchDashboardData();
    setRefreshing(false);
    toast.success('Dashboard refreshed!');
  };

  const copyWalletId = () => {
    if (userProfile?.walletId) {
      navigator.clipboard.writeText(userProfile.walletId);
      toast.success('Wallet ID copied!');
    }
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString();
  };

  const truncateHash = (hash) => {
    if (!hash) return '';
    return `${hash.substring(0, 10)}...${hash.substring(hash.length - 10)}`;
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <RefreshCw className="w-12 h-12 text-blue-600 animate-spin mx-auto mb-4" />
          <p className="text-gray-600">Loading dashboard...</p>
        </div>
      </div>
    );
  }

  return (
    <>
      <Navbar />
      <div className="min-h-screen bg-gradient-to-br from-gray-50 via-blue-50 to-indigo-50">
        {/* Header with Gradient */}
        <div className="bg-gradient-to-r from-blue-600 via-indigo-600 to-purple-600 text-white shadow-2xl">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
            <div className="flex justify-between items-start animate-fade-in-up">
              <div>
                <h1 className="text-4xl font-bold mb-2 drop-shadow-lg">👋 Welcome back, {userProfile?.fullName}!</h1>
                <p className="text-blue-100 text-lg">{userProfile?.email}</p>
              </div>
              <button
                onClick={handleRefresh}
                disabled={refreshing}
                className="bg-white/20 hover:bg-white/30 backdrop-blur-sm px-5 py-2.5 rounded-xl flex items-center gap-2 transition-all duration-300 shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 border border-white/30"
              >
                <RefreshCw className={`w-5 h-5 ${refreshing ? 'animate-spin' : ''}`} />
                Refresh
              </button>
            </div>

            {/* Wallet Info with Glassmorphism */}
            <div className="mt-6 bg-white/10 backdrop-blur-md rounded-2xl p-5 border border-white/20 shadow-xl hover:bg-white/15 transition-all duration-300">
              <p className="text-blue-100 text-sm mb-2 font-semibold">Your Wallet ID</p>
              <div className="flex items-center gap-3">
                <p className="font-mono text-lg font-medium">{userProfile?.walletId}</p>
                <button
                  onClick={copyWalletId}
                  className="hover:bg-white/20 p-2.5 rounded-lg transition-all duration-200 transform hover:scale-110"
                  title="Copy Wallet ID"
                >
                  <Copy className="w-5 h-5" />
                </button>
              </div>
            </div>
          </div>
        </div>

      {/* Main Content */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Balance Card with Gradient */}
        <div className="bg-gradient-to-br from-white via-blue-50 to-indigo-50 rounded-3xl shadow-2xl p-8 mb-8 border-2 border-white/50 transform transition-all duration-500 hover:scale-[1.02] hover:shadow-3xl animate-fade-in-up">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-gray-600 text-sm mb-2 font-semibold uppercase tracking-wide">💰 Total Balance</p>
              <h2 className="text-6xl font-black bg-gradient-to-r from-blue-600 via-indigo-600 to-purple-600 bg-clip-text text-transparent">
                {balance.toFixed(2)} <span className="text-3xl">BC</span>
              </h2>
            </div>
            <div className="bg-gradient-to-br from-blue-500 to-indigo-600 p-5 rounded-2xl shadow-xl animate-float">
              <Wallet className="w-14 h-14 text-white" />
            </div>
          </div>
        </div>

        {/* Quick Actions with Enhanced Styling */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <button
            onClick={() => navigate('/send')}
            className="group bg-gradient-to-br from-white to-green-50 rounded-2xl shadow-lg p-6 hover:shadow-2xl transition-all duration-300 transform hover:-translate-y-2 border-2 border-green-100 hover:border-green-300 relative overflow-hidden"
          >
            <div className="absolute inset-0 bg-gradient-to-r from-green-400/0 via-green-400/10 to-green-400/0 translate-x-[-100%] group-hover:translate-x-[100%] transition-transform duration-1000"></div>
            <div className="flex items-center gap-4 relative z-10">
              <div className="bg-gradient-to-br from-green-500 to-emerald-600 p-4 rounded-xl shadow-lg group-hover:scale-110 transition-transform duration-300">
                <Send className="w-7 h-7 text-white" />
              </div>
              <div className="text-left">
                <h3 className="font-semibold text-gray-800">Send Money</h3>
                <p className="text-sm text-gray-500">Transfer to another wallet</p>
              </div>
            </div>
          </button>

          <button
            onClick={() => navigate('/transactions')}
            className="group bg-gradient-to-br from-white to-purple-50 rounded-2xl shadow-lg p-6 hover:shadow-2xl transition-all duration-300 transform hover:-translate-y-2 border-2 border-purple-100 hover:border-purple-300 relative overflow-hidden"
          >
            <div className="absolute inset-0 bg-gradient-to-r from-purple-400/0 via-purple-400/10 to-purple-400/0 translate-x-[-100%] group-hover:translate-x-[100%] transition-transform duration-1000"></div>
            <div className="flex items-center gap-4 relative z-10">
              <div className="bg-gradient-to-br from-purple-500 to-violet-600 p-4 rounded-xl shadow-lg group-hover:scale-110 transition-transform duration-300">
                <Clock className="w-7 h-7 text-white" />
              </div>
              <div className="text-left">
                <h3 className="font-bold text-gray-800">Transaction History</h3>
                <p className="text-sm text-gray-600">View all transactions</p>
              </div>
            </div>
          </button>

          <button
            onClick={() => navigate('/blockchain')}
            className="group bg-gradient-to-br from-white to-orange-50 rounded-2xl shadow-lg p-6 hover:shadow-2xl transition-all duration-300 transform hover:-translate-y-2 border-2 border-orange-100 hover:border-orange-300 relative overflow-hidden"
          >
            <div className="absolute inset-0 bg-gradient-to-r from-orange-400/0 via-orange-400/10 to-orange-400/0 translate-x-[-100%] group-hover:translate-x-[100%] transition-transform duration-1000"></div>
            <div className="flex items-center gap-4 relative z-10">
              <div className="bg-gradient-to-br from-orange-500 to-amber-600 p-4 rounded-xl shadow-lg group-hover:scale-110 transition-transform duration-300">
                <TrendingUp className="w-7 h-7 text-white" />
              </div>
              <div className="text-left">
                <h3 className="font-bold text-gray-800">Blockchain Explorer</h3>
                <p className="text-sm text-gray-600">Explore blocks</p>
              </div>
            </div>
          </button>
        </div>

        {/* Recent Transactions with Modern Card */}
        <div className="bg-white/80 backdrop-blur-sm rounded-3xl shadow-xl p-8 border border-gray-100 hover:shadow-2xl transition-all duration-300">
          <div className="flex items-center justify-between mb-6">
            <h3 className="text-2xl font-bold bg-gradient-to-r from-gray-800 to-blue-600 bg-clip-text text-transparent flex items-center gap-2">
              <Clock className="w-6 h-6 text-blue-600" />
              Recent Transactions
            </h3>
            <button
              onClick={() => navigate('/transactions')}
              className="text-blue-600 hover:text-blue-700 text-sm font-bold px-4 py-2 rounded-lg hover:bg-blue-50 transition-all duration-200 transform hover:scale-105"
            >
              View All →
            </button>
          </div>

          {transactions.length === 0 ? (
            <div className="text-center py-12">
              <Clock className="w-12 h-12 text-gray-300 mx-auto mb-4" />
              <p className="text-gray-500">No transactions yet</p>
              <button
                onClick={() => navigate('/send')}
                className="mt-4 bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition"
              >
                Make your first transaction
              </button>
            </div>
          ) : (
            <div className="space-y-4">
              {transactions.map((tx) => (
                <div
                  key={tx.hash}
                  className="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition"
                >
                  <div className="flex items-center gap-4">
                    <div className={`p-2 rounded-full ${
                      tx.senderWalletId === userProfile?.walletId
                        ? 'bg-red-100'
                        : 'bg-green-100'
                    }`}>
                      {tx.senderWalletId === userProfile?.walletId ? (
                        <ArrowDownToLine className="w-5 h-5 text-red-600" />
                      ) : (
                        <ArrowDownToLine className="w-5 h-5 text-green-600 rotate-180" />
                      )}
                    </div>
                    <div>
                      <p className="font-semibold text-gray-800">
                        {tx.senderWalletId === userProfile?.walletId ? 'Sent' : 'Received'}
                      </p>
                      <p className="text-sm text-gray-500">{formatDate(tx.timestamp)}</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className={`font-bold ${
                      tx.senderWalletId === userProfile?.walletId
                        ? 'text-red-600'
                        : 'text-green-600'
                    }`}>
                      {tx.senderWalletId === userProfile?.walletId ? '-' : '+'}
                      {tx.amount.toFixed(2)} BC
                    </p>
                    <button
                      onClick={() => navigate(`/transaction/${tx.hash}`)}
                      className="text-sm text-blue-600 hover:text-blue-700 flex items-center gap-1"
                    >
                      <Eye className="w-3 h-3" />
                      View Details
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
      </div>
    </>
  );
};

export default Dashboard;
