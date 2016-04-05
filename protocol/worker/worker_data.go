// automatically generated, do not modify

package worker

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type worker_data struct {
	_tab flatbuffers.Table
}

func GetRootAsworker_data(buf []byte, offset flatbuffers.UOffsetT) *worker_data {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &worker_data{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *worker_data) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *worker_data) Cmd() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 2
}

func (rcv *worker_data) ClientIdf() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *worker_data) WorkerIdf() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *worker_data) Data() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func worker_dataStart(builder *flatbuffers.Builder) { builder.StartObject(4) }
func worker_dataAddCmd(builder *flatbuffers.Builder, cmd int16) { builder.PrependInt16Slot(0, cmd, 2) }
func worker_dataAddClientIdf(builder *flatbuffers.Builder, clientIdf flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(clientIdf), 0) }
func worker_dataAddWorkerIdf(builder *flatbuffers.Builder, workerIdf flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(workerIdf), 0) }
func worker_dataAddData(builder *flatbuffers.Builder, data flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(data), 0) }
func worker_dataEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
