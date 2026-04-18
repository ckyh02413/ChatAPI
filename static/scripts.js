let currentRoomId = null;
let currentRoomName = null;
let currentUsername = null;
let pollingInterval = null;

const avatarColors = [
  { bg: "#dde4f0", color: "#2a7ae4" },
  { bg: "#e1f5ee", color: "#0f6e56" },
  { bg: "#eeedfe", color: "#534ab7" },
  { bg: "#faece7", color: "#993c1d" },
  { bg: "#fbeaf0", color: "#993556" },
  { bg: "#faeeda", color: "#854f0b" },
];

function getAvatarColor(name) {
  let hash = 0;
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash);
  }
  return avatarColors[Math.abs(hash) % avatarColors.length];
}

function getInitials(name) {
  return name.slice(0, 2).toUpperCase();
}

function getToken() {
  return localStorage.getItem("token");
}

// auth toggle
function showRegister() {
  document.getElementById("login-form").style.display = "none";
  document.getElementById("register-form").style.display = "block";
  document.getElementById("login-error").textContent = "";
}

function showLogin() {
  document.getElementById("register-form").style.display = "none";
  document.getElementById("login-form").style.display = "block";
  document.getElementById("reg-error").textContent = "";
}

// auth
async function login() {
  const username = document.getElementById("login-username").value;
  const password = document.getElementById("login-password").value;
  const errorEl = document.getElementById("login-error");

  const response = await fetch("/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });

  if (response.ok) {
    const data = await response.json();
    localStorage.setItem("token", data.token);
    currentUsername = username;
    enterChat();
  } else {
    const error = await response.text();
    errorEl.textContent = error.trim();
  }
}

async function register() {
  const username = document.getElementById("reg-username").value;
  const mail = document.getElementById("reg-mail").value;
  const password = document.getElementById("reg-password").value;
  const errorEl = document.getElementById("reg-error");

  const response = await fetch("/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, mail, password }),
  });

  if (response.ok) {
    const data = await response.json();
    localStorage.setItem("token", data.token);
    currentUsername = username;
    enterChat();
  } else {
    const error = await response.text();
    errorEl.textContent = error.trim();
  }
}

function enterChat() {
  document.getElementById("auth-screen").style.display = "none";
  document.getElementById("chat-screen").style.display = "block";

  const avatar = document.getElementById("user-avatar");
  const ac = getAvatarColor(currentUsername);
  avatar.style.background = ac.bg;
  avatar.style.color = ac.color;
  avatar.textContent = getInitials(currentUsername);
  document.getElementById("user-name").textContent = currentUsername;

  loadRooms();
}

function logout() {
  localStorage.removeItem("token");
  currentUsername = null;
  currentRoomId = null;
  if (pollingInterval) clearInterval(pollingInterval);
  document.getElementById("chat-screen").style.display = "none";
  document.getElementById("auth-screen").style.display = "flex";
  document.getElementById("login-username").value = "";
  document.getElementById("login-password").value = "";
  showLogin();
}

// rooms
async function loadRooms() {
  const response = await fetch("/rooms/", {
    headers: { Authorization: "Bearer " + getToken() },
  });

  const roomList = document.getElementById("room-list");

  if (response.ok) {
    const data = await response.json();
    roomList.innerHTML = "";

    if (data) {
      data.forEach((room) => {
        const div = document.createElement("div");
        div.className =
          "room-item" + (room.id === currentRoomId ? " active" : "");
        div.innerHTML = `
                        <span class="room-name"># ${room.name}</span>
                        <div class="room-actions">
                            <button onclick="event.stopPropagation(); editRoom(${room.id})">edit</button>
                            <button onclick="event.stopPropagation(); deleteRoom(${room.id})">del</button>
                        </div>
                    `;
        div.onclick = () => openRoom(room.id, room.name);
        roomList.appendChild(div);
      });
    }
  }
}

async function createRoom() {
  const input = document.getElementById("room-name-input");
  const name = input.value.trim();
  if (!name) return;

  const response = await fetch("/rooms/", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: "Bearer " + getToken(),
    },
    body: JSON.stringify({ name }),
  });

  if (response.ok) {
    input.value = "";
    loadRooms();
  } else {
    const error = await response.text();
    alert(error.trim());
  }
}

async function deleteRoom(id) {
  if (!confirm("Delete this room?")) return;

  const response = await fetch("/rooms/" + id, {
    method: "DELETE",
    headers: { Authorization: "Bearer " + getToken() },
  });

  if (response.ok) {
    if (currentRoomId === id) {
      currentRoomId = null;
      document.getElementById("chat-active").style.display = "none";
      if (pollingInterval) clearInterval(pollingInterval);
    }
    loadRooms();
  } else {
    const error = await response.text();
    alert(error.trim());
  }
}

async function editRoom(id) {
  const newName = prompt("Enter new room name:");
  if (!newName) return;

  const response = await fetch("/rooms/" + id, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
      Authorization: "Bearer " + getToken(),
    },
    body: JSON.stringify({ name: newName }),
  });

  if (response.ok) {
    if (currentRoomId === id) {
      currentRoomName = newName;
      document.getElementById("chat-room-name").textContent = "# " + newName;
    }
    loadRooms();
  } else {
    const error = await response.text();
    alert(error.trim());
  }
}

// messages
function openRoom(id, name) {
  currentRoomId = id;
  currentRoomName = name;

  document.getElementById("chat-empty").style.display = "none";
  document.getElementById("chat-active").style.display = "flex";
  document.getElementById("chat-room-name").textContent = "# " + name;

  loadMessages();
  loadRooms();

  if (pollingInterval) clearInterval(pollingInterval);
  pollingInterval = setInterval(loadMessages, 3000);
}

async function loadMessages() {
  const response = await fetch("/rooms/" + currentRoomId + "/messages", {
    headers: { Authorization: "Bearer " + getToken() },
  });

  const container = document.getElementById("chat-messages");
  const countEl = document.getElementById("msg-count");

  if (response.ok) {
    const data = await response.json();
    container.innerHTML = "";

    if (data && data.length > 0) {
      countEl.textContent =
        data.length + " message" + (data.length !== 1 ? "s" : "");

      data.forEach((msg) => {
        const ac = getAvatarColor(msg.creator);
        const div = document.createElement("div");
        div.className = "message";
        div.innerHTML = `
                        <div class="msg-avatar" style="background: ${ac.bg}; color: ${ac.color};">${getInitials(msg.creator)}</div>
                        <div class="msg-body">
                            <div class="msg-header">
                                <span class="msg-author">${msg.creator}</span>
                            </div>
                            <p class="msg-content">${msg.content}</p>
                            <div class="msg-actions">
                                <button onclick="editMessage(${msg.id})">edit</button>
                                <button onclick="deleteMessage(${msg.id})">delete</button>
                            </div>
                        </div>
                    `;
        container.appendChild(div);
      });

      container.scrollTop = container.scrollHeight;
    } else {
      countEl.textContent = "no messages yet";
    }
  }
}

async function sendMessage() {
  const input = document.getElementById("msg-input");
  const content = input.value.trim();
  if (!content) return;

  const response = await fetch("/rooms/" + currentRoomId + "/messages", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: "Bearer " + getToken(),
    },
    body: JSON.stringify({ content }),
  });

  if (response.ok) {
    input.value = "";
    loadMessages();
  } else {
    const error = await response.text();
    alert(error.trim());
  }
}

async function deleteMessage(msgId) {
  const response = await fetch(
    "/rooms/" + currentRoomId + "/messages/" + msgId,
    {
      method: "DELETE",
      headers: { Authorization: "Bearer " + getToken() },
    },
  );

  if (response.ok) {
    loadMessages();
  } else {
    const error = await response.text();
    alert(error.trim());
  }
}

async function editMessage(msgId) {
  const newContent = prompt("Edit message:");
  if (!newContent) return;

  const response = await fetch(
    "/rooms/" + currentRoomId + "/messages/" + msgId,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer " + getToken(),
      },
      body: JSON.stringify({ content: newContent }),
    },
  );

  if (response.ok) {
    loadMessages();
  } else {
    const error = await response.text();
    alert(error.trim());
  }
}

async function checkAuth() {
  const token = getToken();
  if (!token) return;

  const response = await fetch("/me", {
    headers: { Authorization: "Bearer " + token },
  });

  if (response.ok) {
    const data = await response.json();
    currentUsername = data.username;
    enterChat();
  } else {
    localStorage.removeItem("token");
  }
}

checkAuth();
