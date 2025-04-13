# 😻 Multiplayer Web Game: "Cat Chase Royale" 😽

## 👀 Overview
"Cat Chase Royale" is a fast-paced, real-time multiplayer game built for 2 to 4 players, playable entirely through web browsers without HTML canvas. Players control adorable cats competing to collect coins while avoiding obstacles and racing against the clock.

## Features
- **Real-Time Multiplayer:** Play with 2 to 4 players simultaneously.
- **Single-Player Mode with NPCs:** Host can select to play alone with 1–3 AI opponents of varying difficulties (Easy, Medium, Hard).
- **Equal Gameplay:** All characters have equal abilities.
- **Unique Player Names:** Each player selects a unique name.
- **Live Scoreboard:** View all players' scores in real-time.
- **Game Timer:** Tracks time until the round ends.
- **Smooth Animations:** Maintains 60 FPS using requestAnimationFrame.
- **In-Game Menu:** Pause, resume, restart, or quit the game with alerts.
- **Sound Effects:** Enjoy immersive game sounds for actions like coin collection and game start.

## 🧠 NPC (AI) Opponents
- **Difficulty Levels:**
  - **Easy** – Moves slowly, randomly fails or does nothing often.
  - **Medium** – Decent movement toward coins, occasionally deviates.
  - **Hard** – Aggressively seeks out coins with near-perfect movement.
- **Customization:** The host can configure each NPC’s color and difficulty before starting.
- **Diverse Behavior:** Random “failure” and path variation prevent NPCs from looping endlessly.

## 🎮 How to Play
1. **Join the Game:** Visit the URL provided by the host.
2. **Choose Your Name:** Enter a unique player name.
3. **Toggle Single-Player or Multi-Player:** If single-player is chosen, NPCs will spawn, and real-player connections are refused.
4. **Start the Match:** The host (lead player) starts the game.
5. **Collect Coins:** Use keyboard controls to move your cat and collect coins.
6. **Compete:** Track scores and time live.
7. **Win:** Highest score when the timer hits zero wins!

## Controls
- **Arrow Keys:** Move Up, Down, Left, Right
- **P:** Pause/Resume the Game
- **R:** Restart (Host Only)
- **Q:** Quit Game

## Technology Stack
- **Frontend:** HTML, CSS, JavaScript (No Canvas)
- **Backend:** Node.js with Socket.IO for real-time communication
- **Networking:** WebSockets for real-time player updates

## Requirements Met
- ✅ Runs at 60 FPS with requestAnimationFrame  
- ✅ Supports 2 to 4 players with real-time multiplayer gameplay  
- ✅ Players join via URL and choose unique names  
- ✅ Single-Player mode with 1–3 virtual opponents (AI)  
- ✅ Equal character abilities with live score tracking  
- ✅ In-game pause, resume, restart, and quit functionality  
- ✅ Real-time game timer  
- ✅ Keyboard controls with no input delays  
- ✅ Sound effects for gameplay events  

## Bonus Features
- 🎉 Power-Ups: Temporary speed boosts for players  
- 🎉 Spectator Mode: View the game if the lobby is full  
- 🎉 Customizable Controls: Players can remap keys  

## ⚙️ Installation & Setup
1. Clone the repository:
   ```bash
   git clone https://gitea.kood.tech/larajohannaluks/web-game.git
   ```
2. Install dependencies:
   ```bash
   cd web-game
   npm install
   ```
3. Run the server:
   ```bash
   node server.js
   ```
4. Share the URL shown in the terminal with players.

## Using Ngrok for External URL Sharing
To allow players to join the game from outside your local network, use ngrok to expose your server to the internet.

### Steps to Use Ngrok:
1. **Install Ngrok:**
   - Download and install from [ngrok.com](https://ngrok.com/)
2. **Expose the Server:**
   - After running the server, open another terminal and run:
   ```bash
   ngrok http 3000
   ```
3. **Share the URL:**
   - Ngrok will provide a public URL. Share this URL with players.
4. **Update the Client:**
   - Ensure your game client connects to the provided ngrok URL.

## ❗️ Hard mode
- Coins spawn less frequently and despawn after 3 seconds.
- AI/human gameplay is otherwise unchanged, but you have less time to grab coins.

## Any questions?
### DM us on Discord!
- @laracroft8
- @katata420