package main

import (
	"context"
	"encoding/hex"
	"github.com/INFURA/go-ethlibs/eth"
	"github.com/INFURA/go-ethlibs/jsonrpc"
	"github.com/INFURA/go-ethlibs/node/websocket"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/op/go-logging"
	"golang.org/x/crypto/sha3"
	"strings"
)

const (
	URL    = "https://rinkeby.infura.io/v3/6141be73a15d47748af0dc14f53d57d7"
	URL_WS = "wss://rinkeby.infura.io/ws/v3/6141be73a15d47748af0dc14f53d57d7"

	LockedABI = "[\n\t\t\t{\n\t\t\t\t\"inputs\": [\n\t\t\t\t\t{\n\t\t\t\t\t\t\"internalType\": \"address\",\n\t\t\t\t\t\t\"name\": \"_relay\",\n\t\t\t\t\t\t\"type\": \"address\"\n\t\t\t\t\t},\n\t\t\t\t\t{\n\t\t\t\t\t\t\"internalType\": \"address[]\",\n\t\t\t\t\t\t\"name\": \"_nfts\",\n\t\t\t\t\t\t\"type\": \"address[]\"\n\t\t\t\t\t},\n\t\t\t\t\t{\n\t\t\t\t\t\t\"internalType\": \"uint256[]\",\n\t\t\t\t\t\t\"name\": \"_networks\",\n\t\t\t\t\t\t\"type\": \"uint256[]\"\n\t\t\t\t\t}\n\t\t\t\t],\n\t\t\t\t\"stateMutability\": \"nonpayable\",\n\t\t\t\t\"type\": \"constructor\"\n\t\t\t},\n\t\t\t{\n\t\t\t\t\"anonymous\": false,\n\t\t\t\t\"inputs\": [\n\t\t\t\t\t{\n\t\t\t\t\t\t\"indexed\": false,\n\t\t\t\t\t\t\"internalType\": \"address\",\n\t\t\t\t\t\t\"name\": \"_nft\",\n\t\t\t\t\t\t\"type\": \"address\"\n\t\t\t\t\t},\n\t\t\t\t\t{\n\t\t\t\t\t\t\"indexed\": false,\n\t\t\t\t\t\t\"internalType\": \"uint256\",\n\t\t\t\t\t\t\"name\": \"_id\",\n\t\t\t\t\t\t\"type\": \"uint256\"\n\t\t\t\t\t},\n\t\t\t\t\t{\n\t\t\t\t\t\t\"indexed\": false,\n\t\t\t\t\t\t\"internalType\": \"uint256\",\n\t\t\t\t\t\t\"name\": \"_networkId\",\n\t\t\t\t\t\t\"type\": \"uint256\"\n\t\t\t\t\t},\n\t\t\t\t\t{\n\t\t\t\t\t\t\"indexed\": false,\n\t\t\t\t\t\t\"internalType\": \"address\",\n\t\t\t\t\t\t\"name\": \"_oldOwner\",\n\t\t\t\t\t\t\"type\": \"address\"\n\t\t\t\t\t}\n\t\t\t\t],\n\t\t\t\t\"name\": \"Locked\",\n\t\t\t\t\"type\": \"event\"\n\t\t\t},\n\t\t\t{\n\t\t\t\t\"anonymous\": false,\n\t\t\t\t\"inputs\": [\n\t\t\t\t\t{\n\t\t\t\t\t\t\"indexed\": false,\n\t\t\t\t\t\t\"internalType\": \"address\",\n\t\t\t\t\t\t\"name\": \"_nft\",\n\t\t\t\t\t\t\"type\": \"address\"\n\t\t\t\t\t},\n\t\t\t\t\t{\n\t\t\t\t\t\t\"indexed\": false,\n\t\t\t\t\t\t\"internalType\": \"uint256\",\n\t\t\t\t\t\t\"name\": \"_id\",\n\t\t\t\t\t\t\"type\": \"uint256\"\n\t\t\t\t\t},\n\t\t\t\t\t{\n\t\t\t\t\t\t\"indexed\": false,\n\t\t\t\t\t\t\"internalType\": \"address\",\n\t\t\t\t\t\t\"name\": \"_newOwner\",\n\t\t\t\t\t\t\"type\": \"address\"\n\t\t\t\t\t}\n\t\t\t\t],\n\t\t\t\t\"name\": \"Unlocked\",\n\t\t\t\t\"type\": \"event\"\n\t\t\t},\n\t\t\t{\n\t\t\t\t\"inputs\": [\n\t\t\t\t\t{\n\t\t\t\t\t\t\"internalType\": \"address\",\n\t\t\t\t\t\t\"name\": \"_nft\",\n\t\t\t\t\t\t\"type\": \"address\"\n\t\t\t\t\t},\n\t\t\t\t\t{\n\t\t\t\t\t\t\"internalType\": \"uint256\",\n\t\t\t\t\t\t\"name\": \"_id\",\n\t\t\t\t\t\t\"type\": \"uint256\"\n\t\t\t\t\t},\n\t\t\t\t\t{\n\t\t\t\t\t\t\"internalType\": \"uint256\",\n\t\t\t\t\t\t\"name\": \"_networkId\",\n\t\t\t\t\t\t\"type\": \"uint256\"\n\t\t\t\t\t}\n\t\t\t\t],\n\t\t\t\t\"name\": \"lock\",\n\t\t\t\t\"outputs\": [],\n\t\t\t\t\"stateMutability\": \"nonpayable\",\n\t\t\t\t\"type\": \"function\"\n\t\t\t},\n\t\t\t{\n\t\t\t\t\"inputs\": [\n\t\t\t\t\t{\n\t\t\t\t\t\t\"internalType\": \"address\",\n\t\t\t\t\t\t\"name\": \"operator\",\n\t\t\t\t\t\t\"type\": \"address\"\n\t\t\t\t\t},\n\t\t\t\t\t{\n\t\t\t\t\t\t\"internalType\": \"address\",\n\t\t\t\t\t\t\"name\": \"from\",\n\t\t\t\t\t\t\"type\": \"address\"\n\t\t\t\t\t},\n\t\t\t\t\t{\n\t\t\t\t\t\t\"internalType\": \"uint256\",\n\t\t\t\t\t\t\"name\": \"tokenId\",\n\t\t\t\t\t\t\"type\": \"uint256\"\n\t\t\t\t\t},\n\t\t\t\t\t{\n\t\t\t\t\t\t\"internalType\": \"bytes\",\n\t\t\t\t\t\t\"name\": \"data\",\n\t\t\t\t\t\t\"type\": \"bytes\"\n\t\t\t\t\t}\n\t\t\t\t],\n\t\t\t\t\"name\": \"onERC721Received\",\n\t\t\t\t\"outputs\": [\n\t\t\t\t\t{\n\t\t\t\t\t\t\"internalType\": \"bytes4\",\n\t\t\t\t\t\t\"name\": \"\",\n\t\t\t\t\t\t\"type\": \"bytes4\"\n\t\t\t\t\t}\n\t\t\t\t],\n\t\t\t\t\"stateMutability\": \"nonpayable\",\n\t\t\t\t\"type\": \"function\"\n\t\t\t},\n\t\t\t{\n\t\t\t\t\"inputs\": [],\n\t\t\t\t\"name\": \"relay\",\n\t\t\t\t\"outputs\": [\n\t\t\t\t\t{\n\t\t\t\t\t\t\"internalType\": \"address payable\",\n\t\t\t\t\t\t\"name\": \"\",\n\t\t\t\t\t\t\"type\": \"address\"\n\t\t\t\t\t}\n\t\t\t\t],\n\t\t\t\t\"stateMutability\": \"view\",\n\t\t\t\t\"type\": \"function\"\n\t\t\t},\n\t\t\t{\n\t\t\t\t\"inputs\": [\n\t\t\t\t\t{\n\t\t\t\t\t\t\"internalType\": \"address\",\n\t\t\t\t\t\t\"name\": \"_nft\",\n\t\t\t\t\t\t\"type\": \"address\"\n\t\t\t\t\t},\n\t\t\t\t\t{\n\t\t\t\t\t\t\"internalType\": \"uint256\",\n\t\t\t\t\t\t\"name\": \"_id\",\n\t\t\t\t\t\t\"type\": \"uint256\"\n\t\t\t\t\t},\n\t\t\t\t\t{\n\t\t\t\t\t\t\"internalType\": \"address\",\n\t\t\t\t\t\t\"name\": \"_newOwner\",\n\t\t\t\t\t\t\"type\": \"address\"\n\t\t\t\t\t}\n\t\t\t\t],\n\t\t\t\t\"name\": \"unlock\",\n\t\t\t\t\"outputs\": [],\n\t\t\t\t\"stateMutability\": \"nonpayable\",\n\t\t\t\t\"type\": \"function\"\n\t\t\t}\n\t\t]"
)

type NewLogsNotificationParams struct {
	Subscription string  `json:"subscription"`
	Result       eth.Log `json:"result"`
}

var log = logging.MustGetLogger("bus")

func main() {
	ctx := context.Background()

	client, err := websocket.NewConnection(ctx, URL_WS)
	if err != nil {
		log.Fatalf("Could not connect to %s: %v", URL_WS, err)
	}

	log.Infof("Connected to %s", client.URL())

	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		log.Fatalf("Could not get block number: %v", err)
	}

	log.Infof("Current block number: %v", blockNumber)

	//filter := &eth.LogFilter{
	//	//ToBlock:   eth.MustBlockNumberOrTag("0x" + strconv.FormatUint(blockNumber-20, 16)),
	//	ToBlock:   eth.MustBlockNumberOrTag(eth.TagLatest.String()),
	//	FromBlock: eth.MustBlockNumberOrTag(eth.TagLatest.String()),
	//	//Address:   []eth.Address{"0xAf95D25652077f5E9412712C57708745E9Fcc9E7"},
	//}

	topic0 := Encode(Keccak256([]byte("Locked(address,uint256,uint256,address)")))
	log.Debugf(topic0)

	filter := &eth.LogFilter{
		ToBlock:   eth.MustBlockNumberOrTag(eth.TagLatest.String()),
		FromBlock: eth.MustBlockNumberOrTag(eth.TagEarliest.String()),
		//BlockHash: eth.MustHash("0xb2e011f4422a743cad50aa1b862568071ca6f7ff3a73541b44a28d86015f153e"),
		Address: []eth.Address{*eth.MustAddress("0x8BB0aE667F011604c9F152879A0aD6Eb12882cC6")},
		//Topics: [][]eth.Topic{
		//	{
		//		eth.Topic(topic0),
		//	},
		//},
	}
	//l, _ := client.GetLogs(ctx, filter)

	r := &jsonrpc.Request{
		JSONRPC: "2.0",
		ID: jsonrpc.ID{
			Num: 3,
		},
		Method: "eth_subscribe",
		Params: jsonrpc.MustParams("logs", filter),
	}

	heads, err := client.Subscribe(ctx, r)
	if err != nil {
		log.Fatalf("[FATAL] Logs subscription error: %v", err)
	}

	for notif := range heads.Ch() {

		l := &NewLogsNotificationParams{}
		err := notif.UnmarshalParamsInto(l)
		if err != nil {
			log.Fatalf("[FATAL] Cannot parse newHeads params: %v", err)
		}

		data := string(l.Result.Data)
		m := decode(data)

		spew.Dump(m)
	}
}

func decode(data string) (m map[string]interface{}) {
	// load contract ABI
	m = make(map[string]interface{})
	locked, err := abi.JSON(strings.NewReader(LockedABI))
	if err != nil {
		log.Fatal(err)
	}

	if err := locked.UnpackIntoMap(m, "Locked", Hex2Bytes(data[2:])); err != nil {
		return nil
	}

	return m
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

// Encode encodes b as a hex string with 0x prefix.
func Encode(b []byte) string {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return string(enc)
}

// FromHex returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func FromHex(s string) []byte {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}
