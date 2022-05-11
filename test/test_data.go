package ldtest

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"time"

	_ "github.com/stretchr/testify/assert"
)

type TestData struct {
	testFileName string
	testFile     *os.File
	csvFile      *csv.Reader
	variants     []string

	reportFile    *os.File
	reportDirPath string
}

func NewTestData() (TestData, error) {

	reportDirPath := fmt.Sprintf("./benchmark-report-%v", time.Now().Format(time.RFC3339))
	err := os.Mkdir(reportDirPath, 0750)

	return TestData{
		testFileName:  "./data/insertion.csv",
		reportDirPath: reportDirPath,
	}, err
}

func (t *TestData) Close() {
	if t.testFile != nil {

		t.testFile.Close()
		t.testFile = nil

		t.csvFile = nil
	}
}

func (t *TestData) Reopen() error {
	if t.testFile != nil {
		t.testFile.Close()
	}
	testFile, err := os.Open(t.testFileName)
	if err != nil {
		return err
	}

	t.csvFile = csv.NewReader(testFile)
	recordHeader, err := t.csvFile.Read()
	if err != nil {
		return err
	}
	t.variants = recordHeader[1:]
	return nil
}

func (t *TestData) Choose(prefix string, idx int) error {
	if t.reportFile != nil {
		t.reportFile.Close()
	}

	var err error
	t.reportFile, err = os.Create(path.Join(t.reportDirPath, fmt.Sprintf("insertion-%v-%v.csv", prefix, t.variants[idx])))
	if err != nil {
		return err
	}
	fmt.Fprintf(t.reportFile, "n,tps\n")
	return nil
}

func (t *TestData) GetVariants() ([]string, error) {
	if t.csvFile == nil {
		err := t.Reopen()
		if err != nil {
			return []string{}, err
		}
	}

	return t.variants, nil
}

func (t *TestData) Initialize(ctx context.Context, idx int, insertFunc func(score float64, user string)) (int, error) {
	err := t.Reopen()
	if err != nil {
		return 0, err
	}
	maxUsers := 10_000

	lastNewUser := 0
	lastLog := 0
	secondsStep := time.Duration(10)
	timer := time.After(time.Second * secondsStep)

	for {
		select {
		case <-ctx.Done():
			return lastNewUser, ctx.Err()
		case <-timer:
			tps := float64(lastNewUser-lastLog) / float64(secondsStep)
			fmt.Fprintf(t.reportFile, "%v,%v\n", lastNewUser, tps)
			lastLog = lastNewUser
			timer = time.After(time.Second * secondsStep)
		default:
		}

		record, err := t.csvFile.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, err
		}

		f, err := strconv.ParseFloat(record[idx], 64)
		if err != nil {
			return 0, err
		}

		insertFunc(f, record[0])
		lastNewUser += 1
		if lastNewUser/1000 > lastLog {
			lastLog = lastNewUser / 1000
		}
		_ = maxUsers
		// if lastNewUser >= maxUsers {
		// 	break
		// }

	}

	return lastNewUser, nil
}
