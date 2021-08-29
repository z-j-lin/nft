package monitor

type txqmon struct {
	db Database
}

func NewTXQmon(rdb Database) *txqmon {
	Mintque := make(chan *blockchain.transaction, 3)
	return &txqmon{
		db: rdb,
	}
}

//function to start the transaction que monitoring loop
func (qmon *txqmon) startTXQmonLoop() {
	go qmon.loop()
}
func (qmon *txqmon) loop() {
	for {
		//check if the transaction queue has a job
		account, resourceID := qmon.db.DQmint()
		if account != "" {
			//initiate a new transaction
			//load tx into buffered channel
		}

	}
}
