export var settings;
(function (settings) {
    class SettingsStruct {
        constructor(source = {}) {
            if ('string' === typeof source)
                source = JSON.parse(source);
            this.version = source["version"];
            this.name = source["name"];
            this.user = source["user"];
            this.songDir = source["songDir"];
            this.dbPath = source["dbPath"];
            this.wsport = source["wsport"];
            this.logPath = source["logPath"];
            this.stateFile = source["stateFile"];
            this.musicdb_token = source["musicdb_token"];
            this.MeilliThreads = source["MeilliThreads"];
            this.MeilliRam = source["MeilliRam"];
            this.MeilliPort = source["MeilliPort"];
            this.AssetDir = source["AssetDir"];
        }
        static createFrom(source = {}) {
            return new SettingsStruct(source);
        }
    }
    settings.SettingsStruct = SettingsStruct;
})(settings || (settings = {}));
export var db;
(function (db) {
    class Song {
        constructor(source = {}) {
            if ('string' === typeof source)
                source = JSON.parse(source);
            this.id = source["id"];
            this.Title = source["Title"];
            this.Artist = source["Artist"];
            this.Album = source["Album"];
            this.Length = source["Length"];
            this.Image = source["Image"];
        }
        static createFrom(source = {}) {
            return new Song(source);
        }
    }
    db.Song = Song;
    class Playlist {
        constructor(source = {}) {
            if ('string' === typeof source)
                source = JSON.parse(source);
            this.id = source["id"];
            this.name = source["name"];
            this.songs = source["songs"];
            this.art = source["art"];
        }
        static createFrom(source = {}) {
            return new Playlist(source);
        }
    }
    db.Playlist = Playlist;
})(db || (db = {}));
export var global_state;
(function (global_state) {
    class Statemap {
        constructor(source = {}) {
            if ('string' === typeof source)
                source = JSON.parse(source);
            this.currentSong = this.convertValues(source["currentSong"], db.Song);
            this.elapsed = source["elapsed"];
            this.paused = source["paused"];
            this.volume = source["volume"];
            this.muted = source["muted"];
            this.repeat = source["repeat"];
            this.shuffle = source["shuffle"];
            this.next = this.convertValues(source["next"], db.Song);
            this.prev = this.convertValues(source["prev"], db.Song);
            this.activePlaylist = this.convertValues(source["activePlaylist"], db.Playlist);
            this.currentIndex = source["currentIndex"];
        }
        static createFrom(source = {}) {
            return new Statemap(source);
        }
        convertValues(a, classs, asMap = false) {
            if (!a) {
                return a;
            }
            if (a.slice) {
                return a.map(elem => this.convertValues(elem, classs));
            }
            else if ("object" === typeof a) {
                if (asMap) {
                    for (const key of Object.keys(a)) {
                        a[key] = new classs(a[key]);
                    }
                    return a;
                }
                return new classs(a);
            }
            return a;
        }
    }
    global_state.Statemap = Statemap;
})(global_state || (global_state = {}));
