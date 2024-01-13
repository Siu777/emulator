package measure

import (
	"blockEmulator/message"
	"blockEmulator/params"
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

// to test average Transaction_Confirm_Latency (TCL)  in this system
type TestModule_TCL_Relay struct {
	epochID           int
	totTxLatencyEpoch []float64          // record the Transaction_Confirm_Latency in each epoch
	txNum             []float64          // record the txNumber in each epoch
	txLatencyMap      map[string]float64 // record the tx confirm latency
	txlatency1        map[string]float64 // record the tx1 latency. add to pool ~ tx1 committed
	txlatency2        map[string]float64 // record the tx2 latency
	isCrossTx         map[string]bool    // whether ctx
}

func NewTestModule_TCL_Relay() *TestModule_TCL_Relay {
	return &TestModule_TCL_Relay{
		epochID:           -1,
		totTxLatencyEpoch: make([]float64, 0),
		txNum:             make([]float64, 0),
		txLatencyMap:      make(map[string]float64),
		txlatency1:        make(map[string]float64),
		txlatency2:        make(map[string]float64),
		isCrossTx:         make(map[string]bool),
	}
}

func (tml *TestModule_TCL_Relay) OutputMetricName() string {
	return "Transaction_Confirm_Latency"
}

// modified latency
func (tml *TestModule_TCL_Relay) UpdateMeasureRecord(b *message.BlockInfoMsg) {
	if b.BlockBodyLength == 0 { // empty block
		return
	}

	epochid := b.Epoch
	txs := b.ExcutedTxs
	mTime := b.CommitTime

	// extend
	for tml.epochID < epochid {
		tml.txNum = append(tml.txNum, 0)
		tml.totTxLatencyEpoch = append(tml.totTxLatencyEpoch, 0)
		tml.epochID++
	}

	for _, tx := range b.Relay1Txs {
		tml.isCrossTx[string(tx.TxHash)] = true
		tml.txlatency1[string(tx.TxHash)] = mTime.Sub(tx.Time).Seconds()
	}

	for _, tx := range txs {
		if !tx.Time.IsZero() {
			tml.txLatencyMap[string(tx.TxHash)] = mTime.Sub(tx.Time).Seconds()
			tml.totTxLatencyEpoch[epochid] += mTime.Sub(tx.Time).Seconds()
			tml.txNum[epochid]++
			if val, ok := tml.isCrossTx[string(tx.TxHash)]; ok && val {
				tml.txlatency2[string(tx.TxHash)] = tml.txLatencyMap[string(tx.TxHash)] - tml.txlatency1[string(tx.TxHash)]
			}
		}
	}
}

func (tml *TestModule_TCL_Relay) HandleExtraMessage([]byte) {}

func (tml *TestModule_TCL_Relay) OutputRecord() (perEpochLatency []float64, totLatency float64) {
	perEpochLatency = make([]float64, 0)
	latencySum := 0.0
	totTxNum := 0.0
	for eid, totLatency := range tml.totTxLatencyEpoch {
		perEpochLatency = append(perEpochLatency, totLatency/tml.txNum[eid])
		latencySum += totLatency
		totTxNum += tml.txNum[eid]
	}
	totLatency = latencySum / totTxNum

	writeToCSV(tml.txLatencyMap, tml.txlatency1, tml.txlatency2, tml.isCrossTx)
	return
}

func writeToCSV(txLatency, t1l, t2l map[string]float64, isCrossTx map[string]bool) {
	dirpath := params.DataWrite_path + "supervisor_measureOutput/"
	err := os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
	targetPath := dirpath + "txlatency.csv"
	file, er := os.Create(targetPath)
	if er != nil {
		panic(er)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	title := []string{"txid", "single_TCL", "isCtx", "ctx1", "ctx2"}
	w.Write(title)
	w.Flush()
	datanum := 0
	for key, val := range txLatency {
		isctx_label := ""
		tx1Latency := ""
		tx2Latency := ""
		if isCrossTx[key] {
			isctx_label = "1"
			tx1Latency = strconv.FormatFloat(t1l[key], 'f', 8, 64)
			tx2Latency = strconv.FormatFloat(t2l[key], 'f', 8, 64)
		}
		resultStr := []string{string(rune(datanum)), strconv.FormatFloat(val, 'f', 8, 64), isctx_label, tx1Latency, tx2Latency}
		datanum++
		err = w.Write(resultStr)
		if err != nil {
			log.Panic()
		}
		w.Flush()
	}
}
