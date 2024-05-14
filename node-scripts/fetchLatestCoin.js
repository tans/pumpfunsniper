// const axios = require('axios');

// const API_URL = process.env.PUMP_CLIENT_API_URL_DEVNET;
// const LOGIN_TOKEN = process.env.LOGIN_TOKEN_DEVNET;

// async function fetchLatestCoin(apiUrl, accessToken) {
//     try {
//         const response = await axios.get(`${apiUrl}coins/latest`, {
//             headers: { 'Authorization': `Bearer ${accessToken}` }
//         });
//         console.log('Latest coin details:', response.data);
//         return response.data;
//     } catch (error) {
//         console.error('Failed to fetch latest coin:', error.response ? error.response.data : error.message);
//         return null;
//     }
// } fetchLatestCoin(API_URL, LOGIN_TOKEN)
