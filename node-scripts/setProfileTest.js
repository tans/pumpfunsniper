// const axios = require('axios');
// const fs = require('fs');
// const path = require('path');

// // interface ProfileResponse {
// //     address: string;
// //     likes_received: number;
// //     mentions_received: number;
// //     username: string;
// //     profile_image: string;
// //     last_username_update_timestamp: number;
// //     followers: number;
// //     following: number;
// //     bio: string | null;
// // }

// // interface ImageUploadResponse {
// //     fileUri: string;
// // }


// const API_URL = process.env.PUMP_DEVNET_CLIENT_API_URL;
// const LOGIN_TOKEN = process.env.LOGIN_TOKEN_DEVNET;
// const PROFILE_IMAGE = 'https://picsum.photos/200';
// const USERNAME = 'thisIsMyUserName';
// const BIO_TEXT = 'GM fuckers, this is my BIO. lets pump it up!';

// const downloadImage = async (imageUrl) => {
//     const response = await axios({
//         url: imageUrl,
//         method: 'GET',
//         responseType: 'arraybuffer' // ensures that the binary data is handled correctly
//     });
//     console.log(Buffer.from(response.data, 'binary'))
//     return Buffer.from(response.data, 'binary');
// };

// const setProfile = async (apiUrl, loginToken, imageUrl, username, bioText) => {
//     const imageUploadUrl = `${apiUrl}api/ipfs-file`;
//     const profileUpdateUrl = `${apiUrl}users`;

//     try {
//         // Downloading the image from the URL
//         const imageData = await downloadImage(imageUrl);

//         // Prepare the image data and multipart form payload
//         const boundary = "----WebKitFormBoundary7MA4YWxkTrZu0gW";
//         const payload = Buffer.concat([
//             Buffer.from(
//                 `--${boundary}\r\nContent-Disposition: form-data; name="file"; filename="image.jpg"\r\nContent-Type: image/jpeg\r\n\r\n`,
//                 'utf-8',
//             ),
//             imageData, 
//             Buffer.from(`\r\n--${boundary}--\r\n`, 'utf-8'),
//         ]);
//         console.log('ImageData', imageData);
//         const headersForImageUpload = {
//             'Content-Type': `multipart/form-data; boundary=${boundary}`,
//             'Authorization': `Bearer ${loginToken}`,
//             'Content-Length': payload.length.toString()
//         };

//         // Uploading the image
//         const uploadResponse = await axios.post(imageUploadUrl, payload, { headers: headersForImageUpload });
//         const fileUri = uploadResponse.data.fileUri;

//         console.log('Image uploaded successfully:', fileUri);

//         // Prepare the profile update payload
//         const profilePayload = {
//             profileImage: fileUri,
//             username: username,
//             bio: bioText
//         };

//         const headersForProfileUpdate = {
//             'Content-Type': 'application/json',
//             'Authorization': `Bearer ${loginToken}`
//         };

//         // Updating the user profile
//         const profileResponse = await axios.post(profileUpdateUrl, profilePayload, { headers: headersForProfileUpdate });
//         console.log('Profile updated successfully:', profileResponse.data);
//     } catch (error) {
//         console.error('Failed to set profile:', error.response ? error.response.data : error.message);
//     }
// };
// console.log('Profile updated successfully:', API_URL, LOGIN_TOKEN, PROFILE_IMAGE, USERNAME, BIO_TEXT);
// setProfile(API_URL, LOGIN_TOKEN, PROFILE_IMAGE, USERNAME, BIO_TEXT);
