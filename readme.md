# 🚀 Ethereum CLI Wallet (Golang)

A simple **Ethereum Wallet CLI** built in **Go** that allows users to:
✅ **Create multiple wallets**  
✅ **Import wallets using private keys**  
✅ **Check ETH balance**  
✅ **Send ETH transactions**  
✅ **Fetch transaction history**  
✅ **Secure wallets using AES encryption**  

---

## 📌 Features
- **🔑 Multi-Wallet Support** → Manage multiple Ethereum wallets.
- **🔒 Secure Storage** → Encrypts private keys with AES (requires an encryption key).
- **⚡ Fast Transactions** → Send ETH instantly using Infura RPC.
- **📊 Transaction History** → Fetch past transactions using Etherscan.
- **🐳 Docker Support** → Run in a **containerized** environment.

---

## ⚙️ 1️⃣ Installation & Setup
### **Prerequisites**
- **Go 1.18+** installed → [Download](https://golang.org/dl/)
- **Docker (optional)** if running in a container
- **Infura & Etherscan API keys** for Ethereum network access

### **Clone the Repository**
```sh
git clone https://github.com/Kol1aUA/go-eth-wallet.git
cd go-eth-wallet


## Setup Env Variables
export WALLET_ENCRYPTION_KEY="your_32_byte_secure_key"
export ETH_RPC_URL="https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID"
export ETHERSCAN_API_KEY="your_etherscan_api_key_here"
