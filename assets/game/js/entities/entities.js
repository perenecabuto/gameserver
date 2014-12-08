/**
 * Player Entity
 */
game.PlayerEntity = me.Entity.extend({

  /**
   * constructor
   */
  init:function (x, y, settings) {
    // call the constructor
    this._super(me.Entity, 'init', [x, y , settings]);

    // set the display to follow our position on both axis
    me.game.viewport.follow(this.pos, me.game.viewport.AXIS.BOTH);
    // ensure the player is updated even when outside of the viewport
    this.alwaysUpdate = true;

    // set the default horizontal & vertical speed (accel vector)
    this.body.setVelocity(3, 15);

    // define a basic walking animation (using all frames)
    this.renderable.addAnimation("walk",  [0, 1, 2, 3, 4, 5, 6, 7]);
    // define a standing animation (using the first frame)
    this.renderable.addAnimation("stand",  [0]);
    // set the standing animation as default
    this.renderable.setCurrentAnimation("stand");
  },

  /**
   * update the entity
   */
  update : function (dt) {

        /* OBS
         * Acho que essa parte pode ser externa, tanto para o player local quanto para o remoto
         * Ou seja a implementação será a mesma
         * somente faz body.update e collision check
         */

        /* OBS this.flipX trocado por this.renderable.flipX
         * procure por flipX em http://blog.ciangames.com/2014/11/upgrading-to-melonjs-20.html
         */
        if (me.input.isKeyPressed('left')) {
            // flip the sprite on horizontal axis
            this.renderable.flipX(true);
            // update the entity velocity
            this.body.vel.x -= this.body.accel.x * me.timer.tick;
            // change to the walking animation
            if (!this.renderable.isCurrentAnimation("walk")) {
                this.renderable.setCurrentAnimation("walk");
            }
        } else if (me.input.isKeyPressed('right')) {
            // unflip the sprite
            this.renderable.flipX(false);
            // update the entity velocity
            this.body.vel.x += this.body.accel.x * me.timer.tick;
            // change to the walking animation
            if (!this.renderable.isCurrentAnimation("walk")) {
                this.renderable.setCurrentAnimation("walk");
            }
        } else {
            this.body.vel.x = 0;
            // change to the standing animation
            this.renderable.setCurrentAnimation("stand");
        }

        if (me.input.isKeyPressed('jump')) {
            // make sure we are not already jumping or falling
            if (!this.body.jumping && !this.body.falling) {
                // set current vel to the maximum defined value
                // gravity will then do the rest
                this.body.vel.y = -this.body.maxVel.y * me.timer.tick;
                // set the jumping flag
                this.body.jumping = true;
            }
        }


    // apply physics to the body (this moves the entity)
    this.body.update(dt);
 
    // handle collisions against other shapes
    me.collision.check(this);

    // return true if we moved or if the renderable was updated
    return (this._super(me.Entity, 'update', [dt]) || this.body.vel.x !== 0 || this.body.vel.y !== 0);
  },

    /**
     * colision handler
     * (called when colliding with other objects)
     */
  onCollision : function (response, other) {
    // Make all other objects solid
    return true;
  }
});
