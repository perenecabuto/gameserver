/**
 * Player Entity
 */
game.PlayerEntity = me.Entity.extend({

  /**
   * constructor
   */
  init: function (name, x, y, settings) {
    settings = game.PlayerEntity.defaultSettings;
    settings.name = String(name);

    // call the constructor
    this._super(me.Entity, 'init', [x, y, settings]);

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
    this.action = null;
  },

  /**
   * update the entity
   */
  update: function (dt) {
    this.updatePosition();

    switch(this.action) {
        case GameMessage.Action.STOP:
          this.stop();
          break;
        case GameMessage.Action.MOVE_LEFT:
          this.moveLeft();
          break;
        case GameMessage.Action.MOVE_RIGHT:
          this.moveRight();
          break;
        case GameMessage.Action.JUMP:
          this.jump();
          break;
        case GameMessage.Action.DIE:
          this.die();
          break;

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
    if (typeof(other.name) == 'string' && other.name.match(/\d+/)) {
        this.stop();
        return false;
    }

    if (other.type == 'spike') {
      this.renderable.flicker(750);
      return false;
    }

    return true;
  },

  updatePosition: function () {
  },

  updateAction: function(action, pos) {
    this.action = action;
    this.pos.x = pos.x;
    this.pos.y = pos.y;
  },
  // player local
  moveLeft: function() {
    // flip the sprite on horizontal axis
    this.renderable.flipX(true);
    // update the entity velocity
    this.body.vel.x -= this.body.accel.x * me.timer.tick;
    // change to the walking animation
    if (!this.renderable.isCurrentAnimation("walk")) {
        this.renderable.setCurrentAnimation("walk");
    }
  },

  // player local
  moveRight: function() {
    // unflip the sprite
    this.renderable.flipX(false);
    // update the entity velocity
    this.body.vel.x += this.body.accel.x * me.timer.tick;
    // change to the walking animation
    if (!this.renderable.isCurrentAnimation("walk")) {
        this.renderable.setCurrentAnimation("walk");
    }
  },

  stop: function() {
    this.body.vel.x = 0;

    // change to the standing animation
    this.renderable.setCurrentAnimation("stand");
  },

  jump: function() {
    // make sure we are not already jumping or falling
    if (!this.body.jumping && !this.body.falling) {
        // set current vel to the maximum defined value
        // gravity will then do the rest
        this.body.vel.y = -this.body.maxVel.y * me.timer.tick;

        // set the jumping flag
        this.body.jumping = true;
    }
  },

  die: function() {
    me.game.world.removeChild(this);
  }

});


game.LocalPlayerEntity = game.PlayerEntity.extend({
    init: function (x, y, settings) {
        game.PlayerEntity.defaultSettings = settings;
        this._super(me.Entity, 'init', [x, y, settings]);

        this.body.setVelocity(3, 15);

        // define a basic walking animation (using all frames)
        this.renderable.addAnimation("walk",  [0, 1, 2, 3, 4, 5, 6, 7]);
        // define a standing animation (using the first frame)
        this.renderable.addAnimation("stand",  [0]);
        // set the standing animation as default
        this.renderable.setCurrentAnimation("stand");
        // set the display to follow our position on both axis
        me.game.viewport.follow(this.pos, me.game.viewport.AXIS.BOTH);
    },

    updatePosition: function () {
        var lastAction = this.action;

        if (me.input.isKeyPressed('left')) {
            this.action = GameMessage.Action.MOVE_LEFT;
        } else if (me.input.isKeyPressed('right')) {
            this.action = GameMessage.Action.MOVE_RIGHT;
        } else if (lastAction && this.action != GameMessage.Action.STOP) {
            this.action = GameMessage.Action.STOP;
        }

        if (me.input.isKeyPressed('jump')) {
            this.action = GameMessage.Action.JUMP;
        }

        if (this.action && this.action != lastAction) {
            me.event.publish("playerAction", [this.action, {x: this.pos.x, y: this.pos.y}]);
        }
    }
});
