/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2014 Ivanix Mobile LLC
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 */

/*
 * Plugin: SaveCPU
 * Author: Ivanix @ Ivanix Mobile LLC
 * Purpose:  Reduce CPU usage caused from redudant rendering
 *           on idle or static display scenes
 *           reduce fps for casual/puzzle games
 *
 *
 * Configurable Properties:
 *                
 * [renderOnFPS]   
 * Constrains maximum FPS to value set. 
 * Reasonable values from 0 to 60 
 * Default value 30
 * Set value to 0 disable rendering based on FPS
 * and use methods described below.
 *
 * [renderOnPointerChange]   
 * Render when pointer movement detected.
 * Possible values  "true" or "false"
 * Default: false
 * Note that renderOnFPS must be set to 0
 *
 *
 * Callable Methods:
 * 
 * [forceRender()]
 * Forces rendering during core game loop
 * Can be called independently or in tandem with above properties.
 * Should be called inside update function.
 * @class Phaser.Plugin.SaveCPU
 */

/*global
    Phaser: true,
    window: true
 */
/*jslint nomen: true */

Phaser.Plugin.SaveCPU = function (game, parent) {
    'use strict';

    Phaser.Plugin.call(this, game, parent);

};
Phaser.Plugin.SaveCPU.prototype = Object.create(Phaser.Plugin.prototype);
Phaser.Plugin.SaveCPU.constructor = Phaser.Plugin.SaveCPU;

Phaser.Plugin.SaveCPU.prototype.init = function () {
    'use strict';
    this.now = window.performance.now();
    this.renderType = this.game.renderType;

    // fps default
    this.renderOnFPS = 30;
    this.renderOnPointerChange = false;
    this.renderDirty = true;
};
Phaser.Plugin.SaveCPU.prototype.setRender = function () {
    'use strict';
    if (this.renderDirty) {
        this.game.renderType = this.renderType;
    } else {
        this.game.renderType = Phaser.HEADLESS;
    }
    this.renderDirty = false;
};
Phaser.Plugin.SaveCPU.prototype.forceRender = function () {
    'use strict';
    this.renderDirty = true;
};
Phaser.Plugin.SaveCPU.prototype.forceRenderOnPointerChange = function () {
    'use strict';
    
    var input = this.game.input;

    if (input.oldx !== input.x || input.oldy !== input.y) {
        this.forceRender();
        input.oldx = input.x;
        input.oldy = input.y;
    }
    if (input.oldDown !== input.activePointer.isDown) {
        this.forceRender();
        input.oldDown = input.activePointer.isDown;
    }
};
Phaser.Plugin.SaveCPU.prototype.forceRenderOnFPS = function () {
    'use strict';
    
    var ts, diff;

    ts = window.performance.now();
    diff = ts - this.now;
    if (diff < (1000 / this.renderOnFPS)) {
        return false;
    }
    this.now = ts;
    this.forceRender();
    return true;

};
Phaser.Plugin.SaveCPU.prototype.postUpdate = function () {
    'use strict';
    if (this.renderOnFPS && this.forceRenderOnFPS()) {
        this.setRender();
        return;
    }
    if (this.renderOnPointerChange && this.forceRenderOnPointerChange()) {
        this.setRender();
        return;
    }
    this.setRender();
};
Phaser.Plugin.SaveCPU.prototype.postRender = function () {
    'use strict';
    if (this.game._paused) {
        this.game.renderType = Phaser.HEADLESS;
    }
};
