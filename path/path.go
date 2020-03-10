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

func PathReplace(replacepattren *[]Replace, trfile string, frfile string, wg *sync.WaitGroup, comChannel chan string, errChannel chan error) {
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
	for _, pattern := range *replacepattren {
		fastresume["qBt-savePath"] = strings.ReplaceAll(fastresume["qBt-savePath"].(string), pattern.From, pattern.To)
		fastresume["save_path"] = strings.ReplaceAll(fastresume["save_path"].(string), pattern.From, pattern.To)
	}
	err = libtorrent.EncodeTorrentFile(frfile, fastresume)
	if err != nil {
		errChannel <- fmt.Errorf("Can't encode fastresume file %v. Error: %v", frfile, err)
		return
	}
	comChannel <- fmt.Sprintf("Changed save path for torrent: %-15v", torrentname)
	return
}
