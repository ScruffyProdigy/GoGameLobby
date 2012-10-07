function sendMessage(type, data) {
	var message = new Object();
	message.type = type;
	message.data = data;
	conn.send(JSON.stringify(message));
}

$(function(){
	var base = window.location.host;
	var path = window.location.pathname
	var isFocused = true;
	var focusActions = new Array();
	$(window).focus(function(){
		isFocused = true;
		$("body").css("background-color","#ffffff");
		while(focusActions.length > 0) {
			var action = focusActions.shift();
			action();
		}
	}).blur(function(){	
		$("body").css("background-color","#eeeeee");
		isFocused = false;
	});
	
	function onGainFocus(action) {
		focusActions.push(action);
	}
	
	function startGame(gameInfo) {
		sendMessage("start",gameInfo)
	}
	
	conn=new WebSocket('ws://'+base+path);
	conn.onmessage = function(msg){
		msg = JSON.parse(msg.data);
		var type = msg.Type;
		msg = msg.Data;
		switch(type) {
			case "start":
				//remove queues from list
				$(".queue-"+msg.game+msg.mode).remove();
				alert("Your game has started");
			break;
			case "roomchat":
				$(".roomchat .chatbox").prepend("<p><a href='\\users\\"+msg.Name+"'>"+msg.Name+"</a>  "+msg.Text+"</p>")
			break;
		}
	};
	
	
	$(".roomchat .chatcontrol .submit").click(function() {
		var message = new Object();
		message.url = window.location.pathname;
		message.text = $(".roomchat .chatcontrol .message").val();
		sendMessage("roomchat",message);
		
		$(".roomchat .chatcontrol .message").val("");
		return false
	});
});