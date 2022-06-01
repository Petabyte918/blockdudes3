import * as THREE from 'three';
import { Sound } from './audio.js';
import { game } from './game.js';
import { Model, loader } from './loader.js';
import { SceneComponentType } from './scene_component.js';
import { RenderCustom } from './render_custom.js';
import { RenderProjectile } from './render_projectile.js';
import { MathUtil } from './util.js';
export class RenderRocket extends RenderProjectile {
    constructor(space, id) {
        super(space, id);
        this._smokeMaterial = new THREE.MeshStandardMaterial({ color: 0xbbbbbb, transparent: true, opacity: 0.5 });
        this._smokeInterval = 15;
        this._rotateZ = 0.12;
        this._lastSmoke = 0;
        this.setSound(Sound.ROCKET);
    }
    initialize() {
        super.initialize();
        loader.load(Model.ROCKET, (mesh) => {
            this.setMesh(mesh);
        });
    }
    setMesh(mesh) {
        super.setMesh(mesh);
        mesh.rotation.y = Math.PI / 2;
    }
    update(msg, seqNum) {
        super.update(msg, seqNum);
        if (!this.hasMesh()) {
            return;
        }
        const pos = this.pos();
        const vel = this.vel();
        const dim = this.dim();
        const dir = vel.clone().normalize();
        const angle = vel.angle() * -1;
        const projectile = this.mesh().getObjectByName("mesh");
        projectile.rotation.x = angle;
        projectile.rotateZ(this._rotateZ);
        if (Date.now() - this._lastSmoke >= this._smokeInterval) {
            const smokeMesh = new THREE.Mesh(new THREE.SphereGeometry(MathUtil.randomRange(0.1, 0.2), 3, 3), this._smokeMaterial);
            smokeMesh.position.x = pos.x - dim.x / 2 * dir.x + MathUtil.randomRange(-0.1, 0.1);
            smokeMesh.position.y = pos.y - dim.y / 2 * dir.y + MathUtil.randomRange(-0.1, 0.1);
            smokeMesh.position.z = this.mesh().position.z + MathUtil.randomRange(-0.1, 0.1);
            const smoke = new RenderCustom();
            smoke.setMesh(smokeMesh);
            smoke.setUpdate(() => {
                smoke.mesh().scale.multiplyScalar(0.9);
            });
            game.sceneComponent(SceneComponentType.PARTICLES).addCustomTemp(smoke, 400);
            this._lastSmoke = Date.now();
        }
    }
}
