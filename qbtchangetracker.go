package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	goflags "github.com/jessevdk/go-flags"
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

type Flags struct {
	QBitDir       string   `short:"d" long:"directory" description:"Destination directory BT_backup (as default)"`
	OldTracker    string   `short:"o" long:"oldtracker" description:"Old tracker"`
	NewTracker    string   `short:"n" long:"newtracker" description:"New tracker"`
	Replaces      []string `short:"r" long:"replace" description:"Replace paths.\n	Delimiter for from/to is comma - ,\n	Example: -r \"D:\\films,/home/user/films\" -r \"D:\\music,/home/user/music\"\n"`
	PathSeparator string   `long:"sep" description:"Default path separator that will use in all paths. You may need use this flag if you migrating from windows to linux in some cases"`
}

func main() {
	var wg sync.WaitGroup
	flags := Flags{PathSeparator: string(os.PathSeparator)}
	sep := string(os.PathSeparator)
	switch OS := runtime.GOOS; OS {
	case "windows":
		flags.QBitDir = os.Getenv("LOCALAPPDATA") + sep + "qBittorrent" + sep + "BT_backup" + sep
		flags.QBitDir = strings.Join([]string{os.Getenv("LOCALAPPDATA"), "qBittorrent", "BT_backup"}, sep)
	case "darwin":
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}
		flags.QBitDir = strings.Join([]string{usr.HomeDir, "Library", "Application Support", "QBittorrent", "BT_backup"}, sep)
	case "linux":
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}
		flags.QBitDir = strings.Join([]string{usr.HomeDir, ".local", "share", "data", "qBittorrent", "BT_backup"}, sep)
	}
	parser := goflags.NewParser(&flags, goflags.Default)
	if _, err := parser.Parse(); err != nil { // https://godoc.org/github.com/jessevdk/go-flags#ErrorType
		if flagsErr, ok := err.(*goflags.Error); ok && flagsErr.Type == goflags.ErrHelp {
			os.Exit(0)
		} else {
			log.Println(err)
			time.Sleep(30 * time.Second)
			os.Exit(1)
		}
	}
	sepdefined := parser.FindOptionByLongName("sep")
	if flags.QBitDir[len(flags.QBitDir)-1] != os.PathSeparator {
		flags.QBitDir += string(os.PathSeparator)
	}

	if _, err := os.Stat(flags.QBitDir); os.IsNotExist(err) {
		log.Println(err)
		time.Sleep(30 * time.Second)
		os.Exit(1)
	}

	files, _ := filepath.Glob(flags.QBitDir + "*fastresume")
	comChannel := make(chan string, len(files))
	errChannel := make(chan error, len(files))

	var replacepatterns []path.Replace

	if len(flags.Replaces) != 0 || sepdefined != nil {
		for _, str := range flags.Replaces {
			patterns := strings.Split(str, ",")
			if len(patterns) < 2 {
				continue
			}
			replacepatterns = append(replacepatterns, path.Replace{From: patterns[0], To: patterns[1]})
		}
	} else {
		if flags.OldTracker == "" || flags.NewTracker == "" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter old tracker: ")
			flags.OldTracker, _ = reader.ReadString('\n')
			flags.OldTracker = flags.OldTracker[:len(flags.OldTracker)-2]
			fmt.Print("Enter new tracker: ")
			flags.NewTracker, _ = reader.ReadString('\n')
			flags.NewTracker = flags.NewTracker[:len(flags.NewTracker)-2]
		}
	}
	color.HiRed("Check that the qBittorrent is turned off and the directory %v and is backed up.\n\n",
		flags.QBitDir)
	fmt.Println("Press Enter to start")
	fmt.Scanln()
	log.Println("Started")
	if len(flags.Replaces) != 0 || sepdefined != nil {
		for _, frfile := range files {
			trfile := strings.TrimSuffix(frfile, filepath.Ext(frfile)) + ".torrent"
			if _, err := os.Stat(trfile); os.IsNotExist(err) {
				continue
			}
			go path.PathReplace(&replacepatterns, trfile, frfile, flags.PathSeparator, &wg, comChannel, errChannel)
			wg.Add(1)
		}
		goto End
	}
	for _, frfile := range files {
		trfile := strings.TrimSuffix(frfile, filepath.Ext(frfile)) + ".torrent"
		if _, err := os.Stat(trfile); os.IsNotExist(err) {
			continue
		}
		go tracker.ChangeTracker(&flags.OldTracker, &flags.NewTracker, trfile, frfile, &wg, comChannel, errChannel)
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
