package build

import (
	"blockEmulator/consensus_shard/pbft_all"
	"blockEmulator/params"
	"blockEmulator/supervisor"
	"strconv"
	"time"
)

//	func initConfig(nid, nnm, sid, snm uint64) *params.ChainConfig {
//		params.ShardNum = int(snm)
//		for i := uint64(0); i < snm; i++ {
//			if _, ok := params.IPmap_nodeTable[i]; !ok {
//				params.IPmap_nodeTable[i] = make(map[uint64]string)
//			}
//			if i < snm/2 {
//				for j := uint64(0); j < nnm; j++ {
//					params.IPmap_nodeTable[i][j] = "192.168.137.1:" + strconv.Itoa(20300+int(i)*100+int(j))
//				}
//			} else {
//				for j := uint64(0); j < nnm; j++ {
//					params.IPmap_nodeTable[i][j] = "192.168.137.195:" + strconv.Itoa(20300+int(i)*100+int(j))
//				}
//			}
//		}
//
// params.IPmap_nodeTable[9][0] = "192.168.137.195:15000"
func initConfig(nid, nnm, sid, snm uint64) *params.ChainConfig {
	params.ShardNum = int(snm)
	for i := uint64(0); i < snm; i++ {
		if _, ok := params.IPmap_nodeTable[i]; !ok {
			params.IPmap_nodeTable[i] = make(map[uint64]string)
		}
		for j := uint64(0); j < nnm; j++ {
			params.IPmap_nodeTable[i][j] = "127.0.0.1:" + strconv.Itoa(10300+int(i)*100+int(j))
		}
	}

	params.IPmap_nodeTable[params.DeciderShard] = make(map[uint64]string)
	params.IPmap_nodeTable[params.DeciderShard][0] = params.SupervisorAddr
	params.NodesinShard = int(nnm)
	params.ShardNum = int(snm)
	pcc := &params.ChainConfig{
		ChainID:        sid,
		NodeID:         nid,
		ShardID:        sid,
		Nodes_perShard: uint64(params.NodesinShard),
		ShardNums:      snm,
		BlockSize:      uint64(params.MaxBlockSize_global),
		BlockInterval:  uint64(params.Block_Interval),
		InjectSpeed:    uint64(params.InjectSpeed),
	}
	return pcc
}

func BuildSupervisor(nnm, snm, mod uint64) {
	var measureMod []string
	if mod == 0 || mod == 2 {
		measureMod = params.MeasureBrokerMod
	} else {
		measureMod = params.MeasureRelayMod
	}

	lsn := new(supervisor.Supervisor)
	lsn.NewSupervisor(params.SupervisorAddr, initConfig(123, nnm, 123, snm), params.CommitteeMethod[mod], measureMod...)
	time.Sleep(22 * time.Second)
	go lsn.SupervisorTxHandling()
	lsn.TcpListen()
}

func BuildNewPbftNode(nid, nnm, sid, snm, mod uint64) {
	worker := pbft_all.NewPbftNode(sid, nid, initConfig(nid, nnm, sid, snm), params.CommitteeMethod[mod])
	if nid == 0 {
		time.Sleep(10 * time.Second)
		go worker.Propose()
		worker.TcpListen()
	} else {
		worker.TcpListen()
	}
}
