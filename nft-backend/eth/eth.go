package eth

import (
	"log"
)

func eth(rpcurl string) {
	//connecting to the rpc server
	cl, err := ethclient.dial(rpcurl)
	if err {
		log.Fatalf("Failked to connect to the ether client: %v", err)
	}

}
func signtx() {

}
