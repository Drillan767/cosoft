// API utilities for making requests with authentication handling
import axios from 'axios';
import chalk from 'chalk';
import { API_CONFIG } from './env.js';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

// Get the directory path for saving the auth file
const __filename = (typeof import.meta !== 'undefined' && import.meta.url) ? fileURLToPath(import.meta.url) : (typeof __filename !== 'undefined' ? __filename : process.cwd() + '/api.js');
const __dirname = (typeof import.meta !== 'undefined' && import.meta.url) ? path.dirname(fileURLToPath(import.meta.url)) : (typeof __dirname !== 'undefined' ? __dirname : process.cwd());
const AUTH_FILE_PATH = path.resolve(__dirname, '../../.auth');

/**
 * Check if the current authentication token is valid and display user information
 * @returns {Promise<boolean>} True if authenticated, false otherwise
 */
export async function checkAuthentication() {
    try {
        console.log(chalk.gray(`Verifying authentication status...`));

        // Use the correct authentication endpoint from the API
        const url = `${API_CONFIG.BASE_URL}/v2/api/api/users/auth`;

        const response = await axios({
            method: 'get',
            url: url,
            headers: {
                'accept': 'application/json',
                'content-type': 'application/json',
                'cookie': `w_auth=${API_CONFIG.AUTH_TOKEN}; w_auth_refresh=${API_CONFIG.REFRESH_TOKEN}`,
                'origin': API_CONFIG.BASE_URL,
                'referer': `${API_CONFIG.BASE_URL}/v2/`
            }
        });

        // Check the isAuth property in the response
        if (response.data && response.data.isAuth === true) {
            // Display logged-in user information if available
            // API returns User with capital 'U', not lowercase 'user'
            if (response.data.User) {
                const user = response.data.User;
                const fullName = user.FirstName && user.LastName
                    ? `${user.FirstName} ${user.LastName}`
                    : user.Email || 'Unknown user';

                console.log(chalk.green(`✓ Authenticated as: ${chalk.bold(fullName)}`));
            } else {
                console.log(chalk.green(`✓ Authentication successful`));
            }
            return true;
        } else {
            console.error(chalk.red.bold('\n⚠️  Authentication Error: Invalid authentication response'));
            return false;
        }
    } catch (error) {
        // Any error means the token is invalid or expired
        console.error(chalk.red.bold('\n⚠️  Authentication Error: Your session has expired'));
        console.log(chalk.yellow('Please login using the following command:'));
        console.log(chalk.blue('  cosoft login'));
        console.log(chalk.gray('\nThis will prompt you for your email and password and save your authentication tokens.'));

        return false;
    }
}

/**
 * Login with email and password
 * @param {string} email - User email
 * @param {string} password - User password
 * @returns {Promise<Object>} Login result with user info
 */
export async function login(email, password) {
    try {
        console.log(chalk.gray(`Authenticating with CoSoft...`));

        const response = await axios({
            method: 'post',
            url: `${API_CONFIG.BASE_URL}/v2/api/api/users/login`,
            headers: {
                'accept': 'application/json',
                'content-type': 'application/json',
                'origin': API_CONFIG.BASE_URL,
                'referer': `${API_CONFIG.BASE_URL}/v2/login`
            },
            data: {
                email,
                password
            }
        });

        if (response.data && response.data.isAuth === true) {
            // Extract tokens from the response
            if (response.data.User && response.data.User.JwtToken) {
                // Save tokens to .auth file
                const authData = {
                    authToken: response.data.User.JwtToken,
                    refreshToken: response.headers['set-cookie'] ?
                        extractRefreshToken(response.headers['set-cookie']) :
                        ''
                };

                // Save auth data to file
                fs.writeFileSync(AUTH_FILE_PATH, JSON.stringify(authData, null, 2));

                return {
                    success: true,
                    user: response.data.User,
                    message: 'Authentication successful'
                };
            } else {
                return {
                    success: false,
                    message: 'Authentication succeeded but no token was provided'
                };
            }
        } else {
            return {
                success: false,
                message: response.data.Message || 'Authentication failed'
            };
        }
    } catch (error) {
        // Handle various error scenarios
        if (error.response) {
            // The request was made but the server responded with an error
            return {
                success: false,
                message: `Authentication failed: ${error.response.status} ${error.response.statusText}`
            };
        } else if (error.request) {
            // The request was made but no response was received
            return {
                success: false,
                message: 'Network error: Could not connect to the CoSoft API'
            };
        } else {
            // Something else happened while setting up the request
            return {
                success: false,
                message: `Error: ${error.message}`
            };
        }
    }
}

/**
 * Extract refresh token from cookies
 * @param {Array<string>} cookies - Cookie headers
 * @returns {string} Refresh token or empty string
 */
function extractRefreshToken(cookies) {
    if (!cookies || !Array.isArray(cookies)) {
        return '';
    }

    for (const cookie of cookies) {
        const match = cookie.match(/w_auth_refresh=([^;]+)/);
        if (match && match[1]) {
            return match[1];
        }
    }

    return '';
}

/**
 * Makes an authenticated API request to the CoSoft API with error handling
 * @param {Object} options - Request options
 * @param {string} options.method - HTTP method (get, post, delete, etc.)
 * @param {string} options.url - The full URL to make the request to
 * @param {Object} [options.data] - Optional data to send with the request
 * @param {Object} [options.extraHeaders] - Additional headers to include
 * @returns {Promise<any>} - The response data
 */
export async function makeApiRequest(options) {
    try {
        const { method, url, data, extraHeaders = {} } = options;

        const headers = {
            'accept': 'application/json',
            'content-type': 'application/json',
            'cookie': `w_auth=${API_CONFIG.AUTH_TOKEN}; w_auth_refresh=${API_CONFIG.REFRESH_TOKEN}`,
            'origin': API_CONFIG.BASE_URL,
            ...extraHeaders
        };

        const response = await axios({
            method,
            url,
            headers,
            data
        });

        return response.data;
    } catch (error) {
        // Handle token expiration
        if (error.response) {
            // Handle HTTP errors
            throw new Error(`API request failed: ${error.response.status} ${error.response.statusText}`);
        }
        // Handle network errors
        if (error.request) {
            throw new Error('Network error: Could not connect to the CoSoft API. Please check your internet connection.');
        }
        // Handle other errors
        throw error;
    }
}

/**
 * Get standard request headers for the referer URL
 * @param {string} path - Path part of the referer URL
 * @returns {Object} - Headers object with referer
 */
export function getRefererHeaders(path) {
    return {
        'referer': `${API_CONFIG.BASE_URL}${path}`
    };
}