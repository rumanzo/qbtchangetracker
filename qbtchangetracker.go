package main

import (
	"bufio"
	"fmt"
	"github.com/zeebo/bencode"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func decodetorrentfile(path string) (map[string]interface{}, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var torrent map[string]interface{}
	if err := bencode.DecodeBytes([]byte(dat), &torrent); err != nil {
		log.Println(err)
		return nil, err
	}
	return torrent, nil
}

func encodetorrentfile(path string, newstructure map[string]interface{}) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Create(path)
	}

	file, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()
	bufferedWriter := bufio.NewWriter(file)
	enc := bencode.NewEncoder(bufferedWriter)
	if err := enc.Encode(newstructure); err != nil {
		log.Println(err)
		return err
	}
	bufferedWriter.Flush()
	return nil
}

func changetracker(oldtracker *string, newtracker *string, path string, wg *sync.WaitGroup) error {
	defer wg.Done()
	fastresume, err := decodetorrentfile(path)
	if err != nil {
		log.Println(err)
		return err
	}
	for num, list := range fastresume["trackers"].([]interface{}) {
		if list.([]interface{})[0] == *oldtracker {
			fastresume["trackers"].([]interface{})[num].([]interface{})[0] = *newtracker
			fmt.Printf("Changed tracker in %v\n", path)
			err = encodetorrentfile(path, fastresume)
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}

func main() {
	var wg sync.WaitGroup
	directory := os.Getenv("LOCALAPPDATA") + "\\qBittorrent\\BT_backup\\"
	files, _ := filepath.Glob(directory + "*fastresume")
	_, err := os.Stat(directory)
	if err != nil {
		log.Println(err)
		time.Sleep(30 * time.Second)
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter old tracker: ")
	oldtracker, _ := reader.ReadString('\n')
	oldtracker = oldtracker[:len(oldtracker)-2]
	fmt.Print("Enter new tracker: ")
	newtracker, _ := reader.ReadString('\n')
	newtracker = newtracker[:len(newtracker)-2]
	for _, file := range files {
		go changetracker(&oldtracker, &newtracker, file, &wg)
		wg.Add(1)
	}
	wg.Wait()
	fmt.Print("Press Enter to exit")
	reader.ReadString('\n')
}
