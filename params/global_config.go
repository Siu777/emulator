package params

var (
	Block_Interval      = 5000        // generate new block interval
	MaxBlockSize_global = 500         // the block contains the maximum number of transactions
	InjectSpeed         = 5000        // the transaction inject speed
	TotalDataSize       = 100000      // the total number of txs
	DataWrite_path      = "./result/" // measurement data result output path
	LogWrite_path       = "./log"     // log output path
	// SupervisorAddr      = "192.168.137.1:28800"
	SupervisorAddr = "127.0.0.1:18800"                         //supervisor ip address
	FileInput      = "./2000000to2999999_BlockTransaction.csv" //the raw BlockTransaction data path
)
