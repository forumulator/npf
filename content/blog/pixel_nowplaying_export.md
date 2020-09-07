+++
title = "Exporting your Pixel 'Now playing' history"
date = 2020-09-07T22:42:25+02:00
type = "post"
tags = ["pixel", "nowplaying", "android"]
+++

Google Pixel phones have the handy 'now playing' feature, which is like auto-Shazam - it allows the Pixel to recognize which song is playing in the background and make it visible to you *right* at your lock screen. It also allows you to view the now playing history and export a song to Play Music / Spotify / Apple Music from settings. Having this enabled by default, I have been accruing songs on the 'Now playing' for about a year and a half, which amounts to about 1500 songs.

But, as I recently discovered to my endless frustration when I decided to move off my Pixel, Google doesn't give you a straightforward way to bulk export the list to CSV / Spotify / text. Individually 'Share'-ing every one of 1500 songs and adding it to a Spotify playlist would take me hours, which is simply unacceptable. 

It turns out that it's fairly simple to use `adb backup` to extract the 'Now playing' history database from the phone (as of the time of this writing, on Android 9, Pixel 3XL), and then use a simple script to extract data in whatever format is most convenient. I'm going to use shell utilities to generate a CSV with `timestamp` and `spotify_track_id`.

The pre-requisites are very simple - you need a working adb setup on your computer handy. There are plenty of guides for all operating systems (for example [this](https://www.xda-developers.com/install-adb-windows-macos-linux/) one), so I won't elaborate on that here. Once adb is set up properly, connecting your phone via a USB cable and then running `./adb devices` on the shell should output something like:

```
$ ./adb.exe devices
List of devices attached
8BPK0YFHJ device
```

(I'm using my Ubuntu in Windows 10 installation so you'll see bash commands with `.exe` executables in this post). Once you have adb showing your device, you can continue with the steps below.

### Backing up the app data
This took me a bit of time because it's not clear which app 'Now playing' is associated with. You can access the history in Settings > Sound, but it really wouldn't make sense for it to be integrated with the settings app. So it should contain, for example, to 'now playing' or 'music' in the name, right? Wrong. After a bit of elimination from system apps plus a bit of Googling tells me that the correct app is **Pixel Ambient Servies**.

Next, to figure out which package the app corresponds to, I list the system packages (`./adb.exe shell pm packages list -s`) and grep for the phrase `com.google`; this makes it fairly easy. The package is `com.google.intelligence.sense`.

The next part requires a bit of luck, and thankfully, Google hasn't disabled the `android:allowBackup` flag in the app. If they had, this would have been a dead end (unless you have a rooted phone). The only option then probably would have been to create a touch command based adb script to export each of the songs to a Spotify playlist. And while that wouldn't have required manual work (beyond the scripting), it would have still taken time. That is the last thing I want to do. 

But luckily adb debugging on `com.google.intelligence.sense` remains enabled, either because of an oversight, or as an intended 'feature'. We can then backup all app data to our pc:

```
$ ./adb.exe backup -f ambient.ab com.google.intelligence.sense
WARNING: adb backup is deprecated and may be removed in a future release
Now unlock your device and confirm the backup operation...
```

(Another useful feature that's deprecated by Google, what a surprise). Confirm the backup on your device, and it should finish in under a minute. My app has a lot of something called Superpacks which inflate the backup size to around 170 Mb.

### Extracting the database file
Now that we've extracted the app data, it's time to un-archive it and extract the db. Android backups are essentially tar gzip archives with a bit of a header twist. There are many ways to extract them (including a java utility called [Android backup extractor](https://github.com/nelenkov/android-backup-extractor)), but I just used the following command:

```sh
( printf "\x1f\x8b\x08\x00\x00\x00\x00\x00" ; tail -c +25 ambient.ab ) | tar xfvz -
```

In the above command we pre-pend the `\x1f\x8b\x08\x00\x00\x00\x00\x00` header to the `ambient.ab` file and then extract it using `tar xvzf -`.

This should create an `apps/com.google.intelligence.sense` in the current directory, which is where the app data resides. We're only interested in the db file `apps/com.google.intelligence.sense/db/history_db`, which is a SQLite database file (another bit of luck is that it's un-encrypted).

### Figuring out the 'Now playing' history database
Install [SQLlite](https://www.sqlite.org/download.html), and let's take a look at the database. The command `.tables` lists all tables in the database:

```
sqlite> .tables
android_metadata recognition_history
```

Great, that seems pretty straightforward. The `android metadata` table has only one row with the locale, so that's pretty much useless. Let's look at the `recognition_history` table:

```sql
sqlite> PRAGMA table_info(recognition_history);
0|timestamp|LONG|0||1
1|history_entry|BLOB|0||0
2|track_id|STRING|0||0
3|shards_region|STRING|0||0
4|downloaded_shards_version|INTEGER|0||0
5|core_shard_version|INTEGER|0||0
```

Ah, pay dirt! The `history_entry` seems to be a BLOB, but if you take a look at a single row:

```
sqlite> select * from recognition_history limit 1;
Europe/WarsawJ
/g/11c6tt7zl9Jyandroid-app://com.apple.android.music/https/itunes.apple.com/us/album/confidently-lost/1209647010?i=1209647307&ign-gact=2/m/013d79knJH:https://www.deezer.com/track/143171428/radio?autoplay=true
/m/03c2ckdJVHandroid-app://com.spotify.music/spotify/track/0RTnMChwd0ARVY3zOyNdaP?v=T
/m/04yhd6cJpbhttps://www.iheart.com/artist/Sabrina+Claudio-31483691/songs/Orion%27s+Belt-45519571?autoplay=true
/m/0c3zxqxJyandroid-app://com.google.android.music/http/play.google.com/music/m/Tigjmk45m2k2yqavlr7edrkcut4?signup_if_needed=1&play=1
t_tfgj6vzLvnUjConfidently Lostp|Tigjmk45m2k2yqavlr7edrkcut4|xa|100|2
```

It looks like the Blob format is just Unicode data. The best part, though, is that this is almost as good as exporting directly to the playlist. Because we have the Spotify track id, the Play Music track id, Deezer id, and Apple Music id here, you can directly (using APIs or a bulk import tool) add the tracks to a playlist regardless of what subscription service you use (as long as it's one of these!). 

### Exporting the database
Finally, let's open the SQLite db, query the `recognition_history` for all rows, export it to a file, and convert each row to the desired format. My first instinct was to write a Python script, but that's lengthier and more error prone than just using some `sed` to convert each record into the desired format. Besides, it doesn't require any environment setup.

I'm going to extract the `timestamp` and `spotify_track_id` into a CSV file. First, we want to export the columns to a file. Because the `history_entry` contains '\n' and possibly '\r' characters, we need to remove them so we can have one record per line and use `sed` conveniently (hence the use of the `replace()` function in the `SELECT` below).

```sql
sqlite> .headers on
sqlite> .mode csv
sqlite> .output results.csv
sqlite> select timestamp, replace(replace(history_entry, CHAR(10), '|'), CHAR(13), '|') as spotify_track_id from recognition_history;
```

This should create a file `results.csv` in your current directory with the query rows. Before continuing, ensure that there's one query row per line (there are still non-printable characters and Unicode in the BLOB, that's fine).

Now let's remove the Unicode characters so `sed` doesn't give us problems. Some versions of `sed` work with Unicode, but I just spent hours debugging locale issues over this so I'm not a fan. To remove Unicode characters, we use a simple `sed` command that searches for characters in the range `\x80-\xFF` (i.e decimal 128 - 255, out of the ASCII range). The `LC_ALL=C` sets the locale to ISO 'C' and this is the only one that worked properly for me.

```sh
LC_ALL=C sed -i -r 's:[\x80-\xFF]:.:g' results.csv
```

The `-i` switch above ensures that the changes are made to the input file in place. To verify that all Unicode characters are gone, you can check that `grep --color=auto -P '[\x80-\xFF]' results.csv` returns no matches. After this, all that's left to do is just `sed` out the timestamp and Spotify track id:

```sh
sed -i -r 's:^([0-9]+).*spotify/track/([A-Za-z0-9]+).*:\1,\2:g' results.csv
```

The regex in the above command matches first a timestamp at the beginning of the line (`^[0-9]+`), then matches the rest of the line, taking note of the string `spotify/track/[A-Za-z0-9]+` (i.e, extracting the track id).

## Conclusion
And that's it. `head -n 5` your `results.csv` file and you should see the timestamp and track ids, one on each line. A bulk add tool or some API calls will help you add the tracks to a playlist (beware of the API quotas though). The above solution is simple to modify (by changing the regex) or extend (by scripting) to fit whatever end goals you have for your tracks in mind.
