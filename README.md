# Blockchain Cryptocurrency Wallet System

A professional, production-ready decentralized cryptocurrency wallet application featuring custom blockchain implementation, UTXO model, Proof-of-Work mining, automated Zakat deduction, and modern React UI with advanced animations.

## 🌟 Features

### Core Blockchain Features
- ✅ **Custom Blockchain** - SHA-256, Merkle trees, Proof-of-Work consensus
- ✅ **UTXO Model** - Bitcoin-style unspent transaction outputs
- ✅ **Digital Signatures** - RSA-2048 cryptographic security
- ✅ **Mining System** - Adjustable difficulty with block rewards
- ✅ **Chain Validation** - Complete blockchain integrity verification
- ✅ **Block Explorer** - Real-time blockchain visualization

### Advanced Features
- 🔐 **Security** - AES-256-GCM encryption, JWT authentication, rate limiting
- 💰 **Zakat System** - Automated 2.5% monthly Islamic charity deductions
- 📊 **Analytics** - Comprehensive logging, reports, and transaction tracking
- 🎨 **Modern UI** - Glassmorphism, gradient effects, smooth animations
- 📱 **Responsive Design** - Mobile-first approach with Tailwind CSS
- 🚀 **Performance** - Optimized Go backend, indexed MongoDB queries

### Initial User Balance
- 💵 **1000 BC** - Every new account receives 1000 BC initial balance to start transactions immediately

## 🏗️ Tech Stack

### Backend
- **Language**: Go 1.25+
- **Framework**: Gin (HTTP router)
- **Database**: MongoDB Atlas
- **Authentication**: JWT tokens
- **Email**: SMTP (Gmail)
- **Security**: RSA-2048, AES-256-GCM, bcrypt

### Frontend
- **Framework**: React 18
- **Styling**: Tailwind CSS
- **Routing**: React Router v6
- **State**: Context API
- **HTTP Client**: Axios
- **Notifications**: React Hot Toast
- **Icons**: Lucide React

## 📋 Project Structure

```
BC_Project/
├── backend/                 # Go server
│   ├── config/             # Database & indexes
│   ├── models/             # Data structures
│   ├── crypto/             # Cryptography (RSA, AES, SHA-256)
│   ├── services/           # Business logic
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # Auth, rate limiting, sanitization
│   ├── routes/             # API routes
│   ├── main.go             # Entry point
│   ├── go.mod              # Dependencies
│   └── .env.example        # Environment template
│
├── frontend/               # React app
│   ├── src/
│   │   ├── pages/         # Page components
│   │   ├── components/    # Reusable UI components
│   │   ├── contexts/      # React contexts
│   │   ├── utils/         # API & utilities
│   │   ├── assets/        # Static assets
│   │   ├── App.jsx        # Root component
│   │   └── main.jsx       # Entry point
│   ├── public/            # Public assets
│   ├── package.json       # Dependencies
│   ├── vite.config.js     # Vite configuration
│   └── .env.example       # Environment template
│
├── README.md              # This file
├── QUICK_START.md         # Quick reference guide
└── .gitignore             # Git ignore rules
```

## 🚀 Quick Start

### Prerequisites

- **Go** 1.25+ - [Download](https://go.dev/dl/)
- **Node.js** 18+ and npm - [Download](https://nodejs.org/)
- **MongoDB Atlas Account** - [Create Free](https://www.mongodb.com/cloud/atlas)
- **Gmail Account** (for SMTP) - For OTP emails
- **Git** - [Download](https://git-scm.com/)

### Step 1: Clone Repository

```bash
git clone <your-repository-url>
cd BC_Project
```

### Step 2: MongoDB Setup

1. Create free MongoDB Atlas cluster at https://www.mongodb.com/cloud/atlas
2. Create database user with read/write permissions
3. Whitelist your IP (or use 0.0.0.0/0 for development)
4. Get connection string from "Connect" → "Connect your application"

### Step 3: Gmail SMTP Setup (For OTP Emails)

1. Enable 2-factor authentication on your Gmail account
2. Go to Google Account → Security → 2-Step Verification
3. Generate App Password:
   - Select App: "Mail"
   - Select Device: "Other" (name it "Blockchain Wallet")
   - Copy the generated 16-character password
4. Use this app password in backend .env file

### Step 4: Backend Setup

```bash
# Navigate to backend
cd backend

# Install dependencies
go mod download

# Create .env file
cp .env.example .env

# Edit .env file with your credentials
# Windows: notepad .env
# Mac/Linux: nano .env
```

**Required .env variables:**
```env
# MongoDB
MONGODB_URI=mongodb+srv://username:password@cluster.mongodb.net/

# JWT Secret (generate a strong random string)
JWT_SECRET=your-super-secret-jwt-key-min-32-chars

# SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-16-char-app-password
FROM_EMAIL=your-email@gmail.com
FROM_NAME=Blockchain Wallet

# Server Configuration
PORT=8080
ENVIRONMENT=development

# Blockchain Settings
MINING_DIFFICULTY=4
MINING_REWARD=50
ZAKAT_PERCENTAGE=2.5

# Encryption (generate with command below)
AES_ENCRYPTION_KEY=<32-byte-base64-key>
```

**Generate AES Encryption Key:**
```powershell
# PowerShell
$bytes = New-Object byte[] 32
[System.Security.Cryptography.RandomNumberGenerator]::Fill($bytes)
$key = [Convert]::ToBase64String($bytes)
Write-Host "AES_ENCRYPTION_KEY=$key"
```

```bash
# Linux/Mac
openssl rand -base64 32
```

**Start Backend:**
```bash
go run main.go
```

Backend will start on `http://localhost:8080`

### Step 5: Frontend Setup

```bash
# Open new terminal
cd frontend

# Install dependencies
npm install

# Create .env file
cp .env.example .env

# Edit .env file
# Windows: notepad .env
# Mac/Linux: nano .env
```

**Frontend .env:**
```env
VITE_API_URL=http://localhost:8080/api
```

**Start Frontend:**
```bash
npm run dev
```

Frontend will start on `http://localhost:5173`

### Step 6: First Time Setup

1. Open browser: `http://localhost:5173`
2. Click "Sign Up" to register
3. Enter email and click "Verify Email"
4. Check your email for OTP (6-digit code)
5. Enter OTP and complete registration
6. **IMPORTANT**: Copy and save your private key securely
7. Login with your credentials
8. You'll start with **1000 BC** balance
# Open new terminal, navigate to frontend
cd frontend

# Install dependencies
npm install

# Create environment file
cp .env.example .env


# Start development server
npm run dev
```

Frontend will start on `http://localhost:5173`

## 🔑 Environment Variables

### Backend (.env)

```env
# MongoDB Configuration
MONGODB_URI=mongodb+srv://username:password@cluster.mongodb.net/blockchain_wallet

# Server Configuration
PORT=8080
ENVIRONMENT=development

# JWT Authentication
JWT_SECRET=your-super-secret-jwt-key-at-least-32-characters-long

# Email/SMTP Configuration (for OTP verification)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-gmail-app-password-16-chars
FROM_EMAIL=your-email@gmail.com
FROM_NAME=Blockchain Wallet

# Blockchain Configuration
MINING_DIFFICULTY=4
MINING_REWARD=50
ZAKAT_PERCENTAGE=2.5
ZAKAT_POOL_WALLET_ID=ZAKAT_POOL_WALLET

# Security & Encryption
AES_ENCRYPTION_KEY=your-32-byte-base64-encoded-encryption-key
```

### Frontend (.env)

```env
# API Configuration
VITE_API_URL=http://localhost:8080/api
```

## 📚 API Endpoints

### Authentication (Public)
```
POST   /api/register                    - Register new user
POST   /api/login                       - User login
POST   /api/otp/generate                - Generate OTP for email verification
POST   /api/otp/verify                  - Verify OTP code
GET    /api/otp/check                   - Check email verification status
```

### Blockchain Explorer (Public)
```
GET    /api/blockchain                  - Get entire blockchain
GET    /api/blockchain/stats            - Blockchain statistics
GET    /api/blockchain/validate         - Validate blockchain integrity
GET    /api/block/hash/:hash            - Get block by hash
GET    /api/block/index/:index          - Get block by index
GET    /api/block/latest                - Get latest block
GET    /api/wallet/validate/:walletId   - Validate wallet ID
GET    /api/transaction/:hash           - Get transaction by hash
GET    /api/transactions/pending        - Get pending transactions
```

### User & Wallet (Protected - Requires JWT)
```
GET    /api/profile                     - Get user profile
PUT    /api/profile                     - Update user profile
GET    /api/wallet                      - Get wallet details
GET    /api/balance                     - Get current balance
GET    /api/wallet/utxos                - Get wallet UTXOs
POST   /api/beneficiary                 - Add beneficiary
DELETE /api/beneficiary/:walletId       - Remove beneficiary
```

### Transactions (Protected)
```
POST   /api/transaction                 - Create new transaction
GET    /api/transactions                - Get transaction history
POST   /api/mine                        - Mine new block (manual)
```

### Zakat (Protected)
```
GET    /api/zakat/history               - Get Zakat deduction history
GET    /api/zakat/summary               - Get Zakat summary
POST   /api/zakat/deduct                - Trigger manual Zakat deduction
```

### Logging & Reports (Protected)
```
GET    /api/logs/system                 - Get system logs
GET    /api/logs/transactions           - Get transaction logs
GET    /api/logs/transactions/all       - Get all transaction logs
GET    /api/reports                     - Get user reports
```

## 🎨 UI Features

### Modern Design Elements
- 🌈 **Gradient Backgrounds** - Smooth color transitions
- 💎 **Glassmorphism** - Frosted glass effects
- ✨ **Animations** - Smooth transitions, hover effects, loading states
- 🎭 **Micro-interactions** - Button animations, card scales, icon transforms
- 📊 **Data Visualization** - Transaction charts, balance displays
- 🌊 **Blob Animations** - Floating background elements
- 🎯 **Focus States** - Enhanced form interactions

### Pages
- 🔐 **Login/Register** - Animated authentication with email OTP
- 📊 **Dashboard** - Balance overview, quick actions, recent transactions
- 💸 **Send Money** - Real-time wallet validation, amount checks
- 📜 **Transactions** - Filterable history with search
- ⛓️ **Blockchain Explorer** - Block details, transaction lookup
- 📈 **Reports** - Zakat history, transaction analytics
- 👤 **Profile** - User information management

## 🔐 Security Features

### Cryptography
- **RSA-2048** - Public/private key generation
- **AES-256-GCM** - Private key encryption
- **SHA-256** - Blockchain hashing
- **Digital Signatures** - Transaction verification
- **Bcrypt** - Password hashing (cost factor 10)

### API Security
- **JWT Tokens** - Stateless authentication
- **Rate Limiting** - 20-100 req/min based on endpoint
- **Input Sanitization** - SQL injection & XSS prevention
- **CORS** - Cross-origin protection
- **Validation** - Request payload validation

### Data Security
- **Encrypted Private Keys** - AES-256-GCM encryption
- **Secure Sessions** - HTTP-only cookies (if applicable)
- **Password Strength** - Minimum length, complexity requirements
- **CNIC Validation** - Format verification (12345-1234567-1)

## 📊 Database Schema

### MongoDB Collections

**users** - User accounts
- ID, FullName, Email, Password (hashed), CNIC
- WalletID, PublicKey, PrivateKey (encrypted)
- Beneficiaries, ZakatTracking
- CreatedAt, UpdatedAt

**wallets** - Wallet information
- WalletID (primary), UserID, PublicKey
- Balance (cached), CreatedAt, UpdatedAt, IsActive

**utxos** - Unspent transaction outputs
- ID, TransactionHash, OutputIndex
- WalletID, Amount, Spent
- CreatedAt

**transactions** - All transactions
- Hash (primary), Sender, Receiver
- Amount, Fee, Signature
- Inputs (UTXOs), Outputs
- BlockHash, Timestamp, Confirmed

**blocks** - Blockchain blocks
- Hash (primary), Index, Timestamp
- PreviousHash, MerkleRoot, Nonce, Difficulty
- Transactions[], MinerWallet

**zakatDeductions** - Zakat records
- ID, UserID, WalletID
- Amount, DeductedAt, TransactionHash

**systemLogs** - System events
- ID, EventType, Message
- UserID, IPAddress, Timestamp

**transactionLogs** - Transaction events
- ID, TransactionHash, EventType
- SenderWalletID, ReceiverWalletID
- Amount, Status, Timestamp

### Indexes
- 30+ database indexes for optimized queries
- Unique constraints on Email, WalletID, Transaction Hash
- Compound indexes for transaction history queries

## 🛠️ Development

### Backend Development
```powershell
cd backend
go run main.go
```

### Frontend Development
```powershell
cd frontend
npm run dev
```

### Build for Production

**Backend:**
```powershell
cd backend
go build -o wallet-api.exe main.go
```

**Frontend:**
```powershell
cd frontend
npm run build
```

## ⚙️ Configuration

### Adjust Mining Difficulty
In `backend/.env`:
```env
MINING_DIFFICULTY=4  # Lower = faster mining (3-5 recommended)
```

### Adjust Zakat Percentage
In `backend/.env`:
```env
ZAKAT_PERCENTAGE=2.5  # Default 2.5%
```

## 📝 Project Requirements Checklist

- ✅ Custom blockchain implementation
- ✅ Real cryptographic digital signatures
- ✅ UTXO-based transaction model
- ✅ Proof-of-Work mining
- ✅ Wallet ID from public/private keys
- ✅ Domestic money transfer system
- ✅ Monthly 2.5% Zakat deduction
- ✅ System logging and analytics
- ✅ React + Tailwind frontend
- ✅ Go backend
- ✅ Email + OTP verification
- ✅ User profiles with CNIC
- ✅ Beneficiary list management
- ✅ Transaction history
- ✅ Block explorer
- ✅ Reports and summaries

## 🐛 Troubleshooting

### Backend won't start
- Check MongoDB credentials in `.env`
- Ensure Go 1.25+ is installed
- Verify port 8080 is available

### Frontend won't connect
- Verify backend is running on port 8080
- Check VITE_API_URL in frontend `.env`
- Clear browser cache

### Mining too slow
- Reduce MINING_DIFFICULTY in backend `.env`
- Try difficulty 3 or 4 instead of 5


## 👥 Support

For issues or questions:
1. Check troubleshooting section
2. Check environment variables
3. Verify all dependencies are installed

## 📄 License

Educational project for academic purposes.

---

