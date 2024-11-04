package generator

import (
	"fmt"
	"sync"
	"time"
)

type Snowflake struct {
	mu							sync.Mutex
	epoch						int64
	dataCenterId		int64
	machineId				int64
	sequence				int64
	lastTimestamp		int64
}

const (
	dataCenterIdBits = 5
	machineIdBits	= 5
	sequenceBits = 12

	maxDataCenterId = -1 ^ (-1 << dataCenterIdBits)
	maxMachineId = -1 ^ (-1 << machineIdBits)
	maxSequence = -1 ^ (-1 << sequenceBits)

	machineIdShift = sequenceBits
	dataCenterIdShift = sequenceBits + machineIdBits
	timestampShift = sequenceBits + machineIdBits + dataCenterIdBits
)

func NewSnowflake(dataCenterId, machineId int64) (*Snowflake, error) {
	if dataCenterId < 0 || dataCenterId > maxDataCenterId {
		return nil, fmt.Errorf("data center id must be between 0 and %d", maxDataCenterId)
	}

	if machineId < 0 || maxMachineId > maxMachineId {
		return nil, fmt.Errorf("machine id must be between 0 and %d", maxMachineId)
	}

	return &Snowflake{
		epoch: 1609459200000,
		dataCenterId: dataCenterId,
		machineId: machineId,
		sequence: 0,
		lastTimestamp: -1,
	}, nil
}

func (s *Snowflake) NextId() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := s.currentTimestamp()

	if timestamp < s.lastTimestamp {
		return 0, fmt.Errorf("clock moved backwards, error generating id")
	}

	if timestamp == s.lastTimestamp {
		// increment the sequence within the same millisecond
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			// sequence overflow; wait for next millisecond
			timestamp = s.waitNextMillis(s.lastTimestamp)
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	id := ((timestamp - s.epoch) << timestampShift) |
				(s.dataCenterId << dataCenterIdShift) |
				(s.machineId << machineIdShift) |
				s.sequence

	return id, nil
}

func (s *Snowflake) currentTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}

func (s *Snowflake) waitNextMillis(lastTimestamp int64) int64 {
	timestamp := s.currentTimestamp()
	for timestamp <= lastTimestamp {
		timestamp = s.currentTimestamp()
	}
	return timestamp
}
