import { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import api from '../utils/api';
import Navbar from '../components/Navbar';
import { User, Mail, Wallet, CreditCard, Edit2, Save, X, Plus, Trash2, Users } from 'lucide-react';
import toast from 'react-hot-toast';

const Profile = () => {
  const { userProfile, refreshProfile } = useAuth();
  const [isEditing, setIsEditing] = useState(false);
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState({
    fullName: '',
    email: '',
    cnic: '',
  });
  const [beneficiaries, setBeneficiaries] = useState([]);
  const [showAddBeneficiary, setShowAddBeneficiary] = useState(false);
  const [newBeneficiary, setNewBeneficiary] = useState({
    name: '',
    relationship: '',
    percentage: '',
  });

  useEffect(() => {
    if (userProfile) {
      setFormData({
        fullName: userProfile.fullName || '',
        email: userProfile.email || '',
        cnic: userProfile.cnic || '',
      });
      setBeneficiaries(userProfile.beneficiaries || []);
    }
  }, [userProfile]);

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleEdit = () => {
    setIsEditing(true);
  };

  const handleCancel = () => {
    setIsEditing(false);
    setFormData({
      fullName: userProfile.fullName || '',
      email: userProfile.email || '',
      cnic: userProfile.cnic || '',
    });
  };

  const handleSave = async () => {
    setLoading(true);
    try {
      await api.put('/profile', {
        fullName: formData.fullName,
        cnic: formData.cnic,
      });
      toast.success('Profile updated successfully');
      await refreshProfile();
      setIsEditing(false);
    } catch (error) {
      const message = error.response?.data?.error || 'Failed to update profile';
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  const handleAddBeneficiary = async () => {
    if (!newBeneficiary.name || !newBeneficiary.relationship || !newBeneficiary.percentage) {
      toast.error('Please fill all beneficiary fields');
      return;
    }

    const percentage = parseFloat(newBeneficiary.percentage);
    if (percentage <= 0 || percentage > 100) {
      toast.error('Percentage must be between 1 and 100');
      return;
    }

    const totalPercentage = beneficiaries.reduce((sum, b) => sum + b.percentage, 0);
    if (totalPercentage + percentage > 100) {
      toast.error('Total beneficiary percentage cannot exceed 100%');
      return;
    }

    setLoading(true);
    try {
      const updatedBeneficiaries = [
        ...beneficiaries,
        {
          ...newBeneficiary,
          percentage: percentage,
        },
      ];

      await api.put('/profile/beneficiaries', {
        beneficiaries: updatedBeneficiaries,
      });

      toast.success('Beneficiary added successfully');
      setBeneficiaries(updatedBeneficiaries);
      setNewBeneficiary({ name: '', relationship: '', percentage: '' });
      setShowAddBeneficiary(false);
      await refreshProfile();
    } catch (error) {
      const message = error.response?.data?.error || 'Failed to add beneficiary';
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  const handleRemoveBeneficiary = async (index) => {
    setLoading(true);
    try {
      const updatedBeneficiaries = beneficiaries.filter((_, i) => i !== index);

      await api.put('/profile/beneficiaries', {
        beneficiaries: updatedBeneficiaries,
      });

      toast.success('Beneficiary removed successfully');
      setBeneficiaries(updatedBeneficiaries);
      await refreshProfile();
    } catch (error) {
      const message = error.response?.data?.error || 'Failed to remove beneficiary';
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  const totalBeneficiaryPercentage = beneficiaries.reduce((sum, b) => sum + b.percentage, 0);

  return (
    <>
      <Navbar />
      <div className="min-h-screen bg-gray-50 py-8">
        <div className="max-w-4xl mx-auto px-4">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-800">Profile</h1>
          <p className="text-gray-600 mt-2">Manage your account information and beneficiaries</p>
        </div>

        {/* Profile Information */}
        <div className="bg-white rounded-xl shadow-lg p-8 mb-6">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-semibold text-gray-800">Personal Information</h2>
            {!isEditing ? (
              <button
                onClick={handleEdit}
                className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
              >
                <Edit2 className="w-4 h-4" />
                Edit
              </button>
            ) : (
              <div className="flex gap-2">
                <button
                  onClick={handleCancel}
                  className="flex items-center gap-2 px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition"
                >
                  <X className="w-4 h-4" />
                  Cancel
                </button>
                <button
                  onClick={handleSave}
                  disabled={loading}
                  className="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition disabled:bg-gray-400"
                >
                  <Save className="w-4 h-4" />
                  Save
                </button>
              </div>
            )}
          </div>

          <div className="space-y-4">
            {/* Full Name */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                <div className="flex items-center gap-2">
                  <User className="w-4 h-4" />
                  Full Name
                </div>
              </label>
              <input
                type="text"
                name="fullName"
                value={formData.fullName}
                onChange={handleChange}
                disabled={!isEditing}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100 disabled:cursor-not-allowed"
              />
            </div>

            {/* Email */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                <div className="flex items-center gap-2">
                  <Mail className="w-4 h-4" />
                  Email
                </div>
              </label>
              <input
                type="email"
                name="email"
                value={formData.email}
                disabled
                className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-100 cursor-not-allowed"
              />
              <p className="text-xs text-gray-500 mt-1">Email cannot be changed</p>
            </div>

            {/* CNIC */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                <div className="flex items-center gap-2">
                  <CreditCard className="w-4 h-4" />
                  CNIC
                </div>
              </label>
              <input
                type="text"
                name="cnic"
                value={formData.cnic}
                onChange={handleChange}
                disabled={!isEditing}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100 disabled:cursor-not-allowed"
              />
            </div>

            {/* Wallet ID */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                <div className="flex items-center gap-2">
                  <Wallet className="w-4 h-4" />
                  Wallet ID
                </div>
              </label>
              <input
                type="text"
                value={userProfile?.walletId || ''}
                disabled
                className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-100 cursor-not-allowed font-mono text-sm"
              />
              <p className="text-xs text-gray-500 mt-1">Your unique wallet identifier</p>
            </div>
          </div>
        </div>

        {/* Beneficiaries */}
        <div className="bg-white rounded-xl shadow-lg p-8">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h2 className="text-xl font-semibold text-gray-800">Beneficiaries</h2>
              <p className="text-sm text-gray-600 mt-1">
                Manage your beneficiaries for inheritance purposes
              </p>
            </div>
            {!showAddBeneficiary && (
              <button
                onClick={() => setShowAddBeneficiary(true)}
                className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
              >
                <Plus className="w-4 h-4" />
                Add Beneficiary
              </button>
            )}
          </div>

          {/* Add Beneficiary Form */}
          {showAddBeneficiary && (
            <div className="bg-blue-50 rounded-lg p-4 mb-6">
              <h3 className="font-semibold text-gray-800 mb-4">Add New Beneficiary</h3>
              <div className="space-y-3">
                <input
                  type="text"
                  placeholder="Name"
                  value={newBeneficiary.name}
                  onChange={(e) =>
                    setNewBeneficiary({ ...newBeneficiary, name: e.target.value })
                  }
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
                <input
                  type="text"
                  placeholder="Relationship (e.g., Son, Daughter, Spouse)"
                  value={newBeneficiary.relationship}
                  onChange={(e) =>
                    setNewBeneficiary({ ...newBeneficiary, relationship: e.target.value })
                  }
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
                <input
                  type="number"
                  placeholder="Percentage (1-100)"
                  value={newBeneficiary.percentage}
                  onChange={(e) =>
                    setNewBeneficiary({ ...newBeneficiary, percentage: e.target.value })
                  }
                  min="1"
                  max="100"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
                <div className="flex gap-2">
                  <button
                    onClick={handleAddBeneficiary}
                    disabled={loading}
                    className="flex-1 bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 transition disabled:bg-gray-400"
                  >
                    Add
                  </button>
                  <button
                    onClick={() => {
                      setShowAddBeneficiary(false);
                      setNewBeneficiary({ name: '', relationship: '', percentage: '' });
                    }}
                    className="flex-1 bg-gray-200 text-gray-700 py-2 rounded-lg hover:bg-gray-300 transition"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            </div>
          )}

          {/* Beneficiaries List */}
          {beneficiaries.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              <Users className="w-12 h-12 mx-auto mb-3 text-gray-400" />
              <p>No beneficiaries added yet</p>
            </div>
          ) : (
            <>
              <div className="space-y-3 mb-4">
                {beneficiaries.map((beneficiary, index) => (
                  <div
                    key={index}
                    className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:border-blue-300 transition"
                  >
                    <div className="flex-1">
                      <h4 className="font-semibold text-gray-800">{beneficiary.name}</h4>
                      <p className="text-sm text-gray-600">{beneficiary.relationship}</p>
                    </div>
                    <div className="flex items-center gap-4">
                      <div className="text-right">
                        <p className="text-2xl font-bold text-blue-600">
                          {beneficiary.percentage}%
                        </p>
                      </div>
                      <button
                        onClick={() => handleRemoveBeneficiary(index)}
                        className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition"
                      >
                        <Trash2 className="w-5 h-5" />
                      </button>
                    </div>
                  </div>
                ))}
              </div>

              {/* Total Percentage */}
              <div className="bg-gray-50 rounded-lg p-4">
                <div className="flex items-center justify-between">
                  <span className="font-semibold text-gray-700">Total Allocation</span>
                  <span
                    className={`text-2xl font-bold ${
                      totalBeneficiaryPercentage === 100
                        ? 'text-green-600'
                        : totalBeneficiaryPercentage > 100
                        ? 'text-red-600'
                        : 'text-yellow-600'
                    }`}
                  >
                    {totalBeneficiaryPercentage}%
                  </span>
                </div>
                {totalBeneficiaryPercentage !== 100 && (
                  <p className="text-sm text-gray-600 mt-2">
                    {totalBeneficiaryPercentage > 100
                      ? '⚠️ Total percentage exceeds 100%'
                      : `💡 Remaining ${100 - totalBeneficiaryPercentage}% to allocate`}
                  </p>
                )}
              </div>
            </>
          )}
        </div>
      </div>
      </div>
    </>
  );
};

export default Profile;
