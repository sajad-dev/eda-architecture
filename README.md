# Go Project Similar to Pusher with EDA Architecture 🚀
This project is written in Go and acts similarly to the Pusher service. The project is designed using Event-Driven Architecture (EDA), which allows you to manage different systems indirectly and scalably using events. 🎯

## Features ✨

- Similar to Pusher: For sending and receiving events in real-time. 🔄


-   EDA Architecture: The system is designed based on event-driven architecture. ⚡

-   Scalable: Designed for high performance and scalability. 📈
-   API Keys: Secure connection to WebSocket using public key and secret key. 🔑

## API Endpoints 🛠️
#### `POST /api/create-key`
This endpoint provides the public key and secret key.
#### Response:
```json
{
  "public_key": "your_public_key",
  "secret_key": "your_secret_key"
}
```

## Pusher Connection 🔌

When connecting to Pusher, you need to send the public key in the route and the secret key in the header.

#### Pusher Connection:

-   Route: Send the public key in the route when establishing the Pusher connection. For example, the Pusher connection URL will look like this:
```bash
ws://ip:port/public_key
```

### Connecting to Pusher WebSocket in Frontend (JavaScript) 🌍
```js
// 🚀 Import the Pusher library
import Pusher from 'pusher-js';

// 🔑 Initialize Pusher WebSocket
const pusher = new Pusher('PUBLIC_KEY', {  // Replace with your public key
    wsHost: IP,  // Running on localhost
    wsPort: PORT,  // Port where WebSocket is running
});

// 🎯 Subscribe to a specific channel
const channel = pusher.subscribe(CHANNEL_NAME);

// 👂 Listen for a specific event
channel.bind(EVENT_NAME, function(data) {
    console.log('📩 Data received:', data);
})
```

**🚨 Important: Ensure that the secret key is only sent in requests where it's necessary. For example, it should only be included in the Pusher connection request.**


## Setup 🛠️
To set up the project, you first need to create the migrations using the following command:

```bash
go run cmd/main.go migration create
```
After the migrations are created, you can run the project with the following command:

```bash
go run cmd/main.go
```

## Usage 📚
For sample code and further details, you can visit the following repository:

https://github.com/sajad-dev/eda-architecture-sample 🔗

https://github.com/sajad-dev/eda-architecture-sample-laravel 🔗

## Author 🖋️

Mohammad Sajad Poorajam / محمد سجاد پورعجم
