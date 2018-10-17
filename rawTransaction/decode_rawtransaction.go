package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type Transcation struct {
	Version   uint32      `json:"version"`
	TxInputs  []*TxInput  `json:"vin"`
	TxOutputs []*TxOutput `json:"vout"`
	LockTime  uint32      `json:"locktime"`
}
type TxInput struct {
	Txid      string `json:"txid"`
	Vout      uint32 `json:"vout"`
	Sequence  uint32 `json:"sequence"`
	ScriptSig []byte `json:"scriptSig"`
}
type TxOutput struct {
	Value    int64  `json:"value"` // satoshis
	PKScript []byte `json:"pk_script"`
}

func ReadCompactSizeUint(buf *bytes.Buffer) uint64 {
	prefix := buf.Next(1)
	switch prefix[0] {
	case 0xfd:
		return uint64(binary.LittleEndian.Uint16(buf.Next(2)))
	case 0xfe:
		return uint64(binary.LittleEndian.Uint32(buf.Next(4)))
	case 0xff:
		return uint64(binary.LittleEndian.Uint64(buf.Next(8)))
	default:
		return uint64(prefix[0])
	}

}

func DecodeRawTransaction(encodedTx []byte) ([]byte, error) {
	var tx Transcation
	byteTx := make([]byte, hex.DecodedLen(len(encodedTx)))
	_, err := hex.Decode(byteTx, encodedTx)
	if err != nil {
		return nil, err
	}
	txBuffer := bytes.NewBuffer(byteTx)

	tx.Version = binary.LittleEndian.Uint32(txBuffer.Next(4))

	txInCount := ReadCompactSizeUint(txBuffer)
	for i := uint64(0); i < txInCount; i++ {
		var txInput TxInput
		txInput.Txid = hex.EncodeToString(txBuffer.Next(32))
		txInput.Vout = binary.LittleEndian.Uint32(txBuffer.Next(4))

		scriptSigLen := ReadCompactSizeUint(txBuffer)
		scriptSig := txBuffer.Next(int(scriptSigLen))
		txInput.ScriptSig = scriptSig

		txInput.Sequence = binary.LittleEndian.Uint32(txBuffer.Next(4))
		tx.TxInputs = append(tx.TxInputs, &txInput)
	}

	txOutCount := ReadCompactSizeUint(txBuffer)
	for i := uint64(0); i < txOutCount; i++ {
		var txOutput TxOutput
		txOutput.Value = int64(
			binary.LittleEndian.Uint64(
				txBuffer.Next(8)))

		pkScriptLen := ReadCompactSizeUint(txBuffer)
		pkScript := txBuffer.Next(int(pkScriptLen))
		txOutput.PKScript = pkScript
		tx.TxOutputs = append(tx.TxOutputs, &txOutput)
	}

	tx.LockTime = binary.LittleEndian.Uint32(txBuffer.Next(4))
	return json.Marshal(&tx)

}

func main() {
	rawtransaction := []byte("0100000001b785e13fb170897ebadd7e563a6fa3a56ad5923e37c471fcecaa0495aeaae79a000000006a47304402202171806c6433ed370b496f4c36b4d552aa690c0f17301016f8acae5dcfba73e602201fcf31a91ab5527881c3aca8227e145ed58a7ab870565d9d7bd40283feb91732012103644b216751be9c1c574f5cc2296a7ddd15663cc2be9a62bbdae726bc6c0fa9ebfdffffff020a060200000000001976a9147935433bf83ea067262da3cdf3f73f3e3b3d6fa788acc7830200000000001976a914570edb3f3bc251f6f05eb9df367f32a46206060888ac474a0800")
	tx, err := DecodeRawTransaction(rawtransaction)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(tx))
}
