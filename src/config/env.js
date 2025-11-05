// Environment variables for the application
import dotenv from 'dotenv';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';
import chalk from 'chalk';

// Load environment variables from .env file (for non-auth config only)
const __filename = (typeof import.meta !== 'undefined' && import.meta.url) ? fileURLToPath(import.meta.url) : (typeof __filename !== 'undefined' ? __filename : process.cwd() + '/env.js');
const __dirname = (typeof import.meta !== 'undefined' && import.meta.url) ? path.dirname(fileURLToPath(import.meta.url)) : (typeof __dirname !== 'undefined' ? __dirname : process.cwd());
dotenv.config({ path: path.resolve(__dirname, '../../.env') });

// Define path to auth file
const AUTH_FILE_PATH = path.resolve(__dirname, '../../.auth');

// Function to load auth tokens from .auth file
function loadAuthTokens() {
    try {
        // Check if .auth file exists
        if (fs.existsSync(AUTH_FILE_PATH)) {
            // Read and parse auth file
            const authContent = fs.readFileSync(AUTH_FILE_PATH, 'utf8');
            const authTokens = JSON.parse(authContent);
            return {
                AUTH_TOKEN: authTokens.authToken || '',
                REFRESH_TOKEN: authTokens.refreshToken || ''
            };
        } else {
            console.error('Error: Authentication file not found.');
            console.log(chalk.yellow('Please login using the following command:'));
            console.log(chalk.blue('  cosoft login'));
            process.exit(1);
        }
    } catch (error) {
        console.error('Error loading auth tokens:', error.message);
        console.log(chalk.yellow('Please login using the following command:'));
        console.log(chalk.blue('  cosoft login'));
        process.exit(1);
    }
}

// Lazy auth token loading
let AUTH_TOKENS = null;

function getAuthTokens() {
    if (!AUTH_TOKENS) {
        AUTH_TOKENS = loadAuthTokens();
    }
    return AUTH_TOKENS;
}

// API configuration
export const API_CONFIG = {
    BASE_URL: process.env.COSOFT_API_BASE_URL || 'https://hub612.cosoft.fr',
    COWORKING_SPACE_ID: process.env.COSOFT_SPACE_ID || 'a4928a70-38c1-42b9-96f9-b2dd00db5b02',
    CATEGORY_ID: process.env.COSOFT_CATEGORY_ID || '7f1e5757-b9b9-4530-84ad-b2dd00db5f0f',
    get AUTH_TOKEN() { return getAuthTokens().AUTH_TOKEN; },
    get REFRESH_TOKEN() { return getAuthTokens().REFRESH_TOKEN; }
};