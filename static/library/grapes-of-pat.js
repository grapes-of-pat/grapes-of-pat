(function(global) {

    var host = window.location.host;

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
        var controllerUrl = protocol + '//' + host + '/controller/' + controller + '.html';

        return fetch('//' + host + '/session/')
            .then(function(response) {
                return response.text();
            })
            .then(function(sessionId) {
                // FIXME wss
                var socketUrl = wsProtocol +'//' + host + '/session/' + sessionId + '/start';
                var socket = new WebSocket(socketUrl);
                socket.onmessage = function(msg) {
                    options.onmessage(JSON.parse(msg.data));
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
        var socketUrl = wsProtocol + '//' + host + '/session/'  + sessionId + '/connect';
        var socket = new WebSocket(socketUrl);

        return new Promise(function(resolve, reject) {
            socket.onopen = function() {
                send({ 'type': 'connected' });
                resolve({
                    send: send
                });
            };
            function send(msg) {
                msg.clientId = clientId;
                socket.send(JSON.stringify(msg));
            }
        });
    }

})(this);