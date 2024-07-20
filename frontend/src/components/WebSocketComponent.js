import React, { useEffect, useState } from 'react';

const WebSocketComponent = () => {
    const [message, setMessage] = useState('');
    const [ws, setWs] = useState(null);

    useEffect(() => {
        const socket = new WebSocket('ws://localhost:8080/ws');
        setWs(socket);

        socket.onopen = () => {
            console.log('Connected to the WebSocket server');
        };

        socket.onmessage = (event) => {
            console.log('Message from server:', event.data);
            setMessage(event.data);
        };

        socket.onclose = () => {
            console.log('Disconnected from the WebSocket server');
        };

        return () => {
            socket.close();
        };
    }, []);

    const sendMessage = () => {
        if (ws) {
            ws.send('Hello Server!');
        }
    };

    return (
        <div>
            <h1>WebSocket Client</h1>
            <button onClick={sendMessage}>Send Message</button>
            <p>Message from server: {message}</p>
        </div>
    );
};

export default WebSocketComponent;
