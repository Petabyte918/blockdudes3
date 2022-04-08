import * as THREE from 'three';

import { Lighting } from './lighting.js'
import { particles } from './particles.js'
import { RenderMesh } from './render_mesh.js'
import { RenderObject } from './render_object.js'
import { RenderPlayer } from './render_player.js'
import { SceneComponent } from './scene_component.js'
import { GameUtil, LogUtil, Util } from './util.js'
import { Weather } from './weather.js'

export class SceneMap {
	private _scene : THREE.Scene;
	private _lighting : Lighting;
	private _weather : Weather;
	private _renders : Map<number, Map<number, RenderObject>>;
	private _deleted : Map<number, Set<number>>;

	private _components : Array<SceneComponent>;

	constructor() {
		this.reset();
	}

	scene() : THREE.Scene { return this._scene; }

	reset() : void {
		this._scene = new THREE.Scene();
		this._renders = new Map();
		this._deleted = new Map();

		this._components = new Array<SceneComponent>();
		this.addComponent(new Lighting());
		this.addComponent(new Weather());
		this.addComponent(particles);
	}

	addComponent(component : SceneComponent) {
		this._components.push(component);
		this._scene.add(component.scene());
	}

	updateComponents(position : THREE.Vector3) : void {
		this._components.forEach((component) => {
			component.update(position);
		});
	}

	add(space : number, id : number, object : RenderObject) : void {
		const map = this.getMap(space);
		if (map.has(id)) {
			LogUtil.d("Overwriting object space " + space + ", id " + id + "!");
		}
		map.set(id, object);
		object.onMeshLoad(() => {
			if (this.has(space, id)) {
				this._scene.add(object.mesh());
			}
		});
	}

	has(space : number, id : number) : boolean {
		const map = this.getMap(space);
		return map.has(id) && Util.defined(map.get(id));	
	}

	deleted(space : number, id : number) : boolean {
		if (!this._deleted.has(space)) {
			return false;
		}
		return this._deleted.get(space).has(id);
	}

	get(space : number, id : number) : RenderObject {
		const map = this.getMap(space);
		return map.get(id);
	}

	delete(space : number, id : number) : void {
		const map = this.getMap(space);
		if (map.has(id)) {
			this._scene.remove(map.get(id).mesh());
			wasmDelete(space, id);
			map.delete(id);

			if (!this._deleted.has(space)) {
				this._deleted.set(space, new Set());
			}
			this._deleted.get(space).add(id);
		}
	}

	clear(space : number) : void {
		const map = this.getMap(space);
		map.forEach((object, id) => {
			this.delete(space, id);
		});
		map.clear();
	}

	clearObjects() : void {
		this._renders.forEach((render, space) => {
			if (space != playerSpace) {
				this.clear(space);
			}
		});
	}

	update(space : number, id : number, msg : Map<number, any>, seqNum?: number) : void {
		const map = this.getMap(space);
		const object = map.get(id);

		if (!Util.defined(object)) {
			this.delete(space, id);
			return;
		}

		object.update(msg, seqNum);

		if (!object.initialized() && object.ready()) {
			object.initialize();
			wasmAdd(space, id, object.data());
		}

		if (wasmHas(space, id)) {
			wasmSetData(space, id, object.data());
		}

		if (object.deleted()) {
			this.delete(space, id);
		}
	}

	private getMap(space : number) : Map<number, RenderObject> {
		if (!this._renders.has(space)) {
			this._renders.set(space, new Map<number, RenderObject>());
		}
		return this._renders.get(space);
	}
}