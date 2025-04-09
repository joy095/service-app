// Connect to your WebSocket server
@ws = ws://localhost:8085/ws

### Connect
WEBSOCKET {{ws}}

### Send message to server
{
  "from": "VS Code",
  "message": "Hello from the WebSocket client!"
}
