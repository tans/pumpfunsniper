// const axios = require('axios');
// const dotenv = require('dotenv');

// dotenv.config();


// // Environment variables
// const API_URL = process.env.PUMP_CLIENT_API_URL;  // Base API URL
// const LOGIN_TOKEN = process.env.LOGIN_TOKEN;      // Authentication token
// const replyPayload = {
//     text: "haha",
//     mint: "v5ShBxGoA8rA54ANpm9xE7Sce3Qr1apBeeZcknhsyjb"
// };

// // Function to post replies
// async function replyComment(apiUrl, payload, token) {
    
//     const url = `${apiUrl}replies`;

//     if (!token) {
//         console.error('Authorization token is required');
//         return;
//     }

//     try {
//         const response = await axios.post(url, payload, {
//             headers: {
//                 'Content-Type': 'application/json',
//                 'Authorization': `Bearer ${token}`,
//             }
//         });

//         if (response.status !== 201) {
//             throw new Error(`Failed to post reply: ${response.status} ${response.statusText}`);
//         }

//         console.log('Reply posted successfully:', response.data);
//     } catch (error) {
//         console.error('Error posting reply:', error);
//     }
// }


// replyComment(API_URL, replyPayload, LOGIN_TOKEN);