import * as THREE from 'three'

import { options } from './options.js'
import { renderer } from './renderer.js'
import { RingBuffer } from './ring_buffer.js'

export class LightBuffer {

	private readonly _bufferSize : number = 10;

	private _spotLights : RingBuffer<THREE.SpotLight>;
	private _pointLights : RingBuffer<THREE.PointLight>;

	constructor(scene : THREE.Scene) {
		this._spotLights = new RingBuffer<THREE.SpotLight>();
		this._pointLights = new RingBuffer<THREE.PointLight>();

		for (let i = 0; i < this._bufferSize; ++i) {
			let spotLight = new THREE.SpotLight(0xFFFFFF, 1, 1, 1);
			this._spotLights.push(spotLight);
			scene.add(spotLight);

			let pointLight = new THREE.PointLight(0xFFFFFF, 1, 1);
			this._pointLights.push(pointLight);
			scene.add(pointLight);
		}
		renderer.compile(scene);

		this._spotLights.asArray().forEach((light) => {
			light.intensity = 0;
		});
		this._pointLights.asArray().forEach((light) => {
			light.intensity = 0;
		})
	}

	getSpotLight() : THREE.SpotLight {
		if (!options.enableDynamicLighting) {
			return null;
		}
		return this._spotLights.getNext();
	}

	getPointLight() : THREE.PointLight {
		if (!options.enableDynamicLighting) {
			return null;
		}
		return this._pointLights.getNext();
	}

}