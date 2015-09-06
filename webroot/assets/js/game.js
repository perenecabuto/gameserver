var STEP_SIZE = 100;
var Game = function() {};

Game.prototype = {
    constructor: Game,

    preload: function () {
        this.game.load.spritesheet('boy', 'assets/sprites/boy.png', 32, 32);
    },

    create: function () {
        this.game.physics.startSystem(Phaser.Physics.ARCADE);
        this.game.physics.arcade.gravity.y = 1200;

        this.player = this.game.add.sprite(this.game.world.centerX, this.game.world.height - 16, 'boy');
        this.player.anchor.setTo(0.5, 0.5);
        this.player.animations.add('walk-left', [4, 5, 6, 7], 20, true);
        this.player.animations.add('walk-right', [8, 9, 10, 11], 20, true);

        this.physics.enable(this.player, Phaser.Physics.ARCADE);
        this.player.body.collideWorldBounds = true;

        var that = this;
        this.game.input.keyboard.onUpCallback = function() {
            var actionEvent = new CustomEvent("player-action", {detail: {
                action: GameMessage.Action.STOP,
                pos: {x: that.player.body.x, y: that.player.body.y}
            }});
            document.dispatchEvent(actionEvent);
        };
    },

    update: function() {
        // TODO send velocity and position to have a minimum prediction and correction
        this.player.jumping = this.player.y + this.player.body.height / 2 < this.world.height;
        if (!this.player.jumping && this.action == GameMessage.Action.JUMP) {
            this.player.body.velocity.y = -400;
        }

        switch(this.action) {
            case GameMessage.Action.STOP:
              this.player.reset(this.player.x, this.player.y);
              break;
            case GameMessage.Action.MOVE_LEFT:
              this.player.body.velocity.x = -STEP_SIZE;
              this.player.play('walk-left');
              break;
            case GameMessage.Action.MOVE_RIGHT:
              this.player.body.velocity.x = STEP_SIZE;
              this.player.play('walk-right');
              break;
            case GameMessage.Action.DIE:
              this.player.die();
              break;

        }

        if (this.game.input.keyboard.isDown(Phaser.Keyboard.SPACEBAR)) {
            this.action = GameMessage.Action.JUMP;
        } else if (this.game.input.keyboard.isDown(Phaser.Keyboard.LEFT)) {
            this.action = GameMessage.Action.MOVE_LEFT;
        } else if (this.game.input.keyboard.isDown(Phaser.Keyboard.RIGHT)) {
            this.action = GameMessage.Action.MOVE_RIGHT;
        } else if (!this.player.jumping) {
            this.action = GameMessage.Action.STOP;
        }
    },

    render: function () {
        //this.game.debug.spriteInfo(this.player, 20, 32);
    }
};
