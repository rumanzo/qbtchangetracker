


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
PS C:\Users\user\go\src\qbtchangetracker> .\qbtchangetracker_v1.0_amd64.exe -h
Usage of C:\Users\user\go\src\qbtchangetracker\qbtchangetracker_v1.0_amd64.exe:
-d, --directory (= "C:\\Users\\user\\AppData\\Local\\qBittorrent\\BT_backup\\")
    Destination directory BT_backup (as default)
-n, --newtracker (= "")
    New tracker
-o, --oldtracker (= "")
    Old tracker
--replace (= "")
    Replace paths.
        Delimiter for replaces - ;
        Delimiter for from/to - ,
        Example: "D:\films,/home/user/films;\,/"
        If you use path separator different from you system, declare it mannually
```

Usage examples:
----------------

- If you just run application, it will processing torrents from %APPDATA%\uTorrent\ to %LOCALAPPDATA%\qBittorrent\BT_BACKUP\ and ask interactively old tracker and new tracker
```
PS C:\Users\user\go\src\qbtchangetracker> .\qbtchangetracker.exe
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
PS C:\Users\user\go\src\qbtchangetracker> .\qbtchangetracker_v1.0_amd64.exe -d C:\temp\BT_backup\ -o oldtracker -n newtracker
Check that the qBittorrent is turned off and the directory C:\temp\BT_backup\ and is backed up.

Press Enter to start

2019/03/13 00:11:39 Started
1  Changed tracker for torrent: torrentname1
2  Changed tracker for torrent: torrentname2
Press Enter to exit
```
Known issuses:
---------------
 - Unknown