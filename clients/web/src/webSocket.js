'use strict';

import { v4 as uuidv4 } from 'uuid';

export default class Websocket {

  static _documentHandler = null;
  static _wsConnection = null;
  static _uuid = null;
  static _nickname = null;

  constructor(documentHandler) {
    
    Websocket._documentHandler = documentHandler;
    Websocket._uuid = uuidv4();
    Websocket._generalHTMLAttributes = {
      error: [{
        name: 'class',
        value: 'error'
      }],
      ok: [{
        name: 'class',
        value: 'ok'
      }],
      chatHistory: function (id, host) {
        return [{
          name: 'class',
          value: `user-message ${id} ${host}`,
        }, {
          name: 'id',
          value: id
        }]
      }
    }

    Websocket._documentHandler.createEvent('init-chat', 'click', function () {

      Websocket._nickname = Websocket._documentHandler.valueByID('nickname');

      if (Websocket._nickname == "") {
        Websocket._documentHandler.createElement(
          'container-information',
          'span',
          'You must put a nickname to start!',
          Websocket._generalHTMLAttributes.error,
        );
      } else {
        Websocket.start(function (status, error) {

          if (status) {
            Websocket._documentHandler.removeElement('container-information');
            const data = JSON.stringify({
              nickname: Websocket._nickname,
              uuid: Websocket._uuid,
              message: 'ping'
            });
            document.getElementById('container-chat').style = 'display: flex';
            Websocket._wsConnection.send(data);
            Websocket.initSender();
            Websocket.initReceiver();

          } else {
            const config = [{
              'class': 'error',
            }];
            Websocket._documentHandler.createElement(
              'container-information',
              'span',
              'Something is wrong:' + error,
              Websocket._generalHTMLAttributes.error,
            );
          }

          Websocket._wsConnection.onerror = function (e) {
            console.log(e)
          }
        });
      }
    });
  }

  static start(cb) {
    const wsConn = new WebSocket('ws://localhost:8000?uuid=' + Websocket._uuid );
    Websocket._wsConnection = wsConn
    wsConn.onopen = function () {
      Websocket._documentHandler.createElement(
        'status-wsConn',
        'div',
        'Now you are connected!',
        Websocket._generalHTMLAttributes.ok,
      );
      cb(true, null)
    };
    wsConn.onerror = function (e) {
      Websocket._documentHandler.createElement(
        'status-wsConn',
        'span',
        'Connection refused!',
        Websocket._generalHTMLAttributes.error,
      );
      cb(false, e)
    };
  }

  static initSender() {
    Websocket._documentHandler.createEvent('send-data', 'click', function () {
      const dataToSend = Websocket._documentHandler.valueByID('message');
      const data = JSON.stringify({
        nickname: Websocket._nickname,
        uuid: Websocket._uuid,
        message: dataToSend
      })
      Websocket._wsConnection.send(data);
    });
  }

  static initReceiver() {
    Websocket._wsConnection.onmessage = function (e) {
      const data = JSON.parse(e.data)
      let message = data.message
      if (message === 'ping') {
        message = `[${data.nickname}]: now is connected!`
      } else {
        message = `[${data.nickname}]: ${message} `
      }
      Websocket._documentHandler.createElement(
        'chat-history',
        'span',
        message,
        Websocket._generalHTMLAttributes.chatHistory(data.uuid, Websocket._uuid),
      );
    };
  }
}