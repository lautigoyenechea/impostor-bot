# 🔪 Dimpostor - Discord Impostor Game Bot

A Discord bot to play the Impostor game with your friends in voice channels!

## 🎮 How to Play

1. **Start the Game**: Join a voice channel with your friends and use the `/start` command in any text channel
2. **Receive Your Word**: The bot will DM each player privately:
   - **Regular players** receive a secret word
   - **The impostor** receives only a hint about the word
3. **Discuss**: Use your voice channel to discuss and describe the word without saying it directly. The impostor must blend in without knowing the actual word!
4. **Vote**: When ready, the admin uses the `/vote` command. Each player receives a DM with voting buttons
5. **Results**: After voting, the bot reveals who was voted out and whether they were the impostor

### Win Conditions
- **Innocents win** if they successfully vote out the impostor
- **Impostor wins** if they avoid being voted out

## 🚀 Setup & Installation

### Prerequisites
- [Go 1.24.6+](https://golang.org/dl/)
- A Discord account and server
- Discord Developer Portal access

### 1. Create a Discord Bot

1. Go to the [Discord Developer Portal](https://discord.com/developers/applications)
2. Click "New Application" and give it a name
3. Navigate to the "Bot" section and click "Add Bot"
4. Under the bot settings, enable these **Privileged Gateway Intents**:
   - Server Members Intent
   - Message Content Intent
5. Copy your bot token (you'll need this later)
6. Go to "OAuth2" → "URL Generator":
   - Select scopes: `bot`, `applications.commands`
   - Select bot permissions: `Send Messages`, `Read Messages/View Channels`, `Use Slash Commands`
   - Copy the generated URL and use it to invite the bot to your server

### 2. Clone and Configure

```bash
# Clone the repository
git clone https://github.com/lautigoyenechea/impostor-bot.git
cd impostor-bot

# Create a .env file
touch .env
```

Add the following to your `.env` file:
```env
AUTH_TOKEN=your_discord_bot_token_here
GUILD_ID=your_discord_server_id_here
```

To get your Guild ID, enable Developer Mode in Discord (Settings → Advanced → Developer Mode), then right-click your server and select "Copy Server ID".

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run the Bot

```bash
go run .
```

You should see:
```
Bot is ready
Impostor Bot is online!
```

## 📝 Commands

- `/start` - Start a new game (must be in a voice channel)
- `/vote` - Begin the voting phase (admin only)

## 🤝 Contributing

Feel free to open issues or submit pull requests!

## TODO:
- Support multiple impostors by an option on the /start command.
- Support draws in voting sessions.


## 📄 License

This project is open source and available for personal use.
