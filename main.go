package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

//Data data
type Data struct {
	CachedSince string `json:"cached_since"`
	Entries     string `json:"entries"`
}

func saveData(body []byte, x int) {
	var tableH string
	var table string
	if x == 0 || x == 3 {
		table = "ggg_ladder"
		tableH = "ggg_ladder_history"
	} else if x == 1 || x == 2 {
		table = "ggg_ssfladder"
		tableH = "ggg_ssfladder_history"
	}

	var league string
	if x == 0 || x == 1 {
		league = "hc"
	} else if x == 2 || x == 3 {
		league = "sc"
	}

	godotenv.Load()
	HOST := os.Getenv("HOST")
	PORT := os.Getenv("PORT")
	DATABASE := os.Getenv("DATABASE")
	USER := os.Getenv("USER")
	PASSWORD := os.Getenv("PASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		HOST, PORT, USER, PASSWORD, DATABASE)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var data Data
	json.Unmarshal(body, &data)
	var dat map[string]interface{}
	json.Unmarshal(body, &dat)
	strs := dat["entries"].([]interface{})
	for i := 0; i < len(strs); i++ {
		// fmt.Println(data.CachedSince)        //<-- SAVE 1
		// fmt.Println(table + tableH + league) //<-- SAVE 2
		temp := strs[i].(map[string]interface{})
		// fmt.Println(temp["rank"])   //<-- SAVE 3
		// fmt.Println(temp["dead"])   //<-- SAVE 4
		// fmt.Println(temp["online"]) //<-- SAVE 5
		character := temp["character"].(map[string]interface{})
		// fmt.Println(character["name"])  //<-- SAVE 6
		// fmt.Println(character["level"]) //<-- SAVE 7
		// fmt.Println(character["class"]) //<-- SAVE 8
		// fmt.Println(character["id"])    //<-- SAVE 9
		exp := int(character["experience"].(float64))
		// fmt.Println(exp) //<-- SAVE 10

		var depth = map[string]int{}
		if character["depth"] != nil {
			depthTemp := character["depth"].(map[string]interface{})
			depth["default"] = int(depthTemp["default"].(float64))
			depth["solo"] = int(depthTemp["solo"].(float64))
			// fmt.Println(depth["default"]) //<-- SAVE11
			// fmt.Println(depth["solo"])    //<-- SAVE12
		} else {
			depth["default"] = int(0)
			depth["solo"] = int(0)
			// fmt.Println(depth["default"]) //<-- SAVE11
			// fmt.Println(depth["solo"])    //<-- SAVE12
		}

		account := temp["account"].(map[string]interface{})
		// fmt.Println(account["name"]) //<-- SAVE 13
		challenges := account["challenges"].(map[string]interface{})
		// fmt.Println(challenges["total"]) //<-- SAVE 14
		tx := db.MustBegin()

		rank := strconv.Itoa(int(temp["rank"].(float64)))
		insertTable := `INSERT INTO ` + table +
			` (cached_since, league, rank, dead, online,
			character_name, character_level,  character_class, character_id, character_experience,
			character_depth_solo, character_depth_default, account_name, account_challenges_total) VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

		insertTableH := `INSERT INTO ` + tableH +
			` (cached_since, league, rank, dead, online,
			character_name, character_level,  character_class, character_id, character_experience,
			character_depth_solo, character_depth_default, account_name, account_challenges_total) VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

		db.MustExec("DELETE FROM " + table + " WHERE rank='" + rank + "' and league='" + league + "'")
		db.MustExec(insertTable, data.CachedSince, league, temp["rank"], temp["dead"], temp["online"], character["name"],
			character["level"], character["class"], character["id"], exp,
			depth["default"], depth["solo"], account["name"], challenges["total"])

		db.MustExec(insertTableH, data.CachedSince, league, temp["rank"], temp["dead"], temp["online"], character["name"],
			character["level"], character["class"], character["id"], exp,
			depth["default"], depth["solo"], account["name"], challenges["total"])
		err = tx.Commit()
		if err != nil {
			log.Fatalln(err)
		}

		elapsed := time.Since(StartTime)
		fmt.Printf("[%v/3][%5v/15000] time: %v\n", x, rank, elapsed)
	}
	waitgroup.Done()
}

//MakeRequest MakeRequest
func MakeRequest(client http.Client, url string, ch chan<- string) {
	reqest, _ := http.NewRequest("GET", url, nil)
	response, _ := client.Do(reqest)
	body, _ := ioutil.ReadAll(response.Body)
	ch <- string(body)
}

//StartTime time
var StartTime = time.Now()

var waitgroup sync.WaitGroup

func main() {
	for x := 0; x < 4; x++ {
		var urls []string
		urlLeague := []string{
			"/Hardcore%20Legion?offset=", "/SSF%20Legion%20HC?offset=",
			"/SSF%20Legion?offset=", "/Legion?offset="}
		for offset := 0; offset <= 14800; offset += 200 {
			url := "http://api.pathofexile.com/ladders" + urlLeague[x] + strconv.Itoa(offset) + "&limit=200"
			urls = append(urls, url)
		}
		client := &http.Client{}
		ch := make(chan string, len(urls))
		for i := 0; i < len(urls); i += 5 {
			for _, url := range urls[i : i+5] {
				go MakeRequest(*client, url, ch)
			}
			for range urls[i : i+5] {
				waitgroup.Add(1)
				go saveData([]byte(<-ch), x)
			}
		}
	}
	waitgroup.Wait()
}
