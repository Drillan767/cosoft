import chalk from 'chalk';
import ora from 'ora';
import select from '@inquirer/select';
import input from '@inquirer/input';
import {
    format,
    addDays,
    parseISO,
    isBefore,
    isAfter
} from 'date-fns';
import { API_CONFIG } from '../config/env.js';
import { makeApiRequest, checkAuthentication, getRefererHeaders } from '../config/api.js';
// Import the fetchRooms function from listRooms.js
import { fetchRooms } from './listRooms.js';
import { bookRoom } from './bookRoom.js';
import { cancelBooking } from './cancelBooking.js';
import { myBookings, fetchUserBookings } from './myBookings.js';
import { listRooms } from './listRooms.js';

// Static list of colors for meeting rooms
const roomColors = [
    'blue',      // ID % 6 = 1
    'green',     // ID % 6 = 2
    'magenta',   // ID % 6 = 3
    'yellow',    // ID % 6 = 4
    'cyan',      // ID % 6 = 5
    'red'        // ID % 6 = 0
];

// Helper function to get color for a room
export function getRoomColor(roomId) {
    let numericId;

    // If roomId is a string that looks like a UUID, extract a numeric value from it
    if (typeof roomId === 'string' && roomId.includes('-')) {
        // Generate a consistent numeric ID from the full UUID 
        // We'll sum the character codes to ensure the same UUID always gets the same color
        numericId = roomId.split('-')
            .join('')  // Remove hyphens
            .split('')
            .reduce((sum, char, index) => sum + char.charCodeAt(0) * (index + 1), 0);
    } else {
        // For numeric IDs, use as is
        numericId = roomId;
    }

    // Ensure we always get a valid index
    return roomColors[(numericId % roomColors.length) || 0];
}

// Helper to convert a time string (HH:MM) to minutes since midnight
function timeToMinutes(timeStr) {
    const [hours, minutes] = timeStr.split(':').map(Number);
    return hours * 60 + minutes;
}

// Function to fetch busy times from the API for a specific room
async function fetchRoomBusyTimes(roomId, date, roomIdMapping) {
    try {
        // Use the provided roomIdMapping instead of the imported ROOM_IDS
        const apiRoomId = roomIdMapping[roomId];
        if (!apiRoomId) {
            throw new Error(`No API room ID mapping found for room ID: ${roomId}`);
        }

        const dateStr = format(date, 'yyyy-MM-dd');
        const url = `${API_CONFIG.BASE_URL}/v2/api/api/CoworkingSpace/${API_CONFIG.COWORKING_SPACE_ID}/category/${API_CONFIG.CATEGORY_ID}/item/${apiRoomId}/busytimes`;

        const refererPath = `/v2/new-reservation/${API_CONFIG.COWORKING_SPACE_ID}/${API_CONFIG.CATEGORY_ID}/${apiRoomId}`;

        const responseData = await makeApiRequest({
            method: 'post',
            url: url,
            extraHeaders: getRefererHeaders(refererPath),
            data: {
                startDate: `${dateStr}T00:00:00`,
                endDate: `${dateStr}T23:59:00`
            }
        });

        // Handle the response data format - the busy times are in response.data.data
        let busyTimesData = [];
        if (responseData && responseData.data && Array.isArray(responseData.data)) {
            busyTimesData = responseData.data;
        } else if (Array.isArray(responseData)) {
            busyTimesData = responseData;
        } else if (responseData && responseData.busyTimes && Array.isArray(responseData.busyTimes)) {
            busyTimesData = responseData.busyTimes;
        } else if (responseData && typeof responseData === 'object') {
            // Try to find an array in the response
            for (const key in responseData) {
                if (Array.isArray(responseData[key])) {
                    busyTimesData = responseData[key];
                    break;
                }
            }
        }

        // Filter out bookings that are outside the requested date
        const currentDateStr = format(date, 'yyyy-MM-dd');
        busyTimesData = busyTimesData.filter(busyTime => {
            // Check if the Start date matches our current date
            if (busyTime.Start && busyTime.Start.startsWith(currentDateStr)) {
                return true;
            }
            if (busyTime.startDate && busyTime.startDate.startsWith(currentDateStr)) {
                return true;
            }
            if (busyTime.StartDate && busyTime.StartDate.startsWith(currentDateStr)) {
                return true;
            }
            return false;
        });

        // Convert API response to our booking format
        return busyTimesData.map(busyTime => {
            // Safely parse dates with error handling
            let startTime = '00:00';
            let endTime = '00:00';

            try {
                // Check for the different possible field names based on the API response
                const startDate = busyTime.Start || busyTime.startDate || busyTime.StartDate;
                const endDate = busyTime.End || busyTime.endDate || busyTime.EndDate;

                if (startDate) {
                    // Try to safely parse the date
                    const parsedStartDate = new Date(startDate);
                    if (!isNaN(parsedStartDate.getTime())) {
                        startTime = format(parsedStartDate, 'HH:mm');
                    }
                }

                if (endDate) {
                    // Try to safely parse the date
                    const parsedEndDate = new Date(endDate);
                    if (!isNaN(parsedEndDate.getTime())) {
                        endTime = format(parsedEndDate, 'HH:mm');
                    }
                }
            } catch (error) {
                console.error(`Error parsing dates for room ${roomId}:`, error.message);
            }

            return {
                id: busyTime.id || busyTime.Id || `BKG-${Math.floor(Math.random() * 10000)}`,
                roomId: roomId,
                date: dateStr,
                startTime,
                endTime,
                title: busyTime.Title || busyTime.title || 'Booked'
            };
        });
    } catch (error) {
        console.error(chalk.red(`Error fetching busy times for room ${roomId}:`, error.message));
        // Return empty array in case of error
        return [];
    }
}

// Helper to create daily calendar view - now using API data
async function createDailyCalendarView(rooms, date, roomIdMapping, userBookings = []) {
    // Note: rooms are already sorted alphabetically by the fetchRooms function

    // Get the date string for filtering bookings
    const dateStr = format(date, 'yyyy-MM-dd');

    // Get current time to determine which slots are in the past
    const now = new Date();
    const isToday = format(now, 'yyyy-MM-dd') === dateStr;
    const currentHour = now.getHours();
    const currentMinute = now.getMinutes();
    const currentTimeInMinutes = currentHour * 60 + currentMinute;

    // User bookings are now pre-filtered and passed as parameter
    // This eliminates the need for an additional API call and filtering

    // Create header for hours (8:00 - 19:00)
    const roomColWidth = 20;
    const infoColWidth = 12;
    let output = chalk.bold('Room'.padEnd(roomColWidth)) + chalk.gray('│') + chalk.bold(' Info'.padEnd(infoColWidth));
    for (let hour = 8; hour < 19; hour++) {
        // Use shorter hour format (e.g., "08" instead of "08:00")
        const hourFormatted = hour.toString().padStart(2, '0');
        output += chalk.gray('│') + chalk.bold(hourFormatted.padEnd(4));
    }
    output += '\n';

    // Add separator line
    output += '─'.repeat(roomColWidth) + '┬' + '─'.repeat(infoColWidth);
    for (let i = 0; i < 11; i++) {
        output += '┼' + '─'.repeat(4);
    }
    output += '\n';

    // Fetch busy times for all rooms in parallel
    const bookingsPromises = rooms.map(room => fetchRoomBusyTimes(room.id, date, roomIdMapping));
    const roomsBookings = await Promise.all(bookingsPromises);

    // Process each room
    rooms.forEach((room, roomIndex) => {
        const roomColor = getRoomColor(room.id);

        // Format room information
        const roomName = chalk[roomColor](room.name.padEnd(roomColWidth));

        // Format capacity and price info
        const capacityStr = room.capacity !== 'N/A' ? `${room.capacity}p` : '?p';
        const priceStr = room.hourlyPrice !== 'N/A' ? room.hourlyPrice.replace(' €', '€') : '?€';
        const infoStr = chalk[roomColor](`${capacityStr} ${priceStr}`.padEnd(infoColWidth));

        // Get bookings for this room 
        const roomBookings = roomsBookings[roomIndex];

        // Convert bookings to a timeline of 15-minute slots (44 slots from 8:00 to 19:00)
        const timeline = Array(44).fill(null);

        // Mark booked slots in the timeline
        roomBookings.forEach(booking => {
            const startMinutes = timeToMinutes(booking.startTime);
            const endMinutes = timeToMinutes(booking.endTime);

            // Convert to slot indices (0-43)
            const startSlot = Math.floor((startMinutes - 8 * 60) / 15);
            const endSlot = Math.ceil((endMinutes - 8 * 60) / 15);

            // Mark slots as booked - we'll determine user ownership per slot
            for (let slot = Math.max(0, startSlot); slot < Math.min(44, endSlot); slot++) {
                timeline[slot] = { ...booking, isUserBooking: false };
            }
        });

        // Now check each slot to see if it's covered by a user booking
        for (let slot = 0; slot < 44; slot++) {
            if (timeline[slot]) {
                // Calculate the time for this slot
                const slotHour = 8 + Math.floor(slot / 4);
                const slotMinute = (slot % 4) * 15;
                const slotStartMinutes = slotHour * 60 + slotMinute;
                const slotEndMinutes = slotStartMinutes + 15;

                // Check if this specific 15-minute slot is covered by any user booking
                const isUserSlot = userBookings.some(userBooking => {
                    // Compare by room ID
                    if (userBooking.roomId !== room.apiId) {
                        return false;
                    }

                    // Check if this slot is within the user's booking time
                    const userStartMinutes = timeToMinutes(format(userBooking.startDate, 'HH:mm'));
                    const userEndMinutes = timeToMinutes(format(userBooking.endDate, 'HH:mm'));

                    return slotStartMinutes >= userStartMinutes && slotEndMinutes <= userEndMinutes;
                });

                timeline[slot].isUserBooking = isUserSlot;
            }
        }

        // Add room name and info to output
        output += roomName + chalk.gray('│') + infoStr;

        // Render the timeline for this room
        for (let hour = 0; hour < 11; hour++) {
            // Add vertical separator at the beginning of each hour
            output += chalk.gray('│');

            // Process 4 slots (15-min each) for this hour
            for (let quarterHour = 0; quarterHour < 4; quarterHour++) {
                const slot = hour * 4 + quarterHour;
                const booking = timeline[slot];

                // Calculate the time for this slot
                const slotHour = 8 + Math.floor(slot / 4);
                const slotMinute = (slot % 4) * 15;
                const slotTimeInMinutes = slotHour * 60 + slotMinute;

                // Check if this slot is in the past
                const isPastSlot = isToday && slotTimeInMinutes < currentTimeInMinutes;

                if (booking) {
                    // Choose the character based on whether it's the user's booking
                    const bookingChar = booking.isUserBooking ? '█' : '=';

                    // Grey out if in the past
                    if (isPastSlot) {
                        output += chalk.gray(bookingChar);
                    } else {
                        // Use bright/bold color for user's own bookings
                        if (booking.isUserBooking) {
                            output += chalk[roomColor].bold(bookingChar);
                        } else {
                            output += chalk[roomColor](bookingChar);
                        }
                    }
                } else {
                    // Use 1 space for available slots, or a grey dot for past slots
                    if (isPastSlot) {
                        output += chalk.gray('·');
                    } else {
                        output += ' ';
                    }
                }
            }
        }

        output += '\n';

        // Add separator after each room (except the last one)
        if (roomIndex < rooms.length - 1) {
            output += '─'.repeat(roomColWidth) + '┼' + '─'.repeat(infoColWidth);
            for (let i = 0; i < 11; i++) {
                output += '┼' + '─'.repeat(4);
            }
            output += '\n';
        }
    });

    return output;
}

export async function calendarView(options = {}) {
    try {

        // By default, use current date
        let currentDate = options.specificDate || new Date();

        // Check if running in interactive mode
        const isInteractive = options.interactive === true;

        const spinner = ora('Loading calendar view...').start();

        try {
            // Fetch rooms and user bookings in parallel for better performance
            const [rooms, userBookingsForDate] = await Promise.all([
                fetchRooms(),
                fetchUserBookings(currentDate) // Only fetch bookings for the specific date
            ]);

            // Create a ROOM_IDS mapping using the fetched rooms' apiId property
            const roomIdMapping = {};
            rooms.forEach(room => {
                roomIdMapping[room.id] = room.apiId;
            });

            // Display date information
            const dateDisplay = format(currentDate, 'EEEE, MMMM do, yyyy');

            // Create and display the calendar view using the mapping we just created
            const calendarOutput = await createDailyCalendarView(rooms, currentDate, roomIdMapping, userBookingsForDate);

            spinner.succeed('Calendar loaded');

            console.log(chalk.bold.blue(`\nDaily Calendar for ${dateDisplay}:`));
            console.log(chalk.gray('Each block represents a 15-minute slot. Colored blocks indicate bookings.'));
            console.log(chalk.gray('█ = Your bookings (bold), = = Other bookings, · = Past time slots'));
            console.log(calendarOutput);
        } catch (error) {
            spinner.fail('Failed to load calendar');
            console.error(chalk.red('Error:', error.message));
            if (!isInteractive) {
                process.exit(1);
            }
            return;
        }

        // If not in interactive mode, exit after showing calendar
        if (!isInteractive) {
            return;
        }

        // Navigation options with expanded choices
        const action = await select({
            message: 'What would you like to do?',
            choices: [
                { name: 'Refresh calendar', value: 'refresh' },
                { name: 'View next day', value: 'next' },
                { name: 'View previous day', value: 'prev' },
                { name: 'Go to specific date', value: 'goto-date' },
                { name: 'List all rooms', value: 'list' },
                { name: 'Book a room', value: 'book' },
                { name: 'Cancel a booking', value: 'cancel' },
                { name: 'View your bookings', value: 'my-bookings' },
                { name: 'View booking details', value: 'details' },
                { name: 'Exit', value: 'exit' }
            ]
        });

        if (action === 'refresh') {
            // Refresh the current calendar view
            console.log(chalk.blue('\n🔄 Refreshing calendar...'));
            return await calendarView({ specificDate: currentDate, interactive: true });

        } else if (action === 'next') {
            // Navigate to the next day
            const nextDate = addDays(currentDate, 1);
            return await calendarView({ specificDate: nextDate, interactive: true });

        } else if (action === 'prev') {
            // Navigate to the previous day
            const prevDate = addDays(currentDate, -1);
            return await calendarView({ specificDate: prevDate, interactive: true });

        } else if (action === 'list') {
            // Call the listRooms function
            await listRooms({ interactive: true });

            // Return to calendar view
            return await calendarView({ specificDate: currentDate, interactive: true });

        } else if (action === 'book') {
            // Call the bookRoom function
            await bookRoom({ interactive: true });

            // Return to calendar view
            return await calendarView({ specificDate: currentDate, interactive: true });

        } else if (action === 'cancel') {
            // Call the cancelBooking function
            await cancelBooking({ interactive: true });

            // Return to calendar view
            return await calendarView({ specificDate: currentDate, interactive: true });

        } else if (action === 'my-bookings') {
            // Call the myBookings function
            await myBookings({ interactive: true });

            // Return to calendar view
            return await calendarView({ specificDate: currentDate, interactive: true });

        } else if (action === 'details') {
            // Show detailed bookings for the day
            console.log(chalk.bold('\nBookings for today:'));

            // Fetch rooms from API again to ensure we have the latest data
            const rooms = await fetchRooms();

            // Create the room ID mapping again
            const roomIdMapping = {};
            rooms.forEach(room => {
                roomIdMapping[room.id] = room.apiId;
            });

            // Fetch bookings for all rooms for this date
            const dateStr = format(currentDate, 'yyyy-MM-dd');
            const bookingsPromises = rooms.map(room =>
                fetchRoomBusyTimes(room.id, currentDate, roomIdMapping)
            );
            const roomsBookings = await Promise.all(bookingsPromises);

            // Flatten the array of bookings
            const allBookings = roomsBookings.flat();

            if (allBookings.length === 0) {
                console.log(chalk.yellow('No bookings for today.'));
            } else {
                allBookings.forEach(booking => {
                    const room = rooms.find(r => r.id === booking.roomId);
                    const roomName = room ? room.name : 'Unknown Room';
                    console.log(
                        chalk[getRoomColor(room.id)](
                            `${roomName}: ${booking.startTime}-${booking.endTime} - ${booking.title}`
                        )
                    );
                });
            }

            // After showing details, prompt for next action
            const nextAction = await select({
                message: 'What would you like to do?',
                choices: [
                    { name: 'Return to calendar view', value: 'calendar' },
                    { name: 'Return to main menu', value: 'main' }
                ]
            });

            if (nextAction === 'calendar') {
                // Use current options to maintain the same day
                return await calendarView({ specificDate: currentDate, interactive: true });
            }
        } else if (action === 'goto-date') {
            // Prompt for a specific date in YYYY-MM-DD format
            const specificDate = await input({
                message: 'Enter date (YYYY-MM-DD):',
                default: format(new Date(), 'yyyy-MM-dd'),
                validate: input => {
                    // Validate the date format using regex
                    const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
                    if (!dateRegex.test(input)) {
                        return 'Please enter a valid date in YYYY-MM-DD format';
                    }

                    // Further validate that it's a valid date
                    try {
                        const parsedDate = parseISO(input);
                        if (isNaN(parsedDate.getTime())) {
                            return 'Please enter a valid date';
                        }
                        return true;
                    } catch (error) {
                        return 'Please enter a valid date';
                    }
                }
            });

            try {
                // Parse the date
                const parsedDate = parseISO(specificDate);

                // Navigate to the specified date
                return await calendarView({ specificDate: parsedDate, interactive: true });
            } catch (error) {
                console.error(chalk.red('Error parsing date:'), error.message);
                return await calendarView({ specificDate: currentDate, interactive: true });
            }
        } else if (action === 'exit') {
            console.log(chalk.blue('Goodbye!'));
            process.exit(0);
        }

    } catch (error) {
        console.error(chalk.red('Error displaying calendar:'), error.message);
        process.exit(1);
    }
}