import { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';
import api from '../utils/api';
import Navbar from '../components/Navbar';
import { ArrowUpRight, ArrowDownLeft, Search, Filter, Calendar, ChevronLeft, ChevronRight } from 'lucide-react';
import toast from 'react-hot-toast';

const Transactions = () => {
  const { userProfile } = useAuth();
  const navigate = useNavigate();
  const [transactions, setTransactions] = useState([]);
  const [filteredTransactions, setFilteredTransactions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState('all'); // all, sent, received
  const [searchTerm, setSearchTerm] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 10;

  useEffect(() => {
    if (!userProfile) {
      navigate('/login');
      return;
    }
    fetchTransactions();
  }, [userProfile, navigate]);

  useEffect(() => {
    applyFilters();
  }, [transactions, filter, searchTerm]);

  const fetchTransactions = async () => {
    try {
      setLoading(true);
      const res = await api.get('/transactions');
      const txData = res.data?.transactions || res.data || [];
      setTransactions(Array.isArray(txData) ? txData : []);
    } catch (error) {
      console.error('Error fetching transactions:', error);
      toast.error('Failed to load transactions');
      setTransactions([]);
    } finally {
      setLoading(false);
    }
  };

  const applyFilters = () => {
    if (!Array.isArray(transactions)) {
      setFilteredTransactions([]);
      return;
    }
    
    let filtered = [...transactions];

    // Filter by type
    if (filter === 'sent') {
      filtered = filtered.filter(
        (tx) => tx.senderWalletId === userProfile?.walletId
      );
    } else if (filter === 'received') {
      filtered = filtered.filter(
        (tx) => tx.senderWalletId !== userProfile?.walletId
      );
    }

    // Search filter
    if (searchTerm) {
      filtered = filtered.filter(
        (tx) =>
          tx.hash?.toLowerCase().includes(searchTerm.toLowerCase()) ||
          tx.receiverWalletId?.toLowerCase().includes(searchTerm.toLowerCase()) ||
          tx.senderWalletId?.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    setFilteredTransactions(filtered);
    setCurrentPage(1);
  };

  const formatDate = (timestamp) => {
    if (!timestamp) return 'N/A';
    return new Date(timestamp).toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const getTransactionType = (tx) => {
    return tx.senderWalletId === userProfile?.walletId ? 'sent' : 'received';
  };

  // Pagination
  const indexOfLastItem = currentPage * itemsPerPage;
  const indexOfFirstItem = indexOfLastItem - itemsPerPage;
  const currentTransactions = filteredTransactions.slice(indexOfFirstItem, indexOfLastItem);
  const totalPages = Math.ceil(filteredTransactions.length / itemsPerPage);

  const paginate = (pageNumber) => setCurrentPage(pageNumber);

  if (loading) {
    return (
      <>
        <Navbar />
        <div className="min-h-screen bg-gradient-to-br from-gray-50 via-blue-50 to-indigo-50 flex items-center justify-center">
          <div className="text-center">
            <div className="relative">
              <div className="animate-spin rounded-full h-16 w-16 border-b-4 border-blue-600 mx-auto"></div>
              <div className="absolute inset-0 rounded-full h-16 w-16 border-b-4 border-indigo-400 mx-auto animate-ping opacity-20"></div>
            </div>
            <p className="mt-6 text-gray-700 font-semibold text-lg">Loading transactions...</p>
          </div>
        </div>
      </>
    );
  }

  return (
    <>
      <Navbar />
      <div className="min-h-screen bg-gradient-to-br from-gray-50 via-blue-50 to-indigo-50 py-8">
        <div className="max-w-7xl mx-auto px-4">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-800">Transaction History</h1>
          <p className="text-gray-600 mt-2">View all your blockchain transactions</p>
        </div>

        {/* Stats Cards with Gradient */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <div className="bg-gradient-to-br from-white to-blue-50 rounded-2xl shadow-lg p-6 border-2 border-blue-100 transform transition-all duration-300 hover:scale-105 hover:shadow-2xl">
            <p className="text-gray-600 text-sm mb-1 font-semibold uppercase tracking-wide">Total Transactions</p>
            <p className="text-4xl font-black bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">{transactions.length}</p>
          </div>
          <div className="bg-gradient-to-br from-white to-red-50 rounded-2xl shadow-lg p-6 border-2 border-red-100 transform transition-all duration-300 hover:scale-105 hover:shadow-2xl">
            <p className="text-gray-600 text-sm mb-1 font-semibold uppercase tracking-wide">Sent</p>
            <p className="text-4xl font-black bg-gradient-to-r from-red-600 to-pink-600 bg-clip-text text-transparent">
              {transactions.filter((tx) => tx.senderWalletId === userProfile?.walletId).length}
            </p>
          </div>
          <div className="bg-gradient-to-br from-white to-green-50 rounded-2xl shadow-lg p-6 border-2 border-green-100 transform transition-all duration-300 hover:scale-105 hover:shadow-2xl">
            <p className="text-gray-600 text-sm mb-1 font-semibold uppercase tracking-wide">Received</p>
            <p className="text-4xl font-black bg-gradient-to-r from-green-600 to-emerald-600 bg-clip-text text-transparent">
              {transactions.filter((tx) => tx.senderWalletId !== userProfile?.walletId).length}
            </p>
          </div>
        </div>

        {/* Filters with Enhanced Design */}
        <div className="bg-white/80 backdrop-blur-sm rounded-2xl shadow-xl p-6 mb-6 border border-gray-100">
          <div className="flex flex-col md:flex-row gap-4">
            {/* Search with Icon Animation */}
            <div className="flex-1 relative group">
              <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                <Search className="h-5 w-5 text-gray-400 group-focus-within:text-blue-600 transition-colors" />
              </div>
              <input
                type="text"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                placeholder="Search by hash or wallet ID..."
                className="pl-11 w-full px-4 py-3 border-2 border-gray-200 rounded-xl focus:ring-4 focus:ring-blue-500/20 focus:border-blue-500 outline-none transition-all duration-300 bg-gray-50/50 hover:bg-white font-medium"
              />
            </div>

            {/* Type Filter Buttons with Gradient */}
            <div className="flex gap-2">
              <button
                onClick={() => setFilter('all')}
                className={`px-5 py-3 rounded-xl font-bold transition-all duration-300 transform hover:scale-105 ${
                  filter === 'all'
                    ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-lg'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                All
              </button>
              <button
                onClick={() => setFilter('sent')}
                className={`px-5 py-3 rounded-xl font-bold transition-all duration-300 transform hover:scale-105 ${
                  filter === 'sent'
                    ? 'bg-gradient-to-r from-red-600 to-pink-600 text-white shadow-lg'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                Sent
              </button>
              <button
                onClick={() => setFilter('received')}
                className={`px-5 py-3 rounded-xl font-bold transition-all duration-300 transform hover:scale-105 ${
                  filter === 'received'
                    ? 'bg-gradient-to-r from-green-600 to-emerald-600 text-white shadow-lg'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                Received
              </button>
            </div>
          </div>
        </div>

        {/* Transactions Table */}
        <div className="bg-white rounded-xl shadow-md overflow-hidden">
          {loading ? (
            <div className="flex items-center justify-center py-12">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
            </div>
          ) : currentTransactions.length === 0 ? (
            <div className="text-center py-12">
              <p className="text-gray-500 text-lg">No transactions found</p>
            </div>
          ) : (
            <>
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead className="bg-gray-50 border-b border-gray-200">
                    <tr>
                      <th className="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                        Type
                      </th>
                      <th className="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                        Hash
                      </th>
                      <th className="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                        From/To
                      </th>
                      <th className="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                        Amount
                      </th>
                      <th className="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                        Date
                      </th>
                      <th className="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                        Status
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {currentTransactions.map((tx) => {
                      const type = getTransactionType(tx);
                      return (
                        <tr
                          key={tx.hash}
                          onClick={() => navigate(`/transaction/${tx.hash}`)}
                          className="hover:bg-gray-50 cursor-pointer transition"
                        >
                          <td className="px-6 py-4 whitespace-nowrap">
                            {type === 'sent' ? (
                              <div className="flex items-center gap-2 text-red-600">
                                <ArrowUpRight className="w-5 h-5" />
                                <span className="font-medium">Sent</span>
                              </div>
                            ) : (
                              <div className="flex items-center gap-2 text-green-600">
                                <ArrowDownLeft className="w-5 h-5" />
                                <span className="font-medium">Received</span>
                              </div>
                            )}
                          </td>
                          <td className="px-6 py-4">
                            <span className="text-sm text-gray-600 font-mono">
                              {tx.hash?.substring(0, 16)}...
                            </span>
                          </td>
                          <td className="px-6 py-4">
                            <span className="text-sm text-gray-600 font-mono">
                              {type === 'sent'
                                ? tx.receiverWalletId?.substring(0, 12)
                                : tx.senderWalletId?.substring(0, 12)}
                              ...
                            </span>
                          </td>
                          <td className="px-6 py-4">
                            <span
                              className={`font-semibold ${
                                type === 'sent' ? 'text-red-600' : 'text-green-600'
                              }`}
                            >
                              {type === 'sent' ? '-' : '+'}
                              {tx.amount?.toFixed(2)} BC
                            </span>
                          </td>
                          <td className="px-6 py-4 text-sm text-gray-600">
                            {formatDate(tx.timestamp)}
                          </td>
                          <td className="px-6 py-4">
                            <span className="px-3 py-1 text-xs font-semibold rounded-full bg-green-100 text-green-800">
                              Confirmed
                            </span>
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>

              {/* Pagination */}
              {totalPages > 1 && (
                <div className="bg-gray-50 px-6 py-4 flex items-center justify-between border-t border-gray-200">
                  <div className="text-sm text-gray-600">
                    Showing {indexOfFirstItem + 1} to{' '}
                    {Math.min(indexOfLastItem, filteredTransactions.length)} of{' '}
                    {filteredTransactions.length} transactions
                  </div>
                  <div className="flex gap-2">
                    <button
                      onClick={() => paginate(currentPage - 1)}
                      disabled={currentPage === 1}
                      className="px-3 py-1 rounded-lg border border-gray-300 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      <ChevronLeft className="w-5 h-5" />
                    </button>
                    {[...Array(totalPages)].map((_, index) => (
                      <button
                        key={index + 1}
                        onClick={() => paginate(index + 1)}
                        className={`px-3 py-1 rounded-lg ${
                          currentPage === index + 1
                            ? 'bg-blue-600 text-white'
                            : 'border border-gray-300 hover:bg-gray-100'
                        }`}
                      >
                        {index + 1}
                      </button>
                    ))}
                    <button
                      onClick={() => paginate(currentPage + 1)}
                      disabled={currentPage === totalPages}
                      className="px-3 py-1 rounded-lg border border-gray-300 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      <ChevronRight className="w-5 h-5" />
                    </button>
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      </div>
      </div>
    </>
  );
};

export default Transactions;
