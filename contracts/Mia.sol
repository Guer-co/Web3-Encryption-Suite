// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;
pragma experimental ABIEncoderV2;

contract Mia {

    fallback() external payable {}
    receive() external payable {}

    address id;
    uint creationDate;
    address[] myAddresses;
    string keystore;
    address creator;

    mapping(uint => EncryptedFile) public encrypteds;
    uint public encryptedCount;

    struct EncryptedFile {
        string location;
        string name;
        uint date;
        string filetype;
        string companion;
        string usage;
    }

    constructor() payable {
        id = address(this);
        creationDate = block.timestamp;
        creator = msg.sender;
    }

    function getKeystore() public view returns (string memory) {
        return (keystore);
    }

    function getEncryptedDetails(uint _id) public view returns (string memory, string memory, uint, string memory, string memory, string memory) {
        return (encrypteds[_id].location, encrypteds[_id].name, encrypteds[_id].date, encrypteds[_id].filetype, encrypteds[_id].companion, encrypteds[_id].usage);
    }

    function getMia() public view returns (address, uint, uint, string memory) {
        return (id, creationDate, encryptedCount, keystore);
    }

    function getEncryptedCount() public view returns (uint) {
        return (encryptedCount);
    }

    function addKeystore(string memory _hash) public {
        keystore = _hash;
    }

    function addEncrypted(string memory _location, string memory _name, string memory _filetype, string memory _companion, string memory _usage) public {
        encryptedCount = encryptedCount + 1;
        encrypteds[encryptedCount].location = _location;
        encrypteds[encryptedCount].name = _name;
        encrypteds[encryptedCount].date = block.timestamp;
        encrypteds[encryptedCount].filetype = _filetype;
        encrypteds[encryptedCount].companion = _companion;
        encrypteds[encryptedCount].usage = _usage;
    }

    function transferFunds(address payable _receiver, uint256 _amount) pure external {
        payable(_receiver).call{value:_amount};
    }

}