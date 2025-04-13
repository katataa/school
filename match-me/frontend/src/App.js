import React, { useState, useEffect, useRef } from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import axios from 'axios';
import './styles/general.css';

import NavBar from './components/NavBar';
import Register from './components/Register';
import Login from './components/Login';
import Profile from './components/Profile';
import HomePage from './components/HomePage';
import Chats from './components/Chats';
import RecommendationsPage from './components/RecommendationsPage';
import ConnectionsPage from './components/ConnectionsPage';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(!!localStorage.getItem('token'));
  const [hasProfile, setHasProfile] = useState(localStorage.getItem('hasProfile') === 'true');
  const [chats, setChats] = useState([]);
  const chatListWS = useRef(null);

  useEffect(() => {
    const handleStorageChange = () => {
      setIsAuthenticated(!!localStorage.getItem('token'));
      setHasProfile(localStorage.getItem('hasProfile') === 'true');
    };

    window.addEventListener('storage', handleStorageChange);
    return () => window.removeEventListener('storage', handleStorageChange);
  }, []);

  const handleLogin = (token, profileStatus) => {
    setIsAuthenticated(!!token);
    setHasProfile(profileStatus);
  };

  useEffect(() => {
    if (!isAuthenticated) return;

    const token = localStorage.getItem('token');
    if (!token) return;

    async function fetchInitialChats() {
      try {
        const response = await axios.get('http://localhost:8080/chats', {
          headers: { Authorization: `Bearer ${token}` },
        });
        setChats(response.data.chats || []);
      } catch (err) {
        console.error('Error fetching initial chats:', err);
      }
    }

    fetchInitialChats();
  }, [isAuthenticated]);

  useEffect(() => {
    if (!isAuthenticated) return;

    const token = localStorage.getItem('token');
    if (!token) return;

    function connectChatListWS() {
      const socket = new WebSocket(`ws://localhost:8080/ws/chat_list?token=${token}`);
      chatListWS.current = socket;

      socket.onopen = () => {
        console.log('Global Chat-List WebSocket connected');
      };

      socket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          if (data.type === 'new_message') {
            setChats((prevChats) => {
              const updatedChats = [...prevChats];
              const idx = updatedChats.findIndex((c) => c.id === data.chat_id);

              if (idx !== -1) {
                updatedChats[idx].latest_message = data.content;
                updatedChats[idx].latest_message_timestamp = data.timestamp;
                updatedChats[idx].unread_count = data.unread_count;

                const [moved] = updatedChats.splice(idx, 1);
                updatedChats.unshift(moved);
              } else {
                updatedChats.unshift({
                  id: data.chat_id,
                  user_id: data.sender_id,
                  name: data.sender_name || 'Unknown',
                  profile_picture: data.sender_profile_picture || 'uploads/default-profile.png',
                  latest_message: data.content,
                  latest_message_timestamp: data.timestamp,
                  unread_count: data.unread_count,
                });
              }
              return updatedChats;
            });
          }
        } catch (err) {
          console.error('Error parsing chat-list WS message:', err);
        }
      };

      socket.onclose = (e) => {
        console.warn('Chat-list WS closed, reconnecting...', e);
        setTimeout(connectChatListWS, 2000);
      };

      socket.onerror = (err) => {
        console.error('Chat-list WS error:', err);
      };
    }

    connectChatListWS();

    return () => {
      if (chatListWS.current) {
        chatListWS.current.close();
      }
    };
  }, [isAuthenticated]);

  return (
    <Router>
      <NavBar chats={chats} />

      <Routes>
        <Route
          path="/"
          element={
            isAuthenticated
              ? hasProfile
                ? <Navigate to="/home" />
                : <Navigate to="/profile" />
              : <Navigate to="/login" />
          }
        />
        <Route path="/register" element={<Register onRegister={handleLogin} />} />
        <Route path="/login" element={<Login onLogin={handleLogin} />} />
        <Route path="/profile" element={isAuthenticated ? <Profile /> : <Navigate to="/login" />} />
        <Route path="/home" element={isAuthenticated ? <HomePage /> : <Navigate to="/login" />} />

        <Route
          path="/chats"
          element={
            isAuthenticated ? (
              <Chats
                chats={chats}
                setChats={setChats}
                chatListWS={chatListWS}
              />
            ) : (
              <Navigate to="/login" />
            )
          }
        />
        <Route
          path="/chats/:userId"
          element={
            isAuthenticated ? (
              <Chats
                chats={chats}
                setChats={setChats}
                chatListWS={chatListWS}
              />
            ) : (
              <Navigate to="/login" />
            )
          }
        />

        <Route
          path="/recommendations"
          element={isAuthenticated ? <RecommendationsPage /> : <Navigate to="/login" />}
        />
        <Route
          path="/connections"
          element={isAuthenticated ? <ConnectionsPage /> : <Navigate to="/login" />}
        />
      </Routes>
    </Router>
  );
}

export default App;
