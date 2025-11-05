#!/usr/bin/env node
import { Command } from 'commander';
import chalk from 'chalk';
import { calendarView } from './commands/calendarView.js';
import { bookRoom } from './commands/bookRoom.js';
import { listRooms } from './commands/listRooms.js';
import { myBookings } from './commands/myBookings.js';
import { cancelBooking } from './commands/cancelBooking.js';
import { checkAuthentication } from './config/api.js';
import { loginCommand } from './commands/login.js';
import { mcpServerCommand } from './commands/mcpServer.js';
import { parseISO } from 'date-fns';

const program = new Command();

// Store the interactive flag at the program level
program
    .name('cosoft')
    .description('CLI tool to interact with the CoSoft API for meeting room reservations')
    .version('1.0.0')
    .option('-i, --interactive', 'Run in interactive mode')
    .allowUnknownOption()
    .action(async (options) => {
        // Check if there are unknown arguments that look like commands
        const args = process.argv.slice(2);
        const knownCommands = ['calendar', 'book', 'list', 'my-bookings', 'cancel', 'auth', 'login', 'mcp-server', 'help'];
        const knownOptions = ['-i', '--interactive', '-V', '--version', '-h', '--help'];

        // Filter out known options and their values
        const filteredArgs = args.filter(arg => {
            if (knownOptions.includes(arg)) return false;
            if (arg.startsWith('--') && knownOptions.some(opt => arg.startsWith(opt + '='))) return false;
            return true;
        });

        // If there are remaining args that aren't known commands, they are unknown commands
        if (filteredArgs.length > 0 && !knownCommands.includes(filteredArgs[0])) {
            console.error(chalk.red.bold(`\nError: Unknown command '${filteredArgs[0]}'`));
            console.log(chalk.yellow('Use "cosoft --help" to see the list of available commands.'));
            process.exit(1);
        }

        // If no command is specified, run calendar view in interactive mode
        // Check authentication first
        const isAuthenticated = await checkAuthentication();
        if (!isAuthenticated) {
            process.exit(1);
        }
        calendarView({ interactive: true });
    });

// Calendar view command
program
    .command('calendar')
    .description('Show the calendar view for meeting rooms')
    .argument('[date]', 'Date in YYYY-MM-DD format')
    .action(async (date, cmdObj) => {
        // Check authentication once at CLI level
        const isAuthenticated = await checkAuthentication();
        if (!isAuthenticated) {
            process.exit(1);
        }

        // Get the interactive flag from the root command
        const isInteractive = program.opts().interactive || false;
        let specificDate = date ? parseISO(date) : new Date();
        calendarView({ specificDate, interactive: isInteractive });
    });

// Book room command
program
    .command('book')
    .description('Book a meeting room or batch of rooms')
    .option('--room-name <name>', 'Name of the room to book (e.g., "CALL BOX 3")')
    .option('--date <date>', 'Date in YYYY-MM-DD format')
    .option('--start-time <time>', 'Start time in HH:MM format')
    .option('--end-time <time>', 'End time in HH:MM format')
    .option('--batch-file <path>', 'Path to JSON file containing batch booking data')
    .option('--batch-json <json>', 'JSON string containing batch booking data')
    .action(async (options) => {
        // Check authentication once at CLI level
        const isAuthenticated = await checkAuthentication();
        if (!isAuthenticated) {
            process.exit(1);
        }

        // Get the interactive flag from the root command
        const isInteractive = program.opts().interactive || false;
        bookRoom({
            interactive: isInteractive,
            roomName: options.roomName,
            date: options.date,
            startTime: options.startTime,
            endTime: options.endTime,
            batchFile: options.batchFile,
            batchJson: options.batchJson
        });
    });

// List rooms command
program
    .command('list')
    .description('List all available meeting rooms')
    .action(async (cmdObj) => {
        // Check authentication once at CLI level
        const isAuthenticated = await checkAuthentication();
        if (!isAuthenticated) {
            process.exit(1);
        }

        // Get the interactive flag from the root command
        const isInteractive = program.opts().interactive || false;
        listRooms({ interactive: isInteractive });
    });

// My bookings command
program
    .command('my-bookings')
    .description('View your bookings')
    .argument('[date]', 'Date in YYYY-MM-DD format (optional)')
    .action(async (date, cmdObj) => {
        // Check authentication once at CLI level
        const isAuthenticated = await checkAuthentication();
        if (!isAuthenticated) {
            process.exit(1);
        }

        // Get the interactive flag from the root command
        const isInteractive = program.opts().interactive || false;
        const bookingOptions = {
            interactive: isInteractive
        };
        if (date) {
            bookingOptions.date = date;
        }
        myBookings(bookingOptions);
    });

// Cancel booking command
program
    .command('cancel')
    .description('Cancel an existing booking or batch of bookings')
    .option('--booking-id <id>', 'ID of the booking to cancel')
    .option('--batch-file <path>', 'Path to JSON file containing batch cancellation data')
    .option('--batch-json <json>', 'JSON string containing batch cancellation data')
    .option('--batch-ids <ids>', 'Comma-separated list of booking IDs to cancel')
    .action(async (options) => {
        // Check authentication once at CLI level
        const isAuthenticated = await checkAuthentication();
        if (!isAuthenticated) {
            process.exit(1);
        }

        // Get the interactive flag from the root command
        const isInteractive = program.opts().interactive || false;
        cancelBooking({
            interactive: isInteractive,
            bookingId: options.bookingId,
            batchFile: options.batchFile,
            batchJson: options.batchJson,
            batchIds: options.batchIds
        });
    });

// Authentication status command
program
    .command('auth')
    .description('Check authentication status and display logged-in user')
    .action(async (cmdObj) => {
        // Run only the authentication check
        const isAuthenticated = await checkAuthentication();
        if (!isAuthenticated) {
            process.exit(1);
        }
        console.log(chalk.blue('Authentication check completed.'));
    });

// Login command
program
    .command('login')
    .description('Login to CoSoft and save authentication tokens')
    .action(async () => {
        // Run the login command
        await loginCommand();
    });

// MCP Server command
program
    .command('mcp-server')
    .description('Start the CoSoft CLI as an MCP (Model Context Protocol) server')
    .action(async () => {
        // Run the MCP server
        await mcpServerCommand();
    });

// Add custom help text with examples

program.addHelpText('after', `
Examples:
    $ cosoft                          # Run in interactive mode starting with calendar view
    $ cosoft list                     # List all meeting rooms and exit
    $ cosoft book                     # Book a room interactively and exit
    $ cosoft book --room-name "CALL BOX 3" --date 2025-09-03 --start-time 13:30 --end-time 14:30  # Book directly with arguments
    $ cosoft book --batch-file bookings.json     # Book multiple rooms from JSON file
    $ cosoft book --batch-json '[{"roomName":"CALL BOX 3","date":"2025-09-03","startTime":"13:30","endTime":"14:30"}]'  # Book multiple rooms from JSON string
    $ cosoft calendar 2025-09-15      # Show calendar for specific date and exit
    $ cosoft my-bookings              # Show all your bookings and exit
    $ cosoft my-bookings 2025-09-15   # Show your bookings for a specific date and exit
    $ cosoft cancel                   # Cancel a booking interactively
    $ cosoft cancel --booking-id "40adf2d0-3a3b-45cd-bb47-b34c011201c1"  # Cancel a specific booking by ID
    $ cosoft cancel --batch-ids "id1,id2,id3"    # Cancel multiple bookings by IDs
    $ cosoft cancel --batch-file cancellations.json  # Cancel multiple bookings from JSON file
    $ cosoft auth                     # Check authentication status and display logged-in user
    $ cosoft login                    # Login to CoSoft interactively
    $ cosoft mcp-server               # Start as MCP (Model Context Protocol) server for stdio communication
`);

// Parse arguments
program.parse();