package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/juju/gnuflag"
	"github.com/zeebo/bencode"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

func decodefastresumefile(path string) (map[string]interface{}, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var torrent map[string]interface{}
	if err := bencode.DecodeBytes([]byte(dat), &torrent); err != nil {
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

func changetracker(oldtracker *string, newtracker *string, trfile string, frfile string, wg *sync.WaitGroup, comChannel chan string, errChannel chan error) {
	defer wg.Done()
	defer func() {
		if r := recover(); r != nil {
			errChannel <- fmt.Errorf(
				"Panic while processing torrent file %v:\n======\nReason: %v.\nText panic:\n%v\n======",
				frfile, r, string(debug.Stack()))
		}
	}()
	fastresume, err := decodefastresumefile(frfile)
	if err != nil {
		errChannel <- fmt.Errorf("Can't decode fastresume file %v. Error: %v", frfile, err)
		return
	}
	torrent, err := decodefastresumefile(trfile)
	if err != nil {
		errChannel <- fmt.Errorf("Can't decode torrent file %v. Error: %v", trfile, err)
		return
	}
	torrentname := torrent["info"].(map[string]interface{})["name"].(string)
	for num, list := range fastresume["trackers"].([]interface{}) {
		for n, tracker := range list.([]interface{}) {
			if tracker == *oldtracker {
				fastresume["trackers"].([]interface{})[num].([]interface{})[n] = *newtracker
				comChannel <- fmt.Sprintf("Changed tracker for torrent: %-15v", torrentname)
				err = encodetorrentfile(frfile, fastresume)
				if err != nil {
					errChannel <- fmt.Errorf("Can't encode fastresume file %v. Error: %v", frfile, err)
					return
				}
			}
		}
	}
	return
}

func main() {
	var wg sync.WaitGroup
	var qbitdir, oldtracker, newtracker string
	gnuflag.StringVar(&qbitdir, "directory", os.Getenv("LOCALAPPDATA")+"\\qBittorrent\\BT_backup\\",
		"Destination directory BT_backup (as default)")
	gnuflag.StringVar(&qbitdir, "d", os.Getenv("LOCALAPPDATA")+"\\qBittorrent\\BT_backup\\",
		"Destination directory BT_backup (as default)")
	gnuflag.StringVar(&oldtracker, "oldtracker", "",
		"Old tracker")
	gnuflag.StringVar(&oldtracker, "o", "",
		"Old tracker")
	gnuflag.StringVar(&newtracker, "newtracker", "",
		"New tracker")
	gnuflag.StringVar(&newtracker, "n", "",
		"New tracker")
	gnuflag.Parse(true)

	if qbitdir[len(qbitdir)-1] != os.PathSeparator {
		qbitdir += string(os.PathSeparator)
	}

	if _, err := os.Stat(qbitdir); os.IsNotExist(err) {
		log.Println(err)
		time.Sleep(30 * time.Second)
		os.Exit(1)
	}
	files, _ := filepath.Glob(qbitdir + "*fastresume")
	if oldtracker == "" || newtracker == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter old tracker: ")
		oldtracker, _ = reader.ReadString('\n')
		oldtracker = oldtracker[:len(oldtracker)-2]
		fmt.Print("Enter new tracker: ")
		newtracker, _ = reader.ReadString('\n')
		newtracker = newtracker[:len(newtracker)-2]
	}
	color.HiRed("Check that the qBittorrent is turned off and the directory %v and is backed up.\n\n",
		qbitdir)
	fmt.Println("Press Enter to start")
	fmt.Scanln()
	comChannel := make(chan string, len(files))
	errChannel := make(chan error, len(files))
	log.Println("Started")
	for _, frfile := range files {
		trfile := strings.TrimSuffix(frfile, filepath.Ext(frfile)) + ".torrent"
		if _, err := os.Stat(trfile); os.IsNotExist(err) {

			continue
		}
		go changetracker(&oldtracker, &newtracker, trfile, frfile, &wg, comChannel, errChannel)
		wg.Add(1)
	}

	waserrors := false
	numjob := 1
	go func() {
		wg.Wait()
		close(comChannel)
		close(errChannel)
	}()
	for message := range comChannel {
		fmt.Printf("%-5v %v \n", numjob, message)
		numjob++
	}
	for message := range errChannel {
		color.HiRed(fmt.Sprintf("%-5v %v \n", numjob, message))
		waserrors = true
		numjob++
	}
	if waserrors {
		log.Println("Not all torrents was processed")
	}
	fmt.Print("Press Enter to exit")
	fmt.Scanln()
}
