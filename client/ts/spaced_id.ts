
export class SpacedId {
	private _space : number;
	private _id : number;

	constructor(space : number, id : number) {
		this._space = space;
		this._id = id;
	}

	space() : number {
		return this._space;
	}

	id() : number {
		return this._id;
	}

	valid() : boolean {
		return this._space > 0 && this._id >= 0;
	}

	toString() : string {
		return this._space + "," + this._id;
	}

	parseSpace(sid : string) : number {
		return Number(sid.split(",")[0]);	
	}

	parseId(sid : string) : number {
		return Number(sid.split(",")[1])	
	}
}
