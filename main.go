package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const RandomLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

type Account struct {
	Id        int
	AccountId string
	LoginId   string
	Password  string
	Retailer  string
}

func main() {
	if len(os.Args) != 2 {
		failOnError("main.exeにCSVファイルをドラッグ&ドロップしてください", nil)
	}
	//convertCsv2Json("./testdata/rakuten_omni_sample.csv")
	convertCsv2Json(os.Args[1])
	waitEnter()
}

func failOnError(errMsg string, err error) {
	//errs := errors.WithStack(err)
	fmt.Println(errMsg)
	if err != nil {
		//fmt.Printf("%+v\n", errs) Stack trace
		fmt.Printf("%s\n", err.Error())
	}
	waitEnter()
	os.Exit(1)
}

func waitEnter() {
	fmt.Println("エンターを押すと処理を終了します。")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
}

func getCsvRecords(csvFilePath string) [][]string {
	fp, err := os.Open(csvFilePath)
	if err != nil {
		failOnError("ファイル読込に失敗しました", err)
	}
	defer fp.Close()

	inputCsv := csv.NewReader(transform.NewReader(fp, japanese.ShiftJIS.NewDecoder()))

	records, err := inputCsv.ReadAll()
	if err != nil {
		failOnError("CSV読み込みエラー", err)
	}

	return records
}

func convertCsv2Json(csvFilePath string) {
	records := getCsvRecords(csvFilePath)
	accountList := make([]*Account, 0)
	id := 0
	for i, items := range records {
		_, err := strconv.Atoi(items[0])
		if err != nil {
			if i != 0 {
				fmt.Printf("%d行目の1列目は管理ID（数値）ではないため処理をJSON出力対象外とします。実際の値->'%s'\n",
					i+1, items[0])
			}
			continue
		}

		retailer := ""
		switch items[1] {
		case "楽天ブックス":
			retailer = "Rakuten"
		case "omni7(セブンネットなど）":
			retailer = "Omni7"
		default:
			fmt.Printf("規定外サイトのため、JSON出力対象外とします。%d行目 値->'%s'\n", i, items[1])
			continue
		}

		account := &Account{Id: id, AccountId: getRandomString(20), LoginId: items[2], Password: items[3], Retailer: retailer}
		accountList = append(accountList, account)
		id++
	}

	outputJson, err := json.MarshalIndent(accountList, "", "\t")
	if err != nil {
		failOnError("JSON変換エラー", err)
	}

	dumpJson(outputJson)
}

func dumpJson(jsonData []byte) {
	exe, err := os.Executable()
	if err != nil {
		failOnError("exeファイル実行パス取得失敗", err)
	}

	outputDirPath := filepath.Dir(exe)
	file, err := os.OpenFile(outputDirPath+"/Account.json", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		failOnError("Account.jsonのオープンに失敗しました", err)
	}
	defer file.Close()

	err = file.Truncate(0) // ファイルを空っぽにする(実行2回目以降用)
	if err != nil {
		failOnError("JSONファイルの初期化に失敗しました", err)
	}

	writer := transform.NewWriter(file, japanese.ShiftJIS.NewEncoder())
	writer.Write(jsonData)

	fmt.Println(outputDirPath + "\\Account.jsonを出力しました")

}

func getRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = RandomLetters[rand.Intn(len(RandomLetters))]
	}
	return string(b)
}
