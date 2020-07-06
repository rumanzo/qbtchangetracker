package path

import (
	"fmt"
	libtorrent "github.com/rumanzo/qbtchangetracker/torrent"
	"runtime/debug"
	"strings"
	"sync"
)

type Replace struct {
	From, To string
}

func PathReplace(replacepattren *[]Replace, trfile string, frfile string, separator string, wg *sync.WaitGroup, comChannel chan string, errChannel chan error) {
	defer wg.Done()
	defer func() {
		if r := recover(); r != nil {
			errChannel <- fmt.Errorf(
				"Panic while processing torrent file %v:\n======\nReason: %v.\nText panic:\n%v\n======",
				frfile, r, string(debug.Stack()))
		}
	}()
	fastresume, err := libtorrent.DecodeFastresumeFile(frfile)
	if err != nil {
		errChannel <- fmt.Errorf("Can't decode fastresume file %v. Error: %v", frfile, err)
		return
	}
	torrent, err := libtorrent.DecodeFastresumeFile(trfile)
	if err != nil {
		errChannel <- fmt.Errorf("Can't decode torrent file %v. Error: %v", trfile, err)
		return
	}
	torrentname := torrent["info"].(map[string]interface{})["name"].(string)
	changed := false
	for _, pattern := range *replacepattren {
		newqBtsavePath := strings.ReplaceAll(fastresume["qBt-savePath"].(string), pattern.From, pattern.To)
		if fastresume["qBt-savePath"] != newqBtsavePath {
			fastresume["qBt-savePath"] = newqBtsavePath
			changed = true
		}
		newsavepath := strings.ReplaceAll(fastresume["save_path"].(string), pattern.From, pattern.To)
		if fastresume["save_path"] != newsavepath {
			fastresume["save_path"] = newsavepath
			changed = true
		}
	}

	var oldsep string
	switch separator {
	case "\\":
		oldsep = "/"
	case "/":
		oldsep = "\\"
	}
	newqBtsavePath := strings.ReplaceAll(fastresume["qBt-savePath"].(string), oldsep, separator)
	if fastresume["qBt-savePath"] != newqBtsavePath {
		fastresume["qBt-savePath"] = newqBtsavePath
		changed = true
	}
	if _, ok := fastresume["mapped_files"]; ok {
		for num, entry := range fastresume["mapped_files"].([]interface{}) {
			newentry := strings.ReplaceAll(entry.(string), oldsep, separator)
			if entry != newentry {
				fastresume["mapped_files"].([]interface{})[num] = newentry
				changed = true
			}
		}
	}

	err = libtorrent.EncodeTorrentFile(frfile, fastresume)
	if err != nil {
		errChannel <- fmt.Errorf("Can't encode fastresume file %v. Error: %v", frfile, err)
		return
	}
	if changed {
		comChannel <- fmt.Sprintf("Changed save path for torrent: %-15v", torrentname)
	}

	return
}
