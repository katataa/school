import React, { useEffect, useState, useCallback } from 'react';
import axios from 'axios';
import '../styles/recommendations.css';

const RecommendationsPage = () => {
    const [recommendations, setRecommendations] = useState([]);
    const [error, setError] = useState(null);
    const token = localStorage.getItem('token');
    const [showFilterPopup, setShowFilterPopup] = useState(false);
    const [filterLocation, setFilterLocation] = useState('');
    const [filterAge, setFilterAge] = useState('');
    const [filterHobbies, setFilterHobbies] = useState('');
    const [filterGender, setFilterGender] = useState('');
    const [filterMode, setFilterMode] = useState('all');

    const loadRecs = useCallback(async () => {
        try {
            const recResp = await axios.get('http://localhost:8080/recommendations', {
                headers: { Authorization: `Bearer ${token}` },
                params: {
                    location: filterLocation,
                    age: filterAge,
                    hobbies: filterHobbies,
                    gender: filterGender,
                    mode: filterMode
                }
            });
            const ids = recResp.data.recommendations;
            const detailedRecs = [];
            for (const item of ids) {
                const userId = item.id;
                const userResp = await axios.get(`http://localhost:8080/users/${userId}`, {
                    headers: { Authorization: `Bearer ${token}` }
                });
                const bioResp = await axios.get(`http://localhost:8080/users/${userId}/bio`, {
                    headers: { Authorization: `Bearer ${token}` }
                });
                detailedRecs.push({
                    id: userId,
                    name: userResp.data.name,
                    profile_picture: userResp.data.profile_picture,
                    location: bioResp.data.location,
                    age: bioResp.data.age,
                    interests: bioResp.data.interests,
                    info: bioResp.data.info,
                    gender: bioResp.data.gender,
                });
            }
            setRecommendations(detailedRecs);
        } catch (err) {
            setError('Failed to load recommendations.');
        }
    }, [token, filterLocation, filterAge, filterHobbies, filterGender, filterMode]);

    useEffect(() => {
        loadRecs();
    }, [loadRecs]);

    const declineRecommendation = async (userId) => {
        try {
            await axios.post('http://localhost:8080/recommendations/decline', { request_id: userId }, {
                headers: { Authorization: `Bearer ${token}` },
            });
            setRecommendations((prev) => prev.filter((user) => user.id !== userId));
            alert('Recommendation declined.');
        } catch {
            alert('Failed to decline recommendation.');
        }
    };

    const connectUser = async (userId) => {
        try {
            await axios.post('http://localhost:8080/connections/request', { receiver_id: userId }, {
                headers: { Authorization: `Bearer ${token}` },
            });
            setRecommendations((prev) => prev.filter((user) => user.id !== userId));
            alert(`Connection request sent to user ${userId}`);
        } catch {
            alert('Failed to send connection request.');
        }
    };

    return (
        <div className="recommendations-page">
            <main className="recommendations-content">
                <div style={{ position: 'relative' }}>
                    <h2 className="recommended-title">Recommended Users</h2>
                    <button
                        className="filter-button"
                        onClick={() => setShowFilterPopup(!showFilterPopup)}
                    >
                        {showFilterPopup ? 'Close Filters' : 'Open Filters'}
                    </button>
                    {showFilterPopup && (
                        <div className="filter-popup">
                            <h4>Filter Options</h4>
                            <div className="filter-field">
                                <label htmlFor="filter-location">Location (Partial match):</label>
                                <input
                                    type="text"
                                    id="filter-location"
                                    placeholder="e.g., Tallinn"
                                    value={filterLocation}
                                    onChange={(e) => setFilterLocation(e.target.value)}
                                />
                            </div>
                            <div className="filter-field">
                                <label htmlFor="filter-age">Age:</label>
                                <input
                                    type="number"
                                    id="filter-age"
                                    placeholder="e.g., 30"
                                    value={filterAge}
                                    onChange={(e) => setFilterAge(e.target.value)}
                                />
                            </div>
                            <div className="filter-field">
                                <label htmlFor="filter-gender">Gender:</label>
                                <select
                                    id="filter-gender"
                                    value={filterGender}
                                    onChange={(e) => setFilterGender(e.target.value)}
                                >
                                    <option value="">Select your gender...</option>
                                    <option value="Male">Male</option>
                                    <option value="Female">Female</option>
                                    <option value="Non-binary">Non-binary</option>
                                    <option value="Other">Other</option>
                                    <option value="Prefer not to say">Prefer not to say</option>
                                </select>
                            </div>
                            <div className="filter-field">
                                <label htmlFor="filter-hobbies">Hobbies (comma-separated):</label>
                                <input
                                    type="text"
                                    id="filter-hobbies"
                                    placeholder="e.g., Football, Minecraft"
                                    value={filterHobbies}
                                    onChange={(e) => setFilterHobbies(e.target.value)}
                                />
                            </div>
                            <div className="filter-field">
                                <label htmlFor="filter-mode">Mode:</label>
                                <select
                                    id="filter-mode"
                                    value={filterMode}
                                    onChange={(e) => setFilterMode(e.target.value)}
                                >
                                    <option value="all">All (default scoring)</option>
                                    <option value="location">Only Location</option>
                                    <option value="age">Only Age</option>
                                    <option value="hobbies">Only Hobbies</option>
                                    <option value="gender">Only Gender</option>
                                </select>
                            </div>
                            <button className="apply-filter-button" onClick={() => {
                                loadRecs();
                                setShowFilterPopup(false);
                            }}>
                                Apply Filter
                            </button>
                        </div>
                    )}
                </div>
                {error && <p className="error-message">{error}</p>}
                <div className="recommendation-list">
                    {recommendations.length > 0 ? (
                        recommendations.map((user) => (
                            <div key={user.id} className="recommendation-card">
                                <img
                                    src={
                                        user.profile_picture.startsWith('http')
                                            ? user.profile_picture
                                            : `http://localhost:8080/${user.profile_picture}`
                                    }
                                    alt={user.name}
                                    className="user-picture"
                                />
                                <div className="recommendation-details">
                                    <h3 className="user-name">{user.name}</h3>
                                    <p><strong>Age:</strong> {user.age || 'Not specified'}</p>
                                    <p><strong>Gender:</strong> {user.gender || 'Not specified'}</p>
                                    <p><strong>Info:</strong> {user.info || 'No info available'}</p>
                                    <p><strong>Location:</strong> {user.location || 'Not specified'}</p>
                                    <p><strong>Interests:</strong> {user.interests || 'Not specified'}</p>
                                </div>
                                <div className="connection-buttons">
                                    <button className="connect-button" onClick={() => connectUser(user.id)}>Connect</button>
                                    <button className="decline-button" onClick={() => declineRecommendation(user.id)}>Decline</button>
                                </div>
                            </div>
                        ))
                    ) : (
                        <p>No recommendations available at this time.</p>
                    )}
                </div>
            </main>
        </div>
    );
};

export default RecommendationsPage;
