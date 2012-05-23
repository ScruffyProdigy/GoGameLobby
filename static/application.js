$(function(){
	var isFocused = true;
	var focusActions = new Array();
	$(window).focus(function(){
		isFocused = true;
		$("body").css("background-color","#ffffff")
		while(focusActions.length > 0) {
			var action = focusActions.shift();
			action()
		}
	}).blur(function(){	
		$("body").css("background-color","#eeeeee")
		isFocused = false;
	});
	
	var onGainFocus = function(action) {
		focusActions.push(action);
	}
	
	var startGame = function(c,gameInfo) {
		var message = new Object()
		message.type = "start"
		message.data = gameInfo.start
		c.send(JSON.stringify(message));
	}
	
	c=new WebSocket('ws://localhost:3000');
	c.onmessage=function(msg){
		msg = JSON.parse(msg.data);
		if(msg.start) {
			//remove queues from list
			$(".queue-"+msg.start.game+msg.start.mode).remove();
			if(isFocused) {
				//start the game
				startGame(c,msg);
			} else {
				//set the callback once we regain focus
				onGainFocus(function(){
					//start the game
					startGame(c,msg);
				});
				//and in the meantime, alert the user that their queue has started
				alert("Your game has started");
			}
		}
		if(msg.startloc) {
			alert("New Window!")
			//open a new window with the start location
			window.open(msg.startloc)
		}
	};
});