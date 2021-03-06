package commander

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/piot/hasty-babel/definition"
	"github.com/piot/hasty-babel/serializers"
	"github.com/piot/hasty-protocol/channel"
	"github.com/piot/hasty-protocol/packet"
	"github.com/piot/hasty-protocol/packetdeserializers"
	"github.com/piot/hasty-protocol/packetserializers"
	"github.com/piot/hasty-protocol/serializer"
	"github.com/piot/hasty-protocol/version"
	"github.com/piot/hasty-repl/handler"
)

// Commander : todo
type Commander struct {
	conn                  *tls.Conn
	protocolDefinition    definition.ProtocolDefinition
	createStreamRequestID uint64
}

func NewCommander(conn *tls.Conn) Commander {
	com := Commander{conn: conn}
	go com.receiveConn()
	return com
}

func checkForPacket(stream *packet.Stream, data []byte) (packet.Packet, error) {
	hexPayload := hex.Dump(data)
	log.Printf("Received: %s", hexPayload)
	stream.Feed(data)
	return stream.FetchPacket()
}

func (in *Commander) receiveConn() {
	connectionIDAlwaysOne := packet.NewConnectionID(1)
	stream := packet.NewPacketStream(connectionIDAlwaysOne)
	packetHandler := replpackethandler.Handler{}

	for {
		buf := make([]byte, 1024)
		octetsRead, err := in.conn.Read(buf)
		if err != nil {
			log.Print(err)
			return
		} else {
			feedBuf := buf[:octetsRead]
			newPacket, packetErr := checkForPacket(&stream, feedBuf)
			if packetErr != nil {
				fmt.Printf("Deserialize error:%s", err)
			}
			deserializeErr := packetdeserializers.Deserialize(newPacket, &packetHandler)
			if deserializeErr != nil {
				return
			}
		}
	}
}

func (in *Commander) Connect(realm string) error {
	log.Printf("Commander connect realm '%s' ", realm)
	protocolVersion, _ := version.New(1, 0, 0)
	octets := packetserializers.ConnectToOctets(protocolVersion, realm)
	in.sendPacket(octets)
	return nil
}

func (in *Commander) Login(username string, password string) error {
	log.Printf("Commander login username '%s' password '%s'", username, password)
	octets := packetserializers.LoginToOctets(username, password)
	in.sendPacket(octets)
	return nil
}

func (in *Commander) SubscribeStream(channel channel.ID, offset uint32) error {
	log.Printf("Commander subscribe %s offset %d", channel, offset)
	octets := packetserializers.SubscribeStreamToOctets(channel, offset)
	in.sendPacket(octets)
	return nil
}

func (in *Commander) CreateStream(path string) error {
	log.Printf("Commander create %s", path)
	in.createStreamRequestID++
	octets := packetserializers.CreateStreamToOctets(in.createStreamRequestID, path)
	in.sendPacket(octets)
	return nil
}

func (in *Commander) LoadDefinition(path string) error {
	definition, err := definition.NewProtocolDefinitionFromFilePath(path)
	if err != nil {
		return err
	}

	log.Printf("Loaded definition '%s'", definition)
	in.protocolDefinition = definition
	return nil
}

func (in *Commander) PublishStream(channel channel.ID, chunk []byte) error {
	octets := packetserializers.PublishStreamToOctets(channel, chunk)
	in.sendPacket(octets)
	return nil
}

func (in *Commander) PublishUsingDefinition(channel channel.ID, cmdName string, data string) error {
	cmd := in.protocolDefinition.FindCommandUsingName(cmdName)
	log.Printf("Sending '%s' to %s", data, cmd)

	octets, toOctetsErr := serializers.StringToOctets(in.protocolDefinition, *cmd, data)
	if toOctetsErr != nil {
		return toOctetsErr
	}
	return in.PublishStream(channel, octets)
}

func (in *Commander) PublishUsingValueFile(channel channel.ID, path string) error {
	fileContents, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		return readErr
	}
	octets, toOctetsErr := serializers.ValueStringToOctets(in.protocolDefinition, string(fileContents))
	if toOctetsErr != nil {
		return toOctetsErr
	}
	return in.PublishStream(channel, octets)
}

func (in *Commander) sendPacket(octets []byte) {
	log.Printf("Sending %X", octets)
	octetCount := len(octets)
	lengthOctets, lengthErr := serializer.SmallLengthToOctets(uint16(octetCount))
	if lengthErr != nil {
		log.Fatalf("Couldn't write length")
		return
	}
	in.conn.Write(lengthOctets)
	in.conn.Write(octets)
}
