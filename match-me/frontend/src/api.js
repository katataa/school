import axios from 'axios';

const API_URL = 'http://localhost:8080';

export const getProfile = async (token) => {
    const response = await axios.get(`${API_URL}/profile`, {
        headers: {
            Authorization: token,
        },
    });
    return response.data;
};
