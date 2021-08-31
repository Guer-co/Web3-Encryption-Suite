# Web3 Encryption Suite (W3ES) 
![logo](https://raw.githubusercontent.com/Guer-co/Web3-Encryption-Suite/master/public/w3es_blawky.png)

Simple Encryption Suite built for Web3 using ECDSA keypairs for Identity, Key Generation, and Key distribution. 

[Introduction](#introduction) - [Getting Started](#getting-started) - [Support & Contact](#support-&-contact) - [License](#license)

## Introduction

A lot of effort has been put into building Web3's decentralized infrastructure. Connecting to Ethereum and other services such as decentralized storage or applications relies on decades-old protocols such as TLS, which in turn require heavily centralized services like Certificate Authorities. This existing centralized web security puts web3 at risk.

These protocols are limited by their dependence on location-based addressing (IP). As the web evolved, centralized services and infrastructure emerged to compensate for the limitations of IP addressing. Most of these services exist to associate and verify identity with location. In addition to Certificate Authorities, these now range from multi-factor authentication schemes (e-mail, SMS), to cookies (validated devices). Each additional service adds attack surface and vulnerability, which is amplified by centralization.

Web3 Encryption Suite (W3ES) uses cryptographic addressing (ECDSA key pairs, such as ETH wallets) as the root digital identity for encryption protocols. This allows Web3 services and applications to circumvent vulnerabilities introduced by IP addressing while still using established cryptographic libraries and protocols.

### Key Features

- Digital Identity Verification using Smart Contracts, and Web3 services such as ENS
- Secure key generation, using verified WASM/eWASM, or remotely in networked Trusted Execution Environment (e.g. Intel SGX)
- Decentralized Key Distribution, via public blockchain networks
- Support for Ethereum and SKALE (more coming soon)
- Integrations with decentralized storage networks (IPFS, Arweave, SIA, Archon, Storj)
- 

### Encryption Services

Currently, the suite consists of three unique encryption services:

- Single-party Encryption
- Multi-party Encryption
- Proxy Re-encryption

The encryption code is originally written in GoLang, and makes use of Go-Ethereum's crypto libraries. It is compiled to WASM for execution in browser, as well as compiled for IntelSGX. 

#### Single-Party Encryption

With single-party encryption, only one address can verify and retrieve data. This is useful for privately storing data on public networks. As this is single party, compiled and verified code is executed in WebAssembly:

*Encryption:*
1. Read plaintext data / file into memory buffer
2. Generate random AES key
3. Encrypt plaintext using AES key
4. Sign ciphertext (ECDSA)
5. Create "Companion" struct, including Signer's Public Key, Ciphertext Signature, Ciphertext Hash, AES Key
6. Encrypt "Companion" struct using Public Key (ECIES)
7. Store encrypted data and companion on public network

*Decryption:*

-   Retrieve encrypted file and "Companion"
-   Decrypt "Companion" using Crypto Identity (ECIES)
-   Verify Signature + Data
-   Recover AES Key
-   Decrypt file
-   Local use of decrypted file

#### Multiparty Encryption

Multiparty Encryption enables multiple parties to securely exchange information without any dependence on centralized or trusted infrastructure.

This is useful for securely connecting to both Web2 and Web3 sites using a modified OpenSSL library, End-to-End Encrypted Messaging over LibP2P, Shared access to files stored on decentralized storage networks, etc. 

This happens in three key phases:

1. Identity Verification through Request and Confirmation
2. Key Generation and Distribution
3. Encryption and Decryption 

![enter image description here](https://guer.co/static/Image_2@2x-cba6f4d76c8758c8b09d27d72bb09a3f.png)

***Identity Verification through Request and Confirmation***

-   Alice forms Request to connect using Bob's Address
-   Bob verifies request, and forms Confirmation

***Key Generation and Distribution***

Within Trusted Execution Environment:
-   Receives Request + Confirmation Messages
-   Verifies both messages for authenticity
-   Generates random AES-128 Key
-   Encrypts AES-128 Key with Alice's Public Key (ECIES)
-   Encrypts AES-128 Key with Bob's Public Key (ECIES)
-   Publishes encrypted keys on-chain

***Symmetric Encryption & Decryption***

*Encryption*
-   Retrieves Encrypted AES Key from public network or storage
-   Decrypts AES Key using Private Key A (ECIES)
-   Encrypts file with AES Key
-   Publishes file

*Decryption*
-   Retrieves Encrypted AES Key from public network or storage
-   Decrypts AES Key using Private Key B (ECIES)
-   Decrypts file with AES Key
-   Prints to terminal + local plaintext file

#### Proxy Re-Encryption (PRES)

*This encryption suite is still in development*

Securely distribute content and manage access using networked Trusted Execution Environments

***Modified Handshake***

-   Submit Request with Pubkey, Signature, Hash of CID
-   PRE Service (SGX) verifies signature, request
-   PRE Service (SGX) generates new keystore
-   PRE Service (SGX) Confirms with SGX Pubkey, Sig, Hash
-   Publishes file

***Proxy Re-encryption Request***

-   Submit request with CID, destination address, AES Key
-   Encrypt PRE Request with SGX Pubkey

***Decentralized Proxy Re-Encryption***

-   PRES receives proxy request
-   PRES decrypts proxy request (ECIES)
-   PRES retrieves file from CID
-   PRES decrypts file using provided AES key
-   PRES generates new AES key
-   PRES encrypts plaintext file with new AES key
-   PRES creates Companion file for target address(es) (ECDSA)
-   PRES encrypts Companion file using target Pubkey (ECIES)
-   PRES publishes newly encrypted file + Companion

### Project Status

- Encryption suite, compiled in WASM for testing, Single- and Multi-party
- Decentralized Authority 1.0 Smart Contracts deployed on Ropsten + SKL Testnet
- GETH compatible wallets via keystore files and local wallets
- APIs with IPFS, SKALE, Arweave, ETH Ropsten 
- Basic Public Docs + Demos

### In Development

- Ethers.JS/Web3.JS Wallet APIs
- Networked SGX Testing for Multi-party Trustless Keygen + PRES  
- Decentralized Authority Contracts v2.0  
- Additional APIs: OpenSSL, ENS, Opera/Brave, LibP2P  

## Getting Started

### Installation

git clone git@github.com:Guer-co/Web3-Encryption-Suite.git  
cd Web3-Encryption-Suite  
npm i && npm start  

### Documentation

[documentation](https://docs.guer.co)

### Dependencies, Credits, and Contributions



### Examples

[Encryption example](https://guer.us/)

[Manage your decentralized data app example](https://skale.guer.co/)

[Stateless ecommerce page](https://future.guer.co/)

## Support & Contact

Gitcoin Grants: https://gitcoin.co/grants/1261/guer-empowering-data-custody-for-web3-through-sec  
Substack: [guer.substack.com](guer.substack.com)  
Website: [guer.co](guer.co)  
ETH/ERC20: 0x13541d356bCC7156e234fF15091BA9B8f1D9B2dd  

## License

*This software is under development and should be used at your own risk. As it is largely experimental, do not use it for any sensitive or otherwise important data*

For now, all rights reserved by Guer Labs, LLC 2021. Upon review and audit, the codebase will be released as Free-and-open-source, with responsible limits to permissible modifications for compatibility and security. 

Any questions, please contact us directly via our website. 
