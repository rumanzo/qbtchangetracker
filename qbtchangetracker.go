package main

import (
	"log"
	"fmt"
	"path/filepath"
	"io/ioutil"
	"github.com/zeebo/bencode"
	"os"
	"bufio"
	"sync"
)

func decodetorrentfile(path string) map[string]interface{} {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	var torrent map[string]interface{}
	if err := bencode.DecodeBytes([]byte(dat), &torrent); err != nil {
		log.Fatal(err)
	}
	return torrent
}

func encodetorrentfile(path string, newstructure map[string]interface{}) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Create(path)
	}

	file, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()
	bufferedWriter := bufio.NewWriter(file)
	enc := bencode.NewEncoder(bufferedWriter)
	if err := enc.Encode(newstructure); err != nil {
		log.Fatal(err)
	}
	bufferedWriter.Flush()
	return nil
}


func changetracker(oldtracker *string, newtracker *string, path string, wg *sync.WaitGroup)  {
	defer wg.Done()
	fastresume := decodetorrentfile(path)
	for num, list := range fastresume["trackers"].([]interface{}) {
		if list.([]interface{})[0] == *oldtracker {
			fastresume["trackers"].([]interface{})[num].([]interface{})[0] = *newtracker
			fmt.Printf("Changed tracker in %v\n", path)
			encodetorrentfile(path, fastresume)
		}
	}
}

func echo(file string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(file)
}

func main() {
	var wg sync.WaitGroup
	directory := os.Getenv("LOCALAPPDATA") + "\\qBittorrent\\BT_backup\\"
	files, err := filepath.Glob(directory + "*fastresume")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter old tracker: ")
	oldtracker, _ := reader.ReadString('\n')
	oldtracker = oldtracker[:len(oldtracker) - 2]
	fmt.Print("Enter new tracker: ")
	newtracker, _ := reader.ReadString('\n')
	newtracker = newtracker[:len(newtracker) - 2]
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		go changetracker(&oldtracker, &newtracker, file, &wg)
		wg.Add(1)
	}
	wg.Wait()
	fmt.Print("Pres Enter to exit")
	reader.ReadString('\n')
}
