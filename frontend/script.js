import { matrix1 } from "./matrix.js";
var socket
function connect() {
    socket = new WebSocket('ws://10.20.219.170:8081/')
    socket.onclose = function() { connect(); };
}
connect()
socket.onopen = () => {
    console.log("Websocket connection opened");
}
var solidIndexes
socket.onmessage = (event) => {
    const message = JSON.parse(event.data)
    const output = document.getElementById("output");
    const messageContent = message.content;
    const messageType = message.type
    const currentPosition = message.position
    const ID = message.clientId
    console.log('id', ID)
    output.innerHTML += `<p>${message}</p>`;
    output.innerHTML += `<p>${messageContent}</p>`;
    output.scrollTop = output.scrollHeight; // Auto-scroll to the bottom
    if(messageType === "AssignPosition") {
        const matrix = message.matrix.matrix;
        solidIndexes = matrix
        console.log('solidIndexes',solidIndexes)
        createMaze();
        createUser(currentPosition, ID );
        
    }
    

}

socket.onclose = (event) => {
    if(event.wasClean) {
        console.log(`Closed cleanly, code = ${event.code}, reason = ${event.reason}`);

    }else {
        console.log("Connection died");

    }
}

socket.onerror = (event) => {
    console.log(`WebSocket connection error: ${event}`);

}

// document.getElementById("message-form").addEventListener("submit", (event) =>{
//     event.preventDefault();
//     const messageInput = document.getElementById("message");
//     const message = messageInput.value;
//     socket.send(message);
//     messageInput.value="";
// })
function createMaze() {
 
    const mazeContainer = document.getElementById("maze");

    for (let i = 0; i < 32 * 16; i++) {
        const cell = document.createElement("div");
        cell.classList.add("cell");
        const _solidIndexes = new Set(solidIndexes);
        if (_solidIndexes.has(i)) {
            cell.classList.add("solid");
        }


        mazeContainer.appendChild(cell);
    }
}
function createUser(position,ID) {
    const user = document.createElement("div");
    user.classList.add("user");
    user.innerHTML = `<p>${ID}</p>`
    const mazeContainer = document.getElementById("maze-container");
    mazeContainer.appendChild(user);
    // user.style.left = "0px";
    // user.style.top = "0px";
  
    const x = Math.floor(position%32) 
    const y = Math.floor(position/32) 
    console.log("x: " + x + " y: " + y)
    user.style.left = 20*x + "px";
    user.style.top = 20*y + "px";
}
var flag = true;

document.addEventListener("DOMContentLoaded", function() {
   
    // Event listener for arrow key presses
    document.addEventListener("keydown", event => {
        event.preventDefault();
        const step = 0;
        switch (event.key) {
            case "ArrowLeft":
                moveLeft(step);
                break;
            case "ArrowRight":
                moveRight(step);
                break;
            case "ArrowUp":
                moveUp(step);
                break;
            case "ArrowDown":
                moveDown(step);
                break;
        }
        event.preventDefault();
     
    });
    
});
var user 


function moveLeft(step) {
    user = document.querySelector(".user");
    const currentPosition = parseInt(user.style.left) || 0;
    user.style.left = currentPosition - step + "px";
}

function moveRight(step) {
    user = document.querySelector(".user");
    const currentPosition = parseInt(user.style.left) || 0;
    user.style.left = currentPosition + step + "px";
}

function moveUp(step) {
    user = document.querySelector(".user");
    const currentPosition = parseInt(user.style.top) || 0;
    user.style.top = currentPosition - step + "px";
}

function moveDown(step) {
    user = document.querySelector(".user");
    const currentPosition = parseInt(user.style.top) || 0;
    user.style.top = currentPosition + step + "px";
}