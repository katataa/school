import React, { useState, useEffect, useRef } from 'react';
import { useParams } from 'react-router-dom';
import ChatList from './ChatList';
import axios from 'axios';
import '../styles/chats.css';

function Chats({ chats, setChats, chatListWS }) {
    const { userId } = useParams();
    const [messages, setMessages] = useState([]);
    const [newMessage, setNewMessage] = useState('');
    const [cursor, setCursor] = useState(null);
    const [loading, setLoading] = useState(false);
    const [hasMore, setHasMore] = useState(true);
    const token = localStorage.getItem('token');
    const chatContainerRef = useRef(null);
    const dmWS = useRef(null);

    useEffect(() => {
        setMessages([]);
        setCursor(null);
        setHasMore(true);

        if (userId) {
            fetchMessages();
        }
    }, [userId]);

    const fetchMessages = async () => {
        if (!userId || !hasMore || loading) return;
        setLoading(true);

        const prevScrollHeight = chatContainerRef.current.scrollHeight;
        const prevScrollTop = chatContainerRef.current.scrollTop;

        try {
            const response = await axios.get(
                `http://localhost:8080/chats/${userId}?cursor=${cursor || ''}&limit=10`,
                {
                    headers: { Authorization: `Bearer ${token}` },
                }
            );

            const fetched = response.data.messages;
            setMessages((prev) => {
                if (!cursor) {
                    return [...fetched.slice(-6)];
                } else {
                    return [...fetched, ...prev];
                }
            });

            setCursor(response.data.nextCursor);
            if (!response.data.nextCursor) {
                setHasMore(false);
            }

            setTimeout(() => {
                const currScrollHeight = chatContainerRef.current.scrollHeight;
                chatContainerRef.current.scrollTop =
                    currScrollHeight - prevScrollHeight + prevScrollTop;
            }, 100);
        } catch (err) {
            if (err.response && err.response.status === 403) {
                if (!window.alertShown) {
                    window.alertShown = true;
                    alert("Oops! You are not connected with this user.");
                    window.location.href = "/connections";
                }
            } else {
                console.error('Error fetching messages:', err);
            }
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        function connectDMWebSocket() {
            const socket = new WebSocket(`ws://localhost:8080/ws/chat?token=${token}`);
            dmWS.current = socket;

            socket.onopen = () => {
                console.log('DM WebSocket connected');
            };

            socket.onmessage = (event) => {
                const msg = JSON.parse(event.data);
                setMessages((prev) => [...prev, msg]);
                setTimeout(() => {
                    if (chatContainerRef.current) {
                        chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
                    }
                }, 100);
            };

            socket.onclose = (e) => {
                console.error('DM WebSocket disconnected:', e);
                setTimeout(connectDMWebSocket, 2000);
            };

            socket.onerror = (err) => {
                console.error('DM WebSocket error:', err);
            };
        }

        if (userId) {
            connectDMWebSocket();
        }

        return () => {
            if (dmWS.current) dmWS.current.close();
        };
    }, [userId, token]);

    useEffect(() => {
        if (!userId) return;

        if (chatListWS?.current && chatListWS.current.readyState === WebSocket.OPEN) {
            chatListWS.current.send(
                JSON.stringify({ type: 'active_chat', chat_id: Number(userId) })
            );
        }

        return () => {
            if (chatListWS?.current && chatListWS.current.readyState === WebSocket.OPEN) {
                chatListWS.current.send(JSON.stringify({ type: 'inactive_chat' }));
            }
        };
    }, [userId, chatListWS]);

    const sendMessage = async () => {
        if (newMessage.trim() === '') return;

        const payload = { receiver_id: parseInt(userId, 10), content: newMessage };

        const temp = {
            chat_id: parseInt(userId, 10),
            sender_id: parseInt(localStorage.getItem('user_id'), 10),
            receiver_id: parseInt(userId, 10),
            content: newMessage,
            timestamp: new Date().toISOString(),
        };

        setMessages((prev) => [...prev, temp]);

        try {
            if (dmWS.current && dmWS.current.readyState === WebSocket.OPEN) {
                dmWS.current.send(JSON.stringify(payload));
            } else {
                const resp = await axios.post('http://localhost:8080/chats/send', payload, {
                    headers: { Authorization: `Bearer ${token}` },
                });
                setMessages((prev) =>
                    prev.map((m) => (m === temp ? resp.data.message : m))
                );
            }
        } catch (err) {
            console.error('Error sending message:', err);
            setMessages((prev) => prev.filter((m) => m !== temp));
        } finally {
            setNewMessage('');
            setTimeout(() => {
                if (chatContainerRef.current) {
                    chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
                }
            }, 100);
        }
    };

    return (
        <div className="chat-container">
            {!userId ? (
                <ChatList chats={chats} setChats={setChats} />
            ) : (
                <div className="chat-window">
                    {hasMore && (
                        <button
                            className="load-more-button"
                            onClick={fetchMessages}
                            disabled={loading}
                        >
                            {loading ? 'Loading...' : 'Load More'}
                        </button>
                    )}
                    <div className="chat-messages" ref={chatContainerRef}>
                        {messages.map((msg, index) => (
                            <div
                                key={index}
                                className={`message ${msg.sender_id === parseInt(userId, 10) ? 'received' : 'sent'
                                    }`}
                            >
                                <p>{msg.content}</p>
                                <span className="timestamp">
                                    {new Date(msg.timestamp).toLocaleString()}
                                </span>
                            </div>
                        ))}
                    </div>
                    <div className="chat-input">
                        <input
                            type="text"
                            value={newMessage}
                            onChange={(e) => setNewMessage(e.target.value)}
                            onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
                            placeholder="Type a message..."
                        />
                        <button onClick={sendMessage}>Send</button>
                    </div>
                </div>
            )}
        </div>
    );
}

export default Chats;
