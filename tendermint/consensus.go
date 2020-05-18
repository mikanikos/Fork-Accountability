package tendermint
//
//import "github.com/mikanikos/Fork-Accountability/common"
//
//type Consensus struct {
//	processId string
//
//	height      uint64
//	round       uint64
//	step        common.MessageType
//	decision    []*common.Value
//	lockedValue *common.Value
//	lockedRound int64
//	validValue  *common.Value
//	validRound  int64
//}
//
//func NewConsensus(id string) *Consensus {
//	return &Consensus{
//		processId:   id,
//		decision:    make([]*common.Value, 0),
//		lockedRound: -1,
//		validRound:  -1,
//	}
//}
//
//func (cons *Consensus) Start() {
//	cons.startRound(0)
//}
//
//func (cons *Consensus) startRound(round uint64) {
//	cons.round = round
//	cons.step = common.Proposal
//
//	if getProposer() == cons.processId {
//		if cons.validValue != nil {
//
//		}
//	}
//}
