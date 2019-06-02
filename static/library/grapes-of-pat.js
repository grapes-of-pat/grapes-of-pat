(function(global) {

    var serverHost = global.GRAPES_HOST || "grapesofpat.com";
    var clientHost = window.location.host;

    global.grapesOfPat = {
        createServer: createServer,
        connect: connect
    }

    /**
     *
     * @param {*} options
     * @returns A promise of the connection details
     */
    // TODO Connect to official server
    function createServer(options){
        var clientIds = options.clientIds || ['Client 1'];
        var controller = options.controller || 'boolean-controller';
        var protocol = window.location.protocol;
        var wsProtocol = protocol === 'https:' ? 'wss:' : 'ws:';
        var controllerUrl = protocol + '//' + serverHost + '/controller/' + controller + '.html';

        return fetch('//' + serverHost + '/session/')
            .then(function(response) {
                return response.text();
            })
            .then(function(sessionId) {
                // FIXME wss
                var socketUrl = wsProtocol +'//' + serverHost + '/session/' + sessionId + '/start';
                var socket = new WebSocket(socketUrl);
                socket.onmessage = function(msg) {
                    var data = JSON.parse(msg.data);
                    if (data.type === "ping") {
                        console.log("Returning ping to " + data.clientId);
                        // Just send client id as text;
                        // TODO Allow the sending of more complex data
                        socket.send(data.clientId);
                    }
                    options.onmessage(data);
                };
                return {
                    urls: clientIds.map(function(id) {
                        return controllerUrl + "#sessionId=" + sessionId + "&clientId=" + id;
                    })
                }
            });
    }

    // TODO Connect to official server
    function connect(options) {
        var hash = window.location.hash
            .slice(1)
            .split("&")
            .reduce(function(hash, part) {
                var kv = part.split("=");
                hash[kv[0]] = kv[1];
                return hash;
            }, {});
        var sessionId = hash.sessionId;
        var clientId = hash.clientId;
        var protocol = window.location.protocol;
        var wsProtocol = protocol === 'https:' ? 'wss:' : 'ws:';
        var socketUrl = wsProtocol + '//' + clientHost + '/session/'  + sessionId
            + '/connect?clientID=' + clientId ;
        var socket = new WebSocket(socketUrl);

        return new Promise(function(resolve, reject) {
            socket.onopen = function() {
                send({ 'type': 'connected' });
                resolve({
                    send: send
                });
            };
            if (options && options.onmessage) {
                socket.onmessage = options.onmessage;
            }
            function send(msg) {
                msg.clientId = clientId;
                socket.send(JSON.stringify(msg));
            }
        });
    }

})(this);