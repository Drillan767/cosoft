import chalk from 'chalk';
import ora from 'ora';
import Table from 'cli-table3';
import select from '@inquirer/select';
import { format } from 'date-fns';
import { API_CONFIG } from '../config/env.js';
import { makeApiRequest, checkAuthentication, getRefererHeaders } from '../config/api.js';
// Import getRoomColor function from calendarView.js and fetchRooms for room mapping
import { getRoomColor } from './calendarView.js';
import { fetchRooms } from './listRooms.js';
import { showStandardMenu } from '../utils/sharedMenus.js';

// Function to fetch user's bookings from the API
// If specificDate is provided, results will be filtered to that date only

// Helper function to create API ID to UI ID mapping for consistent colors
async function createApiToUiIdMapping() {
    try {
        const rooms = await fetchRooms();
        const apiToUiId = {};

        if (rooms && Array.isArray(rooms)) {
            rooms.forEach(room => {
                apiToUiId[room.apiId] = room.id;
            });
        }

        return apiToUiId;
    } catch (error) {
        console.error('Warning: Could not fetch room mapping for consistent colors:', error.message);
        return {};
    }
}

export async function fetchUserBookings(specificDate = null) {
    try {
        // Use the correct endpoint for fetching user's bookings
        // Note: The API doesn't support date filtering, so we fetch all and filter client-side
        const url = `${API_CONFIG.BASE_URL}/v2/api/api/Reservations/get-current-and-incoming?PerPage=100&Page=1`;

        const responseData = await makeApiRequest({
            method: 'get',
            url: url,
            extraHeaders: getRefererHeaders('/v2/my-reservations')
        });

        // Process the bookings data
        let bookingsData = [];
        if (responseData && responseData.data && Array.isArray(responseData.data)) {
            bookingsData = responseData.data;
        } else if (responseData && Array.isArray(responseData)) {
            bookingsData = responseData;
        }

        if (bookingsData.length === 0) {
            return [];
        }

        // Map the API response to a more usable format
        const formattedBookings = bookingsData.map(booking => ({
            id: booking.OrderResourceRentId,
            roomId: booking.ItemId, // This is the correct room ID that matches room.apiId
            room: booking.ItemName,
            startDate: new Date(booking.Start), // Keep as Date object
            endDate: new Date(booking.End),     // Keep as Date object
            date: format(new Date(booking.Start), 'yyyy-MM-dd'), // Formatted date for display
            time: `${format(new Date(booking.Start), 'HH:mm')} - ${format(new Date(booking.End), 'HH:mm')}`, // Formatted time for display
            price: booking.Prices?.EuroTTC ? `${booking.Prices.EuroTTC.toFixed(2)} €` : 'N/A',
            capacity: booking.Capacity || 0
        }));

        // Filter by specific date if provided
        if (specificDate) {
            const targetDate = format(specificDate, 'yyyy-MM-dd');
            return formattedBookings.filter(booking => booking.date === targetDate);
        }

        return formattedBookings;
    } catch (error) {
        console.error('Error fetching user bookings:', error);
        throw error;
    }
}

export async function myBookings(options = {}) {

    // Check if running in interactive mode - default to false (non-interactive)
    const isInteractive = options.interactive === true;

    const spinner = ora('Loading your bookings...').start();

    try {
        // Call the API to get real bookings and create room mapping for consistent colors
        const [bookings, apiToUiIdMapping] = await Promise.all([
            fetchUserBookings(),
            createApiToUiIdMapping()
        ]);

        // Filter by date if option is provided
        let filteredBookings = bookings;
        if (options.date) {
            filteredBookings = bookings.filter(booking => booking.date === options.date);
        }

        spinner.succeed('Bookings loaded');

        if (filteredBookings.length === 0) {
            console.log(chalk.yellow('\nNo bookings found for the specified criteria.'));

            if (!isInteractive) {
                return;
            }

            // If in interactive mode, show a menu to navigate back
            await showStandardMenu();
        }

        console.log(chalk.bold('\nYour Bookings:'));

        // Create a table with colored lines including Booking ID
        const table = new Table({
            head: [
                chalk.white('Booking ID'),
                chalk.white('Room'),
                chalk.white('Date'),
                chalk.white('Time'),
                chalk.white('Capacity'),
                chalk.white('Price')
            ],
            colWidths: [38, 20, 12, 18, 8, 10]
        });

        // Add rows with colored text based on room ID
        filteredBookings.forEach(booking => {
            // Use UI ID for consistent colors, fallback to API ID if mapping not available
            const uiId = apiToUiIdMapping[booking.roomId] || booking.roomId;
            const roomColor = getRoomColor(uiId);

            table.push([
                chalk[roomColor](booking.id),
                chalk[roomColor](booking.room),
                chalk[roomColor](booking.date),
                chalk[roomColor](booking.time),
                chalk[roomColor](booking.capacity),
                chalk[roomColor](booking.price)
            ]);
        });

        console.log(table.toString());

        // If not in interactive mode, exit after displaying bookings
        if (!isInteractive) {
            return;
        }

        // After displaying the table, show a menu with options
        await showStandardMenu({ currentOptions: options });

    } catch (error) {
        spinner.fail('Failed to load bookings');
        console.error(chalk.red(error.message));
        if (!isInteractive) {
            process.exit(1);
        }
    }
}