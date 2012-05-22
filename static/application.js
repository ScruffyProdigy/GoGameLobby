c=new WebSocket('ws://localhost:3000');
c.onmessage=function(msg){alert(msg)};