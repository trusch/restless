class Model {

  constructor(){
    this.__uid = '';
    this.__rev = 0;
    this.__data = {};
    this.__schema = '';
  }


  initFromData(data){
    this.__data = data;
    this.__uid = this._generateUUID();
    return this;
  }

  initFromUID(uid){
    if(typeof this.onPreInit === 'function'){
      this.onPreInit();
    }
    if(typeof this._doInit === 'function'){
      this._doInit(uid); //impl by subclass
    }
    if(typeof this.onPostInit === 'function'){
      this.onPostInit();
    }
    return this;
  }

  remove(){
    var res = false;
    if(typeof this.onPreRemove === 'function'){
      this.onPreRemove();
    }
    if(typeof this._doRemove === 'function'){
      res = this._doRemove(); //impl by subclass
    }
    if(typeof this.onPostRemove === 'function'){
      this.onPostRemove();
    }
    return res;
  }

  commit(){
    var res = false;
    if(typeof this.onPreCommit === 'function'){
      this.onPreCommit();
    }
    if(typeof this._doCommit === 'function'){
      res = this._doCommit(); //impl by subclass
    }
    if(typeof this.onPostCommit === 'function'){
      this.onPostCommit();
    }
    return res;
  }

  get(key=''){
    let parts = key.split('.');
    let res = this.__data;
    for(let i=0;i<parts.length;i++){
      if(typeof res === 'object'){
        res = res[parts[i]]
      }else{
        res = undefined;
        break;
      }
    }
    return res;
  }

  set(key='',value){
    let parts = key.split('.');
    let current = this.__data;
    for(var i=0;i<parts.length-1;i++){
      var next = current[parts[i]];
      if(typeof next !== 'object'){
        current[parts[i]] = {};
        next = current[parts[i]];
      }
      current = next;
    }
    current[parts[parts.length-1]] = value;
  }

  _generateUUID() {
    var d = new Date().getTime();
    var uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = (d + Math.random()*16)%16 | 0;
        d = Math.floor(d/16);
        return (c=='x' ? r : (r&0x3|0x8)).toString(16);
    });
    return uuid;
  }
}
