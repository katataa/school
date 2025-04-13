import ContinuousKeyPressListener from './ContinuousKeyPressListener.js';

function generateUniqueId() {
  return '_' + Math.random().toString(36).substr(2, 9);
}
const playerId = generateUniqueId();


const coinSound = new Audio('/sounds/coin.mp3');
const gameSound = new Audio('/sounds/game.mp3');
const gameOverSound = new Audio('/sounds/game-over.mp3');
const chaseSound = new Audio('/sounds/chase.mp3');

let isHardMode = false;
let lastMoveTime = 0;
let myPosition = { x: 1, y: 4, direction: "right" };
let gameStarted = false;
let gamePausedClient = false;

const mapData = {
  minX: 1,
  maxX: 14,
  minY: 4,
  maxY: 12,
  blockedSpaces: {
    "7x7": true, "7x8": true,
    "6x7": true, "8x7": true,
    "6x8": true, "8x8": true,
    "4x10": true, "5x10": true,
    "10x10": true, "11x10": true,
    "7x4": true,
    "1x11": true,
    "12x10": true,
  },
};

const playerColors = ["brown", "orange", "gray", "black", "white", "tuxedo"];

function randomFromArray(array) {
  return array[Math.floor(Math.random() * array.length)];
}

function getKeyString(x, y) {
  return `${x}x${y}`;
}

function resetGameUI() {
  gameStarted = false;
  gamePausedClient = false;
  const me = players[playerId];
  if (me && me.isLead) {
    startGameBtn.style.display = "inline-block";
  }
  playerNameInput.style.display = "inline-block";
}


function createName() {
  const prefix = randomFromArray([
    "COOL", "SUPER", "HIP", "SMUG", "COOL", "SILKY", "GOOD",
    "SAFE", "DEAR", "DAMP", "WARM", "RICH", "LONG", "DARK",
    "SOFT", "BUFF", "DOPE",
  ]);
  const animal = randomFromArray([
    "BEAR", "DOG", "CAT", "FOX", "LAMB", "LION", "BOAR",
    "GOAT", "VOLE", "SEAL", "PUMA", "MULE", "BULL", "BIRD", "BUG",
  ]);
  return `${prefix} ${animal}`;
}

function isSolid(x, y) {
  const blockedNextSpace = mapData.blockedSpaces[getKeyString(x, y)];
  return (
    blockedNextSpace ||
    x >= mapData.maxX || x < mapData.minX ||
    y >= mapData.maxY || y < mapData.minY
  );
}

function getRandomSafeSpot() {
  const safeSpots = [
    { x: 1, y: 4 }, { x: 2, y: 4 }, { x: 1, y: 5 }, { x: 2, y: 6 },
    { x: 2, y: 8 }, { x: 2, y: 9 }, { x: 4, y: 8 }, { x: 5, y: 5 },
    { x: 5, y: 8 }, { x: 11, y: 7 }, { x: 12, y: 7 }, { x: 13, y: 7 },
    { x: 13, y: 6 }, { x: 13, y: 8 }, { x: 7, y: 6 }, { x: 10, y: 8 },
    { x: 11, y: 4 }, { x: 3, y: 7 }, { x: 9, y: 7 }, { x: 6, y: 5 },
    { x: 8, y: 5 }, { x: 4, y: 6 }, { x: 10, y: 6 }, { x: 12, y: 5 }
  ];
  const validSpots = safeSpots.filter(({ x, y }) => {
    return (
      x >= mapData.minX && x <= mapData.maxX &&
      y >= mapData.minY && y <= mapData.maxY &&
      !mapData.blockedSpaces[getKeyString(x, y)]
    );
  });

  return randomFromArray(validSpots);
}

const gameContainer = document.querySelector(".game-container");
const playerNameInput = document.querySelector("#player-name");
const playerColorBtn = document.querySelector("#player-color");
const startGameBtn = document.querySelector("#start-game");
const timerDisplay = document.querySelector("#timer");
const scoreboardList = document.querySelector("#scoreboard-list");

const inGameMenu = document.querySelector("#in-game-menu");
const pauseBtn = document.querySelector("#pause-btn");
const resumeBtn = document.querySelector("#resume-btn");
const quitBtn = document.querySelector("#quit-btn");
const restartBtn = document.querySelector("#restart-btn");
const menuStatusMsg = document.querySelector("#menu-status-msg");
const singlePlayerCheckbox = document.querySelector("#single-player-checkbox");
const singlePlayerControls = document.querySelector("#single-player-controls");
const singlePlayerLabel = document.querySelector("#single-player-label");
const socket = io();

let players = {};
let playerElements = {};
let coins = {};
let coinElements = {};

function showMenuMessage(text) {
  if (!menuStatusMsg) return;
  menuStatusMsg.textContent = text;
  menuStatusMsg.classList.remove("hidden");
  setTimeout(() => {
    menuStatusMsg.textContent = "";
    menuStatusMsg.classList.add("hidden");
  }, 5000);
}

socket.on('connect', () => {
  console.log('Connected, socket.id:', socket.id);
});

socket.on('name-taken', (takenName) => {
  alert(`Name "${takenName}" is taken. Please choose another one.`);
  playerNameInput.value = '';
});

socket.on('start-failed', (msg) => {
  alert(msg);
});

socket.on('players-update', (serverPlayers) => {
  players = serverPlayers;
  updatePlayers();

  if (gameStarted) {
    startGameBtn.style.display = "none";
  } else {
    const me = players[playerId];
    if (me && me.isLead) {
      startGameBtn.style.display = "inline-block";
    }
  }

  if (singlePlayerCheckbox.checked) {
    updateNPCCustomizationPanel();
  }
});


socket.on('coins-update', (serverCoins) => {
  coins = serverCoins;
  updateCoins();
});

socket.on('game-started', (data) => {
  gameSound.play().catch(console.error);
  gameStarted = true;
  gamePausedClient = false;
  isHardMode = data.hardMode;

  if (data.terrorActive) {
    gameSound.pause();
    gameSound.currentTime = 0;
    chaseSound.play().catch(console.error);
  } else {
    gameSound.play().catch(console.error);
  }
  placeCoin();
  playerNameInput.style.display = "none";
  startGameBtn.style.display = "none";
  singlePlayerControls.style.display = "none";
  singlePlayerCheckbox.checked = false;
  singlePlayerLabel.style.display = "none";
});

socket.on('time-update', (timeLeft) => {
  if (timerDisplay) {
    timerDisplay.textContent = timeLeft + "s";
  }
});

socket.on('game-ended', ({ winnerName, scoreboard }) => {
  gameOverSound.play().catch(console.error);
  gameStarted = false;
  gamePausedClient = false;

  let scoreboardText = "Final Scores:\n";
  scoreboard.forEach((p) => {
    scoreboardText += `${p.name}: ${p.coins}\n`;
  });
  scoreboardText += `\nWinner: ${winnerName}`;
  alert(scoreboardText);

  if (timerDisplay) {
    timerDisplay.textContent = "GAME OVER";
  }

  const me = players[playerId];
  if (me && me.isLead) {
    startGameBtn.style.display = "inline-block";
  }
  startGameBtn.style.display = "inline-block";
  playerNameInput.style.display = "inline-block";
  singlePlayerLabel.style.display = "inline-block";

  resetGameUI();
});

socket.on('game-paused', ({ by }) => {
  gamePausedClient = true;
  showMenuMessage(`Paused by ${by}`);
});

socket.on('game-resumed', ({ by }) => {
  gamePausedClient = false;
  showMenuMessage(`Resumed by ${by}`);
});

socket.on('player-quit', ({ by }) => {
  showMenuMessage(`${by} has quit the game`);
});

socket.on('game-restarted', ({ by }) => {
  showMenuMessage(`${by} restarted the game`);
});

const characterSprites = {
  brown: '/images/cats.png',
  orange: '/images/cats.png',
  gray: '/images/cats.png',
  black: '/images/cats.png',
  white: '/images/cats.png',
  tuxedo: '/images/cats.png',
};

const characterOffsets = {
  brown: '0px 0px',
  orange: '-24px 0px',
  gray: '-48px 0px',
  black: '0px -24px',
  white: '-24px -24px',
  tuxedo: '-48px -24px',
};

function updatePlayers() {
  Object.keys(players).forEach((key) => {
    const st = players[key];
    let el = playerElements[key];
    if (!el) {
      el = document.createElement("div");
      el.classList.add("Character", "grid-cell");
      if (key === playerId) {
        el.classList.add("you");
      }
      el.innerHTML = `
        <div class="Character_shadow grid-cell"></div>
        <div class="Character_sprite grid-cell"></div>
        <div class="Character_name-container">
          <span class="Character_name"></span>
          <span class="Character_coins">0</span>
        </div>
        <div class="Character_you-arrow"></div>
      `;
      playerElements[key] = el;
      gameContainer.appendChild(el);
    }

    const spriteEl = el.querySelector(".Character_sprite");
    spriteEl.style.backgroundImage = `url(${characterSprites[st.color]})`;
    spriteEl.style.backgroundPosition = characterOffsets[st.color] || '0px 0px';

    el.querySelector(".Character_name").innerText = st.name;
    el.querySelector(".Character_coins").innerText = st.coins;
    el.setAttribute("data-color", st.color);
    el.setAttribute("data-direction", st.direction);

    const left = 16 * st.x + "px";
    const top = 16 * st.y - 4 + "px";
    el.style.transform = `translate3d(${left}, ${top}, 0)`;
  });

  Object.keys(playerElements).forEach((key) => {
    if (!players[key]) {
      gameContainer.removeChild(playerElements[key]);
      delete playerElements[key];
    }
  });

  updateScoreboard();
}

function updateScoreboard() {
  if (!scoreboardList) return;
  scoreboardList.innerHTML = "";

  const sorted = Object.values(players).sort((a, b) => b.coins - a.coins);
  sorted.forEach((p) => {
    const li = document.createElement("li");
    li.textContent = `${p.name}: ${p.coins}`;
    scoreboardList.appendChild(li);
  });
}

function updateCoins() {
  Object.keys(coins).forEach((key) => {
    if (!coinElements[key]) {
      const c = coins[key];
      const coinEl = document.createElement("div");
      coinEl.classList.add("Coin", "grid-cell");
      coinEl.innerHTML = `
        <div class="Coin_shadow grid-cell"></div>
        <div class="Coin_sprite grid-cell"></div>
      `;
      const left = 16 * c.x + "px";
      const top = 16 * c.y - 4 + "px";
      coinEl.style.transform = `translate3d(${left}, ${top}, 0)`;
      coinElements[key] = coinEl;
      gameContainer.appendChild(coinEl);
    }
  });
  Object.keys(coinElements).forEach((key) => {
    if (!coins[key]) {
      gameContainer.removeChild(coinElements[key]);
      delete coinElements[key];
    }
  });
}

function placeCoin() {
  if (!gameStarted) return;
  if (gamePausedClient) {
    setTimeout(placeCoin, 2000);
    return;
  }

  const { x, y } = getRandomSafeSpot();
  socket.emit('coin-add', { x, y });

  let timeouts;
  if (isHardMode) {
    timeouts = [1000, 2000, 3000, 4000];
  } else {
    timeouts = [300, 500, 800, 1000];
  }

  setTimeout(placeCoin, randomFromArray(timeouts));
}

function attemptGrabCoin(x, y) {
  const key = getKeyString(x, y);
  if (coins[key]) {
    coinSound.play().catch(console.error);
    socket.emit('coin-remove', { x, y, playerId });
  }
}

let movementCooldown = false;

function handleArrowPress(xChange = 0, yChange = 0) {
  if (!gameStarted || gamePausedClient || movementCooldown) return;

  const nextX = myPosition.x + xChange;
  const nextY = myPosition.y + yChange;

  if (isSolid(nextX, nextY)) {
    return;
  }

  myPosition.x = nextX;
  myPosition.y = nextY;
  if (xChange === 1) myPosition.direction = "right";
  if (xChange === -1) myPosition.direction = "left";

  movementCooldown = true;
  setTimeout(() => {
    movementCooldown = false;
  }, 120);

  socket.emit('player-move', {
    id: playerId,
    name: players[playerId]?.name,
    color: players[playerId]?.color || "brown",
    coins: players[playerId]?.coins,
    ...myPosition,
  });

  attemptGrabCoin(nextX, nextY);
}

function initMenuButtons() {
  pauseBtn.addEventListener("click", () => {
    const me = players[playerId];
    if (!me) return;
    socket.emit('pause-game', { name: me.name });
  });

  resumeBtn.addEventListener("click", () => {
    const me = players[playerId];
    if (!me) return;
    socket.emit('resume-game', { name: me.name });
  });

  quitBtn.addEventListener("click", () => {
    const me = players[playerId];
    if (!me) return;
    socket.emit('quit-game', { playerId, name: me.name });
  });

  restartBtn.addEventListener("click", () => {
    const me = players[playerId];
    if (!me) return;
    socket.emit('restart-game', { name: me.name });
  });

  socket.on('game-paused', ({ by }) => {
    pauseBtn.classList.add("hidden");
    resumeBtn.classList.remove("hidden");
    showMenuMessage(`Paused by ${by}`);
  });

  socket.on('game-resumed', ({ by }) => {
    pauseBtn.classList.remove("hidden");
    resumeBtn.classList.add("hidden");
    showMenuMessage(`Resumed by ${by}`);
  });

  socket.on('player-quit', ({ by }) => {
    showMenuMessage(`${by} quit the game`);
  });

  socket.on('game-restarted', ({ by }) => {
    showMenuMessage(`Game restarted by ${by}`);
  });
}

function setupKeyboard() {
  new ContinuousKeyPressListener("ArrowUp", () => handleArrowPress(0, -1), 150);
  new ContinuousKeyPressListener("ArrowDown", () => handleArrowPress(0, 1), 150);
  new ContinuousKeyPressListener("ArrowLeft", () => handleArrowPress(-1, 0), 150);
  new ContinuousKeyPressListener("ArrowRight", () => handleArrowPress(1, 0), 150);
}

function initGame() {
  let defaultName = createName();
  socket.emit('player-join', {
    id: playerId,
    name: defaultName,
    direction: "right",
    color: "brown",
    x: 1,
    y: 4,
    coins: 0,
  });

  playerNameInput.addEventListener("change", (e) => {
    const newName = e.target.value.trim();
    if (newName.length > 0) {
      socket.emit('player-update', { id: playerId, name: newName });
    }
  });

  playerColorBtn.addEventListener("click", () => {
    const currColor = players[playerId]?.color || "brown";
    const idx = playerColors.indexOf(currColor);
    const next = playerColors[idx + 1] || playerColors[0];
    socket.emit('player-update', { id: playerId, color: next });
  });

  singlePlayerCheckbox.addEventListener("change", () => {
    if (singlePlayerCheckbox.checked) {
      singlePlayerControls.style.display = "block";
    } else {
      singlePlayerControls.style.display = "none";
    }

    socket.emit("toggle-single-player", { enabled: singlePlayerCheckbox.checked });
  });

  document.addEventListener("click", (event) => {
    if (singlePlayerControls.style.display !== "none") {
      const clickedInsidePanel = singlePlayerControls.contains(event.target);
      const clickedOnCheckboxLabel = singlePlayerLabel.contains(event.target);

      if (!clickedInsidePanel && !clickedOnCheckboxLabel) {
        singlePlayerControls.style.display = "none";
      }
    }
  });

  startGameBtn.addEventListener("click", () => {

    const hardModeCheckbox = document.querySelector("#hard-mode-checkbox");
    const hardMode = hardModeCheckbox && hardModeCheckbox.checked;

    const singlePlayerMode = singlePlayerCheckbox.checked;
    let npcCount = 0;
    let npcDifficulty = 'medium';

    if (singlePlayerMode) {
      const countRadios = document.querySelectorAll('input[name="sp-count"]');
      for (let radio of countRadios) {
        if (radio.checked) {
          npcCount = parseInt(radio.value, 10);
          break;
        }
      }
      npcDifficulty = document.querySelector("#npc-difficulty").value;
    }
    let npcConfigs = [];

    for (let i = 0; i < npcCount; i++) {
      const selectElem = document.getElementById(`npc-color-${i}`);
      npcConfigs.push({ color: selectElem.value });
    }


    socket.emit('start-game', {
      hardMode,
      singlePlayerMode,
      npcCount,
      npcDifficulty,
      npcConfigs
    });
  });

  initMenuButtons();
  setupKeyboard();
}

function initPreGameNPCPanel() {
  const panel = document.createElement("div");
  panel.id = "npc-customization";
  panel.style.marginTop = "10px";
  panel.style.border = "1px solid #ccc";
  panel.style.padding = "5px";
  panel.innerHTML = "<h3>NPC Customization</h3>";

  const maxNPCs = 3;
  for (let i = 0; i < maxNPCs; i++) {
    const npcDiv = document.createElement("div");
    npcDiv.id = `npc-config-${i}`;
    npcDiv.innerHTML = `NPC ${i + 1} Color: <select id="npc-color-${i}">
      ${playerColors.map(color => `<option value="${color}">${color}</option>`).join("")}
    </select>`;
    panel.appendChild(npcDiv);
  }

  singlePlayerControls.appendChild(panel);
}

document.addEventListener("DOMContentLoaded", () => {
  initGame();
  if (singlePlayerControls) {
    initPreGameNPCPanel();
  }
});

function updateNPCCustomizationPanel() {
  const panel = document.getElementById("npc-customization");
  if (!panel) return;
  Array.from(panel.querySelectorAll("button.npc-color-btn")).forEach((btn) => btn.remove());

  Object.values(players).forEach((npc) => {
    if (npc.ai) {
      const btn = document.createElement("button");
      btn.classList.add("npc-color-btn");
      btn.style.marginRight = "5px";
      btn.textContent = `NPC ${npc.name} Color: ${npc.color}`;
      btn.addEventListener("click", () => {
        const idx = playerColors.indexOf(npc.color);
        const nextColor = playerColors[(idx + 1) % playerColors.length];
        btn.textContent = `NPC ${npc.name} Color: ${nextColor}`;
        socket.emit('player-update', { id: npc.id, color: nextColor });
      });
      panel.appendChild(btn);
    }
  });
}
