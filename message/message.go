package message

import (
	"distributed/job"
	"distributed/node"
	"distributed/structures"
	"fmt"
	"sync/atomic"
)

// func firstN(s string, n int) string {
// 	i := 0
// 	for j := range s {
// 		if i == n {
// 			return s[:j]
// 		}
// 		i++
// 	}
// 	return s
// }

type MessageType string

const (
	Info                      MessageType = "Info"
	InfoBroadcast             MessageType = "InfoBroadcast"
	Hail                      MessageType = "Hail"
	Contact                   MessageType = "Contact"
	Welcome                   MessageType = "Welcome"
	Join                      MessageType = "Join"
	Leave                     MessageType = "Leave"
	Entered                   MessageType = "Entered"
	ConnectionRequest         MessageType = "ConnectionRequest"
	ConnectionResponse        MessageType = "ConnectionResponse"
	Quit                      MessageType = "Quit"
	ClusterKnock              MessageType = "ClusterKnock"
	EnteredCluster            MessageType = "EnteredCluster"
	ClusterConnectionRequest  MessageType = "ClusterConnectionRequest"
	ClusterConnectionResponse MessageType = "ClusterConnectionResponse"
	ImageInfoRequest          MessageType = "ImageInfoRequest"
	ImageInfo                 MessageType = "ImageInfo"
	SystemKnock               MessageType = "SystemKnock"
	Purge                     MessageType = "Purge"
	StartJob                  MessageType = "StartJob"
	StartJobGenesis           MessageType = "StartJobGenesis"
	ApproachCluster           MessageType = "ApproachCluster"
	ClusterWelcome            MessageType = "ClusterWelcome"
	StopShareJob              MessageType = "StopShareJob"
	StoppedJobInfo            MessageType = "StoppedJobInfo"
	AskForJob                 MessageType = "AskForJob"
	JobStatusRequest          MessageType = "JobStatusRequest"
	JobStatus                 MessageType = "JobStatus"
	UpdatedNode               MessageType = "UpdatedNode"
)

type MessageCounter struct {
	counter int32
}

func (cnt *MessageCounter) Inc() int32 {
	return atomic.AddInt32(&cnt.counter, 1)
}

func (cnt *MessageCounter) Dec() int32 {
	return atomic.AddInt32(&cnt.counter, -1)
}

func (cnt *MessageCounter) Get() int32 {
	return atomic.LoadInt32(&cnt.counter)
}

var MainCounter = MessageCounter{0}

type IMessage interface {
	String() string
	MakeMeASender(node node.INode) IMessage
	Effect(args interface{})
	Log() string
	GetSender() node.NodeInfo
	GetReciver() node.NodeInfo
	GetRoute() []int
	GetMessage() any
}

type Message struct {
	MessageType    MessageType   `json:"MessageType"`
	OriginalSender node.NodeInfo `json:"sender"`
	Reciver        node.NodeInfo `json:"reciver"`
	Route          []int         `json:"route"`
	Message        any           `json:"Message"`
	Id             int64         `json:"id"`
}

func (msg *Message) String() string {
	return fmt.Sprintf("%d¦%d¦%d¦%s", msg.OriginalSender.Id, msg.Reciver.Id, msg.Id, msg.MessageType)
}

func (msg *Message) Effect(args interface{}) {
}

func (msg *Message) GetSender() node.NodeInfo {
	return msg.OriginalSender
}

func (msg *Message) GetReciver() node.NodeInfo {
	return msg.Reciver
}

func (msg *Message) GetRoute() []int {
	return msg.Route
}

func (msg *Message) GetMessage() any {
	return msg.Message
}

func (msg *Message) Log() string {
	return fmt.Sprintf("%d¦%d¦%d¦%s¦%v", msg.OriginalSender.Id, msg.Reciver.Id, msg.Id, msg.MessageType, msg.Message)
}

func (msg *Message) MakeMeASender(node node.INode) IMessage {

	msgReturn := Message{}
	msgReturn.Id = msg.Id
	msgReturn.Message = msg.Message
	msgReturn.MessageType = msg.MessageType

	msgReturn.OriginalSender = msg.OriginalSender
	msgReturn.Reciver = msg.Reciver

	msgReturn.Route = append(msg.Route, node.GetId())

	return &msgReturn

}

func MakeInfoMessage(sender, reciver node.INode, message any) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = message
	msgReturn.MessageType = Info

	msgReturn.OriginalSender = *sender.GetNodeInfo()
	msgReturn.Reciver = *reciver.GetNodeInfo()

	msgReturn.Route = []int{sender.GetId()}

	return &msgReturn
}

func MakeInfoBroadcastMessage(sender node.INode, message any) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = message
	msgReturn.MessageType = InfoBroadcast

	msgReturn.OriginalSender = *sender.GetNodeInfo()
	msgReturn.Reciver = *new(node.NodeInfo)

	msgReturn.Route = []int{sender.GetId()}

	return &msgReturn
}

func MakeHailMessage(sender node.Worker, reciver node.Bootstrap) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = "Hail"
	msgReturn.MessageType = Hail

	msgReturn.OriginalSender = *sender.GetNodeInfo()
	msgReturn.Reciver = *reciver.GetNodeInfo()

	msgReturn.Route = []int{sender.GetId()}

	return &msgReturn
}

func MakeContactMessage(sender node.NodeInfo, reciver, contact node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())

	msgReturn.Message = contact
	msgReturn.MessageType = Contact

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeWelcomeMessage(sender, reciver node.NodeInfo, nodeId int, systemInfo map[int]node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())

	msgMap := map[string]interface{}{"id": nodeId, "systemInfo": systemInfo}

	msgReturn.Message = msgMap
	msgReturn.MessageType = Welcome

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeSystemKnockMessage(sender, reciver node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = "SystemKnock"
	msgReturn.MessageType = SystemKnock

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeJoinMessage(sender, reciver node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = sender.Id
	msgReturn.MessageType = Join

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeLeaveMessage(sender, reciver node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = sender.Id
	msgReturn.MessageType = Leave

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeEnteredMessage(sender node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = sender
	msgReturn.MessageType = Entered

	msgReturn.OriginalSender = sender
	tmpReciver := new(node.NodeInfo)
	tmpReciver.Id = -1
	msgReturn.Reciver = *tmpReciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

type ConnectionSmer string

const (
	Next ConnectionSmer = "NEXT"
	Prev ConnectionSmer = "PREV"
)

func MakeConnectionRequestMessage(sender, reciver node.NodeInfo, smer ConnectionSmer) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = smer
	msgReturn.MessageType = ConnectionRequest

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeConnectionResponseMessage(sender, reciver node.NodeInfo, accepted bool, smer ConnectionSmer) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())

	var tmpMap = map[string]any{"accepted": accepted, "smer": smer}

	msgReturn.Message = tmpMap
	msgReturn.MessageType = ConnectionResponse

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeQuitMessage(sender node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = sender.Id
	msgReturn.MessageType = Quit

	msgReturn.OriginalSender = sender
	tmpReciver := new(node.NodeInfo)
	tmpReciver.Id = -1
	msgReturn.Reciver = *tmpReciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeClusterKnockMessage(sender, reciver node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = "ClusterKnock"
	msgReturn.MessageType = ClusterKnock

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeEnteredClusterMessage(sender, reciver, node node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())

	msgReturn.Message = node
	msgReturn.MessageType = EnteredCluster

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}
	return &msgReturn
}

func MakeClusterConnectionRequestMessage(sender, reciver node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = "ClusterConnectionRequest"
	msgReturn.MessageType = ClusterConnectionRequest

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeClusterConnectionResponseMessage(sender, reciver node.NodeInfo, accepted bool) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = accepted
	msgReturn.MessageType = ClusterConnectionResponse

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeImageInfoRequestMessage(sender, reciver node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = "ImageInfoRequest"
	msgReturn.MessageType = ImageInfoRequest

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeImageInfoMessage(sender, reciver node.NodeInfo, jobName string, points []structures.Point) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())

	outMap := map[string]interface{}{"jobName": jobName, "points": points}

	msgReturn.Message = outMap
	msgReturn.MessageType = ImageInfo

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakePurgeMessage(sender node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = "Purge"
	msgReturn.MessageType = Purge

	tmpReciver := new(node.NodeInfo)
	tmpReciver.Id = -1
	msgReturn.Reciver = *tmpReciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeStartJobMessage(sender, reciver node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = "StartJob"
	msgReturn.MessageType = StartJob

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeStartJobGenesisMessage(sender, reciver node.NodeInfo, jobName string) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = jobName
	msgReturn.MessageType = StartJobGenesis

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeApproachClusterMessage(sender, reciver, contact node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())

	msgReturn.Message = contact
	msgReturn.MessageType = ApproachCluster

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeClusterWelcomeMessage(sender, reciver node.NodeInfo, fractalID, jobName string) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())

	sentMap := map[string]string{"fractalID": fractalID, "jobName": jobName}

	msgReturn.Message = sentMap
	msgReturn.MessageType = ClusterWelcome

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeStopShareJobMessage(sender, reciver node.NodeInfo, jobINput job.Job) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())

	msgReturn.Message = jobINput
	msgReturn.MessageType = StopShareJob

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeStoppedJobInfoMessage(sender, reciver node.NodeInfo, jobName string, points []structures.Point) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())

	outMap := map[string]interface{}{"jobName": jobName, "points": points}

	msgReturn.Message = outMap

	msgReturn.MessageType = StoppedJobInfo

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeJobStatusRequestMessage(sender, reciver node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())

	msgReturn.Message = "JobStatusRequest"

	msgReturn.MessageType = JobStatusRequest

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeJobStatusMessage(sender, reciver node.NodeInfo, jobStatus job.JobStatus) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())

	msgReturn.Message = jobStatus

	msgReturn.MessageType = JobStatus

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeUpdatedNodeMessage(sender, nodeInput node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())

	msgReturn.Message = nodeInput

	msgReturn.MessageType = UpdatedNode

	msgReturn.OriginalSender = sender
	tmpReciver := new(node.NodeInfo)
	tmpReciver.Id = -1
	msgReturn.Reciver = *tmpReciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}
