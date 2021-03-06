//SPDX-License-Identifier: none
pragma solidity ^0.8.0;

import "./contracts/token/ERC721/ERC721.sol";
import "./contracts/access/Ownable.sol";
import "./contracts/utils/Counters.sol";
import "./contracts/access/AccessControlEnumerable.sol";
import "./contracts/token/ERC721/extensions/ERC721URIStorage.sol";


interface ExpiredContracts{
    function expiredContracts(uint256[] memory deleteIds) external;
} 

contract CAToken is 
    Context, 
    AccessControlEnumerable, Ownable, ERC721URIStorage {
    using Counters for Counters.Counter;
    //mapping of used nonces 
    mapping(uint256 => bool) public usedNonce;
    bytes32 public constant SERVER_ROLE = keccak256("SERVER_ROLE");
    Counters.Counter public nextNonce;
    //use counters to make unique tokenID
    Counters.Counter private _tokenIdTracker;
    //mint event to be emitted in the mint function 
    event Minted(address indexed _from, uint256 indexed tokenID);
    //event to tell server the tokens are deleted 
    event DeletedTokens(uint256[] indexed deleteIds);
    
    constructor(
        string memory name,
        string memory symbol,
        address _server
    ) ERC721(name, symbol) Ownable() {
        //give the owner admin role
        _setupRole(DEFAULT_ADMIN_ROLE, _msgSender());
        //set up server role
        _setupRole(SERVER_ROLE, _server);
        //lets give owner server role too
        _setupRole(SERVER_ROLE, _msgSender());
        nextNonce.increment();
    }
    
    function mint(address _to, string memory _resourceID, uint256 _nonce) public virtual{
        require(
            hasRole(SERVER_ROLE, _msgSender()),
            "you don't have access to the minting function"
        );
        require(!usedNonce[_nonce], "nonce is already used by a previous mint call");
        //require unique tokenID
        _mint(_to, _tokenIdTracker.current());
        //set the Resource ID 
        _setTokenURI(_tokenIdTracker.current(), _resourceID);
        usedNonce[_nonce] = true;
        //emit information about the token
        _tokenIdTracker.increment();
        nextNonce.increment();
    }
    function addServerRole(address _serverAddress) public onlyOwner {
        _setupRole(SERVER_ROLE, _serverAddress);
    }
    //delete expired contracts
    function expiredContracts(uint256[] memory deleteIds) public virtual  {
        //only allow the server to do this
        require(
            hasRole(SERVER_ROLE, _msgSender()),
            "this is awkward, you're not allowed to do that"
        );
        //batch delete the tokens
        for (uint256 i = 0; i < deleteIds.length; i++) {
            _burn(deleteIds[i]);
        }
        emit DeletedTokens(deleteIds); 
    }

    function selfDedstruct() public onlyOwner {
        address payable owner = payable(owner());
        selfdestruct(owner);
    }

    function supportsInterface(bytes4 interfaceId)
        public
        view
        virtual
        override(AccessControlEnumerable, ERC721)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}
