//SPDX-License-Identifier: none
pragma solidity ^0.8.0;

import "../node_modules/@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "../node_modules/@openzeppelin/contracts/access/Ownable.sol";
import "../node_modules/@openzeppelin/contracts/utils/Counters.sol";
import "../node_modules/@openzeppelin/contracts/token/ERC721/extensions/ERC721Burnable.sol";
import "../node_modules/@openzeppelin/contracts/access/AccessControlEnumerable.sol";

contract CAToken is Context, ERC721Burnable, AccessControlEnumerable, Ownable {
    using Counters for Counters.Counter;
    //expiring token count
    //might not need this
    bytes32 public constant SERVER_ROLE = keccak256("SERVER_ROLE");
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
    }
    
    function mint(address _to) public {
        require(
            hasRole(SERVER_ROLE, _msgSender()),
            "you don't have access to the minting function"
        );
        //require unique tokenID
        _mint(_to, _tokenIdTracker.current());
        //emit information about the token
        emit Minted (_to, _tokenIdTracker.current());
        _tokenIdTracker.increment();
    }

    //delete expired contracts
    function expiredContracts(uint256[] memory deleteIds) public {
        //only allow the server to do this
        require(
            hasRole(SERVER_ROLE, _msgSender()),
            "this is awkward, you're not allowed to do that"
        );
        //batch delete the tokens
        for (uint256 i = 0; i <= deleteIds.length; i++) {
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
