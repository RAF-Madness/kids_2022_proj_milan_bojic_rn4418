package message

import (
	"distributed/job"
	"distributed/node"
	"encoding/json"
	"fmt"
	"sync/atomic"
)

type MessageType int32

const (
	Info                      MessageType = 0
	InfoBroadcast             MessageType = 1
	Hail                      MessageType = 2
	Contact                   MessageType = 3
	Welcome                   MessageType = 4
	Join                      MessageType = 5
	Leave                     MessageType = 6
	Entered                   MessageType = 7
	ConnectionRequest         MessageType = 8
	ConnectionResponse        MessageType = 9
	Quit                      MessageType = 10
	ClusterKnock              MessageType = 11
	EnterCluster              MessageType = 12
	ClusterConnectionRequest  MessageType = 13
	ClusterConnectionResponse MessageType = 14
	JobSharing                MessageType = 15
	ImageInfoRequest          MessageType = 16
	ImageInfo                 MessageType = 17
	SystemKnock               MessageType = 18
	Purge                     MessageType = 19
	ShareJob                  MessageType = 20
	StartJob                  MessageType = 21
	ApproachCluster           MessageType = 22
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
	GetMessage() string
}

type Message struct {
	MessageType    MessageType   `json:"MessageType"`
	OriginalSender node.NodeInfo `json:"sender"`
	Reciver        node.NodeInfo `json:"reciver"`
	Route          []int         `json:"route"`
	Message        string        `json:"Message"`
	Id             int64         `json:"id"`
}

func (msg *Message) String() string {
	return "Message"
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

func (msg *Message) GetMessage() string {
	return msg.Message
}

func (msg *Message) Log() string {
	return fmt.Sprintf("%d¦%d¦%d¦%d¦%s", msg.OriginalSender.Id, msg.Reciver.Id, msg.Id, msg.MessageType, msg.Message)
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

func MakeInfoMessage(sender, reciver node.INode, message string) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = message
	msgReturn.MessageType = Info

	msgReturn.OriginalSender = *sender.GetNodeInfo()
	msgReturn.Reciver = *reciver.GetNodeInfo()

	msgReturn.Route = []int{sender.GetId()}

	return &msgReturn
}

func MakeInfoBroadcastMessage(sender node.INode, message string) *Message {
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
	contact_byte, _ := json.Marshal(contact)
	msgReturn.Message = string(contact_byte)
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
	msgb, _ := json.Marshal(msgMap)
	msgReturn.Message = string(msgb)
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
	msgReturn.Message = fmt.Sprint(sender.Id)
	msgReturn.MessageType = Join

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeLeaveMessage(sender, reciver node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = fmt.Sprint(sender.Id)
	msgReturn.MessageType = Leave

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeEnteredMessage(sender node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())

	msgjson, _ := json.Marshal(sender)

	msgReturn.Message = string(msgjson)
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
	msgReturn.Message = string(smer)
	msgReturn.MessageType = ConnectionRequest

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeConnectionResponseMessage(sender, reciver node.NodeInfo, accepted bool, smer ConnectionSmer) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = fmt.Sprintf("%t:%v", accepted, smer)
	msgReturn.MessageType = ConnectionResponse

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeQuitMessage(sender node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = fmt.Sprint(sender.Id)
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

func MakeClusterEnterMessage(sender, reciver node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = "EnterCluster"
	msgReturn.MessageType = EnterCluster

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
	msgReturn.Message = fmt.Sprint(accepted)
	msgReturn.MessageType = ClusterConnectionResponse

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeClusterJobSharingMessage(sender, reciver node.NodeInfo, jobInfo string) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = jobInfo
	msgReturn.MessageType = JobSharing

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

func MakeImageInfoMessage(sender, reciver node.NodeInfo, points [][]int) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	points_json, _ := json.Marshal(points)
	msgReturn.Message = string(points_json)
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

func MakeShereJobMessage(sender node.NodeInfo, jobInput job.Job) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	jsonstr, _ := json.Marshal(jobInput)
	msgReturn.Message = string(jsonstr)
	msgReturn.MessageType = ShareJob

	msgReturn.OriginalSender = sender
	tmpReciver := new(node.NodeInfo)
	tmpReciver.Id = -1
	msgReturn.Reciver = *tmpReciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeStartJobMessage(sender, reciver node.NodeInfo, jobName string) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	msgReturn.Message = jobName
	msgReturn.MessageType = StartJob

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}

func MakeApproachClusterMessage(sender, reciver, contact node.NodeInfo) *Message {
	msgReturn := Message{}

	msgReturn.Id = int64(MainCounter.Inc())
	jsonstr, _ := json.Marshal(contact)
	msgReturn.Message = string(jsonstr)
	msgReturn.MessageType = ApproachCluster

	msgReturn.OriginalSender = sender
	msgReturn.Reciver = reciver

	msgReturn.Route = []int{sender.Id}

	return &msgReturn
}
