import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import '../styles/profile.css';

const interestsCategories = {
    Sports: ["Football", "Basketball", "Tennis", "Swimming", "Hiking"],
    Food: ["Chinese", "Italian", "Mexican", "Indian", "Japanese", "Eating out", "Cooking at home"],
    Culture: ["Art", "History", "Museums", "Theater", "Traveling", "Languages", "Religion"],
    Games: ["Counter-Strike", "Minecraft", "Sims 4", "League of Legends", "Valorant", "Genshin Impact"],
    MoviesTV: ["Action", "Comedy", "Drama", "Sci-fi", "Romance", "Documentaries", "Anime", "Horror"],
};

const Profile = () => {
    const navigate = useNavigate();
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [profile, setProfile] = useState({
        name: '',
        age: '',
        location: '',
        info: '',
        selectedHobbies: [],
        profilePicture: '',
        latitude: 0,
        longitude: 0,
        preferredRadius: 50,
        gender: '',
        lookingFor: '',
    });
    const [file, setFile] = useState(null);
    const [error, setError] = useState(null);
    const [query, setQuery] = useState('');
    const [locationResults, setLocationResults] = useState([]);

    useEffect(() => {
        const fetchProfile = async () => {
            const token = localStorage.getItem('token');
            if (!token) {
                alert('Please log in first.');
                navigate('/login');
                return;
            }

            try {
                const response = await axios.get('http://localhost:8080/profile', {
                    headers: { Authorization: `Bearer ${token}` },
                });

                const userData = response.data.data;

                setProfile((prev) => ({
                    ...prev,
                    name: userData.name || '',
                    age: userData.age || '',
                    location: userData.location || '',
                    info: userData.info || '',
                    selectedHobbies: userData.interests
                        ? userData.interests.split(",").map((hobby) => hobby.trim())
                        : [],
                    profilePicture: userData.profilePicture || 'uploads/default-profile.png',
                    latitude: userData.latitude || 0,
                    longitude: userData.longitude || 0,
                    preferredRadius: userData.preferred_radius || 50,
                    gender: userData.gender || '',
                    lookingFor: userData.lookingFor || '',
                }));
            } catch (err) {
                console.error('Error fetching profile:', err);
                setError('Failed to load profile data. Please try again later.');
            }
        };

        fetchProfile();
    }, [navigate]);

    const toggleModal = () => setIsModalOpen(!isModalOpen);

    const handleCheckboxChange = (hobby) => {
        setProfile((prevProfile) => {
            const selectedHobbies = [...prevProfile.selectedHobbies];
            if (selectedHobbies.includes(hobby)) {
                return {
                    ...prevProfile,
                    selectedHobbies: selectedHobbies.filter((selected) => selected !== hobby),
                };
            } else {
                return {
                    ...prevProfile,
                    selectedHobbies: [...selectedHobbies, hobby],
                };
            }
        });
    };

    const handleSearchLocation = async () => {
        if (!query) return;
        try {
            const response = await axios.get('https://nominatim.openstreetmap.org/search', {
                params: {
                    q: query,
                    format: 'json',
                    addressdetails: 1,
                },
            });
            setLocationResults(response.data);
        } catch (err) {
            console.error('Error fetching locations:', err);
        }
    };

    const handleSelectLocation = (loc) => {
        setProfile({
            ...profile,
            location: loc.display_name,
            latitude: parseFloat(loc.lat),
            longitude: parseFloat(loc.lon),
        });
        setLocationResults([]);
    };

    const handleChange = (e) => {
        const { name, value } = e.target;
        setProfile({ ...profile, [name]: value });
    };

    const handleFileChange = (e) => {
        const selectedFile = e.target.files[0];
        setFile(selectedFile);

        if (selectedFile) {
            const previewURL = URL.createObjectURL(selectedFile);
            setProfile((prevProfile) => ({
                ...prevProfile,
                profilePicture: previewURL,
            }));
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError(null);

        const token = localStorage.getItem('token');
        if (!token) {
            alert('Please log in first.');
            return;
        }

        if (!profile.name || !profile.age || !profile.gender || !profile.lookingFor || !profile.info) {
            setError('Please fill out all required fields (including Gender and Looking For).');
            return;
        }

        if (profile.selectedHobbies.length === 0) {
            setError('Please select at least one hobby.');
            return;
        }

        if (!profile.location) {
            setError('Please select or enter a location.');
            return;
        }

        const formData = new FormData();
        formData.append('name', profile.name);
        formData.append('age', profile.age);
        formData.append('location', profile.location);
        formData.append('info', profile.info);
        formData.append('interests', profile.selectedHobbies.join(", "));
        formData.append('latitude', profile.latitude);        // NEW
        formData.append('longitude', profile.longitude);      // NEW
        formData.append('preferredRadius', profile.preferredRadius); // NEW

        formData.append('gender', profile.gender);
        formData.append('lookingFor', profile.lookingFor);
        if (file) formData.append('profilePicture', file);

        try {
            const response = await axios.put('http://localhost:8080/profile', formData, {
                headers: {
                    Authorization: `Bearer ${token}`,
                    'Content-Type': 'multipart/form-data',
                },
            });

            setProfile((prevProfile) => ({
                ...prevProfile,
                profilePicture: response.data.data.profilePicture,
            }));

            alert('Profile updated successfully!');
            navigate('/home');
        } catch (err) {
            console.error('Error updating profile:', err.response || err.message);
            setError('Failed to update profile. Please try again.');
        }
    };

    const handleRemoveProfilePicture = async () => {
        const token = localStorage.getItem('token');
        if (!token) {
            alert('Please log in first.');
            return;
        }

        try {
            await axios.put(
                'http://localhost:8080/profile/remove-picture',
                {},
                {
                    headers: { Authorization: `Bearer ${token}` },
                }
            );
            setProfile((prevProfile) => ({
                ...prevProfile,
                profilePicture: 'uploads/default-profile.png',
            }));
            alert('Profile picture removed successfully!');
        } catch (err) {
            console.error('Error removing profile picture:', err.response || err.message);
            setError('Failed to remove profile picture. Please try again.');
        }
    };

    return (
        <div className="profile-page-container">
            <div className="profile-setup-container">
                <h2>Edit Your Profile</h2>
                {error && <p className="error-message">{error}</p>}
                <form className="profile-form" encType="multipart/form-data">
                    <div className="form-group">
                        <label>Name</label>
                        <input
                            type="text"
                            name="name"
                            value={profile.name || ''}
                            onChange={handleChange}
                            placeholder="Enter your name"
                            required
                        />
                    </div>

                    <div className="form-group">
                        <label>Age</label>
                        <input
                            type="number"
                            name="age"
                            value={profile.age || ''}
                            onChange={handleChange}
                            min="18"
                            max="100"
                            placeholder="Enter your age"
                            required
                        />
                    </div>

                    <div className="form-group">
                        <label>Gender</label>
                        <select
                            name="gender"
                            value={profile.gender || ''}
                            onChange={handleChange}
                            required
                        >
                            <option value="">Select your gender...</option>
                            <option value="Male">Male</option>
                            <option value="Female">Female</option>
                            <option value="Non-binary">Non-binary</option>
                            <option value="Other">Other</option>
                            <option value="Prefer not to say">Prefer not to say</option>
                        </select>
                    </div>

                    <div className="form-group">
                        <label>Looking For</label>
                        <select
                            name="lookingFor"
                            value={profile.lookingFor || ''}
                            onChange={handleChange}
                            required
                        >
                            <option value="">Select what you're looking for...</option>
                            <option value="Male">Male</option>
                            <option value="Female">Female</option>
                            <option value="Non-binary">Non-binary</option>
                            <option value="Any">Any</option>
                            <option value="Other">Other</option>
                        </select>
                    </div>

                    <div className="form-group">
                        <label>Location</label>
                        <input
                            type="text"
                            value={query}
                            onChange={(e) => setQuery(e.target.value)}
                            placeholder="Search for your location..."
                        />
                        <button type="button" className="search-button" onClick={handleSearchLocation}>
                            Search
                        </button>
                        {locationResults.length > 0 && (
                            <div className="suggestions-container">
                                {locationResults.map((result, index) => (
                                    <div
                                        key={index}
                                        className="suggestion-item"
                                        onClick={() => handleSelectLocation(result)}
                                    >
                                        {result.display_name}
                                    </div>
                                ))}
                            </div>
                        )}
                        {profile.location && <p>Selected Location: {profile.location}</p>}
                    </div>

                    <div className="form-group">
                        <label>Preferred Radius (km)</label>
                        <input
                            type="number"
                            name="preferredRadius"
                            value={profile.preferredRadius || 0}
                            onChange={handleChange}
                            placeholder="Distance in km"
                        />
                        <p>We'll only recommend people within this radius if you have lat/lon.</p>
                    </div>

                    <div className="form-group">
                        <label>About Me</label>
                        <textarea
                            name="info"
                            value={profile.info || ''}
                            onChange={handleChange}
                            placeholder="Tell us about yourself"
                            required
                        />
                    </div>

                    <div className="form-group">
                        <label>Hobbies</label>
                        <button type="button" className="open-modal-button" onClick={toggleModal}>
                            Select Hobbies
                        </button>
                        <p>Selected Hobbies: {profile.selectedHobbies.join(", ")}</p>
                    </div>

                    <div className="form-group">
                        <label>Profile Picture</label>
                        <input type="file" accept="image/*" onChange={handleFileChange} />
                        {profile.profilePicture && (
                            <div>
                                <img
                                    src={
                                        profile.profilePicture.startsWith('blob:')
                                            ? profile.profilePicture
                                            : `http://localhost:8080/${profile.profilePicture}`
                                    }
                                    alt="Profile"
                                    className="profile-picture-preview"
                                />
                                <button
                                    type="button"
                                    className="remove-button"
                                    onClick={handleRemoveProfilePicture}
                                >
                                    Remove Profile Picture
                                </button>
                            </div>
                        )}
                    </div>
                </form>

                {isModalOpen && (
                    <div className="modal">
                        <div className="modal-content">
                            <h2>Select Your Hobbies</h2>
                            {Object.entries(interestsCategories).map(([category, hobbies]) => (
                                <div key={category} className="category-section">
                                    <h3>{category}</h3>
                                    {hobbies.map((hobby) => (
                                        <label key={hobby}>
                                            <input
                                                type="checkbox"
                                                value={hobby}
                                                checked={profile.selectedHobbies.includes(hobby)}
                                                onChange={() => handleCheckboxChange(hobby)}
                                            />
                                            {hobby}
                                        </label>
                                    ))}
                                </div>
                            ))}
                            <button type="button" className="close-modal-button" onClick={toggleModal}>
                                Save
                            </button>
                        </div>
                    </div>
                )}

                <div className="save-profile-container">
                    <button type="button" className="save-profile-button" onClick={handleSubmit}>
                        Save Profile
                    </button>
                </div>

            </div>
        </div>
    );
};

export default Profile;
