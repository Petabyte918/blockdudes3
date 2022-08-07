import * as THREE from 'three';
import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader.js';

import { LogUtil, Util } from './util.js'

import { options } from './options.js'

export enum Model {
	UNKNOWN = "",
	CHICKEN = "chicken",
	DUCK = "duck",

	UZI = "uzi",
	BAZOOKA = "bazooka",
	SNIPER = "sniper",
	STAR_GUN = "star_gun",

	BOLT = "bolt",
	ROCKET = "rocket",
}

class Loader {
	private readonly _modelPrefix = "./model/";

	private _loader : GLTFLoader;
	private _cache : Map<Model, any>;
	private _paths : Map<Model, string>;

	constructor() {
		this._loader = new GLTFLoader();
		this._cache = new Map<Model, any>();

		this._paths = new Map<Model, string>();

		for (const model in Model) {
			if (model.length === 0) {
				continue;
			}

			this._paths.set(Model[model], this._modelPrefix + model + ".glb");
		}
	}

	preload(models : Array<Model>) : void {
		models.forEach((model) => {
			this.load(model, (mesh) => {}); 
		});
	}

	load(model : Model, cb : (any) => void) : void {
		if (!this._paths.has(model) || !Util.defined(this._paths.get(model))) {
			LogUtil.d("Tried to cache unknown model " + model);
			return;
		}

		this._loader.load(this._paths.get(model), (data) => {
			this.process(model, data);
			cb(data.scene);
		});
	}

	getWeaponModel(weaponType : number) : Model {
		switch (weaponType) {
			case uziWeapon:
				return Model.UZI;
			case bazookaWeapon:
				return Model.BAZOOKA;
			case sniperWeapon:
				return Model.SNIPER;
			case starWeapon:
				return Model.STAR_GUN;
		}
		return Model.UNKNOWN;
	}

	private process(model : Model, data : any) : void {
		switch (model) {
			case Model.CHICKEN:
			case Model.DUCK:
				data.scene.animations = data.animations;
				if (options.enableShadows) {
					data.scene.getObjectByName("mesh").castShadow = true;
					data.scene.getObjectByName("mesh").receiveShadow = true;
				}
				break;
			case Model.UZI:
			case Model.BAZOOKA:
			case Model.SNIPER:
			case Model.STAR_GUN:
			case Model.ROCKET:
				if (options.enableShadows) {
					data.scene.getObjectByName("mesh").castShadow = true;
					data.scene.getObjectByName("mesh").receiveShadow = true;
				}
				break;
			case Model.BOLT:
				data.scene.getObjectByName("mesh").castShadow = false;
				data.scene.getObjectByName("mesh").receiveShadow = false;
				break;
			default:
				LogUtil.d("Model " + model + " processing skipped.");
				break;
		}
	}
}

export const loader = new Loader();