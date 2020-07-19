'use strict';

import Websocket from './webSocket';
import DocumentHandler from './documentHandler';
const documentHandler = new DocumentHandler();
new Websocket(documentHandler);
