import React, { useEffect, useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import '../styles/home.css';

const HomePage = () => {
    const [connectionRequests, setConnectionRequests] = useState([]);
    const [profile, setProfile] = useState({});
    const [error, setError] = useState(null);
    const navigate = useNavigate();
    const token = localStorage.getItem('token');
    const isAuthenticated = !!token;

    const fetchHomeData = useCallback(async () => {
        if (!token) {
            console.error('No token found. Redirecting to login...');
            navigate('/login');
            return;
        }

        try {
            console.log('Fetching profile data...');
            const profileResponse = await axios.get('http://localhost:8080/profile', {
                headers: { Authorization: `Bearer ${token}` },
            });
            setProfile(profileResponse.data.data);

            console.log('Fetching connection requests...');
            const connectionRequestsResponse = await axios.get('http://localhost:8080/connections/requests', {
                headers: { Authorization: `Bearer ${token}` },
            });
            setConnectionRequests(connectionRequestsResponse.data.requests || []);
        } catch (err) {
            console.error('Error fetching data:', err.response || err.message);
            setError('Failed to load data. Please try again later.');
        }
    }, [navigate, token]);

    useEffect(() => {
        if (!isAuthenticated) {
            navigate('/login');
        } else {
            fetchHomeData();
        }
    }, [isAuthenticated, navigate, fetchHomeData]);


    const acceptConnection = async (requestId) => {
        try {
            await axios.post(`http://localhost:8080/connections/accept`, { request_id: requestId }, {
                headers: { Authorization: `Bearer ${token}` },
            });
            alert('Connection accepted!');
            fetchHomeData();
        } catch (err) {
            console.error('Error accepting connection:', err.response || err.message);
        }
    };

    const declineConnection = async (requestId) => {
        try {
            await axios.post(`http://localhost:8080/connections/decline`, { request_id: requestId }, {
                headers: { Authorization: `Bearer ${token}` },
            });
            alert('Connection declined.');
            fetchHomeData();
        } catch (err) {
            console.error('Error declining connection:', err.response || err.message);
        }
    };

    return (
        <div className="home-container">
            <main className="home-content">
                {error && <p className="error-message">{error}</p>}

                <section className="profile-overview">
                    <h2>Your Profile</h2>
                    <div className="profile-card">
                        <img
                            src={profile.profilePicture ? `http://localhost:8080/${profile.profilePicture}` : 'http://localhost:8080/uploads/default-profile.png'}
                            alt="Profile"
                            className="profile-picture"
                        />
                        <div className="profile-info">
                            <h3>{profile.name || 'Name not provided'}</h3>
                            <p><strong>Age:</strong> {profile.age || 'N/A'}</p>
                            <p><strong>Gender:</strong> {profile.gender || 'Not specified'}</p>
                            <p><strong>Email:</strong> {profile.email || 'Hidden'}</p>
                            <p><strong>Location:</strong> {profile.location || 'Not specified'}</p>
                            <p><strong>Info:</strong> {profile.info || 'No info available'}</p>
                            <p><strong>Interests:</strong> {profile.interests || 'Not specified'}</p>
                        </div>
                    </div>
                </section>

                <section className="recommendations">
                    <h2>Connection Requests</h2>
                    <div className="recommendation-list">
                        {connectionRequests.length > 0 ? (
                            connectionRequests.map((request) => (
                                <div key={request.id} className="recommendation-card">
                                    <img
                                        src={request.sender.profile_picture ? `http://localhost:8080/${request.sender.profile_picture}` : 'http://localhost:8080/uploads/default-profile.png'}
                                        alt={request.sender.name}
                                        className="user-picture"
                                    />
                                    <div className="recommendation-details">
                                        <h3 className="user-name">{request.sender.name}</h3>
                                        <p><strong>Age:</strong> {request.sender.age || 'Not specified'}</p>
                                        <p><strong>Gender:</strong> {request.sender.gender || 'Not specified'}</p>
                                        <p><strong>Info:</strong> {request.sender.info || 'No info available'}</p>
                                        <p><strong>Location:</strong> {request.sender.location || 'Not specified'}</p>
                                        <p><strong>Interests:</strong> {request.sender.interests || 'Not specified'}</p>
                                    </div>
                                    <div className="connection-buttons">
                                        <button className="accept-button" onClick={() => acceptConnection(request.id)}>Accept</button>
                                        <button className="decline-button" onClick={() => declineConnection(request.id)}>Decline</button>
                                    </div>
                                </div>
                            ))
                        ) : (
                            <p>No connection requests at this time.</p>
                        )}
                    </div>
                </section>
            </main>
        </div>
    );
};

export default HomePage;
