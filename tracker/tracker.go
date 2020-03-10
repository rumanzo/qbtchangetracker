package tracker

import (
	"fmt"
	libtorrent "github.com/rumanzo/qbtchangetracker/torrent"
	"runtime/debug"
	"sync"
)

func ChangeTracker(oldtracker *string, newtracker *string, trfile string, frfile string, wg *sync.WaitGroup, comChannel chan string, errChannel chan error) {
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
	for num, list := range fastresume["trackers"].([]interface{}) {
		for n, tracker := range list.([]interface{}) {
			if tracker == *oldtracker {
				fastresume["trackers"].([]interface{})[num].([]interface{})[n] = *newtracker
			}
		}
	}
	err = libtorrent.EncodeTorrentFile(frfile, fastresume)
	if err != nil {
		errChannel <- fmt.Errorf("Can't encode fastresume file %v. Error: %v", frfile, err)
		return
	}
	comChannel <- fmt.Sprintf("Changed tracker for torrent: %-15v", torrentname)
	return
}
