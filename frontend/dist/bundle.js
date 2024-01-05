/*
 * ATTENTION: The "eval" devtool has been used (maybe by default in mode: "development").
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
/******/ (() => { // webpackBootstrap
/******/ 	"use strict";
/******/ 	var __webpack_modules__ = ({

/***/ "./matrix.js":
/*!*******************!*\
  !*** ./matrix.js ***!
  \*******************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   matrix1: () => (/* binding */ matrix1)\n/* harmony export */ });\nconst matrix1 = new Set([\r\n    5, 12, 18, 25, 30, 42, 50, 55, 60, 75,\r\n    80, 90, 95, 110, 120, 130, 140, 150, 160, 170,\r\n    200, 210, 220, 230, 240, 267, 289, 290, 292, 293,\r\n    300, 310, 320, 330, 340, 367, 389, 390, 392, 393,\r\n    400, 410, 420, 430, 440, 467, 489, 490, 492, 493,\r\n    500, 510, 520, 530, 540, 567, 589, 590, 592, 593,\r\n    600, 610, 620, 630, 640, 667, 689, 690, 692, 693,\r\n    700, 710, 720, 730, 740, 767, 789, 790, 792, 793,\r\n    800, 810, 820, 830, 840, 867, 889, 890, 892, 893,\r\n    900, 910, 920, 930, 940, 967, 989, 990, 992, 993,\r\n    12, 14, 34, 32, 38, 43, 51, 56, 72, 82,\r\n    87, 98, 106, 117, 122, 129, 150, 166, 177, 183, \r\n    209, 213, 225, 243, 249, 275, 308, 309, 305, 294, \r\n    312, 317, 321, 344, 352, 361, 386, 391, 409, 404, \r\n    394, 411, 432, 437, 458, 480, 489, 497, 518, 531, \r\n    559, 571, 581, 585, 599, 616, 631, 648, 653, 669, \r\n    691, 696, 708, 721, 744, 763, 782, 784, 799, 816, \r\n    826, 831, 852, 870, 883, 904, 925, 950, 966, 967, 987, \r\n    1003, 1008, 1021, 1035, 1048, 1054, 1067, 1087, 1095, 1110, \r\n    1114, 1126, 1144, 1167, 1183, 1188, 1200, 1215, 1235, 1236, \r\n    1263, 1285, 1300, 1312, 1328, 1352, 1366, 1382, 1388, 1400, \r\n    1414, 1439, 1454, 1468, 1470, 1482, 1483, 1486, 1487, 1490\r\n]);\r\n\n\n//# sourceURL=webpack://frontend/./matrix.js?");

/***/ }),

/***/ "./script.js":
/*!*******************!*\
  !*** ./script.js ***!
  \*******************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var _matrix_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./matrix.js */ \"./matrix.js\");\n\r\nconst socket = new WebSocket(\"ws://localhost:8081/ws\");\r\n\r\nsocket.onopen = () => {\r\n    console.log(\"Websocket connection opened\");\r\n}\r\nvar solidIndexes\r\nsocket.onmessage = (event) => {\r\n    const message = JSON.parse(event.data)\r\n    const output = document.getElementById(\"output\");\r\n    const messageContent = message.content;\r\n    const messageType = message.type\r\n    const currentPosition = message.position\r\n    const ID = message.clientId\r\n    console.log('id', ID)\r\n    output.innerHTML += `<p>${messageContent}</p>`;\r\n    output.scrollTop = output.scrollHeight; // Auto-scroll to the bottom\r\n    if(messageType === \"AssignPosition\") {\r\n        const matrix = message.matrix.matrix;\r\n        solidIndexes = matrix\r\n        console.log('solidIndexes',solidIndexes)\r\n        createMaze();\r\n        createUser(currentPosition, ID );\r\n        \r\n    }\r\n    \r\n\r\n}\r\n\r\nsocket.onclose = (event) => {\r\n    if(event.wasClean) {\r\n        console.log(`Closed cleanly, code = ${event.code}, reason = ${event.reason}`);\r\n\r\n    }else {\r\n        console.log(\"Connection died\");\r\n\r\n    }\r\n}\r\n\r\nsocket.onerror = (event) => {\r\n    console.log(`WebSocket connection error: ${event}`);\r\n\r\n}\r\n\r\n// document.getElementById(\"message-form\").addEventListener(\"submit\", (event) =>{\r\n//     event.preventDefault();\r\n//     const messageInput = document.getElementById(\"message\");\r\n//     const message = messageInput.value;\r\n//     socket.send(message);\r\n//     messageInput.value=\"\";\r\n// })\r\nfunction createMaze() {\r\n \r\n    const mazeContainer = document.getElementById(\"maze\");\r\n\r\n    for (let i = 0; i < 32 * 16; i++) {\r\n        const cell = document.createElement(\"div\");\r\n        cell.classList.add(\"cell\");\r\n        const _solidIndexes = new Set(solidIndexes);\r\n        if (_solidIndexes.has(i)) {\r\n            cell.classList.add(\"solid\");\r\n        }\r\n\r\n\r\n        mazeContainer.appendChild(cell);\r\n    }\r\n}\r\nfunction createUser(position,ID) {\r\n    const user = document.createElement(\"div\");\r\n    user.classList.add(\"user\");\r\n    user.innerHTML = `<p>${ID}</p>`\r\n    const mazeContainer = document.getElementById(\"maze-container\");\r\n    mazeContainer.appendChild(user);\r\n    user.style.left = \"0px\";\r\n    user.style.top = \"0px\";\r\n  \r\n    const x = Math.floor(position%32) \r\n    const y = Math.floor(position/32) \r\n    console.log(\"x: \" + x + \" y: \" + y)\r\n    // user.style.left = 20*x + \"px\";\r\n    // user.style.top = 20*y + \"px\";\r\n}\r\nvar flag = true;\r\n\r\ndocument.addEventListener(\"DOMContentLoaded\", function() {\r\n   event.preventDefault();\r\n    // Event listener for arrow key presses\r\n    document.addEventListener(\"keydown\", event => {\r\n\r\n        const step = 20;\r\n        switch (event.key) {\r\n            case \"ArrowLeft\":\r\n                moveLeft(step);\r\n                break;\r\n            case \"ArrowRight\":\r\n                moveRight(step);\r\n                break;\r\n            case \"ArrowUp\":\r\n                moveUp(step);\r\n                break;\r\n            case \"ArrowDown\":\r\n                moveDown(step);\r\n                break;\r\n        }\r\n        event.preventDefault();\r\n        // var badKey = [37,38,39,40]; //down array keyCode\r\n        // if (flag && badKey.includes(e.keyCode)) {\r\n        //   e.preventDefault();\r\n        // }\r\n    });\r\n    \r\n});\r\nvar user \r\n\r\n\r\nfunction moveLeft(step) {\r\n    user = document.querySelector(\".user\");\r\n    const currentPosition = parseInt(user.style.left) || 0;\r\n    user.style.left = currentPosition - step + \"px\";\r\n}\r\n\r\nfunction moveRight(step) {\r\n    user = document.querySelector(\".user\");\r\n    const currentPosition = parseInt(user.style.left) || 0;\r\n    user.style.left = currentPosition + step + \"px\";\r\n}\r\n\r\nfunction moveUp(step) {\r\n    user = document.querySelector(\".user\");\r\n    const currentPosition = parseInt(user.style.top) || 0;\r\n    user.style.top = currentPosition - step + \"px\";\r\n}\r\n\r\nfunction moveDown(step) {\r\n    user = document.querySelector(\".user\");\r\n    const currentPosition = parseInt(user.style.top) || 0;\r\n    user.style.top = currentPosition + step + \"px\";\r\n}\n\n//# sourceURL=webpack://frontend/./script.js?");

/***/ })

/******/ 	});
/************************************************************************/
/******/ 	// The module cache
/******/ 	var __webpack_module_cache__ = {};
/******/ 	
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/ 		// Check if module is in cache
/******/ 		var cachedModule = __webpack_module_cache__[moduleId];
/******/ 		if (cachedModule !== undefined) {
/******/ 			return cachedModule.exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = __webpack_module_cache__[moduleId] = {
/******/ 			// no module.id needed
/******/ 			// no module.loaded needed
/******/ 			exports: {}
/******/ 		};
/******/ 	
/******/ 		// Execute the module function
/******/ 		__webpack_modules__[moduleId](module, module.exports, __webpack_require__);
/******/ 	
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/ 	
/************************************************************************/
/******/ 	/* webpack/runtime/define property getters */
/******/ 	(() => {
/******/ 		// define getter functions for harmony exports
/******/ 		__webpack_require__.d = (exports, definition) => {
/******/ 			for(var key in definition) {
/******/ 				if(__webpack_require__.o(definition, key) && !__webpack_require__.o(exports, key)) {
/******/ 					Object.defineProperty(exports, key, { enumerable: true, get: definition[key] });
/******/ 				}
/******/ 			}
/******/ 		};
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/hasOwnProperty shorthand */
/******/ 	(() => {
/******/ 		__webpack_require__.o = (obj, prop) => (Object.prototype.hasOwnProperty.call(obj, prop))
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/make namespace object */
/******/ 	(() => {
/******/ 		// define __esModule on exports
/******/ 		__webpack_require__.r = (exports) => {
/******/ 			if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 				Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 			}
/******/ 			Object.defineProperty(exports, '__esModule', { value: true });
/******/ 		};
/******/ 	})();
/******/ 	
/************************************************************************/
/******/ 	
/******/ 	// startup
/******/ 	// Load entry module and return exports
/******/ 	// This entry module can't be inlined because the eval devtool is used.
/******/ 	var __webpack_exports__ = __webpack_require__("./script.js");
/******/ 	
/******/ })()
;