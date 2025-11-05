// Login command for the CoSoft CLI
import chalk from 'chalk';
import ora from 'ora';
import input from '@inquirer/input';
import password from '@inquirer/password';
import { login } from '../config/api.js';
import { checkAuthentication } from '../config/api.js';

export async function loginCommand() {
    console.log(chalk.bold.blue('\nCoSoft Login'));
    console.log(chalk.gray('Please enter your credentials to authenticate with CoSoft.'));

    try {
        // Ask for email
        const email = await input({
            message: 'Email:',
            validate: (value) => {
                if (!value) return 'Email is required';
                // Basic email validation
                const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
                if (!emailRegex.test(value)) return 'Please enter a valid email address';
                return true;
            }
        });

        // Ask for password
        const pwd = await password({
            message: 'Password:',
            mask: true,
            validate: (value) => {
                if (!value) return 'Password is required';
                return true;
            }
        });

        // Show loading spinner
        const spinner = ora('Logging in...').start();

        // Call the login function
        const result = await login(email, pwd);

        if (result.success) {
            spinner.succeed('Login successful');

            // Display user information similar to the auth command
            if (result.user) {
                const fullName = result.user.FirstName && result.user.LastName
                    ? `${result.user.FirstName} ${result.user.LastName}`
                    : result.user.Email || 'Unknown user';

                console.log(chalk.green(`✓ Authenticated as: ${chalk.bold(fullName)}`));
                console.log(chalk.white(`Email: ${result.user.Email}`));

                console.log(chalk.green('\nAuthentication tokens have been saved to .auth file.'));
                console.log(chalk.blue('You can now use other commands that require authentication.'));
            }
        } else {
            spinner.fail('Login failed');
            console.error(chalk.red(`Error: ${result.message}`));
            process.exit(1);
        }
    } catch (error) {
        console.error(chalk.red('Error during login process:'), error.message);
        process.exit(1);
    }
}