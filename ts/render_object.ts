class RenderObject {
	protected readonly _debugMaterial = new THREE.MeshStandardMaterial({color: 0xff0000, wireframe: true });

	protected _mesh : any;

	protected _mixer : any;
	protected _lastMixerUpdate : number;

	protected _actions : Map<any, any>;
	protected _activeActions : Set<any>;
	
	constructor(mesh : any) {
		this._mesh = mesh;
		this._activeActions = new Set();

		this._lastMixerUpdate = Date.now();
	}

	update(msg : any) : void {
		if (msg.hasOwnProperty(posProp)) {
			const pos = msg[posProp]
			this._mesh.position.x = pos.X;
			this._mesh.position.y = pos.Y;
		}
	}

	mesh() : any {
		return this._mesh;
	}

	protected updateMixer() {
		const now = Date.now();
		const seconds = (now - this._lastMixerUpdate) / 1000;
		this._mixer.update(seconds)
		this._lastMixerUpdate = now;
	}

	protected setWeight(action : any, weight : number) : void {
		action.enabled = true;
		action.setEffectiveTimeScale(1);
		action.setEffectiveWeight(weight);
	}

	protected fadeIn(action : any, duration : number) : void {
		if (!this._activeActions.has(action)) {
			this._actions.get(action).reset();
			this._actions.get(action).fadeIn(duration);
			this._activeActions.add(action);
		}
	}

	protected fadeOut(action : any, duration : number) : void {
		if (this._activeActions.has(action)) {
			this._actions.get(action).reset();
			this._actions.get(action).fadeOut(duration);
			this._activeActions.delete(action);
		}
	}

	protected fadeTo(startAction : any, endAction : any, duration : number) : void {
		this.fadeOut(startAction, duration);
		this.fadeIn(endAction, duration);
	}
}