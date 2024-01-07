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
var bulletCount = 0;
var bulletCooldownTime= 2000
socket.onmessage = (event) => {
    const message = JSON.parse(event.data)
    console.log('message: ', message)
    const messageType = message.type
    if(messageType === "updateStatus") {
        console.log('updateStatus')
        const statusContent = message.statusContent
        displayStatusMessage(statusContent) 
        
    }  else {
    const clientInfo = message.clientInfo
    const ID = clientInfo.ID
    const currentPosition= clientInfo.Position
    console.log('id', ID)
    console.log(messageType, clientInfo.ID, currentPosition, ID, currentPosition)
    output.innerHTML += `<p>${message}</p>`;
    output.scrollTop = output.scrollHeight; // Auto-scroll to the bottom
    if(messageType === "assignPositionForNewClient") {
        const matrix = message.matrix.matrix;
        const bullets = message.bullets
        solidIndexes = new Set(matrix)
        currClientInfo = clientInfo
        competitors = message.competitors

        createMaze();
        updateScoreTable();
        createUser(ID, currentPosition, clientInfo.Direction);
        // console.log('first message')
        for (const competitor of Object.values(competitors)) {
            // console.log('competitor', competitor)
            createUser(competitor.ID, competitor.Position, competitor.Direction);
        }
        if (bullets && bullets.length) {
            for (let i = 0; i < bullets.length; i++) {
                createNewBullet(bullets[i]);
            }
        } else {
            console.error('Bullets array is null or empty.');
        }
        
    }
    if(messageType === "hasNewClient") {
        console.log('hasNewClient')
        createUser(ID, currentPosition, clientInfo.Direction);
        insertScoreTable(clientInfo)
        competitors[clientInfo.ID]= clientInfo;
    }    
    if(messageType === "removeOneClient") {
        console.log('removeOneClient')
        removeOneClient(ID);
        deleteRowInScoreTable(clientInfo); 
    }   
    //
    if(messageType === "moveOneClient") {
        console.log('moveOneClient')
        moveOneClient(clientInfo)
        competitors[clientInfo.ID].Position=clientInfo.Position
        competitors[clientInfo.ID].Direction=clientInfo.Direction
        
    }  
    if(messageType === "hasNewBullet") {
        console.log('hasNewBullet')
        const bulletInfo = message.bulletInfo;
        if(currClientInfo.ID === clientInfo.ID) {
            // currClientInfo.BulletCooldown = clientInfo.BulletCooldown
        }
        createNewBullet(bulletInfo)
        
    }  
    if(messageType === "removeOneBullet") {
        console.log('removeOneBullet')
        const bulletInfo = message.bulletInfo;
        removeOneBullet(bulletInfo.ID)
    }  
    if(messageType === "moveOneBullet") {
        console.log('moveOneBullet')
        const bulletInfo = message.bulletInfo;
        if(currClientInfo.ID === clientInfo.ID) {
            //  currClientInfo.BulletCooldown = clientInfo.BulletCooldown
        }
        moveOneBullet(bulletInfo)
    }  
    if(messageType === "updateScoreOfOneClient") {
        console.log('updateScoreOfOneClient')

        updateRowInScoreTable(clientInfo)
        competitors[clientInfo.ID]= clientInfo;
        if (clientInfo.ID == currClientInfo.ID) {
            currClientInfo.Score = clientInfo.Score
        }
    }  
    
    }
    
  
    //return false;  

}
socket.onclose = (event) => {
    console.log('WebSocket connection closed:', event.code, event.reason);
  };
// socket.onclose = (event) => {
  
//     // if(event.wasClean) {
//     //     console.log(`Closed cleanly, code = ${event.code}, reason = ${event.reason}`);

//     // }else {
//     //     console.log("Connection died");

//     // }
// }

socket.onerror = (event) => {
    console.log(`WebSocket connection error: ${event}`);

}


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
    const userId = "user_" + ID;
    const oldUser = document.getElementById(userId)
    if (oldUser) {
        oldUser.remove()
    }
    var _user = document.createElement("div");
    _user.classList.add("user");
    _user.classList.add(direction)
    _user.id = userId

    const mazeContainer = document.getElementById("maze-container");
    mazeContainer.appendChild(_user);

    const x = Math.floor(position % 32);
    const y = Math.floor(position / 32);
    console.log('type of x', typeof(x))

    _user.style.position = "absolute";
    _user.style.left = (20 * x) + "px";
    _user.style.top = (20 * y) + "px";

}
function createNewBullet(bulletInfo) {

    const bulletId = bulletInfo.ID;
    const position = bulletInfo.Position;
    const direction = bulletInfo.Direction;
    const oldBullet = document.getElementById(bulletId)
    if (oldBullet) {
        oldBullet.remove()
    }
    var _bullet = document.createElement("div");
    _bullet.classList.add("bullet");
    _bullet.classList.add(direction)
    _bullet.id = "bullet_" + bulletId

    const mazeContainer = document.getElementById("maze-container");
    mazeContainer.appendChild(_bullet);

    const x = Math.floor(position % 32);
    const y = Math.floor(position / 32);
    console.log('type of x', x, y)

    _bullet.style.position = "absolute";
    _bullet.style.left = (20 * x) + "px";
    _bullet.style.top = (20 * y) + "px";

}
function updateScoreTable() {
    const tableBody = document.querySelector('#scoreTable tbody');
    tableBody.innerHTML = ''; // Clear existing rows

    // Iterate over the values of the competitors object
    for (const competitor of Object.values(competitors)) {
        const rowId = "rowScore_" + competitor.ID
        const row = document.createElement('tr');
        row.id= rowId;
        row.innerHTML = `
            <td>${competitor.ID}</td>
            <td>${competitor.Score}</td>
            <td>${competitor.Position}</td>
            <td>${competitor.Direction}</td>
        `;
        tableBody.appendChild(row);
    }
}
function insertScoreTable(clientInfo) {
    const tableBody = document.querySelector('#scoreTable tbody');
    //tableBody.innerHTML = ''; // Clear existing rows
    const row = document.createElement('tr');
    const rowId = "rowScore_" + clientInfo.ID
    row.id = rowId;
    row.innerHTML = `
        <td>${clientInfo.ID}</td>
        <td>${clientInfo.Score}</td>
        <td>${clientInfo.Position}</td>
        <td>${clientInfo.Direction}</td>
    `;
    tableBody.appendChild(row);
}
function updateRowInScoreTable(clientInfo) {
    const rowId = "rowScore_" + clientInfo.ID
    const row = document.getElementById(rowId);
    if (row) {
        // Update the content of the existing row
        row.innerHTML = `
            <td>${clientInfo.ID}</td>
            <td>${clientInfo.Score}</td>
            <td>${clientInfo.Position}</td>
            <td>${clientInfo.Direction}</td>
        `;
    } else {
        // Row doesn't exist, create a new row
        const tableBody = document.querySelector('#scoreTable tbody');
        const newRow = document.createElement('tr');
        newRow.id = rowId;
        newRow.innerHTML = `
            <td>${clientInfo.ID}</td>
            <td>${clientInfo.Score}</td>
            <td>${clientInfo.Position}</td>
            <td>${clientInfo.Direction}</td>
        `;
        tableBody.appendChild(newRow);
    }
}
function deleteRowInScoreTable(clientInfo) {
    const rowId = "rowScore_" + clientInfo.ID
    const row = document.getElementById(rowId);
    if(row){
        row.remove();
    }
   
}
function displayStatusMessage(message) {
    // Create a new <p> element
    var pElement = document.createElement("p");
    pElement.textContent = message;

    // Get the status div
    var statusDiv = document.getElementById("status");

    // Append the <p> element to the status div
    statusDiv.appendChild(pElement);

    // Scroll to the bottom to show the latest message
    //statusDiv.scrollTop = statusDiv.scrollHeight;
}
function removeOneClient(ID) {

    const userId = "user_" + ID;
    const user = document.getElementById(userId);

    if (user) {
        console.log('remove one client')
        user.remove();
    } else {
        console.error(`Element with ID ${userId} not found.`);
    }
}
function removeOneBullet(ID) {
    const bulletId = "bullet_" + ID;
    const bullet = document.getElementById(bulletId);

    if (bullet) {
        console.log('remove one bullet')
        bullet.remove();
    } else {
        console.error(`Element with ID ${bulletId} not found.`);
    }
}
function moveOneClient(clientInfo){
    const userId = "user_" + clientInfo.ID;
    const _user = document.getElementById(userId)
    if (_user) {
        const position = clientInfo.Position;
        const x = Math.floor(position % 32);
        const y = Math.floor(position / 32);
        console.log('type of x', typeof(x))

        _user.style.position = "absolute";
        _user.style.left = (20 * x) + "px";
        _user.style.top = (20 * y) + "px";
        if(_user.classList.contains("Left")){
            _user.classList.remove("Left");
        }
        if(_user.classList.contains("Right")){
            _user.classList.remove("Right");
        }
        if(_user.classList.contains("Up")){
            _user.classList.remove("Up");
        }
        if(_user.classList.contains("Down")){
            _user.classList.remove("Down");
        }
        _user.classList.add(clientInfo.Direction)
    }
}
function moveOneBullet(bulletInfo){
    const bulletId = "bullet_" + bulletInfo.ID;
    const _bullet = document.getElementById(bulletId)
    console.log('bulletId', bulletId)
    if (_bullet) {
        const position = bulletInfo.Position;
        const x = Math.floor(position % 32);
        const y = Math.floor(position / 32);
        console.log('type of x', typeof(x))

        _bullet.style.position = "absolute";
        _bullet.style.left = (20 * x) + "px";
        _bullet.style.top = (20 * y) + "px";
        // if(_bullet.classList.contains("Left")){
        //     _bullet.classList.remove("Left");
        // }
        // if(_bullet.classList.contains("Right")){
        //     _bullet.classList.remove("Right");
        // }
        // if(_bullet.classList.contains("Up")){
        //     _bullet.classList.remove("Up");
        // }
        // if(_bullet.classList.contains("Down")){
        //     _bullet.classList.remove("Down");
        // }
        //_bullet.classList.add(bulletInfo.Direction)
    }
}
document.addEventListener("DOMContentLoaded", function() {
    let shootingInProgress = false;
    // Event listener for arrow key presses
    document.addEventListener("keydown", event => {
        const step = 20;
        const currentPosition = currClientInfo.Position
        const direction = event.key.slice(5,event.key.length)
        var keyType;
        if(!isPositionOccupiedByCompetitor(determineNewPositionByDirection(currentPosition, direction))
        && !isPositionOccupiedByWall(currentPosition, direction)) {
            const oldDirection = currClientInfo.Direction;
            switch (event.key) {
            
                case "ArrowLeft":
                    moveLeft(step, oldDirection);
                    keyType ="Move"
                    break;
                case "ArrowRight":
                    moveRight(step, oldDirection);
                    keyType ="Move"
                    break;
                case "ArrowUp":
                    moveUp(step, oldDirection);
                    keyType ="Move"
                    break;
                case "ArrowDown":
                    moveDown(step, oldDirection);
                    keyType ="Move"
                    break;
                case "Q" || "q":
                    console.log('quit')
                    // Handle Q key
                    keyType = "Quit";
                    break;   
                default:
                // Handle Q key
                    console.log('Cannot recognize action')
                    break;    
            }
            if (keyType=="Move"){
                currClientInfo.Position = determineNewPositionByDirection(currentPosition, direction); 
                currClientInfo.Direction = direction;
                const message = {
                    type: "clientRequestMoving",
                    clientInfo: currClientInfo
                }
                while (true){
                    if(!shootingInProgress){
                        socket.send(JSON.stringify(message));
                        break;
                    }
                 
                }
               
            }
           
          
           
        }
       
        event.preventDefault();
     
    });
    
    document.addEventListener('keyup', event => {
        if (event.code === 'Space') {
        //shootingInProgress = false;
            console.log('Space pressed')
            console.log('shoot')
            console.log('bullet cooldown', currClientInfo.BulletCooldown);
            if (!shootingInProgress && isAllowedToShoot(currClientInfo)) {
                shootingInProgress = true;
                console.log('is allowed to shoot')
                const bulletId =  bulletCount + currClientInfo.ID*100;
                const bulletDirection = currClientInfo.Direction 
                const bulletPosition = determineNewPositionByDirection(currClientInfo.Position, bulletDirection)
                console.log('bullet info', bulletDirection); //bulletPosition
                console.log('bullet info', bulletPosition); //
                const bulletInfo = {
                    ID: bulletId,
                    ClientID: currClientInfo.ID,
                    Position: bulletPosition,
                    Direction: bulletDirection
                }
                //createNewBullet(bulletInfo)
                const clientRequestShootingMessage = {
                    type: "clientRequestShooting",
                    clientInfo: currClientInfo,
                    bulletInfo: bulletInfo
                }
                
                socket.send(JSON.stringify(clientRequestShootingMessage));
                //currClientInfo.BulletCooldown = 4;

                const operationTimeout = setTimeout(() => {
                console.log("Timeout exceeded. Operation aborted.");
                shootingInProgress = false;
                currClientInfo.BulletCooldown = 0;
                }, bulletCooldownTime);
                
                // Simulate the operation (replace with your actual code)
                // Clear the timeout to prevent it from triggering

                bulletCount++;
            
        
            } else {
                console.log('Cannot shoot. Cooldown in progress.');
            }



        }
        if (event.code==="Q"){
            console.log('want to quiet')
            const message = {
                type: "clientRequestLoggingOut",
                clientInfo: currClientInfo
            }
            socket.send(JSON.stringify(message));
        }

      })
    
});

function isAllowedToShoot (clientInfo){
    const position = clientInfo.Position
    const direction = clientInfo.Direction
    const bulletCooldown = clientInfo.BulletCooldown
    if (!isPositionOccupiedByCompetitor(position) 
    && !isPositionOccupiedByWall(position, direction)
    && bulletCooldown == 0
    ){
        return true;
    }
    return false
}
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
function isPositionOccupiedByCompetitor(position) {
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