package gameserver

import (
	"github.com/sydnash/lotou/log"
	"github.com/sydnash/lotou/simple/btype"
	"github.com/sydnash/lotou/simple/global"
	"math/rand"
)

type MjLogicInfo struct {
	marker      int32
	dice1       int32
	dice2       int32
	curCQPos    int32
	mjCQVec     [108]int32
	curMoPaiPos int32
	deskOpDesc  [4]*PosOpDesc
	dc          *DeskConrol
}

func (mj *MjLogicInfo) init() {
	for k, _ := range mj.deskOpDesc {
		mj.deskOpDesc[k] = newPosOpDesc()
	}
}

func (mj *MjLogicInfo) Start() {
	mj.marker = 0
	mj.dice1 = rand.Int31n(6) + 1
	mj.dice2 = rand.Int31n(6) + 1
	mj.curCQPos = 0
	mj.curMoPaiPos = mj.marker
	mj.disturbMJ()
}

func (mj *MjLogicInfo) disturbMJ() {
	for i, _ := range mj.mjCQVec {
		mj.mjCQVec[i] = int32(i)
	}
	for i := 0; i < 108; i++ {
		pos := rand.Int31n(int32(108-i)) + int32(i)
		mj.mjCQVec[i], mj.mjCQVec[pos] = mj.mjCQVec[pos], mj.mjCQVec[i]
	}
}

func (mj *MjLogicInfo) getMj(cnt int32) []int32 {
	ret := mj.mjCQVec[mj.curCQPos : mj.curCQPos+cnt]
	mj.curCQPos += cnt
	return ret
}
func (mj *MjLogicInfo) getLeftMjCnt() int32 {
	return 108 - mj.curCQPos
}
func (mj *MjLogicInfo) addOpDesc(op *OpDesc, pos int32) {
	mj.deskOpDesc[pos].ops = append(mj.deskOpDesc[pos].ops, op)
}

type MJLogicPosInfo struct {
	shouPai []int32
	queId   int32
}

func (this MJLogicPosInfo) chuPai(mjId int32) int32 {
	var pos int
	for k, v := range this.shouPai {
		if v == mjId {
			pos = k
			break
		}
	}
	if pos == len(this.shouPai)-1 {
		this.shouPai = this.shouPai[:len(this.shouPai)-1]
	} else {
		mjId = this.shouPai[pos]
		this.shouPai = append(this.shouPai[:pos], this.shouPai[pos+1:len(this.shouPai)]...)
	}
	return mjId
}

func NewMJLogicPosInfo() *MJLogicPosInfo {
	ret := &MJLogicPosInfo{}
	ret.shouPai = make([]int32, 0, 14)
	ret.queId = -1
	return ret
}
func (self *MJLogicPosInfo) fapai(mjs []int32) {
	log.Debug("fapai : %v", mjs)
	self.shouPai = append(self.shouPai, mjs...)
}

func (dc *DeskConrol) getNextPos(cur int32) int32 {
	var i int32
	for i = 0; i < 4; i++ {
		pos := i + cur + 1
		if pos >= 4 {
			pos -= 4
		}
		if dc.posInfos[pos].hasPeople {
			return pos
		}
	}
	return -1
}

func (dc *DeskConrol) faPai() {
	curPos := dc.mjLogicInfo.marker
	var i int32
	for i = 0; i < 4; i++ {
		pos := i + curPos
		if pos >= 4 {
			pos -= 4
		}
		deskPosInfo := dc.posInfos[pos]
		if deskPosInfo.hasPeople {
			deskPosInfo.mjLogicPosInfo = NewMJLogicPosInfo()
			deskPosInfo.mjLogicPosInfo.fapai(dc.mjLogicInfo.getMj(13))
		}
	}
	dc.posInfos[curPos].mjLogicPosInfo.fapai(dc.mjLogicInfo.getMj(1))

	for _, info := range dc.posInfos {
		if info.hasPeople {
			client := info.client
			client.gs.encoder.Reset()
			head := btype.PHead{}
			head.Type = btype.S_MSG_FAPAI
			client.gs.encoder.Encode(head)
			client.gs.encoder.Encode(dc.posInfos[dc.mjLogicInfo.marker].client.acId)
			client.gs.encoder.Encode(dc.mjLogicInfo.dice1)
			client.gs.encoder.Encode(dc.mjLogicInfo.dice2)
			client.gs.encoder.Encode(int32(dc.curHasPeople))
			for _, v := range dc.posInfos {
				if v.hasPeople {
					client.gs.encoder.Encode(v.client.acId)
					client.gs.encoder.Encode(v.client.playerInfo.GetPropertyInt64(global.KPropertyTypeQueBiDC))
					mjCnt := len(v.mjLogicPosInfo.shouPai)
					client.gs.encoder.Encode(mjCnt)
					for i := 0; i < mjCnt; i++ {
						if v.client.acId == client.acId {
							client.gs.encoder.Encode(v.mjLogicPosInfo.shouPai[i])
						} else {
							client.gs.encoder.Encode(-1)
						}
					}
				}
			}
			client.gs.encoder.Encode(0)
			client.gs.sendToAgent(client.agentId)
		}
	}
}

func (dc *DeskConrol) checkIsCanStart() {
	var clientNum int16 = 0
	for _, v := range dc.posInfos {
		if v.hasPeople && v.isReady {
			clientNum++
		}
	}
	if clientNum == dc.rc.roomInfo.StartMinMax {
		dc.mjLogicInfo.Start()
		dc.faPai()
		dc.deskState = KDeskStateDingQue
	}
}
func (dc *DeskConrol) dingQue(client *GameClient) {
	if dc.deskState != KDeskStateDingQue {
		return
	}
	pos := client.deskPos
	info := dc.posInfos[pos]
	var que int32
	client.gs.decoder.Decode(&que)

	if que >= 0 && que <= 2 {
		if info.mjLogicPosInfo.queId < 0 {
			info.mjLogicPosInfo.queId = que
			dc.checkIsAllDingQue()
		}
	}
}

func (dc *DeskConrol) checkIsAllDingQue() {
	var gs *GameService
	for _, v := range dc.posInfos {
		if v.hasPeople {
			gs = v.client.gs
			if v.mjLogicPosInfo.queId < 0 {
				return
			}
		}
	}

	gs.encoder.Reset()
	head := btype.PHead{}
	head.Type = btype.S_MSG_DINGQUE
	gs.encoder.Encode(head)
	gs.encoder.Encode(int32(dc.curHasPeople))
	for _, v := range dc.posInfos {
		if v.hasPeople {
			gs.encoder.Encode(v.client.acId)
			gs.encoder.Encode(v.mjLogicPosInfo.queId)
		}
	}

	for _, v := range dc.posInfos {
		if v.hasPeople {
			gs.sendToAgent(v.client.agentId)
		}
	}

	dc.checkOPAfterMoPai()
}

const (
	KOPTypeNone = iota
	KOPTypeChuPai
	KOPTypePeng
	KOPTypeGang
	KOPTypeHu
	KOPTypePass
	_
	_
	KOPTypeDuoXiang
)

type OpDesc struct {
	OpType  int32
	SubType int32
	mjIdxs  []int32
	bePos   []int32
	whoPos  int32
}

func newOpDesc() *OpDesc {
	ret := &OpDesc{}
	ret.mjIdxs = make([]int32, 0, 4)
	ret.bePos = make([]int32, 0, 4)
	return ret
}

type PosOpDesc struct {
	isChoosed bool
	choosedOp *OpDesc
	mjId      int32
	ops       []*OpDesc
}

func (this *PosOpDesc) hasOp() bool {
	if len(this.ops) > 0 {
		return true
	}
	return false
}
func (this *PosOpDesc) getOpByOpType(optype int32) *OpDesc {
	for _, op := range this.ops {
		if op.OpType == optype {
			return op
		}
	}
	return nil
}
func (this *PosOpDesc) chooseOp(op int32) bool {
	if this.isChoosed {
		return false
	}
	this.choosedOp = this.getOpByOpType(op)
	if this.choosedOp != nil {
		this.isChoosed = true
		return true
	}
	return false
}

func newPosOpDesc() *PosOpDesc {
	ret := &PosOpDesc{}
	ret.ops = make([]*OpDesc, 0)
	return ret
}

func (dc *DeskConrol) clearAllOP() {
	for _, v := range dc.mjLogicInfo.deskOpDesc {
		v.isChoosed = false
		v.choosedOp = nil
		v.mjId = -1
		v.ops = v.ops[:0]
	}
}
func (dc *DeskConrol) checkOPAfterMoPai() {
	dc.clearAllOP()
	curMoPaiPos := dc.mjLogicInfo.curMoPaiPos
	opDesc := newOpDesc()
	opDesc.OpType = KOPTypeChuPai
	opDesc.whoPos = curMoPaiPos
	dc.mjLogicInfo.addOpDesc(opDesc, curMoPaiPos)

	dc.sendOpHint()
}
func (dc *DeskConrol) sendOpHint() {
	for _, v := range dc.posInfos {
		if v.hasPeople {
			posOpDesc := dc.mjLogicInfo.deskOpDesc[v.pos]
			if len(posOpDesc.ops) > 0 {
				gs := v.client.gs
				gs.encoder.Reset()
				head := btype.PHead{}
				head.Type = btype.S_MSG_OPHINT
				gs.encoder.Encode(head)
				gs.encoder.Encode(len(posOpDesc.ops))
				for _, op := range posOpDesc.ops {
					gs.encoder.Encode(op.OpType)
					gs.encoder.Encode(len(op.mjIdxs))
					for _, idx := range op.mjIdxs {
						gs.encoder.Encode(idx)
					}
				}
				gs.sendToAgent(v.client.agentId)
			}
		}
	}
}
func (dc *DeskConrol) opDo(client *GameClient) {
	var opType int32
	var mjId int32
	client.gs.decoder.Decode(&opType)
	client.gs.decoder.Decode(&mjId)
	switch opType {
	case KOPTypeChuPai:
		dc.opChuPai(client, mjId)
	}
}

func (dc *DeskConrol) opChuPai(client *GameClient, mjId int32) {
	pos := client.deskPos
	posOpDesc := dc.mjLogicInfo.deskOpDesc[pos]
	posOpDesc.chooseOp(KOPTypeChuPai)
	if posOpDesc.choosedOp != nil {
		posOpDesc.mjId = mjId
		mjId = dc.posInfos[pos].mjLogicPosInfo.chuPai(mjId)
		gs := client.gs
		gs.encoder.Reset()
		head := btype.PHead{}
		head.Type = btype.S_MSG_OPDO
		gs.encoder.Encode(head)
		gs.encoder.Encode(1)
		gs.encoder.Encode(client.acId)
		gs.encoder.Encode(KOPTypeChuPai)
		gs.encoder.Encode(false)
		gs.encoder.Encode(mjId)
		gs.encoder.Encode(0)
		gs.encoder.Encode(1)
		gs.encoder.Encode(mjId)
		gs.encoder.Encode(0)
		gs.encoder.Encode(0)
		for _, v := range dc.posInfos {
			if v.hasPeople {
				gs.sendToAgent(v.client.agentId)
			}
		}
		nextPos := dc.getNextPos(pos)
		moPaiIds := dc.mjLogicInfo.getMj(1)
		leftCnt := dc.mjLogicInfo.getLeftMjCnt()
		dc.posInfos[nextPos].mjLogicPosInfo.fapai(moPaiIds)
		dc.mjLogicInfo.curMoPaiPos = nextPos

		nc := dc.posInfos[nextPos]
		gs.encoder.Reset()
		head.Type = btype.S_MSG_MOPAI
		gs.encoder.Encode(head)
		gs.encoder.Encode(nc.client.acId)
		gs.encoder.Encode(1)
		gs.encoder.Encode(moPaiIds[0])
		gs.encoder.Encode(leftCnt)
		for _, v := range dc.posInfos {
			if v.hasPeople {
				gs.sendToAgent(v.client.agentId)
			}
		}

		dc.checkOPAfterMoPai()
	}
}
