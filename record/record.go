package record

import (
	"encoding/json"
)

type AvailableRecord struct {
	records []Record
}

type Record struct {
	Address       string
	CountTrying   uint64
	CountSucceed  uint64
	LastAvailable bool
}

type jsonRecord struct {
	Records []Record
}

func (arcd *AvailableRecord) Record(addr string) (rcd Record) {
	rcd = Record{Address: addr, LastAvailable: true}
	for _, r := range arcd.records {
		if r.Address == addr {
			rcd = r
			break
		}
	}
	return
}

func (arcd *AvailableRecord) Write(addr string, nowAvailable bool) (switched bool) {
	rcd := arcd.Record(addr)

	// right expr is same as `now XOR last`
	switched = nowAvailable != rcd.LastAvailable

	rcd.LastAvailable = nowAvailable
	rcd.CountTrying++
	if nowAvailable {
		rcd.CountSucceed++
	}

	arcd.put(rcd)
	return
}

func (arcd *AvailableRecord) put(rcd Record) {
	for i, r := range arcd.records {
		if r.Address == rcd.Address {
			arcd.records[i] = rcd
			return
		}
	}

	arcd.records = append(arcd.records, rcd)
}

func (arcd AvailableRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonRecord{arcd.records})
}

func (arcd *AvailableRecord) UnmarshalJSON(data []byte) error {
	var jsonRcd jsonRecord
	err := json.Unmarshal(data, &jsonRcd)
	if err != nil {
		return err
	}

	arcd.records = jsonRcd.Records
	return nil
}
