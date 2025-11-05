import chalk from 'chalk';
import ora from 'ora';
import { format } from 'date-fns';
import select from '@inquirer/select';
import input from '@inquirer/input';
import confirm from '@inquirer/confirm';
import fs from 'fs';
import { API_CONFIG } from '../config/env.js';
import { makeApiRequest, checkAuthentication, getRefererHeaders } from '../config/api.js';
import { fetchRooms } from './listRooms.js';
import { showStandardMenu } from '../utils/sharedMenus.js';
// Import getRoomColor function for consistent colors
import { getRoomColor } from './calendarView.js';



export async function bookRoom(options = {}) {

    // Check if running in interactive mode
    const isInteractive = options.interactive === true;

    // Check for batch processing options
    const batchBooking = options.batchFile || options.batchJson;

    if (batchBooking) {
        return await processBatchBookings(options, isInteractive);
    }

    // For non-interactive mode, check if all required options are provided
    const nonInteractiveBooking = options.roomName && options.date && options.startTime && options.endTime;

    try {
        console.log(chalk.bold.blue('\nBook a Meeting Room'));
        const spinner = ora('Loading meeting rooms...').start();

        // Fetch rooms from the API
        let rooms;
        try {
            rooms = await fetchRooms();
            spinner.succeed('Meeting rooms loaded');
        } catch (error) {
            spinner.fail('Failed to load meeting rooms');
            console.error(chalk.red(error.message));
            if (!isInteractive) {
                process.exit(1);
            }
            return;
        }

        // Filter available rooms
        const availableRooms = rooms.filter(room => room.available);

        if (availableRooms.length === 0) {
            console.log(chalk.yellow('No available meeting rooms found.'));
            if (!isInteractive) {
                process.exit(1);
            }
            return;
        }

        let selectedRoom;
        let date, startTime, endTime;

        if (nonInteractiveBooking) {
            // Use command-line arguments
            selectedRoom = availableRooms.find(room =>
                room.name.toLowerCase() === options.roomName.toLowerCase()
            );

            if (!selectedRoom) {
                console.error(chalk.red(`Room "${options.roomName}" not found or not available.`));
                console.log(chalk.yellow('Available rooms:'));
                availableRooms.forEach(room => {
                    console.log(`- ${room.name}`);
                });
                process.exit(1);
            }

            date = options.date;
            startTime = options.startTime;
            endTime = options.endTime;

            // Validate date and time formats
            const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
            const timeRegex = /^([01]\d|2[0-3]):([0-5]\d)$/;

            if (!dateRegex.test(date)) {
                console.error(chalk.red('Invalid date format. Use YYYY-MM-DD format.'));
                process.exit(1);
            }

            if (!timeRegex.test(startTime) || !timeRegex.test(endTime)) {
                console.error(chalk.red('Invalid time format. Use HH:MM format.'));
                process.exit(1);
            }

            // Check if end time is after start time
            const [startHour, startMinute] = startTime.split(':').map(Number);
            const [endHour, endMinute] = endTime.split(':').map(Number);

            if (endHour < startHour || (endHour === startHour && endMinute <= startMinute)) {
                console.error(chalk.red('End time must be after start time.'));
                process.exit(1);
            }

            console.log(chalk.gray(`Using room: ${selectedRoom.name}`));
            console.log(chalk.gray(`Date: ${date}`));
            console.log(chalk.gray(`Time: ${startTime} - ${endTime}`));
        } else {
            // Step 1: Select a room interactively
            const roomId = await select({
                message: 'Select a meeting room:',
                choices: availableRooms.map(room => ({
                    name: `${room.name} (Capacity: ${room.capacity}, Floor: ${room.floor})`,
                    value: room.id
                }))
            });

            selectedRoom = rooms.find(room => room.id === roomId);

            // Step 2: Enter booking details
            date = await input({
                message: 'Enter date (YYYY-MM-DD):',
                default: format(new Date(), 'yyyy-MM-dd'),
                validate: input => {
                    const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
                    return dateRegex.test(input) ? true : 'Please enter a valid date in YYYY-MM-DD format';
                }
            });

            startTime = await input({
                message: 'Enter start time (HH:MM):',
                validate: input => {
                    const timeRegex = /^([01]\d|2[0-3]):([0-5]\d)$/;
                    return timeRegex.test(input) ? true : 'Please enter a valid time in HH:MM format';
                }
            });

            endTime = await input({
                message: 'Enter end time (HH:MM):',
                validate: input => {
                    const timeRegex = /^([01]\d|2[0-3]):([0-5]\d)$/;
                    if (!timeRegex.test(input)) {
                        return 'Please enter a valid time in HH:MM format';
                    }

                    // Check if end time is after start time
                    const [startHour, startMinute] = startTime.split(':').map(Number);
                    const [endHour, endMinute] = input.split(':').map(Number);

                    if (endHour < startHour || (endHour === startHour && endMinute <= startMinute)) {
                        return 'End time must be after start time';
                    }

                    return true;
                }
            });
        }

        // Enable debug mode for conflict detection
        if (!process.env.DEBUG_BOOKING) {
            console.log(chalk.gray('💡 Tip: Set DEBUG_BOOKING=1 to see detailed API responses if booking fails'));
        }

        // Generate a cart ID
        const cartId = generateCartId();

        // Step 4: Confirm booking (only in interactive mode)
        const isConfirmed = nonInteractiveBooking ? true : await confirm({
            message: `Confirm booking for ${selectedRoom.name} on ${date} from ${startTime} to ${endTime}?`,
            default: true
        });

        if (isConfirmed) {
            const spinner = ora('Processing booking...').start();

            try {
                // Convert date and time to API format
                const startDateTime = `${date}T${startTime}:00`;
                const endDateTime = `${date}T${endTime}:00`;

                // Step 1: Add to cart
                const cartUrl = `${API_CONFIG.BASE_URL}/v2/api/api/Cart/%7BSkipVisibleFilter%7D`;
                const cartRefererPath = `/v2/new-reservation/${API_CONFIG.COWORKING_SPACE_ID}/${API_CONFIG.CATEGORY_ID}/${selectedRoom.apiId}`;

                const cartData = {
                    orders: [{
                        coworkingSpaceId: API_CONFIG.COWORKING_SPACE_ID,
                        categoryId: API_CONFIG.CATEGORY_ID,
                        itemId: selectedRoom.apiId,
                        startenddate_: {
                            date: startDateTime,
                            times: [{
                                start: startTime,
                                end: endTime
                            }]
                        },
                        startenddate: [{
                            start: startDateTime,
                            end: endDateTime,
                            type: "hour",
                            timeSlotId: null,
                            id: selectedRoom.apiId
                        }],
                        cartId: cartId
                    }]
                };

                spinner.text = 'Adding to cart...';
                const cartResponse = await makeApiRequest({
                    method: 'post',
                    url: cartUrl,
                    extraHeaders: getRefererHeaders(cartRefererPath),
                    data: cartData
                });

                // Enhanced debug logging to understand API response structure
                if (process.env.DEBUG_BOOKING) {
                    console.log(chalk.gray('\n=== CART API DEBUG INFO ==='));
                    console.log(chalk.gray('Request URL:'), cartUrl);
                    console.log(chalk.gray('Request Data:'), JSON.stringify(cartData, null, 2));
                    console.log(chalk.gray('Cart API Response:'), JSON.stringify(cartResponse, null, 2));
                    console.log(chalk.gray('=== END CART DEBUG ===\n'));
                } else {
                    // Always show key indicators even without debug mode
                    if (cartResponse) {
                        console.log(chalk.gray(`Cart Response - HasError: ${!!cartResponse.CartHasError}, Items: ${cartResponse.ItemsCategory?.length || 0}, Total: ${cartResponse.Total?.EuroTTC || 0}€`));
                    }
                }

                // Enhanced error detection for cart response
                if (!cartResponse) {
                    throw new Error('No response received from booking API');
                }

                // Check for missing cart items first (indicates booking was rejected)
                if (!cartResponse.ItemsCategory || cartResponse.ItemsCategory.length === 0) {
                    throw new Error('⚠️  Room could not be added to cart. This usually indicates the room is unavailable or the time slot is invalid.');
                }

                // Analyze cart items for specific failure types first (more specific than CartHasError)
                if (cartResponse.ItemsCategory && cartResponse.ItemsCategory.length > 0) {
                    const schedulingConflicts = [];
                    const validationFailures = [];

                    cartResponse.ItemsCategory.forEach(category => {
                        if (category.Items) {
                            category.Items.forEach(item => {
                                // Type 1: Scheduling conflict (room already booked)
                                if (item.HasAlreadyOrdered === true) {
                                    schedulingConflicts.push(item.ItemName || 'Unknown room');
                                }
                                // Type 2: Validation failure (time passed, room disabled, etc.)
                                else if (item.DisabledItem === true) {
                                    const reason = item.InfoMessage === 'cart.unavailable.passed' ?
                                        'the requested time has already passed' :
                                        'the room is currently unavailable';
                                    validationFailures.push(`${item.ItemName || 'Unknown room'} (${reason})`);
                                }
                            });
                        }
                    });

                    // Prioritize scheduling conflicts as they're more specific
                    if (schedulingConflicts.length > 0) {
                        throw new Error(`🚫 Scheduling conflict: The room "${schedulingConflicts.join(', ')}" is already booked during your requested time slot. Please choose a different time or room.`);
                    }

                    if (validationFailures.length > 0) {
                        throw new Error(`❌ Booking validation failed: ${validationFailures.join(', ')}. Please check your booking details and try again.`);
                    }

                    // Generic check for any blocked/unavailable items
                    const invalidItems = cartResponse.ItemsCategory.filter(item =>
                        !item.Items || item.Items.length === 0 ||
                        item.Items.some(subItem => subItem.IsBlocked || subItem.IsUnavailable)
                    );

                    if (invalidItems.length > 0) {
                        throw new Error('⚠️  Room booking was rejected due to unavailability. Please try a different time slot or room.');
                    }
                }

                // Check for explicit API errors after specific item analysis
                if (cartResponse.CartHasError || cartResponse.Error || cartResponse.ErrorMessage) {
                    const errorMsg = cartResponse.Error || cartResponse.ErrorMessage || 'API validation failed';
                    throw new Error(`⚠️  Booking validation failed: ${errorMsg}. Please check your booking details (date, time, room availability).`);
                }

                // Check for zero total which might indicate silent failure
                if (cartResponse.Total && cartResponse.Total.EuroTTC === 0 && cartResponse.Total.Credits === 0) {
                    // Additional check: see if there are actually valid bookable items
                    const hasValidItems = cartResponse.ItemsCategory &&
                        cartResponse.ItemsCategory.some(item =>
                            item.Items && item.Items.length > 0 &&
                            item.Items.some(subItem => !subItem.IsBlocked && !subItem.IsUnavailable)
                        );

                    if (!hasValidItems) {
                        throw new Error('Booking failed - room appears to be unavailable for the requested time slot');
                    }

                    console.log(chalk.yellow('⚠️  Warning: Cart total is 0, which may indicate a booking conflict or free booking'));
                }

                // Extract the booking price
                let bookingPrice = 'N/A';
                if (cartResponse.Total && cartResponse.Total.EuroTTC) {
                    bookingPrice = `${cartResponse.Total.EuroTTC.toFixed(2)} €`;
                } else if (cartResponse.Total && cartResponse.Total.Credits) {
                    bookingPrice = `${cartResponse.Total.Credits.toFixed(2)} credits`;
                }

                spinner.text = `Processing payment (${bookingPrice})...`;

                // Step 2: Make the payment
                const paymentUrl = `${API_CONFIG.BASE_URL}/v2/api/api/Payment/pay`;
                const paymentRefererPath = `/v2/cart/validate`;

                const paymentData = {
                    isUser: true,
                    isPerson: true,
                    isVatRequired: true,
                    isStatusRequired: true,
                    cgv: true,
                    societyname: "",
                    societyvat: "",
                    societysiret: "",
                    societystatus: "",
                    firstname: "",
                    lastname: "",
                    address: "",
                    city: "",
                    zipCode: "",
                    phone: "",
                    email: "",
                    cart: [{
                        coworkingSpaceId: API_CONFIG.COWORKING_SPACE_ID,
                        categoryId: API_CONFIG.CATEGORY_ID,
                        itemId: selectedRoom.apiId,
                        startenddate_: {
                            date: startDateTime,
                            times: [{
                                start: startTime,
                                end: endTime
                            }]
                        },
                        startenddate: [{
                            start: startDateTime,
                            end: endDateTime,
                            type: "hour",
                            timeSlotId: null,
                            id: selectedRoom.apiId
                        }],
                        cartId: cartId
                    }],
                    paymentType: "credit"
                };

                const paymentResponse = await makeApiRequest({
                    method: 'post',
                    url: paymentUrl,
                    extraHeaders: getRefererHeaders(paymentRefererPath),
                    data: paymentData
                });

                // Enhanced debug logging for payment response
                if (process.env.DEBUG_BOOKING) {
                    console.log(chalk.gray('\n=== PAYMENT API DEBUG INFO ==='));
                    console.log(chalk.gray('Request URL:'), paymentUrl);
                    console.log(chalk.gray('Payment API Response:'), JSON.stringify(paymentResponse, null, 2));
                    console.log(chalk.gray('=== END PAYMENT DEBUG ===\n'));
                } else {
                    // Always show key indicators even without debug mode
                    if (paymentResponse) {
                        console.log(chalk.gray(`Payment Response - RedirectUrl: ${!!paymentResponse.RedirectUrl}, Error: ${paymentResponse.Error || 'none'}`));
                    }
                }

                // Enhanced error detection for payment response
                if (!paymentResponse) {
                    throw new Error('No response received from payment API');
                }

                // Check for explicit error messages
                if (paymentResponse.Error || paymentResponse.ErrorMessage) {
                    const errorMsg = paymentResponse.Error || paymentResponse.ErrorMessage;
                    throw new Error(`Payment failed: ${errorMsg}`);
                }

                // Check for failed payment status indicators
                if (paymentResponse.Status === 'Failed' || paymentResponse.status === 'failed') {
                    throw new Error('Payment was declined - booking was not completed');
                }

                // Check for booking-specific failures
                if (paymentResponse.BookingFailed || paymentResponse.bookingFailed) {
                    throw new Error('Booking was rejected during payment processing - likely due to room unavailability');
                }

                // Traditional check for redirect URL
                if (!paymentResponse.RedirectUrl) {
                    // Check for other success indicators
                    if (!paymentResponse.Success && !paymentResponse.IsSuccess && !paymentResponse.success) {
                        // If no success indicator and no redirect, this might be a silent failure
                        console.log(chalk.yellow('⚠️  Warning: Payment completed but no redirect URL received - booking status uncertain'));
                        console.log(chalk.yellow('Please verify your booking was created by checking "My Bookings"'));
                    }
                }

                // Additional validation: if we have a redirect URL, it should not be an error page
                if (paymentResponse.RedirectUrl && (
                    paymentResponse.RedirectUrl.includes('error') ||
                    paymentResponse.RedirectUrl.includes('failed') ||
                    paymentResponse.RedirectUrl.includes('declined')
                )) {
                    throw new Error('Payment was redirected to an error page - booking failed');
                }

                spinner.succeed('Booking confirmed!');

                // Use consistent room coloring for the room name
                const roomColor = getRoomColor(selectedRoom.id);

                console.log(chalk.green.bold('\nBooking Details:'));
                console.log(`Room: ${chalk[roomColor](selectedRoom.name)}`);
                console.log(`Date: ${date}`);
                console.log(`Time: ${startTime} - ${endTime}`);
                console.log(`Price: ${bookingPrice}`);
                console.log(chalk.green(`\nBooking completed successfully!`));

            } catch (error) {
                spinner.fail('Booking failed');
                console.error(chalk.red(`Error: ${error.message}`));
                if (error.response && error.response.data) {
                    console.error(chalk.red(`API response: ${JSON.stringify(error.response.data)}`));
                }
            }
        } else {
            console.log(chalk.yellow('Booking cancelled'));
        }

        // If not in interactive mode, exit after booking
        if (!isInteractive) {
            return;
        }

        // In interactive mode, show a menu
        await showStandardMenu();

    } catch (error) {
        console.error(chalk.red('Error during booking process:'), error.message);
        if (!isInteractive) {
            process.exit(1);
        }
    }
}

/**
 * Generate a random cart ID
 * @returns {string} A random cart ID
 */
function generateCartId() {
    // Generate a random string for cart ID
    return Math.random().toString(36).substring(2, 12);
}

/**
 * Process batch bookings from file or JSON string
 * @param {Object} options - Batch booking options
 * @param {boolean} isInteractive - Whether running in interactive mode
 */
async function processBatchBookings(options, isInteractive) {
    try {
        let bookingData = [];

        // Load booking data from file or JSON string
        if (options.batchFile) {
            console.log(chalk.gray(`Loading batch bookings from file: ${options.batchFile}`));

            if (!fs.existsSync(options.batchFile)) {
                console.error(chalk.red(`Error: Batch file '${options.batchFile}' not found.`));
                process.exit(1);
            }

            const fileContent = fs.readFileSync(options.batchFile, 'utf8');
            try {
                bookingData = JSON.parse(fileContent);
            } catch (parseError) {
                console.error(chalk.red(`Error: Invalid JSON in batch file '${options.batchFile}'.`));
                console.error(chalk.red(parseError.message));
                process.exit(1);
            }
        } else if (options.batchJson) {
            console.log(chalk.gray('Processing batch bookings from JSON string...'));

            try {
                bookingData = JSON.parse(options.batchJson);
            } catch (parseError) {
                console.error(chalk.red('Error: Invalid JSON in batch string.'));
                console.error(chalk.red(parseError.message));
                process.exit(1);
            }
        }

        if (!Array.isArray(bookingData)) {
            console.error(chalk.red('Error: Batch data must be an array of booking objects.'));
            process.exit(1);
        }

        if (bookingData.length === 0) {
            console.log(chalk.yellow('No bookings found in batch data.'));
            return;
        }

        console.log(chalk.blue(`\nProcessing ${bookingData.length} booking(s)...`));

        // Validate all booking entries first
        const validationErrors = [];
        bookingData.forEach((booking, index) => {
            const errors = validateBookingData(booking, index);
            validationErrors.push(...errors);
        });

        if (validationErrors.length > 0) {
            console.error(chalk.red('\nValidation errors found:'));
            validationErrors.forEach(error => console.error(chalk.red(`  - ${error}`)));
            process.exit(1);
        }

        // Fetch rooms once for all bookings
        const spinner = ora('Loading meeting rooms...').start();
        let rooms;
        try {
            rooms = await fetchRooms();
            spinner.succeed('Meeting rooms loaded');
        } catch (error) {
            spinner.fail('Failed to load meeting rooms');
            console.error(chalk.red(error.message));
            process.exit(1);
        }

        const availableRooms = rooms.filter(room => room.available);

        // Process each booking
        const results = [];
        for (let i = 0; i < bookingData.length; i++) {
            const booking = bookingData[i];
            console.log(chalk.cyan(`\nProcessing booking ${i + 1}/${bookingData.length}:`));
            console.log(chalk.gray(`  Room: ${booking.roomName}`));
            console.log(chalk.gray(`  Date: ${booking.date}`));
            console.log(chalk.gray(`  Time: ${booking.startTime} - ${booking.endTime}`));

            try {
                const result = await processSingleBooking(booking, availableRooms);
                results.push({ ...booking, success: true, result });
                console.log(chalk.green(`  ✓ Booking completed successfully`));
            } catch (error) {
                results.push({ ...booking, success: false, error: error.message });
                console.log(chalk.red(`  ✗ Booking failed: ${error.message}`));
            }
        }

        // Summary
        console.log(chalk.bold('\n📊 Batch Booking Summary:'));
        const successful = results.filter(r => r.success).length;
        const failed = results.filter(r => !r.success).length;

        console.log(chalk.green(`✓ Successful: ${successful}`));
        if (failed > 0) {
            console.log(chalk.red(`✗ Failed: ${failed}`));

            console.log(chalk.yellow('\nFailed bookings:'));
            results.filter(r => !r.success).forEach(booking => {
                console.log(chalk.red(`  - ${booking.roomName} on ${booking.date} at ${booking.startTime}-${booking.endTime}: ${booking.error}`));
            });
        }

        const totalCost = results
            .filter(r => r.success && r.result.price)
            .reduce((sum, r) => sum + parseFloat(r.result.price.replace(' €', '')), 0);

        if (totalCost > 0) {
            console.log(chalk.blue(`💰 Total cost: ${totalCost.toFixed(2)} €`));
        }

    } catch (error) {
        console.error(chalk.red('Error processing batch bookings:'), error.message);
        if (!isInteractive) {
            process.exit(1);
        }
    }
}

/**
 * Validate booking data object
 * @param {Object} booking - Booking data to validate
 * @param {number} index - Index of the booking in the batch
 * @returns {Array} Array of validation error messages
 */
function validateBookingData(booking, index) {
    const errors = [];
    const prefix = `Booking ${index + 1}`;

    if (!booking.roomName) {
        errors.push(`${prefix}: roomName is required`);
    }

    if (!booking.date) {
        errors.push(`${prefix}: date is required`);
    } else {
        const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
        if (!dateRegex.test(booking.date)) {
            errors.push(`${prefix}: date must be in YYYY-MM-DD format`);
        }
    }

    if (!booking.startTime) {
        errors.push(`${prefix}: startTime is required`);
    } else {
        const timeRegex = /^([01]\d|2[0-3]):([0-5]\d)$/;
        if (!timeRegex.test(booking.startTime)) {
            errors.push(`${prefix}: startTime must be in HH:MM format`);
        }
    }

    if (!booking.endTime) {
        errors.push(`${prefix}: endTime is required`);
    } else {
        const timeRegex = /^([01]\d|2[0-3]):([0-5]\d)$/;
        if (!timeRegex.test(booking.endTime)) {
            errors.push(`${prefix}: endTime must be in HH:MM format`);
        }

        // Check if end time is after start time
        if (booking.startTime && booking.endTime) {
            const [startHour, startMinute] = booking.startTime.split(':').map(Number);
            const [endHour, endMinute] = booking.endTime.split(':').map(Number);

            if (endHour < startHour || (endHour === startHour && endMinute <= startMinute)) {
                errors.push(`${prefix}: endTime must be after startTime`);
            }
        }
    }

    return errors;
}

/**
 * Process a single booking from batch data
 * @param {Object} booking - Booking data
 * @param {Array} availableRooms - Array of available rooms
 * @returns {Promise<Object>} Booking result
 */
async function processSingleBooking(booking, availableRooms) {
    // Find the room
    const selectedRoom = availableRooms.find(room =>
        room.name.toLowerCase() === booking.roomName.toLowerCase()
    );

    if (!selectedRoom) {
        throw new Error(`Room "${booking.roomName}" not found or not available`);
    }

    // Generate a cart ID
    const cartId = generateCartId();

    // Convert date and time to API format
    const startDateTime = `${booking.date}T${booking.startTime}:00`;
    const endDateTime = `${booking.date}T${booking.endTime}:00`;

    // Step 1: Add to cart
    const cartUrl = `${API_CONFIG.BASE_URL}/v2/api/api/Cart/%7BSkipVisibleFilter%7D`;
    const cartRefererPath = `/v2/new-reservation/${API_CONFIG.COWORKING_SPACE_ID}/${API_CONFIG.CATEGORY_ID}/${selectedRoom.apiId}`;

    const cartData = {
        orders: [{
            coworkingSpaceId: API_CONFIG.COWORKING_SPACE_ID,
            categoryId: API_CONFIG.CATEGORY_ID,
            itemId: selectedRoom.apiId,
            startenddate_: {
                date: startDateTime,
                times: [{
                    start: booking.startTime,
                    end: booking.endTime
                }]
            },
            startenddate: [{
                start: startDateTime,
                end: endDateTime,
                type: "hour",
                timeSlotId: null,
                id: selectedRoom.apiId
            }],
            cartId: cartId
        }]
    };

    const cartResponse = await makeApiRequest({
        method: 'post',
        url: cartUrl,
        extraHeaders: getRefererHeaders(cartRefererPath),
        data: cartData
    });

    // Apply the same enhanced error detection as the main booking function
    if (!cartResponse) {
        throw new Error('No response received from booking API');
    }

    if (!cartResponse.ItemsCategory || cartResponse.ItemsCategory.length === 0) {
        throw new Error('⚠️  Room could not be added to cart. This usually indicates the room is unavailable.');
    }

    // Analyze cart items for specific failure types first
    if (cartResponse.ItemsCategory && cartResponse.ItemsCategory.length > 0) {
        const schedulingConflicts = [];
        const validationFailures = [];

        cartResponse.ItemsCategory.forEach(category => {
            if (category.Items) {
                category.Items.forEach(item => {
                    // Type 1: Scheduling conflict (room already booked)
                    if (item.HasAlreadyOrdered === true) {
                        schedulingConflicts.push(item.ItemName || 'Unknown room');
                    }
                    // Type 2: Validation failure (time passed, room disabled, etc.)
                    else if (item.DisabledItem === true) {
                        const reason = item.InfoMessage === 'cart.unavailable.passed' ?
                            'the requested time has already passed' :
                            'the room is currently unavailable';
                        validationFailures.push(`${item.ItemName || 'Unknown room'} (${reason})`);
                    }
                });
            }
        });

        // Prioritize scheduling conflicts as they're more specific
        if (schedulingConflicts.length > 0) {
            throw new Error(`🚫 Scheduling conflict: The room "${schedulingConflicts.join(', ')}" is already booked during your requested time slot.`);
        }

        if (validationFailures.length > 0) {
            throw new Error(`❌ Booking validation failed: ${validationFailures.join(', ')}.`);
        }

        // Generic check for any blocked/unavailable items
        const invalidItems = cartResponse.ItemsCategory.filter(item =>
            !item.Items || item.Items.length === 0 ||
            item.Items.some(subItem => subItem.IsBlocked || subItem.IsUnavailable)
        );

        if (invalidItems.length > 0) {
            throw new Error('⚠️  Room booking was rejected due to unavailability.');
        }
    }

    // Check for explicit API errors after specific item analysis
    if (cartResponse.CartHasError || cartResponse.Error || cartResponse.ErrorMessage) {
        const errorMsg = cartResponse.Error || cartResponse.ErrorMessage || 'API validation failed';
        throw new Error(`⚠️  Booking validation failed: ${errorMsg}. Please check your booking details.`);
    }

    // Extract the booking price
    let bookingPrice = 'N/A';
    if (cartResponse.Total && cartResponse.Total.EuroTTC) {
        bookingPrice = `${cartResponse.Total.EuroTTC.toFixed(2)} €`;
    } else if (cartResponse.Total && cartResponse.Total.Credits) {
        bookingPrice = `${cartResponse.Total.Credits.toFixed(2)} credits`;
    }

    // Step 2: Make the payment
    const paymentUrl = `${API_CONFIG.BASE_URL}/v2/api/api/Payment/pay`;
    const paymentRefererPath = `/v2/cart/validate`;

    const paymentData = {
        isUser: true,
        isPerson: true,
        isVatRequired: true,
        isStatusRequired: true,
        cgv: true,
        societyname: "",
        societyvat: "",
        societysiret: "",
        societystatus: "",
        firstname: "",
        lastname: "",
        address: "",
        city: "",
        zipCode: "",
        phone: "",
        email: "",
        cart: [{
            coworkingSpaceId: API_CONFIG.COWORKING_SPACE_ID,
            categoryId: API_CONFIG.CATEGORY_ID,
            itemId: selectedRoom.apiId,
            startenddate_: {
                date: startDateTime,
                times: [{
                    start: booking.startTime,
                    end: booking.endTime
                }]
            },
            startenddate: [{
                start: startDateTime,
                end: endDateTime,
                type: "hour",
                timeSlotId: null,
                id: selectedRoom.apiId
            }],
            cartId: cartId
        }],
        paymentType: "credit"
    };

    const paymentResponse = await makeApiRequest({
        method: 'post',
        url: paymentUrl,
        extraHeaders: getRefererHeaders(paymentRefererPath),
        data: paymentData
    });

    // Apply the same enhanced error detection as the main booking function
    if (!paymentResponse) {
        throw new Error('No response received from payment API');
    }

    if (paymentResponse.Error || paymentResponse.ErrorMessage) {
        const errorMsg = paymentResponse.Error || paymentResponse.ErrorMessage;
        throw new Error(`Payment failed: ${errorMsg}`);
    }

    if (paymentResponse.Status === 'Failed' || paymentResponse.status === 'failed') {
        throw new Error('Payment was declined - booking was not completed');
    }

    if (paymentResponse.BookingFailed || paymentResponse.bookingFailed) {
        throw new Error('Booking was rejected during payment processing - likely due to room unavailability');
    }

    if (!paymentResponse.RedirectUrl) {
        if (!paymentResponse.Success && !paymentResponse.IsSuccess && !paymentResponse.success) {
            throw new Error('Payment processing failed - booking was not completed');
        }
    }

    if (paymentResponse.RedirectUrl && (
        paymentResponse.RedirectUrl.includes('error') ||
        paymentResponse.RedirectUrl.includes('failed') ||
        paymentResponse.RedirectUrl.includes('declined')
    )) {
        throw new Error('Payment was redirected to an error page - booking failed');
    }

    return {
        room: selectedRoom.name,
        date: booking.date,
        startTime: booking.startTime,
        endTime: booking.endTime,
        price: bookingPrice
    };
}