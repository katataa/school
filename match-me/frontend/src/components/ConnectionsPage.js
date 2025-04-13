import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import '../styles/connections.css';

const ConnectionsPage = () => {
    const [connections, setConnections] = useState([]);
    const [error, setError] = useState(null);
    const navigate = useNavigate();
    const token = localStorage.getItem('token');

    useEffect(() => {
        const fetchConnections = async () => {
            try {
                const resp = await axios.get('http://localhost:8080/connections', {
                    headers: { Authorization: `Bearer ${token}` },
                });

                const connectionIDs = resp.data.connections || [];
                const detailedConnections = [];

                for (const item of connectionIDs) {
                    const userId = item.id;

                    const userResp = await axios.get(`http://localhost:8080/users/${userId}`, {
                        headers: { Authorization: `Bearer ${token}` },
                    });

                    const bioResp = await axios.get(`http://localhost:8080/users/${userId}/bio`, {
                        headers: { Authorization: `Bearer ${token}` },
                    });

                    const userData = {
                        id: userId,
                        name: userResp.data.name,
                        profile_picture: userResp.data.profile_picture,
                        age: bioResp.data.age,
                        location: bioResp.data.location,
                        interests: bioResp.data.interests,
                        info: bioResp.data.info,
                        gender: bioResp.data.gender,
                    };

                    detailedConnections.push(userData);
                }

                setConnections(detailedConnections);
            } catch (err) {
                console.error('Error fetching connections:', err.response || err.message);
                setError('Failed to load connections. Please try again later.');
            }
        };

        fetchConnections();
    }, [token]);

    const handleDisconnect = async (userId) => {
        try {
            await axios.post(
                'http://localhost:8080/connections/disconnect',
                { user_id: userId },
                { headers: { Authorization: `Bearer ${token}` } }
            );
            setConnections((prev) => prev.filter((c) => c.id !== userId));
            alert(`Disconnected from user ${userId}`);
        } catch (err) {
            console.error('Error disconnecting:', err.response || err.message);
            alert('Failed to disconnect. Please try again.');
        }
    };

    return (
        <div className="connections-container">
            <main className="connections-content">
                <h2>Your Connections</h2>
                {error && <p className="error-message">{error}</p>}
                <div className="connections-list">
                    {connections.length > 0 ? (
                        connections.map((user) => (
                            <div key={user.id} className="connection-card">
                                <img
                                    src={user.profile_picture
                                        ? `http://localhost:8080/${user.profile_picture}`
                                        : 'http://localhost:8080/uploads/default-profile.png'
                                    }
                                    alt={user.name}
                                    className="user-picture"
                                />
                                <div className="connection-details">
                                    <h3 className="user-name">{user.name}</h3>
                                    <p><strong>Age:</strong> {user.age || 'Not specified'}</p>
                                    <p><strong>Gender:</strong> {user.gender || 'Not specified'}</p>
                                    <p><strong>Info:</strong> {user.info || 'Not specified'}</p>
                                    <p><strong>Location:</strong> {user.location || 'Not specified'}</p>
                                    <p><strong>Interests:</strong> {user.interests || 'Not specified'}</p>
                                </div>
                                <div className="button-group">
                                    <button
                                        className="message-button"
                                        onClick={() => navigate(`/chats/${user.id}`)}
                                    >
                                        Message
                                    </button>
                                    <button
                                        className="message-button disconnect-button"
                                        onClick={() => handleDisconnect(user.id)}
                                    >
                                        Disconnect
                                    </button>
                                </div>
                            </div>
                        ))
                    ) : (
                        <p>You have no connections yet.</p>
                    )}
                </div>
            </main>
        </div>
    );
};

export default ConnectionsPage;
