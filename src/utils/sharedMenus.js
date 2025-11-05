import select from '@inquirer/select';
import chalk from 'chalk';

/**
 * Standard menu options for all views except calendar
 */
const STANDARD_MENU_CHOICES = [
    { name: 'Return to calendar view', value: 'calendar' },
    { name: 'Book a room', value: 'book' },
    { name: 'Cancel a booking', value: 'cancel' },
    { name: 'View your bookings', value: 'my-bookings' },
    { name: 'List all rooms', value: 'list' },
    { name: 'Exit', value: 'exit' }
];

/**
 * Shows the standard menu and handles navigation
 * @param {Object} options - Options for the menu
 * @param {string} options.message - Custom message for the menu (optional)
 * @param {Object} options.currentOptions - Current command options to pass through
 * @returns {Promise<void>} - Handles navigation or exits
 */
export async function showStandardMenu(options = {}) {
    const { message = 'What would you like to do?', currentOptions = {} } = options;

    const action = await select({
        message,
        choices: STANDARD_MENU_CHOICES
    });

    await handleStandardMenuAction(action, currentOptions);
}

/**
 * Handles the navigation logic for standard menu actions
 * @param {string} action - The selected action
 * @param {Object} currentOptions - Current command options to pass through
 * @returns {Promise<void>} - Handles navigation or exits
 */
export async function handleStandardMenuAction(action, currentOptions = {}) {
    if (action === 'calendar') {
        const { calendarView } = await import('../commands/calendarView.js');
        return await calendarView({ interactive: true });
    } else if (action === 'book') {
        const { bookRoom } = await import('../commands/bookRoom.js');
        return await bookRoom({ interactive: true });
    } else if (action === 'cancel') {
        const { cancelBooking } = await import('../commands/cancelBooking.js');
        return await cancelBooking({ interactive: true });
    } else if (action === 'my-bookings') {
        const { myBookings } = await import('../commands/myBookings.js');
        return await myBookings({ ...currentOptions, interactive: true });
    } else if (action === 'list') {
        const { listRooms } = await import('../commands/listRooms.js');
        return await listRooms({ interactive: true });
    } else if (action === 'exit') {
        console.log(chalk.blue('Goodbye!'));
        process.exit(0);
    }
}