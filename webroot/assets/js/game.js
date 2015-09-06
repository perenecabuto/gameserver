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
        this.onComplete();
    },

    update: function() {
        // TODO send velocity and position to have a minimum prediction and correction
        var jumping = this.player.y + this.player.body.height / 2 < this.world.height;
        var old_inaction = Boolean(this.player.inaction);

        if (!jumping && this.game.input.keyboard.isDown(Phaser.Keyboard.SPACEBAR)) {
            this.player.inaction = true;
            jumping = true;
            this.player.body.velocity.y = -400;
        }

        if (this.game.input.keyboard.isDown(Phaser.Keyboard.LEFT)) {
            this.player.inaction = true;
            this.player.body.velocity.x = -STEP_SIZE;
            this.player.play('walk-left');
        } else if (this.game.input.keyboard.isDown(Phaser.Keyboard.RIGHT)) {
            this.player.inaction = true;
            this.player.body.velocity.x = STEP_SIZE;
            this.player.play('walk-right');
        } else if (!jumping) {
            this.player.inaction = false;
            this.player.reset(this.player.x, this.player.y);
        }

        if (old_inaction && !this.player.inaction) {
            $(document).trigger("player-action", [GameMessage.Action.STOP, {x: this.player.body.x, y: this.player.body.y}]);
        }
    },

    render: function () {
        this.game.debug.spriteInfo(this.player, 20, 32);
    },

    onComplete: function(callback) {
        this._onComplete = this._onComplete || function() {};

        if (callback) {
            this._onComplete = callback;
        } else {
            this._onComplete(this);
        }

        return this;
    }
};
