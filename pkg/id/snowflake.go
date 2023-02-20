// implement https://en.wikipedia.org/wiki/Snowflake_ID
//
// | 41 timestamp | 10 instance | 12 sequence |

package id

import (
	"encoding/base64"
	"encoding/binary"
	"strconv"
	"sync"
	"time"
)

const (
	// 2023-02-17 15:47:00.347516086 +0800 CST
	epoch int64 = 1676620020347

	nodeBits    = 10
	sequenceBit = 12
)

// Instance is a snowflake instance, 10bit, 1024 max
type Instance struct {
	mu        sync.Mutex
	timestamp int64
	instance  int
	sequence  int
}

type ID int64

func NewInstance(instanceId int) *Instance {
	return &Instance{
		instance: instanceId,
	}
}

func (i *Instance) Next() ID {
	i.mu.Lock()
	defer i.mu.Unlock()

	nowTimestamp := time.Now().UnixMilli() - epoch

	// we step forward timestamp when sequance conflict,
	// so we need to make sure timestamp is always forward
	if nowTimestamp < i.timestamp {
		nowTimestamp = i.timestamp
	}

	sequence := i.sequence + 1
	if nowTimestamp == i.timestamp {
		if sequence >= (1 << sequenceBit) {
			// if sequence conflict, step forward timestamp
			nowTimestamp++
			sequence = 0
		}
	} else {
		sequence = 0
	}

	i.sequence = sequence
	i.timestamp = nowTimestamp

	return ID((nowTimestamp << (nodeBits + sequenceBit)) |
		(int64(i.instance) << sequenceBit) |
		int64(sequence))
}

func (i ID) Int64() int64 {
	return int64(i)
}

func (i ID) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i ID) Bytes() []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(i))
	return b[:]
}

func (i ID) Base64() string {
	return base64.StdEncoding.EncodeToString(i.Bytes())
}

func (i ID) Base64Url() string {
	return base64.URLEncoding.EncodeToString(i.Bytes())
}

func Parse(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	return ID(i), err
}

func ParseBase64(s string) (ID, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return 0, err
	}
	return ID(binary.BigEndian.Uint64(b)), nil
}

func ParseBase64Url(s string) (ID, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return 0, err
	}
	return ID(binary.BigEndian.Uint64(b)), nil
}
