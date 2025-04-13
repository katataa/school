import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import '../styles/login.css';

const Login = ({ onLogin }) => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState(null);
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError(null);

        try {
            // Send login request
            const response = await axios.post('http://localhost:8080/login', { email, password });
            const token = response.data.token;

            if (!token) {
                setError('Unexpected error: No token received.');
                return;
            }

            // Save token to localStorage
            localStorage.setItem('token', token);

            // Fetch profile status
            const profileResponse = await axios.get('http://localhost:8080/profile', {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });

            const hasProfile = profileResponse.data.data?.info || false; // Check if profile is complete
            localStorage.setItem('hasProfile', hasProfile ? 'true' : 'false');

            // Update app state dynamically via the callback
            if (onLogin) {
                onLogin(token, hasProfile);
            }

            // Redirect to the appropriate page
            navigate(hasProfile ? '/home' : '/profile');
        } catch (err) {
            console.error('Login error:', err);

            // Set a more descriptive error message based on the response
            if (err.response?.status === 401) {
                setError('Invalid email or password. Please try again.');
            } else if (err.response?.status === 404) {
                setError('User not found. Please register.');
            } else {
                setError('Something went wrong. Please try again later.');
            }
        }
    };

    return (
        <div className="login-container">
            <h2>Login</h2>
            <form onSubmit={handleSubmit}>
                <div className="form-group">
                    <label>Email</label>
                    <input
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        placeholder="Enter your email"
                        required
                    />
                </div>
                <div className="form-group">
                    <label>Password</label>
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        placeholder="Enter your password"
                        required
                    />
                </div>
                <button type="submit">Login</button>
            </form>
            {error && <p className="error-message">{error}</p>}
            <p>
                Don't have an account? <a href="/register">Register here</a>
            </p>
        </div>
    );
};

export default Login;
