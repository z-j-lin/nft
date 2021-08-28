package contract

import (
	"fmt"

	"github.com/z-j-lin/nft/tree/main/nft-backend/pkg/monitor"
)

func main() {
	//newDB instance
	rdb, err := monitor.NewDBinstance()
	if err != nil {
		fmt.Errorf("unable to connect to db: %v", err)
	}

	//start a connection to the redis server

	//check redis transaction list for stuff

}
