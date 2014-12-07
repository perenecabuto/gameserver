"use strict";

var ProtoBuf = dcodeIO.ProtoBuf;

ProtoBuf.loadProtoFile("protobuf/object.proto", function(err, builder) {
    var Obj = builder.build('protobuf.Object');
    var Position = Obj.Position;
    var Status = Obj.Status;

    var player = new Obj(1, new Position(0, 0), Status.IDLE);

    console.log(player, player.encode(), player.toArrayBuffer());

    player.position.x = 10;
    player.position.y = 20;

    var decodedPlayer = Obj.decode(player.toArrayBuffer());
    console.log(decodedPlayer, decodedPlayer.encode(), decodedPlayer.toArrayBuffer());
});
