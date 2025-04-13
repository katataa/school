import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import '../styles/chats.css';

function ChatList({ chats, setChats }) {
    const navigate = useNavigate();
    const token = localStorage.getItem('token');

    useEffect(() => {
        const fetchChats = async () => {
            try {
                const response = await axios.get('http://localhost:8080/chats', {
                    headers: { Authorization: `Bearer ${token}` },
                });
                setChats(response.data.chats || []);
            } catch (err) {
                console.error('Error fetching chats:', err.response || err.message);
            }
        };

        fetchChats();
    }, [token, setChats]);

    // Mark as read in DB + local state
    const markAsRead = async (chatId) => {
        try {
            await axios.post(
                `http://localhost:8080/chats/${chatId}/read`,
                {},
                { headers: { Authorization: `Bearer ${token}` } }
            );
            setChats((prevChats) =>
                prevChats.map((chat) =>
                    chat.id === chatId ? { ...chat, unread_count: 0 } : chat
                )
            );
        } catch (err) {
            console.error('Failed to mark messages as read:', err);
        }
    };

    return (
        <div className="chat-list">
            <h2>Your Chats</h2>
            {chats.length > 0 ? (
                chats.map((chat) =>
                    chat.user_id && chat.name ? (
                        <div
                            key={chat.id}
                            className="chat-summary"
                            onClick={() => {
                                markAsRead(chat.id);
                                navigate(`/chats/${chat.user_id}`);
                            }}
                        >
                            <div className="chat-image-container">
                                <img
                                    src={`http://localhost:8080/${chat.profile_picture || 'uploads/default-profile.png'
                                        }`}
                                    alt={chat.name}
                                    className="chat-profile-picture"
                                />
                                {chat.unread_count > 0 && (
                                    <span className="button__badge">{chat.unread_count}</span>
                                )}
                            </div>
                            <div className="chat-details">
                                <h3>{chat.name}</h3>
                                <p>{chat.latest_message}</p>
                            </div>
                        </div>
                    ) : null
                )
            ) : (
                <p>No active chats</p>
            )}
        </div>
    );
}

export default ChatList;
