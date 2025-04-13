import React from 'react';
import { useNavigate } from 'react-router-dom';
import '../styles/navbar.css';

const NavBar = ({ chats = [] }) => {
    const token = localStorage.getItem('token');
    const navigate = useNavigate();

    const handleLogout = () => {
        localStorage.removeItem('token');
        navigate('/login');
    };

    const totalUnread = chats.reduce((acc, chat) => acc + (chat.unread_count || 0), 0);

    return (
        <header className="navbar">
            <h1>Match-Me Web</h1>
            <nav>
                <button onClick={handleLogout}>Log Out</button>
                <button onClick={() => navigate('/home')}>Home</button>
                <button onClick={() => navigate('/recommendations')}>Recommendations</button>
                <button onClick={() => navigate('/connections')}>Connections</button>
                <button onClick={() => navigate('/chats')} className="chats-button">
                    Chats
                    {totalUnread > 0 && (
                        <span className="button__badge">{totalUnread}</span>
                    )}
                </button>

                <button onClick={() => navigate('/profile')}>Edit Profile</button>
            </nav>
        </header>
    );
};

export default NavBar;
