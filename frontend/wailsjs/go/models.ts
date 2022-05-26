export namespace db {
	
	export class Song {
	    id: string;
	    Title: string;
	    Artist: string;
	    Album: string;
	    Length: number;
	    Image: string;
	
	    static createFrom(source: any = {}) {
	        return new Song(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.Title = source["Title"];
	        this.Artist = source["Artist"];
	        this.Album = source["Album"];
	        this.Length = source["Length"];
	        this.Image = source["Image"];
	    }
	}
	export class Playlist {
	    id: string;
	    name: string;
	    songs: string[];
	    art: string;
	
	    static createFrom(source: any = {}) {
	        return new Playlist(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.songs = source["songs"];
	        this.art = source["art"];
	    }
	}

}

export namespace global_state {
	
	export class Statemap {
	    currentSong: db.Song;
	    elapsed: number;
	    paused: boolean;
	    volume: number;
	    muted: boolean;
	    repeat: boolean;
	    shuffle: boolean;
	    next: db.Song;
	    prev: db.Song;
	    activePlaylist: db.Playlist;
	    queue: Song[];
	    currentIndex: number;
	
	    static createFrom(source: any = {}) {
	        return new Statemap(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
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
	        this.queue = this.convertValues(source["queue"], Song);
	        this.currentIndex = source["currentIndex"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
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

}

export namespace settings {
	
	export class SettingsStruct {
	    version: string;
	    name: string;
	    user: string;
	    songDir: string;
	    dbPath: string;
	    wsport: number;
	    logPath: string;
	    stateFile: string;
	    musicdb_token: string;
	    MeilliThreads: number;
	    MeilliRam: number;
	    MeilliPort: number;
	    AssetDir: string;
	
	    static createFrom(source: any = {}) {
	        return new SettingsStruct(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
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
	}

}

