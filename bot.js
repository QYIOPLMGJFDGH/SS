const TelegramBot = require('node-telegram-bot-api');
const readlineSync = require('readline-sync');
const MTProto = require('telegram-mtproto');
const express = require('express');

// Hardcoded values for API ID, API Hash, and Bot Token
const botToken = '7426075639:AAE854r2874ZJAVat6zVUeSR4IYBnsW-y-w'; // Replace with actual Bot Token
const apiId = '16457832';                   // Your Telegram API ID
const apiHash = '3030874d0befdb5d05597deacc3e83ab'; // Your Telegram API Hash

// Initialize Telegram Bot
const bot = new TelegramBot(botToken, { polling: true });

// Initialize Express app
const app = express();
const port = 8080;

// MTProto Client Configuration
const settings = {
  api_id: apiId,
  api_hash: apiHash,
  phone_number: null,
};

const client = MTProto(settings);

// Validate phone number format
function validatePhoneNumber(phone) {
  const phoneRegex = /^\+?[1-9]\d{1,14}$/;
  return phoneRegex.test(phone);
}

// Handle incoming messages
bot.on('message', async (msg) => {
  const chatId = msg.chat.id;

  if (msg.text === '/start') {
    const options = {
      reply_markup: {
        inline_keyboard: [
          [
            { text: 'Generate v1 Session', callback_data: 'generate_v1' },
            { text: 'Generate v2 Session', callback_data: 'generate_v2' },
          ],
        ],
      },
    };
    bot.sendMessage(chatId, 'Welcome to the String Session Generator Bot! Select a session type below to begin.', options);
  }
});

// Handle callback queries (button clicks)
bot.on('callback_query', async (query) => {
  const chatId = query.message.chat.id;
  const data = query.data;

  // Acknowledge the callback query
  bot.answerCallbackQuery(query.id);

  if (data === 'generate_v1' || data === 'generate_v2') {
    let sessionType = data === 'generate_v2' ? 'v2' : 'v1';
    bot.sendMessage(chatId, `You selected ${sessionType} session.\n\nPlease send your API ID:`);

    // Collect user inputs for API ID, API Hash, and Phone Number
    await waitForUserInputs(chatId, sessionType);
  }
});

// Wait for user inputs (API ID, API Hash, Phone Number)
async function waitForUserInputs(chatId, sessionType) {
  const apiIdInput = readlineSync.question('Enter your Telegram API ID: ');
  const apiHashInput = readlineSync.question('Enter your Telegram API Hash: ');
  const phoneInput = readlineSync.question('Enter your phone number (with country code): ');

  if (!validatePhoneNumber(phoneInput)) {
    console.log("Invalid phone number format. Please enter a valid phone number.");
    return;
  }

  const stringSession = await generateStringSession(apiIdInput, apiHashInput, phoneInput, sessionType);
  bot.sendMessage(chatId, `Your ${sessionType} string session:\n\n\`${stringSession}\`\n\nSave it securely!`);
}

// Generate String Session using MTProto
async function generateStringSession(apiId, apiHash, phone, sessionType) {
  try {
    settings.phone_number = phone;

    // Send verification code to phone
    const phone_code_hash = await client('auth.sendCode', {
      phone_number: phone,
      settings: { api_id: apiId, api_hash: apiHash },
    });

    // Get OTP from user
    const otp = readlineSync.question('Enter the OTP sent to your phone: ');

    // Sign in the user
    await client('auth.signIn', {
      phone_number: phone,
      phone_code_hash: phone_code_hash.phone_code_hash,
      phone_code: otp,
    });

    console.log("Successfully authenticated!");

    // Get string session
    const stringSession = await client.session.getStringSession();
    return stringSession;
  } catch (error) {
    console.error('Error generating string session:', error);
    return 'Failed to generate string session. Please try again.';
  }
}

// Start Express HTTP server
app.get('/', (req, res) => {
  res.send('Bot is running!');
});

app.listen(port, () => {
  console.log(`Server running on port ${port}`);
});
