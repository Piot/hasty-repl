package repl

import (
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/piot/hasty-protocol/channel"
	"github.com/piot/hasty-repl/commander"
)

type Evaluator struct {
	commander *commander.Commander
}

func NewEvaluator(commander *commander.Commander) Evaluator {
	return Evaluator{commander: commander}
}

func parseChannel(channelString string) channel.ID {
	v, _ := strconv.ParseUint(channelString, 16, 32)
	channel, _ := channel.NewFromID(uint32(v))
	return channel
}

func parseOffset(channelString string) uint32 {
	v, _ := strconv.ParseUint(channelString, 16, 32)
	return uint32(v)
}

func (in Evaluator) parseSubscribe(args []string) {
	log.Printf("Parsing subscribe")
	channel := parseChannel(args[0])
	offset := parseOffset(args[1])
	in.commander.SubscribeStream(channel, offset)
}

func (in Evaluator) parseCreateStream(args []string) {
	log.Printf("Parsing subscribe")
	path := args[0]
	in.commander.CreateStream(path)
}

func (in Evaluator) parseLoadDefinition(args []string) error {
	path := args[0]
	return in.commander.LoadDefinition(path)
}

func (in Evaluator) parsePublishStream(args []string) error {
	channel := parseChannel(args[0])
	hexString := args[1]
	log.Printf("Publish to stream '%s' data:'%s'", channel, hexString)
	decodedChunk, err := hex.DecodeString(hexString)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return in.commander.PublishStream(channel, decodedChunk)
}

func (in Evaluator) parsePublishUsingDefinition(args []string) error {
	channel := parseChannel(args[0])
	cmd := args[1]
	data := strings.Replace(args[2], ",", "\n", -1)
	data = strings.Replace(data, ";", ": ", -1)
	log.Printf("Converted: '%s'", data)
	return in.commander.PublishUsingDefinition(channel, cmd, data)
}

func (in Evaluator) parsePublishUsingValueFile(args []string) error {
	channel := parseChannel(args[0])
	filePath := args[1]
	return in.commander.PublishUsingValueFile(channel, filePath)
}

func (in Evaluator) Eval(input string) error {
	fmt.Printf("You wanted to say '%s'\n", input)
	input = strings.TrimSpace(input)
	commands := strings.Split(input, " ")
	if len(commands) == 0 {
		return nil
	}

	cmd := commands[0]
	args := commands[1:]

	switch cmd {
	case "subs":
		in.parseSubscribe(args)
	case "crs":
		in.parseCreateStream(args)
	case "lddef":
		return in.parseLoadDefinition(args)
	case "pubs":
		return in.parsePublishStream(args)
	case "pub":
		return in.parsePublishUsingDefinition(args)
	case "pubf":
		return in.parsePublishUsingValueFile(args)

	}
	return nil
}
