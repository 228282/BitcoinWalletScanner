package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/piotrnar/gocoin/lib/btc"
)

type Answer struct {
	Address      string  `json:"address"`
	FinalBalance float64 `json:"final_balance"`
}

func MakeRequest(address string) string {
	resp, err := http.Get(fmt.Sprintf(`https://blockchain.info/address/%s?format=json`, address))
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return string(body[:])
}

func FindBalance(wg *sync.WaitGroup) {
	wg.Add(1)

	fi, err := os.Create(time.Now().Format("2006.01.02 15.04.05") + ".txt")
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	var hasBalance []string

	privateKey := make([]byte, 32)

	for {

		rand.Read(privateKey)
		publicKey := btc.PublicFromPrivate(privateKey, false)
		address := btc.NewAddrFromPubkey(publicKey, 0x00).String()
		//address := "16jY7qLJnxb7CHZyqBP8qca9d51gAjyXQN"

		body := MakeRequest(address)
		var ans Answer
		json.Unmarshal([]byte(body), &ans)

		if ans.Address == "" {
			fmt.Println("Error")
			time.Sleep(time.Second * 60)
			continue

		}

		fmt.Printf("%s %x %s\n", address, privateKey, fmt.Sprintf("%f", ans.FinalBalance))

		if ans.FinalBalance != 0 {
			_, err := fi.WriteString(fmt.Sprintf("%s %x\n", address, privateKey))
			hasBalance = append(hasBalance, address)
			hasBalance = append(hasBalance, string(privateKey))
			if err != nil {
				bufio.NewReader(os.Stdin).ReadBytes('\n')
			}
		}

		fmt.Println(hasBalance)
	}
}

func main() {
	wg := &sync.WaitGroup{}

	for i := 0; i < 1; i++ {
		go FindBalance(wg)
		time.Sleep(time.Second)
	}

	wg.Wait()

}
