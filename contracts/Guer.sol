// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;
pragma experimental ABIEncoderV2;
import './Mia.sol';

contract Guer {

    address docId;
    mapping(address => address[]) MiasByUser;
    address payable ourWallet;
    constructor() payable {
        ourWallet = payable(msg.sender);
    }

    event keyData(
        address indexed _from, //assume sender = key1
        bytes32 indexed _key1,
        bytes32 indexed _key2  //in the frontend code, if your address doesn't match sender, then key 2 is yours.
    );

    function doPayRoyalty() public payable returns(bool){
        if (msg.value >= .0001 ether)
        {
            ourWallet.transfer(msg.value);
            return true;
        }
        else {
            return false;
        }
    }

    function createNFT() public payable returns (bool){
        ourWallet.transfer(msg.value);
        Mia a = new Mia();
        MiasByUser[msg.sender].push(address(a));        
        return true;
    }

    function getUserNFTs() public view returns (address[] memory){
        return MiasByUser[msg.sender];
    }

    function getNFTInfo(address payable docAddress) public view returns (address, uint, uint, string memory){
        return Mia(docAddress).getMia();
    }

    function getGuerOwnerAddress() public view returns(address){
        return ourWallet;
    }

    function getKeystore(address payable docAddress) public view returns (string memory) {
        return Mia(docAddress).getKeystore();
    }

    function getEncryptedCount (address payable docAddress) public view returns (uint) {
        return Mia(docAddress).getEncryptedCount();
    }

    function getEncrypted(address payable docAddress, uint _id) public view returns (string memory, string memory, uint, string memory, string memory, string memory) {
        return Mia(docAddress).getEncryptedDetails(_id);
    }
    
    function doAddKeystore(address payable docAddress, string memory _hash) public {
        Mia(docAddress).addKeystore(_hash);
    }

    function doAddEncrypted(address payable docAddress, string memory _location, string memory _name, string memory _filetype, string memory _companion, string memory _usage) public {
        Mia(docAddress).addEncrypted(_location, _name, _filetype, _companion, _usage);
    }
    
    function postKeyData(address _sender, bytes32 _key1, bytes32 _key2) external {
        emit keyData(_sender, _key1, _key2);
    }

    function doTransferFunds(address payable docAddress, address payable _receiver, uint256 _amount) public payable {
        Mia(docAddress).transferFunds(_receiver, _amount);
    }

}
