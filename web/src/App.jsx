import React, { useState, useEffect, useRef } from 'react';
import './App.css';

function App() {
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const ws = useRef(null);

  useEffect(() => {
    // Connect to WebSocket server
    ws.current = new WebSocket('ws://localhost:8080/ws');
    ws.current.onopen = () => console.log('ws opened');
    ws.current.onclose = () => console.log('ws closed');

    const wsCurrent = ws.current;

    // Listen for messages
    wsCurrent.onmessage = (e) => {
      const message = JSON.parse(e.data);
      setMessages((prevMessages) => [...prevMessages, message]);
    };

    return () => {
      wsCurrent.close();
    };
  }, []);

  const handleSendMessage = () => {
    if (newMessage !== '') {
      // Send message to WebSocket server
      ws.current.send(JSON.stringify({ username: 'User', message: newMessage }));
      setNewMessage('');
    }
  };

  const handleNewMessageChange = (event) => {
    setNewMessage(event.target.value);
  };

  return (
    <div className="App">
      <header className="App-header">
        <h2>Global Chat</h2>
      </header>
      <div className="messages">
        {messages.map((message, index) => (
          <p key={index}><strong>{message.username}:</strong> {message.message}</p>
        ))}
      </div>
      <textarea
        value={newMessage}
        onChange={handleNewMessageChange}
        placeholder="Write your message..."
        className="new-message-input-field"
      />
      <button onClick={handleSendMessage} className="send-message-button">
        Send
      </button>
    </div>
  );
}

export default App;