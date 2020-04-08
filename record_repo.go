package main

import (
	"encoding/json"
	"io/ioutil"
)

const repoFileName = "record.json"

type RecordRepository struct {
	Records []PingRecord
}

type PingRecord struct {
	Address      string
	CountTrying  uint64
	CountSucceed uint64
	LastAchieved bool
}

func LoadRecordRepository() (*RecordRepository, error) {
	repo := RecordRepository{}
	path := GetCurrentPath() + "/" + repoFileName

	if FileExists(path) {
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(buf, &repo)
		if err != nil {
			return nil, err
		}
	}

	return &repo, nil
}

func (repo *RecordRepository) Flush() error {
	bytes, err := json.Marshal(repo)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(GetCurrentPath()+"/"+repoFileName, bytes, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (repo *RecordRepository) Record(addr string, achieved bool) (switched bool) {
	rcd := repo.GetRecord(addr)

	switched = achieved != rcd.LastAchieved // XOR

	rcd.LastAchieved = achieved
	rcd.CountTrying++
	if achieved {
		rcd.CountSucceed++
	}

	repo.putRecord(rcd)
	return
}

func (repo *RecordRepository) GetRecord(addr string) PingRecord {
	for _, repoRcd := range repo.Records {
		if repoRcd.Address == addr {
			return repoRcd
		}
	}

	return PingRecord{Address: addr, LastAchieved: true}
}

func (repo *RecordRepository) putRecord(rcd PingRecord) {
	for i, repoRcd := range repo.Records {
		if repoRcd.Address == rcd.Address {
			repo.Records[i] = rcd
			return
		}
	}

	repo.Records = append(repo.Records, rcd)
}

type jsonRecord struct {
	Records []PingRecord
}

func (repo *RecordRepository) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonRecord{repo.Records})
}

func (repo *RecordRepository) UnmarshalJSON(data []byte) error {
	var jsonRcd jsonRecord
	err := json.Unmarshal(data, &jsonRcd)
	if err != nil {
		return err
	}

	repo.Records = jsonRcd.Records
	return nil
}
