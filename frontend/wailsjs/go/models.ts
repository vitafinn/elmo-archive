export namespace handler {
	
	export class AcquisitionLoss {
	    ID: string;
	    Name: string;
	    Count: number;
	    BannerID: number;
	    BannerName: string;
	
	    static createFrom(source: any = {}) {
	        return new AcquisitionLoss(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Count = source["Count"];
	        this.BannerID = source["BannerID"];
	        this.BannerName = source["BannerName"];
	    }
	}
	export class BannerPeriod {
	    id: number;
	    name: string;
	    startTime: number;
	    endTime: number;
	    lastTimestamp: number;
	    status: string;
	    totalCount: number;
	    rank5Count: number;
	    completedPulls: number;
	    rank5Avg: number;
	    rank5Spent: number;
	    tailCount: number;
	    targetAttempts: number;
	    targetWins: number;
	    unknownOutcomes: number;
	    unknownRankCount: number;
	
	    static createFrom(source: any = {}) {
	        return new BannerPeriod(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.startTime = source["startTime"];
	        this.endTime = source["endTime"];
	        this.lastTimestamp = source["lastTimestamp"];
	        this.status = source["status"];
	        this.totalCount = source["totalCount"];
	        this.rank5Count = source["rank5Count"];
	        this.completedPulls = source["completedPulls"];
	        this.rank5Avg = source["rank5Avg"];
	        this.rank5Spent = source["rank5Spent"];
	        this.tailCount = source["tailCount"];
	        this.targetAttempts = source["targetAttempts"];
	        this.targetWins = source["targetWins"];
	        this.unknownOutcomes = source["unknownOutcomes"];
	        this.unknownRankCount = source["unknownRankCount"];
	    }
	}
	export class DisplayRecord {
	    ID: string;
	    ItemID: number;
	    ItemType: number;
	    PoolType: number;
	    BannerID: number;
	    BannerName: string;
	    Outcome: string;
	    Annotation: string;
	    Name: string;
	    SearchNames: string[];
	    Icon: string;
	    Count: number;
	    Timestamp: number;
	    IsMissing: boolean;
	    TargetKnown: boolean;
	    AcquisitionCount: number;
	    PreviousLoss?: AcquisitionLoss;
	    CostBasis: string;
	
	    static createFrom(source: any = {}) {
	        return new DisplayRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.ItemID = source["ItemID"];
	        this.ItemType = source["ItemType"];
	        this.PoolType = source["PoolType"];
	        this.BannerID = source["BannerID"];
	        this.BannerName = source["BannerName"];
	        this.Outcome = source["Outcome"];
	        this.Annotation = source["Annotation"];
	        this.Name = source["Name"];
	        this.SearchNames = source["SearchNames"];
	        this.Icon = source["Icon"];
	        this.Count = source["Count"];
	        this.Timestamp = source["Timestamp"];
	        this.IsMissing = source["IsMissing"];
	        this.TargetKnown = source["TargetKnown"];
	        this.AcquisitionCount = source["AcquisitionCount"];
	        this.PreviousLoss = this.convertValues(source["PreviousLoss"], AcquisitionLoss);
	        this.CostBasis = source["CostBasis"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class DollWeaponHistory {
	    ItemID: number;
	    Name: string;
	    Records: DisplayRecord[];
	
	    static createFrom(source: any = {}) {
	        return new DollWeaponHistory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ItemID = source["ItemID"];
	        this.Name = source["Name"];
	        this.Records = this.convertValues(source["Records"], DisplayRecord);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class ImportPreview {
	    Ticket: string;
	    UID: string;
	    Server: string;
	    IdentitySource: string;
	
	    static createFrom(source: any = {}) {
	        return new ImportPreview(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Ticket = source["Ticket"];
	        this.UID = source["UID"];
	        this.Server = source["Server"];
	        this.IdentitySource = source["IdentitySource"];
	    }
	}
	export class OutfitRecord {
	    id: string;
	    itemId: number;
	    name: string;
	    searchNames: string[];
	    category: string;
	    accessNumber: number;
	    timestamp: number;
	
	    static createFrom(source: any = {}) {
	        return new OutfitRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.itemId = source["itemId"];
	        this.name = source["name"];
	        this.searchNames = source["searchNames"];
	        this.category = source["category"];
	        this.accessNumber = source["accessNumber"];
	        this.timestamp = source["timestamp"];
	    }
	}
	export class OutfitSeason {
	    id: number;
	    name: string;
	    startTime: number;
	    endTime: number;
	    status: string;
	    totalCount: number;
	    limitedCount: number;
	    rareCount: number;
	    normalCount: number;
	    unknownCount: number;
	    mainOutfit: string;
	    mainOutfitAt: number;
	    records: OutfitRecord[];
	
	    static createFrom(source: any = {}) {
	        return new OutfitSeason(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.startTime = source["startTime"];
	        this.endTime = source["endTime"];
	        this.status = source["status"];
	        this.totalCount = source["totalCount"];
	        this.limitedCount = source["limitedCount"];
	        this.rareCount = source["rareCount"];
	        this.normalCount = source["normalCount"];
	        this.unknownCount = source["unknownCount"];
	        this.mainOutfit = source["mainOutfit"];
	        this.mainOutfitAt = source["mainOutfitAt"];
	        this.records = this.convertValues(source["records"], OutfitRecord);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class OutfitInfo {
	    totalCount: number;
	    seasons: OutfitSeason[];
	
	    static createFrom(source: any = {}) {
	        return new OutfitInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalCount = source["totalCount"];
	        this.seasons = this.convertValues(source["seasons"], OutfitSeason);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	
	
	export class PoolInfo {
	    bannerPeriods: BannerPeriod[];
	    rank5Spent: number;
	    expiredPity: number;
	    pendingPity: number;
	    unknownRankCount: number;
	    targetAttempts: number;
	    targetWins: number;
	    unknownOutcomes: number;
	    storedCount: number;
	    recordList: DisplayRecord[];
	    totalCount: number;
	    rank5Count: number;
	    rank4Count: number;
	    rank3Count: number;
	    rank5Rate: number;
	    rank4Rate: number;
	    rank3Rate: number;
	    rank5Avg: number;
	    rank5UpAvg: number;
	    missingCount: number;
	    missingRate: number;
	
	    static createFrom(source: any = {}) {
	        return new PoolInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bannerPeriods = this.convertValues(source["bannerPeriods"], BannerPeriod);
	        this.rank5Spent = source["rank5Spent"];
	        this.expiredPity = source["expiredPity"];
	        this.pendingPity = source["pendingPity"];
	        this.unknownRankCount = source["unknownRankCount"];
	        this.targetAttempts = source["targetAttempts"];
	        this.targetWins = source["targetWins"];
	        this.unknownOutcomes = source["unknownOutcomes"];
	        this.storedCount = source["storedCount"];
	        this.recordList = this.convertValues(source["recordList"], DisplayRecord);
	        this.totalCount = source["totalCount"];
	        this.rank5Count = source["rank5Count"];
	        this.rank4Count = source["rank4Count"];
	        this.rank3Count = source["rank3Count"];
	        this.rank5Rate = source["rank5Rate"];
	        this.rank4Rate = source["rank4Rate"];
	        this.rank3Rate = source["rank3Rate"];
	        this.rank5Avg = source["rank5Avg"];
	        this.rank5UpAvg = source["rank5UpAvg"];
	        this.missingCount = source["missingCount"];
	        this.missingRate = source["missingRate"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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

export namespace model {
	
	export class User {
	    id: number;
	    uid: number;
	    server: string;
	    gameDataDir: string;
	    lastBBSToken: string;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.uid = source["uid"];
	        this.server = source["server"];
	        this.gameDataDir = source["gameDataDir"];
	        this.lastBBSToken = source["lastBBSToken"];
	    }
	}

}

export namespace service {
	
	export class CaptureSecurityStatus {
	    Active: boolean;
	    Certificate?: util.CaptureCertificate;
	    Hosts: string[];
	
	    static createFrom(source: any = {}) {
	        return new CaptureSecurityStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Active = source["Active"];
	        this.Certificate = this.convertValues(source["Certificate"], util.CaptureCertificate);
	        this.Hosts = source["Hosts"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class SyncDiff {
	    PoolType: number;
	    Count: number;
	
	    static createFrom(source: any = {}) {
	        return new SyncDiff(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.PoolType = source["PoolType"];
	        this.Count = source["Count"];
	    }
	}
	export class SyncResult {
	    Server: string;
	    Uid: number;
	    DiffList: SyncDiff[];
	
	    static createFrom(source: any = {}) {
	        return new SyncResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Server = source["Server"];
	        this.Uid = source["Uid"];
	        this.DiffList = this.convertValues(source["DiffList"], SyncDiff);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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

export namespace util {
	
	export class CaptureCertificate {
	    Directory: string;
	    Name: string;
	    Thumbprint: string;
	    Trusted: boolean;
	    Legacy: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CaptureCertificate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Directory = source["Directory"];
	        this.Name = source["Name"];
	        this.Thumbprint = source["Thumbprint"];
	        this.Trusted = source["Trusted"];
	        this.Legacy = source["Legacy"];
	    }
	}

}

