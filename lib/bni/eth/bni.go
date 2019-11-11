package bni

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	base_raw_transaction "github.com/HyperService-Consortium/go-uip/base-raw-transaction"
	payment_option "github.com/HyperService-Consortium/go-ves/lib/bni/payment-option"
	"github.com/HyperService-Consortium/go-ves/types"
	"math/big"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/HyperService-Consortium/go-ethabi"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/sha3"

	TransType "github.com/HyperService-Consortium/go-uip/const/trans_type"
	valuetype "github.com/HyperService-Consortium/go-uip/const/value_type"
	"github.com/HyperService-Consortium/go-uip/uiptypes"

	ethclient "github.com/HyperService-Consortium/go-ves/lib/net/eth-client"
)

func decoratePrefix(hexs string) string {
	if !strings.HasPrefix(hexs, "0x") {
		hexs =  "0x" + hexs
	}
	for strings.HasPrefix(hexs, "0x0") && len(hexs) > 3 {
		hexs = "0x" + hexs[3:]
	}
	return hexs
}

type BN struct {
	dns types.ChainDNSInterface
	signer uiptypes.Signer
}

func (bn *BN) Deserialize(raw []byte) (rawTransaction uiptypes.RawTransaction, err error) {
	// todo
	return base_raw_transaction.Transaction(raw), nil
}

func (bn *BN) Translate(intent *uiptypes.TransactionIntent, storage uiptypes.Storage) (uiptypes.RawTransaction, error) {
	switch intent.TransType {
	case TransType.Payment:
		meta := gjson.ParseBytes(intent.Meta)
		value, err := payment_option.ParseInconsistentValueOption(meta, storage, intent.Amt)
		if err != nil {
			return nil, err
		}

		//fmt.Println(value, ".........")

		b, err := json.Marshal(map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "eth_sendTransaction",
			"params": []interface{}{
				map[string]interface{}{
					"from":  decoratePrefix(hex.EncodeToString(intent.Src)),
					"to":    decoratePrefix(hex.EncodeToString(intent.Dst)),
					"value": decoratePrefix(value),
				},
			},
			"id": 1,
		})
		//fmt.Println("...", string(b))
		return base_raw_transaction.Transaction(b), err
	case TransType.ContractInvoke:
		var meta uiptypes.ContractInvokeMeta
		err := json.Unmarshal(intent.Meta, &meta)
		if err != nil {
			return nil, err
		}
		//_ = meta
		// todo, test
		data, err := ContractInvocationDataABI(intent.ChainID, &meta, storage)
		if err != nil {
			return nil, err
		}

		hexdata := hex.EncodeToString(data)
		// meta.FuncName

		b, err := json.Marshal(map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "eth_sendTransaction",
			"params": []interface{}{
				map[string]interface{}{
					"from":  decoratePrefix(hex.EncodeToString(intent.Src)),
					"to":    decoratePrefix(hex.EncodeToString(intent.Dst)),
					// todo
					//"value": decoratePrefix(intent.Amt),
					"data":  decoratePrefix(hexdata),
				},
			},
			"id": 1,
		})
		return base_raw_transaction.Transaction(b), err
	default:
		return nil, errors.New("cant translate")
	}
}


var ErrNotSigned  = errors.New("not signed raw transaction")

func (bn *BN) RouteRaw(destination uiptypes.ChainID, rawTransaction uiptypes.RawTransaction) (
	transactionReceipt uiptypes.TransactionReceipt, err error) {
	if !rawTransaction.Signed() {
		return nil, ErrNotSigned
	}
	ci, err := bn.dns.GetChainInfo(destination)
	if err != nil {
		return nil, err
	}
	// todo receipt
	b, err := rawTransaction.Serialize()
	if err != nil {
		return nil, err
	}
	b, err = ethclient.Do((&url.URL{Scheme: "http", Host: ci.GetChainHost(), Path: "/"}).String(), b)
	if err != nil {
		return nil, err
	}
	var x string
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}

	if x == "" {
		return nil, errors.New("deploy failed?")
	}

	b, err = hex.DecodeString(x[2:])
	if err != nil {
		return nil, err
	}
	return b, nil
}

type options struct {
	timeout time.Duration
}

func parseOptions(rOption []interface{}) options {
	var parsedOptions options
	for i := range rOption {
		switch option := rOption[i].(type) {
		case time.Duration:
			parsedOptions.timeout = option
		case uiptypes.RouteOptionTimeout:
			parsedOptions.timeout = time.Duration(option)
		}
	}
	return parsedOptions
}

func (bn *BN) WaitForTransact(chainID uiptypes.ChainID, transactionReceipt uiptypes.TransactionReceipt,
	rOptions ...interface{}) (blockID uiptypes.BlockID, color []byte, err error) {
	options := parseOptions(rOptions)
	chainInfo, err := bn.dns.GetChainInfo(chainID)
	if err != nil {
		return nil, nil, err
	}
	worker := ethclient.NewEthClient(chainInfo.GetChainHost())
	ddl := time.Now().Add(options.timeout)
	for time.Now().Before(ddl) {
		tx, err := worker.GetTransactionByHash(transactionReceipt)
		if err != nil {
			return nil, nil, err
		}
		fmt.Println(string(tx))
		if gjson.GetBytes(tx, "blockNumber").Type != gjson.Null {
			b, _ := hex.DecodeString(gjson.GetBytes(tx, "blockHash").String()[2:])
			idx, _ := strconv.ParseUint(gjson.GetBytes(tx, "transactionIndex").String()[2:], 16, 64)
			var a = make([]byte, 8)
			binary.BigEndian.PutUint64(a, idx)
			return b, a, nil
		}
		time.Sleep(time.Millisecond * 500)
	}
	return nil, nil, errors.New("timeout")
}

func (bn *BN) RouteWithSigner(signer uiptypes.Signer) (uiptypes.Router, error) {
	nbn :=  *bn
	nbn.signer = signer
	return &nbn, nil
}

func NewBN(dns types.ChainDNSInterface) *BN {
	return &BN{dns: dns}
}

func (bn *BN) MustWithSigner() bool {
	return true
}

func (bn *BN) Route(intent *uiptypes.TransactionIntent, storage uiptypes.Storage) ([]byte, error) {
	// todo
	rawTransaction, err := bn.Translate(intent, storage)
	if err != nil {
		return nil, err
	}
	if !rawTransaction.Signed() {
		rawTransaction, err = rawTransaction.Sign(bn.signer)
		if err != nil {
			return nil, err
		}
	}
	return bn.RouteRaw(intent.ChainID, rawTransaction)
}


func (bn *BN) CheckAddress(addr []byte) bool {
	return len(addr) == 32
}

func convertToEthVariable(variable uiptypes.Variable) interface{} {
	return variable.GetValue()
}

func ContractInvocationDataABI(chainID uiptypes.ChainID, meta *uiptypes.ContractInvokeMeta, storage uiptypes.Storage) ([]byte, error) {

	abiencoder := ethabi.NewEncoder()
	//Encodes(descs []string, vals []interface{})
	keccak := sha3.NewLegacyKeccak256()
	var descs []string
	var vals []interface{} = make([]interface{}, 0)
	var funcsig string = meta.FuncName + "("
	//var err error
	for id, param := range meta.Params {
		t := param.Type
		if t == "address payable" || t == "contract" {
			t = "address"
		}
		funcsig += t
		descs = append(descs, t)
		constant := gjson.Get(string(param.Value), "constant")
		if !constant.Exists() {
			field := gjson.Get(string(param.Value), "field")
			if !field.Exists() {
				return nil, errors.New("field not found")
			}
			pos := gjson.Get(string(param.Value), "pos")
			if !pos.Exists() {
				return nil, errors.New("pos not found")
			}
			contract := gjson.Get(string(param.Value), "contract")
			if !contract.Exists() {
				return nil, errors.New("contract not found")
			}
			var contractAddress, err = hex.DecodeString(contract.String())
			if err != nil {
				return nil, err
			}
			v, err := storage.GetStorageAt(chainID, valuetype.FromString(t), contractAddress, []byte(pos.Str), []byte(field.String()))
			if err != nil {
				return nil, err
			}
			vals = append(vals, convertToEthVariable(v))
		} else {
			val, err := appendVal(t, constant)
			if err != nil {
				return nil, err
			}
			vals = append(vals, val)
		}

		if id == len(meta.Params)-1 {
			funcsig += ")"
		} else {
			funcsig += ","
		}
	}
	keccak.Write([]byte(funcsig))
	result := keccak.Sum([]byte(""))[0:4]
	if len(descs) > 0 {
		result2, err := abiencoder.Encodes(descs, vals)
		if err != nil {
			return nil, err
		}
		result = append(result, result2...)
	}
	return result, nil
}

func appendSliceVal(brcnt int, t string, rawval gjson.Result) (interface{}, error) {
	var err error
	var ret interface{}
	i := strings.LastIndex(t, "[")
	// grab the last cell and create a type from there
	sliced := t[i:]
	// grab the slice size with regexp
	re := regexp.MustCompile("[0-9]+")
	intz := re.FindAllString(sliced, -1)
	var arr []gjson.Result = rawval.Array()
	if len(intz) == 1 {
		// is a array
		siz, err := strconv.Atoi(intz[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing variable size: %v", err)
		}
		if siz != len(arr) {
			return nil, fmt.Errorf("array length not match")
		}
	} else if len(intz) != 0 {
		return nil, fmt.Errorf("invalid formatting of array type")
	}
	/////////////////////////////////////////////////////////
	t = t[:i]
	typeRegex := regexp.MustCompile("([a-zA-Z]+)(([0-9]+)(x([0-9]+))?)?")
	matches := typeRegex.FindAllStringSubmatch(t, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("invalid type '%v'", t)
	}
	parsedType := matches[0]
	var varSize int
	if len(parsedType[3]) > 0 {
		var err error
		varSize, err = strconv.Atoi(parsedType[2]) //ParseInt(sparsedType[2], 10, 0) //strconv.Atoi()
		if err != nil {
			return nil, fmt.Errorf("error parsing variable size: %v", err)
		}
	} else {
		if parsedType[0] == "uint" || parsedType[0] == "int" {
			// this should fail because it means that there's something wrong with
			// the abi type (the compiler should always format it to the size...always)
			return nil, fmt.Errorf("unsupported arg type: %s", t)
		}
	}
	switch varType := parsedType[1]; varType {
	case "int":
		switch varSize {
		case 8:
			ret = make([]int8, len(arr))
			for id, obj := range arr {
				elem, err := appendVal(t, obj)
				ret.([]int8)[id] = elem.(int8)
				if err != nil {
					return nil, err
				}
			}
		case 16:
			ret = make([]int16, len(arr))
			for id, obj := range arr {
				elem, err := appendVal(t, obj)
				ret.([]int16)[id] = elem.(int16)
				if err != nil {
					return nil, err
				}
			}
		case 32:
			ret = make([]int32, len(arr))
			for id, obj := range arr {
				elem, err := appendVal(t, obj)
				ret.([]int32)[id] = elem.(int32)
				if err != nil {
					return nil, err
				}
			}
		case 64:
			ret = make([]int64, len(arr))
			for id, obj := range arr {
				elem, err := appendVal(t, obj)
				ret.([]int64)[id] = elem.(int64)
				if err != nil {
					return nil, err
				}
			}
		case 256:
			ret = make([]big.Int, len(arr))
			for id, obj := range arr {
				elem, err := appendVal(t, obj)
				ret.([]big.Int)[id] = elem.(big.Int)
				if err != nil {
					return nil, err
				}
			}
		}
	case "uint":
		switch varSize {
		case 8:
			ret = make([]uint8, len(arr))
			for id, obj := range arr {
				elem, err := appendVal(t, obj)
				ret.([]uint8)[id] = elem.(uint8)
				if err != nil {
					return nil, err
				}
			}
		case 16:
			ret = make([]uint16, len(arr))
			for id, obj := range arr {
				elem, err := appendVal(t, obj)
				ret.([]uint16)[id] = elem.(uint16)
				if err != nil {
					return nil, err
				}
			}
		case 32:
			ret = make([]uint32, len(arr))
			for id, obj := range arr {
				elem, err := appendVal(t, obj)
				ret.([]uint32)[id] = elem.(uint32)
				if err != nil {
					return nil, err
				}
			}
		case 64:
			ret = make([]uint64, len(arr))
			for id, obj := range arr {
				elem, err := appendVal(t, obj)
				ret.([]uint64)[id] = elem.(uint64)
				if err != nil {
					return nil, err
				}
			}
		case 256:
			ret = make([]big.Int, len(arr))
			for id, obj := range arr {
				elem, err := appendVal(t, obj)
				ret.([]big.Int)[id] = elem.(big.Int)
				if err != nil {
					return nil, err
				}
			}
		}
	case "bool":
		ret = make([]bool, len(arr))
		for id, obj := range arr {
			elem, err := appendVal(t, obj)
			ret.([]bool)[id] = elem.(bool)
			if err != nil {
				return nil, err
			}
		}
	case "address":
		ret = make([][20]byte, len(arr))
		for id, obj := range arr {
			elem, err := appendVal(t, obj)
			ret.([][20]byte)[id] = elem.([20]byte)
			if err != nil {
				return nil, err
			}
		}
	case "string":
		ret = make([]string, len(arr))
		for id, obj := range arr {
			elem, err := appendVal(t, obj)
			ret.([]string)[id] = elem.(string)
			if err != nil {
				return nil, err
			}
		}
	case "bytes":
		ret = make([][]byte, len(arr))
		for id, obj := range arr {
			elem, err := appendVal(t, obj)
			ret.([][]byte)[id] = elem.([]byte)
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("unsupported arg type: %s", t)
	}
	///////////////////////////////////////////
	//fmt.Println("ARR", ret, reflect.TypeOf(ret).Elem())
	return ret, err
}

func appendVal(t string, rawval gjson.Result) (interface{}, error) {
	// check that array brackets are equal if they exist
	brcnt := strings.Count(t, "[")
	if brcnt != strings.Count(t, "]") {
		return nil, fmt.Errorf("invalid arg type in abi")
	}
	var err error
	var ret interface{}

	// if there are brackets, get ready to go into slice/array mode and
	// recursively create the type
	if brcnt != 0 {
		if brcnt != 1 {
			return nil, fmt.Errorf("unsupported arg type: %s", t)
		}
		return appendSliceVal(brcnt, t, rawval)
	}
	typeRegex := regexp.MustCompile("([a-zA-Z]+)(([0-9]+)(x([0-9]+))?)?")
	matches := typeRegex.FindAllStringSubmatch(t, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("invalid type '%v'", t)
	}
	parsedType := matches[0]
	var varSize int
	if len(parsedType[3]) > 0 {
		var err error
		varSize, err = strconv.Atoi(parsedType[2]) //ParseInt(sparsedType[2], 10, 0) //strconv.Atoi()
		if err != nil {
			return nil, fmt.Errorf("error parsing variable size: %v", err)
		}
	} else {
		if parsedType[0] == "uint" || parsedType[0] == "int" {
			// this should fail because it means that there's something wrong with
			// the abi type (the compiler should always format it to the size...always)
			return nil, fmt.Errorf("unsupported arg type: %s", t)
		}
	}
	switch varType := parsedType[1]; varType {
	case "int":
		switch varSize {
		case 8:
			ret = int8(rawval.Int())
		case 16:
			ret = int16(rawval.Int())
		case 32:
			ret = int32(rawval.Int())
		case 64:
			ret = int64(rawval.Int())
		case 256:
			str := rawval.String()
			val, ok := big.NewInt(0).SetString(str, 10)
			if !ok {
				return nil, fmt.Errorf("Invalid int value")
			}
			ret = val
		}
	case "uint":
		switch varSize {
		case 8:
			ret = uint8(rawval.Int())
		case 16:
			ret = uint16(rawval.Int())
		case 32:
			ret = uint32(rawval.Int())
		case 64:
			ret = uint64(rawval.Int())
		case 256:
			str := rawval.String()
			val, ok := big.NewInt(0).SetString(str, 10)
			if !ok {
				return nil, fmt.Errorf("Invalid int value")
			}
			ret = val
		}
	case "bool":
		ret = rawval.Bool()
	case "address":
		addr := rawval.String()
		retsli, err := hex.DecodeString(addr[2:])
		if err != nil {
			return nil, err
		}
		if len(retsli) != 20 {
			return nil, fmt.Errorf("invalid address value")
		}
		var rett [20]byte
		for i := 0; i < 20; i++ {
			rett[i] = retsli[i]
		}
		ret = rett
	case "string":
		ret = rawval.String()
	case "bytes":
		if varSize != 0 {
			return nil, fmt.Errorf("unsupported arg type: %s", t)
		}
		tmpret := []byte("\"" + rawval.String() + "\"")
		var retval []byte
		err = json.Unmarshal(tmpret, &retval)
		if err != nil {
			return nil, err
		}
		ret = retval
	default:
		return nil, fmt.Errorf("unsupported arg type: %s", t)
	}
	return ret, nil
}
