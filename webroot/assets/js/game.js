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
        this.stage.disableVisibilityChange = true;
        this.playerPool = this.buildPlayerPool();
    },

    update: function() {
        var that = this;
        this.playerPool.forEach(function(player) {
            // TODO send velocity and position to have a minimum prediction and correction
            player.jumping = player.y + player.body.height / 2 < that.world.height;
            if (!player.jumping && player.action == GameMessage.Action.JUMP) {
                player.body.velocity.y = -400;
            }

            switch(player.action) {
                case GameMessage.Action.STOP:
                    if (!player.jumping) player.reset(player.x, player.y);
                    break;
                case GameMessage.Action.MOVE_LEFT:
                    player.body.velocity.x = -STEP_SIZE;
                    player.play('walk-left');
                    break;
                case GameMessage.Action.MOVE_RIGHT:
                    player.body.velocity.x = STEP_SIZE;
                    player.play('walk-right');
                    break;
                case GameMessage.Action.DIE:
                    player.kill();
                    break;
            }
        });
    },

    render: function () {
        //this.game.debug.spriteInfo(this.player, 20, 32);
    },

    buildPlayerPool: function() {
        var playerPool = this.add.group();
        playerPool.enableBody = true;
        playerPool.physicsBodyType = Phaser.Physics.ARCADE;
        playerPool.createMultiple(50, 'boy');
        playerPool.setAll('anchor.x', 0.5);
        playerPool.setAll('anchor.y', 0.5);
        playerPool.setAll('body.collideWorldBounds', true);
        playerPool.setAll('body.collideWorldBounds', true);
        playerPool.forEach(function (player) {
            player.animations.add('walk-left', [4, 5, 6, 7], 20, true);
            player.animations.add('walk-right', [8, 9, 10, 11], 20, true);
        });

        return playerPool;
    }
};
