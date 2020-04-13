package repo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

var recordRepo *RecordRepository

const repoFileName = "repo.json"

type RecordRepository struct {
	Records []PingRecord
}

type PingRecord struct {
	Address      string
	CountTrying  uint64
	CountSucceed uint64
	LastAchieved bool
	FirstBootAt  time.Time
	LastBootAt   time.Time
}

func init() {
	recordRepo = &RecordRepository{}
	path := getCurrentPath() + "/" + repoFileName

	if fileExists(path) {
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal("Failed to read a config")
		}
		err = json.Unmarshal(buf, &recordRepo)
		if err != nil {
			log.Fatal("Failed to unmarshal json in a config")
		}
	}
}

func GetRecordRepository() *RecordRepository {
	return recordRepo
}

func (repo *RecordRepository) Flush() error {
	bytes, err := json.Marshal(repo)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(getCurrentPath()+"/"+repoFileName, bytes, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (repo *RecordRepository) Record(addr string, achieved bool) (switched bool) {
	rcd := repo.GetOrNewRecord(addr)

	switched = achieved != rcd.LastAchieved // XOR

	rcd.LastAchieved = achieved
	rcd.CountTrying++
	if achieved {
		rcd.CountSucceed++
	}

	// down -> up
	if switched && achieved {
		rcd.LastBootAt = time.Now()
	}

	repo.putRecord(rcd)
	return
}

func (repo *RecordRepository) GetOrNewRecord(addr string) PingRecord {
	for _, repoRcd := range repo.Records {
		if repoRcd.Address == addr {
			return repoRcd
		}
	}

	return PingRecord{
		Address:      addr,
		LastAchieved: true,
		FirstBootAt:  time.Now(),
		LastBootAt:   time.Now(),
	}
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

func getCurrentPath() string {
	exec, err := os.Executable()
	if err != nil {
		log.Fatal("Failed to get an executable")
	}

	return filepath.Dir(exec)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
