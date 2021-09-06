package blockchain

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	CAToken "github.com/z-j-lin/nft/tree/main/nft-backend/pkg/Token"
)

type Contract struct {
	ContractAddress common.Address
	eth             *Ethereum
	Instance        *CAToken.CAToken
	MintEvent       string
	DeleteEvent     string
}

func (c *Contract) init() {
	instance, err := CAToken.NewCAToken(c.ContractAddress, c.eth.Client)
	if err != nil {
		log.Fatalf("failed to initiate contract instance: %v", err)
	}
	LogMintedSig := []byte("Minted(address,uint256)")
	LogDeletedTokensSig := []byte("DeletedTokens(uint256[])")
	LogMintedSigHash := crypto.Keccak256Hash(LogMintedSig).Hex()
	LogDeletedTokensSigHash := crypto.Keccak256Hash(LogDeletedTokensSig).Hex()
	c.MintEvent = LogMintedSigHash
	c.DeleteEvent = LogDeletedTokensSigHash
	c.Instance = instance

}
func (c *Contract) MintToken(Auth *bind.TransactOpts, RecipientAddr common.Address) (*types.Transaction, error) {
	tx, err := c.Instance.Mint(Auth, RecipientAddr)
	if err != nil {
		log.Println("failed to send transaction @ MintToken:", err)
		return nil, err
	}
	return tx, err
}

func (c *Contract) DeleteTokens(auth *bind.TransactOpts, IDs []*big.Int) (*types.Transaction, error) {
	//gasPrice, err := c.eth.Client.SuggestGasPrice(context.Background())
	//fromAddress := c.eth.Account.Address
	//nonce, err := c.eth.Client.PendingNonceAt(context.Background(), fromAddress)
	//options for transaction
	//auth.Nonce = big.NewInt(int64(nonce))
	//auth.Value = big.NewInt(0)
	//auth.GasLimit = uint64(300000)
	//auth.GasPrice = gasPrice
	tx, err := c.Instance.ExpiredContracts(auth, IDs)
	if err != nil {
		log.Println("failed to send transaction @ DeleteTokens:", err)
		return nil, err
	}
	return tx, err
}

//method to check if owner has token??
