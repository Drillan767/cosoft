import chalk from 'chalk';
import ora from 'ora';
import select from '@inquirer/select';
import confirm from '@inquirer/confirm';
import fs from 'fs';
import { API_CONFIG } from '../config/env.js';
import { makeApiRequest, checkAuthentication, getRefererHeaders } from '../config/api.js';
import { fetchUserBookings } from './myBookings.js';
import { showStandardMenu } from '../utils/sharedMenus.js';
// Import getRoomColor function for consistent colors
import { getRoomColor } from './calendarView.js';

export async function cancelBooking(options = {}) {

    // Check if running in interactive mode
    const isInteractive = options.interactive === true;

    // Check for batch processing options
    const batchCancellation = options.batchFile || options.batchJson || options.batchIds;

    if (batchCancellation) {
        return await processBatchCancellations(options, isInteractive);
    }

    // Check if booking ID is provided directly (non-interactive mode)
    const directCancellation = options.bookingId !== undefined;

    try {
        console.log(chalk.bold.blue('\nCancel a Booking'));

        // If booking ID is provided directly, cancel it without prompts
        if (directCancellation) {
            console.log(chalk.gray(`Cancelling booking with ID: ${options.bookingId}`));

            const spinner = ora('Cancelling booking...').start();

            try {
                await cancelBookingById(options.bookingId, spinner);
                return; // Exit after direct cancellation
            } catch (error) {
                spinner.fail('Cancellation failed');
                console.error(chalk.red(`Error: ${error.message}`));
                process.exit(1);
            }
        }

        const spinner = ora('Loading your bookings...').start();

        // Fetch user's bookings from the API
        let userBookings;
        try {
            userBookings = await fetchUserBookings();
            spinner.succeed('Bookings loaded');
        } catch (error) {
            spinner.fail('Failed to load bookings');
            console.error(chalk.red(error.message));
            if (!isInteractive) {
                process.exit(1);
            }
            return;
        }

        if (userBookings.length === 0) {
            console.log(chalk.yellow('You have no active bookings to cancel.'));

            if (!isInteractive) {
                return;
            }

            // If in interactive mode, show a menu to navigate back
            await showStandardMenu();
        }

        // Step 1: Select a booking to cancel
        const bookingId = await select({
            message: 'Select a booking to cancel:',
            choices: userBookings.map(booking => {
                const roomColor = getRoomColor(booking.roomId);
                return {
                    name: `${booking.id}: ${chalk[roomColor](booking.room)} on ${booking.date} at ${booking.time}`,
                    value: booking.id
                };
            })
        });

        const selectedBooking = userBookings.find(booking => booking.id === bookingId);
        const roomColor = getRoomColor(selectedBooking.roomId);

        // Step 2: Confirm cancellation
        const isConfirmed = await confirm({
            message: `Are you sure you want to cancel booking ${selectedBooking.id} for ${chalk[roomColor](selectedBooking.room)} on ${selectedBooking.date} at ${selectedBooking.time}?`,
            default: false
        });

        if (isConfirmed) {
            const spinner = ora('Cancelling booking...').start();

            try {
                await cancelBookingById(selectedBooking.id, spinner);
            } catch (error) {
                spinner.fail('Cancellation failed');
                console.error(chalk.red(`Error: ${error.message}`));
                if (error.response && error.response.data) {
                    console.error(chalk.red(`API response: ${JSON.stringify(error.response.data)}`));
                }

                if (!isInteractive) {
                    process.exit(1);
                }
            }
        } else {
            console.log(chalk.yellow('Cancellation aborted'));
        }

        // If not in interactive mode, exit after cancellation
        if (!isInteractive) {
            return;
        }

        // In interactive mode, show a menu
        await showStandardMenu();
    } catch (error) {
        console.error(chalk.red('Error during cancellation process:'), error.message);
        if (!isInteractive) {
            process.exit(1);
        }
    }
}

/**
 * Cancel a booking by its ID
 * @param {string} bookingId - The ID of the booking to cancel
 * @param {Object} spinner - Ora spinner object for showing progress
 * @returns {Promise<void>}
 */
async function cancelBookingById(bookingId, spinner) {
    // Call the API to cancel the booking using the correct endpoint
    const url = `${API_CONFIG.BASE_URL}/v2/api/api/Reservation/cancel-order`;
    const refererPath = `/v2/my-reservations`;

    await makeApiRequest({
        method: 'post',
        url: url,
        extraHeaders: getRefererHeaders(refererPath),
        data: {
            Id: bookingId
        }
    });

    spinner.succeed('Booking cancelled successfully');
    console.log(chalk.green(`Booking ${bookingId} has been cancelled`));
}

/**
 * Process batch cancellations from file, JSON string, or comma-separated IDs
 * @param {Object} options - Batch cancellation options
 * @param {boolean} isInteractive - Whether running in interactive mode
 */
async function processBatchCancellations(options, isInteractive) {
    try {
        let bookingIds = [];

        // Load booking IDs from different sources
        if (options.batchIds) {
            console.log(chalk.gray('Processing batch cancellations from comma-separated IDs...'));
            bookingIds = options.batchIds.split(',').map(id => id.trim()).filter(id => id.length > 0);
        } else if (options.batchFile) {
            console.log(chalk.gray(`Loading batch cancellations from file: ${options.batchFile}`));

            if (!fs.existsSync(options.batchFile)) {
                console.error(chalk.red(`Error: Batch file '${options.batchFile}' not found.`));
                process.exit(1);
            }

            const fileContent = fs.readFileSync(options.batchFile, 'utf8');
            try {
                const data = JSON.parse(fileContent);
                if (Array.isArray(data)) {
                    // Array of booking IDs
                    bookingIds = data;
                } else if (data.bookingIds && Array.isArray(data.bookingIds)) {
                    // Object with bookingIds property
                    bookingIds = data.bookingIds;
                } else {
                    throw new Error('File must contain an array of booking IDs or an object with a bookingIds property');
                }
            } catch (parseError) {
                console.error(chalk.red(`Error: Invalid JSON in batch file '${options.batchFile}'.`));
                console.error(chalk.red(parseError.message));
                process.exit(1);
            }
        } else if (options.batchJson) {
            console.log(chalk.gray('Processing batch cancellations from JSON string...'));

            try {
                const data = JSON.parse(options.batchJson);
                if (Array.isArray(data)) {
                    // Array of booking IDs
                    bookingIds = data;
                } else if (data.bookingIds && Array.isArray(data.bookingIds)) {
                    // Object with bookingIds property
                    bookingIds = data.bookingIds;
                } else {
                    throw new Error('JSON must contain an array of booking IDs or an object with a bookingIds property');
                }
            } catch (parseError) {
                console.error(chalk.red('Error: Invalid JSON in batch string.'));
                console.error(chalk.red(parseError.message));
                process.exit(1);
            }
        }

        if (bookingIds.length === 0) {
            console.log(chalk.yellow('No booking IDs found in batch data.'));
            return;
        }

        // Validate booking IDs
        const invalidIds = bookingIds.filter(id => typeof id !== 'string' || id.length === 0);
        if (invalidIds.length > 0) {
            console.error(chalk.red('Error: All booking IDs must be non-empty strings.'));
            console.error(chalk.red(`Invalid IDs found: ${invalidIds.length}`));
            process.exit(1);
        }

        console.log(chalk.blue(`\nProcessing ${bookingIds.length} cancellation(s)...`));

        // Optionally fetch user bookings to show details
        let userBookings = [];
        try {
            const spinner = ora('Loading your bookings for verification...').start();
            userBookings = await fetchUserBookings();
            spinner.succeed('Bookings loaded');
        } catch (error) {
            console.log(chalk.yellow('Warning: Could not load booking details for verification.'));
        }

        // Show confirmation if in interactive mode
        if (isInteractive) {
            console.log(chalk.yellow('\nBookings to be cancelled:'));
            bookingIds.forEach((id, index) => {
                const booking = userBookings.find(b => b.id === id);
                if (booking) {
                    const roomColor = getRoomColor(booking.roomId);
                    console.log(chalk.gray(`  ${index + 1}. ${chalk[roomColor](booking.room)} on ${booking.date} at ${booking.time} (ID: ${id})`));
                } else {
                    console.log(chalk.gray(`  ${index + 1}. Booking ID: ${id} (details not available)`));
                }
            });

            const isConfirmed = await confirm({
                message: `Are you sure you want to cancel all ${bookingIds.length} booking(s)?`,
                default: false
            });

            if (!isConfirmed) {
                console.log(chalk.yellow('Batch cancellation cancelled.'));
                return;
            }
        }

        // Process each cancellation
        const results = [];
        for (let i = 0; i < bookingIds.length; i++) {
            const bookingId = bookingIds[i];
            console.log(chalk.cyan(`\nProcessing cancellation ${i + 1}/${bookingIds.length}:`));
            console.log(chalk.gray(`  Booking ID: ${bookingId}`));

            try {
                const spinner = ora('Cancelling booking...').start();
                await cancelBookingById(bookingId, spinner);
                results.push({ bookingId, success: true });
                console.log(chalk.green(`  ✓ Cancellation completed successfully`));
            } catch (error) {
                results.push({ bookingId, success: false, error: error.message });
                console.log(chalk.red(`  ✗ Cancellation failed: ${error.message}`));
            }
        }

        // Summary
        console.log(chalk.bold('\n📊 Batch Cancellation Summary:'));
        const successful = results.filter(r => r.success).length;
        const failed = results.filter(r => !r.success).length;

        console.log(chalk.green(`✓ Successful: ${successful}`));
        if (failed > 0) {
            console.log(chalk.red(`✗ Failed: ${failed}`));

            console.log(chalk.yellow('\nFailed cancellations:'));
            results.filter(r => !r.success).forEach(result => {
                console.log(chalk.red(`  - ${result.bookingId}: ${result.error}`));
            });
        }

    } catch (error) {
        console.error(chalk.red('Error processing batch cancellations:'), error.message);
        if (!isInteractive) {
            process.exit(1);
        }
    }
}