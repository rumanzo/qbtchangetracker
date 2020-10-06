


# qbtchangetracker
qbtchangetracker is cli tool for iterative change trackers in qBittorrent fastresume files
- [qbtchangetracker](#qbtchangetracker)
	- [Feature](#user-content-feature)
	- [Help](#user-content-help)
	- [Usage examples](#user-content-usage-examples)
	- [Known issuses](#user-content-known-issuses)
	
Feature:
---------
 - Processing all fastresume files
 - Multithreading


> [!IMPORTANT]
> Don't forget before use make backup qbittorrent folder. Close qBittorrent before use.

Help:
-------

Help (from cmd or powerwhell)

```
PS C:\Downloads>qbtchangetracker_v1.4_amd64.exe -h
Usage:
  qbtchangetracker_v1.4_amd64.exe [OPTIONS]

Application Options:
  -d, --directory=  Destination directory BT_backup (as default) (default:
                    C:\Users\user\AppData\Local\qBittorrent\BT_backup)
  -o, --oldtracker= Old tracker
  -n, --newtracker= New tracker
  -r, --replace=    Replace paths.
                    Delimiter for from/to is comma - ,
                    Example: -r "D:\films,/home/user/films" -r
                    "D:\music,/home/user/music"

      --sep=        Default path separator that will use in all paths. You may
                    need use this flag if you migrating from windows to linux
                    in some cases (default: \)

Help Options:
  -h, --help        Show this help message


```

Usage examples:
----------------

- If you just run application, it will processing torrents from %APPDATA%\uTorrent\ to %LOCALAPPDATA%\qBittorrent\BT_BACKUP\ and ask interactively old tracker and new tracker
```
PS C:\qbtchangetracker> .\qbtchangetracker_v1.4_amd64.exe
Enter old tracker: oldtracker
Enter new tracker: newtracker
Check that the qBittorrent is turned off and the directory C:\Users\user\AppData\Local\qBittorrent\BT_backup\ and is backed up.

Press Enter to start

2019/03/13 00:11:39 Started
1  Changed tracker for torrent: torrentname1
2  Changed tracker for torrent: torrentname2
```

- Run application from cmd or powershell with keys, if you want change source dir or destination dir, or export/import behavior
```
PS C:\qbtchangetracker> .\qbtchangetracker_v1.4_amd64.exe -d C:\temp\BT_backup\ -o oldtracker -n newtracker
Check that the qBittorrent is turned off and the directory C:\temp\BT_backup\ and is backed up.

Press Enter to start

2019/03/13 00:11:39 Started
1  Changed tracker for torrent: torrentname1
2  Changed tracker for torrent: torrentname2
Press Enter to exit
```
- Run with replace and\or sep key (if you migrate to different OS)
```
PS C:\qbtchangetracker> .\qbtchangetracker_v1.4_amd64.exe -d C:\temp\BT_backup\ --sep '/' -r 'D:/films,/mnt/d/films' -r 'D:\music,/mnt/d/music'
Check that the qBittorrent is turned off and the directory C:\temp\BT_backup\ and is backed up.

Press Enter to start

2019/03/13 00:11:39 Started
1  Changed save path for torrent: torrentname1
2  Changed save path for torrent: torrentname2
Press Enter to exit
```
Known issuses:
---------------
 - Unknown
