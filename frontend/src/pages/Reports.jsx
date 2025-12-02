import { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';
import api from '../utils/api';
import Navbar from '../components/Navbar';
import { DollarSign, Calendar, FileText, TrendingUp, Download, Filter } from 'lucide-react';
import toast from 'react-hot-toast';

const Reports = () => {
  const { userProfile } = useAuth();
  const navigate = useNavigate();
  const [zakatHistory, setZakatHistory] = useState([]);
  const [transactionLogs, setTransactionLogs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('zakat'); // zakat, transactions
  const [dateRange, setDateRange] = useState('all'); // all, month, year

  useEffect(() => {
    if (!userProfile) {
      navigate('/login');
      return;
    }
    fetchReports();
  }, [userProfile, navigate]);

  const fetchReports = async () => {
    try {
      setLoading(true);
      // Fetch Zakat History
      const zakatRes = await api.get('/zakat/history');
      console.log('Zakat API Response:', zakatRes.data);
      const zakatData = zakatRes.data?.history || zakatRes.data || [];
      console.log('Extracted zakat history:', zakatData);
      setZakatHistory(Array.isArray(zakatData) ? zakatData : []);

      // Fetch Transaction Logs (user-specific endpoint)
      const logsRes = await api.get('/logs/transactions');
      console.log('Transaction Logs API Response:', logsRes.data);
      const logsData = logsRes.data?.logs || logsRes.data || [];
      console.log('Extracted transaction logs:', logsData);
      setTransactionLogs(Array.isArray(logsData) ? logsData : []);
    } catch (error) {
      console.error('Error fetching reports:', error);
      toast.error('Failed to load reports');
      setZakatHistory([]);
      setTransactionLogs([]);
    } finally {
      setLoading(false);
    }
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

  const formatCurrency = (amount) => {
    return amount?.toFixed(2) || '0.00';
  };

  const filterByDateRange = (items) => {
    if (dateRange === 'all') return items;

    const now = new Date();
    const filtered = items.filter((item) => {
      const itemDate = new Date(item.timestamp || item.deductedAt);
      if (dateRange === 'month') {
        return (
          itemDate.getMonth() === now.getMonth() &&
          itemDate.getFullYear() === now.getFullYear()
        );
      } else if (dateRange === 'year') {
        return itemDate.getFullYear() === now.getFullYear();
      }
      return true;
    });

    return filtered;
  };

  const getTotalZakat = () => {
    const filtered = filterByDateRange(zakatHistory);
    return filtered.reduce((sum, z) => sum + (z.amount || 0), 0);
  };

  const handleDownloadReport = () => {
    toast.success('Report download feature coming soon!');
    // In a real implementation, this would generate a PDF or CSV
  };

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
            <p className="mt-6 text-gray-700 font-semibold text-lg">Loading reports...</p>
          </div>
        </div>
      </>
    );
  }

  return (
    <>
      <Navbar />
      <div className="min-h-screen bg-gradient-to-br from-gray-50 via-green-50 to-blue-50 py-8">
        <div className="max-w-7xl mx-auto px-4">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-800">Reports & Analytics</h1>
          <p className="text-gray-600 mt-2">View your Zakat deductions and transaction logs</p>
        </div>

        {/* Summary Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <div className="bg-gradient-to-r from-green-500 to-green-600 text-white rounded-xl shadow-md p-6">
            <div className="flex items-center gap-3 mb-2">
              <DollarSign className="w-8 h-8" />
              <p className="text-green-100 text-sm">Total Zakat Paid</p>
            </div>
            <p className="text-4xl font-bold">{formatCurrency(getTotalZakat())} BC</p>
          </div>

          <div className="bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-xl shadow-md p-6">
            <div className="flex items-center gap-3 mb-2">
              <FileText className="w-8 h-8" />
              <p className="text-blue-100 text-sm">Zakat Transactions</p>
            </div>
            <p className="text-4xl font-bold">{filterByDateRange(zakatHistory).length}</p>
          </div>

          <div className="bg-gradient-to-r from-purple-500 to-purple-600 text-white rounded-xl shadow-md p-6">
            <div className="flex items-center gap-3 mb-2">
              <TrendingUp className="w-8 h-8" />
              <p className="text-purple-100 text-sm">Transaction Logs</p>
            </div>
            <p className="text-4xl font-bold">{filterByDateRange(transactionLogs).length}</p>
          </div>
        </div>

        {/* Filters */}
        <div className="bg-white rounded-xl shadow-md p-4 mb-6">
          <div className="flex flex-col md:flex-row items-center justify-between gap-4">
            {/* Date Range Filter */}
            <div className="flex items-center gap-2">
              <Filter className="w-5 h-5 text-gray-600" />
              <select
                value={dateRange}
                onChange={(e) => setDateRange(e.target.value)}
                className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="all">All Time</option>
                <option value="month">This Month</option>
                <option value="year">This Year</option>
              </select>
            </div>

            {/* Download Button */}
            <button
              onClick={handleDownloadReport}
              className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
            >
              <Download className="w-4 h-4" />
              Download Report
            </button>
          </div>
        </div>

        {/* Tabs */}
        <div className="bg-white rounded-xl shadow-lg overflow-hidden">
          <div className="flex border-b border-gray-200">
            <button
              onClick={() => setActiveTab('zakat')}
              className={`flex-1 py-4 px-6 font-semibold transition ${
                activeTab === 'zakat'
                  ? 'bg-blue-50 text-blue-600 border-b-2 border-blue-600'
                  : 'text-gray-600 hover:bg-gray-50'
              }`}
            >
              Zakat History
            </button>
            <button
              onClick={() => setActiveTab('transactions')}
              className={`flex-1 py-4 px-6 font-semibold transition ${
                activeTab === 'transactions'
                  ? 'bg-blue-50 text-blue-600 border-b-2 border-blue-600'
                  : 'text-gray-600 hover:bg-gray-50'
              }`}
            >
              Transaction Logs
            </button>
          </div>

          <div className="p-6">
            {loading ? (
              <div className="flex items-center justify-center py-12">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
              </div>
            ) : (
              <>
                {/* Zakat History Tab */}
                {activeTab === 'zakat' && (
                  <div>
                    {filterByDateRange(zakatHistory).length === 0 ? (
                      <div className="text-center py-12 text-gray-500">
                        <DollarSign className="w-12 h-12 mx-auto mb-3 text-gray-400" />
                        <p>No Zakat deductions found</p>
                      </div>
                    ) : (
                      <div className="overflow-x-auto">
                        <table className="w-full">
                          <thead className="bg-gray-50 border-b border-gray-200">
                            <tr>
                              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase">
                                Date
                              </th>
                              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase">
                                Amount
                              </th>
                              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase">
                                Transaction Hash
                              </th>
                              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase">
                                Status
                              </th>
                            </tr>
                          </thead>
                          <tbody className="divide-y divide-gray-200">
                            {filterByDateRange(zakatHistory).map((zakat, index) => (
                              <tr key={index} className="hover:bg-gray-50">
                                <td className="px-6 py-4 text-sm text-gray-600">
                                  {formatDate(zakat.deductedAt)}
                                </td>
                                <td className="px-6 py-4">
                                  <span className="font-semibold text-green-600">
                                    {formatCurrency(zakat.amount)} BC
                                  </span>
                                </td>
                                <td className="px-6 py-4">
                                  <span className="text-sm text-gray-600 font-mono">
                                    {zakat.transactionHash?.substring(0, 16)}...
                                  </span>
                                </td>
                                <td className="px-6 py-4">
                                  <span className="px-3 py-1 text-xs font-semibold rounded-full bg-green-100 text-green-800">
                                    Completed
                                  </span>
                                </td>
                              </tr>
                            ))}
                          </tbody>
                        </table>
                      </div>
                    )}
                  </div>
                )}

                {/* Transaction Logs Tab */}
                {activeTab === 'transactions' && (
                  <div>
                    {filterByDateRange(transactionLogs).length === 0 ? (
                      <div className="text-center py-12 text-gray-500">
                        <FileText className="w-12 h-12 mx-auto mb-3 text-gray-400" />
                        <p>No transaction logs found</p>
                      </div>
                    ) : (
                      <div className="overflow-x-auto">
                        <table className="w-full">
                          <thead className="bg-gray-50 border-b border-gray-200">
                            <tr>
                              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase">
                                Date
                              </th>
                              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase">
                                Event Type
                              </th>
                              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase">
                                Transaction Hash
                              </th>
                              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase">
                                Details
                              </th>
                            </tr>
                          </thead>
                          <tbody className="divide-y divide-gray-200">
                            {filterByDateRange(transactionLogs).map((log, index) => (
                              <tr key={index} className="hover:bg-gray-50">
                                <td className="px-6 py-4 text-sm text-gray-600">
                                  {formatDate(log.timestamp)}
                                </td>
                                <td className="px-6 py-4">
                                  <span className="px-3 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800">
                                    {log.eventType || 'Transaction'}
                                  </span>
                                </td>
                                <td className="px-6 py-4">
                                  <span className="text-sm text-gray-600 font-mono">
                                    {log.transactionHash?.substring(0, 16)}...
                                  </span>
                                </td>
                                <td className="px-6 py-4 text-sm text-gray-600">
                                  {log.details || 'N/A'}
                                </td>
                              </tr>
                            ))}
                          </tbody>
                        </table>
                      </div>
                    )}
                  </div>
                )}
              </>
            )}
          </div>
        </div>

        {/* Info Card */}
        <div className="bg-blue-50 border-l-4 border-blue-500 p-4 mt-6">
          <div className="flex">
            <div className="ml-3">
              <h3 className="text-sm font-medium text-blue-800">About Reports</h3>
              <div className="mt-2 text-sm text-blue-700">
                <ul className="list-disc list-inside space-y-1">
                  <li>Zakat is automatically deducted from your balance</li>
                  <li>Transaction logs track all activities on your account</li>
                  <li>Use filters to view reports for specific time periods</li>
                  <li>Download reports for your personal records</li>
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

export default Reports;
