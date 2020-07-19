'use strict';

export default class DocumentHandler {
 
  createEvent(id, type, cb) {
    document.getElementById(id).addEventListener(type, cb);
  }

  valueByID(id) {
    const value = document.getElementById(id).value;
    return value;
  }

  removeElement(id) {
    document.getElementById(id).remove();
  }


  /*createElement(id, type, status, value) {
    const config = [{
      'class': status,
    }];
    DocumentHandler.createEl(id, type, value, config);
  }*/

  createElement(id, type, value, config) {
    let element = document.createElement(type);
    config.map(function (obj) {
      element.setAttribute(obj.name, obj.value)
    });
    let node = document.createTextNode(value);
    element.appendChild(node);
    document.getElementById(id).appendChild(element);
  }
}


