import chalk from 'chalk';
import ora from 'ora';
import Table from 'cli-table3';
import select from '@inquirer/select';
import { API_CONFIG } from '../config/env.js';
import { makeApiRequest, checkAuthentication, getRefererHeaders } from '../config/api.js';
// Import getRoomColor function from calendarView.js
import { getRoomColor } from './calendarView.js';
import { showStandardMenu } from '../utils/sharedMenus.js';

// Helper function to fetch rooms from the API
export async function fetchRooms() {
    try {
        const url = `${API_CONFIG.BASE_URL}/v2/api/api/CoworkingSpace/${API_CONFIG.COWORKING_SPACE_ID}/category/${API_CONFIG.CATEGORY_ID}/items?price=null`;

        const refererPath = `/v2/new-reservation/${API_CONFIG.COWORKING_SPACE_ID}/${API_CONFIG.CATEGORY_ID}`;

        const responseData = await makeApiRequest({
            method: 'post',
            url: url,
            extraHeaders: getRefererHeaders(refererPath),
            data: {
                price: null
            }
        });

        // Extract rooms from both VisitedItems and UnvisitedItems
        let roomsData = [];

        // Add VisitedItems if they exist
        if (responseData && responseData.VisitedItems && Array.isArray(responseData.VisitedItems)) {
            roomsData = [...responseData.VisitedItems];
        }

        // Add UnvisitedItems if they exist
        if (responseData && responseData.UnvisitedItems && Array.isArray(responseData.UnvisitedItems)) {
            roomsData = [...roomsData, ...responseData.UnvisitedItems];
        }

        if (roomsData.length === 0) {
            console.warn('Warning: No rooms found in either VisitedItems or UnvisitedItems.');
            return []; // Return empty array if no rooms found
        }

        // Transform API response to our format
        const rooms = roomsData.map((room, index) => {
            // Extract hourly price from the Prices array
            let hourlyPrice = 'N/A';
            if (room.Prices && Array.isArray(room.Prices) && room.Prices.length > 0) {
                // Find the hourly price option
                const hourlyPriceOption = room.Prices.find(price =>
                    price.DurationType === 'hour' || price.Description === 'heure');

                if (hourlyPriceOption) {
                    // Format the price with 2 decimal places and € symbol - now using EuroHT instead of EuroTTC
                    hourlyPrice = `${hourlyPriceOption.EuroHT.toFixed(2)} €`;
                }
            }

            // First try to get equipment from the Equipments field
            let equipments = [];
            if (room.Equipments && Array.isArray(room.Equipments) && room.Equipments.length > 0) {
                equipments = room.Equipments;
            }
            // If no equipment found, try to extract from ShortDescription
            else if (room.ShortDescription) {
                // Look for "Equipement : " or similar patterns in the description
                const equipmentMatch = room.ShortDescription.match(/quipement\s*:\s*([^<]+)/i);
                if (equipmentMatch && equipmentMatch[1]) {
                    // Split by commas if multiple equipments are listed
                    equipments = equipmentMatch[1].split(',').map(item => item.trim());
                }
            }

            return {
                id: index + 1, // Assign sequential ID for UI purposes
                apiId: room.Id, // API uses uppercase 'Id'
                name: room.Name || `Room ${index + 1}`, // API uses uppercase 'Name'
                capacity: room.NbUsers || 'N/A', // Use NbUsers field for capacity
                hourlyPrice: hourlyPrice,
                floor: room.Floor || 'N/A',
                available: !room.IsLocked,
                description: room.ShortDescription ? room.ShortDescription.replace(/<\/?p>/g, '') : '',
                equipments: equipments
            };
        });

        // Sort rooms alphabetically by name
        return rooms.sort((a, b) => a.name.localeCompare(b.name));
    } catch (error) {
        console.error('Error fetching rooms:', error);
        throw error;
    }
}

export async function listRooms(options = {}) {

    // Check if running in interactive mode
    const isInteractive = options.interactive === true;

    const spinner = ora('Loading meeting rooms...').start();

    try {
        const rooms = await fetchRooms();
        spinner.succeed('Meeting rooms loaded');

        // Create a table to display room information
        const table = new Table({
            head: [
                chalk.white('Name'),
                chalk.white('Capacity'),
                chalk.white('Floor'),
                chalk.white('Price'),
                chalk.white('Status')
            ],
            colWidths: [25, 10, 10, 15, 15]
        });

        // Add room data to the table
        rooms.forEach(room => {
            const roomColor = getRoomColor(room.id);
            table.push([
                chalk[roomColor](room.name),
                chalk[roomColor](room.capacity.toString()),
                chalk[roomColor](room.floor.toString()),
                chalk[roomColor](room.hourlyPrice),
                room.available ? chalk.green('Available') : chalk.red('Booked')
            ]);
        });

        console.log(chalk.bold.blue('\nAvailable Meeting Rooms:'));
        console.log(table.toString());

        // If not in interactive mode, exit after displaying rooms
        if (!isInteractive) {
            return;
        }

        // In interactive mode, show a menu
        await showStandardMenu();

    } catch (error) {
        spinner.fail('Failed to load meeting rooms');
        console.error(chalk.red('Error:', error.message));
        if (!isInteractive) {
            process.exit(1);
        }
    }
}