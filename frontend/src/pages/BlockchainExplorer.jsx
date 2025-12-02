import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import api from '../utils/api';
import Navbar from '../components/Navbar';
import { Box, Hash, Calendar, ArrowRight, Search, ChevronLeft, ChevronRight, Copy, Check } from 'lucide-react';
import toast from 'react-hot-toast';

const BlockchainExplorer = () => {
  const navigate = useNavigate();
  const { hash: urlHash } = useParams();
  const [blocks, setBlocks] = useState([]);
  const [selectedBlock, setSelectedBlock] = useState(null);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [copiedField, setCopiedField] = useState(null);
  const blocksPerPage = 10;

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      navigate('/login');
      return;
    }
    fetchBlocks();
  }, []);

  useEffect(() => {
    if (urlHash) {
      fetchBlockByHash(urlHash);
    }
  }, [urlHash]);

  const fetchBlocks = async () => {
    try {
      setLoading(true);
      const res = await api.get('/blockchain');
      console.log('Blockchain API Response:', res.data);
      const blocksData = res.data?.blockchain || res.data?.blocks || res.data || [];
      console.log('Extracted blocks:', blocksData);
      setBlocks(Array.isArray(blocksData) ? blocksData : []);
    } catch (error) {
      console.error('Error fetching blocks:', error);
      toast.error('Failed to load blockchain');
      setBlocks([]);
    } finally {
      setLoading(false);
    }
  };

  const fetchBlockByHash = async (hash) => {
    try {
      const res = await api.get(`/block/hash/${hash}`);
      const blockData = res.data?.block || res.data;
      setSelectedBlock(blockData);
    } catch (error) {
      console.error('Error fetching block:', error);
      toast.error('Block not found');
    }
  };

  const handleBlockClick = (block) => {
    setSelectedBlock(block);
    navigate(`/blockchain/${block.hash}`);
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

  const copyToClipboard = (text, field) => {
    navigator.clipboard.writeText(text);
    setCopiedField(field);
    toast.success('Copied to clipboard');
    setTimeout(() => setCopiedField(null), 2000);
  };

  // Filter blocks based on search
  const filteredBlocks = blocks.filter(
    (block) =>
      block.hash?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      block.index?.toString().includes(searchTerm)
  );

  // Pagination
  const indexOfLastBlock = currentPage * blocksPerPage;
  const indexOfFirstBlock = indexOfLastBlock - blocksPerPage;
  const currentBlocks = filteredBlocks.slice(indexOfFirstBlock, indexOfLastBlock);
  const totalPages = Math.ceil(filteredBlocks.length / blocksPerPage);

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
            <p className="mt-6 text-gray-700 font-semibold text-lg">Loading blockchain...</p>
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
          <h1 className="text-3xl font-bold text-gray-800">Blockchain Explorer</h1>
          <p className="text-gray-600 mt-2">Explore all blocks in the blockchain</p>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <div className="bg-white rounded-xl shadow-md p-6">
            <p className="text-gray-600 text-sm mb-1">Total Blocks</p>
            <p className="text-3xl font-bold text-gray-800">{blocks.length}</p>
          </div>
          <div className="bg-white rounded-xl shadow-md p-6">
            <p className="text-gray-600 text-sm mb-1">Total Transactions</p>
            <p className="text-3xl font-bold text-gray-800">
              {blocks.reduce((sum, block) => sum + (block.transactions?.length || 0), 0)}
            </p>
          </div>
          <div className="bg-white rounded-xl shadow-md p-6">
            <p className="text-gray-600 text-sm mb-1">Latest Block</p>
            <p className="text-3xl font-bold text-gray-800">
              #{blocks.length > 0 ? blocks[blocks.length - 1]?.index : 0}
            </p>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Blocks List */}
          <div className="bg-white rounded-xl shadow-lg p-6">
            <h2 className="text-xl font-semibold text-gray-800 mb-4">Blocks</h2>

            {/* Search */}
            <div className="mb-4 relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <Search className="h-5 w-5 text-gray-400" />
              </div>
              <input
                type="text"
                value={searchTerm}
                onChange={(e) => {
                  setSearchTerm(e.target.value);
                  setCurrentPage(1);
                }}
                placeholder="Search by hash or index..."
                className="pl-10 w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>

            {loading ? (
              <div className="flex items-center justify-center py-12">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
              </div>
            ) : currentBlocks.length === 0 ? (
              <div className="text-center py-12">
                <p className="text-gray-500">No blocks found</p>
              </div>
            ) : (
              <>
                <div className="space-y-3 mb-4">
                  {currentBlocks.map((block) => (
                    <div
                      key={block.hash}
                      onClick={() => handleBlockClick(block)}
                      className={`p-4 border rounded-lg cursor-pointer transition ${
                        selectedBlock?.hash === block.hash
                          ? 'border-blue-500 bg-blue-50'
                          : 'border-gray-200 hover:border-blue-300 hover:bg-gray-50'
                      }`}
                    >
                      <div className="flex items-center justify-between mb-2">
                        <div className="flex items-center gap-2">
                          <Box className="w-5 h-5 text-blue-600" />
                          <span className="font-semibold text-gray-800">Block #{block.index}</span>
                        </div>
                        <span className="text-sm text-gray-600">
                          {block.transactions?.length || 0} txs
                        </span>
                      </div>
                      <p className="text-xs text-gray-500 font-mono truncate">
                        {block.hash?.substring(0, 40)}...
                      </p>
                      <p className="text-xs text-gray-500 mt-1">{formatDate(block.timestamp)}</p>
                    </div>
                  ))}
                </div>

                {/* Pagination */}
                {totalPages > 1 && (
                  <div className="flex items-center justify-between pt-4 border-t border-gray-200">
                    <div className="text-sm text-gray-600">
                      Page {currentPage} of {totalPages}
                    </div>
                    <div className="flex gap-2">
                      <button
                        onClick={() => paginate(currentPage - 1)}
                        disabled={currentPage === 1}
                        className="px-3 py-1 rounded-lg border border-gray-300 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        <ChevronLeft className="w-5 h-5" />
                      </button>
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

          {/* Block Details */}
          <div className="bg-white rounded-xl shadow-lg p-6">
            <h2 className="text-xl font-semibold text-gray-800 mb-4">Block Details</h2>

            {selectedBlock ? (
              <div className="space-y-4">
                {/* Block Index */}
                <div className="pb-4 border-b border-gray-200">
                  <p className="text-sm text-gray-600 mb-1">Block Index</p>
                  <p className="text-2xl font-bold text-gray-800">#{selectedBlock.index}</p>
                </div>

                {/* Hash */}
                <div className="pb-4 border-b border-gray-200">
                  <div className="flex items-center justify-between mb-1">
                    <p className="text-sm text-gray-600">Hash</p>
                    <button
                      onClick={() => copyToClipboard(selectedBlock.hash, 'hash')}
                      className="p-1 hover:bg-gray-100 rounded"
                    >
                      {copiedField === 'hash' ? (
                        <Check className="w-4 h-4 text-green-600" />
                      ) : (
                        <Copy className="w-4 h-4 text-gray-600" />
                      )}
                    </button>
                  </div>
                  <p className="text-sm text-gray-800 font-mono break-all">{selectedBlock.hash}</p>
                </div>

                {/* Previous Hash */}
                {selectedBlock.previousHash && (
                  <div className="pb-4 border-b border-gray-200">
                    <div className="flex items-center justify-between mb-1">
                      <p className="text-sm text-gray-600">Previous Hash</p>
                      <button
                        onClick={() => copyToClipboard(selectedBlock.previousHash, 'prevHash')}
                        className="p-1 hover:bg-gray-100 rounded"
                      >
                        {copiedField === 'prevHash' ? (
                          <Check className="w-4 h-4 text-green-600" />
                        ) : (
                          <Copy className="w-4 h-4 text-gray-600" />
                        )}
                      </button>
                    </div>
                    <p className="text-sm text-gray-800 font-mono break-all">
                      {selectedBlock.previousHash}
                    </p>
                  </div>
                )}

                {/* Merkle Root */}
                {selectedBlock.merkleRoot && (
                  <div className="pb-4 border-b border-gray-200">
                    <div className="flex items-center justify-between mb-1">
                      <p className="text-sm text-gray-600">Merkle Root</p>
                      <button
                        onClick={() => copyToClipboard(selectedBlock.merkleRoot, 'merkle')}
                        className="p-1 hover:bg-gray-100 rounded"
                      >
                        {copiedField === 'merkle' ? (
                          <Check className="w-4 h-4 text-green-600" />
                        ) : (
                          <Copy className="w-4 h-4 text-gray-600" />
                        )}
                      </button>
                    </div>
                    <p className="text-sm text-gray-800 font-mono break-all">
                      {selectedBlock.merkleRoot}
                    </p>
                  </div>
                )}

                {/* Timestamp */}
                <div className="pb-4 border-b border-gray-200">
                  <p className="text-sm text-gray-600 mb-1">Timestamp</p>
                  <p className="text-sm text-gray-800">{formatDate(selectedBlock.timestamp)}</p>
                </div>

                {/* Nonce */}
                <div className="pb-4 border-b border-gray-200">
                  <p className="text-sm text-gray-600 mb-1">Nonce</p>
                  <p className="text-sm text-gray-800 font-mono">{selectedBlock.nonce}</p>
                </div>

                {/* Difficulty */}
                {selectedBlock.difficulty && (
                  <div className="pb-4 border-b border-gray-200">
                    <p className="text-sm text-gray-600 mb-1">Difficulty</p>
                    <p className="text-sm text-gray-800">{selectedBlock.difficulty}</p>
                  </div>
                )}

                {/* Miner */}
                {selectedBlock.minedBy && (
                  <div className="pb-4 border-b border-gray-200">
                    <div className="flex items-center justify-between mb-1">
                      <p className="text-sm text-gray-600">Mined By</p>
                      <button
                        onClick={() => copyToClipboard(selectedBlock.minedBy, 'miner')}
                        className="p-1 hover:bg-gray-100 rounded"
                      >
                        {copiedField === 'miner' ? (
                          <Check className="w-4 h-4 text-green-600" />
                        ) : (
                          <Copy className="w-4 h-4 text-gray-600" />
                        )}
                      </button>
                    </div>
                    <p className="text-sm text-gray-800 font-mono break-all">
                      {selectedBlock.minedBy}
                    </p>
                  </div>
                )}

                {/* Transactions */}
                <div>
                  <p className="text-sm text-gray-600 mb-2">Transactions</p>
                  <p className="text-2xl font-bold text-blue-600">
                    {selectedBlock.transactions?.length || 0}
                  </p>
                  {selectedBlock.transactions && selectedBlock.transactions.length > 0 && (
                    <div className="mt-4 space-y-2">
                      <p className="text-xs font-semibold text-gray-700 uppercase">
                        Transaction Hashes
                      </p>
                      {selectedBlock.transactions.map((tx, index) => (
                        <div
                          key={index}
                          onClick={() => navigate(`/transaction/${tx.hash}`)}
                          className="p-2 bg-gray-50 rounded border border-gray-200 hover:border-blue-300 cursor-pointer transition"
                        >
                          <p className="text-xs font-mono text-gray-600 truncate">{tx.hash}</p>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            ) : (
              <div className="text-center py-12 text-gray-500">
                Select a block to view details
              </div>
            )}
          </div>
        </div>
      </div>
      </div>
    </>
  );
};

export default BlockchainExplorer;
