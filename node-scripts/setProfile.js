// const axios = require('axios');
// const dotenv = require('dotenv');

// dotenv.config();

// const API_URL = process.env.PUMP_DEVNET_CLIENT_API_URL;
// const LOGIN_TOKEN = process.env.LOGIN_TOKEN_DEVNET;
// const USERNAME = 'Dicky69';
// const BIO_TEXT = 'GM fuckers, thats my new bio. Lets pump it!';

// async function updateProfile(apiUrl, username, bioText, token) {
//   if (!username) {
//     console.error('Username is required');
//     return;
//   }
//   if (!token) {
//     console.error('Authorization token is required');
//     return;
//   }

//   const url = `${apiUrl}users`;

//   try {
//     const response = await axios.post(url, {
//       username: username,
//       bio: bioText
//     }, {
//       headers: {
//         'Content-Type': 'application/json',
//         'Authorization': `Bearer ${token}`,
//       }
//     });

//     if (response.status !== 201) {
//       throw new Error(`Failed to update profile: ${response.status} ${response.statusText}`);
//     }

//     console.log('Profile updated successfully:', response.data);
//   } catch (error) {
//     console.error('Error updating profile:', error);
//   }
// }

// updateProfile(API_URL, USERNAME, BIO_TEXT, LOGIN_TOKEN);
