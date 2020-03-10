package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/juju/gnuflag"
	"github.com/rumanzo/qbtchangetracker/path"
	"github.com/rumanzo/qbtchangetracker/tracker"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var qbitdir, oldtracker, newtracker, replace string
	sep := string(os.PathSeparator)
	switch OS := runtime.GOOS; OS {
	case "windows":
		qbitdir = os.Getenv("LOCALAPPDATA") + sep + "qBittorrent" + sep + "BT_backup" + sep
		qbitdir = strings.Join([]string{os.Getenv("LOCALAPPDATA"), "qBittorrent", "BT_backup"}, sep)
	case "darwin":
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}
		qbitdir = strings.Join([]string{usr.HomeDir, "Library", "Application Support", "QBittorrent", "BT_backup"}, sep)
	case "linux":
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}
		qbitdir = strings.Join([]string{usr.HomeDir, ".local", "share", "data", "qBittorrent", "BT_backup"}, sep)
	}

	gnuflag.StringVar(&qbitdir, "directory", qbitdir,
		"Destination directory BT_backup (as default)")
	gnuflag.StringVar(&qbitdir, "d", qbitdir,
		"Destination directory BT_backup (as default)")
	gnuflag.StringVar(&oldtracker, "oldtracker", "",
		"Old tracker")
	gnuflag.StringVar(&oldtracker, "o", "",
		"Old tracker")
	gnuflag.StringVar(&newtracker, "newtracker", "",
		"New tracker")
	gnuflag.StringVar(&newtracker, "n", "",
		"New tracker")
	gnuflag.StringVar(&replace, "replace", "", "Replace paths.\n	"+
		"Delimiter for replaces - ;\n	"+
		"Delimiter for from/to - ,\n	Example: \"D:\\films,/home/user/films;\\,/\"\n	"+
		"If you use path separator different from you system, declare it mannually")
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
	comChannel := make(chan string, len(files))
	errChannel := make(chan error, len(files))

	var replacepatterns []path.Replace

	if replace != "" {
		for _, str := range strings.Split(replace, ";") {
			patterns := strings.Split(str, ",")
			if len(patterns) < 2 {
				log.Println("Bad replace pattern")
				time.Sleep(30 * time.Second)
				os.Exit(1)
			}
			replacepatterns = append(replacepatterns, path.Replace{
				From: patterns[0],
				To:   patterns[1],
			})
		}
	} else {
		if oldtracker == "" || newtracker == "" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter old tracker: ")
			oldtracker, _ = reader.ReadString('\n')
			oldtracker = oldtracker[:len(oldtracker)-2]
			fmt.Print("Enter new tracker: ")
			newtracker, _ = reader.ReadString('\n')
			newtracker = newtracker[:len(newtracker)-2]
		}
	}
	color.HiRed("Check that the qBittorrent is turned off and the directory %v and is backed up.\n\n",
		qbitdir)
	fmt.Println("Press Enter to start")
	fmt.Scanln()
	log.Println("Started")
	if replace != "" {
		for _, frfile := range files {
			trfile := strings.TrimSuffix(frfile, filepath.Ext(frfile)) + ".torrent"
			if _, err := os.Stat(trfile); os.IsNotExist(err) {
				continue
			}
			go path.PathReplace(&replacepatterns, trfile, frfile, &wg, comChannel, errChannel)
			wg.Add(1)
		}
		goto End
	}
	for _, frfile := range files {
		trfile := strings.TrimSuffix(frfile, filepath.Ext(frfile)) + ".torrent"
		if _, err := os.Stat(trfile); os.IsNotExist(err) {
			continue
		}
		go tracker.ChangeTracker(&oldtracker, &newtracker, trfile, frfile, &wg, comChannel, errChannel)
		wg.Add(1)
	}
End:
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
