import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import api from '../utils/api';
import Navbar from '../components/Navbar';
import {
  ArrowLeft,
  Copy,
  Check,
  ArrowUpRight,
  ArrowDownLeft,
  Calendar,
  Hash,
  Wallet,
  FileText,
  Box,
} from 'lucide-react';
import toast from 'react-hot-toast';

const TransactionDetails = () => {
  const { hash } = useParams();
  const navigate = useNavigate();
  const { userProfile } = useAuth();
  const [transaction, setTransaction] = useState(null);
  const [block, setBlock] = useState(null);
  const [loading, setLoading] = useState(true);
  const [copiedField, setCopiedField] = useState(null);

  useEffect(() => {
    fetchTransactionDetails();
  }, [hash]);

  const fetchTransactionDetails = async () => {
    try {
      // Fetch transaction
      const txRes = await api.get(`/transaction/${hash}`);
      const txData = txRes.data?.transaction || txRes.data;
      setTransaction(txData);

      // Fetch block if transaction has blockHash
      if (txData.blockHash) {
        try {
          const blockRes = await api.get(`/block/hash/${txData.blockHash}`);
          const blockData = blockRes.data?.block || blockRes.data;
          setBlock(blockData);
        } catch (error) {
          console.error('Error fetching block:', error);
        }
      }
    } catch (error) {
      console.error('Error fetching transaction:', error);
      toast.error('Failed to load transaction details');
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = (text, field) => {
    navigator.clipboard.writeText(text);
    setCopiedField(field);
    toast.success('Copied to clipboard');
    setTimeout(() => setCopiedField(null), 2000);
  };

  const formatDate = (timestamp) => {
    if (!timestamp) return 'N/A';
    return new Date(timestamp).toLocaleString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    });
  };

  const getTransactionType = () => {
    if (!transaction) return null;
    return transaction.senderWalletId === userProfile?.walletId ? 'sent' : 'received';
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!transaction) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-800 mb-4">Transaction Not Found</h2>
          <button
            onClick={() => navigate('/transactions')}
            className="text-blue-600 hover:text-blue-700"
          >
            Back to Transactions
          </button>
        </div>
      </div>
    );
  }

  const type = getTransactionType();

  return (
    <>
      <Navbar />
      <div className="min-h-screen bg-gray-50 py-8">
        <div className="max-w-4xl mx-auto px-4">
        {/* Header */}
        <div className="mb-8">
          <button
            onClick={() => navigate('/transactions')}
            className="flex items-center gap-2 text-gray-600 hover:text-gray-800 mb-4"
          >
            <ArrowLeft className="w-5 h-5" />
            Back to Transactions
          </button>
          <h1 className="text-3xl font-bold text-gray-800">Transaction Details</h1>
        </div>

        {/* Type Badge */}
        <div className="mb-6">
          {type === 'sent' ? (
            <div className="inline-flex items-center gap-2 bg-red-100 text-red-700 px-4 py-2 rounded-lg">
              <ArrowUpRight className="w-5 h-5" />
              <span className="font-semibold">Sent Transaction</span>
            </div>
          ) : (
            <div className="inline-flex items-center gap-2 bg-green-100 text-green-700 px-4 py-2 rounded-lg">
              <ArrowDownLeft className="w-5 h-5" />
              <span className="font-semibold">Received Transaction</span>
            </div>
          )}
        </div>

        {/* Amount Card */}
        <div className="bg-gradient-to-r from-blue-600 to-indigo-700 text-white rounded-xl p-8 mb-8 text-center">
          <p className="text-blue-100 text-sm mb-2">Amount</p>
          <h2 className="text-5xl font-bold">
            {type === 'sent' ? '-' : '+'}
            {transaction.amount?.toFixed(2)} <span className="text-2xl text-blue-100">BC</span>
          </h2>
        </div>

        {/* Transaction Information */}
        <div className="bg-white rounded-xl shadow-lg p-8 mb-6">
          <h3 className="text-xl font-semibold text-gray-800 mb-6">Transaction Information</h3>
          <div className="space-y-4">
            {/* Hash */}
            <div className="flex items-start justify-between pb-4 border-b border-gray-200">
              <div className="flex-1">
                <div className="flex items-center gap-2 text-gray-600 mb-2">
                  <Hash className="w-4 h-4" />
                  <span className="text-sm font-medium">Transaction Hash</span>
                </div>
                <p className="text-gray-800 font-mono text-sm break-all">{transaction.hash}</p>
              </div>
              <button
                onClick={() => copyToClipboard(transaction.hash, 'hash')}
                className="ml-4 p-2 hover:bg-gray-100 rounded-lg transition"
              >
                {copiedField === 'hash' ? (
                  <Check className="w-5 h-5 text-green-600" />
                ) : (
                  <Copy className="w-5 h-5 text-gray-600" />
                )}
              </button>
            </div>

            {/* Sender */}
            <div className="flex items-start justify-between pb-4 border-b border-gray-200">
              <div className="flex-1">
                <div className="flex items-center gap-2 text-gray-600 mb-2">
                  <Wallet className="w-4 h-4" />
                  <span className="text-sm font-medium">From</span>
                </div>
                <p className="text-gray-800 font-mono text-sm break-all">
                  {transaction.senderWalletId}
                </p>
              </div>
              <button
                onClick={() => copyToClipboard(transaction.senderWalletId, 'sender')}
                className="ml-4 p-2 hover:bg-gray-100 rounded-lg transition"
              >
                {copiedField === 'sender' ? (
                  <Check className="w-5 h-5 text-green-600" />
                ) : (
                  <Copy className="w-5 h-5 text-gray-600" />
                )}
              </button>
            </div>

            {/* Receiver */}
            <div className="flex items-start justify-between pb-4 border-b border-gray-200">
              <div className="flex-1">
                <div className="flex items-center gap-2 text-gray-600 mb-2">
                  <Wallet className="w-4 h-4" />
                  <span className="text-sm font-medium">To</span>
                </div>
                <p className="text-gray-800 font-mono text-sm break-all">
                  {transaction.receiverWalletId}
                </p>
              </div>
              <button
                onClick={() => copyToClipboard(transaction.receiverWalletId, 'receiver')}
                className="ml-4 p-2 hover:bg-gray-100 rounded-lg transition"
              >
                {copiedField === 'receiver' ? (
                  <Check className="w-5 h-5 text-green-600" />
                ) : (
                  <Copy className="w-5 h-5 text-gray-600" />
                )}
              </button>
            </div>

            {/* Timestamp */}
            <div className="pb-4 border-b border-gray-200">
              <div className="flex items-center gap-2 text-gray-600 mb-2">
                <Calendar className="w-4 h-4" />
                <span className="text-sm font-medium">Timestamp</span>
              </div>
              <p className="text-gray-800">{formatDate(transaction.timestamp)}</p>
            </div>

            {/* Note */}
            {transaction.note && (
              <div className="pb-4 border-b border-gray-200">
                <div className="flex items-center gap-2 text-gray-600 mb-2">
                  <FileText className="w-4 h-4" />
                  <span className="text-sm font-medium">Note</span>
                </div>
                <p className="text-gray-800">{transaction.note}</p>
              </div>
            )}

            {/* Signature */}
            {transaction.signature && (
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-2 text-gray-600 mb-2">
                    <FileText className="w-4 h-4" />
                    <span className="text-sm font-medium">Signature</span>
                  </div>
                  <p className="text-gray-800 font-mono text-sm break-all">
                    {transaction.signature}
                  </p>
                </div>
                <button
                  onClick={() => copyToClipboard(transaction.signature, 'signature')}
                  className="ml-4 p-2 hover:bg-gray-100 rounded-lg transition"
                >
                  {copiedField === 'signature' ? (
                    <Check className="w-5 h-5 text-green-600" />
                  ) : (
                    <Copy className="w-5 h-5 text-gray-600" />
                  )}
                </button>
              </div>
            )}
          </div>
        </div>

        {/* Block Information */}
        {block && (
          <div className="bg-white rounded-xl shadow-lg p-8">
            <h3 className="text-xl font-semibold text-gray-800 mb-6">Block Information</h3>
            <div className="space-y-4">
              {/* Block Hash */}
              <div className="flex items-start justify-between pb-4 border-b border-gray-200">
                <div className="flex-1">
                  <div className="flex items-center gap-2 text-gray-600 mb-2">
                    <Box className="w-4 h-4" />
                    <span className="text-sm font-medium">Block Hash</span>
                  </div>
                  <p className="text-gray-800 font-mono text-sm break-all">{block.hash}</p>
                </div>
                <button
                  onClick={() => copyToClipboard(block.hash, 'blockHash')}
                  className="ml-4 p-2 hover:bg-gray-100 rounded-lg transition"
                >
                  {copiedField === 'blockHash' ? (
                    <Check className="w-5 h-5 text-green-600" />
                  ) : (
                    <Copy className="w-5 h-5 text-gray-600" />
                  )}
                </button>
              </div>

              {/* Block Index */}
              <div className="pb-4 border-b border-gray-200">
                <div className="flex items-center gap-2 text-gray-600 mb-2">
                  <Hash className="w-4 h-4" />
                  <span className="text-sm font-medium">Block Index</span>
                </div>
                <p className="text-gray-800 font-semibold">#{block.index}</p>
              </div>

              {/* Block Timestamp */}
              <div>
                <div className="flex items-center gap-2 text-gray-600 mb-2">
                  <Calendar className="w-4 h-4" />
                  <span className="text-sm font-medium">Block Mined At</span>
                </div>
                <p className="text-gray-800">{formatDate(block.timestamp)}</p>
              </div>
            </div>

            {/* View Block Button */}
            <button
              onClick={() => navigate(`/blockchain/${block.hash}`)}
              className="w-full mt-6 bg-blue-600 text-white py-3 rounded-lg font-semibold hover:bg-blue-700 transition"
            >
              View Block Details
            </button>
          </div>
        )}

        {/* Status */}
        <div className="mt-6 bg-green-50 border-l-4 border-green-500 p-4">
          <div className="flex items-center">
            <Check className="w-6 h-6 text-green-600 mr-3" />
            <div>
              <h4 className="text-green-800 font-semibold">Transaction Confirmed</h4>
              <p className="text-green-700 text-sm mt-1">
                This transaction has been confirmed and added to the blockchain
              </p>
            </div>
          </div>
        </div>
      </div>
      </div>
    </>
  );
};

export default TransactionDetails;
