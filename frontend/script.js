import { matrix1 } from "./matrix.js";
console.log("Start of client script");
console.log("Before WebSocket connection");
const socket = new WebSocket("ws://localhost:8081/ws");
console.log("After WebSocket connection");

socket.addEventListener("error", (event) => {
    console.log("WebSocket error: ", event);
  });
// Event handler when the WebSocket connection is opened.
socket.onopen = () => {
    console.log("WebSocket connection opened");
};
var solidIndexes
var currClientInfo
var competitors
socket.onmessage = (event) => {
    const message = JSON.parse(event.data)
    console.log('message: ', message)
    const output = document.getElementById("output");
    const messageType = message.type
    const clientInfo = message.clientInfo
    const ID = clientInfo.ID
    const currentPosition= clientInfo.Position
    console.log('id', ID)
    console.log(messageType, clientInfo.ID, currentPosition, ID, currentPosition)
    output.innerHTML += `<p>${message}</p>`;
    output.scrollTop = output.scrollHeight; // Auto-scroll to the bottom
    if(messageType === "assignPositionForNewClient") {
        const matrix = message.matrix.matrix;
        solidIndexes = new Set(matrix)
        currClientInfo = clientInfo
        competitors = message.competitors
         createMaze();
         createUser(ID, currentPosition, clientInfo.Direction);
    }
    if(messageType === "hasNewClient") {
        console.log('hasnewclient')
         createUser(ID, currentPosition, clientInfo.Direction);
    }    
    return false;  

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
        const _solidIndexes = solidIndexes;

        if (_solidIndexes.has(i)) {
            cell.classList.add("solid");
        }


        mazeContainer.appendChild(cell);
    }
}
function createUser(ID, position, direction) {
    var _user = document.createElement("div");
    _user.classList.add("user");
    _user.classList.add(direction)
    _user.id = "user_" + ID;
   // _user.innerHTML = `<p>${ID}</p>`;
    const mazeContainer = document.getElementById("maze-container");
    mazeContainer.appendChild(_user);

    const x = Math.floor(position % 32);
    const y = Math.floor(position / 32);
    console.log('type of x', typeof(x))
    // console.log("x: " + x + " y: " + y);
    // console.log("x: " + x + " y: " + y);
    // console.log("x*20")
    _user.style.position = "absolute";
    _user.style.left = (20 * x) + "px";
    _user.style.top = (20 * y) + "px";

}


document.addEventListener("DOMContentLoaded", function() {
   
    // Event listener for arrow key presses
    document.addEventListener("keydown", event => {
        event.preventDefault();
        const step = 20;
        const currentPosition = currClientInfo.Position
        const direction = event.key.slice(5,event.key.length)
        if(!isPositionOccupiedByCompetitor(determineNewPositionByDirection(currentPosition, direction), competitors)
        && !isPositionOccupiedByWall(currentPosition, direction)) {
            const oldDirection = currClientInfo.Direction;
            switch (event.key) {
            
                case "ArrowLeft":
                    moveLeft(step, oldDirection);
                    break;
                case "ArrowRight":
                    moveRight(step, oldDirection);
                    break;
                case "ArrowUp":
                    moveUp(step, oldDirection);
                    break;
                case "ArrowDown":
                    moveDown(step, oldDirection);
                    break;
            }
            currClientInfo.Position = determineNewPositionByDirection(currentPosition, direction); 
           
            currClientInfo.Direction = direction;
        }
       
        event.preventDefault();
     
    });
    
});


function determineNewPositionByDirection(currPosition, direction){
    var newPosition
    switch(direction){
        case "Up":
            newPosition = currPosition - 32;
            break;
        case "Down":
            newPosition = currPosition + 32
            break;
        case "Left":
            newPosition = currPosition - 1
            break;
        case "Right":
            newPosition = currPosition + 1
            break;
    }
    return newPosition;
}
function isPositionOccupiedByCompetitor(position, competitorsArray) {
    for (let i = 0; i < competitors.length; i++) {
      if (competitors[i].Position === position) {
        console.log('competitor position', competitors[i].Position);
        return true; // Position is occupied
      }
    }
    return false; // Position is not occupied
}
function isPositionOccupiedByWall(currPosition, direction) {

    const newPosition = determineNewPositionByDirection(currPosition, direction) 
    switch(direction){
        case "Up":
            if(newPosition >=0 && !solidIndexes.has(newPosition)){
                return false;
            }
            break;
        case "Down":
            if(newPosition < 32*16 && !solidIndexes.has(newPosition)){
                return false;
            }
            break;
        case "Left":
            if(currPosition % 32!==0 && !solidIndexes.has(newPosition)){
                return false;
            }
            break;
          
        case "Right":
            if(currPosition % 32!==31 && !solidIndexes.has(newPosition)){
                return false;
            }
            break;
    }

    return true;
}
function moveLeft(step, oldDirection) {
   
    const userId = "user_" + currClientInfo.ID;
    const user = document.getElementById(userId);
    const currentPosition = parseInt(user.style.left) || 0;
    user.style.left = currentPosition - step + "px";
    user.classList.remove(oldDirection);
    user.classList.add("Left")
}

function moveRight(step, oldDirection) {
    const userId = "user_" + currClientInfo.ID;
    const user = document.getElementById(userId);
    const currentPosition = parseInt(user.style.left) || 0;
    user.style.left = currentPosition + step + "px";
    user.classList.remove(oldDirection);
    user.classList.add("Right")
}

function moveUp(step, oldDirection) {
    const userId = "user_" + currClientInfo.ID;
    const user = document.getElementById(userId);
    const currentPosition = parseInt(user.style.top) || 0;
    user.style.top = currentPosition - step + "px";
    user.classList.remove(oldDirection);
    user.classList.add("Up")
    
}

function moveDown(step, oldDirection) {
    const userId = "user_" + currClientInfo.ID;
    const user = document.getElementById(userId);
    const currentPosition = parseInt(user.style.top) || 0;
    user.style.top = currentPosition + step + "px";
    user.classList.remove(oldDirection);
    user.classList.add("Down")
}