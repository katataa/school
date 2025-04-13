const express = require('express');
const http = require('http');
const { Server } = require('socket.io');

const app = express();
const server = http.createServer(app);
const io = new Server(server, {
  cors: {
    origin: "*",
  },
});

app.get('/', (req, res) => {
  res.sendFile(__dirname + '/public/index.html');
});
app.use('/styles', express.static(__dirname + '/styles'));
app.use('/images', express.static(__dirname + '/public/images'));
app.use(express.static(__dirname + '/public'));

let hardModeEnabled = false;
let players = {};
let coins = {};
let leadPlayerSocketId = null;
let gameStarted = false;
let gamePaused = false;
let timeLeft = 0;
let timerInterval = null;
let singlePlayerMode = false;
let npcInterval = null;

const mapInfo = {
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

function isBlocked(x, y) {
  const key = `${x}x${y}`;
  if (mapInfo.blockedSpaces[key]) return true;
  if (x < mapInfo.minX || x > mapInfo.maxX) return true;
  if (y < mapInfo.minY || y > mapInfo.maxY) return true;
  return false;
}

function randomDir() {
  const r = Math.random();
  if (r < 0.25) return [1, 0];
  if (r < 0.5) return [-1, 0];
  if (r < 0.75) return [0, 1];
  return [0, -1];
}

function tryMove(npc, [dx, dy]) {
  if (dx === 0 && dy === 0) return;
  const nx = npc.x + dx;
  const ny = npc.y + dy;
  if (!isBlocked(nx, ny)) {
    npc.x = nx;
    npc.y = ny;
    if (dx > 0) npc.direction = "right";
    if (dx < 0) npc.direction = "left";
  }
}

function moveNPC(npc) {
  if (npc.difficulty === 'terror') {
    if (npc.terrorCooldown > 0) {
      npc.terrorCooldown--;
      return;
    }

    npc.terrorMoveFrame = (npc.terrorMoveFrame || 0) + 1;
    if (npc.terrorMoveFrame % 3 === 0) {
      return;
    }

    const target = getClosestRealPlayer(npc);
    if (!target) {
      tryMove(npc, randomDir());
      return;
    }

    let dx = 0, dy = 0;
    if (target.x > npc.x) dx = 1;
    else if (target.x < npc.x) dx = -1;
    if (target.y > npc.y) dy = 1;
    else if (target.y < npc.y) dy = -1;

    tryMove(npc, [dx, dy]);

    if (npc.x === target.x && npc.y === target.y) {
      players[target.id].coins = 0;
      io.emit('players-update', players);
      npc.terrorCooldown = 20;
    }

    return;
  }

  if (npc.difficulty === 'hard') {
    let closestCoin = null;
    let minDist = Infinity;
    for (const cKey in coins) {
      const c = coins[cKey];
      const dist = Math.abs(c.x - npc.x) + Math.abs(c.y - npc.y);
      if (dist < minDist) {
        minDist = dist;
        closestCoin = c;
      }
    }
    if (closestCoin) {
      let dx = 0, dy = 0;
      if (closestCoin.x > npc.x) dx = 1;
      else if (closestCoin.x < npc.x) dx = -1;
      if (closestCoin.y > npc.y) dy = 1;
      else if (closestCoin.y < npc.y) dy = -1;

      tryMove(npc, [dx, dy]);

      const coinKey = `${npc.x}x${npc.y}`;
      if (coins[coinKey]) {
        delete coins[coinKey];
        npc.coins++;
      }
      return;
    } else {
      tryMove(npc, randomDir());
      return;
    }
  }

  if (npc.difficulty === 'easy' && Math.random() < 0.5) {
    return;
  }

  let closestCoin = null;
  let minDist = Infinity;
  for (const cKey in coins) {
    const c = coins[cKey];
    const dist = Math.abs(c.x - npc.x) + Math.abs(c.y - npc.y);
    if (dist < minDist) {
      minDist = dist;
      closestCoin = c;
    }
  }

  if (!closestCoin) {
    tryMove(npc, randomDir());
    return;
  }

  let dx = 0, dy = 0;
  if (closestCoin.x > npc.x) dx = 1;
  else if (closestCoin.x < npc.x) dx = -1;
  if (closestCoin.y > npc.y) dy = 1;
  else if (closestCoin.y < npc.y) dy = -1;

  let successChance = 1.0;
  if (npc.difficulty === 'easy') {
    successChance = 0.2;
  } else if (npc.difficulty === 'medium') {
    successChance = 0.85;
  }

  if (Math.random() > successChance) {
    [dx, dy] = randomDir();
  }

  if (Math.random() < 0.3) {
    dx += Math.floor(Math.random() * 3) - 1;
    dy += Math.floor(Math.random() * 3) - 1;
    dx = Math.max(-1, Math.min(1, dx));
    dy = Math.max(-1, Math.min(1, dy));
  }

  tryMove(npc, [dx, dy]);
  const coinKey = `${npc.x}x${npc.y}`;
  if (coins[coinKey]) {
    delete coins[coinKey];
    npc.coins++;
  }
}

function createNPC(index, difficulty, npcColor) {
  const npcId = 'npc_' + index + '_' + Math.floor(Math.random() * 100000);
  const x = 10 + index;
  const y = 5;

  players[npcId] = {
    id: npcId,
    socketId: null,
    name: `NPC-${index}`,
    color: npcColor,
    x,
    y,
    direction: 'left',
    coins: 0,
    isLead: false,
    ai: true,
    difficulty,
    terrorCooldown: 0,
    terrorMoveDelay: 0,
  };
}

function isNameTaken(name) {
  return Object.values(players).some((p) => p.name === name);
}

function resetPlayersAndCoins() {
  coins = {};
  for (const pid in players) {
    if (players[pid].ai) {
      delete players[pid];
    } else {
      players[pid].coins = 0;
      players[pid].x = 1;
      players[pid].y = 4;
      players[pid].direction = "right";
    }
  }
  io.emit('coins-update', coins);
  io.emit('players-update', players);
  singlePlayerMode = false;
}

function endGame() {
  gameStarted = false;
  if (timerInterval) {
    clearInterval(timerInterval);
    timerInterval = null;
  }
  if (npcInterval) {
    clearInterval(npcInterval);
    npcInterval = null;
  }

  const playersArray = Object.values(players);
  playersArray.sort((a, b) => b.coins - a.coins);

  const winner = playersArray[0];
  const scoreboard = playersArray.map((p) => ({
    name: p.name,
    coins: p.coins,
  }));

  if (!players[leadPlayerSocketId]) {
    const firstPlayer = Object.values(players)[0];
    if (firstPlayer) {
      leadPlayerSocketId = firstPlayer.socketId;
      firstPlayer.isLead = true;
    } else {
      leadPlayerSocketId = null;
    }
  }

  resetPlayersAndCoins();
  io.emit('game-ended', {
    winnerName: winner ? winner.name : "No winner",
    scoreboard,
  });
}

function startGame(terrorActive = false) {
  gameStarted = true;
  gamePaused = false;
  timeLeft = 30;
  io.emit('game-started', { hardMode: hardModeEnabled, terrorActive });
  io.emit('time-update', timeLeft);

  timerInterval = setInterval(() => {
    if (!gamePaused) {
      timeLeft--;
      io.emit('time-update', timeLeft);
      if (timeLeft <= 0) {
        clearInterval(timerInterval);
        timerInterval = null;
        endGame();
      }
    }
  }, 1000);

  if (singlePlayerMode) {
    npcInterval = setInterval(() => {
      if (!gameStarted || gamePaused) return;
      Object.values(players).forEach((p) => {
        if (p.ai) {
          moveNPC(p);
          if (p.difficulty === 'hard') {
            moveNPC(p);
          }
        }
      });
      io.emit('players-update', players);
      io.emit('coins-update', coins);
    }, 250);
  }
}

function getClosestRealPlayer(npc) {
  let closest = null;
  let minDist = Infinity;

  for (const pid in players) {
    const p = players[pid];
    if (!p.ai) {
      const dist = Math.abs(p.x - npc.x) + Math.abs(p.y - npc.y);
      if (dist < minDist) {
        minDist = dist;
        closest = p;
      }
    }
  }
  return closest;
}

io.on('connection', (socket) => {
  console.log('A user connected:', socket.id);

  if (!leadPlayerSocketId) {
    leadPlayerSocketId = socket.id;
    console.log('Lead player is now:', socket.id);
  }

  if (singlePlayerMode && socket.id !== leadPlayerSocketId) {
    socket.emit('start-failed', 'Single-player mode active. No new connections allowed.');
    socket.disconnect();
    return;
  }

  socket.on('player-join', (playerData) => {
    if (isNameTaken(playerData.name)) {
      socket.emit('name-taken', playerData.name);
      return;
    }
    playerData.socketId = socket.id;
    playerData.isLead = (socket.id === leadPlayerSocketId);

    players[playerData.id] = playerData;
    io.emit('players-update', players);
  });

  socket.on('player-update', (playerData) => {
    const existingPlayer = players[playerData.id];
    if (!existingPlayer) return;

    if (playerData.name && playerData.name !== existingPlayer.name) {
      if (isNameTaken(playerData.name)) {
        socket.emit('name-taken', playerData.name);
        return;
      }
    }
    players[playerData.id] = {
      ...existingPlayer,
      ...playerData,
      color: playerData.color || "brown",
    };
    io.emit('players-update', players);
  });

  socket.on('toggle-single-player', ({ enabled }) => {
    if (socket.id === leadPlayerSocketId) {
      singlePlayerMode = !!enabled;
      console.log("Lead toggled singlePlayerMode to:", singlePlayerMode);
    }

    if (singlePlayerMode) {
      for (const pid in players) {
        const p = players[pid];
        if (p.socketId !== leadPlayerSocketId) {
          const otherSocket = io.sockets.sockets.get(p.socketId);
          if (otherSocket) {
            otherSocket.emit('start-failed', 'Single-player mode selected. You have been disconnected.');
            otherSocket.disconnect(true);
          }
          delete players[pid];
        }
      }

      io.emit('players-update', players);
    }
  });

  socket.on('player-move', (playerData) => {
    if (players[playerData.id]) {
      players[playerData.id] = playerData;
      io.emit('players-update', players);
    }
  });

  socket.on('coin-add', (coinData) => {
    const key = `${coinData.x}x${coinData.y}`;
    coins[key] = coinData;
    io.emit('coins-update', coins);

    if (hardModeEnabled) {
      // coin vanishes after 3s in hard mode
      setTimeout(() => {
        if (coins[key]) {
          delete coins[key];
          io.emit('coins-update', coins);
        }
      }, 3000);
    }
  });

  socket.on('coin-remove', ({ x, y, playerId }) => {
    const key = `${x}x${y}`;
    if (coins[key]) {
      delete coins[key];
      if (players[playerId]) {
        players[playerId].coins++;
      }
      io.emit('coins-update', coins);
      io.emit('players-update', players);
    }
  });

  socket.on('start-game', (data) => {
    if (socket.id === leadPlayerSocketId) {
      if (gameStarted) {
        socket.emit('start-failed', 'Game is already running!');
        return;
      }
      hardModeEnabled = data && data.hardMode;
      singlePlayerMode = !!data.singlePlayerMode;

      let terrorActive = false;

      if (singlePlayerMode) {
        const npcCount = data.npcCount || 1;
        const npcDifficulty = data.npcDifficulty || 'medium';
        const npcConfigs = data.npcConfigs || [];

        for (let i = 0; i < npcCount; i++) {
          const npcColor = npcConfigs[i]?.color || 'gray';
          const difficulty = i === 0 ? npcDifficulty : npcDifficulty;
          createNPC(i, difficulty, npcColor);

          if (difficulty === 'terror') {
            terrorActive = true;
          }
        }
        const totalPlayers = Object.keys(players).length;
        if (totalPlayers >= 2 && totalPlayers <= 4) {
          startGame(terrorActive);
        } else {
          socket.emit('start-failed', 'Need between 2 and 4 total players to start single-player!');
        }

      } else {
        const numPlayers = Object.keys(players).length;
        if (numPlayers >= 2 && numPlayers <= 4) {
          startGame(false);
        } else {
          socket.emit('start-failed', 'Need between 2 and 4 players to start!');
        }
      }
    }
  });

  socket.on('pause-game', ({ name }) => {
    if (!gameStarted || gamePaused) return;
    gamePaused = true;
    io.emit('game-paused', { by: name });
  });

  socket.on('resume-game', ({ name }) => {
    if (!gameStarted || !gamePaused) return;
    gamePaused = false;
    io.emit('game-resumed', { by: name });
  });

  socket.on('quit-game', ({ playerId, name }) => {
    if (players[playerId]) {
      delete players[playerId];
    }
    io.emit('player-quit', { by: name });
    io.emit('players-update', players);

    if (socket.id === leadPlayerSocketId) {
      leadPlayerSocketId = Object.values(players)[0]?.socketId || null;
      if (leadPlayerSocketId) {
        const newLeadPlayer = Object.values(players).find(
          (p) => p.socketId === leadPlayerSocketId
        );
        if (newLeadPlayer) {
          newLeadPlayer.isLead = true;
        }
      }
    }

    const numPlayers = Object.keys(players).length;
    if (numPlayers < 2 && gameStarted) {
      clearInterval(timerInterval);
      timerInterval = null;
      endGame();
    }
  });

  socket.on('restart-game', ({ name }) => {
    io.emit('game-restarted', { by: name });
    if (timerInterval) {
      clearInterval(timerInterval);
      timerInterval = null;
    }
    if (npcInterval) {
      clearInterval(npcInterval);
      npcInterval = null;
    }
    endGame();

    const numPlayers = Object.keys(players).length;
    if (numPlayers >= 2 && numPlayers <= 4) {
      startGame();
    } else {
      io.emit('start-failed', 'Not enough players to restart.');
    }
  });

  socket.on('disconnect', () => {
    console.log('User disconnected:', socket.id);
    for (const id in players) {
      if (players[id].socketId === socket.id) {
        delete players[id];
      }
    }

    if (socket.id === leadPlayerSocketId) {
      leadPlayerSocketId = Object.values(players)[0]?.socketId || null;
      if (leadPlayerSocketId) {
        const newLeadPlayer = Object.values(players).find(
          (p) => p.socketId === leadPlayerSocketId
        );
        if (newLeadPlayer) {
          newLeadPlayer.isLead = true;
        }
      }
    }

    io.emit('players-update', players);

    const numPlayers = Object.keys(players).length;
    if (numPlayers < 2 && gameStarted) {
      if (timerInterval) {
        clearInterval(timerInterval);
        timerInterval = null;
      }
      if (npcInterval) {
        clearInterval(npcInterval);
        npcInterval = null;
      }
      endGame();
    }
  });
});

const PORT = process.env.PORT || 3000;
server.listen(PORT, () => {
  console.log(`Server running at http://localhost:${PORT}`);
});
