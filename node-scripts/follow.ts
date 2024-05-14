// const axios = require('axios');
// const dotenv = require('dotenv');

// dotenv.config();

// const API_URL = process.env.PUMP_CLIENT_API_URL;
// const USER_ID = process.env.WALLET_ADDRESS;
// const LOGIN_TOKEN = process.env.LOGIN_TOKEN;

// async function followUser(apiUrl, userId, token) {
//   if (!userId) {
//     console.error('User ID is required');
//     return;
//   }
//   if (!token) {
//     console.error('Authorization token is required');
//     return;
//   }

//   const url = `${apiUrl}` + `following/${userId}`;
  
//   try {
//     const response = await axios.post(url, {}, {
//       headers: {
//         'Content-Type': 'application/json',
//         'Authorization': `Bearer ${token}`,
//       }
//     });

//     if (response.status !== 201) {
//       throw new Error(`Failed to create following relationship: ${response.status} ${response.statusText}`);
//     }

//     console.log('Following relationship successfully created:', response.data);
//   } catch (error) {
//     console.error('Error creating following relationship:', error);
//   }
// }

// followUser(API_URL, USER_ID, LOGIN_TOKEN);

